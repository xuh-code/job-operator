package job

import "sigs.k8s.io/controller-runtime/pkg/predicate"

type ResourceLabelChangedPredicate struct {
	predicate.Funcs
}

//
//func (rl *ResourceLabelChangedPredicate) Update(e event.UpdateEvent) bool {
//	if e.ObjectNew.GetUID() != e.ObjectOld.GetUID() {
//		return true
//	}
//	return false
//}
