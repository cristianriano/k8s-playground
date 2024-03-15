/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/batch/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	mydomainv1 "k8s.example.com/v2/api/v1"
)

// PdfDocumentReconciler reconciles a PdfDocument object
type PdfDocumentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=pdfdocuments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=pdfdocuments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=pdfdocuments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PdfDocument object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *PdfDocumentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx, "pdfdocument", req.NamespacedName)

	// Start 2 init containers, the first one will persist the text and the second one will convert the text to a PDF
	// Main container will just sleep

	// 1. Get the PdfDocument resource
	var pdfDoc mydomainv1.PdfDocument
	if err := r.Get(ctx, req.NamespacedName, &pdfDoc); err != nil {
		log.Log.Error(err, "unable to fetch PdfDocument")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. Create job spec with the init containers
	jobSpec, err := r.createJob(pdfDoc)
	if err != nil {
		log.Log.Error(err, "failed to create Job spec")
		return ctrl.Result{}, err
	}

	// 3. Create the job
	if err := r.Create(ctx, &jobSpec); err != nil {
		log.Log.Error(err, "unable to create Job")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PdfDocumentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1.PdfDocument{}).
		Complete(r)
}

func (r *PdfDocumentReconciler) createJob(pdfDoc mydomainv1.PdfDocument) (v1.Job, error) {
	image := "knsit/pandoc"
	//base64text := base64.StdEncoding.EncodeToString([]byte(pdfDoc.Spec.Text))

	jobSpec := v1.Job{
		TypeMeta: metav1.TypeMeta{APIVersion: v1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pdfDoc.Name + "-job",
			Namespace: pdfDoc.Namespace,
		},
		Spec: v1.JobSpec{
			Template: v12.PodTemplateSpec{
				Spec: v12.PodSpec{
					RestartPolicy: v12.RestartPolicyOnFailure,
					InitContainers: []v12.Container{
						{
							Name:    "persist-md",
							Image:   "alpine",
							Command: []string{"/bin/sh"},
							Args:    []string{"-c", fmt.Sprintf("echo %s > /data/text.md", pdfDoc.Spec.Text)},
							VolumeMounts: []v12.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
						{
							Name:    "convert-to-pdf",
							Image:   image,
							Command: []string{"/bin/sh"},
							Args:    []string{"-c", "pandoc", "/data/text.md", "-o", fmt.Sprintf("/data/%s.pdf", pdfDoc.Spec.Title)},
							VolumeMounts: []v12.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Containers: []v12.Container{
						{
							Name:    "main",
							Image:   "alpine",
							Command: []string{"/bin/sh", "-c", "sleep 1800"},
							VolumeMounts: []v12.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []v12.Volume{
						{
							Name: "data-volume",
							VolumeSource: v12.VolumeSource{
								EmptyDir: &v12.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	return jobSpec, nil
}
