package cmd

import (
	"fmt"

	"github.com/fabric-cli/fabric-cli-go/internal/generator"
	"github.com/fabric-cli/fabric-cli-go/internal/prompts"
	"github.com/spf13/cobra"
)

var (
	modID               string
	modName             string
	description         string
	packageName         string
	author              string
	license             string
	mcVersion           string
	javaVersion         string
	loaderVersion       string
	loomVersion         string
	language            string
	environment         string
	template            string
	includeFabricAPI    bool
	useOfficialMappings bool
	useMixins           bool
	useDatagen          bool
	interactive         bool
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Fabric mod project",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *generator.ProjectConfig
		var err error

		// Interactive mode if no args or --interactive flag
		if interactive || len(args) == 0 {
			cfg, err = prompts.GetProjectConfigInteractive()
			if err != nil {
				return err
			}
		} else {
			// Non-interactive mode with args
			cfg, err = getConfigFromFlags(args)
			if err != nil {
				return err
			}
		}

		if err := generator.GenerateProject(cfg); err != nil {
			return err
		}

		fmt.Printf("\n✓ Created Fabric mod project: %s\n", cfg.ModID)
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", cfg.ModID)
		fmt.Println("  ./gradlew build")
		return nil
	},
}

func getConfigFromFlags(args []string) (*generator.ProjectConfig, error) {
	cfg := &generator.ProjectConfig{}

	if len(args) > 0 {
		cfg.ModID = args[0]
	}

	if cfg.ModID == "" {
		return nil, fmt.Errorf("mod ID is required")
	}

	cfg.ModName = modName
	cfg.ModDescription = description
	cfg.PackageName = packageName
	cfg.Author = author
	cfg.License = license
	cfg.MCVersion = mcVersion
	cfg.JavaVersion = javaVersion
	cfg.LoaderVersion = loaderVersion
	cfg.LoomVersion = loomVersion
	cfg.Language = language
	cfg.Environment = environment
	cfg.Template = template
	cfg.IncludeFabricAPI = includeFabricAPI
	cfg.UseOfficialMappings = useOfficialMappings
	cfg.UseMixins = useMixins
	cfg.UseDatagen = useDatagen

	// Set defaults
	if cfg.ModName == "" {
		cfg.ModName = cfg.ModID
	}
	if cfg.PackageName == "" {
		cfg.PackageName = "com.example." + cfg.ModID
	}
	if cfg.License == "" {
		cfg.License = "MIT"
	}
	if cfg.MCVersion == "" {
		cfg.MCVersion = "1.20.4"
	}
	if cfg.JavaVersion == "" {
		cfg.JavaVersion = "17"
	}
	if cfg.LoaderVersion == "" {
		cfg.LoaderVersion = "0.15.11"
	}
	if cfg.LoomVersion == "" {
		cfg.LoomVersion = "1.4.0"
	}
	if cfg.Language == "" {
		cfg.Language = "Java"
	}
	if cfg.Environment == "" {
		cfg.Environment = "*"
	}

	return cfg, nil
}

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate code (block/item/enchantment/effect)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &generator.CodeConfig{
			PackageName: packageName,
		}

		if cfg.PackageName == "" {
			cfg.PackageName = "com.example.mymod"
		}

		outputPath, err := generator.GenerateCode(cfg, args[0], args[1])
		if err != nil {
			return err
		}

		fmt.Printf("✓ Generated %s: %s\n", args[0], args[1])
		fmt.Printf("  Location: %s\n", outputPath)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVarP(&modName, "name", "n", "", "Mod name")
	newCmd.Flags().StringVarP(&description, "description", "d", "", "Mod description")
	newCmd.Flags().StringVarP(&packageName, "package", "p", "", "Package name (com.example.modid)")
	newCmd.Flags().StringVarP(&author, "author", "a", "Anonymous", "Author name")
	newCmd.Flags().StringVarP(&license, "license", "l", "MIT", "License (MIT/Apache-2.0/LGPL-3.0)")
	newCmd.Flags().StringVarP(&mcVersion, "mc-version", "m", "", "Minecraft version")
	newCmd.Flags().StringVarP(&javaVersion, "java-version", "j", "", "Java version")
	newCmd.Flags().StringVarP(&loaderVersion, "loader-version", "v", "", "Fabric loader version")
	newCmd.Flags().StringVarP(&loomVersion, "loom-version", "", "", "Loom version")
	newCmd.Flags().StringVarP(&language, "language", "g", "Java", "Language (Java/Kotlin)")
	newCmd.Flags().StringVarP(&environment, "environment", "e", "*", "Environment (*, client, server)")
	newCmd.Flags().StringVarP(&template, "template", "t", "", "Template (basic/mixin/datagen)")
	newCmd.Flags().BoolVarP(&includeFabricAPI, "fabric-api", "f", true, "Include Fabric API")
	newCmd.Flags().BoolVarP(&useOfficialMappings, "official-mappings", "o", false, "Use Mojang official mappings")
	newCmd.Flags().BoolVarP(&useMixins, "mixins", "", false, "Use Mixins")
	newCmd.Flags().BoolVarP(&useDatagen, "datagen", "", false, "Use Data Generation")
	newCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")

	generateCmd.Flags().StringVarP(&packageName, "package", "p", "", "Package name")
}
