package kubernetes

import (
	"os/user"

	"github.com/sukeesh/k8s-job-notify/env"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client represents the wrapper of kubernetes API client
type Client struct {
	clientset kubernetes.Interface
}

// NewClient returns Client struct
func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientSet,
	}, nil
}

func getConfig() (config *rest.Config, err error) {
	var path string

	if !env.IsInCluster() {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}

		path = usr.HomeDir + "/.kube/config"
	}

	config, err = clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ListJobs returns the list of Jobs
func (c *Client) ListJobs(namespace string) (*batchv1.JobList, error) {
	jobs, err := c.clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return jobs, err
}
