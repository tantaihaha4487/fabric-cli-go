package resolve

import (
	"testing"

	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/config"
)

func TestQuickContextFor26Profile(t *testing.T) {
	ctx, buildProfile, err := QuickContext(QuickOptions{
		MCVersion:      "26.1",
		ModName:        "t",
		ModVersion:     "1.0.0",
		License:        "MIT",
		Environment:    "*",
		JavaVersion:    21,
		JavaVersionSet: false,
	}, config.DefaultConfig(), VersionData{
		Fabric: &api.FabricVersions{
			Loader: []api.LoaderVersion{{Version: "0.18.5", Stable: true}},
		},
		Modrinth: []api.ModrinthVersion{{VersionNumber: "0.144.3+26.1", GameVersions: []string{"26.1"}}},
	})
	if err != nil {
		t.Fatalf("QuickContext returned error: %v", err)
	}
	if buildProfile.SupportsExplicitMapping {
		t.Fatalf("expected implicit mappings profile for 26.1")
	}
	if ctx.ModID != "t_mod" {
		t.Fatalf("ModID = %q, want %q", ctx.ModID, "t_mod")
	}
	if ctx.UseOfficialMappings {
		t.Fatalf("expected UseOfficialMappings to be false for implicit profile")
	}
	if ctx.YarnMappings != "" {
		t.Fatalf("YarnMappings = %q, want empty", ctx.YarnMappings)
	}
	if ctx.JavaVersion != 25 {
		t.Fatalf("JavaVersion = %d, want 25", ctx.JavaVersion)
	}
}

func TestQuickContextForYarnProfile(t *testing.T) {
	ctx, buildProfile, err := QuickContext(QuickOptions{
		MCVersion:      "1.20.1",
		ModName:        "My Mod",
		ModVersion:     "1.0.0",
		License:        "MIT",
		Environment:    "client",
		JavaVersion:    17,
		JavaVersionSet: true,
	}, config.DefaultConfig(), VersionData{
		Fabric: &api.FabricVersions{
			Mappings: []api.Mapping{{GameVersion: "1.20.1", Version: "1.20.1+build.3", Build: 3}},
			Loader:   []api.LoaderVersion{{Version: "0.15.11", Stable: true}},
		},
		Modrinth: []api.ModrinthVersion{{VersionNumber: "0.91.0+1.20.1", GameVersions: []string{"1.20.1"}}},
	})
	if err != nil {
		t.Fatalf("QuickContext returned error: %v", err)
	}
	if !buildProfile.SupportsExplicitMapping {
		t.Fatalf("expected explicit mappings profile for 1.20.1")
	}
	if ctx.YarnMappings != "1.20.1+build.3" {
		t.Fatalf("YarnMappings = %q, want 1.20.1+build.3", ctx.YarnMappings)
	}
	if ctx.JavaVersion != 17 {
		t.Fatalf("JavaVersion = %d, want 17", ctx.JavaVersion)
	}
}
