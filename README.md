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
