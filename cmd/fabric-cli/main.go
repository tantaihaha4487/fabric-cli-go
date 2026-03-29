package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/config"
	"github.com/tantaihaha4487/fabric-cli-go/internal/app"
	"github.com/tantaihaha4487/fabric-cli-go/internal/profile"
	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
	"github.com/tantaihaha4487/fabric-cli-go/internal/resolve"
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
	svc := app.NewService(api.NewClient())

	fmt.Println("[Fabric] Fabric Project Generator")
	fmt.Println("Fetching latest versions...")

	versions, err := svc.FetchVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching versions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[OK] Found %d Minecraft versions\n", len(versions.Fabric.Game))
	fmt.Printf("[OK] Found %d mapping versions\n", len(versions.Fabric.Mappings))
	fmt.Printf("[OK] Found %d Fabric API versions\n", len(versions.Modrinth))
	fmt.Println()

	wiz := svc.NewWizard(versions)
	ctx, err := wiz.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in wizard: %v\n", err)
		os.Exit(1)
	}

	generateProject(svc, ctx)
}

func runQuickMode(args []string, cliCfg *CLIConfig) {
	svc := app.NewService(api.NewClient())

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

	versions, err := svc.FetchVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching versions: %v\n", err)
		os.Exit(1)
	}

	// Load saved config for defaults
	savedCfg, _ := config.Load()
	if savedCfg == nil {
		savedCfg = config.DefaultConfig()
	}

	// Determine environment
	environment := "*"
	if cliCfg.ClientOnly {
		environment = "client"
	} else if cliCfg.ServerOnly {
		environment = "server"
	}

	ctx, buildProfile, err := resolve.QuickContext(resolve.QuickOptions{
		MCVersion:           mcVersion,
		ModName:             modName,
		ModVersion:          modVersion,
		GroupID:             groupID,
		License:             cliCfg.License,
		NoMixins:            cliCfg.NoMixins,
		Environment:         environment,
		JavaVersion:         cliCfg.JavaVersion,
		JavaVersionSet:      cliCfg.JavaVersionSet,
		UseOfficialMappings: cliCfg.OfficialMappings,
	}, savedCfg, versions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if ctx.APIVersion == "" {
		fmt.Fprintf(os.Stderr, "Warning: No Fabric API found for Minecraft %s\n", mcVersion)
	}

	// Save config with preferences
	savedCfg.GroupID = ctx.GroupID
	savedCfg.Version = modVersion
	savedCfg.Save()

	printSummary(ctx, buildProfile)
	generateProject(svc, ctx)
}

func printSummary(ctx *project.Context, buildProfile profile.BuildProfile) {
	fmt.Println("\n[Summary] Configuration Summary:")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("Minecraft:    %s\n", ctx.MCVersion)
	if !buildProfile.SupportsExplicitMapping {
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
}

func generateProject(svc *app.Service, ctx *project.Context) {
	projectPath := ctx.ModID
	fmt.Printf("[Generate] Generating project in '%s/'...\n\n", projectPath)

	if err := svc.Generate(ctx); err != nil {
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
