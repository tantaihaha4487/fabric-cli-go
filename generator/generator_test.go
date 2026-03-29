package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
)

func TestRenderExplicitYarnProfile(t *testing.T) {
	ctx := &project.Context{
		MCVersion:     "1.20.1",
		YarnMappings:  "1.20.1+build.3",
		LoaderVersion: "0.15.11",
		APIVersion:    "0.91.0+1.20.1",
		ModID:         "testmod",
		ModName:       "Test Mod",
		GroupID:       "com.example",
		Version:       "1.0.0",
		Environment:   "*",
		JavaVersion:   17,
		UseMixins:     true,
	}

	buildGradle, fabricMod, mixins := renderTemplateFiles(t, ctx)
	if !strings.Contains(buildGradle, `mappings "net.fabricmc:yarn:${project.yarn_mappings}:v2"`) {
		t.Fatalf("expected yarn mappings in build.gradle")
	}
	if !strings.Contains(buildGradle, `modImplementation "net.fabricmc:fabric-loader:${project.loader_version}"`) {
		t.Fatalf("expected modImplementation in build.gradle")
	}
	if !strings.Contains(fabricMod, `"fabric": "*"`) {
		t.Fatalf("expected fabric dependency in fabric.mod.json")
	}
	if strings.Contains(mixins, `"overwrites"`) {
		t.Fatalf("did not expect implicit-profile overwrite block")
	}
}

func TestRenderOfficialProfile(t *testing.T) {
	ctx := &project.Context{
		MCVersion:           "1.21.4",
		LoaderVersion:       "0.16.9",
		APIVersion:          "0.110.0",
		ModID:               "testmod",
		ModName:             "Test Mod",
		GroupID:             "com.example",
		Version:             "1.0.0",
		Environment:         "*",
		JavaVersion:         21,
		UseMixins:           true,
		UseOfficialMappings: true,
	}

	buildGradle, _, _ := renderTemplateFiles(t, ctx)
	if !strings.Contains(buildGradle, `mappings loom.officialMojangMappings()`) {
		t.Fatalf("expected official mappings in build.gradle")
	}
}

func TestRenderImplicit26Profile(t *testing.T) {
	ctx := &project.Context{
		MCVersion:     "26.1",
		LoaderVersion: "0.18.5",
		APIVersion:    "0.144.3+26.1",
		ModID:         "testmod",
		ModName:       "Test Mod",
		GroupID:       "com.example",
		Version:       "1.0.0",
		Environment:   "*",
		JavaVersion:   25,
		UseMixins:     true,
	}

	buildGradle, fabricMod, mixins := renderTemplateFiles(t, ctx)
	if strings.Contains(buildGradle, `mappings `) {
		t.Fatalf("did not expect explicit mappings line in 26.x build.gradle")
	}
	if !strings.Contains(buildGradle, `implementation "net.fabricmc:fabric-loader:${project.loader_version}"`) {
		t.Fatalf("expected implementation dependency in 26.x build.gradle")
	}
	if !strings.Contains(fabricMod, `"fabric-api": "*"`) {
		t.Fatalf("expected fabric-api dependency in 26.x fabric.mod.json")
	}
	if !strings.Contains(mixins, `"overwrites"`) {
		t.Fatalf("expected implicit-profile overwrite block in mixins file")
	}
}

func renderTemplateFiles(t *testing.T, ctx *project.Context) (string, string, string) {
	t.Helper()

	gen := NewGenerator(ctx)
	data := buildTemplateData(ctx)
	projectPath := t.TempDir()

	if err := gen.generateFile(projectPath, "build.gradle", buildGradleTemplate, data); err != nil {
		t.Fatalf("generateFile(build.gradle) failed: %v", err)
	}
	if err := gen.generateFile(projectPath, "src/main/resources/fabric.mod.json", fabricModJsonTemplate, data); err != nil {
		t.Fatalf("generateFile(fabric.mod.json) failed: %v", err)
	}
	if err := gen.generateFile(projectPath, "src/main/resources/"+ctx.ModID+".mixins.json", mixinsJsonTemplate, data); err != nil {
		t.Fatalf("generateFile(mixins) failed: %v", err)
	}

	buildGradle, err := os.ReadFile(filepath.Join(projectPath, "build.gradle"))
	if err != nil {
		t.Fatalf("ReadFile(build.gradle) failed: %v", err)
	}
	fabricMod, err := os.ReadFile(filepath.Join(projectPath, "src/main/resources/fabric.mod.json"))
	if err != nil {
		t.Fatalf("ReadFile(fabric.mod.json) failed: %v", err)
	}
	mixins, err := os.ReadFile(filepath.Join(projectPath, "src/main/resources", ctx.ModID+".mixins.json"))
	if err != nil {
		t.Fatalf("ReadFile(mixins) failed: %v", err)
	}

	return string(buildGradle), string(fabricMod), string(mixins)
}
