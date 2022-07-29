# Lagoon Synchronise Keycloak Attributes (SKA)

[![Release](https://github.com/smlx/lagoon-ska/actions/workflows/release.yaml/badge.svg)](https://github.com/smlx/lagoon-ska/actions/workflows/release.yaml)
[![Coverage](https://coveralls.io/repos/github/smlx/lagoon-ska/badge.svg?branch=main)](https://coveralls.io/github/smlx/lagoon-ska?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/smlx/lagoon-ska)](https://goreportcard.com/report/github.com/smlx/lagoon-ska)

### What this does

Checks Keycloak groups in the Lagoon realm for the `group-lagoon-project-ids` attribute required by the `ssh-portal`.
If the attribute doesn't exist, it is synthesised from the group's `lagoon-projects` attribute.

### How to use

Copy the binary into the Lagoon API pod.
```bash
kubectl cp ~/download/lagoon-ska lagoon-core-api-74bffc6876-lklzf:/tmp/lagoon-ska
```

Dump the groups in JSON format as a backup.
```bash
$ kubectl exec -it lagoon-core-api-74bffc6876-lklzf -- /tmp/lagoon-ska dump > /tmp/groups.before.json
```

Run in dry-run mode (the default) first to see what will change.
```bash
$ kubectl exec -it lagoon-core-api-74bffc6876-lklzf -- /tmp/lagoon-ska
{"level":"info","ts":1659065140.0172124,"caller":"lagoon-ska/sync.go:34","msg":"not making any changes in dry-run mode"}
{"level":"info","ts":1659065140.0215838,"caller":"keycloak/sync.go:78","msg":"groups found with missing attributes","count":2}
{"level":"info","ts":1659065140.021606,"caller":"keycloak/sync.go:83","msg":"missing group-lagoon-project-ids attribute on group","name":"ci-group","ID":"01f8c347-62ef-4f6c-b1dc-5c1050bf5155"}
{"level":"info","ts":1659065140.0216126,"caller":"keycloak/sync.go:83","msg":"missing group-lagoon-project-ids attribute on group","name":"project-ci-drush-la-control-k8s","ID":"24ad6257-502c-407a-bfac-52d75b9ed8a2"}
```

Run the sync to update the group attributes.
```bash
$ kubectl exec -it lagoon-core-api-74bffc6876-lklzf -- /tmp/lagoon-ska sync --dry-run=false
{"level":"info","ts":1659065225.1785626,"caller":"keycloak/sync.go:78","msg":"groups found with missing attributes","count":2}
{"level":"info","ts":1659065225.181286,"caller":"keycloak/sync.go:97","msg":"updated group-lagoon-project-ids attribute on group","name":"ci-group","ID":"01f8c347-62ef-4f6c-b1dc-5c1050bf5155"}
{"level":"info","ts":1659065225.1840253,"caller":"keycloak/sync.go:97","msg":"updated group-lagoon-project-ids attribute on group","name":"project-ci-drush-la-control-k8s","ID":"24ad6257-502c-407a-bfac-52d75b9ed8a2"}
```

Dump the groups again (and diff against the previous dump to confirm the changes).
```bash
$ kubectl exec -it lagoon-core-api-74bffc6876-lklzf -- /tmp/lagoon-ska dump > /tmp/groups.after.json
```
