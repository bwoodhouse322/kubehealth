package kubehealth

import (
	"fmt"
	"os"
	"text/tabwriter"

	v1 "k8s.io/api/apps/v1"
	v1beta1 "k8s.io/api/apps/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
)

func checkJobStatus(jobList *batchv1.JobList) (healthCheckPrinter *tabwriter.Writer, healthCheckStatus bool, hardFail bool) {
	printer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	healthCheckStatus = true
	hardFail = false
	fmt.Fprintln(printer, "EVENT\tJOB\tCONTAINER\tREASON\tSTATUS")
	for _, job := range jobList.Items {
		if job.Status.Succeeded != int32(1) {
			fmt.Fprintf(printer, "[ERROR]\t%s\t%s\tJOB IS NOT COMPLETE\tThe pod is still running\n", job.Name, job.Name)

			healthCheckStatus = false
		}

		if int(job.Status.Failed) == int(*job.Spec.BackoffLimit+1) {
			hardFail = true
		}
	}

	return printer, healthCheckStatus, hardFail
}

func checkDeploymentStatus(deploymentList *v1.DeploymentList) (healthCheckPrinter *tabwriter.Writer, healthCheckStatus bool) {

	printer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	healthCheckStatus = true

	for _, deployment := range deploymentList.Items {
		if deployment.Status.AvailableReplicas < 1 {
			healthCheckStatus = false
			fmt.Fprintf(printer, "[ERROR]\t%s\t%s\tDEPLOYMENT IS NOT AVAILABLE\tNo pods are ready\n", deployment.Name, deployment.Name)

		}
	}

	return printer, healthCheckStatus
}

func checkStatefulSetStatus(statefulSetList *v1beta1.StatefulSetList) (healthCheckPrinter *tabwriter.Writer, healthCheckStatus bool) {

	printer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	healthCheckStatus = true

	for _, statefulSet := range statefulSetList.Items {

		if statefulSet.Status.CurrentReplicas != statefulSet.Status.ReadyReplicas {
			healthCheckStatus = false
			fmt.Fprintf(printer, "[ERROR]\t%s\t%s\tSTATEFUL SETS ARE NOT READY\tpods are not ready\n", statefulSet.Name, statefulSet.Name)

		}
	}
	return printer, healthCheckStatus
}
