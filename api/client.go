package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	FabricMetaURL  = "https://meta.fabricmc.net/v2/versions"
	ModrinthAPIURL = "https://api.modrinth.com/v2/project/P7dR8mSH/version"
)

// GameVersion represents a Minecraft version
type GameVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// Mapping represents a Yarn mapping version
type Mapping struct {
	GameVersion string `json:"gameVersion"`
	Version     string `json:"version"`
	Build       int    `json:"build"`
}

// LoaderVersion represents a Fabric loader version
type LoaderVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// FabricVersions holds all version data from Fabric Meta API
type FabricVersions struct {
	Game     []GameVersion   `json:"game"`
	Mappings []Mapping       `json:"mappings"`
	Loader   []LoaderVersion `json:"loader"`
}

// ModrinthVersion represents a version from Modrinth API
type ModrinthVersion struct {
	VersionNumber string   `json:"version_number"`
	GameVersions  []string `json:"game_versions"`
	Files         []struct {
		Filename string `json:"filename"`
	} `json:"files"`
}

// Client handles API requests
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchFabricVersions fetches version data from Fabric Meta API
func (c *Client) FetchFabricVersions() (*FabricVersions, error) {
	resp, err := c.httpClient.Get(FabricMetaURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Fabric versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Fabric Meta API returned status %d", resp.StatusCode)
	}

	var versions FabricVersions
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode Fabric versions: %w", err)
	}

	return &versions, nil
}

// FetchModrinthVersions fetches Fabric API versions from Modrinth
func (c *Client) FetchModrinthVersions() ([]ModrinthVersion, error) {
	resp, err := c.httpClient.Get(ModrinthAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Modrinth versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Modrinth API returned status %d", resp.StatusCode)
	}

	var versions []ModrinthVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("failed to decode Modrinth versions: %w", err)
	}

	// Filter to only include fabric-api files
	var filtered []ModrinthVersion
	for _, v := range versions {
		for _, f := range v.Files {
			if len(f.Filename) > 11 && f.Filename[:11] == "fabric-api-" {
				filtered = append(filtered, v)
				break
			}
		}
	}

	return filtered, nil
}

// FetchAllVersions fetches both Fabric and Modrinth versions concurrently
func (c *Client) FetchAllVersions() (*FabricVersions, []ModrinthVersion, error) {
	type result struct {
		fabric   *FabricVersions
		modrinth []ModrinthVersion
		err      error
	}

	ch := make(chan result, 2)

	// Fetch Fabric versions concurrently
	go func() {
		fabric, err := c.FetchFabricVersions()
		ch <- result{fabric: fabric, err: err}
	}()

	// Fetch Modrinth versions concurrently
	go func() {
		modrinth, err := c.FetchModrinthVersions()
		ch <- result{modrinth: modrinth, err: err}
	}()

	var fabricVersions *FabricVersions
	var modrinthVersions []ModrinthVersion

	// Wait for both goroutines
	for i := 0; i < 2; i++ {
		r := <-ch
		if r.err != nil {
			return nil, nil, r.err
		}
		if r.fabric != nil {
			fabricVersions = r.fabric
		}
		if r.modrinth != nil {
			modrinthVersions = r.modrinth
		}
	}

	return fabricVersions, modrinthVersions, nil
}

// GetMappingsForVersion returns yarn mappings for a specific Minecraft version
func GetMappingsForVersion(mappings []Mapping, mcVersion string) []Mapping {
	var result []Mapping
	for _, m := range mappings {
		if m.GameVersion == mcVersion {
			result = append(result, m)
		}
	}
	// If no mappings found, return all (graceful fallback)
	if len(result) == 0 {
		return mappings
	}
	return result
}

// GetAPIVersionsForMCVersion returns Fabric API versions compatible with a Minecraft version
func GetAPIVersionsForMCVersion(versions []ModrinthVersion, mcVersion string) []ModrinthVersion {
	var result []ModrinthVersion
	for _, v := range versions {
		for _, gv := range v.GameVersions {
			if gv == mcVersion {
				result = append(result, v)
				break
			}
		}
	}
	// If no versions found, return all (graceful fallback)
	if len(result) == 0 {
		return versions
	}
	return result
}
