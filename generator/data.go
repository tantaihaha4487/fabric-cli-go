package generator

import (
	"github.com/tantaihaha4487/fabric-cli-go/internal/profile"
	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
)

type templateData struct {
	MCVersion        string
	YarnMappings     string
	LoaderVersion    string
	APIVersion       string
	ModID            string
	ModName          string
	ModDescription   string
	ModEnvironment   string
	License          string
	Mixins           bool
	MixinPackageName string
	GroupID          string
	ArtifactID       string
	Version          string
	JavaVersion      int
	LoomVersion      string
	OfficialMappings bool
	ImplicitMappings bool
	DependencyConfig string
	GradleVersion    string
	APIDependencyKey string
}

func buildTemplateData(ctx *project.Context) templateData {
	buildProfile := profile.ForMinecraftVersion(ctx.MCVersion)

	return templateData{
		MCVersion:        ctx.MCVersion,
		YarnMappings:     ctx.YarnMappings,
		LoaderVersion:    ctx.LoaderVersion,
		APIVersion:       ctx.APIVersion,
		ModID:            ctx.ModID,
		ModName:          ctx.ModName,
		ModDescription:   ctx.ModDescription,
		ModEnvironment:   ctx.Environment,
		License:          ctx.License,
		Mixins:           ctx.UseMixins,
		MixinPackageName: ctx.GroupID + "." + ctx.ModID + ".mixin",
		GroupID:          ctx.GroupID,
		ArtifactID:       ctx.ModID,
		Version:          ctx.Version,
		JavaVersion:      ctx.JavaVersion,
		LoomVersion:      buildProfile.LoomVersion,
		OfficialMappings: ctx.UseOfficialMappings,
		ImplicitMappings: !buildProfile.SupportsExplicitMapping,
		DependencyConfig: buildProfile.DependencyConfiguration,
		GradleVersion:    buildProfile.GradleVersion,
		APIDependencyKey: buildProfile.APIDependencyKey,
	}
}
