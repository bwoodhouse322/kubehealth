package kubehealth

import (
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"
)

// Check - struct containing all configurations for healthcheck
type Check struct {
	Namespace          string
	KubeContext        string
	Labels             string
	IsVerbose          bool
	HealthCheckRetries int
	HealthCheckSleep   int
	SkipEmptyNamespace bool
	Resources          string
}

func (healthCheck *Check) healthCheckJob(clientset *kubernetes.Clientset, wg *sync.WaitGroup, check *bool) {

	jobList := healthCheck.getJobList(healthCheck.Namespace, healthCheck.Labels, clientset)

	if !jobHealthCheckStatus {
		printer, jobHealthCheckStatus, hardFail = checkJobStatus(jobList)

	} else {
		*check = false
	}
	wg.Done()

}

func (healthCheck *Check) healthCheckDeployment(clientset *kubernetes.Clientset, wg *sync.WaitGroup, check *bool) {

	deploymentList := healthCheck.getDeploymentList(healthCheck.Namespace, healthCheck.Labels, clientset)

	if !deploymentHealthCheckStatus {
		printer, deploymentHealthCheckStatus = checkDeploymentStatus(deploymentList)

	} else {
		*check = false

	}
	wg.Done()

}

func (healthCheck *Check) healthCheckStatefulSet(clientset *kubernetes.Clientset, wg *sync.WaitGroup, check *bool) {

	statefulSetList := healthCheck.getStatefulSetList(healthCheck.Namespace, healthCheck.Labels, clientset)

	if !statefulSetHealthcheckStatus {
		printer, statefulSetHealthcheckStatus = checkStatefulSetStatus(statefulSetList)
	} else {
		*check = false
	}
	wg.Done()

}

func (healthCheck Check) parseResources() ResourceList {
	parsedResources := strings.Split(healthCheck.Resources, ",")
	resourceList := ResourceList{Deployment: false, StatefulSet: false, Job: false}

	for _, resource := range parsedResources {
		switch resource {
		case "deployment", "deploy", "deployments":
			resourceList.Deployment = true
		case "job", "jobs":
			resourceList.Job = true
		case "statefulset", "ss", "statefulsets":
			resourceList.StatefulSet = true

		}

	}
	if resourceList.Deployment == false && resourceList.Job == false && resourceList.StatefulSet == false {

		panic("No Resources Selected")
	}
	return resourceList
}
