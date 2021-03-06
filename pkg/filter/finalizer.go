package filter

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type Filter func(object runtime.Object) runtime.Object

type Filters []Filter

func (f Filters) See(obj runtime.Object) error {
	for _, f := range f {
		if obj == nil {
			continue
		}
		obj = f(obj)
	}
	return nil
}

func HasFinalizer(obj runtime.Object) runtime.Object {
	m, err := meta.Accessor(obj)
	if err != nil {
		return nil
	}

	if len(m.GetFinalizers()) == 0 {
		return nil
	}

	return obj
}

func IsDeleted(obj runtime.Object) runtime.Object {
	m, err := meta.Accessor(obj)
	if err != nil {
		return nil
	}

	if m.GetDeletionTimestamp() == nil {
		return nil
	}

	return obj
}

func IsDeletedOutsideWindow(window time.Duration) func(obj runtime.Object) runtime.Object {
	return func(obj runtime.Object) runtime.Object {
		m, err := meta.Accessor(obj)
		if err != nil {
			return nil
		}

		deletion := m.GetDeletionTimestamp()

		if deletion == nil {
			return nil
		}

		if time.Since(deletion.Time) < window {
			return nil
		}

		return obj
	}
}
