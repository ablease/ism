package actors

import "github.com/pivotal-cf/ism/osbapi"

//go:generate counterfeiter . PlatformRepository

type PlatformRepository interface {
	FindAll() ([]*osbapi.Platform, error)
	Register(*osbapi.Platform) error
}

type PlatformsActor struct {
	Repository PlatformRepository
}

func (a *PlatformsActor) GetPlatforms() ([]*osbapi.Platform, error) {
	return a.Repository.FindAll()
}

//TODO: Make names consistent
func (a *PlatformsActor) Register(broker *osbapi.Platform) error {
	return a.Repository.Register(broker)
}
