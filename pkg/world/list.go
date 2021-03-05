package world

import (
	"context"
	"sync"

	"github.com/rancher/wrangler/pkg/ratelimit"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/pager"
)

type Traveler interface {
	See(obj runtime.Object) error
}

type TravelerFunc func(obj runtime.Object) error

func (t TravelerFunc) See(obj runtime.Object) error {
	return t(obj)
}

type Trip struct {
	namespace  string
	list       metav1.ListOptions
	restConfig *rest.Config
	k8s        kubernetes.Interface
	dynamic    dynamic.Interface
	sem        *semaphore.Weighted
	writeLock  sync.Mutex
}

type Options struct {
	Namespace   string
	List        metav1.ListOptions
	Parallelism int64
}

func NewTrip(restConfig *rest.Config, opts *Options) (*Trip, error) {
	if opts == nil {
		opts = &Options{}
	}

	if restConfig.RateLimiter == nil {
		restConfig.RateLimiter = ratelimit.None
	}

	k8s, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	dynamic, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	var sem *semaphore.Weighted
	if opts.Parallelism > 0 {
		sem = semaphore.NewWeighted(opts.Parallelism)
	} else {
		sem = semaphore.NewWeighted(10)
	}

	return &Trip{
		namespace:  opts.Namespace,
		list:       opts.List,
		restConfig: restConfig,
		k8s:        k8s,
		dynamic:    dynamic,
		sem:        sem,
	}, nil
}

func (t *Trip) Wander(ctx context.Context, traveler Traveler) error {
	_, apis, err := t.k8s.Discovery().ServerGroupsAndResources()
	if err != nil {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)
	for _, api := range apis {
		for _, api2 := range api.APIResources {
			listable := false
			for _, verb := range api2.Verbs {
				if verb == "list" {
					listable = true
					break
				}
			}
			if !listable {
				break
			}

			if t.namespace != "" && !api2.Namespaced {
				continue
			}

			gvr := schema.FromAPIVersionAndKind(api.GroupVersion, api2.Kind).
				GroupVersion().
				WithResource(api2.Name)

			var client dynamic.ResourceInterface
			if t.namespace == "" {
				client = t.dynamic.Resource(gvr)
			} else {
				client = t.dynamic.Resource(gvr).Namespace(t.namespace)
			}

			eg.Go(func() error {
				if err := t.listAll(ctx, traveler, client, api, api2); err != nil {
					logrus.Warn("Failed to list", api.GroupVersion, api2.Kind, err)
				}
				return nil
			})
		}
	}

	return eg.Wait()
}

func (t *Trip) listAll(ctx context.Context, traveler Traveler, client dynamic.ResourceInterface, api *metav1.APIResourceList, api2 metav1.APIResource) error {
	t.sem.Acquire(ctx, 1)
	defer t.sem.Release(1)

	pager := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		objs, err := client.List(ctx, t.list)
		return objs, err
	}))

	return pager.EachListItem(ctx, t.list, func(obj runtime.Object) error {
		t.writeLock.Lock()
		defer t.writeLock.Unlock()
		if err := traveler.See(obj); err != nil {
			m, err2 := meta.Accessor(obj)
			if err2 == nil {
				logrus.Warn("Failed to process", api.GroupVersion, api2.Kind, m.GetNamespace(), m.GetName(), err)
			}
		}
		return nil
	})
}
