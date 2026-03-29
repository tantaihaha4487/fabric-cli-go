package profile

import "testing"

func TestForMinecraftVersion(t *testing.T) {
	tests := []struct {
		mcVersion          string
		wantExplicit       bool
		wantDefaultMapping MappingMode
		wantGradle         string
		wantDependency     string
		wantAPIKey         string
	}{
		{"1.20.1", true, MappingModeYarn, "9.2.1", "modImplementation", "fabric"},
		{"1.21.4", true, MappingModeOfficial, "9.2.1", "modImplementation", "fabric"},
		{"26.1", false, MappingModeImplicit, "9.3.0", "implementation", "fabric-api"},
	}

	for _, tt := range tests {
		t.Run(tt.mcVersion, func(t *testing.T) {
			got := ForMinecraftVersion(tt.mcVersion)
			if got.SupportsExplicitMapping != tt.wantExplicit {
				t.Fatalf("SupportsExplicitMapping = %v, want %v", got.SupportsExplicitMapping, tt.wantExplicit)
			}
			if got.DefaultMappingMode != tt.wantDefaultMapping {
				t.Fatalf("DefaultMappingMode = %q, want %q", got.DefaultMappingMode, tt.wantDefaultMapping)
			}
			if got.GradleVersion != tt.wantGradle {
				t.Fatalf("GradleVersion = %q, want %q", got.GradleVersion, tt.wantGradle)
			}
			if got.DependencyConfiguration != tt.wantDependency {
				t.Fatalf("DependencyConfiguration = %q, want %q", got.DependencyConfiguration, tt.wantDependency)
			}
			if got.APIDependencyKey != tt.wantAPIKey {
				t.Fatalf("APIDependencyKey = %q, want %q", got.APIDependencyKey, tt.wantAPIKey)
			}
		})
	}
}
