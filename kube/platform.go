package kube

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pivotal-cf/ism/osbapi"
	"github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
)

type PlatformRegisterTimeoutErr struct {
	platformName string
}

func (e PlatformRegisterTimeoutErr) Error() string {
	return fmt.Sprintf("timed out waiting for platform '%s' to be registered", e.platformName)
}

type Platform struct {
	KubeClient          client.Client
	RegistrationTimeout time.Duration
}

func (b *Platform) FindAll() ([]*osbapi.Platform, error) {
	list := &v1alpha1.PlatformList{}
	if err := b.KubeClient.List(context.TODO(), &client.ListOptions{}, list); err != nil {
		return []*osbapi.Platform{}, err
	}

	platforms := []*osbapi.Platform{}
	for _, platform := range list.Items {
		platforms = append(platforms, &osbapi.Platform{
			ID:   string(platform.UID),
			Name: platform.Spec.Name,
			URL:  platform.Spec.URL,
		})
	}

	return platforms, nil
}

func (b *Platform) Register(platform *osbapi.Platform) error {
	platformResource := &v1alpha1.Platform{
		ObjectMeta: metav1.ObjectMeta{
			Name:      platform.Name,
			Namespace: "default",
		},
		Spec: v1alpha1.PlatformSpec{
			Name: platform.Name,
			URL:  platform.URL,
		},
	}

	if err := b.KubeClient.Create(context.TODO(), platformResource); err != nil {
		return err
	}

	return nil
}

// func (b *Platform) waitForPlatformRegistration(platform *v1alpha1.Platform) error {
// 	err := wait.Poll(time.Second/2, b.RegistrationTimeout, func() (bool, error) {
// 		fetchedPlatform := &v1alpha1.Platform{}
//
// 		err := b.KubeClient.Get(context.TODO(), types.NamespacedName{Name: platform.Name, Namespace: platform.Namespace}, fetchedPlatform)
// 		if err == nil && fetchedPlatform.Status.State == v1alpha1.PlatformStateRegistered {
// 			return true, nil
// 		}
//
// 		return false, nil
// 	})
//
// 	if err != nil {
// 		if err == wait.ErrWaitTimeout {
// 			return PlatformRegisterTimeoutErr{platformName: platform.Name}
// 		}
//
// 		return err
// 	}
//
// 	return nil
// }
