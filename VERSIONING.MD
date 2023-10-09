# Keptn Versioning

## Compatibility Requirements

A **backward-compatible change** is a modification that does not disrupt or break the functionality of existing applications or systems that rely on a specific public REST API. It allows for the codebase to evolve while ensuring that existing functionality remains intact. Such changes typically involve adding new features or options without affecting existing ones.

This codebase adheres to strict compatibility requirements, which are crucial for maintaining the stability of its artifacts. These requirements are based on a versioning scheme using the format `MAJOR.MINOR.PATCH` or `MAJOR.MINOR.PATCH-next.ID`. Here are the key principles:

**Backward-Incompatible Change**: A **backward-incompatible change** is a modification that can potentially disrupt or break the functionality of applications or systems relying on a specific public REST API. These changes affect the public-facing components of the codebase and may include actions such as:

- Removing an endpoint
- Altering the response code
- Requiring additional parameters
- Changing the structure of entity models within the API

**Compatibility Assurance**: To ensure compatibility and minimize disruptions for users of the codebase, the following rules apply:

1. For artifacts with stable releases (versions not ending in `-next.ID`), no backward-incompatible changes will be introduced unless there is an increment in the `MINOR` version number.

2. Patch releases, which include urgent cherry-picked bugfixes, are created by incrementing the `PATCH` version.

**Public REST API**: The term **public REST API** encompasses the following:

- Publicly reachable endpoints that are explicitly marked as `internal` are not considered part of the public API.
- The CloudEvents definition is also included in the public REST API definition.

**Handling Backward-Incompatible Changes**: In the event that a backward-incompatible change becomes necessary, the following steps are taken to mitigate its impact:

1. The affected endpoint is clearly marked as `deprecated`. This serves as a warning to users that the endpoint will undergo significant changes in the future.

2. The deprecated endpoint remains fully supported for one additional `MINOR` version of the codebase.

3. After this transitional period, the affected API component may undergo backward-incompatible modifications, as indicated by the versioning scheme.

**Expectation for Internal API**: It's important to note that **internal APIs**, which are not intended for public use, are inherently subject to backward-incompatible changes as part of the normal development process. Users should anticipate such changes when interacting with internal APIs.

## Stable vs Next

Not all of our artifacts are published as stable artifacts - any non-stable artifact has the suffix
`-next.ID` on its version and they are meant for testing new features. NONE of the guarantees described
above apply to non-stable artifacts. They may require code or environment changes on every release and
are not meant for consumption for users where versioning stability is important. When a `next` artifact
is ready to be made stable, the release will be made by keeping the version, while the `-next.ID` suffix
will be removed.
