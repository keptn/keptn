# Keptn Versioning

## Compatibility requirements

This codebase marks all artifacts with a version of the format `MAJOR.MINOR.PATCH` or
`MAJOR.MINOR.PATCH-next.ID`. For any artifact with a stable release, that is its version does not
end in `-next.ID`, no backwards-incompatible changes will be made unless incrementing the `MINOR`
version number. Patch releases with urgent cherry-picked bugfixes will be made by incrementing the
`PATCH` version.

A backwards-incompatible change affects the public API. The public API is any public reachable endpoint
that is not marked as `internal`. We refer to a backward-incompatible API change, a change that could
potentially cause failures in the applications that consume that API, e.g., removing an endpoint,
changing the response code, requiring additional parameters, or changing the entity models.

Such changes will be avoided - if they must be made, the affected endpoint will be marked as `deprecated`
and supported for one further `MINOR` version. After that, the affected API may be made backward-incompatible.

Backwards-incompatible changes to `internal` API are to be expected.

## Stable vs Next

Not all of our artifacts are published as stable artifacts - any non-stable artifact has the suffix
`-next.ID` on its version and they are meant for testing new features. NONE of the guarantees described
above apply to non-stable artifacts. They may require code or environment changes on every release and
are not meant for consumption for users where versioning stability is important. When a `next` artifact
is ready to be made stable, the release will be made by keeping the version, while the `-next.ID` suffix
will be removed.
