# deployd

Update Kubernetes deployments and Docker Swarm services using webhooks


## Kubernetes

- Create a service account:
  `kubectl create serviceaccount deployd`
- Create a role binding allowing the service account to perform edits:
  `kubectl apply -f role-binding.yaml`
- Create a secret named `deployd-secret` with key `token`
- Create the deployment:
  `kubectl apply -f deployment.yaml`
- Create the service:
  `kubectl apply -f service.yaml`
- Create an ingress

## Docker Swarm

- See `docker-compose.yml`
- Make sure to set environment variable `DEPLOYD_TOKEN` to a new secret


## Updating deployments/services

Updates of deployments/services can be triggered using, e.g.

```sh
wget --content-on-error --header "Authorization: $DEPLOYD_TOKEN" -qO- \
    https://$DEPLOYD_HOST/update\?name\=$DEPLOYD_NAME\&image\=$DEPLOYD_IMAGE:$DEPLOYD_TAG
```

Here, `$DEPLOYD_TOKEN` is the secret token that was specified in the configuration, `$DEPLOYD_NAME` is the name of the service/deployment to be updated, `$DEPLOYD_IMAGE` is the name of the Docker image that should be used, and `$DEPLOYD_TAG` is the tag of the Docker image that should be used.