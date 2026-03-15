package main

import (
	"fmt"
	"os"

	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/generator"
	"github.com/tantaihaha4487/fabric-cli-go/wizard"
)

func main() {
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
	fmt.Printf("[OK] Found %d Yarn mappings\n", len(fabricVersions.Mappings))
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

	// Show summary
	fmt.Println("\n[Summary] Configuration Summary:")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("Minecraft:    %s\n", ctx.MCVersion)
	fmt.Printf("Yarn:         %s\n", ctx.YarnMappings)
	fmt.Printf("Loader:       %s\n", ctx.LoaderVersion)
	fmt.Printf("Fabric API:   %s\n", ctx.APIVersion)
	fmt.Println()
	fmt.Printf("Mod ID:       %s\n", ctx.ModID)
	fmt.Printf("Mod Name:     %s\n", ctx.ModName)
	fmt.Printf("Group ID:     %s\n", ctx.GroupID)
	fmt.Printf("Version:      %s\n", ctx.Version)
	fmt.Printf("Mixins:       %v\n", ctx.UseMixins)
	fmt.Println()

	// Generate project
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
