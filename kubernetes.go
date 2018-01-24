package main

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"fmt"
	"k8s.io/client-go/rest"
)

type KubernetesDeployer struct {
	clientset *kubernetes.Clientset
}

func NewKubernetesDeployer() (Deployer, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubernetesDeployer{clientset}, nil
}

func (deployer KubernetesDeployer) Update(name, image string) error {
	deploymentsClient := deployer.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, err := deploymentsClient.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		updated := false
		for i := range deployment.Spec.Template.Spec.Containers {
			if deployment.Spec.Template.Spec.Containers[i].Name == name {
				deployment.Spec.Template.Spec.Containers[i].Image = image
				updated = true
				break
			}
		}

		if !updated {
			return fmt.Errorf("couldn't find container spec with name: %v", name)
		}

		_, updateErr := deploymentsClient.Update(deployment)
		return updateErr
	})
}
