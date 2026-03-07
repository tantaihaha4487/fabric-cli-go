package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fabric-cli/fabric-cli-go/internal/generator"
	"github.com/fabric-cli/fabric-cli-go/internal/versions"
)

var (
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	magenta = "\033[35m"
	white   = "\033[37m"
	gray    = "\033[90m"
	bold    = "\033[1m"
	reset   = "\033[0m"
	check   = "✓"
	arrow   = "❯"
)

func printHeader(title string) {
	fmt.Println()
	fmt.Println(bold + cyan + "╭───────────────────────────────────────────────" + reset)
	fmt.Println(bold + cyan + "│" + reset + " " + bold + white + title + reset)
	fmt.Println(bold + cyan + "╰───────────────────────────────────────────────" + reset)
}

func printSubHeader(title string) {
	fmt.Println()
	fmt.Println(bold + yellow + "━━ " + reset + bold + white + title + reset)
}

func promptInput(prompt, defaultVal, hint string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(reset + "  " + arrow + " " + bold + white + prompt + reset)
	if defaultVal != "" {
		fmt.Print(gray + " [" + defaultVal + "]" + reset)
	}
	if hint != "" {
		fmt.Print(gray + " " + hint + reset)
	}
	fmt.Print(": ")

	if !scanner.Scan() {
		return ""
	}
	result := strings.TrimSpace(scanner.Text())
	if result == "" {
		return defaultVal
	}
	return result
}

func promptSelect(prompt string, options []string, defaultIdx int) int {
	fmt.Println()
	fmt.Println("  " + bold + white + prompt + reset)
	fmt.Println("  " + gray + "  Use number to select or press Enter for default" + reset)
	fmt.Println()

	for i, opt := range options {
		if i == defaultIdx {
			fmt.Printf("  %s %s %s %s\n", green+"▸"+reset, green+"["+reset+green+"✓"+reset+green+"]"+reset, bold+white+opt+reset, gray+"(default)"+reset)
		} else {
			fmt.Printf("    %s [ ] %s\n", gray, gray+opt+reset)
		}
	}

	fmt.Print("  " + arrow + " Select: ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return defaultIdx
		}
		var n int
		fmt.Sscanf(input, "%d", &n)
		if n > 0 && n <= len(options) {
			return n - 1
		}
	}
	return defaultIdx
}

func promptYesNo(prompt string, defaultYes bool) bool {
	fmt.Println()
	fmt.Print("  " + arrow + " " + bold + white + prompt + reset)
	if defaultYes {
		fmt.Print(gray + " [Y/n]" + reset)
	} else {
		fmt.Print(gray + " [y/N]" + reset)
	}
	fmt.Print(": ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			return defaultYes
		}
		return strings.ToLower(input) == "y" || strings.ToLower(input) == "yes"
	}
	return defaultYes
}

func GetProjectConfigInteractive() (*generator.ProjectConfig, error) {
	cfg := &generator.ProjectConfig{}

	printHeader("Fabric Mod Project Generator")
	fmt.Println()
	fmt.Println("  " + gray + "Create a new Fabric mod project with ease" + reset)

	// Mod ID
	printSubHeader("Project Identity")
	cfg.ModID = promptInput("Mod ID", "mymod", "(lowercase, no spaces)")
	if cfg.ModID == "" {
		cfg.ModID = "mymod"
	}

	// Auto-suggest mod name from mod ID
	suggestedName := strings.Title(cfg.ModID)
	cfg.ModName = promptInput("Mod Name", suggestedName, "(display name)")
	if cfg.ModName == "" {
		cfg.ModName = suggestedName
	}

	cfg.ModDescription = promptInput("Description", "My awesome Fabric mod", "")

	// Package name (auto-suggest from mod ID)
	suggestedPkg := "com.example." + cfg.ModID
	cfg.PackageName = promptInput("Package Name", suggestedPkg, "")
	if cfg.PackageName == "" {
		cfg.PackageName = suggestedPkg
	}

	printSubHeader("Versions")

	// Minecraft Version (default to latest = 0)
	mcIdx := promptSelect("Minecraft Version", versions.McVersions, 0)
	cfg.MCVersion = versions.McVersions[mcIdx]

	// Loom Version
	loomIdx := promptSelect("Loom Version", versions.LoomVersions, 0)
	cfg.LoomVersion = versions.LoomVersions[loomIdx]

	// Loader Version
	loaderIdx := promptSelect("Loader Version", versions.LoaderVersions, 0)
	cfg.LoaderVersion = versions.LoaderVersions[loaderIdx]

	// Java Version (default 17)
	javaIdx := 1
	for i, v := range versions.JavaVersions {
		if v == "17" {
			javaIdx = i
			break
		}
	}
	javaIdx = promptSelect("JDK Version", versions.JavaVersions, javaIdx)
	cfg.JavaVersion = versions.JavaVersions[javaIdx]

	// Yarn Mappings (auto-set based on MC version)
	cfg.YarnMappings = versions.YarnMappings[cfg.MCVersion]
	fmt.Println()
	fmt.Println("  " + cyan + "▸" + reset + " Yarn Mappings: " + green + cfg.YarnMappings + reset + " " + gray + "(auto-set)")

	// Fabric API Version (auto-set based on MC version)
	cfg.APIVersion = versions.FabricAPIVersions[cfg.MCVersion]
	fmt.Println("  " + cyan + "▸" + reset + " Fabric API: " + green + cfg.APIVersion + reset + " " + gray + "(auto-set)")

	printSubHeader("Environment & Features")

	// Environment
	envOptions := []string{"* (Both)", "client", "server"}
	envIdx := promptSelect("Environment", envOptions, 0)
	envMap := map[int]string{0: "*", 1: "client", 2: "server"}
	cfg.Environment = envMap[envIdx]

	// Language
	langOptions := []string{"Java", "Kotlin"}
	langIdx := promptSelect("Language", langOptions, 0)
	cfg.Language = langOptions[langIdx]

	// Use Mixins
	cfg.UseMixins = promptYesNo("Use Mixins?", true)

	// Use Datagen
	cfg.UseDatagen = promptYesNo("Use Data Generation?", false)

	// Use Official Mappings
	cfg.UseOfficialMappings = promptYesNo("Use Mojang Official Mappings?", false)

	// Include Fabric API (only if not using official mappings)
	if !cfg.UseOfficialMappings {
		cfg.IncludeFabricAPI = promptYesNo("Include Fabric API?", true)
	} else {
		cfg.IncludeFabricAPI = false
	}

	printSubHeader("Metadata")

	cfg.Author = promptInput("Author", "Anonymous", "")
	if cfg.Author == "" {
		cfg.Author = "Anonymous"
	}

	licenseOptions := []string{"MIT", "Apache-2.0", "LGPL-3.0", "GPL-3.0", "All Rights Reserved"}
	licenseIdx := promptSelect("License", licenseOptions, 0)
	cfg.License = licenseOptions[licenseIdx]

	// Summary
	printHeader("Summary")
	fmt.Println()
	fmt.Printf("  %s %s %s\n", green+check+reset, "Mod ID:", cyan+cfg.ModID+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Mod Name:", cyan+cfg.ModName+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Package:", cyan+cfg.PackageName+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Language:", cyan+cfg.Language+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Minecraft:", cyan+cfg.MCVersion+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Loom:", cyan+cfg.LoomVersion+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Loader:", cyan+cfg.LoaderVersion+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "JDK:", cyan+"Java "+cfg.JavaVersion+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Environment:", cyan+cfg.Environment+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Mixins:", cyan+boolToStr(cfg.UseMixins)+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Datagen:", cyan+boolToStr(cfg.UseDatagen)+reset)
	mappings := "Yarn"
	if cfg.UseOfficialMappings {
		mappings = "Official"
	}
	fmt.Printf("  %s %s %s\n", green+check+reset, "Mappings:", cyan+mappings+reset)
	fmt.Printf("  %s %s %s\n", green+check+reset, "Fabric API:", cyan+boolToStr(cfg.IncludeFabricAPI)+reset)
	fmt.Println()

	confirm := promptYesNo("Create project?", true)
	if !confirm {
		return nil, fmt.Errorf("cancelled by user")
	}

	return cfg, nil
}

func boolToStr(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func GetProjectConfig(defaultName string) (*generator.ProjectConfig, error) {
	return GetProjectConfigInteractive()
}

func GetCodeConfig() (*generator.CodeConfig, error) {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &generator.CodeConfig{}

	fmt.Println("\n=== Code Generation ===\n")

	fmt.Print("Package Name: ")
	if !scanner.Scan() {
		return nil, scanner.Err()
	}
	cfg.PackageName = strings.TrimSpace(scanner.Text())

	return cfg, nil
}
