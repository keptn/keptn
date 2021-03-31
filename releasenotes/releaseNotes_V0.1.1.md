# Release Notes v0.1.1

This is a bugfix release of version 0.1.0. It does not add any usecases but does address several bugs that prevented automated setup of some components and prevented frictionless operations of the use cases.

The updated instructions can be found on the [keptn.sh website](https://keptn.sh/docs/0.1.1/).

Fixed bugs:
- fixed outdated Jenkins configuration
- fixed pipelines that where tagging docker images incorrectly
- fixed wrong URL in Ansible script
- fixed null value bug in production pipelines

Improvements:
- added logging to `setupInfrastructure.sh` and `forkGitHubRepositories.sh` 
- using Ansible vault for storing credentials
- automatically create Dynatrace request attributes during setup
