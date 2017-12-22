package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Config struct {
	Address string `default:":7070"`
	Token   string `default:""`
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

type Deployer interface {
	Update(service, image string) error
}

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

type Handler struct {
	deployer Deployer
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		log.Error(err)
		return
	}

	service := r.Form.Get("service")
	image := r.Form.Get("image")
	if service == "" || image == "" {
		http.Error(w, "Empty parameters", http.StatusBadRequest)
		log.Error(err)
		return
	}

	err = h.deployer.Update(service, image)
	if err != nil {
		message := fmt.Sprintf("Update failed: %s", err.Error())
		http.Error(w, message, http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	log.Infof("Started update of %s to %s", service, image)
}

func authorized(validToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != validToken {
			log.Warnf("Unauthorized request from %s", r.RemoteAddr)
			w.WriteHeader(http.StatusAccepted)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	var config Config
	err := envconfig.Process("deployd", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	deployer, err := NewDockerSwarmDeployer()
	if err != nil {
		log.Fatal(err)
	}

	handler := &Handler{deployer}

	var update http.Handler
	update = http.HandlerFunc(handler.HandleUpdate)
	if config.Token != "" {
		update = authorized(config.Token, update)
	}
	http.Handle("/update", update)

	log.Infof("Starting server listening at %s", config.Address)

	log.Fatal(http.ListenAndServe(config.Address, nil))
}
