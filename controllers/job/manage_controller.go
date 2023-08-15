/*
Copyright 2023 xuh.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package job

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	jobv1alpha2 "gitee.com/xuh-code/job/apis/job/v1alpha2"
)

// ManageReconciler reconciles a Manage object
type ManageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=job.wajc.com,resources=manages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=job.wajc.com,resources=manages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=job.wajc.com,resources=manages/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Manage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ManageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name, "Request.Name", req.String())

	// TODO(user): your logic here
	logger.Info("Reconcile Manage v2 Info")
	//logger.Info(req.String())

	job := &jobv1alpha2.Manage{}

	// 获取Manage 实例
	err := r.Client.Get(context.TODO(), req.NamespacedName, job)
	if err != nil {
		if errors.IsNotFound(err) {
			//	找不到请求对象, 可以在协调请求后删除.
			//	拥有的对象会自动进行垃圾回收. 对于其他清理逻辑, 请使用终结器
			//	返回而且不排队
			return ctrl.Result{}, nil
		} else {
			status := jobv1alpha2.ManageStatus{
				Message: err.Error(),
			}
			job.Status = status
			_ = r.Status().Update(ctx, job)

			// 读取对象时报错- 重新排队请求.
			return ctrl.Result{}, err
		}

	}

	// 根据CRD生成一个新的Deployment资源对象
	deploy := newPodForCR(job)
	if err := controllerutil.SetControllerReference(job, deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// 检查现有deployment资源是否存在
	found := &v1.Deployment{}
	if err = r.Client.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found); err != nil {
		//	如果不存在 deployment, name, if中创建
		if errors.IsNotFound(err) {
			logger.Info("Create new deployment", "Deployment.Name", deploy.Name)
			err := r.Client.Create(context.TODO(), deploy)
			if err != nil {
				status := jobv1alpha2.ManageStatus{
					Message: err.Error(),
				}
				job.Status = status

				_ = r.Status().Update(ctx, job)
				return ctrl.Result{}, nil
			}
		} else {
			status := jobv1alpha2.ManageStatus{
				Message: err.Error(),
			}
			job.Status = status
			_ = r.Status().Update(ctx, job)

			logger.Info("Create new deployment", "Deployment.Name", deploy.Name, "error info : ", err.Error())
			return ctrl.Result{}, err
		}
	} else {
		err := r.Client.Update(context.TODO(), deploy)
		if err != nil {

			status := jobv1alpha2.ManageStatus{
				Message: err.Error(),
			}
			job.Status = status
			_ = r.Status().Update(ctx, job)

			return ctrl.Result{}, err
		} else {
			status := jobv1alpha2.ManageStatus{
				Message: "success",
			}
			job.Status = status
			_ = r.Status().Update(ctx, job)
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{}, nil
}

func newPodForCR(cr *jobv1alpha2.Manage) *v1.Deployment {

	labels := map[string]string{
		"app": cr.Name,
	}

	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: "default",
			Labels:    labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas: cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v12.PodSpec{
					Containers: []v12.Container{
						{
							Name:  "user",
							Image: cr.Spec.Image,
							Ports: []v12.ContainerPort{
								{
									//Name:          "tcp80",
									ContainerPort: 80,
									Protocol:      v12.ProtocolTCP,
								},
							},
						},
					},
				},
			},
			//Strategy: v1.DeploymentStrategy{
			//	Type: v1.RollingUpdateDeploymentStrategyType,
			//	RollingUpdate: &v1.RollingUpdateDeployment{
			//		MaxUnavailable: &intstr.IntOrString{
			//			Type:   2,
			//			IntVal: 0,
			//			StrVal: "0",
			//		},
			//		MaxSurge: &intstr.IntOrString{
			//			Type:   3,
			//			StrVal: "25%",
			//		},
			//	},
			//},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jobv1alpha2.Manage{}).
		Complete(r)
}
