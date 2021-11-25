# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [0.11.1](https://github.com/keptn/keptn/compare/0.11.0...0.11.1) (2021-11-25)


### Bug Fixes

* **bridge:** Fix problem with redirect and headers on cluster ([#6169](https://github.com/keptn/keptn/issues/6169)) ([8d17870](https://github.com/keptn/keptn/commit/8d17870d8f0ea8e88fe7d2ba7bd4441e004bf5f8))

## [0.11.0](https://github.com/keptn/keptn/compare/0.10.1...0.11.0) (2021-11-24)


### ⚠ BREAKING CHANGES

* MongoDB was updated from 3.6 to 4.4, also the custom helm chart was switched out for the Bitnami MongoDB Helm Chart. This means that a manual database migration is needed to preserve data during the keptn upgrade process! Steps to upgrade keptn with the manual migration can be found on the [Keptn Upgrade page](https://keptn.sh/docs/0.11.x/operate/upgrade/).

### Features

* Added context with cancel function to sdk ([#4552](https://github.com/keptn/keptn/issues/4552)) ([#5972](https://github.com/keptn/keptn/issues/5972)) ([d21e682](https://github.com/keptn/keptn/commit/d21e68230a789ac29dd47bf2b284928ded89a464))
* added probes for readiness and liveness ([#5303](https://github.com/keptn/keptn/issues/5303))  ([#5534](https://github.com/keptn/keptn/issues/5534)) ([6899ee7](https://github.com/keptn/keptn/commit/6899ee7dd6f23724c4c4dc6e16c4a218cbf2453c)), closes [#5533](https://github.com/keptn/keptn/issues/5533)
* **bridge:** Add 404 page ([#4983](https://github.com/keptn/keptn/issues/4983)) ([#6004](https://github.com/keptn/keptn/issues/6004)) ([aa7b4fa](https://github.com/keptn/keptn/commit/aa7b4fa0810da9c40f96996cb9a2fe34f8b66929))
* **bridge:** Add checkbox to set the `sendFinished` flag ([#5735](https://github.com/keptn/keptn/issues/5735)) ([#5989](https://github.com/keptn/keptn/issues/5989)) ([89598f8](https://github.com/keptn/keptn/commit/89598f8aaff5543171910c54b5dceffece4a1029))
* **bridge:** Prevent cut off of evaluation board ([#5279](https://github.com/keptn/keptn/issues/5279)) ([1ae06a1](https://github.com/keptn/keptn/commit/1ae06a13a53ab21df183b858544f3ced115897ba))
* **cli:** created user warning about changed database model in keptn 0.11.* ([#6071](https://github.com/keptn/keptn/issues/6071)) ([7f3447c](https://github.com/keptn/keptn/commit/7f3447c8eab6476ba5b484bae9f1344457f5ed02))
* handle time consistently ([#4788](https://github.com/keptn/keptn/issues/4788)) ([#5971](https://github.com/keptn/keptn/issues/5971)) ([e284d72](https://github.com/keptn/keptn/commit/e284d72268f8fa5b68faf00e23896dd731205b73))
* **lighthouse-service:** Added SIGTERM for lighthouse handlers ([#5304](https://github.com/keptn/keptn/issues/5304)) ([#5558](https://github.com/keptn/keptn/issues/5558)) ([ca9742c](https://github.com/keptn/keptn/commit/ca9742cc6a841bcb59ade85789e2155fbc7d8693))
* Switch mongoDB image to bitnami mongoDB chart ([#4801](https://github.com/keptn/keptn/issues/4801)) ([b3dabd6](https://github.com/keptn/keptn/commit/b3dabd6297bd0ddd0b6e2e5815c53919892045c2))


### Bug Fixes

* Adapt log level of SDK logs ([#5920](https://github.com/keptn/keptn/issues/5920)) ([#5921](https://github.com/keptn/keptn/issues/5921)) ([d314008](https://github.com/keptn/keptn/commit/d314008f3026c90252533fdf1c8ebe46538f9e42))
* **api:** Remove multiple types in event model ([#5948](https://github.com/keptn/keptn/issues/5948)) ([#5957](https://github.com/keptn/keptn/issues/5957)) ([30d5556](https://github.com/keptn/keptn/commit/30d5556e2da731a3d517c14ba36895c4bbaff11a))
* **approval-service:** Fall back to manual strategy when there is no result available ([#6012](https://github.com/keptn/keptn/issues/6012)) ([#6017](https://github.com/keptn/keptn/issues/6017)) ([9617814](https://github.com/keptn/keptn/commit/9617814ecbd5dcc853861d5b9939ab04c3d4a772))
* **bridge:** Add empty state to sequence-view ([#5084](https://github.com/keptn/keptn/issues/5084)) ([#5693](https://github.com/keptn/keptn/issues/5693)) ([b7c10df](https://github.com/keptn/keptn/commit/b7c10dfe0bccce073ef62f81806e0edf1166465a))
* **bridge:** Correctly show warning state ([#6003](https://github.com/keptn/keptn/issues/6003)) ([9a21d19](https://github.com/keptn/keptn/commit/9a21d19a698127801c8dbf9ca35d6e0f0e0531f3))
* **bridge:** don't log err (contains the x-token), only log err.message ([#6047](https://github.com/keptn/keptn/issues/6047)) ([#6052](https://github.com/keptn/keptn/issues/6052)) ([3eea6e3](https://github.com/keptn/keptn/commit/3eea6e37e8682d0a580313b0462cd02447925fcc))
* **bridge:** Fix integration curl commands for api ([#5941](https://github.com/keptn/keptn/issues/5941)) ([d76eccc](https://github.com/keptn/keptn/commit/d76ecccbcacc2b6db3ccbb06aec50942303d2bf7))
* **bridge:** Fixed missing problem title and decode of remediation config ([#6053](https://github.com/keptn/keptn/issues/6053)) ([ea0c53f](https://github.com/keptn/keptn/commit/ea0c53f5b429327c3780ffc6fd2b97522d25b44f))
* **bridge:** Fixed overwriting of data in environment screen ([#5841](https://github.com/keptn/keptn/issues/5841)) ([74a9a3d](https://github.com/keptn/keptn/commit/74a9a3ded3184ebe057fd918c7c1b23fd25d86c2))
* **bridge:** Fixed wrong weight of SLI ([#5987](https://github.com/keptn/keptn/issues/5987)) ([e536dbc](https://github.com/keptn/keptn/commit/e536dbc9bd8602deab4a1037a0f8bd775b7b16b7))
* **bridge:** Possible fix for flaky clicks in UI tests ([#5909](https://github.com/keptn/keptn/issues/5909)) ([58c5deb](https://github.com/keptn/keptn/commit/58c5deb78b5467d3e2b304447fcbed8b0589f984))
* **bridge:** Remove inline script for base url and upgrade-insecure-requests header part ([#6019](https://github.com/keptn/keptn/issues/6019)) ([b2e9960](https://github.com/keptn/keptn/commit/b2e9960a42e2b14f49a888619395855ede3a7c44))
* **bridge:** Show right event type ([#5828](https://github.com/keptn/keptn/issues/5828)) ([316d117](https://github.com/keptn/keptn/commit/316d117aa8c4c6d9f3f66d746be702bda6c53f02))
* **bridge:** Take SLI-weight out of SLO-file ([#5782](https://github.com/keptn/keptn/issues/5782)) ([f961ce1](https://github.com/keptn/keptn/commit/f961ce1fc126130047c0f5971ef7683c5e49a50f))
* **bridge:** Use helmet middlewares to prevent XSS ([8a58fb3](https://github.com/keptn/keptn/commit/8a58fb30389f4433dc049e7da37c03814a712d4b))
* **cli:** Make sure the release version is set in command descriptions ([#5762](https://github.com/keptn/keptn/issues/5762)) ([#5888](https://github.com/keptn/keptn/issues/5888)) ([24110c0](https://github.com/keptn/keptn/commit/24110c067ea8a698b5651d9b777a620765009c53))
* **cli:** problem with missing http(s) in endpoint flag during keptn auth ([#6039](https://github.com/keptn/keptn/issues/6039)) ([e4164db](https://github.com/keptn/keptn/commit/e4164db156cf0d767032d37cc59bd57091b6be10))
* **configuration-service:** changed bad order of extracting and adding resources to services ([#6006](https://github.com/keptn/keptn/issues/6006)) ([35605b7](https://github.com/keptn/keptn/commit/35605b7618c23271eadde447f7b8f12119f8db8b))
* **configuration-service:** Completely replace previous helm chart directory when updating ([#6050](https://github.com/keptn/keptn/issues/6050)) ([#6058](https://github.com/keptn/keptn/issues/6058)) ([74eefdf](https://github.com/keptn/keptn/commit/74eefdf5c12716c466f6ff64f6de4451a41ac650))
* **configuration-service:** Fix order of extracting and adding files to the repo ([#6041](https://github.com/keptn/keptn/issues/6041)) ([#6045](https://github.com/keptn/keptn/issues/6045)) ([4a3bf22](https://github.com/keptn/keptn/commit/4a3bf22bd560f0492e98ba19ea427162b9b549df))
* **distributor:** Fix message filtering in distributor ([#6074](https://github.com/keptn/keptn/issues/6074)) ([#6075](https://github.com/keptn/keptn/issues/6075)) ([602eb37](https://github.com/keptn/keptn/commit/602eb37609f65bc418db165a0679e2c0ab3edd78))
* **distributor:** fix subscription handling after message broker reconnect ([#5823](https://github.com/keptn/keptn/issues/5823)) ([49b1051](https://github.com/keptn/keptn/commit/49b1051fa31f8860fd490d4d5e1c20ffb6465048))
* **distributor:** Sanitized logs and cleaned up forwarder lifecycle ([#6036](https://github.com/keptn/keptn/issues/6036)) ([be5adb5](https://github.com/keptn/keptn/commit/be5adb5b77d10a15a8fa13afa5674460b52bd929))
* **distributor:** Set default timeout of Uniform API requests to 5s ([#6011](https://github.com/keptn/keptn/issues/6011)) ([#6015](https://github.com/keptn/keptn/issues/6015)) ([d89cab9](https://github.com/keptn/keptn/commit/d89cab9b11ebd83cc73f9c63d94ff23ae170818e))
* Fix bug where approval and remediation service would not run through unit tests anymore ([495654c](https://github.com/keptn/keptn/commit/495654cd0c59677527bd22e9764eecafc8b38c5f))
* Fix bug where DCO check always fails on dependabot PRs ([6a4b58d](https://github.com/keptn/keptn/commit/6a4b58d29a5da25206d021780aea4d8d6b1a762d))
* Fix multiple issues found by Sonatype Lift static analysis ([#5934](https://github.com/keptn/keptn/issues/5934)) ([dd93b4e](https://github.com/keptn/keptn/commit/dd93b4ea1f261435e5b995f10c5005aab2be06ca))
* Fix sub-project change detection for build-everything and master builds ([db808d6](https://github.com/keptn/keptn/commit/db808d68ee1cc02a104674078ad3b09d9ac5cc12))
* Fix version not showing up anymore in API ([#5783](https://github.com/keptn/keptn/issues/5783)) ([1eea3f9](https://github.com/keptn/keptn/commit/1eea3f94e62b7ded029c54f36673bb577fee1977))
* Fixed bug where MongoDB would not come up in airgapped setup ([#5939](https://github.com/keptn/keptn/issues/5939)) ([079a6b4](https://github.com/keptn/keptn/commit/079a6b43c33c80f1739a8c2d8985b4c73337f517))
* Handle upstream not found ([#5977](https://github.com/keptn/keptn/issues/5977)) ([#5994](https://github.com/keptn/keptn/issues/5994)) ([77240d4](https://github.com/keptn/keptn/commit/77240d4fa9e3f3f24629319f007051c1cc387a0c))
* **shipyard-controller:** cleanup uniform subscriptions when service is deleted ([#5725](https://github.com/keptn/keptn/issues/5725)) ([#5766](https://github.com/keptn/keptn/issues/5766)) ([d95f7a6](https://github.com/keptn/keptn/commit/d95f7a665f9ffaa8929fa65d84a0c634591c0a9c))
* **shipyard-controller:** migrate and avoid mongodb keys containing dots ([#6065](https://github.com/keptn/keptn/issues/6065)) ([5259bcf](https://github.com/keptn/keptn/commit/5259bcfd534690c596383f4eae93576e292f8b03))
* **shipyard-controller:** removed error shadowing ([#6048](https://github.com/keptn/keptn/issues/6048)) ([04416da](https://github.com/keptn/keptn/commit/04416da01497d561300583770d01914ec3af20c8))
* **shipyard-controller:** Store `lastEventTypes` only for events that belong to a sequence controlled by the shipyard controller ([#5309](https://github.com/keptn/keptn/issues/5309)) ([#5777](https://github.com/keptn/keptn/issues/5777)) ([ee27c62](https://github.com/keptn/keptn/commit/ee27c62f07e04791892c90b39ade410cc345034e))
* Update auto-update pipelines to follow keptns semantic PR guidelines ([#5931](https://github.com/keptn/keptn/issues/5931)) ([280fa4e](https://github.com/keptn/keptn/commit/280fa4ef4829e4d9a3968b7049befe7d3fd87304))
* **webhook-service:** Avoid .finished.finished/.started.finished events ([#5954](https://github.com/keptn/keptn/issues/5954)) ([#6000](https://github.com/keptn/keptn/issues/6000)) ([fbe01a8](https://github.com/keptn/keptn/commit/fbe01a88a642904eaecf507d22f9615a560f124f))
* **webhook-service:** invalid conversion of eventType ([#5998](https://github.com/keptn/keptn/issues/5998)) ([67dba55](https://github.com/keptn/keptn/commit/67dba55d8fed3e978ce596a02073ff059581e165))


### Docs

* Add release notes for 0.10.0 ([d748251](https://github.com/keptn/keptn/commit/d7482517e74b8b12038c9438daa01f3172ad81e3))


### Other

* Add @RealAnna to Maintainers list ([34175bb](https://github.com/keptn/keptn/commit/34175bb5aa9cb03b1ba2600ae1d93e5e8602d13d))
* Add environment variables for setting log levels of Keptn services ([#5373](https://github.com/keptn/keptn/issues/5373)) ([#5911](https://github.com/keptn/keptn/issues/5911)) ([809baea](https://github.com/keptn/keptn/commit/809baea2672fdbcebe236a9b4a6760223cb84870))
* Add flowcharts that describe components of the shipyard controller ([#5919](https://github.com/keptn/keptn/issues/5919)) ([8aa4dd8](https://github.com/keptn/keptn/commit/8aa4dd85fc24d09fa362812518d1314ec4a41c79))
* add missing release notes ([#5781](https://github.com/keptn/keptn/issues/5781)) ([dab9844](https://github.com/keptn/keptn/commit/dab9844c7213638663faf9b37879432e85d4f312))
* Add odubajDT as maintainer ([#60](https://github.com/keptn/keptn/issues/60)) ([#6049](https://github.com/keptn/keptn/issues/6049)) ([65ae6cf](https://github.com/keptn/keptn/commit/65ae6cfa505dc07a5f7207b7625b36dfac8549cc))
* Add TannerGilbert as project member ([#5899](https://github.com/keptn/keptn/issues/5899)) ([65148be](https://github.com/keptn/keptn/commit/65148be195779baeed0a7462cfc3d753ec88c861))
* Add the correct label for bug reports ([#5908](https://github.com/keptn/keptn/issues/5908)) ([dc296a5](https://github.com/keptn/keptn/commit/dc296a5a504a453925c5a761690554a9b2fd896b))
* Added go-sdk and webhook-service to dependencies-and-licenses check ([#5898](https://github.com/keptn/keptn/issues/5898)) ([6481ca3](https://github.com/keptn/keptn/commit/6481ca30fdf8068d3626c32de8b967606023a0e2))
* Cancel integration tests when mismatch between CLI and kube context is detected ([#5743](https://github.com/keptn/keptn/issues/5743)) ([#5824](https://github.com/keptn/keptn/issues/5824)) ([5596611](https://github.com/keptn/keptn/commit/559661149928ed7cfc19ea1e99da60e2fb1322a9)), closes [#5734](https://github.com/keptn/keptn/issues/5734)
* fixing imports according to snyc ([#5936](https://github.com/keptn/keptn/issues/5936)) ([391ace2](https://github.com/keptn/keptn/commit/391ace23fa8df83be6415d61c65ba9b5de66f637))
* **helm-service:** More meaningful error messages ([#6089](https://github.com/keptn/keptn/issues/6089)) ([80d59cb](https://github.com/keptn/keptn/commit/80d59cbb68633659cd26ded62a18c12c92f84155))
* Increase timeout of DeliveryAssistant integration test ([#6067](https://github.com/keptn/keptn/issues/6067)) ([b141ce4](https://github.com/keptn/keptn/commit/b141ce435491b7e01a3c1f499da90a125a1c27af))
* **jmeter-service:** bump version of jmeter binary to 5.4.1 ([#6032](https://github.com/keptn/keptn/issues/6032)) ([3c250d2](https://github.com/keptn/keptn/commit/3c250d206c6694202b3116bdf28207778ecaa67b))
* **jmeter-service:** cleanups ([#6014](https://github.com/keptn/keptn/issues/6014)) ([5e779eb](https://github.com/keptn/keptn/commit/5e779eb91ac0cbdc8d64729bfcbe9c8d739b1b31))
* Mitigating racecondition in unit tests ([#5901](https://github.com/keptn/keptn/issues/5901)) ([5a642a5](https://github.com/keptn/keptn/commit/5a642a538aee5d4a4e162de7dbb722216b5794a4))
* **mongodb-datastore:** Refactoring ([#5917](https://github.com/keptn/keptn/issues/5917)) ([#6002](https://github.com/keptn/keptn/issues/6002)) ([3242094](https://github.com/keptn/keptn/commit/324209413043e79561198075b0da71eb42de063c))
* Polish HTTP(S) headers ([a4f52b4](https://github.com/keptn/keptn/commit/a4f52b409aa77a64e89eba637e71b5cdeefb22e1))
* Remove sequence migration integration test because component has been removed ([#6101](https://github.com/keptn/keptn/issues/6101)) ([afeb7fc](https://github.com/keptn/keptn/commit/afeb7fc718ecad0babf6cdb4d7e8517e5c3cf721))
* removed cluster role binding ([#5955](https://github.com/keptn/keptn/issues/5955)) ([391a3ba](https://github.com/keptn/keptn/commit/391a3bae209254dd9c104c2e4782f89804f14900))
* Removed obsolete files ([#4818](https://github.com/keptn/keptn/issues/4818)) ([#5932](https://github.com/keptn/keptn/issues/5932)) ([588a76d](https://github.com/keptn/keptn/commit/588a76ded965af36f25790e22a5e6a82d11c66a9))
* **shipyard-controller:** Adapted log level ([#5978](https://github.com/keptn/keptn/issues/5978)) ([3cbfcd7](https://github.com/keptn/keptn/commit/3cbfcd7bb5eceea810095964f7892eb0886a18c3))
* **shipyard-controller:** cleaning up package(s) ([#5786](https://github.com/keptn/keptn/issues/5786)) ([a6e51d4](https://github.com/keptn/keptn/commit/a6e51d42f9cfe633437900a23d5d2d7d0f6ae317))
* **shipyard-controller:** cleanups & refactorings 2 ([#5937](https://github.com/keptn/keptn/issues/5937)) ([adf4078](https://github.com/keptn/keptn/commit/adf4078fc277076ea3e057b2f4448cdc2b44243e))
* **shipyard-controller:** Do not interpret absence of configurationChange property as an error ([#5979](https://github.com/keptn/keptn/issues/5979)) ([#5982](https://github.com/keptn/keptn/issues/5982)) ([28a9a92](https://github.com/keptn/keptn/commit/28a9a92dfdda33d4952d4322c14c2363401a38d3))
* **shipyard-controller:** Extract shipyard retrieval into its own component ([#5243](https://github.com/keptn/keptn/issues/5243)) ([#5821](https://github.com/keptn/keptn/issues/5821)) ([a1d18ae](https://github.com/keptn/keptn/commit/a1d18ae4c277bfaebe0c19f408a764033057f1f2))
* **shipyard-controller:** move event operations to event repo ([#5902](https://github.com/keptn/keptn/issues/5902)) ([730864b](https://github.com/keptn/keptn/commit/730864b9e57c8d4a5cfb82567a575eb8e9661cb7))
* Updated dependencies according to ArtifactHub and Snyk ([#5543](https://github.com/keptn/keptn/issues/5543)) ([#5951](https://github.com/keptn/keptn/issues/5951)) ([48fc51c](https://github.com/keptn/keptn/commit/48fc51c1c2b4689b6658e36a8751d0d020c1912c))
* Updated go-utils dependency ([#5968](https://github.com/keptn/keptn/issues/5968)) ([#5969](https://github.com/keptn/keptn/issues/5969)) ([f2c796e](https://github.com/keptn/keptn/commit/f2c796e19b8cac6ae737b2a85e7868b763816e0b))
* use correct link in CLI upgrade message ([961ea2a](https://github.com/keptn/keptn/commit/961ea2a6f9572d3289abb61b2827289e5e2ac224))
* Version 0.10.0 into master ([9eb12ec](https://github.com/keptn/keptn/commit/9eb12ec3f0119e9e0ea30050019e3acf4cca535a))


### Refactoring

* **bridge:** Reduce number of API calls for project endpoint ([#5450](https://github.com/keptn/keptn/issues/5450)) ([25fd876](https://github.com/keptn/keptn/commit/25fd8766b781797b0e043dc929ec92841745982c))
* **bridge:** Refactoring of project settings / create project ([#5100](https://github.com/keptn/keptn/issues/5100)) ([03fc3d2](https://github.com/keptn/keptn/commit/03fc3d298592423cef72a12aaac2a1a8001c2a66))
* **bridge:** Refactoring of service screen ([#4918](https://github.com/keptn/keptn/issues/4918)) ([#5244](https://github.com/keptn/keptn/issues/5244)) ([8f3b810](https://github.com/keptn/keptn/commit/8f3b810e60a6ca7e3cfb02adffb62375cabc9a27))
* **bridge:** Refactoring of services settings ([#5100](https://github.com/keptn/keptn/issues/5100)) ([771ec59](https://github.com/keptn/keptn/commit/771ec595b482d24c9e0d2473d2eabf303e139bc2))
* **cli:** use viper to manage keptn config ([#5694](https://github.com/keptn/keptn/issues/5694)) ([498d893](https://github.com/keptn/keptn/commit/498d89345e42605fe24aea731e66b03dea722be7))
