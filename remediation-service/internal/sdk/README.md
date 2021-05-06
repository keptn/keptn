# Experimental Keptn Service SDK

**TODOs**:
* implement "early-return", i.e. the processing of the business logic shall be done in the
  background, while the actual HTTP response for the triggering cloud event shall return early.
  This. should be the default behavior. However, for being able to write more deterministic
  unit tests, also the "sequential" flavor of processing events shall be possible to enable
* Move code into own go module
* Support for streamlined database connections
* Generator for creating project structure
* Documentation on how to use the SDK
* ...