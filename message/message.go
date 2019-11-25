package message

func JobSuccess(jobName string, completedTime string) string {
	return "*" + jobName + "* succeeded at " + completedTime + " :tada:"
}

func JobFailure(jobName string) string {
	return "*" + jobName + "* failed :alert:"
}
