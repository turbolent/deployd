package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		log.Error(err)
		return
	}

	name := r.Form.Get("name")
	image := r.Form.Get("image")
	if name == "" || image == "" {
		http.Error(w, "Empty parameters", http.StatusBadRequest)
		log.Error(err)
		return
	}

	err = h.deployer.Update(name, image)
	if err != nil {
		message := fmt.Sprintf("Update failed: %s", err.Error())
		http.Error(w, message, http.StatusInternalServerError)
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	log.Infof("Started update of %s to %s", name, image)
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

type DeployerConstructor func()(Deployer, error)

var modes = map[string]DeployerConstructor {
	"docker": NewDockerSwarmDeployer,
	"kubernetes": NewKubernetesDeployer,
}

func main() {
	var config Config
	err := envconfig.Process("deployd", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	constructor, ok := modes[config.Mode]
	if !ok {
		log.Fatalf("Invalid mode: %v", config.Mode)
	}

	deployer, err := constructor()
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
