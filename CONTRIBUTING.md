# Contributing

First off, thank you for considering contributing to Keptn! It's people like you who make Keptn great.

* **Create an issue**: If you have noticed a bug, want to contribute features, or simply ask a question that for whatever reason you do not want to ask in the [Keptn Slack workspace](https://slack.keptn.sh), please [search the issue tracker](https://github.com/keptn/keptn/issues?q=something) to see if someone else in the community has already created a ticket. If not, go ahead and [create an issue](https://github.com/keptn/keptn/issues/new).

* **Start contributing**: We also have a list of [good first issues](https://github.com/keptn/keptn/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22). If you want to work on it, just post a comment on the issue.

* **Add yourself**: Add yourself to the [list of contributors](CONTRIBUTORS.md) along with your first pull request.

This document lays out how to get you started in contributing to Keptn, so please read on.

## Documentation

The Keptn repository is a monorepo with multiple sub-folders that contain microservices, the installer, Keptn Bridge, 
and the `keptn` CLI. 
As a starting point, please read the README files in the respective sub-folder. Also, take a look at the docs within the [docs/](docs/) folder in this repository.

We are aware that not every microservice has comprehensive documentation, so if you have any questions feel free to ask us in the [Keptn Slack workspace](https://slack.keptn.sh).

## Development Process

Within the Keptn project several features of GitHub are used to control the development process. For instance, you can
find the [Roadmap](https://github.com/orgs/keptn/projects/1) as a [GitHub Project Board](https://github.com/orgs/keptn/projects).

In addition, we use GitHub Issues as the primary method for tracking concrete implementation tasks, bugs, etc...

### Keptn Enhancement Proposals and GitHub Issues

Proposing new functionality for Keptn is a transparent process done via a so-called [Keptn Enhancement Proposal](https://github.com/keptn/enhancement-proposals).
This is required when the proposed functionality is intended to introduce new behaviour, change desired behaviour, or otherwise modify the requirements of Keptn.

If all you want to propose is a simple addition to a Keptn-service, Bridge or the CLI, a typo, or just a bugfix, please open [a new issue on GitHub](https://github.com/keptn/keptn/issues/new/choose) and try to label it accordingly.

If you are unsure, just open a [new issue](https://github.com/keptn/keptn/issues/new/choose) and we will label it for you or direct you to an enhancement proposal if necessary.

### Working items

We use GitHub Issues as the primary method for tracking concrete implementation tasks (e.g., an enhancement proposal can
lead to multiple issues). In addition, we use a [GitHub Project Board](https://github.com/keptn/keptn/projects) to 
plan and assign the working items for the next upcoming days. This board is discussed and presented publicly in our
[weekly community meeting](https://github.com/keptn/community#%EF%B8%8F-community-meetings).

### Bug first policy

We derived the following labels based on the Blog Post [Zero-Bug Software Development](https://www.xolv.io/blog/zero-bug-software-development/)
to classify the importance of work tasks with GitHub issues:

* `type:critical` - This is a **defect** that causes us and/or our end-users/stakeholders to lose a significant amount 
   of time, money and/or value.
* `type:bug` - Something is not working as intended or documented and therefore needs to be fixed (either by writing
   code or by adapting the documentation).
* `type:feature` - This indicates that the issue contains a complete new feature that provides value for our 
   end-users/stakeholders, most likely in combination with a `kep:XYZ` label
* `type:improvement` - This label should be used for indicating that something is just an improvement of an existing 
   feature.
* `type:chore` - Chores are issues that provide value to the (developer) team, e.g., cleanups, automation, etc... 
   (things that keep our team productive and efficient)

**Please note**: The Keptn maintainers reserve the right to change the labels to the best of their knowledge.

Based on the labels above, we try to adhere to the following bug first policy:

* Issues labeled with `type:critical` need to be addressed right away (either with a comment, a work-around, or a fix).
  The rule is: stop what you are doing, and fix it.
* Issues labeled with `type:bug` need to have a response within 3 business days, and a workaround or fix needs to be 
  provided within 10 business days.
* Issues labeled with `type:feature`, `type:improvement` or `type:chore` will be worked on as per the backlog priority
  order.

**Please note**: This is a policy, not an enforced rule. The Keptn maintainers will do their best to follow this policy.

### Other labels

There are many other labels that either indicate that a feature belongs to a certain part of Keptn (e.g., `bridge`, 
`cli`, `core`, ...) or to an enhancement proposal (e.g., `kep:06`). In addition, we use the labels `good first issue`
to indicate that this issue is good to get started with contributing to Keptn. 


### Fork Repositories and work in a branch

For contributions to Keptn, please fork the Keptn repository and work in a branch. More information about forking is written down in the [docs/fork](docs/fork.md).

We are loosely following the [Git Flow branching model](https://nvie.com/posts/a-successful-git-branching-model/) however, we try to avoid the extra step of the `develop` branch and instead work directly on the `master` branch.

* If you work on a new feature, [fork the repository](docs/fork.md), branch away from the `master` branch and use the following schema for naming your branches:
```
<ticket-type>/<github-issue-number>/<descriptive-name-with-dashes>

Examples:
feat/1234/my-new-feature
fix/2235/important-bugfix
```

All possible ticket/commit/PR types and scopes can be found [here](#commit-types-and-scopes).

* If you work on a hotfix and a maintenance branch for the related release(s) doesn't exist yet, please create a maintenance branch with the following naming scheme:
If the release was e.g. `v0.8.6`, the maintenance branch should be named `0.8.x` and  it should have the faulty release as its base instead of master.
Then, you can create a bugfix branch from the **master** branch and work on the hotfix. When you are done, please create a PR against the master branch.
The merged commit should then cherry-picked to the maintenance branch and then the hotfix release can be created.


**Note:** The `###` part of the branch has to reference the GitHub issue number, e.g., if you work on feature described in issue `#1234`, the branch name would be: `feature/1234/my-feature`.


### Follow Coding Style

When contributing code to Keptn, we politely ask you to follow the coding style suggested by the [Golang community](https://github.com/golang/go/wiki/CodeReviewComments).
We are running automated code style checks for pull requests using the following tools:

* [reviewdog](.reviewdog.yml) - automatic code review
  * ``golint ./...``
  * ``gofmt -l -s .`` 
* [codecov.io](codecov.yml) - tracks code coverage based on unit tests

### Boy Scout Rule

> Always leave the code better than you found it!

While we would all like to make the world better all the time, please apply common sense to this rule, for instance:

* When you're fixing a bug, don't refactor code around the bug (in fact, this could be counter-productive for backports/cherry-picks).
* Do not change or extend the scope of your issue (e.g., stay within the microservice/sub-directory).
* Don't refactor multiple parts at the same time - change one thing at a time.
* Make sure code-changes are balanced! Don't refactor 200 lines of codes when a simple change would only require 10 lines of code.
* If you file a PR, always ask yourself: What changes within the PR would you expect as a reviewer?

These rules help us to keep the scope of issues and PRs assessable.

### Run Tests

Keptn currently has two types of tests:

* *end-to-end tests* (integration tests) - run automatically once per day (see [test/](test/) for more details on the integration tests)
* *unit tests* - run for every Pull Request in every directory that contains Go code

Before pushing your code to the repository, please run unit tests locally. When creating a new feature, please consider writing unit tests.

### Deploy your local changes to an existing cluster

If you are changing behaviour or a large part of the code, please verify Keptn still works as it is supposed to do. This can be done, by following the tutorials available at [Keptn tutorials](https://tutorials.keptn.sh).

To deploy your local changes to an existing Kubernetes cluster with Keptn running on it, we recommend using [skaffold](https://skaffold.dev). We provide a `skaffold.yaml` file in every repository/directory, which you can be used to automatically deploy the service using:

```console
skaffold run --tail --default-repo=your-docker-registry
```

Please replace `your-docker-registry` with your DockerHub username and repository name.
Skaffold does then: 
* automatically build the docker image for the service,
* push the docker image to the defined container registry/repository,
* deploy the service to the Kubernetes cluster using the image that was just built, and
* print log output to your terminal.

In case you are using JetBrains GoLand, you can also use the built-in debugging features using Skaffold and Google Cloud Code as described [here](docs/debugging.md).


### Make a Pull Request

At this point, you should switch back to the `master` branch in your repository, and make sure it is up to date with `master` branch of Keptn:

```bash
git remote add upstream git@github.com:keptn/keptn.git
git checkout master
git pull upstream master
```

Then update your feature branch from your local copy of `master` and push it:

```bash
git checkout feature/123/foo
git rebase master
git push --set-upstream origin feature/123/foo
```

Finally, go to GitHub and create a Pull Request. There should be a PR template already prepared for you.
If not, you will find it at `.github/pull_request_template.md`.
Please describe what this PR is about and add a link to relevant GitHub issues.
If you changed something that is visible to the user, please add a screenshot.
Please follow the [conventional commit guidelines](https://www.conventionalcommits.org/en/v1.0.0/) for your PR title.

If you only have one commit in your PR, please follow the guidelines for the message of that single commit, otherwise the PR title is enough.
You can find a list of all possible feature types [here](#commit-types-and-scopes).

An example for a pull request title would be:
```
feat(api): New endpoint for feature X (#1234)
```
This would be a PR that adds a new endpoint to the keptn API and the issue number related to this PR is #1234.

If you have **breaking changes** in your PR, it is important to note them in the PR description but also in the merge commit for that PR.
When pressing "squash and merge", you have the option to fill out the commit message. Please use that feature to add the breaking changes according to the conventional commit guidelines.
Also, please remove the PR number at the end and just add the issue number.

An example for a PR with breaking changes and the according merge commit:
```
feat(bridge): New button that breaks other things (#345)

BREAKING CHANGE: The new button added with #345 introduces new functionality that is not compatible with the previous type of sent events.
```

If your breaking change can be explained in a single line you can also use this form:
```
feat(bridge)!: New button that breaks other things (#345)
```

Following those guidelines helps us create automated releases where the commit and PR messages are directly used in the changelog.

In addition, please always ask yourself the following questions:

**Based on the linked issue, what changes within the PR would you expect as a reviewer?**

Your PR will usually be reviewed by the Keptn team within a couple of days, but feel free to let us know about your PR [via Slack](https://slack.keptn.sh).

### Commit Types and Scopes
**Please find and up-to-date list of types and scopes [here](https://github.com/keptn/keptn/blob/master/.github/semantic.yml).**

#### Types

- `feat    `: A new feature
- `fix     `: A bug fix
- `build   `: Changes that affect the build system or external dependencies
- `chore   `: Other changes that don't modify source or test files
- `ci      `: Changes to our CI configuration files and scripts
- `docs    `: Documentation only changes
- `perf    `: A code change that improves performance
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `revert  `: Reverts a previous commit
- `style   `: Changes that do not affect the meaning of the code
- `test    `: Adding missing tests or correcting existing tests

#### Scopes

- api
- approval-service
- bridge
- cli
- configuration-service
- distributor
- docs
- helm-service
- installer
- jmeter-service
- lighthouse-service
- mongodb-datestore
- remediation-service
- secret-service
- shipyard-controller
- statistics-service
