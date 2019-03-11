package commands

import (
	"github.com/pivotal-cf/ism/osbapi"
)

//go:generate counterfeiter . PlatformRegistrar

type PlatformRegistrar interface {
	Register(*osbapi.Platform) error
}

//go:generate counterfeiter . PlatformFetcher

type PlatformFetcher interface {
	GetPlatforms() ([]*osbapi.Platform, error)
}

type PlatformCommand struct {
	PlatformRegisterCommand PlatformRegisterCommand `command:"register" platform:"Register a platform into the marketplace"`
	PlatformListCommand     PlatformListCommand     `command:"list" long-description:"Lists the platforms in the marketplace"`
}

type PlatformListCommand struct {
	UI              UI
	PlatformFetcher PlatformFetcher
}

type PlatformRegisterCommand struct {
	Name string `long:"name" description:"Name of the platform" required:"true"`
	URL  string `long:"url" description:"URL of the platform" required:"true"`

	UI                UI
	PlatformRegistrar PlatformRegistrar
}

func (cmd *PlatformRegisterCommand) Execute([]string) error {
	newPlatform := &osbapi.Platform{
		Name: cmd.Name,
		URL:  cmd.URL,
	}

	if err := cmd.PlatformRegistrar.Register(newPlatform); err != nil {
		return err
	}

	cmd.UI.DisplayText("Platform '{{.PlatformName}}' registered.", map[string]interface{}{"PlatformName": cmd.Name})

	return nil
}

func (cmd *PlatformListCommand) Execute([]string) error {
	platforms, err := cmd.PlatformFetcher.GetPlatforms()
	if err != nil {
		return err
	}

	if len(platforms) == 0 {
		cmd.UI.DisplayText("No platforms found.")
		return nil
	}

	platformsTable := buildPlatformTableData(platforms)
	cmd.UI.DisplayTable(platformsTable)
	return nil
}

func buildPlatformTableData(platforms []*osbapi.Platform) [][]string {
	headers := []string{"NAME", "URL"}
	data := [][]string{headers}

	for _, p := range platforms {
		row := []string{p.Name, p.URL}
		data = append(data, row)
	}
	return data
}
