package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fabric-cli/fabric-cli-go/cmd"
	"github.com/fabric-cli/fabric-cli-go/internal/versions"
)

const (
	FabricMetaURL  = "https://meta.fabricmc.net/v2/versions"
	ModrinthAPIURL = "https://api.modrinth.com/v2/project/fabric-api/version"
)

type FabricVersions struct {
	Game     []GameVersion     `json:"game"`
	Mappings []MappingsVersion `json:"mappings"`
	Loader   []LoaderVersion   `json:"loader"`
}

type GameVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

type MappingsVersion struct {
	GameVersion string `json:"gameVersion"`
	Version     string `json:"version"`
	Build       int    `json:"build"`
}

type LoaderVersion struct {
	Version string `json:"version"`
}

type ModrinthVersion struct {
	GameVersions []string `json:"game_versions"`
	Version      string   `json:"version_number"`
}

func fetchJSON(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func isMajorVersion(v string) bool {
	// Major versions like 1.21, 1.20, 1.19, etc.
	// Not snapshots like 21w13a, 24w14a, experimental snapshots, etc.
	if len(v) < 4 {
		return false
	}
	// Must start with "1."
	if v[0] != '1' || v[1] != '.' {
		return false
	}
	// Must have second digit after dot (like 1.21, 1.20)
	return v[2] >= '1' && v[2] <= '9'
}

func loadVersions() {
	// Try to fetch from Fabric Meta API
	var fabricVersions FabricVersions
	if err := fetchJSON(FabricMetaURL, &fabricVersions); err == nil {
		// Filter only stable major versions (like 1.21, 1.20, etc.)
		seen := make(map[string]bool)
		for _, v := range fabricVersions.Game {
			// Only stable releases
			if !v.Stable {
				continue
			}
			// Only major versions (like 1.21, 1.20, not snapshots like 21w13a)
			if !isMajorVersion(v.Version) {
				continue
			}
			if !seen[v.Version] {
				versions.McVersions = append(versions.McVersions, v.Version)
				seen[v.Version] = true
			}
		}

		// Get latest loader versions (just take top stable ones)
		seenLoader := make(map[string]bool)
		for _, v := range fabricVersions.Loader {
			if !seenLoader[v.Version] && len(versions.LoaderVersions) < 15 {
				versions.LoaderVersions = append(versions.LoaderVersions, v.Version)
				seenLoader[v.Version] = true
			}
		}

		for _, v := range fabricVersions.Mappings {
			key := v.GameVersion
			// Version from API already includes build number (e.g., "1.21.11+build.4")
			if _, ok := versions.YarnMappings[key]; !ok {
				versions.YarnMappings[key] = v.Version
			}
		}
	}

	// Try to fetch Fabric API versions from Modrinth
	var modrinthVersions []ModrinthVersion
	if err := fetchJSON(ModrinthAPIURL, &modrinthVersions); err == nil {
		for _, v := range modrinthVersions {
			if len(v.GameVersions) > 0 {
				key := v.GameVersions[0]
				if _, ok := versions.FabricAPIVersions[key]; !ok {
					versions.FabricAPIVersions[key] = v.Version
				}
			}
		}
	}

	// If no versions loaded, use defaults
	if len(versions.McVersions) == 0 {
		versions.McVersions = []string{"1.21.4", "1.21.3", "1.21.2", "1.21.1", "1.21", "1.20.6", "1.20.5", "1.20.4", "1.20.3", "1.20.2", "1.20.1", "1.20"}
	}
	if len(versions.LoaderVersions) == 0 {
		versions.LoaderVersions = []string{"0.18.1", "0.18.0", "0.17.2", "0.17.1", "0.17.0", "0.16.14", "0.16.13", "0.16.12", "0.16.11", "0.16.10", "0.15.11", "0.15.10"}
	}
	if len(versions.YarnMappings) == 0 {
		versions.YarnMappings = map[string]string{
			"1.21.4": "1.21.4+build.3", "1.21.3": "1.21.3+build.8", "1.21.2": "1.21.2+build.4",
			"1.21.1": "1.21.1+build.10", "1.21": "1.21+build.1-v2", "1.20.6": "1.20.6+build.3",
			"1.20.5": "1.20.5+build.8", "1.20.4": "1.20.4+build.8", "1.20.3": "1.20.3+build.7",
			"1.20.2": "1.20.2+build.8", "1.20.1": "1.20.1+build.18", "1.20": "1.20+build.3",
		}
	}
	if len(versions.FabricAPIVersions) == 0 {
		versions.FabricAPIVersions = map[string]string{
			"1.21.4": "0.100.8+1.21.4", "1.21.3": "0.100.4+1.21.3", "1.21.2": "0.100.0+1.21.2",
			"1.21.1": "0.99.0+1.21.1", "1.21": "0.98.0+1.21", "1.20.6": "0.92.2+1.20.6",
			"1.20.5": "0.91.0+1.20.5", "1.20.4": "0.91.0+1.20.4", "1.20.3": "0.90.0+1.20.3",
			"1.20.2": "0.89.0+1.20.2", "1.20.1": "0.88.0+1.20.1", "1.20": "0.87.0+1.20",
		}
	}

	// Also update cmd package versions
	cmd.McVersions = versions.McVersions
	cmd.LoaderVersions = versions.LoaderVersions
}

func main() {
	loadVersions()
	cmd.Execute()
}
