package env

import (
	"errors"
	"os"
)

func GetNamespace() (namespace string) {
	if namespace = os.Getenv("namespace"); namespace == "" {
		namespace = "default"
	}
	return namespace
}

func GetSlackWebHookURL() (webhook string, err error) {
	if webhook = os.Getenv("webhook"); webhook == "" {
		return "", errors.New("environment variable $webhook not set")
	}
	return webhook, nil
}

func GetNotificationLevel() (level string) {
	if level = os.Getenv("notification_level"); level == "" {
		level = "all"
	}
	return level
}


func IsInCluster() bool {
	inCluster := os.Getenv("incluster")
	return inCluster == "1"
}
