## Contributing

First off, thank you for considering contributing to Keptn. It's people like you that make Keptn great.

If you've noticed a bug, want to contribute features, or simply ask a question that for whatever reason you don't want to ask in the [Keptn Slack workspace](keptn.slack.com), please [search the issue tracker](https://github.com/keptn/keptn/issues?q=something) to see if someone else in the community has already created a ticket. If not, go ahead and [make one](https://github.com/keptn/keptn/issues/new).

If you want to work on an issue and contribute code, this is the right document to get started.

### Read the docs

We are in the process of making sure that each repository and each service within the Keptn organization is documented properly. 
We are aware that some parts are currently missing, in the meantime please get in touch with us through the [Keptn Slack workspace](keptn.slack.com) if you have any questions.

As a starting point, please read the docs within the [docs/](docs/) folder in this repository.


### Tell us if you start working on an issue

In case you want to work on an unassigned [GitHub issue](https://github.com/keptn/keptn/issues), please let us know via
 a quick comment in the issue itself. We will then try to assign it to you.

### Propose/design new features

Proposing new functionality for Keptn is done via a so called [Keptn Enhancement Proposal (KEP)](https://github.com/keptn/enhancement-proposals).
This is required when it is intended to introduce new behaviour, change desired behaviour, or otherwise modify requirements of Keptn.

If all you want to propose is a simple addition to a Keptn-service or the CLI, or even a bugfix, you can just open 
[a new issue on GitHub](https://github.com/keptn/keptn/issues/new/choose).

### Follow Coding style

When contributing code to Keptn, we politely ask you to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).
We are running automated code style checks for pull requests using the following tools:

* [reviewdog](.reviewdog.yml) - automatic code review
  * ``golint ./...``
  * ``gofmt -l -s .`` 
* [codecov.io](codecov.yml) - tracks code coverage based on unit tests

### Fork and work in a branch

For contributions to Keptn, please fork the Keptn repo and work in a branch. More information about forking is written
down in our [docs](docs/fork.md).

We are following the [git branching model described in this blog post](https://nvie.com/posts/a-successful-git-branching-model/), however we try to avoid the extra step of the `develop` branch, and instead work directly on the `master` branch.

If you work on a new feature or on a bugfix, [fork the repository](docs/fork.md), branch away from the `master` branch
and use the following schema for naming your branches:

* `feature/###/name` for new features,
* `patch/###/name` for patches,
* `bug/###/name` for bugfixes, and
* `hotfix/###/name` for hotfixes (e.g., for releases),

where `###` is the github issue id. E.g., if you work on feature described on issue #123, the branch name could be

* `feature/123/foo`.

For Hotfixes please branch away from the master branch, and create a PR to master and the respective release branches.

### Run Tests

Keptn currently has two types of tests:

* end-to-end tests (integration tests) - run automatically once per day (see [test/](test/))
* unit tests - run for every Pull Request in every directory that contains Go code

Before pushing your code to the repository, please run unit tests locally. When creating features, please also consider
 writing unit tests.

### Deploy your local changes to an existing cluster

If you are changing behaviour or a large part of the code, please verify Keptn still works as it's supposed to, by following the tutorials described on the [Keptn website](https://tutorials.keptn.sh).

To deploy your local changes to an existing Kubernetes cluster with Keptn running on it, we recommend using [skaffold](https://skaffold.dev).
We provide a `skaffold.yaml` file in every repository/directory, which you can use to automatically deploy the service using
```console
skaffold run --tail --default-repo=your-docker-registry
```

Please replace `your-docker-registry` with your DockerHub username and repository name.
This should 

* automatically build the docker image for the service,
* push the docker image to the defined container registry/repository,
* deploy the service to the Kubernetes cluster using the image that was just built, and
* print log output to your terminal.

In case you are using JetBrains GoLand, you can also use the built-in debugging features using Skaffold and Google Cloud Code as described [here](docs/debugging.md).


### Make a Pull Request

At this point, you should switch back to the `master` branch in your repo, and make sure it's up to date with Keptn's `master` branch:

```bash
git remote add upstream git@github.com:keptn/keptn.git
git checkout master
git pull upstream master
```

Then update your feature branch from your local copy of `master`, and push it.

```bash
git checkout feature/123/foo
git rebase master
git push --set-upstream origin feature/123/foo
```

Finally, go to GitHub and make a Pull Request. Please describe what this PR is about and add a link to relevant GitHub issues.
Your PR will usually be reviewed by the Keptn team within a couple of days, but feel free to let us know about your PR [via Slack](https://slack.keptn.sh).
