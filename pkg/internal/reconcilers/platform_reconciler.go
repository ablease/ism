package reconcilers

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	v1alpha1 "github.com/pivotal-cf/ism/pkg/apis/osbapi/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//go:generate counterfeiter . KubeBrokerRepo

type KubePlatformRepo interface {
	Get(resource types.NamespacedName) (*v1alpha1.Platform, error)
	UpdateState(platform *v1alpha1.Platform, newState v1alpha1.PlatformState) error
}

type PlatformReconciler struct {
	kubePlatformRepo KubePlatformRepo
}

func NewPlatformReconciler(kubePlatformRepo KubePlatformRepo) *PlatformReconciler {
	return &PlatformReconciler{
		kubePlatformRepo: kubePlatformRepo,
	}
}

func (r *PlatformReconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	platform, err := r.kubePlatformRepo.Get(request.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, err
	}

	if platform.Status.State == v1alpha1.PlatformStateRegistered {
		return reconcile.Result{}, nil
	}

	if err := r.kubePlatformRepo.UpdateState(platform, v1alpha1.PlatformStateRegistered); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
