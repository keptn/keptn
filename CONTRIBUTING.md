## Contributing

First off, thank you for considering contributing to Keptn. It's people like you that make Keptn great.

### Where do I go from here?

If you've noticed a bug, want to contribute features, or simply ask a question that for whatever reason you don't want to ask in the [Keptn Slack workspace](keptn.slack.com), please [search the issue tracker](https://github.com/keptn/keptn/issues?q=something) to see if someone else in the community has already created a ticket. If not, go ahead and [make one](https://github.com/keptn/keptn/issues/new).

### Fork and create a branch

If you work on a new feature or a bug, then fork Keptn and create a branch with a descriptive name. A good branch name would be, e.g. where issue #123 is the ticket you're working on:

```
git checkout -b 123-feature-i-am-adding
```

### Run Tests

Currently Keptn has no automated end-to-end test suite. To verify Keptn still works as it's supposed to, please make sure that the use cases described on the [Keptn website](https://keptn.sh/docs/) can be completed successfully.

### Make a Pull Request

At this point, you should switch back to your master branch and make sure it's up to date with Keptn's master branch:

```
git remote add upstream git@github.com:keptn/keptn.git
git checkout master
git pull upstream master
```

Then update your feature branch from your local copy of master, and push it.

```
git checkout 123-feature-i-am-adding
git rebase master
git push --set-upstream origin 123-feature-i-am-adding
```

Finally, go to GitHub and make a Pull Request.
