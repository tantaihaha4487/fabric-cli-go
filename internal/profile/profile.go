package profile

import (
	"fmt"

	"github.com/tantaihaha4487/fabric-cli-go/internal/javaport"
)

type MappingMode string

const (
	MappingModeImplicit MappingMode = "implicit"
	MappingModeYarn     MappingMode = "yarn"
	MappingModeOfficial MappingMode = "official"
)

type BuildProfile struct {
	MCVersion               string
	SupportsExplicitMapping bool
	DefaultMappingMode      MappingMode
	LoomVersion             string
	GradleVersion           string
	DependencyConfiguration string
	APIDependencyKey        string
	RecommendedJavaVersion  int
}

func parseMCVersion(mcVersion string) (major, minor, patch int, count int) {
	count, _ = fmt.Sscanf(mcVersion, "%d.%d.%d", &major, &minor, &patch)
	return major, minor, patch, count
}

func ForMinecraftVersion(mcVersion string) BuildProfile {
	profile := BuildProfile{
		MCVersion:               mcVersion,
		SupportsExplicitMapping: true,
		DefaultMappingMode:      MappingModeYarn,
		LoomVersion:             "1.15-SNAPSHOT",
		GradleVersion:           "9.2.1",
		DependencyConfiguration: "modImplementation",
		APIDependencyKey:        "fabric",
		RecommendedJavaVersion:  javaport.GetRecommendedJava(mcVersion),
	}

	major, minor, _, count := parseMCVersion(mcVersion)
	if count >= 1 && major >= 26 {
		profile.SupportsExplicitMapping = false
		profile.DefaultMappingMode = MappingModeImplicit
		profile.GradleVersion = "9.3.0"
		profile.DependencyConfiguration = "implementation"
		profile.APIDependencyKey = "fabric-api"
		return profile
	}

	if count >= 2 && major == 1 && minor >= 21 {
		profile.DefaultMappingMode = MappingModeOfficial
	}

	return profile
}
