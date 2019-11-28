package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"time"

	"go.uber.org/zap"

	"github.com/sukeesh/k8s-job-notify/env"
	"github.com/sukeesh/k8s-job-notify/message"
	"github.com/sukeesh/k8s-job-notify/slack"

	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
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
		log.Printf("using inClusterConfig")
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
	log.Printf("fetching jobs from %s namespace", namespace)
	for {
		jobs, err := clientSet.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
		if err != nil {
			log.Fatalf("failed to list all jobs in the namespace %v", zap.Error(err))
			os.Exit(1)
		}
		for _, job := range jobs.Items {
			// job.Name can be unique, so using job.Name+CreationTimeStamp for checking uniqueness of the job
			// so that duplicated messages to slack can be avoided
			jobUniqueHash := job.Name + job.CreationTimestamp.String()
			if pastJobs[jobUniqueHash] == false && job.Status.StartTime.Time.Add(time.Minute*20).After(time.Now()) {
				if job.Status.Succeeded > 0 {
					timeSinceCompletion := time.Now().Sub(job.Status.CompletionTime.Time).Minutes()
					err = slack.SendSlackMessage(message.JobSuccess(job.Name, timeSinceCompletion))
					if err != nil {
						log.Fatalf("sending a message to slack failed %v", zap.Error(err))
					}
					pastJobs[jobUniqueHash] = true
				} else if job.Status.Failed > 0 {
					err = slack.SendSlackMessage(message.JobFailure(job.Name))
					if err != nil {
						log.Fatalf("sending a message to slack failed %v", zap.Error(err))
					}
					pastJobs[jobUniqueHash] = true
				}
			}
		}
		time.Sleep(time.Minute * 1)
		log.Printf("end of 1 minute wait.. fetching new jobs")
	}
	os.Exit(0)
}
