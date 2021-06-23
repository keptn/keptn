# Release Notes 0.8.4-hotfix1

This is a hotfix release for Keptn 0.8.4, containing an updated Helm Chart *keptn-0.8.4-hotfix1.tgz* for setting the Bridge Look and Feel URL.

---

## Fixes

- Bridge LOOK_AND_FEEL_URL is missing in Keptn installer helm chart  [4476](https://github.com/keptn/keptn/issues/4476)

## Installation and upgrade instructions

Installing and/or upgrading to this release should be done using the `helm upgrade` command, e.g.:

**Installation** (using `--install`)
```console
helm upgrade keptn keptn --install -n keptn --create-namespace --wait --version=0.8.4-hotfix1 --repo=https://storage.googleapis.com/keptn-installer --set=control-plane.bridge.lookAndFeelUrl=https://example.com/bridge-look-and-feel.zip
```

**Upgrade** (using `--reuse-values`)
```console
helm upgrade keptn keptn -n keptn --wait --version=0.8.4-hotfix1 --repo=https://storage.googleapis.com/keptn-installer --reuse-values --set=control-plane.bridge.lookAndFeelUrl=https://example.com/bridge-look-and-feel.zip
```

It is not required to upgrade the CLI.
