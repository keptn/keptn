# Release Notes 0.8.6

This is a hotfix release for Keptn 0.8.5, containing an updated Helm Chart *keptn-0.8.6.tgz* for setting the node/namespace of a distributor.

---

## Fixes

- Allow to set the metadata for the node and namespace [4591](https://github.com/keptn/keptn/issues/4591)
- Correct syntax [4609](https://github.com/keptn/keptn/issues/4609)

## Installation and upgrade instructions

Installing and/or upgrading to this release should be done using the `helm upgrade` command, e.g.:

**Installation** (using `--install`)
```console
helm upgrade keptn keptn --install -n keptn --create-namespace --wait --version=0.8.6 --repo=https://storage.googleapis.com/keptn-installer
```

**Upgrade** (using `--reuse-values`)
```console
helm upgrade keptn keptn -n keptn --wait --version=0.8.6 --repo=https://storage.googleapis.com/keptn-installer --reuse-values
```

It is not required to upgrade the CLI.