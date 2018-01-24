package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"context"
)

type DockerSwarmDeployer struct {
	client *docker.Client
}

func NewDockerSwarmDeployer() (Deployer, error) {
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &DockerSwarmDeployer{client}, nil
}

func (deployer DockerSwarmDeployer) getService(name string) (*swarm.Service, error) {
	ctx := context.Background()
	service, _, err := deployer.client.ServiceInspectWithRaw(ctx, name, types.ServiceInspectOptions{})
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (deployer DockerSwarmDeployer) Update(name, image string) error {
	// get service
	service, err := deployer.getService(name)
	if err != nil {
		return err
	}

	newSpec := service.Spec
	newSpec.TaskTemplate.ContainerSpec.Image = image

	// perform change
	ctx := context.Background()
	options := types.ServiceUpdateOptions{}
	_, err = deployer.client.ServiceUpdate(ctx, service.ID, service.Version, newSpec, options)
	return err
}