package main

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	k8s "./kubernetes"
	"github.com/sukeesh/k8s-job-notify/env"
	"github.com/sukeesh/k8s-job-notify/message"
	"github.com/sukeesh/k8s-job-notify/slack"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	pastJobs := make(map[string]bool)

	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("failed to create client %v", zap.Error(err))
		os.Exit(1)
	}

	namespace := env.GetNamespace()
	log.Printf("fetching jobs from %s namespace", namespace)
	for {
		jobs, err := client.ListJobs(namespace)
		if err != nil {
			log.Fatalf("failed to list all jobs in the namespace %v", zap.Error(err))
			os.Exit(1)
		}

		for _, job := range jobs.Items {
			// job.Name can be same for different jobs, so using job.Name+CreationTimeStamp for checking
			// uniqueness of the job. so that duplicated messages to slack can be avoided
			jobUniqueHash := job.Name + job.CreationTimestamp.String()
			if pastJobs[jobUniqueHash] == false && job.Status.StartTime.Time.Add(time.Minute*20).After(time.Now()) {
				if job.Status.Succeeded > 0 {
					timeSinceCompletion := time.Now().Sub(job.Status.CompletionTime.Time).Minutes()
					err := slack.SendSlackMessage(message.JobSuccess(job.Name, timeSinceCompletion))
					if err != nil {
						log.Fatalf("sending a message to slack failed %v", zap.Error(err))
					}
					pastJobs[jobUniqueHash] = true
				} else if job.Status.Failed > 0 {
					err := slack.SendSlackMessage(message.JobFailure(job.Name))
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
}
