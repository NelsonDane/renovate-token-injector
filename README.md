# Renovate GitHub App Token Injector for Kubernetes

Generate GitHub App Tokens for Renovate

For some reason Renovate only supports GitHub App auth on their closed source version. See [this issue](https://github.com/renovatebot/renovate/discussions/21035) for more details. Because Renovate only reads its config once at startup, it was impossible to "patch" the config with a new token (since they're short-lived). But using a [new feature](https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-via-file/) in Kubernetes `v1.34`, init containers can now inject env variables into the main container BEFORE it starts. This allows us to generate a fresh GitHub App token for Renovate on each pod cron job start.

Used in:
- [Escargatoire Homelab Cluster](https://github.com/NelsonDane/escargatoire)
