# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [0.18.0](https://github.com/keptn/keptn/compare/0.17.0...0.18.0) (2022-07-28)


### ⚠ BREAKING CHANGES

* **resource-service:** Trailing `/` chars in the resource APIs will return a 404. This way, the difference between an empty URI and getting all the resources is explicit.
* All Keptn core now depends on resource-service. From this moment on resource-service installation is mandatory.

### Features

* **api:** Add create-secret api action to import endpoint ([#8348](https://github.com/keptn/keptn/issues/8348)) ([df9c42b](https://github.com/keptn/keptn/commit/df9c42b2463f3ded77a1f0eb742ddd0e68c7d7d3))
* **api:** Implement create webhook subscription action ([73133f0](https://github.com/keptn/keptn/commit/73133f04eb933d2117a18704c9b188cabc59dc08))
* **api:** Import package manifest templating ([96035b9](https://github.com/keptn/keptn/commit/96035b9e0218e1678d8f31d1bb52c6a2044ae057))
* **api:** Process import package manifest and execute API tasks ([74744aa](https://github.com/keptn/keptn/commit/74744aa7d6fc60f187d7065ba135a9ad36da8892))
* **api:** Support simple templating for resource and api tasks in import manifests ([#8456](https://github.com/keptn/keptn/issues/8456)) ([02fd6d5](https://github.com/keptn/keptn/commit/02fd6d520946edbce8890df6b2367f533309b1a6))
* **api:** Upload resources from import package ([67339ea](https://github.com/keptn/keptn/commit/67339ea85f12978a820295b5d7bd14f177893ae7))
* **bridge:** Add ktb-chart ([#8420](https://github.com/keptn/keptn/issues/8420)) ([9d55c35](https://github.com/keptn/keptn/commit/9d55c352ebc3394c00d70bd63a0c672462a8be6b))
* **bridge:** Modularize dashboard view and introduce lazy loading ([#8315](https://github.com/keptn/keptn/issues/8315)) ([a6326ca](https://github.com/keptn/keptn/commit/a6326ca66b70e63802b66d6f4a6ce8915d05bc1c))
* **bridge:** Modularize environment view and introduce lazy loading ([#8313](https://github.com/keptn/keptn/issues/8313)) ([4c1ad1a](https://github.com/keptn/keptn/commit/4c1ad1a6c9f188a8cc3d79de941fcd92a580c9f7))
* **bridge:** Modularize evaluation-board and introduce lazy loading ([#8340](https://github.com/keptn/keptn/issues/8340)) ([60309c5](https://github.com/keptn/keptn/commit/60309c5867748c1281b75951dd090c37055ad6c7))
* **bridge:** Modularize integration and common use case views and introduce lazy loading ([#8305](https://github.com/keptn/keptn/issues/8305)) ([609602a](https://github.com/keptn/keptn/commit/609602a636adb19865c4b18b3b55a472f00707a4))
* **bridge:** Modularize project board and introduce lazy loading ([#8342](https://github.com/keptn/keptn/issues/8342)) ([63d61fb](https://github.com/keptn/keptn/commit/63d61fb35e084e4e31a72b1c1440582cc2188430))
* **bridge:** Modularize sequence and logout view and introduce lazy loading ([#8289](https://github.com/keptn/keptn/issues/8289)) ([6cc2e2c](https://github.com/keptn/keptn/commit/6cc2e2c2c88cff5a02a80ecaa6202ff2101e1098))
* **bridge:** Modularize services view and introduce lazy loading ([#8325](https://github.com/keptn/keptn/issues/8325)) ([e1f18d4](https://github.com/keptn/keptn/commit/e1f18d4eafaf4055f1b0bc180e5bb9d545ced27c))
* **bridge:** Modularize settings view ([#8397](https://github.com/keptn/keptn/issues/8397)) ([4373f21](https://github.com/keptn/keptn/commit/4373f21edc4d8c2382c21ddefc4566643bf16e5e))
* **bridge:** Preselect date for datetime picker ([#8450](https://github.com/keptn/keptn/issues/8450)) ([2817781](https://github.com/keptn/keptn/commit/2817781cb7bf885f15c02cc8c0c0a63d1326e880))
* **bridge:** Select project stage from project overview ([#7736](https://github.com/keptn/keptn/issues/7736)) ([e05415c](https://github.com/keptn/keptn/commit/e05415c64e40444f42aec1beacb482170052d2f3))
* **bridge:** Show pause icon if sequence is paused ([#8471](https://github.com/keptn/keptn/issues/8471)) ([6b2669b](https://github.com/keptn/keptn/commit/6b2669bf81097477ad431aeb2589b92e42321337))
* **bridge:** Show user info for OAuth "Insufficient permission" error message ([#8403](https://github.com/keptn/keptn/issues/8403)) ([b2afdf9](https://github.com/keptn/keptn/commit/b2afdf9d9ecb687fca52a7852c15ecf90741bea8))
* **cli:** Introduce webhookConfig migrator ([#8396](https://github.com/keptn/keptn/issues/8396)) ([917e056](https://github.com/keptn/keptn/commit/917e056f28228de99cd3e927a562771acd194296))
* **cli:** Removed install/uninstall/upgrade commands ([#8302](https://github.com/keptn/keptn/issues/8302)) ([bb8015c](https://github.com/keptn/keptn/commit/bb8015ca9bbe2fdca847b09458d36dd587ae4524))
* **installer:** Add options for setting image repository and tag globally ([#8152](https://github.com/keptn/keptn/issues/8152)) ([100eae9](https://github.com/keptn/keptn/commit/100eae9680b4e96b6404c68bb3a66c4a33d2f9da))
* **installer:** Enable clustered NATS ([#8464](https://github.com/keptn/keptn/issues/8464)) ([3c1ae2b](https://github.com/keptn/keptn/commit/3c1ae2b0a2c3cf4130aab80e76265bc1fc9f6431))
* **installer:** Introduce flags to enable / disable Keptn services ([#8316](https://github.com/keptn/keptn/issues/8316)) ([6ccc7b1](https://github.com/keptn/keptn/commit/6ccc7b1fcead29f5c67b3a81f79b35e1ed9b29cd))
* **installer:** More Security Improvements for NATS ([#8421](https://github.com/keptn/keptn/issues/8421)) ([42e9fad](https://github.com/keptn/keptn/commit/42e9fadefe5b70fed8e07623713ff590c25e723e))
* **installer:** Remove configuration-service references and resourceService.enabled option ([#8296](https://github.com/keptn/keptn/issues/8296)) ([8d8eb99](https://github.com/keptn/keptn/commit/8d8eb99f14c035e86007dce9241ca820e7068921))
* **installer:** Security Improvements ([#8373](https://github.com/keptn/keptn/issues/8373)) ([d946f67](https://github.com/keptn/keptn/commit/d946f67af3827cad0c46beabbffaf39a669d98bc))
* **shipyard-controller:** Introduce API Endpoint for retrieving Sequence Executions ([#8430](https://github.com/keptn/keptn/issues/8430)) ([ac326c7](https://github.com/keptn/keptn/commit/ac326c780b2886a35588b8d1c1869e418ba566b8))
* **shipyard-controller:** Introduce RemoteURL denyList ([#8490](https://github.com/keptn/keptn/issues/8490)) ([6db8f3d](https://github.com/keptn/keptn/commit/6db8f3db00944eb80c94fac92c1d62e9ae9dc551))


### Bug Fixes

* **bridge:** Added missing wait to view more services ui test ([#8320](https://github.com/keptn/keptn/issues/8320)) ([f2bce6b](https://github.com/keptn/keptn/commit/f2bce6b1665fd67b107193890da0794241796ea5))
* **bridge:** check if configurationChange has image ([#8507](https://github.com/keptn/keptn/issues/8507)) ([16ec462](https://github.com/keptn/keptn/commit/16ec46205361c8b130a6c947e8a20b062d1e75e0))
* **bridge:** Evaluation info misleading if failed because of key sli ([#8250](https://github.com/keptn/keptn/issues/8250)) ([a5d79d0](https://github.com/keptn/keptn/commit/a5d79d03b5d526033d73395a3fc708846b1c21c3))
* **bridge:** Fix detection of pending changes when automatic provisioning active ([#8531](https://github.com/keptn/keptn/issues/8531)) ([0d4c7d2](https://github.com/keptn/keptn/commit/0d4c7d28329cafcf69f3063ec5d72e47638b2505))
* **bridge:** Fix error on viewing service deployment ([#8332](https://github.com/keptn/keptn/issues/8332)) ([9e9f776](https://github.com/keptn/keptn/commit/9e9f7769daaea303310ad7a76ab8dc880de6e6fa))
* **bridge:** Fix evaluation badge wrapping ([#8524](https://github.com/keptn/keptn/issues/8524)) ([d8f75ea](https://github.com/keptn/keptn/commit/d8f75ea2ca217a0f90dbb45bec498039e8b4f1fc))
* **bridge:** Fix incorrect API URL for auth command ([#8386](https://github.com/keptn/keptn/issues/8386)) ([9ea6132](https://github.com/keptn/keptn/commit/9ea613270e159faae3c9059c42decab50b4cae14))
* **bridge:** Navigating to service from stage-details ([#8399](https://github.com/keptn/keptn/issues/8399)) ([e0ce5bd](https://github.com/keptn/keptn/commit/e0ce5bde8d73437aa1bd1bd92a1c435316ce58a9))
* **cli:** Fix broken xref in CLI command reference docs ([#8374](https://github.com/keptn/keptn/issues/8374)) ([cb92bf5](https://github.com/keptn/keptn/commit/cb92bf530757ebac76a7382de8948b01fb51934a))
* **cli:** Print ID of Keptn context after sending events ([#8392](https://github.com/keptn/keptn/issues/8392)) ([65ce578](https://github.com/keptn/keptn/commit/65ce57807ed557ef97067ee9ef075690a2b515c8))
* **installer:** Disable nats cluster due to unreliable behavior ([#8523](https://github.com/keptn/keptn/issues/8523)) ([36cdb07](https://github.com/keptn/keptn/commit/36cdb0735d0b9338664d5c1768e5a9c845276ff7))
* **installer:** Fix NATS clustering settings ([#8484](https://github.com/keptn/keptn/issues/8484)) ([af15cbe](https://github.com/keptn/keptn/commit/af15cbe6902ac4730f8404618dcc97e3677c2ba1))
* **installer:** Fix Nginx not starting when statistics service is disabled ([#8326](https://github.com/keptn/keptn/issues/8326)) ([cde5942](https://github.com/keptn/keptn/commit/cde5942a070f54ba0c373a715acc847649371f02))
* **installer:** Remove configuration service from airgapped installer scripts ([#8376](https://github.com/keptn/keptn/issues/8376)) ([772ebd6](https://github.com/keptn/keptn/commit/772ebd612f9719eec5944c07e718b12e804aae07))
* **installer:** RoleBinding is not installed if not needed for shippy leader election ([#8535](https://github.com/keptn/keptn/issues/8535)) ([e90e94b](https://github.com/keptn/keptn/commit/e90e94bafd456879ff4389e5ae4cd2e9a0ba06df))
* **resource-service:** Return 404 with trailing slashes ([#8265](https://github.com/keptn/keptn/issues/8265)) ([785a39c](https://github.com/keptn/keptn/commit/785a39c60a2056d37fc7c02ff52151a375916e95))
* **resource-service:** Unescape resourceURI before updating single resource ([#8441](https://github.com/keptn/keptn/issues/8441)) ([a73af9e](https://github.com/keptn/keptn/commit/a73af9e2b3db856fdd47777a08ed27e06f85283c))
* **shipyard-controller:** Handling error messages ([#8480](https://github.com/keptn/keptn/issues/8480)) ([dbcb214](https://github.com/keptn/keptn/commit/dbcb214a7fd8efbd191161659739102c1f6dd8ad))
* **webhook-service:** Do not respond to anything else than .triggered event on pre execution error ([#8337](https://github.com/keptn/keptn/issues/8337)) ([4430a13](https://github.com/keptn/keptn/commit/4430a13184e5d4a98eda089e642e587434323492))
* **webhook-service:** Typo in component tests ([#8409](https://github.com/keptn/keptn/issues/8409)) ([7d77b7d](https://github.com/keptn/keptn/commit/7d77b7dbea3f24e642c1146f75902c0e44db687a))
* Zero Downtime test for the webhook-service ([#8408](https://github.com/keptn/keptn/issues/8408)) ([9212fb2](https://github.com/keptn/keptn/commit/9212fb29f07dd877e58b6c9db9e55a69354d453b))


### Docs

* **cli:** Fix typo in create secret command ([#8498](https://github.com/keptn/keptn/issues/8498)) ([36d373f](https://github.com/keptn/keptn/commit/36d373f0516bab7d07ee17208b282b380e484571))
* Fix instructions to install master ([#8429](https://github.com/keptn/keptn/issues/8429)) ([ac943cc](https://github.com/keptn/keptn/commit/ac943cc14c8f05296177fe603c8563e3ff6505df))


### Refactoring

* **bridge:** Refactor project settings ([#8510](https://github.com/keptn/keptn/issues/8510)) ([f10880b](https://github.com/keptn/keptn/commit/f10880bb4a1318aa5f90946e05fdeb7dcc429c9d))
* **bridge:** Refactor secrets page ([#8300](https://github.com/keptn/keptn/issues/8300)) ([66b1dfc](https://github.com/keptn/keptn/commit/66b1dfcdf658c74afaa53c8e32ab343b563f1dd5))
* **bridge:** Refactor the project settings of services ([#8323](https://github.com/keptn/keptn/issues/8323)) ([7bb4122](https://github.com/keptn/keptn/commit/7bb412215d77d6f1f5c1d66adf55f74f26cb1e5f))
* **bridge:** Remove global project polling and remove project dependency in integrations view ([#8412](https://github.com/keptn/keptn/issues/8412)) ([c4845c9](https://github.com/keptn/keptn/commit/c4845c9b568351ab2c93320095d4f061584b7463))
* **bridge:** Rename D3 feature flag ([#8499](https://github.com/keptn/keptn/issues/8499)) ([6a389df](https://github.com/keptn/keptn/commit/6a389df6d44ff561c11d72dc727b49247b5ad9da))
* Refactor all services to use resource-service ([#8271](https://github.com/keptn/keptn/issues/8271)) ([f866d09](https://github.com/keptn/keptn/commit/f866d094a7eba76306e31c7e12d9bf80aa9ae046))


### Other

*  Added new component test in remediation service ([#8343](https://github.com/keptn/keptn/issues/8343)) ([a0c22f9](https://github.com/keptn/keptn/commit/a0c22f9c950176c54e2e25d2183bff8e39b38181))
*  Fix dev repo registry in zd test ([#8411](https://github.com/keptn/keptn/issues/8411)) ([1d17283](https://github.com/keptn/keptn/commit/1d17283ab1db21e8d79b23402de9973c9aac0794))
* Add helm dependencies directly to repository charts ([#8472](https://github.com/keptn/keptn/issues/8472)) ([e6669a4](https://github.com/keptn/keptn/commit/e6669a49ce43ff6d65afcf622a4f78e8eb0dcf47))
* Added repo to resource-service.yaml ([#8382](https://github.com/keptn/keptn/issues/8382)) ([d70d82d](https://github.com/keptn/keptn/commit/d70d82d39fc42d6530bd23ee1726bad9b573a584))
* Added timeout to keptn install ([#8383](https://github.com/keptn/keptn/issues/8383)) ([e2837bb](https://github.com/keptn/keptn/commit/e2837bb62d8ca72ee641f5a006ae17cf950f2293))
* **bridge:** Enable resource-service by default ([#8432](https://github.com/keptn/keptn/issues/8432)) ([40d75d1](https://github.com/keptn/keptn/commit/40d75d1a0ee4a93ab0c44d67c133137ee3fa8479))
* **bridge:** Fix Sonar issues ([#8384](https://github.com/keptn/keptn/issues/8384)) ([b389f67](https://github.com/keptn/keptn/commit/b389f67e9c70b260e36fbd16811ee712aa1035f0))
* **bridge:** Fix Sonar issues part 2 ([#8398](https://github.com/keptn/keptn/issues/8398)) ([ce80143](https://github.com/keptn/keptn/commit/ce8014361af1eb7fc54027f026529e1dd01b2adb))
* **bridge:** Generalization of showing a running sequence ([#8379](https://github.com/keptn/keptn/issues/8379)) ([73e4634](https://github.com/keptn/keptn/commit/73e4634f33a9849e68c868a9240f669e808ecbfc))
* **bridge:** Remove loading of integrations on common-use-case page ([#8344](https://github.com/keptn/keptn/issues/8344)) ([77560f5](https://github.com/keptn/keptn/commit/77560f5a61e802fabb94e0b9d48a31f4db050386))
* **bridge:** Remove no Git upstream is set warning ([#8447](https://github.com/keptn/keptn/issues/8447)) ([ab35607](https://github.com/keptn/keptn/commit/ab356077803b8b9e67a54bc393bcd97128090b0b))
* **bridge:** Remove second labels tag list for remediation sequences ([#8410](https://github.com/keptn/keptn/issues/8410)) ([5bb977e](https://github.com/keptn/keptn/commit/5bb977e281f9b027ba1a0ce627534a2dd50041d0))
* **bridge:** Remove unused service page env var ([#8356](https://github.com/keptn/keptn/issues/8356)) ([7098fdb](https://github.com/keptn/keptn/commit/7098fdba69079e0eb1c36ca0db047e8bc8ae3152))
* **bridge:** Removed obsolete common use cases page ([#8419](https://github.com/keptn/keptn/issues/8419)) ([98e477b](https://github.com/keptn/keptn/commit/98e477b4c276a92033fef6d5f3867e174583f4ec))
* **cli:** Remove warning that no Git upstream is set ([#8518](https://github.com/keptn/keptn/issues/8518)) ([ff49bad](https://github.com/keptn/keptn/commit/ff49bad611f1ed165373b907d14699a78be79501))
* Fix ZeroDowntime registry ([#8434](https://github.com/keptn/keptn/issues/8434)) ([c89506d](https://github.com/keptn/keptn/commit/c89506d088f9b5244955dfa1a721f3a5b81f2a7a))
* Increased coverage for remediation-service ([#8357](https://github.com/keptn/keptn/issues/8357)) ([867d947](https://github.com/keptn/keptn/commit/867d947728ec3399885b54f7880dfd73e880fbd3))
* **installer:** Improve NATS configuration ([#8475](https://github.com/keptn/keptn/issues/8475)) ([0c8a964](https://github.com/keptn/keptn/commit/0c8a964f6bdc038aa0698c7ccc31c911da1cf9db))
* Remove go mod files of configuration service ([#8341](https://github.com/keptn/keptn/issues/8341)) ([1c74388](https://github.com/keptn/keptn/commit/1c743885a14c935d11ce6a015b12b793359c5bd5))
* Remove reference to removed test ([#8369](https://github.com/keptn/keptn/issues/8369)) ([03aec7b](https://github.com/keptn/keptn/commit/03aec7b3698540c1a0cdfc73304d0dd53a0db7b8))
* Removed configuration-service module ([#8294](https://github.com/keptn/keptn/issues/8294)) ([bd3c9af](https://github.com/keptn/keptn/commit/bd3c9af302af6b9523c77d738f6fecca51264136))
* Removed redundant integration tests ([#8324](https://github.com/keptn/keptn/issues/8324)) ([44764cd](https://github.com/keptn/keptn/commit/44764cd784479341bac0c0fac67e7c811bb06d5d))
* **shipyard-controller:** Add extra debug logging to the Git Automatic Provisioner call ([#8440](https://github.com/keptn/keptn/issues/8440)) ([cc9a212](https://github.com/keptn/keptn/commit/cc9a212970034a108314709c11b67dc62bbc72a5))
* **webhook-service:** Slimmed down integration tests for webhook service ([#8339](https://github.com/keptn/keptn/issues/8339)) ([7a01bd0](https://github.com/keptn/keptn/commit/7a01bd0fcbd410c47810adb4d0b5876a5b9e5fe6))

## [0.17.0](https://github.com/keptn/keptn/compare/0.16.0...0.17.0) (2022-07-06)


### ⚠ BREAKING CHANGES

* Git credentials for git authentication were moved to a separate sub-structure in go-utils package and split to either ssh or http sub-structures depending on the used authentication method. This leads to new models for creating, updating and retrieving the project information.
 * **installer:** Keptn's Helm charts were reworked and some values were changed and/or moved. Please consult the upgrade guide to make sure your installation can be upgraded successfully.
With this change, users now have the option to customise resource limits/requests and to add custom sidecars and extra volumes from the Helm values.

### Features

* Adapt go-utils changes in git credentials models ([#8020](https://github.com/keptn/keptn/issues/8020)) ([e8e2e6c](https://github.com/keptn/keptn/commit/e8e2e6cfe7af4b30f4071de6aef38ecdc12907c7))
* Add headers to git provisioner ([#8132](https://github.com/keptn/keptn/issues/8132)) ([f02aeba](https://github.com/keptn/keptn/commit/f02aeba6d384b9d7bb2aa50ce1796fc9e4cb3d6d))
* Add OAuth scopes to swagger and add possibility to disable deprecated APIs ([#8051](https://github.com/keptn/keptn/issues/8051)) ([0dc1203](https://github.com/keptn/keptn/commit/0dc12034548d1003bd02f331979e7219f4557d73))
* **api:** Create import endpoint ([#8137](https://github.com/keptn/keptn/issues/8137)) ([75ae009](https://github.com/keptn/keptn/commit/75ae0093c16ebb0396ee0eca3ac6c47592510a19))
* **bridge:** Add approval-item-module ([#8069](https://github.com/keptn/keptn/issues/8069)) ([15050ba](https://github.com/keptn/keptn/commit/15050ba60ce48a171d2ccf9926c480243e623f37))
* **bridge:** Add deletion dialog module ([#8060](https://github.com/keptn/keptn/issues/8060)) ([bac2bc8](https://github.com/keptn/keptn/commit/bac2bc8e2fd1870ccf188e504c5b02eab1d155f5))
* **bridge:** Add ktb-confirmation-dialog module ([#8058](https://github.com/keptn/keptn/issues/8058)) ([dfc286e](https://github.com/keptn/keptn/commit/dfc286e516e2ed185335c94832f202fbddc0b752))
* **bridge:** Add ktb-copy-to-clipboard module ([#8072](https://github.com/keptn/keptn/issues/8072)) ([473fce5](https://github.com/keptn/keptn/commit/473fce5baa601fee106af1d1c4b409e7cc74911c))
* **bridge:** Add ktb-create-service-module ([#8073](https://github.com/keptn/keptn/issues/8073)) ([ff73348](https://github.com/keptn/keptn/commit/ff73348b5a817981195d76ade00334463fbc541a))
* **bridge:** Add ktb-loading module ([#8048](https://github.com/keptn/keptn/issues/8048)) ([b6717fd](https://github.com/keptn/keptn/commit/b6717fdd1a7beee657a9fbccd5c609ee900974ef))
* **bridge:** Add modules in a bulk (01) ([#8077](https://github.com/keptn/keptn/issues/8077)) ([eeef827](https://github.com/keptn/keptn/commit/eeef827260015b8aa4b0bbf669b4f3926e9f4f3f))
* **bridge:** Add modules in a bulk (02) ([#8091](https://github.com/keptn/keptn/issues/8091)) ([1cc9a44](https://github.com/keptn/keptn/commit/1cc9a44e47f6505af3e6124d6e76f54580ab628d))
* **bridge:** Add modules in a bulk (03) ([#8125](https://github.com/keptn/keptn/issues/8125)) ([a28be76](https://github.com/keptn/keptn/commit/a28be76bca7bc05d3af5139d9afba376d2433c5e))
* **bridge:** Add sli-breakdown-module ([#8062](https://github.com/keptn/keptn/issues/8062)) ([dcd09da](https://github.com/keptn/keptn/commit/dcd09da0b72ff98eff4f2c274a9e2ebb947d5399))
* **bridge:** Bundle size report ([#8274](https://github.com/keptn/keptn/issues/8274)) ([ef3c504](https://github.com/keptn/keptn/commit/ef3c504f3d7cf47fa7ef6d9d1fe1ddb917461e63))
* **bridge:** Cleanup app modules and fix missing modules ([#8199](https://github.com/keptn/keptn/issues/8199)) ([58ada1e](https://github.com/keptn/keptn/commit/58ada1ee68176522dd6aaaa37433ce50476f7a54))
* **bridge:** Introduce Configuration and ComponentLogger ([#8042](https://github.com/keptn/keptn/issues/8042)) ([aa4bcf0](https://github.com/keptn/keptn/commit/aa4bcf0eda42faa8d108bafee7a711ad60f73ad0))
* **bridge:** introduce modules for ktb-proxy-input and others ([#8127](https://github.com/keptn/keptn/issues/8127)) ([258c5a6](https://github.com/keptn/keptn/commit/258c5a603498ef89d7c563298c68b9b273ec4a0f)), closes [#7932](https://github.com/keptn/keptn/issues/7932) [#7932](https://github.com/keptn/keptn/issues/7932) [#7932](https://github.com/keptn/keptn/issues/7932)
* **bridge:** Introduce modules for ktb-sequence-controls and others [#7932](https://github.com/keptn/keptn/issues/7932) ([#8139](https://github.com/keptn/keptn/issues/8139)) ([448e53f](https://github.com/keptn/keptn/commit/448e53fd888a37ca865b431dbdab1db3ce711f57))
* **bridge:** Introduce modules for ktb-sequence-state-info and others [#7932](https://github.com/keptn/keptn/issues/7932) ([#8119](https://github.com/keptn/keptn/issues/8119)) ([e9ff5cc](https://github.com/keptn/keptn/commit/e9ff5cc47a3a5c9d4d414d1dea9c4c91397facf0))
* **bridge:** Introduce modules for sequence-timeline and others [#7932](https://github.com/keptn/keptn/issues/7932) ([#8153](https://github.com/keptn/keptn/issues/8153)) ([e7b2ec6](https://github.com/keptn/keptn/commit/e7b2ec6527f429a8ea28088c96f44a42d5463eb1))
* **bridge:** ktb-certificate-input module ([#8071](https://github.com/keptn/keptn/issues/8071)) ([8cf36aa](https://github.com/keptn/keptn/commit/8cf36aac4155d4061fdaf7b90fb600f647a52903))
* **bridge:** ktb-evaluation-details module ([#8066](https://github.com/keptn/keptn/issues/8066)) ([e7640dd](https://github.com/keptn/keptn/commit/e7640dddc420a9e45990e7973068a90f9180d334))
* **bridge:** Refactor dashboard to use interfaces ([#8205](https://github.com/keptn/keptn/issues/8205)) ([2cbbc2d](https://github.com/keptn/keptn/commit/2cbbc2da7c59791662af4932bf134c6257b6c788))
* **bridge:** Refactor DataService's loadProjects ([#8268](https://github.com/keptn/keptn/issues/8268)) ([8c55b1b](https://github.com/keptn/keptn/commit/8c55b1beb3f3c0cb0e5199bbf8c439f27ae9f342))
* **bridge:** Rename dashboard to dashboard-legacy ([#8097](https://github.com/keptn/keptn/issues/8097)) ([cffbf50](https://github.com/keptn/keptn/commit/cffbf502edb3c128c99ef79b6c98386300c79d43))
* **bridge:** RX-ify the dashboard component ([#8167](https://github.com/keptn/keptn/issues/8167)) ([6d1c05d](https://github.com/keptn/keptn/commit/6d1c05d9ad4047b5bdb941562f27abdc149aab52))
* **bridge:** Support configured AUTH_MSG ([#8043](https://github.com/keptn/keptn/issues/8043)) ([0589b26](https://github.com/keptn/keptn/commit/0589b262fa1844097ee499c8cc153fdefca4f614))
* **bridge:** Support new webhook.yaml version v1beta1 ([#8247](https://github.com/keptn/keptn/issues/8247)) ([bad1ee7](https://github.com/keptn/keptn/commit/bad1ee757c456d1d959ed80934053fb08e1e6c6a))
* **bridge:** Use Configuration Object instead of Env Var ([#8096](https://github.com/keptn/keptn/issues/8096)) ([6a3bc4d](https://github.com/keptn/keptn/commit/6a3bc4daca9f588936af026558acf39e46795cc6))
* **cp-connector:** Connect to NATS only at event source startup ([#8064](https://github.com/keptn/keptn/issues/8064)) ([9793f4e](https://github.com/keptn/keptn/commit/9793f4e252a8cff7621f4f0dffb1d799118b09c9))
* **cp-connector:** HTTP based EventSource implementation ([#8140](https://github.com/keptn/keptn/issues/8140)) ([5e2f548](https://github.com/keptn/keptn/commit/5e2f548ac192064c86414d64abf39fccd8e71654))
* **cp-connector:** Injectable logger implementation ([#8024](https://github.com/keptn/keptn/issues/8024)) ([d074978](https://github.com/keptn/keptn/commit/d0749780a7e57cdf433e0120bcc6699efda0778d))
* **go-sdk:** Use APISet instead of resource handler ([#8059](https://github.com/keptn/keptn/issues/8059)) ([8e00834](https://github.com/keptn/keptn/commit/8e00834cbe93c7cd854c415f5cfd38b9b24817f5))
* **installer:** Helm Chart revamp ([#7678](https://github.com/keptn/keptn/issues/7678)) ([f78f867](https://github.com/keptn/keptn/commit/f78f867743e6477c28b311f4cc84db60bc1f5df3))


### Bug Fixes

* Added longer retry in provisioning URL test ([#8074](https://github.com/keptn/keptn/issues/8074)) ([2d97f9c](https://github.com/keptn/keptn/commit/2d97f9c451397432a910224a22ac862c93a325d0))
* Added proxy to integration test ([#8052](https://github.com/keptn/keptn/issues/8052)) ([52509d6](https://github.com/keptn/keptn/commit/52509d6a1ed4907dc8aeb47af2139f63e52b26ee))
* **bridge:** Corrected misleading message in creating project ([#8142](https://github.com/keptn/keptn/issues/8142)) ([6a1d013](https://github.com/keptn/keptn/commit/6a1d013ea6b0e51aa11858ba8197a507021f392d))
* **bridge:** Fix 'view more' of quick filter ([#8306](https://github.com/keptn/keptn/issues/8306)) ([9453e5b](https://github.com/keptn/keptn/commit/9453e5ba7a852d7e129221f82dc5805e22157b6a))
* **bridge:** Fix approval being sent twice ([#8004](https://github.com/keptn/keptn/issues/8004)) ([3a31f55](https://github.com/keptn/keptn/commit/3a31f552ccb610af982792b91e91aa216ecf71e3))
* **bridge:** Fix broken UI if connection was lost ([#8050](https://github.com/keptn/keptn/issues/8050)) ([746be23](https://github.com/keptn/keptn/commit/746be23f7f86fe5a5b9c4b3173a81c061168749c))
* **bridge:** Fix incorrect selected stage on refresh ([#7974](https://github.com/keptn/keptn/issues/7974)) ([9abd6a3](https://github.com/keptn/keptn/commit/9abd6a31731687646f155b69bb069e051384a921))
* **bridge:** Fix missing evaluation score of sequence ([#8032](https://github.com/keptn/keptn/issues/8032)) ([3fe27e0](https://github.com/keptn/keptn/commit/3fe27e039875bf01e33ddffd3de878bf25df20bb))
* **bridge:** Fix missing sequence menu icon selection ([#8308](https://github.com/keptn/keptn/issues/8308)) ([d841387](https://github.com/keptn/keptn/commit/d8413876fa724291f74c5095d203489297057463))
* **bridge:** Handle invalid bridge versions ([#8283](https://github.com/keptn/keptn/issues/8283)) ([7a17271](https://github.com/keptn/keptn/commit/7a17271df64f3fb4947c6147a22ace3e29d16918))
* **bridge:** Remove previous filter from URL ([#7998](https://github.com/keptn/keptn/issues/7998)) ([fcd19ac](https://github.com/keptn/keptn/commit/fcd19ac27f8270ed1387abed2bbf42276f658d61))
* **bridge:** Respond with a default version payload, when the call to get.keptn.sh/version.json fails ([#8037](https://github.com/keptn/keptn/issues/8037)) ([b4be4ca](https://github.com/keptn/keptn/commit/b4be4cac68f35df15a7f959680e99302f7808320))
* **bridge:** Save client secret in k8s secret ([#8269](https://github.com/keptn/keptn/issues/8269)) ([27f1b6a](https://github.com/keptn/keptn/commit/27f1b6a646743c1984ec2d92c5ca4940c5cae015))
* **bridge:** Settings view overflow problem ([#8291](https://github.com/keptn/keptn/issues/8291)) ([f473eb6](https://github.com/keptn/keptn/commit/f473eb6b24a4a2a3c68ca64452fa2de6ed0f810e))
* **bridge:** Show all evaluations in the environment screen ([#8090](https://github.com/keptn/keptn/issues/8090)) ([ffb937c](https://github.com/keptn/keptn/commit/ffb937ce40d0204daedcfaa37400d334ee693143))
* **bridge:** Show loading indicator for sequences before filters are applied the first time ([#8033](https://github.com/keptn/keptn/issues/8033)) ([04a7eb8](https://github.com/keptn/keptn/commit/04a7eb88b12156941dd78fccb515eca4c9b4a6bd))
* **bridge:** Show the heatmap even if the SLO of an evaluation is invalid ([#7965](https://github.com/keptn/keptn/issues/7965)) ([d0edcbc](https://github.com/keptn/keptn/commit/d0edcbc1052b4dff02741f2d715ce3302bb97c7a))
* **bridge:** Update projects if dashboard is visited ([#7997](https://github.com/keptn/keptn/issues/7997)) ([e201bc1](https://github.com/keptn/keptn/commit/e201bc1aba518ffa6521a5ea097d597ea9ac733b))
* Change name label to respect the nameOverride ([#8249](https://github.com/keptn/keptn/issues/8249)) ([6f6af8b](https://github.com/keptn/keptn/commit/6f6af8bcf38ea70752bdf5c05ec8c79bb96a3c95))
* **cli:** Skip version check for auth sub command ([#8126](https://github.com/keptn/keptn/issues/8126)) ([0b03dd0](https://github.com/keptn/keptn/commit/0b03dd07a3ccc21fef8626d4aa491a7bdae6e1e5))
* **cp-connector:** Added return of error in queue subscribe function ([#8101](https://github.com/keptn/keptn/issues/8101)) ([7285f51](https://github.com/keptn/keptn/commit/7285f5120cd65df0c0d9cfc550c1517cedd32b91))
* **cp-connector:** Synchronized shutdown of cp-connector during cancellation ([#8063](https://github.com/keptn/keptn/issues/8063)) ([a3f3010](https://github.com/keptn/keptn/commit/a3f301010f4f4bc48aa1fdae9c39ecd34ba2c582))
* **distributor:** Limit payload size sent to the distributor's API proxy ([#8200](https://github.com/keptn/keptn/issues/8200)) ([d40ee5b](https://github.com/keptn/keptn/commit/d40ee5bb5585f5d8110b6f1241fe13d285047d2a))
* **installer:** Add missing quotes to env var for distributor ([#8157](https://github.com/keptn/keptn/issues/8157)) ([4fcf792](https://github.com/keptn/keptn/commit/4fcf792d9197872df5f9ea918cdf49baa947fd13))
* **installer:** Revert immutable k8s labels ([#8213](https://github.com/keptn/keptn/issues/8213)) ([bed7b04](https://github.com/keptn/keptn/commit/bed7b045a32dc3f4ff77410dfcda6dafe7c40f9b))
* Integration tests ([#8198](https://github.com/keptn/keptn/issues/8198)) ([23038a1](https://github.com/keptn/keptn/commit/23038a14775f2ad82191c900bcc9789757d73548))
* Only trigger CLI command docs auto-generation for full release tags ([#8120](https://github.com/keptn/keptn/issues/8120)) ([8ffe5fc](https://github.com/keptn/keptn/commit/8ffe5fcdb85418d00585faab742a8a591fbb84e6))
* **resource-service:** Always delete local project directory when project creation fails ([#8123](https://github.com/keptn/keptn/issues/8123)) ([44cbcb3](https://github.com/keptn/keptn/commit/44cbcb33ec2f33e2240e7067f67cf6e6ffca1002))
* **resource-service:** Remove token enforcement ([#8040](https://github.com/keptn/keptn/issues/8040)) ([44f9a4a](https://github.com/keptn/keptn/commit/44f9a4a059addd2eb7a85d7d0a5fe6e428eba813))
* **shipyard-controller:** Add time property to EventFilter ([#8134](https://github.com/keptn/keptn/issues/8134)) ([37bb437](https://github.com/keptn/keptn/commit/37bb437670e5174a19a7eac368b9c9c67b2a892e))
* **shipyard-controller:** Fix project deletion unit test ([#8231](https://github.com/keptn/keptn/issues/8231)) ([12a60f2](https://github.com/keptn/keptn/commit/12a60f267115dc88275c411719adf575a5bafd50))
* **shipyard-controller:** Include namespace in call to provisioning service ([#8041](https://github.com/keptn/keptn/issues/8041)) ([9429678](https://github.com/keptn/keptn/commit/94296789a43ae6a7a4ad8491261a1d0b9ac7e8fb))
* **shipyard-controller:** Project should be deleted even if upstream delete fails ([#8204](https://github.com/keptn/keptn/issues/8204)) ([314c93a](https://github.com/keptn/keptn/commit/314c93abc506882131c32913e57433ab0905eded))
* **shipyard-controller:** Return `ErrProjectNotFound` instead of `nil, nil` when project is not in the db ([#8266](https://github.com/keptn/keptn/issues/8266)) ([2d20f6f](https://github.com/keptn/keptn/commit/2d20f6f891436b81ef168193420addbe55ee4a47))
* Use distributor values namespace and hostname in svc env vars ([#8297](https://github.com/keptn/keptn/issues/8297)) ([7140f5b](https://github.com/keptn/keptn/commit/7140f5ba6d01ec357bc38e32ebd441b2cf47103a))


### Docs

* **cli:** Improve CLI Documentation ([#8061](https://github.com/keptn/keptn/issues/8061)) ([922ba5b](https://github.com/keptn/keptn/commit/922ba5b16aac953b82224a833391291b012b9659))
* Typo: we are using swagger.yaml not swagger.json ([#8099](https://github.com/keptn/keptn/issues/8099)) ([ee6e18b](https://github.com/keptn/keptn/commit/ee6e18bba072c1b491961ca32774213e5e976cd8))


### Refactoring

* **bridge:** Introduce modules for app-header and environment components ([#8158](https://github.com/keptn/keptn/issues/8158)) ([c2174cf](https://github.com/keptn/keptn/commit/c2174cfb977075b9bf0540dca0845e3c433936ca))
* **bridge:** Make use of new Git API model ([#8180](https://github.com/keptn/keptn/issues/8180)) ([8da8df8](https://github.com/keptn/keptn/commit/8da8df85287e01b5bb4cc75bc14ef8fd2041abad))
* **bridge:** update sequence screen data model ([#8083](https://github.com/keptn/keptn/issues/8083)) ([e031b2f](https://github.com/keptn/keptn/commit/e031b2f0c7c8a4b0c478344ef623b031411d51ac))


### Other

*  Updated webhook and remediation to new sdk ([#8170](https://github.com/keptn/keptn/issues/8170)) ([adfa700](https://github.com/keptn/keptn/commit/adfa7001d6dc2d04cd0489ebf3ba7adf7d2d4ed1))
* Add [@sarahhuber001](https://github.com/sarahhuber001) as member ([#7893](https://github.com/keptn/keptn/issues/7893)) ([1709806](https://github.com/keptn/keptn/commit/17098067f3be2830a2ead91da07aaffad4a8a1cb))
* Add @STRRL to CONTRIBUTORS ([#8149](https://github.com/keptn/keptn/issues/8149)) ([a2745b8](https://github.com/keptn/keptn/commit/a2745b81dd948cad50d51623c4cd5377868aa776))
* **bridge:** Added missing modules ([#8147](https://github.com/keptn/keptn/issues/8147)) ([f436de5](https://github.com/keptn/keptn/commit/f436de55c0942d4831a140337e8e18707d67f4cd))
* **bridge:** Added missing modules for evaluation-details ([#8156](https://github.com/keptn/keptn/issues/8156)) ([c4d75c2](https://github.com/keptn/keptn/commit/c4d75c21e69d1c878463afb259d5573350a4e717))
* **bridge:** Improve has-logs polling ([#8039](https://github.com/keptn/keptn/issues/8039)) ([8b67a23](https://github.com/keptn/keptn/commit/8b67a23746d7fd27891f16b2c4c1e8feaf138890))
* **bridge:** Removed remediation config, only poll remediations when needed ([#8217](https://github.com/keptn/keptn/issues/8217)) ([63bb742](https://github.com/keptn/keptn/commit/63bb7421ed39049fd426760966e4557ae779dcca))
* Bump swagger-ui to version 4.12.0 ([#8279](https://github.com/keptn/keptn/issues/8279)) ([7a9997a](https://github.com/keptn/keptn/commit/7a9997a5aa0d4d7b6ca03e3f1c56159966f38fe2))
* **cli:** Deprecate install uninstall and upgrade commands ([#8103](https://github.com/keptn/keptn/issues/8103)) ([d9c8d58](https://github.com/keptn/keptn/commit/d9c8d585a001ff4885bfb3c506cc08e5ee0bf5cf))
* cp-connector package restructuring ([#7910](https://github.com/keptn/keptn/issues/7910)) ([9072004](https://github.com/keptn/keptn/commit/90720040a24d5436634e931c42dc7f6ba48fd5f2))
* **cp-connector:** Added debug logs to controlplane ([#8012](https://github.com/keptn/keptn/issues/8012)) ([4f4069f](https://github.com/keptn/keptn/commit/4f4069f42d380975be6f217b66f97fc9c833119f))
* **cp-connector:** Additional debug logs ([#8016](https://github.com/keptn/keptn/issues/8016)) ([efe9ad5](https://github.com/keptn/keptn/commit/efe9ad550f8e4e7d01be7fe762bf94f5004fce54))
* **cp-connector:** Fixed missing error in queuesubscribe for nats ([#8122](https://github.com/keptn/keptn/issues/8122)) ([d57cd8c](https://github.com/keptn/keptn/commit/d57cd8c91758a1afe271ef6c302923adbca3d02d))
* **installer:** Added  API_PROXY_HTTP_TIMEOUT to distributor helm values ([#8138](https://github.com/keptn/keptn/issues/8138)) ([b84391f](https://github.com/keptn/keptn/commit/b84391fef7cc07f5eb26528a15a36224a147f726))
* **installer:** Moved automaticProvisionMsg under features ([#8145](https://github.com/keptn/keptn/issues/8145)) ([d1dcecb](https://github.com/keptn/keptn/commit/d1dcecb18af5d9f0b141fb1e6c833340000d2a0f))
* Mark kubernetes-utils as deprecated ([#8117](https://github.com/keptn/keptn/issues/8117)) ([9ba17c0](https://github.com/keptn/keptn/commit/9ba17c008bd26189ee4af0150a5a69b1ad75d965))
* Remove configuration-service from pipelines ([#8284](https://github.com/keptn/keptn/issues/8284)) ([6b136eb](https://github.com/keptn/keptn/commit/6b136ebcb57b259ac8279d58f528cc621dd2e8f7))
* Remove reference to go-sdk from renovate.json ([#8229](https://github.com/keptn/keptn/issues/8229)) ([5d14929](https://github.com/keptn/keptn/commit/5d14929bf74869920938bff6e966bb77154ea1bf))
* Removed BETA from uniform API ([#8135](https://github.com/keptn/keptn/issues/8135)) ([f1c6c7d](https://github.com/keptn/keptn/commit/f1c6c7dfbb2654cac2fa8f7f1ff763040fb1b72e))
* **shipyard-controller:** Move integration tests to faster component tests ([#8087](https://github.com/keptn/keptn/issues/8087)) ([4303cff](https://github.com/keptn/keptn/commit/4303cfffd79db3b1dad73000e3315a8f3900bec1))
* **shipyard-controller:** Remove references to deprecated subscription ([#8035](https://github.com/keptn/keptn/issues/8035)) ([18afeb4](https://github.com/keptn/keptn/commit/18afeb45fb41acd240d088fe970415eec851bf1f))
* Update cp connector ([#8133](https://github.com/keptn/keptn/issues/8133)) ([38cd84b](https://github.com/keptn/keptn/commit/38cd84bea1f816021072076774301c19e254a825))
* Update cp-connector ref in go-sdk ([#8094](https://github.com/keptn/keptn/issues/8094)) ([23d1878](https://github.com/keptn/keptn/commit/23d18789f3c812e489cd63f5d8d8443301edd4a5))
* Updated k8s dependencies ([#8173](https://github.com/keptn/keptn/issues/8173)) ([87cc798](https://github.com/keptn/keptn/commit/87cc79834333dfa965050f107caa95fc94c80fb6))
* Use logrus StandardLogger in webhook and remediation service ([#8292](https://github.com/keptn/keptn/issues/8292)) ([fd5c201](https://github.com/keptn/keptn/commit/fd5c201599049c2bd0808a551d7a60170bcbd394))

## [0.16.0](https://github.com/keptn/keptn/compare/0.15.0...0.16.0) (2022-06-07)


### ⚠ BREAKING CHANGES

* The `resource-service` replaces the old `configuration-service`. The new service always requires a Git upstream to be configured for a Keptn project. The new service will bring many advantages, such as faster response times and the possibility to upgrade Keptn without any downtime.

### Features

* Add ability to customize client_max_body_size in the nginx gateway ([#7727](https://github.com/keptn/keptn/issues/7727)) ([d27033b](https://github.com/keptn/keptn/commit/d27033bcbc20770f73bd758b9e1181b49b62e344))
* **api:** Send events directly to nats instead via distributor ([#7672](https://github.com/keptn/keptn/issues/7672)) ([58f9615](https://github.com/keptn/keptn/commit/58f9615679006fd74c3fb23bceb448499b7aba02))
* **approval-service:** Consider nats connection in readiness probe ([#7723](https://github.com/keptn/keptn/issues/7723)) ([d170354](https://github.com/keptn/keptn/commit/d1703544a4999feeac7ac0f4c49d141b01f680d6))
* **approval-service:** Run approval-service without distributor sideCar ([#7689](https://github.com/keptn/keptn/issues/7689)) ([bceaf4b](https://github.com/keptn/keptn/commit/bceaf4b5e298170073c566fc28361ae474e80898))
* **bridge:** Automatic provisioning url makes git form optional and ap message can be set and displayed ([60bd257](https://github.com/keptn/keptn/commit/60bd2573c02a875c44cfaa876be23d36ad6597d5))
* **bridge:** Implement heatmap with d3 ([#7658](https://github.com/keptn/keptn/issues/7658)) ([84dc4a0](https://github.com/keptn/keptn/commit/84dc4a08e806446edb6cef2da3ac9435d0af042b))
* **bridge:** Introduce Module (ktb-notification) ([#7897](https://github.com/keptn/keptn/issues/7897)) ([a87254a](https://github.com/keptn/keptn/commit/a87254a911001d86c197f4e3546f262e71fe9168))
* **bridge:** Make filters in sequence view stable across page refresh ([#7526](https://github.com/keptn/keptn/issues/7526)) ([0b18e45](https://github.com/keptn/keptn/commit/0b18e45695648e96d37f5d61b621fe0301de6e64))
* **bridge:** Remove millis and micros from evaluation time frame ([#7774](https://github.com/keptn/keptn/issues/7774)) ([15b4735](https://github.com/keptn/keptn/commit/15b4735e43a2d3b380342c1939b4af7aac6de25c))
* **bridge:** Remove polling for evaluation history in environment screen ([#7851](https://github.com/keptn/keptn/issues/7851)) ([71874bd](https://github.com/keptn/keptn/commit/71874bdac5da1ab4cc1c8fe7ffb961cd49858880))
* **bridge:** Remove polling for services in settings screen ([#7853](https://github.com/keptn/keptn/issues/7853)) ([b99032c](https://github.com/keptn/keptn/commit/b99032c8c9d808535df6e02a4943cf05c7e3300f))
* **bridge:** Remove tag input field for creating a sequence ([#7757](https://github.com/keptn/keptn/issues/7757)) ([2e16548](https://github.com/keptn/keptn/commit/2e16548b8f3775e7252d60b6645254b2841920f9))
* **bridge:** Removes projects polling on dashboard [#7796](https://github.com/keptn/keptn/issues/7796) ([#7812](https://github.com/keptn/keptn/issues/7812)) ([7a71e05](https://github.com/keptn/keptn/commit/7a71e056f59e65edc99918763aa71759b16af6c5))
* **bridge:** Trigger sequence - Remove polling for custom sequences ([#7813](https://github.com/keptn/keptn/issues/7813)) ([138a773](https://github.com/keptn/keptn/commit/138a77374deb34e69adb67299d7cff22b67c193b))
* **bridge:** Use ktb-heatmap component ([#7816](https://github.com/keptn/keptn/issues/7816)) ([5bca4bd](https://github.com/keptn/keptn/commit/5bca4bd3a077f8eef400508d61b81d9126611ab7))
* Change default values of preStop hook time and grace period ([#7682](https://github.com/keptn/keptn/issues/7682)) ([a31023b](https://github.com/keptn/keptn/commit/a31023b46f36e7566886767423c23681177d6214))
* **cp-connector:** Ensure mandatory CloudEvent attributes are set ([#7744](https://github.com/keptn/keptn/issues/7744)) ([becb01f](https://github.com/keptn/keptn/commit/becb01f3a1b68d666f4f0c96b611975a20c3574a))
* **cp-connector:** Introduce log forwarding to `cp-connector` library ([#7713](https://github.com/keptn/keptn/issues/7713)) ([c36faf0](https://github.com/keptn/keptn/commit/c36faf0071dc689abbc692c460dfd35981f4b962))
* **cp-connector:** Make sure event timestamp is always set  ([#7743](https://github.com/keptn/keptn/issues/7743)) ([6473142](https://github.com/keptn/keptn/commit/64731420f4215b1cc45a0f75f2740612e1169ae5))
* Enable Resource-Service by default ([#7826](https://github.com/keptn/keptn/issues/7826)) ([73d264b](https://github.com/keptn/keptn/commit/73d264b376a20f52f71d63e6960a58eeb5dcdb34))
* **lighthouse-service:** Adapt readiness probe of lighthouse service to consider nats subscription ([#7735](https://github.com/keptn/keptn/issues/7735)) ([51837d7](https://github.com/keptn/keptn/commit/51837d7e3ddd893b61476fb6c0e23c734ece5108))
* **lighthouse-service:** Run lighthouse-service without distributor sidecar ([#7691](https://github.com/keptn/keptn/issues/7691)) ([b2ad6ad](https://github.com/keptn/keptn/commit/b2ad6adf4cdfa11aab37e5031b116fd586a0f43c))
* **mongodb-datastore:** Use cp-connector library ([#7685](https://github.com/keptn/keptn/issues/7685)) ([defee50](https://github.com/keptn/keptn/commit/defee50b0a039ecec36312a315dc466141ebb31c))
* Refactor `go-sdk` to use `cp-connector` internally ([#7686](https://github.com/keptn/keptn/issues/7686)) ([1712149](https://github.com/keptn/keptn/commit/171214903be538942fbc32e41435a1746abd56cc))
* **resource-service:** Removed NATS ([#7694](https://github.com/keptn/keptn/issues/7694)) ([fa48649](https://github.com/keptn/keptn/commit/fa48649afe87d11bdebc773d37f695ca6fca1c87))


### Bug Fixes

* Added retry to url provisioning integration test ([#7815](https://github.com/keptn/keptn/issues/7815)) ([93095eb](https://github.com/keptn/keptn/commit/93095ebcc9bc544377d081730b54adca0fdca7c9))
* **approval-service:** Use deployment name for registration name to fix queue group functionality ([#7718](https://github.com/keptn/keptn/issues/7718)) ([42cf370](https://github.com/keptn/keptn/commit/42cf370969c21b46d000098bc660ae368e84120b))
* **bridge:** Add missing update project notification ([#7770](https://github.com/keptn/keptn/issues/7770)) ([4bdaa71](https://github.com/keptn/keptn/commit/4bdaa717415d34ed2866b32224e3fb7a53b916a2))
* **bridge:** Allow Webhook configuration URL to be a secret ([#7728](https://github.com/keptn/keptn/issues/7728)) ([0372484](https://github.com/keptn/keptn/commit/0372484907c9a7ce0a70f1397fb672205b2c4e52))
* **bridge:** Duplicate headline in project settings page ([#7988](https://github.com/keptn/keptn/issues/7988)) ([1645230](https://github.com/keptn/keptn/commit/1645230fffa3df299c6270b13bb8196ce5e3eb23))
* **bridge:** Fix D3 heatmap selection ([#7842](https://github.com/keptn/keptn/issues/7842)) ([c15740a](https://github.com/keptn/keptn/commit/c15740a9cc6e72531b430bede317509ff59e4fc0))
* **bridge:** Fix flickering of filter in sequence view ([#8009](https://github.com/keptn/keptn/issues/8009)) ([3e6ec42](https://github.com/keptn/keptn/commit/3e6ec428adab2b52db2c7ce1d7f40bdef47cca56))
* **bridge:** Fix logout not being visible if metadata is not returned ([#7794](https://github.com/keptn/keptn/issues/7794)) ([1c2b196](https://github.com/keptn/keptn/commit/1c2b1967442db69e6d398e07631370d8b3a1e63e))
* **bridge:** Fixed D3 heatmap issues ([#7833](https://github.com/keptn/keptn/issues/7833)) ([3e697bf](https://github.com/keptn/keptn/commit/3e697bf16d700d6398ca3c907b0b1a7c14068843))
* **bridge:** Fixed missing 'View service details' button ([#7806](https://github.com/keptn/keptn/issues/7806)) ([41cb52d](https://github.com/keptn/keptn/commit/41cb52de9552d1d9b56d1625be3474b4360924e2))
* **bridge:** Fixed triggering of validation on token reset ([#7766](https://github.com/keptn/keptn/issues/7766)) ([85dc15b](https://github.com/keptn/keptn/commit/85dc15b92eb222c098aa457166cb354812c2fc15))
* **cli:** Remove unnecessary `--sequence` flag from `keptn trigger sequence` ([#7902](https://github.com/keptn/keptn/issues/7902)) ([b252b6d](https://github.com/keptn/keptn/commit/b252b6d70ec70d34afc2dc02709e9cde30a6544f))
* Correctly match nginx location for Bridge ([#7729](https://github.com/keptn/keptn/issues/7729)) ([dd236ef](https://github.com/keptn/keptn/commit/dd236ef4fa97ece1b9ba3703c2b66e939dbc9215))
* **cp-connector:** Added missing Register() method to FixedSubscriptionSource ([#7731](https://github.com/keptn/keptn/issues/7731)) ([fe5b978](https://github.com/keptn/keptn/commit/fe5b9782f67d7ae14524d38f1712ea0cf3eb9703))
* **cp-connector:** fix passing deduplicated subjects to nats subscriber ([#7782](https://github.com/keptn/keptn/issues/7782)) ([39124e1](https://github.com/keptn/keptn/commit/39124e19d65f7455ad2aaffee0110fd33b503cba))
* **cp-connector:** Flaky unit test ([#7976](https://github.com/keptn/keptn/issues/7976)) ([be9cafb](https://github.com/keptn/keptn/commit/be9cafbd2b97d36e86198b2f8ca78993639ccdc4))
* **cp-connector:** Subscribe to integrations before creating a job ([#7952](https://github.com/keptn/keptn/issues/7952)) ([ccc4f26](https://github.com/keptn/keptn/commit/ccc4f2648267e54bc6b92893120183420064c409))
* **cp-connector:** Unsubscribe instead of disconnect from nats on cancel ([#7795](https://github.com/keptn/keptn/issues/7795)) ([8854339](https://github.com/keptn/keptn/commit/885433905cc0bc0eb727d89986c41d048c5cfd46))
* Disallow calls to `SendEvent` or `GetMetaData` when used via `InternalAPISet` ([#7939](https://github.com/keptn/keptn/issues/7939)) ([d683005](https://github.com/keptn/keptn/commit/d683005ae55f7f871821e0f492e040d651369050))
* Ensure that all mongodb cursors are being closed after use ([#7909](https://github.com/keptn/keptn/issues/7909)) ([01c0a9d](https://github.com/keptn/keptn/commit/01c0a9df26ca7222611aec26aa2c365ac17fe5f2))
* Fixed zd test to run without UI tests ([#7908](https://github.com/keptn/keptn/issues/7908)) ([bd8fb20](https://github.com/keptn/keptn/commit/bd8fb209478eea9536002697bd7f59e9234d878f))
* **go-sdk:** Return from event handler when wg for graceful shutdown cannot be retrieved from context ([#7810](https://github.com/keptn/keptn/issues/7810)) ([2c2ed2c](https://github.com/keptn/keptn/commit/2c2ed2c2c496c360c1b082a187ed54952e4e000e))
* **go-sdk:** Use the correct env var for setting the integration version ([#7930](https://github.com/keptn/keptn/issues/7930)) ([cd130b7](https://github.com/keptn/keptn/commit/cd130b7b8291b16fc3aecac8347b06e15a4ee533))
* **installer:** Adapt preStop hook times for lighthouse, statistics and webhook ([#7947](https://github.com/keptn/keptn/issues/7947)) ([3e9f9b5](https://github.com/keptn/keptn/commit/3e9f9b51839aa0da3fdea1372a4698d3a8a3ec97))
* **installer:** Add resource service to airgapped installer script ([#7869](https://github.com/keptn/keptn/issues/7869)) ([2196c11](https://github.com/keptn/keptn/commit/2196c11be8ceddb2935317675d760a13e47efb09))
* **installer:** Configure default preStopHook and grafefulPeriods timeouts ([#7926](https://github.com/keptn/keptn/issues/7926)) ([7a6489c](https://github.com/keptn/keptn/commit/7a6489c3360521a9294676f1b382a20b2fdee4c8))
* **installer:** Fix airgapped setup not finding correct nginx image ([#7935](https://github.com/keptn/keptn/issues/7935)) ([2ee4bab](https://github.com/keptn/keptn/commit/2ee4bab3b1a885cc749be0b8ed7acf9e896fba45))
* **installer:** Fix wrong regex for log location ([#7921](https://github.com/keptn/keptn/issues/7921)) ([295099d](https://github.com/keptn/keptn/commit/295099dc4de84f95832f1f644c5268f0874e7c33))
* **installer:** Redirect output of preStop hook command to /dev/null ([#7837](https://github.com/keptn/keptn/issues/7837)) ([117f1fb](https://github.com/keptn/keptn/commit/117f1fb9a35540a76e7c5c20e8908e81af660568))
* **installer:** Use exec preStop hook for shipyard controller ([#7768](https://github.com/keptn/keptn/issues/7768)) ([283f72f](https://github.com/keptn/keptn/commit/283f72f71950c0341967a3c26fb85b03c716151d))
* **lighthouse-service:** Ensure sloFileContent property is always a base64 encoded string ([#7892](https://github.com/keptn/keptn/issues/7892)) ([e19fcfc](https://github.com/keptn/keptn/commit/e19fcfc9291c8751f0cdc34109b6f6ab5f48c197))
* Make sure that all events are being processed before shutting down lighthouse/approval service ([#7787](https://github.com/keptn/keptn/issues/7787)) ([0facb58](https://github.com/keptn/keptn/commit/0facb58c8a295aec4a594db7c965552056ff46fb))
* **mongodb-datastore:** Change name of integration to name of service instead of pod name ([#7777](https://github.com/keptn/keptn/issues/7777)) ([21d2774](https://github.com/keptn/keptn/commit/21d2774b751b7dde6b7bf341d4ecce5200eb9797))
* **mongodb-datastore:** Return [] instead of nil from get methods ([#7919](https://github.com/keptn/keptn/issues/7919)) ([4992bc5](https://github.com/keptn/keptn/commit/4992bc56f6632f8c76e91c1c81baef8409be4d1c))
* Removed wrong beta11 from webhook integration test ([#7861](https://github.com/keptn/keptn/issues/7861)) ([08ee81d](https://github.com/keptn/keptn/commit/08ee81d2bec8d452e0433bae2b559c2157044a9e))
* Return missing error in test-utils ([#7928](https://github.com/keptn/keptn/issues/7928)) ([d42af22](https://github.com/keptn/keptn/commit/d42af221680d7f24d49d8a52a47ab9b20f84ac91))
* **secret-service:** Deleting a secret does not remove references from related roles ([#7789](https://github.com/keptn/keptn/issues/7789)) ([56786b8](https://github.com/keptn/keptn/commit/56786b8b1aef7b87e12ce4b3ca26d1cb9b1ff6a8))
* **shipyard-controller:** Allow parallel sequence execution if the service is different ([#7775](https://github.com/keptn/keptn/issues/7775)) ([5f2dc74](https://github.com/keptn/keptn/commit/5f2dc7495ec33202a01712ef767ccbc41f872cfd))
* **shipyard-controller:** Avoid lost writes to subscriptions due to concurrent writes ([#7960](https://github.com/keptn/keptn/issues/7960)) ([1c9b40b](https://github.com/keptn/keptn/commit/1c9b40b7d11d6659fd208e40b56e2d27d829ca8f))
* **shipyard-controller:** Dispatch new sequence directly only if no older sequence is waiting ([#7793](https://github.com/keptn/keptn/issues/7793)) ([b8bad71](https://github.com/keptn/keptn/commit/b8bad7162821aba8c7d09794a4a1337756adee5e))
* **shipyard-controller:** Make sure result and status are set if sequence is timed out ([#7901](https://github.com/keptn/keptn/issues/7901)) ([81858c0](https://github.com/keptn/keptn/commit/81858c0372ee9d494f8976594598fe989b3731b9))
* **shipyard-controller:** Set the sequence execution state back to `started` when approval task has been finished ([#7838](https://github.com/keptn/keptn/issues/7838)) ([8444b48](https://github.com/keptn/keptn/commit/8444b481dbb315f234116530e0c1d03040436446))
* **shipyard-controller:** Update registration info based on integration name/namespace ([#8006](https://github.com/keptn/keptn/issues/8006)) ([d277a83](https://github.com/keptn/keptn/commit/d277a83915c65aaaf42544cd371d7f8a972e8494))
* **webhook-service:** Added denied curl in webhook beta based on host ([#7716](https://github.com/keptn/keptn/issues/7716)) ([d194367](https://github.com/keptn/keptn/commit/d1943671b8ccc5a13c50dafd833c486a54aedb9b))
* **webhook-service:** Added missing webhook-config version check ([#7832](https://github.com/keptn/keptn/issues/7832)) ([445000a](https://github.com/keptn/keptn/commit/445000a88258801aae32ad395b073d47dff9ffc7))


### Performance

* **bridge:** Use adapted sequence endpoint for project endpoint of bridge server ([#7696](https://github.com/keptn/keptn/issues/7696)) ([5bed56d](https://github.com/keptn/keptn/commit/5bed56d4b62cd996b09f8c15ab5a81e02aa03d70))


### Docs

* Added zero downtime tests documentation ([#7895](https://github.com/keptn/keptn/issues/7895)) ([cefdab5](https://github.com/keptn/keptn/commit/cefdab5fbb76b4d24e22b2a30549867880731ab7))
* Improve developer API + integration tests docs ([#7771](https://github.com/keptn/keptn/issues/7771)) ([b6fb2d6](https://github.com/keptn/keptn/commit/b6fb2d64afad324f75a85fec0ec24f6acd9d1cec))
* Improve documentation for resource-service ([#7765](https://github.com/keptn/keptn/issues/7765)) ([0995fda](https://github.com/keptn/keptn/commit/0995fda85ac681a4d82219109ad801c9132af553))
* Update version for installation of Helm and JMeter services ([#7700](https://github.com/keptn/keptn/issues/7700)) ([788366a](https://github.com/keptn/keptn/commit/788366aa1af7ba4414d4fced30cb3bbe0f7b3080))


### Other

* Add [@heinzburgstaller](https://github.com/heinzburgstaller) as member ([#7847](https://github.com/keptn/keptn/issues/7847)) ([e3ac5fc](https://github.com/keptn/keptn/commit/e3ac5fcb3f36cedc0fa4ebbbb28e11510dc100a9))
* Add New Integration and Keptn Slack to the new issue screen ([#7669](https://github.com/keptn/keptn/issues/7669)) ([48ba7aa](https://github.com/keptn/keptn/commit/48ba7aaaedb19d9426d961b4f7f4c02067ae5ea6))
* Added cp-common to keptn repo and to pipeline ([#7814](https://github.com/keptn/keptn/issues/7814)) ([05ef470](https://github.com/keptn/keptn/commit/05ef470de8840c1de94e3c620a13db6e86029626))
* **bridge:** Only update sequence metadata when needed ([#7733](https://github.com/keptn/keptn/issues/7733)) ([e2473ec](https://github.com/keptn/keptn/commit/e2473ec7993b4e23384e972076088ee16f81d836))
* **bridge:** Remove dev-dependency jest-fetch-mock ([#7703](https://github.com/keptn/keptn/issues/7703)) ([d130add](https://github.com/keptn/keptn/commit/d130adddcc742df80ec423d30b4a33a506eca013))
* **bridge:** Remove sequence metadata polling ([#7870](https://github.com/keptn/keptn/issues/7870)) ([91360bc](https://github.com/keptn/keptn/commit/91360bca828f6b8ce7ba98fdcea161a268c402a6))
* **bridge:** Upgrade to Angular 12 ([#7724](https://github.com/keptn/keptn/issues/7724)) ([34434be](https://github.com/keptn/keptn/commit/34434be3737c61a232b137bc711c630bc09e54ee))
* Bump go-sdk version ([#7931](https://github.com/keptn/keptn/issues/7931)) ([f9cc0e7](https://github.com/keptn/keptn/commit/f9cc0e7b7613a2ecb2d95421defdb7ba2393b990))
* **cli:** Clean up auth messages ([#7911](https://github.com/keptn/keptn/issues/7911)) ([4d013cc](https://github.com/keptn/keptn/commit/4d013cc9ca43681b9128aa065e03e6048d951aeb))
* **cp-connector:** Remove unnecesarry logs ([#7966](https://github.com/keptn/keptn/issues/7966)) ([92d5991](https://github.com/keptn/keptn/commit/92d59919c8786f4e7445ed17a3712265d6ac90b2))
* Executed swag fmt to format swag annotations ([#7871](https://github.com/keptn/keptn/issues/7871)) ([8a7c809](https://github.com/keptn/keptn/commit/8a7c8093a6990addbeac5cac5ea624f7027f3bb2))
* **installer:** Adapted default values for preStop hook times and resource-service resource limits ([#7894](https://github.com/keptn/keptn/issues/7894)) ([51000d2](https://github.com/keptn/keptn/commit/51000d256ea7d2770310726b49c7d126f81e9afa))
* Reduce execution time of shipyard-controller tests ([#7929](https://github.com/keptn/keptn/issues/7929)) ([3562e44](https://github.com/keptn/keptn/commit/3562e44ca8359e1ce77137480dfbeef5b2b02db8))
* **shipyard-controller:** Improve logs with ctx of blocking sequence ([#7948](https://github.com/keptn/keptn/issues/7948)) ([6cc9544](https://github.com/keptn/keptn/commit/6cc9544f7a456b571cedb9d3061984e3a7ef89af))
* Update references to cp-common ([#7823](https://github.com/keptn/keptn/issues/7823)) ([a07259f](https://github.com/keptn/keptn/commit/a07259f33b53c38984a707bd1519f3dbbe36f8ea))
* updated refs to go-sdk ([#7811](https://github.com/keptn/keptn/issues/7811)) ([5a03c55](https://github.com/keptn/keptn/commit/5a03c555f674138af9c39ad938817a87de505290))

## [0.15.0](https://github.com/keptn/keptn/compare/0.14.0...0.15.0) (2022-05-06)


### ⚠ BREAKING CHANGES

* **cli:** The deprecated `send event new-artifact` command has been removed from the CLI
 * Update go-utils to a version without GitCommit in the finished events

### Features

* add `datadog` flag to `keptn configure monitoring` ([#7285](https://github.com/keptn/keptn/issues/7285)) ([bfcb352](https://github.com/keptn/keptn/commit/bfcb3524d5d0a6e32085196dca4458d5e1fef1f0))
* **bridge:** Configure Git upstream via SSH/HTTPS ([#7330](https://github.com/keptn/keptn/issues/7330)) ([0aaeded](https://github.com/keptn/keptn/commit/0aaededa6057f09e28dd4f6b0b90e2f9cb3dcec5))
* **bridge:** Consider real waiting state for sequences ([#7399](https://github.com/keptn/keptn/issues/7399)) ([f8a5bf0](https://github.com/keptn/keptn/commit/f8a5bf0cb2157d838155aaa2ed9fbfa136cb59e5))
* **bridge:** Create styled loading indicator component ([3c91f7d](https://github.com/keptn/keptn/commit/3c91f7d4c889aca5a4902f3fa9256cd2c4ce3f24)), closes [#5568](https://github.com/keptn/keptn/issues/5568)
* **bridge:** Custom sequence depends on selected stage ([#7463](https://github.com/keptn/keptn/issues/7463)) ([577b8f1](https://github.com/keptn/keptn/commit/577b8f1c31ab96051b06ad253dc891a581956ba7))
* **bridge:** Format trigger sequence date with `YYYY-MM-DD HH:mm:ss` ([#7522](https://github.com/keptn/keptn/issues/7522)) ([096c7a1](https://github.com/keptn/keptn/commit/096c7a161c93924371a9e82af27e9facf9263617))
* **bridge:** OAUTH error handling polished ([#7397](https://github.com/keptn/keptn/issues/7397)) ([0b89a37](https://github.com/keptn/keptn/commit/0b89a37de7996fcac30bc69de717a8b0e7bea13a))
* **bridge:** Open project in new tab ([#7629](https://github.com/keptn/keptn/issues/7629)) ([ba102d5](https://github.com/keptn/keptn/commit/ba102d551e219c6de23a29ac6922c6c3eab0fa9d))
* **bridge:** Show service and time stamp on sequence details page ([#7283](https://github.com/keptn/keptn/issues/7283)) ([d03ab0c](https://github.com/keptn/keptn/commit/d03ab0c78426d201c0df3bb3769c5bb598cb52ad))
* **bridge:** Stop event propergation when clicking on external link ([#7632](https://github.com/keptn/keptn/issues/7632)) ([e93ba8b](https://github.com/keptn/keptn/commit/e93ba8b31e6c5edf5141a54ea5c492a60cfe25cc))
* **bridge:** Unify loading indicators ([#5568](https://github.com/keptn/keptn/issues/5568)) ([#7527](https://github.com/keptn/keptn/issues/7527)) ([b90ac83](https://github.com/keptn/keptn/commit/b90ac831ba39e410325e2482e7ef6f071d6b5ac2))
* Configure terminationGracePeriod, preStop hooks and upgrade strategy for core deployments ([#7466](https://github.com/keptn/keptn/issues/7466)) ([44dbbe1](https://github.com/keptn/keptn/commit/44dbbe17f2a14a8f779eb0463972761b7c77d920))
* **cp-connector:** Added `FixedSubscriptionSource` ([#7528](https://github.com/keptn/keptn/issues/7528)) ([1bfaa27](https://github.com/keptn/keptn/commit/1bfaa2752f62a42351ee94940a02447ee3a590ab))
* **cp-connector:** Forward subscription id to event receiver ([#7655](https://github.com/keptn/keptn/issues/7655)) ([b88db17](https://github.com/keptn/keptn/commit/b88db17d2ead797b42aca2a5b50b8b2ada9bebce))
* Initial implementation of `cp-connector` library ([#7418](https://github.com/keptn/keptn/issues/7418)) ([367e859](https://github.com/keptn/keptn/commit/367e8592633262268c7a7096e7bbf778e5918595))
* **installer:** Add option to define nodeSelector globally or per service ([#7336](https://github.com/keptn/keptn/issues/7336)) ([8b257fa](https://github.com/keptn/keptn/commit/8b257fa56a36cf970a69723f9c3a51c2bcbe4436))
* **installer:** Create separate helm chart for commonly used functionality ([#7568](https://github.com/keptn/keptn/issues/7568)) ([8c93343](https://github.com/keptn/keptn/commit/8c9334390a39b02076b07eeff64b75970e8483f5))
* Introducing ZeroDowntime tests ([#7479](https://github.com/keptn/keptn/issues/7479)) ([71d2c94](https://github.com/keptn/keptn/commit/71d2c94c36d24bbccdac953733774d69c3362f4f))
* **secret-service:** provide HTTP 400 when scope is not found ([#7325](https://github.com/keptn/keptn/issues/7325)) ([8cf10b6](https://github.com/keptn/keptn/commit/8cf10b69731f094fa131ae7d2d5e00e7082ee261))
* **shipyard-controller:** Introduce automatic provisioning of gitRemoteURI ([#7276](https://github.com/keptn/keptn/issues/7276)) ([59778e0](https://github.com/keptn/keptn/commit/59778e0cfe61ba63e040c0ce4f7fceaa856e2d24))
* **shipyard-controller:** stop pulling messages after receiving sigterm ([#7464](https://github.com/keptn/keptn/issues/7464)) ([f04874a](https://github.com/keptn/keptn/commit/f04874a6ecf2cd6b9ecb53400da1a53cd5ee5b02))
* **shipyard-controller:** Store sequence executions in new format without potential dots (.) in property names ([#7605](https://github.com/keptn/keptn/issues/7605)) ([1bc93f3](https://github.com/keptn/keptn/commit/1bc93f339b43f82c1735d59041f9358837f93ae5))
* **webhook-service:** Implement v1beta1 webhook config version support ([#7329](https://github.com/keptn/keptn/issues/7329)) ([56c082f](https://github.com/keptn/keptn/commit/56c082fa971eda89b4bc826b4d014e4aa5c049f0))
* **webhook-service:** Introduce keptn-webhook-config ConfigMap with denyList ([#7548](https://github.com/keptn/keptn/issues/7548)) ([b392dc0](https://github.com/keptn/keptn/commit/b392dc025a893d69e87dd7ccf209d5ffe93fbb92))


### Bug Fixes

*  Added correct error and test to jmeter exec ([#7377](https://github.com/keptn/keptn/issues/7377)) ([f689877](https://github.com/keptn/keptn/commit/f68987703d3ab7b3a9a6e821f800cf631e9d0826))
*  Resource-service clean-up  ([#7427](https://github.com/keptn/keptn/issues/7427)) ([0e75970](https://github.com/keptn/keptn/commit/0e7597043d35c0f0f9d11f6179a8dec732b1a026))
* add support for ingress class name ([#7324](https://github.com/keptn/keptn/issues/7324)) ([2fe45a8](https://github.com/keptn/keptn/commit/2fe45a872e6247a1703bd270ac503c0f763350dd))
* Added default user string ([#7430](https://github.com/keptn/keptn/issues/7430)) ([3b8f1ca](https://github.com/keptn/keptn/commit/3b8f1caed9dcdb49e40007cf9fc604bb76ce1ce7))
* Added missing UpdateProject parameters ([#7362](https://github.com/keptn/keptn/issues/7362)) ([ae5b81c](https://github.com/keptn/keptn/commit/ae5b81c82e55de2f4c92493ac0ab068b10ea1ce1))
* Added validation of uniform subscriptions ([#7366](https://github.com/keptn/keptn/issues/7366)) ([c9670c7](https://github.com/keptn/keptn/commit/c9670c716508d39f31976cbd474e283fe045e10b))
* **api:** Allow to enable/disable rate limit ([#7534](https://github.com/keptn/keptn/issues/7534)) ([b36816c](https://github.com/keptn/keptn/commit/b36816ce83773fc804517c2e3540a7e67a63b85a))
* **api:** Metadata model update ([#7349](https://github.com/keptn/keptn/issues/7349)) ([f93c920](https://github.com/keptn/keptn/commit/f93c92031bc4c5a8c16b72f0ab8a565ea88602e7))
* **bridge:** Copy to clipboard button rendering ([#7571](https://github.com/keptn/keptn/issues/7571)) ([f2f236f](https://github.com/keptn/keptn/commit/f2f236fe963f1d1d664adc69d26c6ac932684ef2))
* **bridge:** Do not send a start date for evaluation if none is given by the user ([43f053c](https://github.com/keptn/keptn/commit/43f053c8327f433ffcb0475cd740415df9fd9c3a))
* **bridge:** Fix update of git upstream without a user ([#7519](https://github.com/keptn/keptn/issues/7519)) ([4a05795](https://github.com/keptn/keptn/commit/4a05795acd224911a9c695893e9e3b7f0d5784e2))
* **bridge:** Fixed incorrect selected stage in sequence timeline ([#7394](https://github.com/keptn/keptn/issues/7394)) ([558e491](https://github.com/keptn/keptn/commit/558e4914f936f377a5931d1f18c0f63609571e1a))
* **bridge:** Pretty-print request errors ([#7652](https://github.com/keptn/keptn/issues/7652)) ([5b395b9](https://github.com/keptn/keptn/commit/5b395b97595bcc026a437671773a67b28041ecdc))
* **bridge:** Render html in notifications ([#7523](https://github.com/keptn/keptn/issues/7523)) ([5ae2853](https://github.com/keptn/keptn/commit/5ae2853f3a728d5233e22b9715819ea0be9cc9a9))
* **bridge:** Show remediation sequence in default color while running ([#7300](https://github.com/keptn/keptn/issues/7300)) ([6cf6f6b](https://github.com/keptn/keptn/commit/6cf6f6b9fa546c9f4d7b45d7c0a5b3acb6b7cd14))
* **bridge:** Subscription filter now correctly updates on delete/create service ([#7480](https://github.com/keptn/keptn/issues/7480)) ([fc7d3b4](https://github.com/keptn/keptn/commit/fc7d3b4390546746bba2f14bd51bde7aa7e9c20a))
* Changed help messages in labels ([#7491](https://github.com/keptn/keptn/issues/7491)) ([0a2ca97](https://github.com/keptn/keptn/commit/0a2ca97b982cedd781e8ca203b2fa4196b6adcd6))
* **cli:** Cleaned up Oauth command ([#7307](https://github.com/keptn/keptn/issues/7307)) ([c4c9cd1](https://github.com/keptn/keptn/commit/c4c9cd1a9b7046530596de1869cbacdbc66ac18a))
* **cli:** Provide values needed for upgrading the nats dependency ([#7316](https://github.com/keptn/keptn/issues/7316)) ([#7321](https://github.com/keptn/keptn/issues/7321)) ([8962936](https://github.com/keptn/keptn/commit/89629360f4b54300fa923b99d0ad58b8dcaa45f1))
* **cli:** Remove --tag option from trigger delivery command, remove deprecated new-artifact command ([#7376](https://github.com/keptn/keptn/issues/7376)) ([787f08b](https://github.com/keptn/keptn/commit/787f08ba1f6fa3897eb9582c7655fa270ac947d2))
* Disconnect MongoDB client before reconnecting ([#7416](https://github.com/keptn/keptn/issues/7416)) ([a90d39c](https://github.com/keptn/keptn/commit/a90d39c33ddd248f4c19fc3713ab50121b5763d1))
* **distributor:** Parsing of url forces scheme to http or https ([#7641](https://github.com/keptn/keptn/issues/7641)) ([9240659](https://github.com/keptn/keptn/commit/9240659031ec117bf481cee7543742e95ffd48b3))
* Do not require git user being set when updating project upstream credentials ([#7533](https://github.com/keptn/keptn/issues/7533)) ([ccbf2f1](https://github.com/keptn/keptn/commit/ccbf2f179564741dcd41022fd5ea9840171c4cf8))
* **installer:** Custom readiness probe for MongoDB to fix default one ([#7663](https://github.com/keptn/keptn/issues/7663)) ([0c8b879](https://github.com/keptn/keptn/commit/0c8b87950aa15b3c89c037d8664d6d4846375b45))
* **installer:** Quote value of MAX_AUTH_ENABLED ([#7549](https://github.com/keptn/keptn/issues/7549)) ([b3a4cb9](https://github.com/keptn/keptn/commit/b3a4cb9270eae64ca149bd5fc9e267436d26c75a))
* **installer:** Revert configuration-service back to update strategy recreate ([#7650](https://github.com/keptn/keptn/issues/7650)) ([c4ab18d](https://github.com/keptn/keptn/commit/c4ab18d941600e592d26e75989d6298a30705ccb))
* **jmeter-service:** Avoid nil pointer access when logging results ([#7391](https://github.com/keptn/keptn/issues/7391)) ([c981022](https://github.com/keptn/keptn/commit/c981022228bf35641fc3722c06e54ceb810a7486))
* Rename GitProxyInsecure to InsecureSkpTLS and pass it to upstream interactions ([#7410](https://github.com/keptn/keptn/issues/7410)) ([07d2ad9](https://github.com/keptn/keptn/commit/07d2ad909eb88641ebb3adfe66ede38dec67a902))
* **resource-service:** Fixed unit test ([#7443](https://github.com/keptn/keptn/issues/7443)) ([8f6dbb5](https://github.com/keptn/keptn/commit/8f6dbb5e3274b9f891a4aaab9cb43f39433d12c2))
* **shipyard-controller:** Added option to configure maximum service name length, adapted returned http status code ([#7445](https://github.com/keptn/keptn/issues/7445)) ([26bc02a](https://github.com/keptn/keptn/commit/26bc02ab7016f8d40153e8849115fb4ef05c99a3))
* **shipyard-controller:** Fix order of merging properties for event payload ([#7631](https://github.com/keptn/keptn/issues/7631)) ([#7651](https://github.com/keptn/keptn/issues/7651)) ([640b80e](https://github.com/keptn/keptn/commit/640b80e9e499722ad3f3db845950032d94ac7fa5))
* **shipyard-controller:** Proceed with service deletion if the service is not present on the configuration service anymore ([#7461](https://github.com/keptn/keptn/issues/7461)) ([6ee9f48](https://github.com/keptn/keptn/commit/6ee9f4851ba498d8948e60d006bd7e6459802154))
* URL-provisioning test should wait for pod restart([#7411](https://github.com/keptn/keptn/issues/7411)) ([966a549](https://github.com/keptn/keptn/commit/966a549600d6c8a4f0f50ddca5e515014d3d4b00))


### Refactoring

* **bridge:** Move static server pages to client ([#7369](https://github.com/keptn/keptn/issues/7369)) ([0ff21f3](https://github.com/keptn/keptn/commit/0ff21f3a335379f32afa3b6bc715e574f3ec886d))


### Other

* Add [@renepanzar](https://github.com/renepanzar) as member ([#7612](https://github.com/keptn/keptn/issues/7612)) ([a99e889](https://github.com/keptn/keptn/commit/a99e8890095bb3bb6422c3e3cfd6e953b9449ef6))
* **cli:** Polish upgrade message when no upstream is present ([#7310](https://github.com/keptn/keptn/issues/7310)) ([bdda191](https://github.com/keptn/keptn/commit/bdda1917ca758ef7cf93b08eb1bfc276e2c5faed))
* **installer:** Upgrade MongoDB to v11 ([#7444](https://github.com/keptn/keptn/issues/7444)) ([9346d41](https://github.com/keptn/keptn/commit/9346d41f851300bf308fcc8fe1112ee875924506))
* Make filter a mandatory field in mongo datastore get event by type ([#7355](https://github.com/keptn/keptn/issues/7355)) ([117f904](https://github.com/keptn/keptn/commit/117f904ccb1d270e9cc093b5a346b30803c0892c))
* Updated go-utils to version removing gitcommit from finished events ([#7320](https://github.com/keptn/keptn/issues/7320)) ([c241059](https://github.com/keptn/keptn/commit/c24105911e36b1c9695b5b424ab66740db586bc9))


### Docs

* Add conventions for logging and env var naming ([#7611](https://github.com/keptn/keptn/issues/7611)) ([90f8536](https://github.com/keptn/keptn/commit/90f8536f8b38b667b88cbe12600270fa9e8c44a1))
* **cli:** Add missing/remove unsupported commands from README ([#7544](https://github.com/keptn/keptn/issues/7544)) ([bea81f1](https://github.com/keptn/keptn/commit/bea81f1dcb76e93411f59ee63991b954d83991c8))
* **distributor:** Fixed broken link to cloud events docs ([#7441](https://github.com/keptn/keptn/issues/7441)) ([5ee6f28](https://github.com/keptn/keptn/commit/5ee6f28ff8ccd6aabc405e0405115eab2235a4f9))
* Fix hyperlink to references to docs folder ([#7327](https://github.com/keptn/keptn/issues/7327)) ([5d8b4eb](https://github.com/keptn/keptn/commit/5d8b4eb711b479d1195ee059f790368d3d4e0507))

## [0.14.0](https://github.com/keptn/keptn/compare/0.13.0...0.14.0) (2022-03-30)


### ⚠ BREAKING CHANGES

* **cli:** The CLI does not require anymore passing git_user as a parameter to create or upgrade a project. In case you are experiencing issues with the command, we suggest trying it without specifying the user.
* **configuration-service:** adding invalid token results now to 404 error code (424 was used previously) 

* fix: Return 404 when token is invalid
 
### Features

* Add prometheus-service scope to secret-service ([#6938](https://github.com/keptn/keptn/issues/6938)) ([b2993f2](https://github.com/keptn/keptn/commit/b2993f223444dca7722b204a9d2307ebdb081195))
* Add SSH publicKey auth support ([#6855](https://github.com/keptn/keptn/issues/6855)) ([b1b3d11](https://github.com/keptn/keptn/commit/b1b3d11c6d0ed6dea1016b0757ce4a1d0bddbc85))
* **api:** Added automaticProvisioning helm value to api-service ([#7269](https://github.com/keptn/keptn/issues/7269)) ([0bda1c7](https://github.com/keptn/keptn/commit/0bda1c78c4f6c553109177bbc2b87e088c5dd23f))
* **bridge:** Allow to configure sendStarted flag for webhook config ([#7183](https://github.com/keptn/keptn/issues/7183)) ([7117204](https://github.com/keptn/keptn/commit/7117204fffeab07af57cdbc6b881763057bf6ff5))
* **bridge:** Make secret selection dynamic ([#6940](https://github.com/keptn/keptn/issues/6940)) ([be8394d](https://github.com/keptn/keptn/commit/be8394de2f7bc7a9d5abc5b47375e7f76ce85378))
* **bridge:** Show history of quality gates in environment details ([#7009](https://github.com/keptn/keptn/issues/7009)) ([d1b96ef](https://github.com/keptn/keptn/commit/d1b96ef72ed369e71fcca90290d869ea803193a7))
* **bridge:** Trigger a sequence from Bridge ([#4507](https://github.com/keptn/keptn/issues/4507)) ([84322f3](https://github.com/keptn/keptn/commit/84322f37e19e50e96757f35643287e81530b1b13))
* **bridge:** Use new evaluation finished payload and UI adoptions in SLI breakdown ([#6813](https://github.com/keptn/keptn/issues/6813)) ([711b845](https://github.com/keptn/keptn/commit/711b84512ab47fd3b6e9f066afadb8b92da0b462))
* **cli:** Added keptn trigger `sequence` cmd ([#7070](https://github.com/keptn/keptn/issues/7070)) ([80f2f7d](https://github.com/keptn/keptn/commit/80f2f7d1e2f4fbac9af222fe546e927baf5ce691))
* **cli:** trigger authorization code flow when refresh token is expired ([#7014](https://github.com/keptn/keptn/issues/7014)) ([d596efb](https://github.com/keptn/keptn/commit/d596efbe44b4fbfa182797705a91a293b88ad1fe))
* **distributor:** Added preamble to distributor logs ([#7296](https://github.com/keptn/keptn/issues/7296)) ([1413ad6](https://github.com/keptn/keptn/commit/1413ad6d7f3b5deb93d7f846ec055bea92fa3cd1))
* Improve unallowed URLs of webhook-service ([#7147](https://github.com/keptn/keptn/issues/7147)) ([d5c1d3c](https://github.com/keptn/keptn/commit/d5c1d3c8ab2573719ad8ba275cfce11b61d3c2ba))
* **resource-service:** Delete project via cloud events ([#7024](https://github.com/keptn/keptn/issues/7024)) ([86b0cb9](https://github.com/keptn/keptn/commit/86b0cb940e69b6cb70019ae6f8538c6ef4499c1b))
* **shipyard-controller:** Added leader election ([#6967](https://github.com/keptn/keptn/issues/6967)) ([c5264bd](https://github.com/keptn/keptn/commit/c5264bd67ba52b65affed9cc8029daa45cfdb10f))
* **shipyard-controller:** Introduce new data model ([#6977](https://github.com/keptn/keptn/issues/6977)) ([f46905a](https://github.com/keptn/keptn/commit/f46905ad97ba5d566737e5703a7a5593b0c2fe1b))
* **shipyard-controller:** Subscribe to events using Jetstream ([#6834](https://github.com/keptn/keptn/issues/6834)) ([753547b](https://github.com/keptn/keptn/commit/753547b592dfd588a51aed939c1e6a5a1d11df43))
* Support --git-commit-id flag in CLI trigger evaluation ([#6956](https://github.com/keptn/keptn/issues/6956)) ([f98155c](https://github.com/keptn/keptn/commit/f98155c54c8732a5caf408ccd12b8c14ed4f2cde))
* Support auth via proxy ([#6984](https://github.com/keptn/keptn/issues/6984)) ([63fca54](https://github.com/keptn/keptn/commit/63fca54f18379d98dba21ed2d5121dc23bb82f05))


### Bug Fixes

* add default helm value for project name max size ([#7289](https://github.com/keptn/keptn/issues/7289)) ([1b016a1](https://github.com/keptn/keptn/commit/1b016a164e2b5ed812b19ff88896c2395fa7d05c))
* Backup git-credentials when using resource-service in integration tests ([#7111](https://github.com/keptn/keptn/issues/7111)) ([cafab72](https://github.com/keptn/keptn/commit/cafab722d95da8579960ac46d85362afdf6e9f76))
* **bridge:** Add latestEvaluationTrace to every stage ([8048020](https://github.com/keptn/keptn/commit/8048020f7f5387c255e6fbcb25f61a1851f12c60))
* **bridge:** Break words in project tile, to keep fix width ([#7214](https://github.com/keptn/keptn/issues/7214)) ([3227f8a](https://github.com/keptn/keptn/commit/3227f8aa02861383d9e9e5163fbc2fd71660dafa))
* **bridge:** Fix duplicated sequence and incorrect show older sequences ([#7054](https://github.com/keptn/keptn/issues/7054)) ([95c5bdc](https://github.com/keptn/keptn/commit/95c5bdc300dd6d3112578c205ad61cead8d1da6c))
* **bridge:** Fix incorrect content security policy ([e575943](https://github.com/keptn/keptn/commit/e5759437196cc189edce635762e1d616812f2d3e))
* **bridge:** Fix no-services message and link ([#7035](https://github.com/keptn/keptn/issues/7035)) ([c9e58a7](https://github.com/keptn/keptn/commit/c9e58a7df8091276c1323250d8911faa0f062388))
* **bridge:** Fix quick filter overflow ([#7077](https://github.com/keptn/keptn/issues/7077)) ([2dff06a](https://github.com/keptn/keptn/commit/2dff06afaba6ea440c4432a69de10da8ea8ea3e9))
* **bridge:** Fix wrong time in sequence timeline ([#7036](https://github.com/keptn/keptn/issues/7036)) ([76811ec](https://github.com/keptn/keptn/commit/76811ece751193cf62dd9d8f38d541771a677b40))
* **bridge:** load projects, also if version.json could not be loaded ([#7241](https://github.com/keptn/keptn/issues/7241)) ([50acc3a](https://github.com/keptn/keptn/commit/50acc3ace3058b3716cf0cdd8b98a420fc6f682c))
* **bridge:** Prevent spaces in URL ([#7023](https://github.com/keptn/keptn/issues/7023)) ([0d01639](https://github.com/keptn/keptn/commit/0d016390bf3f88f6e93493f50f9828ce8d463f79))
* **bridge:** Remove line breaks and unnecessary escaping in strings in webhook.yaml ([#7025](https://github.com/keptn/keptn/issues/7025)) ([23ac339](https://github.com/keptn/keptn/commit/23ac339e9b0a42d72f50d613e1fc42499f98bc99))
* **bridge:** Rounding evaluation score correctly ([#6976](https://github.com/keptn/keptn/issues/6976)) ([5b89a91](https://github.com/keptn/keptn/commit/5b89a916b5542af2e21b016edffa4147a3a90d68))
* **bridge:** Truncate evaluation score ([#6993](https://github.com/keptn/keptn/issues/6993)) ([df8e03a](https://github.com/keptn/keptn/commit/df8e03a4cef074595940be83fb2c8818d8cb29ce))
* **bridge:** Validate start end date duration ([0596eae](https://github.com/keptn/keptn/commit/0596eaec6e5beb363ba3e122af60eea2b45d0456))
* **cli:** Added missing env variables to tests ([#6867](https://github.com/keptn/keptn/issues/6867)) ([33feef1](https://github.com/keptn/keptn/commit/33feef190a54d6c8414d897ddf9604af6b912034))
* **cli:** Fixed parsing of image option in trigger delivery ([#7302](https://github.com/keptn/keptn/issues/7302)) ([171a979](https://github.com/keptn/keptn/commit/171a979e5f510c25f0d17ae8f0f81824c9c93dc9))
* **cli:** Removed user check from create/update project and added simple tests ([#7193](https://github.com/keptn/keptn/issues/7193)) ([2b490d5](https://github.com/keptn/keptn/commit/2b490d597e4718b76954d0a1b0179148bcaddb64))
* **configuration-service:** Return 404 when token is invalid ([#7121](https://github.com/keptn/keptn/issues/7121)) ([6805da2](https://github.com/keptn/keptn/commit/6805da214c6c620ffab5edbbd152681c24c9dd6c))
* correct passing of projectNameMaxSize helm value with quotes ([#7288](https://github.com/keptn/keptn/issues/7288)) ([517e2a2](https://github.com/keptn/keptn/commit/517e2a2b74d7bd67320a5aae999b8582daf5294d))
* **distributor:** Added longer sleep for Nats down test in forwarder ([#7205](https://github.com/keptn/keptn/issues/7205)) ([3fff36d](https://github.com/keptn/keptn/commit/3fff36dd8ddaa0d7fd6d27f9a90e7ec9c2fff27c))
* **distributor:** Fixed reconnection issue of (re)used ce clients ([#7109](https://github.com/keptn/keptn/issues/7109)) ([9b69d64](https://github.com/keptn/keptn/commit/9b69d648055b6131a3cc49e7655b4fbfc8e61659))
* **distributor:** Include event filter for project, stage, service ([#6968](https://github.com/keptn/keptn/issues/6968)) ([#6972](https://github.com/keptn/keptn/issues/6972)) ([6ab050d](https://github.com/keptn/keptn/commit/6ab050d6bbc37a02ea8506d6c3fc5dd2472805c0))
* **distributor:** Increase timout of http client to 30s ([#6948](https://github.com/keptn/keptn/issues/6948)) ([#6954](https://github.com/keptn/keptn/issues/6954)) ([3ccbd77](https://github.com/keptn/keptn/commit/3ccbd77f32d95bf3540817a9d59f89591e88a3fb))
* **distributor:** shut down distributor when not able to send heartbeat to control plane ([#7263](https://github.com/keptn/keptn/issues/7263)) ([7c50feb](https://github.com/keptn/keptn/commit/7c50feb198a95d8663bd2dfa4bb7f6a839237011))
* ensure indicators are set in computeObjectives ([#6922](https://github.com/keptn/keptn/issues/6922)) ([b1cc56d](https://github.com/keptn/keptn/commit/b1cc56d543982212772acf32ef4ca398943822a0))
* Forbid project names longer than a certain size ([#7277](https://github.com/keptn/keptn/issues/7277)) ([237c4cf](https://github.com/keptn/keptn/commit/237c4cf2e32567e928ddd18c9ac29574c09df6b9))
* hardening of oauth in distributor and cli ([#6917](https://github.com/keptn/keptn/issues/6917)) ([b73a379](https://github.com/keptn/keptn/commit/b73a3798aa393edef7d17b6b577683415ca3bfae))
* **helm-service:** Handling of helm-charts loading problems ([#7108](https://github.com/keptn/keptn/issues/7108)) ([3a60e50](https://github.com/keptn/keptn/commit/3a60e50d2bb35f6ef704e8335c2a329012150cd9))
* **installer:** Make securityContext configurable ([#6932](https://github.com/keptn/keptn/issues/6932)) ([#6949](https://github.com/keptn/keptn/issues/6949)) ([b711b0a](https://github.com/keptn/keptn/commit/b711b0a1b1fa4d137eb9177015726de8f134e128))
* **lighthouse-service:** Fail sequence when evaluation is aborted/errored ([#7211](https://github.com/keptn/keptn/issues/7211)) ([1faca09](https://github.com/keptn/keptn/commit/1faca099c982b4536748d8559ef438f664d0d056))
* Normalize error messages ([#7080](https://github.com/keptn/keptn/issues/7080)) ([0730f1d](https://github.com/keptn/keptn/commit/0730f1d1cb33bf604893b55aba5922365b50455d))
* **resource-service:** fix nats subscription and added retry logic ([#7215](https://github.com/keptn/keptn/issues/7215)) ([180d833](https://github.com/keptn/keptn/commit/180d833bcbc3cdd35f3d71694a653d9550e9ce9e))
* **resource-service:** Make sure to delete "/" prefix in resourcePath when resolving git commit id ([#6919](https://github.com/keptn/keptn/issues/6919)) ([2ae4c52](https://github.com/keptn/keptn/commit/2ae4c5223a59f774c040b65f2fd38df2cc3756f4))
* **shipyard-controller:** Abort multi-stage sequences ([#7175](https://github.com/keptn/keptn/issues/7175)) ([d06aefb](https://github.com/keptn/keptn/commit/d06aefb519108436840be23b566a27046345ea72))
* **shipyard-controller:** Consider parallel stages when trying to set overall sequence state to finished ([#7250](https://github.com/keptn/keptn/issues/7250)) ([9550f59](https://github.com/keptn/keptn/commit/9550f5986e20ad70b5d5d00bf58dc055462d7fe5))
* **shipyard-controller:** Do not exit pull subscription loop when invalid event has been received ([#7255](https://github.com/keptn/keptn/issues/7255)) ([75c5971](https://github.com/keptn/keptn/commit/75c59716d6a042a923c8b4557dfa2f7f02f39544))
* **shipyard-controller:** Do not reset subscriptions when updating distributor/integration version ([#7046](https://github.com/keptn/keptn/issues/7046)) ([#7059](https://github.com/keptn/keptn/issues/7059)) ([5865cf1](https://github.com/keptn/keptn/commit/5865cf1c3a538c332e5522dce307a578f5dc60fd)), closes [#6598](https://github.com/keptn/keptn/issues/6598) [#6613](https://github.com/keptn/keptn/issues/6613) [#6618](https://github.com/keptn/keptn/issues/6618) [#6619](https://github.com/keptn/keptn/issues/6619) [#6634](https://github.com/keptn/keptn/issues/6634) [#6559](https://github.com/keptn/keptn/issues/6559) [#6642](https://github.com/keptn/keptn/issues/6642) [#6643](https://github.com/keptn/keptn/issues/6643) [#6659](https://github.com/keptn/keptn/issues/6659) [#6670](https://github.com/keptn/keptn/issues/6670) [#6632](https://github.com/keptn/keptn/issues/6632) [#6718](https://github.com/keptn/keptn/issues/6718) [#6816](https://github.com/keptn/keptn/issues/6816) [#6819](https://github.com/keptn/keptn/issues/6819) [#6820](https://github.com/keptn/keptn/issues/6820) [#6875](https://github.com/keptn/keptn/issues/6875) [#6763](https://github.com/keptn/keptn/issues/6763) [#6857](https://github.com/keptn/keptn/issues/6857) [#6804](https://github.com/keptn/keptn/issues/6804) [#6931](https://github.com/keptn/keptn/issues/6931) [#6944](https://github.com/keptn/keptn/issues/6944) [#6966](https://github.com/keptn/keptn/issues/6966) [#6971](https://github.com/keptn/keptn/issues/6971)
* **webhook-service:** Disallow `@` file uploads inside data block ([#7158](https://github.com/keptn/keptn/issues/7158)) ([aa0f71e](https://github.com/keptn/keptn/commit/aa0f71e4fffda8c0959d7e7ef32dd90f6f9914f5))
* **webhook-service:** enhance denylist of disallowed urls ([#7191](https://github.com/keptn/keptn/issues/7191)) ([048dbe4](https://github.com/keptn/keptn/commit/048dbe45685b3b383cea052f42612f37079bd323))
* **webhook-service:** Fix retrieval of webhook config ([#7144](https://github.com/keptn/keptn/issues/7144)) ([08ae798](https://github.com/keptn/keptn/commit/08ae798e5436055e936f60628ca2c3b41fdce341))


### Docs

* **bridge:** Add documentation for environment variables ([0bb45a9](https://github.com/keptn/keptn/commit/0bb45a9475a4d4411e1d0b0ee86ae468a9b03e39))
* Reference the code of conduct in the .github repository ([#7110](https://github.com/keptn/keptn/issues/7110)) ([3dbd75c](https://github.com/keptn/keptn/commit/3dbd75c52f99fb0a7864801021eef037e4aa2342))
* Stop-gap info about filtering by stage, project,service ([#7155](https://github.com/keptn/keptn/issues/7155)) ([ee03d92](https://github.com/keptn/keptn/commit/ee03d9260d55c197d7b7aed7b54b707adedf0b9c))
* Use K3d 5.3.0 in README for developing ([#6926](https://github.com/keptn/keptn/issues/6926)) ([f02cad5](https://github.com/keptn/keptn/commit/f02cad5de1c584621504fdd4b3fe7bf4c19870e2))


### Other

*  Changed all integration tests to use go utils ([#7165](https://github.com/keptn/keptn/issues/7165)) ([d926eb4](https://github.com/keptn/keptn/commit/d926eb429404f892c1628862a3d5ff6bf075d4d8))
* Add [@j-poecher](https://github.com/j-poecher) as member ([#7294](https://github.com/keptn/keptn/issues/7294)) ([979e81d](https://github.com/keptn/keptn/commit/979e81daa1803f2b21069ba12274fe24275968ad))
* Add [@pchila](https://github.com/pchila) as member to maintainers.md ([#6946](https://github.com/keptn/keptn/issues/6946)) ([b919720](https://github.com/keptn/keptn/commit/b9197205ce633f6b0dd277337d72aff1840b1931))
* Add [@raffy23](https://github.com/raffy23) as member ([#7174](https://github.com/keptn/keptn/issues/7174)) ([67fa5a5](https://github.com/keptn/keptn/commit/67fa5a5e4c139e6672c28e88015af62118366593))
* Add Slack issue link ([#7181](https://github.com/keptn/keptn/issues/7181)) ([33bd789](https://github.com/keptn/keptn/commit/33bd7896038684343ec779edd581f55f28d4ac83))
* **bridge:** Remove unused dependencies ([#7012](https://github.com/keptn/keptn/issues/7012)) ([9be7608](https://github.com/keptn/keptn/commit/9be760883d44ed513f391a2fef1fee4bef109659))
* **distributor:** cleanup of package structure ([#7028](https://github.com/keptn/keptn/issues/7028)) ([e97875c](https://github.com/keptn/keptn/commit/e97875cad6031bbe175aaedf865a9ebec1ea2c58))
* **distributor:** hardening of unit test stability ([#6992](https://github.com/keptn/keptn/issues/6992)) ([f4f1365](https://github.com/keptn/keptn/commit/f4f13650de9b5a9efb984bb466b9397bd32bab77))
* **installer:** Cleaned up common labels ([#6796](https://github.com/keptn/keptn/issues/6796)) ([1f6f6dc](https://github.com/keptn/keptn/commit/1f6f6dcb77f2cb5bd548ac7b11ce1b5f74ea4f42))
* **jmeter-service:** Updated Dynatrace JMeter extension to 1.8.0 ([#6879](https://github.com/keptn/keptn/issues/6879)) ([89b2ba1](https://github.com/keptn/keptn/commit/89b2ba170deb10f400295c64103c53dffcb7a452))
* Move Stage API endpoint into the correct subsection ([#6994](https://github.com/keptn/keptn/issues/6994)) ([bac751d](https://github.com/keptn/keptn/commit/bac751d74a3e60f107d9824c341a1bdc2be555f9))
* Removed makefile and all usages of it ([#6804](https://github.com/keptn/keptn/issues/6804)) ([e55355f](https://github.com/keptn/keptn/commit/e55355ffaa7b19c9ce2394400e8405876778c03a))
* Replace the Security guidelines by the hyperlink ([#7145](https://github.com/keptn/keptn/issues/7145)) ([f640e2c](https://github.com/keptn/keptn/commit/f640e2c308ef56e1641507f7059912852018bd73))
* Upgrade to Go 1.17 ([#7095](https://github.com/keptn/keptn/issues/7095)) ([9deafc9](https://github.com/keptn/keptn/commit/9deafc95f694796ebdfc6fc9388878eccea348ca))
* Use correct Keptn branding logo and spelling ([#7240](https://github.com/keptn/keptn/issues/7240)) ([376ce36](https://github.com/keptn/keptn/commit/376ce36952f1e3e632a5ea972417b9648db41563))
* **webhook-service:** added test for being able to use @ char inside payload ([#7166](https://github.com/keptn/keptn/issues/7166)) ([68db33c](https://github.com/keptn/keptn/commit/68db33cd041f55ae82d4f745b482d73c635517e3))
* **webhook-service:** replaced "unallowed" with "denied" ([#7286](https://github.com/keptn/keptn/issues/7286)) ([ac3e52e](https://github.com/keptn/keptn/commit/ac3e52e274f1c61e5119afac3f4b3435cf817214))

## [0.13.0](https://github.com/keptn/keptn/compare/0.12.0...0.13.0) (2022-02-17)


### ⚠ BREAKING CHANGES

* **bridge:** The uniform screen has been moved into the settings screen.
* in keptn sdk the keptn_fake interfaces have been updated to have api.GetOption in their signature (see https://github.com/keptn/go-utils/pull/375/files#diff-245aca76b6ab2043d44c217312e1b9d487545aca0dd53418fb2106efacaec7b3
* The sequence control now supports also a `waiting` state. 
* Several API endpoints have been marked as internal. For more information, please check [#6303](https://github.com/keptn/keptn/issues/6303).

### Features

* Added commitID to webhook and jmeter, updated go-utils dependencies ([#6567](https://github.com/keptn/keptn/issues/6567)) ([#6787](https://github.com/keptn/keptn/issues/6787)) ([5ad04fa](https://github.com/keptn/keptn/commit/5ad04fadad24d06616880b3538a907cb1dcdfd46))
* Added get options to fake keptn in go-sdk ([#6742](https://github.com/keptn/keptn/issues/6742)) ([c6f298c](https://github.com/keptn/keptn/commit/c6f298c8c6ea06ff0b599182738d86fc653a3f9f))
* Block external traffic to internal API endpoints ([#6625](https://github.com/keptn/keptn/issues/6625)) ([7f6a864](https://github.com/keptn/keptn/commit/7f6a8649561f9de2aaa4e00e3b6a384194234944))
* **bridge:** Login via OpenID ([#6076](https://github.com/keptn/keptn/issues/6076)) ([#6077](https://github.com/keptn/keptn/issues/6077)) ([1a657c8](https://github.com/keptn/keptn/commit/1a657c853c8f495ad931dabe014522b86bf919cf))
* **bridge:** Poll sequence metadata for filters and deployed artifacts ([#5246](https://github.com/keptn/keptn/issues/5246)) ([4c5b9df](https://github.com/keptn/keptn/commit/4c5b9dfcb93728920f081fe26078872303e0e1e7))
* **bridge:** Replace memory store with MongoDB ([8d7708f](https://github.com/keptn/keptn/commit/8d7708f736eac08967d801cb59b400b3c2835b94)), closes [#6076](https://github.com/keptn/keptn/issues/6076) [#6688](https://github.com/keptn/keptn/issues/6688) [#6784](https://github.com/keptn/keptn/issues/6784)
* **bridge:** Send access token for each request ([#6078](https://github.com/keptn/keptn/issues/6078)) ([6726306](https://github.com/keptn/keptn/commit/6726306f11a9156d5a095ba86e57d738769eb68a)), closes [#6076](https://github.com/keptn/keptn/issues/6076)
* **bridge:** Show secret scope and keys on overview table ([#6296](https://github.com/keptn/keptn/issues/6296)) ([39fef32](https://github.com/keptn/keptn/commit/39fef32852d20549734d28208d654d857027736c))
* **bridge:** Show specific error message if secret already exists ([#6297](https://github.com/keptn/keptn/issues/6297)) ([fbf7eb0](https://github.com/keptn/keptn/commit/fbf7eb07a1bcaf0342d89ac7e00efc1309b68501))
* **bridge:** Unify notifications ([#5087](https://github.com/keptn/keptn/issues/5087)) ([11941fd](https://github.com/keptn/keptn/commit/11941fdc729871e408d2f9d28c2302c61176a9ed)), closes [#6076](https://github.com/keptn/keptn/issues/6076)
* **cli:** Added `--sso-scopes` flag to cli ([#6708](https://github.com/keptn/keptn/issues/6708)) ([e6e11ba](https://github.com/keptn/keptn/commit/e6e11baf2a62ebe85b899578db9bde2df244893a))
* **cli:** Allow to skip sending the API token when using an SSO integration ([#6675](https://github.com/keptn/keptn/issues/6675)) ([5644e03](https://github.com/keptn/keptn/commit/5644e03e1839228f7b21d3964447ca2297dbeadd))
* **cli:** SSO integration ([#6549](https://github.com/keptn/keptn/issues/6549)) ([2c5f3f7](https://github.com/keptn/keptn/commit/2c5f3f76fa9edee0b73cbdfb0ed07896496ea8ec))
* **cli:** Use `state` param during Oauth flow ([#6701](https://github.com/keptn/keptn/issues/6701)) ([02aecbc](https://github.com/keptn/keptn/commit/02aecbc112fb4db8f60ec55192e10540932f89f6))
* Get and post with commitid ([#6349](https://github.com/keptn/keptn/issues/6349)) ([#6567](https://github.com/keptn/keptn/issues/6567)) ([c3496c0](https://github.com/keptn/keptn/commit/c3496c0f8e3304f916aabebe247b469c618212f5))
* **installer:** Allow API token to be pulled from pre-defined secret ([#6312](https://github.com/keptn/keptn/issues/6312)) ([dc1037a](https://github.com/keptn/keptn/commit/dc1037a1838421210ee363f4bb53822c90e7451c))
* Introduce 'waiting' status to sequences ([#6603](https://github.com/keptn/keptn/issues/6603)) ([e63f312](https://github.com/keptn/keptn/commit/e63f312d1aa4a84beeb4fef2cbecd66933d6c3c9))
* Introduce Oauth integration for distributor and Oauth enhancements for CLI ([#6729](https://github.com/keptn/keptn/issues/6729)) ([7245013](https://github.com/keptn/keptn/commit/7245013d44c45b8785f46da1c131900eae1a53dd))
* Mark endpoints as internal in swagger doc ([#6599](https://github.com/keptn/keptn/issues/6599)) ([3785eed](https://github.com/keptn/keptn/commit/3785eedd1a9581878b70edb7d801fce5c337e7d4))
* **mongodb-datastore:** Use simple join query instead of uncorrelated sub-query ([#6612](https://github.com/keptn/keptn/issues/6612)) ([f57412a](https://github.com/keptn/keptn/commit/f57412a00fccae0b1e293475d464211daa25f388))
* Release helm charts on GitHub pages ([#6559](https://github.com/keptn/keptn/issues/6559)) ([efc285e](https://github.com/keptn/keptn/commit/efc285e65f15b2a6ac6672ffdfe672a5cf4fb7c5))
* **resource-service:** Added support for directory based git model ([#6397](https://github.com/keptn/keptn/issues/6397)) ([#6714](https://github.com/keptn/keptn/issues/6714)) ([ddd5585](https://github.com/keptn/keptn/commit/ddd5585bc78d03073156ef09fab8ba50a871fc24))
* **shipyard-controller:** Propagate git commit ID passed in sequence.triggered events ([#6348](https://github.com/keptn/keptn/issues/6348)) ([#6597](https://github.com/keptn/keptn/issues/6597)) ([ac1f44e](https://github.com/keptn/keptn/commit/ac1f44e648268570da85bc92a7ed73c9e76868c4))
* Update pod config to be more strict w.r.t. security ([#6020](https://github.com/keptn/keptn/issues/6020)) ([6d69563](https://github.com/keptn/keptn/commit/6d6956332ad2259a57b965574c0a411e26bf285e))
* **webhook-service:** Allow disabling .started events ([#6524](https://github.com/keptn/keptn/issues/6524)) ([#6664](https://github.com/keptn/keptn/issues/6664)) ([e07091f](https://github.com/keptn/keptn/commit/e07091f2aa883b1250bbdd66c5618b167f500b30))


### Bug Fixes

* Adapt http status code for not found upstream repositories ([#6641](https://github.com/keptn/keptn/issues/6641)) ([a3ad118](https://github.com/keptn/keptn/commit/a3ad118f4d80ee44addbe39ab11945cd3c8c4548))
* Avoid nil pointer access for undefined value in helm charts ([#6863](https://github.com/keptn/keptn/issues/6863)) ([d845ea6](https://github.com/keptn/keptn/commit/d845ea67ca6df0477629ac3f083795c1a70af4b4))
* **bridge:** Add message that no events are available when sequence has no traces ([#5985](https://github.com/keptn/keptn/issues/5985)) ([64540b9](https://github.com/keptn/keptn/commit/64540b983da60eeadc7b9bc3911129d740b6217c))
* **bridge:** Display additional error information when creating a project ([#6715](https://github.com/keptn/keptn/issues/6715)) ([e8b263f](https://github.com/keptn/keptn/commit/e8b263f08fd9b7f74f64c5bf318f137375023ecc))
* **bridge:** Fix failed sequence shown as succeeded ([#6896](https://github.com/keptn/keptn/issues/6896)) ([e723398](https://github.com/keptn/keptn/commit/e723398a55a7e38779b69a3b52bcdbee6d187548))
* **bridge:** Fix style content security policy ([#6750](https://github.com/keptn/keptn/issues/6750)) ([bd0d569](https://github.com/keptn/keptn/commit/bd0d569f7161dd9ff809a6b44f7e7f8289bfb941))
* **bridge:** Fixed incorrectly shown loading indicator in sequence list ([#6579](https://github.com/keptn/keptn/issues/6579)) ([f238cf4](https://github.com/keptn/keptn/commit/f238cf44e7f29d2c50e43d17cb0d5674f1d50ccf))
* **bridge:** Show error when having problems parsing shipyard.yaml ([#6592](https://github.com/keptn/keptn/issues/6592)) ([#6606](https://github.com/keptn/keptn/issues/6606)) ([0ceb54d](https://github.com/keptn/keptn/commit/0ceb54dfefbd7df2defe1e74e2bcd4c0da0cad91))
* **bridge:** When updating an all events subscription, keep sh.keptn.> format ([#6628](https://github.com/keptn/keptn/issues/6628)) ([1e83fb7](https://github.com/keptn/keptn/commit/1e83fb7967b11264ad7cfb0849d52fd1f4c43a1a))
* **cli:** Added missing command description for `keptn create secret` ([#6621](https://github.com/keptn/keptn/issues/6621)) ([22bddf9](https://github.com/keptn/keptn/commit/22bddf9486f2ad7aeb15e93a085ed3f6371f5820))
* **cli:** Check for unknown subcommands ([#6698](https://github.com/keptn/keptn/issues/6698)) ([c1782c0](https://github.com/keptn/keptn/commit/c1782c01f73ff2a1a4bbe32631f92c1ad19b63bf))
* **cli:** CLI new version checker message ([#6864](https://github.com/keptn/keptn/issues/6864)) ([d836e89](https://github.com/keptn/keptn/commit/d836e890f38c43bad8ad6a2443293448adac69e5))
* **configuration-service:** Adapt to different response from git CLI when upstream branch is not there yet ([#6876](https://github.com/keptn/keptn/issues/6876)) ([#6882](https://github.com/keptn/keptn/issues/6882)) ([c9f0b78](https://github.com/keptn/keptn/commit/c9f0b78063d89d50eac4620d860869664568bd2a))
* **configuration-service:** Ensure that git user and email are set before committing ([#6645](https://github.com/keptn/keptn/issues/6645)) ([#6654](https://github.com/keptn/keptn/issues/6654)) ([d38bb6e](https://github.com/keptn/keptn/commit/d38bb6ef6f0a6eaed95190bcf14fce193c79bee6))
* Fix container image OCI labels ([#6878](https://github.com/keptn/keptn/issues/6878)) ([0f759d4](https://github.com/keptn/keptn/commit/0f759d469c19115b5ca9f507b4b5d50b33f6d688))
* Fixed wrong nginx location for bridge urls ([#6696](https://github.com/keptn/keptn/issues/6696)) ([700895e](https://github.com/keptn/keptn/commit/700895e91b85febbd1ee6c09531a6203aa644a04))
* **installer:** External connection string not used while helm upgrade ([#6760](https://github.com/keptn/keptn/issues/6760)) ([6d04780](https://github.com/keptn/keptn/commit/6d047806f21c2fcd474f6f018c73d3e4bfe47c00))
* **installer:** Fixed helm/jmeter service helm values schema ([#6629](https://github.com/keptn/keptn/issues/6629)) ([085edf1](https://github.com/keptn/keptn/commit/085edf19409ffddceccaec0346090c3cee565873))
* **installer:** Set distributor version in helm chart ([#6652](https://github.com/keptn/keptn/issues/6652)) ([#6658](https://github.com/keptn/keptn/issues/6658)) ([8c2d8de](https://github.com/keptn/keptn/commit/8c2d8dec3d5b80bd59b97fb745cb8db21158067f))
* **jmeter-service:** Finish processes when '.finished' event is sent ([#6786](https://github.com/keptn/keptn/issues/6786)) ([4484a80](https://github.com/keptn/keptn/commit/4484a80b2eb6cb393e013129b0e1fd7c36205163))
* **resource-service:** Fix git-id based file retrieval ([#6616](https://github.com/keptn/keptn/issues/6616)) ([6ba0165](https://github.com/keptn/keptn/commit/6ba01658a2c54b1efa86efd7e86ecee98e4f0a58))
* revert intaller mongoDB version dump ([#6733](https://github.com/keptn/keptn/issues/6733)) ([d96495b](https://github.com/keptn/keptn/commit/d96495bfc5481acd70233e9f7ff0b7c42c01c4f4))
* **shipyard-controller:** Reflect cancellation in sequence state even when no triggered event is there anymore ([#6837](https://github.com/keptn/keptn/issues/6837)) ([bdcd95e](https://github.com/keptn/keptn/commit/bdcd95e5c5857a7c5ef0abff49524cafdf2a8b86))
* Support Openshift 3.11 ([#6578](https://github.com/keptn/keptn/issues/6578)) ([c72dbf2](https://github.com/keptn/keptn/commit/c72dbf2aca410359baa90c52e2cc541ff9ce77f8))
* Trim Incoming Keptn Context and Triggered ID via API ([#6845](https://github.com/keptn/keptn/issues/6845)) ([32d98d9](https://github.com/keptn/keptn/commit/32d98d9ae55a9ad1dd0f61dac20aa56cf865a85a))


### Performance

* Directly return Bridge version ([#6764](https://github.com/keptn/keptn/issues/6764)) ([345469c](https://github.com/keptn/keptn/commit/345469c15106510e786eee1a6e7ce87d7a18840c))


* **bridge:** Restructuring of Bridge settings for project ([75e2842](https://github.com/keptn/keptn/commit/75e284268271070918ec5541997b8db4d6ef1d54))


### Other

* adapted CLI to newest state of APISet in go-utils ([#6655](https://github.com/keptn/keptn/issues/6655)) ([f86774d](https://github.com/keptn/keptn/commit/f86774db1c3b411f3aaf75a73e010cd52a3a4e85))
* Add [@lmmarsano](https://github.com/lmmarsano) as a contributor ([#6046](https://github.com/keptn/keptn/issues/6046)) ([8bfdfd0](https://github.com/keptn/keptn/commit/8bfdfd0b75c3c76890fac905c10192b30c22bea9))
* Add @Im5tu as a contributor ([#6622](https://github.com/keptn/keptn/issues/6622)) ([4dcb4c8](https://github.com/keptn/keptn/commit/4dcb4c8d0c0b62be86952a7490c8b91fe87d263e))
* Add k8s resource stats to release notes ([#6718](https://github.com/keptn/keptn/issues/6718)) ([5ed8ba5](https://github.com/keptn/keptn/commit/5ed8ba50d38661dfbf09b1682623de4dfab22a38))
* adjustments to recent changes in go-utils ([#6706](https://github.com/keptn/keptn/issues/6706)) ([e1f2fd7](https://github.com/keptn/keptn/commit/e1f2fd7ad8e9ca6ddfbe0067fb900f396fd8a6aa))
* **bridge:** Added log for used OAuth scope ([c65fd48](https://github.com/keptn/keptn/commit/c65fd489eac4f73cc6199061b7f609afd45adfc2))
* **bridge:** Remove unused dependency puppeteer ([#6762](https://github.com/keptn/keptn/issues/6762)) ([9224afe](https://github.com/keptn/keptn/commit/9224afe051a45e09240ebe2a748e7b3273cb57b9))
* **installer:** Added metadata to keptn helm chart ([#6624](https://github.com/keptn/keptn/issues/6624)) ([88c3e2b](https://github.com/keptn/keptn/commit/88c3e2bc51b30cd9956aa946eb427610e0cffbac))
* promote [@thschue](https://github.com/thschue) to maintainers ([#6640](https://github.com/keptn/keptn/issues/6640)) ([fb06427](https://github.com/keptn/keptn/commit/fb06427e36ab03371fdd717463d161e2632eb79a))


### Docs

* Add structure for developer documentation ([#6671](https://github.com/keptn/keptn/issues/6671)) ([3fdc8b7](https://github.com/keptn/keptn/commit/3fdc8b78b907f2ecc8d1b8a3146466fd6959d012))
* Updated instructions to install master ([#6889](https://github.com/keptn/keptn/issues/6889)) ([2d4f1be](https://github.com/keptn/keptn/commit/2d4f1be3dc94536dae50be051ee51557100688f9))

## [0.12.0](https://github.com/keptn/keptn/compare/0.11.0...0.12.0) (2022-01-17)


### Features

*  Added problem title in sequence endpoint ([#4206](https://github.com/keptn/keptn/issues/4206)) ([#6258](https://github.com/keptn/keptn/issues/6258)) ([130f3d4](https://github.com/keptn/keptn/commit/130f3d4b3ff886716e4caea790644844f5fd86c4))
* Add helm value for configuration-service version selection ([#6387](https://github.com/keptn/keptn/issues/6387)) ([6e85346](https://github.com/keptn/keptn/commit/6e853469d18e8d6e55964fd60361fb8a030b3b07))
* Add warning for missing upstream ([#6433](https://github.com/keptn/keptn/issues/6433)) ([7a25226](https://github.com/keptn/keptn/commit/7a2522682ce7fd06985a9c53c503e67fa4874d9f))
* Added unit tests for webhook-service and improved error messages ([#6064](https://github.com/keptn/keptn/issues/6064)) ([5b4516e](https://github.com/keptn/keptn/commit/5b4516e7835b02757499c36f946f31ebf3ff5653))
* **bridge:** Add event picker to webhook configuration ([#5260](https://github.com/keptn/keptn/issues/5260)) ([a3c30cc](https://github.com/keptn/keptn/commit/a3c30ccb978d47890dc9d0c7caf706e5cad69f65))
* **bridge:** Add hint for tree list select component and change texts ([#5576](https://github.com/keptn/keptn/issues/5576)) ([0707ff7](https://github.com/keptn/keptn/commit/0707ff7dce5da6fd748985627623a011f3e66a42))
* **bridge:** Add hint that secrets are shared among all projects ([#5832](https://github.com/keptn/keptn/issues/5832)) ([0e40acf](https://github.com/keptn/keptn/commit/0e40acf0ea19ca1912cfcc12cd0008c2b02a527f))
* **bridge:** Add validation for payload if it contains specific characters ([#5950](https://github.com/keptn/keptn/issues/5950)) ([5d5b388](https://github.com/keptn/keptn/commit/5d5b3880742eef9eda49ec7769f14d7e921c397a))
* **bridge:** Consider abort state of sequence ([#6215](https://github.com/keptn/keptn/issues/6215)) ([f86f2f2](https://github.com/keptn/keptn/commit/f86f2f2fa920f7a7dacbcf16fb4f1e6d607d2ba9))
* **bridge:** Feature flags for bridge server ([#6073](https://github.com/keptn/keptn/issues/6073)) ([6351a58](https://github.com/keptn/keptn/commit/6351a58213db91bbe1fff551d8131218af8d961b))
* **bridge:** Introduce keptn color scheme [#5081](https://github.com/keptn/keptn/issues/5081) ([#6577](https://github.com/keptn/keptn/issues/6577)) ([9915466](https://github.com/keptn/keptn/commit/9915466adbb639dd40688791b23e62b86f2469f9))
* **bridge:** Login via OpenID ([#6076](https://github.com/keptn/keptn/issues/6076)) ([#6077](https://github.com/keptn/keptn/issues/6077)) ([8762c83](https://github.com/keptn/keptn/commit/8762c830607c8ff010bb9f3182bf4b29aaf124b5))
* **bridge:** Read problem title from GET /sequence endpoint ([#6388](https://github.com/keptn/keptn/issues/6388)) ([a4086c8](https://github.com/keptn/keptn/commit/a4086c8c0a76d2bbee39898bf11d7eee635f4434)), closes [#5526](https://github.com/keptn/keptn/issues/5526)
* **bridge:** Show dialog to prevent data loss on unsaved form ([#6526](https://github.com/keptn/keptn/issues/6526)) ([c7e7273](https://github.com/keptn/keptn/commit/c7e72732980af6755ee1514d362331d3afad6071))
* **bridge:** Use textarea for webhook url config and adapt styles ([#5706](https://github.com/keptn/keptn/issues/5706)) ([d2a8509](https://github.com/keptn/keptn/commit/d2a850914e2099205403530f4b4118d70b7748ba))
* **cli:** Add missing upstream warning during keptn upgrade ([#6434](https://github.com/keptn/keptn/issues/6434)) ([4867fa5](https://github.com/keptn/keptn/commit/4867fa52cc14f59051b1ef403774feb83373e205))
* **cli:** remove legacy code responsible for shipyard file upgrade from version 0.1.* to 0.2.0 ([#6270](https://github.com/keptn/keptn/issues/6270)) ([8b67626](https://github.com/keptn/keptn/commit/8b6762648791294fc0afd8e56ab3190d1b4994d7))
* Create integration test for resource-service ([#6430](https://github.com/keptn/keptn/issues/6430)) ([220f1ed](https://github.com/keptn/keptn/commit/220f1edbf86fd2a0a22654bcc507d709c51bca9f))
* Enable Nats queue groups in helm chart ([#4519](https://github.com/keptn/keptn/issues/4519)) ([#6062](https://github.com/keptn/keptn/issues/6062)) ([9c493c7](https://github.com/keptn/keptn/commit/9c493c70e0af69c0ce8ec89af82b55ce18af00f3))
* Finalize graceful shutdown ([#4522](https://github.com/keptn/keptn/issues/4522)) ([#6079](https://github.com/keptn/keptn/issues/6079)) ([b5c5d8d](https://github.com/keptn/keptn/commit/b5c5d8da52fa4ba448baf0e5ed9f31fa7460b320))
* graceful shutdown for jmeter and helm service ([#4522](https://github.com/keptn/keptn/issues/4522)) ([#5973](https://github.com/keptn/keptn/issues/5973)) ([41df113](https://github.com/keptn/keptn/commit/41df1138cb2e0bbbc7e32345a99244fc0079dcdc))
* Improve error reporting for CLI trigger cmd ([#6516](https://github.com/keptn/keptn/issues/6516)) ([4904c19](https://github.com/keptn/keptn/commit/4904c19d3829d490910f1241f5ac616fac255932))
* Introduce new-configuration-service ([#6400](https://github.com/keptn/keptn/issues/6400)) ([447f7a0](https://github.com/keptn/keptn/commit/447f7a02cfc735400da24e2044acbcc15542e7c8))
* introduce swappable logger in `go-sdk` ([#6284](https://github.com/keptn/keptn/issues/6284)) ([12d222b](https://github.com/keptn/keptn/commit/12d222b906e8a0899e038419ea550c3da9c6d263))
* **jmeter-service:** Improve error reporting for JMeter-service ([#6511](https://github.com/keptn/keptn/issues/6511)) ([c7d8224](https://github.com/keptn/keptn/commit/c7d8224e3a9aaa9898f10496368a9963aad88160))
* **lighthouse-service:** Add compared value to payload ([#5496](https://github.com/keptn/keptn/issues/5496)) ([#6194](https://github.com/keptn/keptn/issues/6194)) ([f5af13c](https://github.com/keptn/keptn/commit/f5af13c6ba48250ab71432c8edc7aac5666a4ecf))
* Resource service endpoint handler implementation([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6461](https://github.com/keptn/keptn/issues/6461)) ([e19ed7b](https://github.com/keptn/keptn/commit/e19ed7b25fba4a52eedfa2b8dc68dea5c7864c5b)), closes [#6448](https://github.com/keptn/keptn/issues/6448) [#6448](https://github.com/keptn/keptn/issues/6448)
* Resource service first working version ([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6517](https://github.com/keptn/keptn/issues/6517)) ([00f81f1](https://github.com/keptn/keptn/commit/00f81f1186704bf9d115064f4dd95cb1f1ac42f2))
* **resource-service:** Add common interface for interacting with Git/Secrets ([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6411](https://github.com/keptn/keptn/issues/6411)) ([14af1d8](https://github.com/keptn/keptn/commit/14af1d881b4385db91bfbe352c4b36e1b076a8f9))
* **resource-service:** Complete implementation of Service ([#6530](https://github.com/keptn/keptn/issues/6530)) ([a91c116](https://github.com/keptn/keptn/commit/a91c116acb997429dbab429bab8a0bf352817422))
* **resource-service:** Improve git implementation ([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6510](https://github.com/keptn/keptn/issues/6510)) ([2f31d44](https://github.com/keptn/keptn/commit/2f31d44d1fd443746940039f00a497f00d02ad2e))
* **resource-service:** Improve git implementation and testing ([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6529](https://github.com/keptn/keptn/issues/6529)) ([91c5417](https://github.com/keptn/keptn/commit/91c54173b782692128cec8a7e9c0ac1c0e449270))
* **resource-service:** Resource service git implementation ([#6346](https://github.com/keptn/keptn/issues/6346)) ([#6438](https://github.com/keptn/keptn/issues/6438)) ([b6b4ac2](https://github.com/keptn/keptn/commit/b6b4ac26fa3b95f1337193de0592b245918bfe06)), closes [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462) [#6462](https://github.com/keptn/keptn/issues/6462)
* **shipyard-controller:** Add .triggered events to `lastEventTypes` property of materialized view ([#6220](https://github.com/keptn/keptn/issues/6220)) ([#6235](https://github.com/keptn/keptn/issues/6235)) ([dca04ea](https://github.com/keptn/keptn/commit/dca04eaa4538aab26700da180edfee42c60dcdd8))
* **shipyard-controller:** Allow to filter sequence states by multiple keptnContext IDs ([#6056](https://github.com/keptn/keptn/issues/6056)) ([#6093](https://github.com/keptn/keptn/issues/6093)) ([8a21919](https://github.com/keptn/keptn/commit/8a2191994f279d4ca5c4a1cc7f28dbb9c6c74213))
* **shipyard-controller:** handle sigterm ([#6051](https://github.com/keptn/keptn/issues/6051)) ([adfba40](https://github.com/keptn/keptn/commit/adfba40e4b21f01f0a806b5d525d3944305e6ca3))
* **shipyard-controller:** introduced sequence aborted state ([#6214](https://github.com/keptn/keptn/issues/6214)) ([02ab54b](https://github.com/keptn/keptn/commit/02ab54bfd949f720b41d7f71e3ff0aff06b754c5))
* validate shipyard.yaml when updating project ([#6222](https://github.com/keptn/keptn/issues/6222)) ([499352d](https://github.com/keptn/keptn/commit/499352d322a64b4a0207fb0be48d093198f6dcc1))


### Bug Fixes

* added delay and logging to graceful shutdown ([#6485](https://github.com/keptn/keptn/issues/6485)) ([#6486](https://github.com/keptn/keptn/issues/6486)) ([313db7f](https://github.com/keptn/keptn/commit/313db7ffacd9a202fa96eb95d80a3c08b3dc4eb5))
* Backup and restore integration test ([#6224](https://github.com/keptn/keptn/issues/6224)) ([7d622f8](https://github.com/keptn/keptn/commit/7d622f8f577d721a642094cb3de541b7a88a8671))
* **bridge:** Allow server inline script for base href ([#6248](https://github.com/keptn/keptn/issues/6248)) ([adebbbb](https://github.com/keptn/keptn/commit/adebbbbde0a8dbd0fa39d0e54766a6b4f34e0389))
* **bridge:** Fix problem with redirect and headers on cluster ([7407bcd](https://github.com/keptn/keptn/commit/7407bcdcfa996583d6e028afe4adf10950fd21ab))
* **bridge:** fix showing error message for OAUTH ([#6294](https://github.com/keptn/keptn/issues/6294)) ([6120087](https://github.com/keptn/keptn/commit/61200870c3f75ac0471b90640962dcb2a1d6f6ba))
* **bridge:** Fixed bridge server test ([#6314](https://github.com/keptn/keptn/issues/6314)) ([2d59f64](https://github.com/keptn/keptn/commit/2d59f6477ebe105f14989e5566f2fe181bb65fbb))
* **bridge:** Fixed bridge server tests ([#6261](https://github.com/keptn/keptn/issues/6261)) ([9f02adc](https://github.com/keptn/keptn/commit/9f02adc49695fada972039ac310f04e8c8907e69))
* **bridge:** Fixed environment screen update issues ([#6271](https://github.com/keptn/keptn/issues/6271)) ([0d5ff40](https://github.com/keptn/keptn/commit/0d5ff4002f9066f7ade3addd22da530c38b803f3))
* **bridge:** Fixed incorrect deployment link title ([#6304](https://github.com/keptn/keptn/issues/6304)) ([f237520](https://github.com/keptn/keptn/commit/f23752077c53e06204b14475b4f82b4403e9d913))
* **bridge:** Fixed removal of sequences if project endpoint of bridge server responds before projects endpoint of shipyard ([#6183](https://github.com/keptn/keptn/issues/6183)) ([8153fea](https://github.com/keptn/keptn/commit/8153fea4d707b6102c83292907a1301b6c7c4404))
* **bridge:** Remove .event suffix from payload variables in webhook ([#6396](https://github.com/keptn/keptn/issues/6396)) ([f67e5da](https://github.com/keptn/keptn/commit/f67e5daaa3315728d17a049a577d086ffffe0c2d))
* **bridge:** Update services on project change ([#6252](https://github.com/keptn/keptn/issues/6252)) ([65d4437](https://github.com/keptn/keptn/commit/65d4437d39b9889e28c1895e5e045a1af555d191))
* **cli:** Added rollback events to generated spec ([#3722](https://github.com/keptn/keptn/issues/3722)) ([#6161](https://github.com/keptn/keptn/issues/6161)) ([15ff2c6](https://github.com/keptn/keptn/commit/15ff2c63525a17fecac2e6ab222b8630a06007db))
* **cli:** Fix error handling during helm installation ([#6437](https://github.com/keptn/keptn/issues/6437)) ([#6583](https://github.com/keptn/keptn/issues/6583)) ([88c418b](https://github.com/keptn/keptn/commit/88c418b8ab56fc24a86d36163fd853a6ffabc1fe))
* **cli:** Print error message if service does not exist with `trigger delivery` ([#6351](https://github.com/keptn/keptn/issues/6351)) ([1d994a4](https://github.com/keptn/keptn/commit/1d994a4c1541879d3e6d50e4c4ca7b8a5c3c1c05))
* **cli:** project without upstream is defined as project without ([#6584](https://github.com/keptn/keptn/issues/6584)) ([aaf0a61](https://github.com/keptn/keptn/commit/aaf0a61f853c088ea82bc7de2eabde6a28560032))
* **cli:** Set CLI context before attempting to check for K8s context change ([#6449](https://github.com/keptn/keptn/issues/6449)) ([#6458](https://github.com/keptn/keptn/issues/6458)) ([3c2236d](https://github.com/keptn/keptn/commit/3c2236de055fc9f40183ce3e8c74c54efa1b81f5))
* **cli:** wrong handling of HTTPS in auth command ([#6268](https://github.com/keptn/keptn/issues/6268)) ([fa8fd1c](https://github.com/keptn/keptn/commit/fa8fd1cae9e81b351c378fbe38e7032f60d4c59f))
* **configuration-service:** Create tmpdir for unarchiving in /data/config ([#6329](https://github.com/keptn/keptn/issues/6329)) ([#6331](https://github.com/keptn/keptn/issues/6331)) ([a1f04af](https://github.com/keptn/keptn/commit/a1f04af05fae9366ec976e9b23239682901e74a6))
* **configuration-service:** Creating projects from empty upstream ([#6398](https://github.com/keptn/keptn/issues/6398)) ([#6399](https://github.com/keptn/keptn/issues/6399)) ([dc8337e](https://github.com/keptn/keptn/commit/dc8337e57605bd0760724ab128510810584d84b5))
* **configuration-service:** Fix permission issues for configuration service ([#6315](https://github.com/keptn/keptn/issues/6315)) ([#6317](https://github.com/keptn/keptn/issues/6317)) ([#6321](https://github.com/keptn/keptn/issues/6321)) ([61d9914](https://github.com/keptn/keptn/commit/61d99147389966c4db51e8eff0c30a347ca1e951))
* **configuration-service:** Make check for helm chart archives more strict ([#6447](https://github.com/keptn/keptn/issues/6447)) ([#6457](https://github.com/keptn/keptn/issues/6457)) ([babb3cd](https://github.com/keptn/keptn/commit/babb3cd11c76fe00562e2e83ddb8edcfb591b461))
* Dependencies cleanup ([#6369](https://github.com/keptn/keptn/issues/6369)) ([a38507b](https://github.com/keptn/keptn/commit/a38507b5d2fca643b69748f004835c5174fdfa88))
* Dependencies in lighthouse and remediation services ([#6368](https://github.com/keptn/keptn/issues/6368)) ([3f1646c](https://github.com/keptn/keptn/commit/3f1646cabf7460fef572e9ef3ebe3d0b3cc13cbc))
* Disable gitea installation on k3d ([#6408](https://github.com/keptn/keptn/issues/6408)) ([#6409](https://github.com/keptn/keptn/issues/6409)) ([cd984d4](https://github.com/keptn/keptn/commit/cd984d4f33b8e23df886e244045c7f7d9ba276ad))
* **distributor:** forcing restart if integration is malformed ([#6309](https://github.com/keptn/keptn/issues/6309)) ([#6363](https://github.com/keptn/keptn/issues/6363)) ([308261e](https://github.com/keptn/keptn/commit/308261e20b63160030c461eab36794f96c66dc62))
* fix graceful shutdown in sdk ([#6234](https://github.com/keptn/keptn/issues/6234)) ([a8db696](https://github.com/keptn/keptn/commit/a8db696c0f1df2eca3ef6cd283c3d5337ffd3d3d))
* Fix uniform integration test ([#6171](https://github.com/keptn/keptn/issues/6171)) ([#6174](https://github.com/keptn/keptn/issues/6174)) ([e55c398](https://github.com/keptn/keptn/commit/e55c3982b819de799fe1139eb4092aa72ff4a8d8))
* Graceful shutdown failing test ([#6462](https://github.com/keptn/keptn/issues/6462)) ([#6427](https://github.com/keptn/keptn/issues/6427)) ([4a28d73](https://github.com/keptn/keptn/commit/4a28d731d09c073137c388a1f003f6abfec511e8))
* Increase Bridge memory limits to avoid OOM ([#6562](https://github.com/keptn/keptn/issues/6562)) ([7f8d1a5](https://github.com/keptn/keptn/commit/7f8d1a5168ee3962ea3e6c61cbe1ca792b76736b))
* **installer:** Disable nats config reloader per default ([#6316](https://github.com/keptn/keptn/issues/6316)) ([#6318](https://github.com/keptn/keptn/issues/6318)) ([#6322](https://github.com/keptn/keptn/issues/6322)) ([d9263cf](https://github.com/keptn/keptn/commit/d9263cf574239bd518c517125d9e5ca80bf0e73f))
* **installer:** Remove obsolete openshift-route-service ([#6272](https://github.com/keptn/keptn/issues/6272)) ([#6389](https://github.com/keptn/keptn/issues/6389)) ([508dc25](https://github.com/keptn/keptn/commit/508dc2506414dd3bdfa1f6a290de8d0f1085294b))
* **installer:** Remove unneeded helm chart values ([#6419](https://github.com/keptn/keptn/issues/6419)) ([e5e508e](https://github.com/keptn/keptn/commit/e5e508e08dcba5d2c99d442c5a1beb406e0c16d5))
* **installer:** Use previous fsGroup per default, provide option to execute init container ([#6385](https://github.com/keptn/keptn/issues/6385)) ([#6386](https://github.com/keptn/keptn/issues/6386)) ([91eca02](https://github.com/keptn/keptn/commit/91eca02be99c02dd7df8f5b686b2d9e0368ba0b0))
* **lighthouse-service:** Lighthouse now fails if SLI fails  ([#6096](https://github.com/keptn/keptn/issues/6096)) ([#6281](https://github.com/keptn/keptn/issues/6281)) ([218cc39](https://github.com/keptn/keptn/commit/218cc390846aa79e2fc661a68a6006211960db95))
* **lighthouse-service:** Modified criteria example in SLO ([#6106](https://github.com/keptn/keptn/issues/6106)) ([#6404](https://github.com/keptn/keptn/issues/6404)) ([5b7bd19](https://github.com/keptn/keptn/commit/5b7bd198475703390ae55d21b72f9fcf83cebc76))
* minor fix in integration tests + added configuration-service securityContext ([#6540](https://github.com/keptn/keptn/issues/6540)) ([00cfe26](https://github.com/keptn/keptn/commit/00cfe26b0c0199a07a5728192c79c660c8711ce7))
* **mongodb-datastore:** Ensure MongoDB Client is not nil before retrieving database ([#6251](https://github.com/keptn/keptn/issues/6251)) ([#6255](https://github.com/keptn/keptn/issues/6255)) ([#6257](https://github.com/keptn/keptn/issues/6257)) ([fbaf0f0](https://github.com/keptn/keptn/commit/fbaf0f0b364f1c20c3303ffdba81b7465421ab12))
* **remediation-service:** add problemTitle to event message ([#5719](https://github.com/keptn/keptn/issues/5719)) ([c7d09d8](https://github.com/keptn/keptn/commit/c7d09d8a945c2a693fd52bc3dd9e339cafedae9e))
* Remove deprecated commands from CLI ([#6435](https://github.com/keptn/keptn/issues/6435)) ([d1625a7](https://github.com/keptn/keptn/commit/d1625a70512538f6593f44c67312f45fa97d8be5))
* Remove hardcoded namespace reference in installer ([#6286](https://github.com/keptn/keptn/issues/6286)) ([5396d6d](https://github.com/keptn/keptn/commit/5396d6d9030625dc0f23607abf01a9d879abd5be))
* Removed path issue within tests ([#6523](https://github.com/keptn/keptn/issues/6523)) ([#6525](https://github.com/keptn/keptn/issues/6525)) ([4295e2e](https://github.com/keptn/keptn/commit/4295e2e42abe8d1d3ffb7d8f97cf4d821c059874))
* Stabilize BackupRestore integration test ([#6344](https://github.com/keptn/keptn/issues/6344)) ([6fbd8cb](https://github.com/keptn/keptn/commit/6fbd8cb13f47ebce9fa85bdb4c047e2bea8527be))
* **statistics-service:** migrate data containing dots in keys ([#6266](https://github.com/keptn/keptn/issues/6266)) ([663c2bc](https://github.com/keptn/keptn/commit/663c2bc58d98aee11c9bee0952b6f04cd36314a0))
* **statistics-service:** migration of keptn service execution data ([#6324](https://github.com/keptn/keptn/issues/6324)) ([766a8e3](https://github.com/keptn/keptn/commit/766a8e335586917665bce8d2c73e498fe407cf81))
* Unit test in shipyard-controller ([#6370](https://github.com/keptn/keptn/issues/6370)) ([491a19a](https://github.com/keptn/keptn/commit/491a19a96c50a760868be63a5c66c35f2a9becda))
* Update dependencies ([#6381](https://github.com/keptn/keptn/issues/6381)) ([65a229a](https://github.com/keptn/keptn/commit/65a229aa31b71afab0b05789652c1a98309335e3))
* Update error messages ([#6197](https://github.com/keptn/keptn/issues/6197)) ([d43188e](https://github.com/keptn/keptn/commit/d43188e701e071877bc7cc97aaf20ff51f0a3ae9))
* Update go.sum of distributor ([#6367](https://github.com/keptn/keptn/issues/6367)) ([fc2b60a](https://github.com/keptn/keptn/commit/fc2b60a17dbe1486dc34afdf8fed5db562d6e07e))
* Update the JMeter Service to JMeter 5.4.2 ([#6405](https://github.com/keptn/keptn/issues/6405)) ([ccef405](https://github.com/keptn/keptn/commit/ccef4050a7a2497e6da01bf7d2f0632e72206a20))
* **webhook-service:** Disallow requests to loopback addresses ([#6361](https://github.com/keptn/keptn/issues/6361)) ([e7f814e](https://github.com/keptn/keptn/commit/e7f814e8bd6f975ff4a10c2bb7b056daaeae016f))


### Refactoring

* **bridge:** Move secret picker in own component ([#5733](https://github.com/keptn/keptn/issues/5733)) ([#6099](https://github.com/keptn/keptn/issues/6099)) ([a54f6a7](https://github.com/keptn/keptn/commit/a54f6a7ea6508eac4e27af61877c04cd6a52aa30))
* **bridge:** Replace data service mock with api service mock ([#5093](https://github.com/keptn/keptn/issues/5093)) ([101e472](https://github.com/keptn/keptn/commit/101e4728c2c272658497fe9c4023398757f703b7))


### Docs

* Add Keptn versioning and version compatibility document ([#5489](https://github.com/keptn/keptn/issues/5489)) ([c6e8a5c](https://github.com/keptn/keptn/commit/c6e8a5cda2e17ac6626448669f5572196f8d5511))
* **configuration-service:** Update API documentation ([#6008](https://github.com/keptn/keptn/issues/6008)) ([76f9ef2](https://github.com/keptn/keptn/commit/76f9ef2d24cc0a86e0e73d465cf6a6921b4de0cf))
* Update Integration Tests Developer documentation ([#6548](https://github.com/keptn/keptn/issues/6548)) ([d34b70c](https://github.com/keptn/keptn/commit/d34b70c4d5f0dd63909146d1b205c8056a910e1d))


### Performance

* **shipyard-controller:** Remove DB connection locking ([#6326](https://github.com/keptn/keptn/issues/6326)) ([690ce1c](https://github.com/keptn/keptn/commit/690ce1cacfe77529ff4189463ff5f3cd37496181))


### Other

*  update go_utils to safe version ([#6289](https://github.com/keptn/keptn/issues/6289)) ([f331482](https://github.com/keptn/keptn/commit/f3314822fdf91ef07da9cd1267c1ecc5bf656734))
* Add [@oleg-nenashev](https://github.com/oleg-nenashev) to the list of contributors ([#6256](https://github.com/keptn/keptn/issues/6256)) ([6817795](https://github.com/keptn/keptn/commit/6817795a3121ba651d52c8f0d42a23a830da36c3))
* **bridge:** Revert PR [#6341](https://github.com/keptn/keptn/issues/6341) ([#6585](https://github.com/keptn/keptn/issues/6585)) ([71c1e19](https://github.com/keptn/keptn/commit/71c1e19a45a2f706f69bcb6f75d77dc34faf871d))
* Bump JMeter to latest version ([307abf9](https://github.com/keptn/keptn/commit/307abf9142d52d97c2a5373c5cf6f61273d711f1))
* Correct example lighthouse criteria ([#6160](https://github.com/keptn/keptn/issues/6160)) ([#6406](https://github.com/keptn/keptn/issues/6406)) ([2d432eb](https://github.com/keptn/keptn/commit/2d432eb08105102750047cfa1e04beab71c3ae82)), closes [#6106](https://github.com/keptn/keptn/issues/6106)
* **distributor:** Upgrade go-utils, use thread safe fake EventSender in unit tests ([#6153](https://github.com/keptn/keptn/issues/6153)) ([da6fef0](https://github.com/keptn/keptn/commit/da6fef0c4bea6f26957c27757fd3832aace973fb))
* **helm-service:** Remove `service.create.finished` subscription ([#6181](https://github.com/keptn/keptn/issues/6181)) ([dc21c46](https://github.com/keptn/keptn/commit/dc21c46613625c4ff392db8d60b08b7cfb574387))
* Promote [@oleg-nenashev](https://github.com/oleg-nenashev) to maintainers ([#6522](https://github.com/keptn/keptn/issues/6522)) ([40e2deb](https://github.com/keptn/keptn/commit/40e2deb0fb13827efe26bd374dd67ade631e78f6))
* **secret-service:** updated README.md ([#6156](https://github.com/keptn/keptn/issues/6156)) ([d600e55](https://github.com/keptn/keptn/commit/d600e556236888dc32bab9c670d7565cab0dc9b0))
* update affiliation ([#6521](https://github.com/keptn/keptn/issues/6521)) ([642410f](https://github.com/keptn/keptn/commit/642410f4cc275fcc1900cc8ac6f354cee1b8a723))
* Update contributor lists ([#6450](https://github.com/keptn/keptn/issues/6450)) ([809532b](https://github.com/keptn/keptn/commit/809532b0f49fce2c4abf84e06ec1b4df62d1fbc6))

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
