package kubehealth

import (
	"fmt"
	"log"
	"sync"
	"text/tabwriter"
	"time"
)

var (
	// DefaultHealthCheckRetries - amount of times to retry if healthcheck fails
	DefaultHealthCheckRetries = 50
	// DefaultHealthCheckSleep - amount of time to wait between retries (seconds)
	DefaultHealthCheckSleep      = 5
	jobHealthCheckStatus         = false
	deploymentHealthCheckStatus  = false
	statefulSetHealthcheckStatus = false
	healthCheckRetries           = 5
	healthCheckBanner            = `
==============================================================
	PERFORMING HEALTHCHECK`
	printer  *tabwriter.Writer
	hardFail = false
)

// ResourceList - Resources to conduct healthcheck against
type ResourceList struct {
	StatefulSet bool
	Deployment  bool
	Job         bool
}

// Initialise - main method for health check
func (healthCheck *Check) Initialise() {

	if healthCheck.Namespace == "" {
		log.Fatalln(`No namespace provided.
Example: ./helmplugin healthcheck start -n iptcwidev4 -l key=value,key2=value2 -r 50 -s 5 -R deployments`)
	}

	clientset := healthCheck.configureClient()

	resourceCheck := healthCheck.parseResources()

	fmt.Println(healthCheckBanner)
	fmt.Println("[INFO] Healthcheck verbosity:", healthCheck.IsVerbose)
	fmt.Println("[CONTEXT]:", healthCheck.KubeContext, "[RETRIES]:", healthCheck.HealthCheckRetries, "[RETRYSLEEP]:", healthCheck.HealthCheckSleep, "seconds")

	sleepDuration := time.Duration(int64(healthCheck.HealthCheckSleep)) * time.Second

	// Retry loop
	for index := 0; index < healthCheck.HealthCheckRetries; index++ {
		var wg sync.WaitGroup
		try := index + 1
		fmt.Printf("Healthcheck try %v of %v\n", try, healthCheck.HealthCheckRetries)

		//if type is one shot only get oneshot pods
		//else get all pods and oneshot pods
		//Get pod list

		// Check health of pods for all cases
		if resourceCheck.Deployment {
			wg.Add(1)
			go healthCheck.healthCheckDeployment(clientset, &wg, &resourceCheck.Deployment)
		}
		if resourceCheck.Job {
			wg.Add(1)
			go healthCheck.healthCheckJob(clientset, &wg, &resourceCheck.Job)
		}
		if resourceCheck.StatefulSet {
			wg.Add(1)
			go healthCheck.healthCheckStatefulSet(clientset, &wg, &resourceCheck.StatefulSet)

		}
		wg.Wait()

		if hardFail {
			break
		}
		if resourceCheck.Deployment || resourceCheck.Job || resourceCheck.StatefulSet {
			if healthCheck.IsVerbose {
				printer.Flush()
			}
			fmt.Printf("Healthcheck failed on try %v/%v - retrying after sleep: %v seconds\n", try, healthCheck.HealthCheckRetries, healthCheck.HealthCheckSleep)

			// Dont sleep on last try
			if try != healthCheck.HealthCheckRetries {
				time.Sleep(sleepDuration)
			}

			continue

		} else {

			break
		}

	}

	if resourceCheck.Deployment || resourceCheck.Job || resourceCheck.StatefulSet {
		printer.Flush()
		log.Fatalln("Healthcheck failed --------- details above")
	} else {
		fmt.Println("Healthcheck succeeded")
	}
}
