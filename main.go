package main

import (
	"flag"
	"os"
	"os/user"
	"time"

	"github.com/sukeesh/k8s-job-notify/env"
	"github.com/sukeesh/k8s-job-notify/message"
	"github.com/sukeesh/k8s-job-notify/slack"

	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	pastJobs map[string]bool
)

func main() {
	var kubeconfig *string
	var config *rest.Config
	var err error

	pastJobs := make(map[string]bool)
	if env.IsInCluster() {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		usr, err := user.Current()
		if err != nil {
			panic(err.Error())
		}
		filePath := usr.HomeDir + "/.kube/config"
		kubeconfig = flag.String("kubeconfig", filePath, "absolute path to file")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := env.GetNamespace()
	for {
		jobs, err := clientSet.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, job := range jobs.Items {
			if pastJobs[job.Name] == false && job.Status.StartTime.Time.Add(time.Minute*10).After(time.Now()) {
				if job.Status.Succeeded > 0 {
					timeSinceCompletion := time.Now().Sub(job.Status.CompletionTime.Time).Minutes()
					err = slack.SendSlackMessage(message.JobSuccess(job.Name, timeSinceCompletion))
					if err != nil {
						panic(err.Error())
					}
					pastJobs[job.Name] = true
				} else if job.Status.Failed > 0 {
					err = slack.SendSlackMessage(message.JobFailure(job.Name))
					if err != nil {
						panic(err.Error())
					}
					pastJobs[job.Name] = true
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
	os.Exit(0)
}
