package wizard

import "fmt"

func parseMCVersion(mcVersion string) (major, minor, patch int, count int) {
	count, _ = fmt.Sscanf(mcVersion, "%d.%d.%d", &major, &minor, &patch)
	return major, minor, patch, count
}

func UsesImplicitMappingsProfile(mcVersion string) bool {
	major, _, _, count := parseMCVersion(mcVersion)
	return count >= 1 && major >= 26
}

func ShouldDefaultToOfficialMappings(mcVersion string) bool {
	if UsesImplicitMappingsProfile(mcVersion) {
		return false
	}

	major, minor, _, count := parseMCVersion(mcVersion)
	return count >= 2 && major == 1 && minor >= 21
}
