package kubehealth

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Note to self: Comment later please Abdul - use cobra instead
func checkParams() (namespace, labels string) {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: %v namespace 'labelkey:labelvalue'\n", os.Args[0])
		fmt.Println("Label args can be left blank. Accept a comma seperatedlist of 'key=value,key2=value'")
		os.Exit(1)
	} else {
		fmt.Println("Performing pods for namespace:", os.Args[1])
	}

	namespace = os.Args[1]

	if len(os.Args) < 3 {
		labels = ""
	} else {
		labels = os.Args[2]
	}

	return namespace, labels
}

func logResourceCounts(resourceType string, size int, labels string, namespace string) {

	if size < 1 {
		log.Fatalf("[ERROR] No %s found in the namespace matching given criteria.", resourceType)
	}

	fmt.Printf("[NAMESPACE]: %s, [ %s COUNT]: %v, [LABELS]: %v\n", namespace, strings.ToUpper(resourceType), size, labels)
	fmt.Println("==============================================================")

}
