package message

import (
	"fmt"
)

func JobSuccess(clusterName, jobName string, timeSinceCompletion float64) string {
	return "*" + clusterName + ": " + jobName + "* succeeded " + fmt.Sprintf("%f", timeSinceCompletion) + " minutes ago :tada:"
}

func JobFailure(clusterName, jobName string) string {
	return "*" + clusterName + ": " + jobName + "* failed :alert:"
}
