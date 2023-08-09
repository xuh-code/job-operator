package v1alpha1

import (
	"gitee.com/xuh-code/job/apis/job/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Manage to the Hub version (v1alpha1).
func (src *Manage) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.Manage)

	image := src.Spec.Image
	name := src.Spec.Name
	replicas := src.Spec.Replicas
	//scheduleParts := []string{"*", "*", "*", "*", "*"}
	//if sched.Minute != nil {
	//	scheduleParts[0] = string(*sched.Minute)
	//}
	//if sched.Hour != nil {
	//	scheduleParts[1] = string(*sched.Hour)
	//}
	//if sched.DayOfMonth != nil {
	//	scheduleParts[2] = string(*sched.DayOfMonth)
	//}
	//if sched.Month != nil {
	//	scheduleParts[3] = string(*sched.Month)
	//}
	//if sched.DayOfWeek != nil {
	//	scheduleParts[4] = string(*sched.DayOfWeek)
	//}
	//dst.Spec.Schedule = strings.Join(scheduleParts, " ")
	dst.Spec.Image = image
	dst.Spec.Name = name
	dst.Spec.Replicas = replicas

	// +kubebuilder:docs-gen:collapse=rote conversion

	return nil
}

// ConvertFrom 从 Hub 版本 (v1) 转换到这个版本。
func (dst *Manage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.Manage)

	dst.Spec.Image = src.Spec.Image
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Spec.Name = src.Spec.Name

	// +kubebuilder:docs-gen:collapse=rote conversion

	return nil
}
