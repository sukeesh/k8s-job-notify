package main

import (
	"flag"
	"time"

	"github.com/sukeesh/cron-k8s-watch/env"
	"github.com/sukeesh/cron-k8s-watch/message"
	"github.com/sukeesh/cron-k8s-watch/slack"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "/Users/sukeesh/.kube/config", "absolute path to file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		namespace := env.GetNamespace()
		jobs, err := clientSet.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, job := range jobs.Items {
			if job.Status.StartTime.Time.Add(time.Hour * 12).After(time.Now()) {
				if job.Status.Succeeded > 0 {
					err = slack.SendSlackMessage(message.JobSuccess(job.Name, job.Status.CompletionTime.String()))
					if err != nil {
						panic(err.Error())
					}
				} else if job.Status.Failed > 0 {
					err = slack.SendSlackMessage(message.JobFailure(job.Name))
					if err != nil {
						panic(err.Error())
					}
				}
			}
		}
		break
	}
}
