package kubehealth

import (
	"log"

	v1 "k8s.io/api/apps/v1"
	v1beta1 "k8s.io/api/apps/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (healthcheck *Check) getJobList(namespace, labels string, clientSet *kubernetes.Clientset) (jobList *batchv1.JobList) {

	opts := metav1.ListOptions{LabelSelector: labels}

	jobList, err := clientSet.BatchV1().Jobs(namespace).List(opts)
	if err != nil {
		log.Fatalf("[ERROR] Failed to list jobs. Stacktrace below:\n %v", err)
	}

	logResourceCounts("jobs", len(jobList.Items), labels, namespace)

	return jobList
}

func (healthcheck *Check) getDeploymentList(namespace, labels string, clientSet *kubernetes.Clientset) (deploymentList *v1.DeploymentList) {

	opts := metav1.ListOptions{LabelSelector: labels}

	deploymentList, err := clientSet.AppsV1().Deployments(namespace).List(opts)

	if err != nil {
		log.Fatalf("[ERROR] Failed to list jobs. Stacktrace below:\n %v", err)
	}

	logResourceCounts("deployments", len(deploymentList.Items), labels, namespace)

	return deploymentList
}

func (healthcheck *Check) getStatefulSetList(namespace, labels string, clientSet *kubernetes.Clientset) (statefulSetList *v1beta1.StatefulSetList) {

	opts := metav1.ListOptions{LabelSelector: labels}

	statefulSetList, err := clientSet.AppsV1beta1().StatefulSets(namespace).List(opts)
	if err != nil {
		log.Fatalf("[ERROR] Failed to list statefulSets. Stacktrace below:\n %v", err)
	}

	logResourceCounts("statefulsets", len(statefulSetList.Items), labels, namespace)

	return statefulSetList
}
