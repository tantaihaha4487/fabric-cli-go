package app

import (
	"github.com/tantaihaha4487/fabric-cli-go/api"
	"github.com/tantaihaha4487/fabric-cli-go/generator"
	"github.com/tantaihaha4487/fabric-cli-go/internal/project"
	"github.com/tantaihaha4487/fabric-cli-go/internal/resolve"
	"github.com/tantaihaha4487/fabric-cli-go/wizard"
)

type Service struct {
	client *api.Client
}

func NewService(client *api.Client) *Service {
	return &Service{client: client}
}

func (s *Service) FetchVersions() (resolve.VersionData, error) {
	fabricVersions, modrinthVersions, err := s.client.FetchAllVersions()
	if err != nil {
		return resolve.VersionData{}, err
	}

	return resolve.VersionData{
		Fabric:   fabricVersions,
		Modrinth: modrinthVersions,
	}, nil
}

func (s *Service) NewWizard(versions resolve.VersionData) *wizard.Wizard {
	wiz := wizard.NewWizard()
	wiz.AddStep(wizard.NewVersionStep(versions.Fabric, versions.Modrinth))
	wiz.AddStep(&wizard.MetadataStep{})
	wiz.AddStep(&wizard.OptionsStep{})
	return wiz
}

func (s *Service) Generate(ctx *project.Context) error {
	return generator.NewGenerator(ctx).Generate(ctx.ModID)
}
