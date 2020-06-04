package message

import (
	"fmt"
)

func JobSuccess(jobName string, timeSinceCompletion float64) string {
	return "*" + jobName + "* succeeded " + fmt.Sprintf("%.1f", timeSinceCompletion) + " minutes ago :tada:"
}

func JobFailure(jobName string) string {
	return "*" + jobName + "* failed :alert:"
}
