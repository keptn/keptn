## Contributing

First off, thank you for considering contributing to Keptn. It's people like you that make Keptn great.

### Where do I go from here?

If you've noticed a bug, want to contribute features, or simply ask a question that for whatever reason you don't want to ask in the [Keptn Slack workspace](keptn.slack.com), please [search the issue tracker](https://github.com/keptn/keptn/issues?q=something) to see if someone else in the community has already created a ticket. If not, go ahead and [make one](https://github.com/keptn/keptn/issues/new).

### Read the docs

We are in the process of making sure that each repository and each service within the Keptn organization is documented properly. 
We are aware that some parts are currently missing, in the meantime please get in touch with us through the [Keptn Slack workspace](keptn.slack.com) if you have any questions.

As a starting point, please read the docs within the [docs/](docs/) folder in this repository.

### Fork and create a branch

We are following the [git branching model described in this blog post](https://nvie.com/posts/a-successful-git-branching-model/).

If you work on a new feature or on a bugfix, then fork the repository and branch away from the `develop` branch according to the following schema:

* `feature/###/name` for new features,
* `bug/###/name` for bugfixes, and
* `hotfix/###/name` for hotfixes.

where `###` is the github issue id. E.g., if you work on feature based on issue #123, the branch name could be

* `feature/123/foo`.

```bash
git checkout develop
git pull
git checkout -b feature/123/foo
```

For Hotfixes please branch away from the master branch, and create a PR to master and develop.

### Run Tests

Currently Keptn has no automated end-to-end test suite. To verify Keptn still works as it's supposed to, please make sure that the use cases described on the [Keptn website](https://keptn.sh/docs/) can be completed successfully.

### Make a Pull Request

At this point, you should switch back to the `develop` branch and make sure it's up to date with Keptn's `develop` branch:

```bash
git remote add upstream git@github.com:keptn/keptn.git
git checkout develop
git pull upstream develop
```

Then update your feature branch from your local copy of master, and push it.

```bash
git checkout feature/123/foo
git rebase develop
git push --set-upstream origin feature/123/foo
```

Finally, go to GitHub and make a Pull Request. Please describe what this PR is about and add a link to relevant GitHub issues.
Your PR will usually be reviewed automatically, but feel free to let us know about your PR via Slack.
