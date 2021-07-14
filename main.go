package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"

	"github.com/sukeesh/k8s-job-notify/env"
	k8s "github.com/sukeesh/k8s-job-notify/kubernetes"
	"github.com/sukeesh/k8s-job-notify/message"
	"github.com/sukeesh/k8s-job-notify/slack"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var opts struct {
	ClusterName string `long:"cluster-name" description:"Show cluster name in message (optional)"`
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	clusterName := opts.ClusterName
	pastJobs := make(map[string]bool)

	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("failed to create client %v", zap.Error(err))
		os.Exit(1)
	}

	namespace := env.GetNamespace()
	log.Printf("fetching jobs from %s namespace", namespace)
	level := env.GetNotificationLevel()
	log.Printf("notification_level set at '%s'", level)
	for {
		jobs, err := client.ListJobs(namespace)
		if err != nil {
			log.Fatalf("failed to list all jobs in the namespace %v", zap.Error(err))
			continue
		}

		for _, job := range jobs.Items {
			// job.Name can be same for different jobs, so using job.Name+CreationTimeStamp for checking
			// uniqueness of the job. so that duplicated messages to slack can be avoided
			jobUniqueHash := job.Name + job.CreationTimestamp.String()
			if pastJobs[jobUniqueHash] == false {
				if level == "all" {
					// Send success notifications.
					if job.Status.Succeeded > 0 &&  (job.Status.CompletionTime.Add(20*time.Minute).Unix() > time.Now().Unix()) {
						timeSinceCompletion := time.Now().Sub(job.Status.CompletionTime.Time).Minutes()
						err = slack.SendSlackMessage(message.JobSuccess(clusterName, job.Name, timeSinceCompletion))
						if err != nil {
							log.Fatalf("sending a message to slack failed %v", zap.Error(err))
						}
						pastJobs[jobUniqueHash] = true
					}
				}

				if level == "failed" || level == "all" {
					// Send failed notifications.
				  if job.Status.Failed > 0 {
						if job.Status.StartTime.Add(5*time.Hour).Unix() > time.Now().Unix() {
							err = slack.SendSlackMessage(message.JobFailure(clusterName, job.Name))
							if err != nil {
								log.Fatalf("sending a message to slack failed %v", zap.Error(err))
							}
							pastJobs[jobUniqueHash] = true
						}
					}
				}
			}
		}
		time.Sleep(time.Minute * 1)
		log.Printf("end of 1 minute wait.. fetching new jobs")
	}
}
