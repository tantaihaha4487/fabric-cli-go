package wizard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/config"
)

// ProjectContext holds all configuration data collected from the wizard
type ProjectContext struct {
	// Version data
	MCVersion     string
	YarnMappings  string
	LoaderVersion string
	APIVersion    string

	// Mod metadata
	ModID          string
	ModName        string
	ModDescription string
	License        string
	GroupID        string
	Version        string

	// Options
	UseMixins           bool
	UseOfficialMappings bool
	Environment         string // * or client or server
	JavaVersion         int

	// Templates
	Templates map[string]string
}

// Step represents a single step in the wizard
type Step interface {
	Name() string
	Execute(ctx *ProjectContext) error
}

// Wizard manages the step-by-step configuration process
type Wizard struct {
	steps []Step
	ctx   *ProjectContext
}

// NewWizard creates a new wizard instance
func NewWizard() *Wizard {
	// Load saved config
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	return &Wizard{
		steps: make([]Step, 0),
		ctx: &ProjectContext{
			Templates:   make(map[string]string),
			JavaVersion: 21,
			Environment: "*",
			License:     "MIT",
			GroupID:     cfg.GroupID,
			Version:     cfg.Version,
		},
	}
}

// AddStep adds a step to the wizard
func (w *Wizard) AddStep(step Step) {
	w.steps = append(w.steps, step)
}

// Execute runs all steps in sequence
func (w *Wizard) Execute() (*ProjectContext, error) {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00D8BD")).
		MarginBottom(1)

	fmt.Println(titleStyle.Render("╔════════════════════════════════════════════════════════╗"))
	fmt.Println(titleStyle.Render("║     Fabric Mod Project Generator                       ║"))
	fmt.Println(titleStyle.Render("╚════════════════════════════════════════════════════════╝"))
	fmt.Println()

	for _, step := range w.steps {
		if err := step.Execute(w.ctx); err != nil {
			return nil, err
		}
	}

	// Save config with user's preferences
	if err := w.SaveConfig(); err != nil {
		fmt.Printf("[!]  Warning: Could not save config: %v\n", err)
	}

	return w.ctx, nil
}

// SaveConfig saves the current GroupID and Version to config file
func (w *Wizard) SaveConfig() error {
	cfg := &config.Config{
		GroupID: w.ctx.GroupID,
		Version: w.ctx.Version,
	}
	return cfg.Save()
}

// VersionStep handles Minecraft and dependency version selection
type VersionStep struct {
	FabricVersions   *api.FabricVersions
	ModrinthVersions []api.ModrinthVersion
}

// NewVersionStep creates a new version selection step
func NewVersionStep(fabric *api.FabricVersions, modrinth []api.ModrinthVersion) *VersionStep {
	return &VersionStep{
		FabricVersions:   fabric,
		ModrinthVersions: modrinth,
	}
}

func (s *VersionStep) Name() string {
	return "Version Selection"
}

func (s *VersionStep) Execute(ctx *ProjectContext) error {
	// Prepare Minecraft version options (stable versions only)
	var allStableVersions []api.GameVersion
	for _, g := range s.FabricVersions.Game {
		if g.Stable {
			allStableVersions = append(allStableVersions, g)
		}
	}

	// Prepare all version options
	var mcOptions []huh.Option[string]
	for _, g := range allStableVersions {
		mcOptions = append(mcOptions, huh.NewOption(g.Version, g.Version))
	}

	// Show total count
	fmt.Printf("[OK] Loaded %d Minecraft versions\n\n", len(mcOptions))

	// Ask for search term
	var searchTerm string
	searchForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Type to search versions (or leave empty to show all)").
				Placeholder("e.g., 1.21, 1.20.1").
				Value(&searchTerm),
		),
	)

	if err := searchForm.Run(); err != nil {
		return err
	}

	// Filter versions based on search term
	searchTerm = strings.ToLower(strings.TrimSpace(searchTerm))
	if searchTerm != "" {
		var filteredOptions []huh.Option[string]
		for _, opt := range mcOptions {
			if strings.Contains(strings.ToLower(opt.Key), searchTerm) {
				filteredOptions = append(filteredOptions, opt)
			}
		}
		if len(filteredOptions) > 0 {
			mcOptions = filteredOptions
			fmt.Printf("[OK] Found %d matching versions\n\n", len(mcOptions))
		} else {
			fmt.Printf("[!] No matches found, showing all versions\n\n")
		}
	}

	// Minecraft version selection
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Minecraft Version").
				Options(mcOptions...).
				Value(&ctx.MCVersion),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Filter yarn mappings for selected MC version
	mappings := api.GetMappingsForVersion(s.FabricVersions.Mappings, ctx.MCVersion)
	var yarnOptions []huh.Option[string]
	for _, m := range mappings {
		yarnOptions = append(yarnOptions, huh.NewOption(m.Version, m.Version))
	}

	// Yarn version selection
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Yarn Mappings").
				Options(yarnOptions[:min(5, len(yarnOptions))]...).
				Value(&ctx.YarnMappings),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Loader version options
	var loaderOptions []huh.Option[string]
	for _, l := range s.FabricVersions.Loader {
		loaderOptions = append(loaderOptions, huh.NewOption(l.Version, l.Version))
	}

	// Loader version selection
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Fabric Loader").
				Options(loaderOptions[:min(5, len(loaderOptions))]...).
				Value(&ctx.LoaderVersion),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Filter API versions for selected MC version
	apiVersions := api.GetAPIVersionsForMCVersion(s.ModrinthVersions, ctx.MCVersion)
	var apiOptions []huh.Option[string]
	for _, v := range apiVersions {
		apiOptions = append(apiOptions, huh.NewOption(v.VersionNumber, v.VersionNumber))
	}

	// API version selection
	form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Fabric API").
				Options(apiOptions[:min(5, len(apiOptions))]...).
				Value(&ctx.APIVersion),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	return nil
}

// MetadataStep handles mod metadata input
type MetadataStep struct{}

func (s *MetadataStep) Name() string {
	return "Mod Metadata"
}

func (s *MetadataStep) Execute(ctx *ProjectContext) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Mod ID").
				Placeholder("lowercase, no spaces").
				Description("Unique identifier for your mod").
				Value(&ctx.ModID).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("mod ID is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Mod Name").
				Placeholder("My Awesome Mod").
				Description("Display name for your mod").
				Value(&ctx.ModName),

			huh.NewText().
				Title("Description").
				Placeholder("What does your mod do?").
				Lines(3).
				Value(&ctx.ModDescription),

			huh.NewInput().
				Title("Maven Group ID").
				Placeholder("com.example").
				Description("Java package namespace").
				Value(&ctx.GroupID),

			huh.NewInput().
				Title("Mod Version").
				Placeholder("1.0.0").
				Value(&ctx.Version),

			huh.NewSelect[string]().
				Title("License").
				Options(
					huh.NewOption("MIT", "MIT"),
					huh.NewOption("Apache-2.0", "Apache-2.0"),
					huh.NewOption("GPL-3.0", "GPL-3.0"),
					huh.NewOption("LGPL-3.0", "LGPL-3.0"),
					huh.NewOption("BSD-3-Clause", "BSD-3-Clause"),
					huh.NewOption("CC0-1.0", "CC0-1.0"),
					huh.NewOption("Other", "Other"),
				).
				Value(&ctx.License),
		),
	)

	return form.Run()
}

// OptionsStep handles additional configuration options
type OptionsStep struct{}

func (s *OptionsStep) Name() string {
	return "Additional Options"
}

func (s *OptionsStep) Execute(ctx *ProjectContext) error {
	var envChoice string

	// Set mixins to enabled by default
	ctx.UseMixins = true

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Enable Mixins?").
				Description("Mixins allow you to modify Minecraft's code").
				Value(&ctx.UseMixins).
				Affirmative("Yes").
				Negative("No"),

			huh.NewSelect[string]().
				Title("Environment").
				Description("Where will your mod run?").
				Options(
					huh.NewOption("Both (Client + Server)", "*"),
					huh.NewOption("Client Only", "client"),
					huh.NewOption("Server Only", "server"),
				).
				Value(&envChoice),

			huh.NewSelect[int]().
				Title("Java Version").
				Description("Target Java version for compilation").
				Options(
					huh.NewOption("Java 21", 21),
					huh.NewOption("Java 17", 17),
					huh.NewOption("Java 11", 11),
					huh.NewOption("Java 8", 8),
				).
				Value(&ctx.JavaVersion),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	ctx.Environment = envChoice

	// Set defaults for empty values
	if ctx.ModID == "" {
		ctx.ModID = "mymod"
	}
	if ctx.ModName == "" {
		ctx.ModName = "My Mod"
	}
	if ctx.ModDescription == "" {
		ctx.ModDescription = "A Fabric mod"
	}
	if ctx.GroupID == "" {
		ctx.GroupID = "com.example"
	}
	if ctx.Version == "" {
		ctx.Version = "1.0.0"
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
