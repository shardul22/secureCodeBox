// SPDX-FileCopyrightText: the secureCodeBox authors
//
// SPDX-License-Identifier: Apache-2.0

package scancontrollers

import (
	"context"
	"fmt"
	"strings"

	executionv1 "github.com/secureCodeBox/secureCodeBox/operator/apis/execution/v1"
	util "github.com/secureCodeBox/secureCodeBox/operator/utils"
	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ScanReconciler) startParser(scan *executionv1.Scan) error {
	ctx := context.Background()
	namespacedName := fmt.Sprintf("%s/%s", scan.Namespace, scan.Name)
	log := r.Log.WithValues("scan_parse", namespacedName)

	jobs, err := r.getJobsForScan(scan, client.MatchingLabels{"securecodebox.io/job-type": "parser"})
	if err != nil {
		return err
	}
	if len(jobs.Items) > 0 {
		log.V(8).Info("Job already exists. Doesn't need to be created.")
		return nil
	}

	parseType := scan.Status.RawResultType

	// get the parse definition matching the parseType of the scan result
	var parseDefinitionSpec executionv1.ParseDefinitionSpec
	if scan.Spec.ResourceMode == nil || *scan.Spec.ResourceMode == executionv1.NamespaceLocal {
		var parseDefinition executionv1.ParseDefinition
		if err := r.Get(ctx, types.NamespacedName{Name: parseType, Namespace: scan.Namespace}, &parseDefinition); err != nil {
			log.V(7).Info("Unable to fetch ParseDefinition")

			scan.Status.State = executionv1.ScanStateErrored
			scan.Status.ErrorDescription = fmt.Sprintf("No ParseDefinition for ResultType '%s' found in Scans Namespace.", parseType)
			if err := r.updateScanStatus(ctx, scan); err != nil {
				r.Log.Error(err, "unable to update Scan status")
				return err
			}

			return fmt.Errorf("no ParseDefinition of type '%s' found", parseType)
		}
		log.Info("Matching ParseDefinition Found", "ParseDefinition", parseType)
		parseDefinitionSpec = parseDefinition.Spec
	} else if *scan.Spec.ResourceMode == executionv1.ClusterWide {
		var clusterParseDefinition executionv1.ClusterParseDefinition
		if err := r.Get(ctx, types.NamespacedName{Name: parseType}, &clusterParseDefinition); err != nil {
			log.V(7).Info("Unable to fetch ClusterParseDefinition")

			scan.Status.State = executionv1.ScanStateErrored
			scan.Status.ErrorDescription = fmt.Sprintf("No ClusterParseDefinition for ResultType '%s' found.", parseType)
			if err := r.updateScanStatus(ctx, scan); err != nil {
				r.Log.Error(err, "unable to update Scan status")
				return err
			}

			return fmt.Errorf("no ClusterParseDefinition of type '%s' found", parseType)
		}
		log.Info("Matching ClusterParseDefinition Found", "ClusterParseDefinition", parseType)
		parseDefinitionSpec = clusterParseDefinition.Spec
	}

	urlExpirationDuration, err := util.GetUrlExpirationDuration(util.ParserController)
	if err != nil {
		r.Log.Error(err, "Failed to parse parser url expiration")
		panic(err)
	}

	findingsUploadURL, err := r.PresignedPutURL(*scan, "findings.json", urlExpirationDuration)
	if err != nil {
		r.Log.Error(err, "Could not get presigned url from s3 or compatible storage provider")
		return err
	}
	rawResultDownloadURL, err := r.PresignedGetURL(*scan, scan.Status.RawResultFile, urlExpirationDuration)
	if err != nil {
		return err
	}

	rules := []rbacv1.PolicyRule{
		{
			APIGroups: []string{"execution.securecodebox.io"},
			Resources: []string{"scans"},
			Verbs:     []string{"get"},
		},
		{
			APIGroups: []string{"execution.securecodebox.io"},
			Resources: []string{"scans/status"},
			Verbs:     []string{"get", "patch"},
		},
		{
			APIGroups: []string{"execution.securecodebox.io"},
			Resources: []string{"parsedefinitions"},
			Verbs:     []string{"get"},
		},
	}
	r.ensureServiceAccountExists(
		scan.Namespace,
		"parser",
		"Parser need to access the status of Scans to update how many findings have been identified",
		rules,
	)

	labels := scan.ObjectMeta.DeepCopy().Labels
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["securecodebox.io/job-type"] = "parser"
	automountServiceAccountToken := true
	var backOffLimit int32 = 3
	truePointer := true
	falsePointer := false

	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("200m"),
			corev1.ResourceMemory: resource.MustParse("100Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("400m"),
			corev1.ResourceMemory: resource.MustParse("200Mi"),
		},
	}
	if len(parseDefinitionSpec.Resources.Requests) != 0 || len(parseDefinitionSpec.Resources.Limits) != 0 {
		resources = parseDefinitionSpec.Resources
	}

	job := &batch.Job{
		ObjectMeta: metav1.ObjectMeta{
			Annotations:  make(map[string]string),
			GenerateName: util.TruncateName(fmt.Sprintf("parse-%s", scan.Name)),
			Namespace:    scan.Namespace,
			Labels:       labels,
		},
		Spec: batch.JobSpec{
			TTLSecondsAfterFinished: parseDefinitionSpec.TTLSecondsAfterFinished,
			BackoffLimit:            &backOffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "securecodebox",
					},
					Annotations: map[string]string{
						"auto-discovery.securecodebox.io/ignore": "true",
						"sidecar.istio.io/inject":                allowIstioSidecarInjectionInJobs,
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: "parser",
					ImagePullSecrets:   parseDefinitionSpec.ImagePullSecrets,
					Containers: []corev1.Container{
						{
							Name:  "parser",
							Image: parseDefinitionSpec.Image,
							Env: []corev1.EnvVar{
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name:  "SCAN_NAME",
									Value: scan.Name,
								},
							},
							Args: []string{
								rawResultDownloadURL,
								findingsUploadURL,
							},
							ImagePullPolicy: parseDefinitionSpec.ImagePullPolicy,
							Resources:       resources,
							SecurityContext: &corev1.SecurityContext{
								RunAsNonRoot:             &truePointer,
								AllowPrivilegeEscalation: &falsePointer,
								ReadOnlyRootFilesystem:   &truePointer,
								Privileged:               &falsePointer,
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
							},
						},
					},
					AutomountServiceAccountToken: &automountServiceAccountToken,
				},
			},
		},
	}
	job.Spec.Template.Labels = util.MergeStringMaps(job.Spec.Template.Labels, scan.ObjectMeta.DeepCopy().Labels)

	// Merge Env from ParserTemplate
	job.Spec.Template.Spec.Containers[0].Env = append(
		job.Spec.Template.Spec.Containers[0].Env,
		parseDefinitionSpec.Env...,
	)
	// Merge VolumeMounts from ParserTemplate
	job.Spec.Template.Spec.Containers[0].VolumeMounts = append(
		job.Spec.Template.Spec.Containers[0].VolumeMounts,
		parseDefinitionSpec.VolumeMounts...,
	)
	// Merge Volumes from ParserTemplate
	job.Spec.Template.Spec.Volumes = append(
		job.Spec.Template.Spec.Volumes,
		parseDefinitionSpec.Volumes...,
	)

	// Merge NodeSelectors from ParseDefinition & Scan into Parse Job
	job.Spec.Template.Spec.NodeSelector = util.MergeStringMaps(job.Spec.Template.Spec.NodeSelector, parseDefinitionSpec.NodeSelector, scan.Spec.NodeSelector)

	// Set affinity based on scan, if defined, or parseDefinition if not overridden by scan
	if scan.Spec.Affinity != nil {
		job.Spec.Template.Spec.Affinity = scan.Spec.Affinity
	} else {
		job.Spec.Template.Spec.Affinity = parseDefinitionSpec.Affinity
	}

	// Set tolerations, either from parseDefinition or from scan
	if scan.Spec.Tolerations != nil {
		job.Spec.Template.Spec.Tolerations = scan.Spec.Tolerations
	} else {
		job.Spec.Template.Spec.Tolerations = parseDefinitionSpec.Tolerations
	}

	r.Log.V(8).Info("Configuring customCACerts for Parser")
	injectCustomCACertsIfConfigured(job)

	if err := ctrl.SetControllerReference(scan, job, r.Scheme); err != nil {
		return err
	}

	log.V(7).Info("Constructed Job object", "job args", strings.Join(job.Spec.Template.Spec.Containers[0].Args, ", "))

	if err := r.Create(ctx, job); err != nil {
		log.Error(err, "unable to create Job for Parser", "job", job)
		return err
	}

	scan.Status.State = executionv1.ScanStateParsing
	if err := r.updateScanStatus(ctx, scan); err != nil {
		log.Error(err, "unable to update Scan status")
		return err
	}

	log.V(7).Info("created Parse Job for Scan", "job", job)
	return nil
}

func (r *ScanReconciler) checkIfParsingIsCompleted(scan *executionv1.Scan) error {
	ctx := context.Background()

	status, err := r.checkIfJobIsCompleted(scan, client.MatchingLabels{"securecodebox.io/job-type": "parser"})
	if err != nil {
		return err
	}

	switch status {
	case completed:
		r.Log.V(7).Info("Parsing is completed")
		scan.Status.State = executionv1.ScanStateParseCompleted
		if err := r.updateScanStatus(ctx, scan); err != nil {
			r.Log.Error(err, "unable to update Scan status")
			return err
		}
	case failed:
		scan.Status.State = executionv1.ScanStateErrored
		scan.Status.ErrorDescription = "Failed to run the Parser. This is likely a Bug, we would like to know about. Please open up a Issue on GitHub."
		if err := r.updateScanStatus(ctx, scan); err != nil {
			r.Log.Error(err, "unable to update Scan status")
			return err
		}
	}

	return nil
}
