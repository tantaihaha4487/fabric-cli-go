package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/config"
	"github.com/tantaihaha4487/fabric-cli-go/generator"
	"github.com/tantaihaha4487/fabric-cli-go/internal/javaport"
	"github.com/tantaihaha4487/fabric-cli-go/wizard"
)

const version = "1.0.0"

func printUsage() {
	fmt.Println("Fabric Mod Project Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  fabric-cli                                    # Interactive wizard mode")
	fmt.Println("  fabric-cli [options] <mc_version> <mod_name> <mod_version> [group_id]")
	fmt.Println()
	fmt.Println("Quick Mode Arguments:")
	fmt.Println("  mc_version    Minecraft version (e.g., 1.21.4, 1.20.1)")
	fmt.Println("  mod_name      Display name for your mod (auto-converted to mod ID)")
	fmt.Println("  mod_version   Version of your mod (e.g., 1.0.0)")
	fmt.Println("  group_id      Maven group ID (optional, uses saved config or 'com.example')")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("  -v, --version        Show version information")
	fmt.Println("  --no-mixins          Disable Mixins support")
	fmt.Println("  --client-only        Set environment to client only")
	fmt.Println("  --server-only        Set environment to server only")
	fmt.Println("  --java-version=N     Set Java version (default: recommended for MC version)")
	fmt.Println("  --license=TYPE       Set license (default: MIT)")
	fmt.Println("  --official-mappings  Use official Mojang mappings (default: Yarn, ignored for 26.x)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Interactive wizard")
	fmt.Println("  fabric-cli")
	fmt.Println()
	fmt.Println("  # Quick mode with defaults")
	fmt.Println(`  fabric-cli 1.21.4 "My Cool Mod" 1.0.0`)
	fmt.Println()
	fmt.Println("  # Quick mode with all options")
	fmt.Println(`  fabric-cli 1.21.4 "My Cool Mod" 1.0.0 com.example \`)
	fmt.Println(`    --no-mixins --client-only --java-version=17 --license=Apache-2.0`)
}

func printVersion() {
	fmt.Printf("fabric-cli version %s\n", version)
}

type CLIConfig struct {
	ShowHelp         bool
	ShowVersion      bool
	NoMixins         bool
	ClientOnly       bool
	ServerOnly       bool
	OfficialMappings bool
	JavaVersion      int
	JavaVersionSet   bool
	License          string
	Args             []string
}

func parseArgs() *CLIConfig {
	cfg := &CLIConfig{
		JavaVersion: 21,
		License:     "MIT",
	}

	i := 1
	for i < len(os.Args) {
		arg := os.Args[i]

		switch arg {
		case "-h", "--help":
			cfg.ShowHelp = true
			i++
		case "-v", "--version":
			cfg.ShowVersion = true
			i++
		case "--no-mixins":
			cfg.NoMixins = true
			i++
		case "--client-only":
			cfg.ClientOnly = true
			i++
		case "--server-only":
			cfg.ServerOnly = true
			i++
		case "--official-mappings":
			cfg.OfficialMappings = true
			i++
		default:
			if strings.HasPrefix(arg, "--java-version=") {
				val := strings.TrimPrefix(arg, "--java-version=")
				if v, err := parseInt(val); err == nil {
					cfg.JavaVersion = v
					cfg.JavaVersionSet = true
				}
				i++
			} else if strings.HasPrefix(arg, "--license=") {
				cfg.License = strings.TrimPrefix(arg, "--license=")
				i++
			} else if !strings.HasPrefix(arg, "-") {
				// This is a positional argument
				cfg.Args = append(cfg.Args, arg)
				i++
			} else {
				// Unknown flag
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n\n", arg)
				printUsage()
				os.Exit(1)
			}
		}
	}

	return cfg
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

func main() {
	cfg := parseArgs()

	// Handle --help
	if cfg.ShowHelp {
		printUsage()
		os.Exit(0)
	}

	// Handle --version
	if cfg.ShowVersion {
		printVersion()
		os.Exit(0)
	}

	// Get positional arguments
	args := cfg.Args

	// Determine mode: wizard or quick
	if len(args) == 0 {
		// Interactive wizard mode
		runWizardMode()
	} else if len(args) >= 3 && len(args) <= 4 {
		// Quick mode
		runQuickMode(args, cfg)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid number of arguments\n\n")
		printUsage()
		os.Exit(1)
	}
}

func runWizardMode() {
	fmt.Println("[Fabric] Fabric Project Generator")
	fmt.Println("Fetching latest versions...")

	// Fetch versions from APIs
	client := api.NewClient()
	fabricVersions, modrinthVersions, err := client.FetchAllVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching versions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[OK] Found %d Minecraft versions\n", len(fabricVersions.Game))
	fmt.Printf("[OK] Found %d mapping versions\n", len(fabricVersions.Mappings))
	fmt.Printf("[OK] Found %d Fabric API versions\n", len(modrinthVersions))
	fmt.Println()

	// Create wizard
	wiz := wizard.NewWizard()
	wiz.AddStep(wizard.NewVersionStep(fabricVersions, modrinthVersions))
	wiz.AddStep(&wizard.MetadataStep{})
	wiz.AddStep(&wizard.OptionsStep{})

	// Execute wizard
	ctx, err := wiz.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in wizard: %v\n", err)
		os.Exit(1)
	}

	// Generate project
	generateProject(ctx)
}

func runQuickMode(args []string, cliCfg *CLIConfig) {
	mcVersion := args[0]
	modName := args[1]
	modVersion := args[2]

	// Optional group_id
	var groupID string
	if len(args) >= 4 {
		groupID = args[3]
	}

	fmt.Println("[Fabric] Fabric Project Generator (Quick Mode)")
	fmt.Println("Fetching latest versions...")

	// Fetch versions from APIs
	client := api.NewClient()
	fabricVersions, modrinthVersions, err := client.FetchAllVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching versions: %v\n", err)
		os.Exit(1)
	}

	// Auto-complete mod ID from mod name
	modID := wizard.NormalizeAutoModID(modName)

	useImplicitMappings := wizard.UsesImplicitMappingsProfile(mcVersion)

	// Determine mappings based on version profile and --official-mappings flag
	var yarnMappings string
	if !useImplicitMappings && !cliCfg.OfficialMappings {
		// Auto-select latest yarn mappings for MC version
		yarnMappings = getLatestYarnMapping(fabricVersions.Mappings, mcVersion)
		if yarnMappings == "" {
			fmt.Fprintf(os.Stderr, "Error: No yarn mappings found for Minecraft %s\n", mcVersion)
			os.Exit(1)
		}
	}

	// Auto-select latest loader version
	loaderVersion := getLatestLoaderVersion(fabricVersions.Loader)
	if loaderVersion == "" {
		fmt.Fprintf(os.Stderr, "Error: No loader versions found\n")
		os.Exit(1)
	}

	// Auto-select latest fabric API for MC version
	apiVersion := getLatestAPIVersion(modrinthVersions, mcVersion)
	if apiVersion == "" {
		fmt.Fprintf(os.Stderr, "Warning: No Fabric API found for Minecraft %s\n", mcVersion)
	}

	// Load saved config for defaults
	savedCfg, _ := config.Load()
	if savedCfg == nil {
		savedCfg = config.DefaultConfig()
	}

	// Use command-line group_id if provided, otherwise use config or default
	if groupID == "" {
		groupID = savedCfg.GroupID
	}

	// Determine environment
	environment := "*"
	if cliCfg.ClientOnly {
		environment = "client"
	} else if cliCfg.ServerOnly {
		environment = "server"
	}

	javaVersion := cliCfg.JavaVersion
	if !cliCfg.JavaVersionSet {
		javaVersion = javaport.GetRecommendedJava(mcVersion)
	}

	// Create context
	ctx := &wizard.ProjectContext{
		MCVersion:           mcVersion,
		YarnMappings:        yarnMappings,
		LoaderVersion:       loaderVersion,
		APIVersion:          apiVersion,
		ModID:               modID,
		ModName:             modName,
		ModDescription:      fmt.Sprintf("A Fabric mod: %s", modName),
		License:             cliCfg.License,
		GroupID:             groupID,
		Version:             modVersion,
		UseMixins:           !cliCfg.NoMixins,
		UseOfficialMappings: !useImplicitMappings && cliCfg.OfficialMappings,
		Environment:         environment,
		JavaVersion:         javaVersion,
		Templates:           make(map[string]string),
	}

	// Save config with preferences
	savedCfg.GroupID = groupID
	savedCfg.Version = modVersion
	savedCfg.Save()

	// Show summary
	fmt.Println("\n[Summary] Configuration Summary:")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("Minecraft:    %s\n", ctx.MCVersion)
	if useImplicitMappings {
		fmt.Printf("Mappings:     Template default (26.x)\n")
	} else if ctx.UseOfficialMappings {
		fmt.Printf("Mappings:     Official Mojang\n")
	} else {
		fmt.Printf("Mappings:     Yarn %s\n", ctx.YarnMappings)
	}
	fmt.Printf("Loader:       %s\n", ctx.LoaderVersion)
	fmt.Printf("Fabric API:   %s\n", ctx.APIVersion)
	fmt.Println()
	fmt.Printf("Mod ID:       %s\n", ctx.ModID)
	fmt.Printf("Mod Name:     %s\n", ctx.ModName)
	fmt.Printf("Group ID:     %s\n", ctx.GroupID)
	fmt.Printf("Version:      %s\n", ctx.Version)
	fmt.Printf("Mixins:       %v\n", ctx.UseMixins)
	fmt.Printf("Environment:  %s\n", ctx.Environment)
	fmt.Printf("Java:         %d\n", ctx.JavaVersion)
	fmt.Printf("License:      %s\n", ctx.License)
	fmt.Println()

	// Generate project
	generateProject(ctx)
}

func getLatestYarnMapping(mappings []api.Mapping, mcVersion string) string {
	var latest string
	var latestBuild int
	for _, m := range mappings {
		if m.GameVersion == mcVersion && m.Build > latestBuild {
			latest = m.Version
			latestBuild = m.Build
		}
	}
	return latest
}

func getLatestLoaderVersion(loaders []api.LoaderVersion) string {
	for _, l := range loaders {
		if l.Stable {
			return l.Version
		}
	}
	// If no stable, return first
	if len(loaders) > 0 {
		return loaders[0].Version
	}
	return ""
}

func getLatestAPIVersion(versions []api.ModrinthVersion, mcVersion string) string {
	for _, v := range versions {
		for _, gv := range v.GameVersions {
			if gv == mcVersion {
				return v.VersionNumber
			}
		}
	}
	return ""
}

func generateProject(ctx *wizard.ProjectContext) {
	projectPath := ctx.ModID
	fmt.Printf("[Generate] Generating project in '%s/'...\n\n", projectPath)

	gen := generator.NewGenerator(ctx)
	if err := gen.Generate(projectPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("[Success] Project generated successfully!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectPath)
	fmt.Println("  ./gradlew build")
}
