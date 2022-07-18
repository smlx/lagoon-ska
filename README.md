# Lagoon Synchronise Keycloak Annotations (SKA)

[![Release](https://github.com/smlx/lagoon-ska/actions/workflows/release.yaml/badge.svg)](https://github.com/smlx/lagoon-ska/actions/workflows/release.yaml)
[![Coverage](https://coveralls.io/repos/github/smlx/lagoon-ska/badge.svg?branch=main)](https://coveralls.io/github/smlx/lagoon-ska?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/smlx/lagoon-ska)](https://goreportcard.com/report/github.com/smlx/lagoon-ska)

### What this does

Checks Keycloak groups in the Lagoon realm for the `group-lagoon-project-ids` annotation required by the `ssh-portal`.
If the annotation doesn't exist, it is synthesised from the group's `lagoon-projects` annotation.

### How to use

Run the binary in the Lagoon API pod.
Use `./lagoon-ska sync --dry-run` to see what will change without updating any annotations.
