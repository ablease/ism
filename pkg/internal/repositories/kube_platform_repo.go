package repositories

import (
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubePlatformRepo struct {
	client client.Client
}

func NewKubePlatformRepo(client client.Client) *KubePlatformRepo {
	return &KubePlatformRepo{
		client: client,
	}
}

func (repo *KubePlatformRepo) Get(resource types.NamespacedName) (*v1alpha1.Platform, error) {
	broker := &v1alpha1.Platform{}

	err := repo.client.Get(ctx, resource, broker)
	if err != nil {
		return nil, err
	}

	return broker, nil
}

func (repo *KubePlatformRepo) UpdateState(broker *v1alpha1.Platform, newState v1alpha1.PlatformState) error {
	broker.Status.State = newState

	return repo.client.Status().Update(ctx, broker)
}
