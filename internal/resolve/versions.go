package resolve

import (
	"fmt"

	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/config"
	"github.com/tantaihaha4487/fabric-cli-go/internal/profile"
	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
)

type QuickOptions struct {
	MCVersion           string
	ModName             string
	ModVersion          string
	GroupID             string
	License             string
	NoMixins            bool
	Environment         string
	JavaVersion         int
	JavaVersionSet      bool
	UseOfficialMappings bool
}

type VersionData struct {
	Fabric   *api.FabricVersions
	Modrinth []api.ModrinthVersion
}

func LatestYarnMapping(mappings []api.Mapping, mcVersion string) string {
	var latest string
	var latestBuild int
	for _, m := range mappings {
		if m.GameVersion == mcVersion && m.Build > latestBuild {
			latest = m.Version
			latestBuild = m.Build
		}
	}
	return latest
}

func LatestLoaderVersion(loaders []api.LoaderVersion) string {
	for _, l := range loaders {
		if l.Stable {
			return l.Version
		}
	}
	if len(loaders) > 0 {
		return loaders[0].Version
	}
	return ""
}

func LatestAPIVersion(versions []api.ModrinthVersion, mcVersion string) string {
	for _, v := range versions {
		for _, gv := range v.GameVersions {
			if gv == mcVersion {
				return v.VersionNumber
			}
		}
	}
	return ""
}

func QuickContext(opts QuickOptions, savedCfg *config.Config, versions VersionData) (*project.Context, profile.BuildProfile, error) {
	buildProfile := profile.ForMinecraftVersion(opts.MCVersion)

	loaderVersion := LatestLoaderVersion(versions.Fabric.Loader)
	if loaderVersion == "" {
		return nil, buildProfile, fmt.Errorf("no loader versions found")
	}

	apiVersion := LatestAPIVersion(versions.Modrinth, opts.MCVersion)

	groupID := opts.GroupID
	if groupID == "" {
		groupID = savedCfg.GroupID
	}

	javaVersion := opts.JavaVersion
	if !opts.JavaVersionSet {
		javaVersion = buildProfile.RecommendedJavaVersion
	}

	ctx := project.NewContext()
	ctx.MCVersion = opts.MCVersion
	ctx.LoaderVersion = loaderVersion
	ctx.APIVersion = apiVersion
	ctx.ModID = project.NormalizeAutoModID(opts.ModName)
	ctx.ModName = opts.ModName
	ctx.ModDescription = fmt.Sprintf("A Fabric mod: %s", opts.ModName)
	ctx.License = opts.License
	ctx.GroupID = groupID
	ctx.Version = opts.ModVersion
	ctx.UseMixins = !opts.NoMixins
	ctx.Environment = opts.Environment
	ctx.JavaVersion = javaVersion

	if buildProfile.SupportsExplicitMapping {
		ctx.UseOfficialMappings = opts.UseOfficialMappings
		if !ctx.UseOfficialMappings {
			ctx.YarnMappings = LatestYarnMapping(versions.Fabric.Mappings, opts.MCVersion)
			if ctx.YarnMappings == "" {
				return nil, buildProfile, fmt.Errorf("no yarn mappings found for Minecraft %s", opts.MCVersion)
			}
		}
	}

	return ctx, buildProfile, nil
}
