# Working on a fork

For contributions, we expect you to work on a fork of the keptn/keptn repo. If you are unfamiliar with forks or Pull 
Request on GitHub, here is some information on the GitHub help pages:

* [Fork a repo](https://help.github.com/en/github/getting-started-with-github/fork-a-repo)
* [Configuring a remote for a fork](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/configuring-a-remote-for-a-fork)
* [Creating a pull request from a fork](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork)

*Warning*: If you are committing **any changes within your fork**, you **must not** commit them to the **master** branch. 
Changes on the master branch (even within your fork) can lead to losing your changes when syncing with the upstream repo
or to unresolvable merge conflicts.

## Forking the Keptn Repo

1. Fork the *keptn/keptn* repository on GitHub to your account, then clone the forked repo (**do not clone the original keptn/keptn repo**), e.g.:
    ```console
    git clone git@github.com:YOUR_USERNAME/keptn.git
    ```

2. Add an upstream to keptn/keptn (e.g., for syncing with the upstream afterwards):
    ```console
    git remote add upstream https://github.com/keptn/keptn.git
    ```
    and verify that there are two upstreams (on your local computer):
    ```console
    git remote -v
    ```
    Expected output:
    ```
    origin    https://github.com/YOUR_USERNAME/keptn.git (fetch)
    origin    https://github.com/YOUR_USERNAME/keptn.git (push)
    upstream  https://github.com/keptn/keptn.git (fetch)
    upstream  https://github.com/keptn/keptn.git (push)
    ```


## Keep your fork synced (merge master from keptn/keptn to your repos master branch)

Whenever you start working on a new feature, make sure that you branch away from the current master branch. 
The only exception of this is when you are working on a hotfix for a release, then branch away from one of the release branches.

1. Ensure that you have two upstreams: origin (your repo) and upstream (the keptn/keptn repo):
    ```console
    git remote -v
    ```
    Expected output:
    ```
    origin    https://github.com/YOUR_USERNAME/keptn.git (fetch)
    origin    https://github.com/YOUR_USERNAME/keptn.git (push)
    upstream  https://github.com/keptn/keptn.git (fetch)
    upstream  https://github.com/keptn/keptn.git (push)
    ```

2. Fetch changes from the upstream repo:
   ```console
   git fetch upstream
   ```
   
3. Check out your fork's local `master` (or `release-*`) branch:
   ```console
   git checkout master
   ```
   
4. Fetch all changes from the upstreams `master` (or `release-*`) branch (**this overwrites any changes in your local master branch!**):
   ```console
   git reset --hard upstream/master
   ```
 
5. Push those changes to your repository:
   ```console
   git push -f -u origin master
   ```
   
## Working on a branch
   
1. Ensure your forks master branch is up to date (see info above).

2. Create a new branch:
   ```console
   git checkout -b feature/1234/short_feature_description
   ```
   
3. Make changes to that branch

4. Push those changes to your repository:
   ```console
   git push -u origin feature/1234/short_feature_description
   ```

5. Create a pull request (see [Creating a pull request from a fork](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request-from-a-fork))
   using the GitHub Web UI


