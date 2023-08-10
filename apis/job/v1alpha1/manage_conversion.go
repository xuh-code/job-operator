package v1alpha1

import (
	"gitee.com/xuh-code/job/apis/job/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this Manage to the Hub version (v1alpha1).
func (src *Manage) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.Manage)

	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Image = src.Spec.Image
	dst.Spec.Name = src.Spec.Name
	dst.Spec.Replicas = src.Spec.Replicas

	// +kubebuilder:docs-gen:collapse=rote conversion

	return nil
}

// ConvertFrom 从 Hub 版本 (v1) 转换到这个版本。
func (dst *Manage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.Manage)

	dst.Spec.Image = src.Spec.Image
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Spec.Name = src.Spec.Name

	return nil
}
