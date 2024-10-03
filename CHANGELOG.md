# Changelog

## [1.28.4](https://github.com/momentohq/client-sdk-go/compare/v1.28.3...v1.28.4) (2024-10-03)


### Miscellaneous

* **protos:** update protos to v0.119.0 and regenerate code ([#533](https://github.com/momentohq/client-sdk-go/issues/533)) ([f1ac517](https://github.com/momentohq/client-sdk-go/commit/f1ac5177f57799a697cd5a05e014a6b09b7e6a34))

## [1.28.3](https://github.com/momentohq/client-sdk-go/compare/v1.28.2...v1.28.3) (2024-10-03)


### Miscellaneous

* add makefile targets to update and build protos ([#525](https://github.com/momentohq/client-sdk-go/issues/525)) ([86b747d](https://github.com/momentohq/client-sdk-go/commit/86b747d002f9489eb184165901c079f80e8e4389))
* enable consistent reads on all canary test targets ([#523](https://github.com/momentohq/client-sdk-go/issues/523)) ([a6948b9](https://github.com/momentohq/client-sdk-go/commit/a6948b98e18ec458f64591c0a59ef9a6c3e0fc7b))
* fix formatting in Makefile ([#535](https://github.com/momentohq/client-sdk-go/issues/535)) ([d97f229](https://github.com/momentohq/client-sdk-go/commit/d97f229a0c8b269f612b3101b5f3674c8554aef7))

## [1.28.2](https://github.com/momentohq/client-sdk-go/compare/v1.28.1...v1.28.2) (2024-09-20)


### Miscellaneous

* Add a flag and make command to use consistent reads in tests ([#518](https://github.com/momentohq/client-sdk-go/issues/518)) ([79a60df](https://github.com/momentohq/client-sdk-go/commit/79a60dfef3d19003bf7c67836ec829cef032aa4b))
* improve logging if we exceed the max number of concurrent streams ([#520](https://github.com/momentohq/client-sdk-go/issues/520)) ([2cc3ec2](https://github.com/momentohq/client-sdk-go/commit/2cc3ec2be2e686870790c9fb355d027da2c78cbe))
* log whether consistent reads are enabled for integration tests ([#519](https://github.com/momentohq/client-sdk-go/issues/519)) ([22af385](https://github.com/momentohq/client-sdk-go/commit/22af385fe32439ea83aca465844176ba08229131))
* Use the consistent reads flag in integration tests ([#515](https://github.com/momentohq/client-sdk-go/issues/515)) ([0f25f1c](https://github.com/momentohq/client-sdk-go/commit/0f25f1c89d964ddc431bf1e0a00980b26f24c64c))

## [1.28.1](https://github.com/momentohq/client-sdk-go/compare/v1.28.0...v1.28.1) (2024-09-17)


### Miscellaneous

* minor tweaks to logging inside of topic client ([#513](https://github.com/momentohq/client-sdk-go/issues/513)) ([fccfc5e](https://github.com/momentohq/client-sdk-go/commit/fccfc5e3e7a260c2ffb6c899d9e9e583d8bf3f34))

## [1.28.0](https://github.com/momentohq/client-sdk-go/compare/v1.27.7...v1.28.0) (2024-09-17)


### Features

* improve logger interface to accept `any` ([#507](https://github.com/momentohq/client-sdk-go/issues/507)) ([5dd3929](https://github.com/momentohq/client-sdk-go/commit/5dd39296028e1fea7ef7b4b77d741bcdace8e137))

NOTE: for existing users, if you have written a custom logger that implements the `MomentoLogger` interface, you will need to make a minor change to convert the signatures of the log functions from accepting `...string` to `...any`. This improves interopability with existing go logging libraries, because it allows you to pass in non-string data types and use other formatting strings like `%d` and `%v` to interpolate them into your log messages.



### Miscellaneous

* update topics example to illustrate how to wait for subscription close ([#511](https://github.com/momentohq/client-sdk-go/issues/511)) ([38843bf](https://github.com/momentohq/client-sdk-go/commit/38843bffef807ecdc29f4ba86761ba8202107767))

## [1.27.7](https://github.com/momentohq/client-sdk-go/compare/v1.27.6...v1.27.7) (2024-09-16)


### Miscellaneous

* add logging to list fetch test ([#506](https://github.com/momentohq/client-sdk-go/issues/506)) ([c2601d1](https://github.com/momentohq/client-sdk-go/commit/c2601d1733a1a9732cc1d913edb9b16aa7df74b7))
* debug sorted fetch by rank - remove elements test ([#509](https://github.com/momentohq/client-sdk-go/issues/509)) ([46ac38b](https://github.com/momentohq/client-sdk-go/commit/46ac38b59897213adea49eb865d0673c0122bc81))

## [1.27.6](https://github.com/momentohq/client-sdk-go/compare/v1.27.5...v1.27.6) (2024-09-12)


### Miscellaneous

* better debug logging on discontinuities ([#503](https://github.com/momentohq/client-sdk-go/issues/503)) ([0133f21](https://github.com/momentohq/client-sdk-go/commit/0133f2124dc1c52bde70707ca5e3200d2cc15a86))

## [1.27.5](https://github.com/momentohq/client-sdk-go/compare/v1.27.4...v1.27.5) (2024-09-10)


### Miscellaneous

* add epsilon sleep durations/reduce default ttl on tests ([#500](https://github.com/momentohq/client-sdk-go/issues/500)) ([b1241a9](https://github.com/momentohq/client-sdk-go/commit/b1241a954f9cf9f0f77de68a837eed075561e261))

## [1.27.4](https://github.com/momentohq/client-sdk-go/compare/v1.27.3...v1.27.4) (2024-09-10)


### Bug Fixes

* remaining ttl &lt;= ttl ([#497](https://github.com/momentohq/client-sdk-go/issues/497)) ([e3593cd](https://github.com/momentohq/client-sdk-go/commit/e3593cd3282ede37660c9ad8a9e493bf1efc392b))


### Miscellaneous

* add logging to failing canary tests ([#499](https://github.com/momentohq/client-sdk-go/issues/499)) ([bfe58f2](https://github.com/momentohq/client-sdk-go/commit/bfe58f21167bcda5ef2c319f00313a3be95c261c))

## [1.27.3](https://github.com/momentohq/client-sdk-go/compare/v1.27.2...v1.27.3) (2024-09-05)


### Miscellaneous

* debug canary tests ([#494](https://github.com/momentohq/client-sdk-go/issues/494)) ([8a5d399](https://github.com/momentohq/client-sdk-go/commit/8a5d39922d95bc3b70f729f7bca2f95bc96fe5d4))

## [1.27.2](https://github.com/momentohq/client-sdk-go/compare/v1.27.1...v1.27.2) (2024-09-04)


### Bug Fixes

* detailed sub item test canary error ([#491](https://github.com/momentohq/client-sdk-go/issues/491)) ([8adee94](https://github.com/momentohq/client-sdk-go/commit/8adee949b0bb845dbd3f0880b1a1ac8b584041bd))

## [1.27.1](https://github.com/momentohq/client-sdk-go/compare/v1.27.0...v1.27.1) (2024-08-30)


### Bug Fixes

* use ginkgo label matching to match on storage filter ([#487](https://github.com/momentohq/client-sdk-go/issues/487)) ([7f36e7d](https://github.com/momentohq/client-sdk-go/commit/7f36e7d82367a498060813861d58097a470e380e))
* use unique test keys ([#488](https://github.com/momentohq/client-sdk-go/issues/488)) ([f8fbd5e](https://github.com/momentohq/client-sdk-go/commit/f8fbd5ea7b666ab39793895d66fd42d5cb2d1cf0))


### Miscellaneous

* remove publisher id from simple topics example and use only one polling function at a time ([#489](https://github.com/momentohq/client-sdk-go/issues/489)) ([d97e718](https://github.com/momentohq/client-sdk-go/commit/d97e7182b473e92349f2d9e7fad5c6e3a51e6bc7))
* set topics resubscribe delay to 500ms ([#490](https://github.com/momentohq/client-sdk-go/issues/490)) ([8d70251](https://github.com/momentohq/client-sdk-go/commit/8d702518c8ea6fc2fb681b530aba6bb8ff7f11b2))
* update go topics examples ([#483](https://github.com/momentohq/client-sdk-go/issues/483)) ([ff0332f](https://github.com/momentohq/client-sdk-go/commit/ff0332f2d0efa2d9deb82915cc7ed9988f35bc65))

## [1.27.0](https://github.com/momentohq/client-sdk-go/compare/v1.26.2...v1.27.0) (2024-08-20)


### Features

* return topics subscription items with value, publisher id, and sequence number ([#476](https://github.com/momentohq/client-sdk-go/issues/476)) ([4808388](https://github.com/momentohq/client-sdk-go/commit/4808388cc0ce89a61fafa522ae7a38ad2fe2b1a1))


### Miscellaneous

* add per-service makefile targets ([#481](https://github.com/momentohq/client-sdk-go/issues/481)) ([4e758ce](https://github.com/momentohq/client-sdk-go/commit/4e758ce806769f0276b732bcde17b31e6ee2850d))
* clean up and try to speed up tests ([#477](https://github.com/momentohq/client-sdk-go/issues/477)) ([e2a9870](https://github.com/momentohq/client-sdk-go/commit/e2a987047625ae04368606520678414bab025c3b))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#474](https://github.com/momentohq/client-sdk-go/issues/474)) ([55be217](https://github.com/momentohq/client-sdk-go/commit/55be21725315cc60fa3aad8a557c9fb2dc3ea3d1))
* try to make flaky tests less flaky ([#482](https://github.com/momentohq/client-sdk-go/issues/482)) ([10500ea](https://github.com/momentohq/client-sdk-go/commit/10500ea2014a3ddeae45bc665a28a7f9cc1f9677))

## [1.26.2](https://github.com/momentohq/client-sdk-go/compare/v1.26.1...v1.26.2) (2024-08-07)


### Bug Fixes

* use sync.Map to track first-time agent headers ([#473](https://github.com/momentohq/client-sdk-go/issues/473)) ([d7fffa9](https://github.com/momentohq/client-sdk-go/commit/d7fffa9ae7db6aa3d5e9532b6e6604e3c5a1d697))


### Miscellaneous

* add GenerateApiKey dev docs snippet ([#463](https://github.com/momentohq/client-sdk-go/issues/463)) ([97ba7b3](https://github.com/momentohq/client-sdk-go/commit/97ba7b3d3ae49aea1a335ca84637d515a6ab2279))
* remove tests that require account session token ([#472](https://github.com/momentohq/client-sdk-go/issues/472)) ([d574a74](https://github.com/momentohq/client-sdk-go/commit/d574a7457581f6df8776d4ce936f751e1a9d1783))

## [1.26.1](https://github.com/momentohq/client-sdk-go/compare/v1.26.0...v1.26.1) (2024-07-19)


### Bug Fixes

* package api key and endpoint correctly in generate token/api key responses ([#467](https://github.com/momentohq/client-sdk-go/issues/467)) ([2a4ee7f](https://github.com/momentohq/client-sdk-go/commit/2a4ee7fb83ea0f4b917a1988d731da4eeb5b7b01))

## [1.26.0](https://github.com/momentohq/client-sdk-go/compare/v1.25.0...v1.26.0) (2024-07-19)


### Features

* add RefreshApiKey ([#465](https://github.com/momentohq/client-sdk-go/issues/465)) ([e1eb693](https://github.com/momentohq/client-sdk-go/commit/e1eb6938742616bf88210b5d4e48eccf41b561ff))

## [1.25.0](https://github.com/momentohq/client-sdk-go/compare/v1.24.2...v1.25.0) (2024-07-18)


### Features

* implement GenerateApiKey ([#460](https://github.com/momentohq/client-sdk-go/issues/460)) ([1371155](https://github.com/momentohq/client-sdk-go/commit/137115593779786fb701e9aee35639957d551580))


### Bug Fixes

* make sure to clean up caches and stores that can be leaked in tests ([#461](https://github.com/momentohq/client-sdk-go/issues/461)) ([f98bcd9](https://github.com/momentohq/client-sdk-go/commit/f98bcd9ccbc7bfe49de8308a2de002518a1b2eb5))


### Miscellaneous

* add dev docs snippets for storage client ([#457](https://github.com/momentohq/client-sdk-go/issues/457)) ([769cfc8](https://github.com/momentohq/client-sdk-go/commit/769cfc898060c56da7462eb36bae2bb7fd1a8a17))
* make sure push-to-main has test session token too ([#462](https://github.com/momentohq/client-sdk-go/issues/462)) ([268e771](https://github.com/momentohq/client-sdk-go/commit/268e7713ad2be16dd9f595f37d41e4a1c4547d65))

## [1.24.2](https://github.com/momentohq/client-sdk-go/compare/v1.24.1...v1.24.2) (2024-07-11)


### Miscellaneous

* add accessor for storage client logger ([#455](https://github.com/momentohq/client-sdk-go/issues/455)) ([27b9907](https://github.com/momentohq/client-sdk-go/commit/27b99075974a00dbb62b012d9de529cebd4d7152))
* fix example for clarity ([#453](https://github.com/momentohq/client-sdk-go/issues/453)) ([30511d0](https://github.com/momentohq/client-sdk-go/commit/30511d0eda679d94d7e7e233c765a9a853d64e17))

## [1.24.1](https://github.com/momentohq/client-sdk-go/compare/v1.24.0...v1.24.1) (2024-07-09)


### Bug Fixes

* ensure one-time headers are actually sent on only first request with non-empty info ([#447](https://github.com/momentohq/client-sdk-go/issues/447)) ([1ccf140](https://github.com/momentohq/client-sdk-go/commit/1ccf140a5d1cfff84fd725c2e3b365653a097e8f))
* lowercase runtime-version and agent headers ([#440](https://github.com/momentohq/client-sdk-go/issues/440)) ([a4f4bb9](https://github.com/momentohq/client-sdk-go/commit/a4f4bb95318c8226ec1d6ec5f7ec8ddd9173e2b1))


### Miscellaneous

* add SDK version and runtime version headers for store client ([#436](https://github.com/momentohq/client-sdk-go/issues/436)) ([05b7b53](https://github.com/momentohq/client-sdk-go/commit/05b7b53d7d299046a1d7038be691eb373428cde8))
* add strings to Describe()s and It()s for test segmentation ([#444](https://github.com/momentohq/client-sdk-go/issues/444)) ([a8f5e04](https://github.com/momentohq/client-sdk-go/commit/a8f5e04d440b4a91ceb0ef39b77039b82e29b4e0))
* **deps-dev:** bump braces in /examples/aws-lambda/infrastructure ([#441](https://github.com/momentohq/client-sdk-go/issues/441)) ([f440d6d](https://github.com/momentohq/client-sdk-go/commit/f440d6d868b9fe0332fd8325bdde9bce51494ed2))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#437](https://github.com/momentohq/client-sdk-go/issues/437)) ([9610373](https://github.com/momentohq/client-sdk-go/commit/9610373e17a9b751d8e0fae6cf15add65ab076fa))
* **deps:** bump golang.org/x/net from 0.21.0 to 0.23.0 ([#442](https://github.com/momentohq/client-sdk-go/issues/442)) ([03bc194](https://github.com/momentohq/client-sdk-go/commit/03bc194785acbc1653dd15264848f005e7a50cd7))
* **deps:** bump golang.org/x/net in /examples/aws-lambda/lambda ([#450](https://github.com/momentohq/client-sdk-go/issues/450)) ([ce9ce33](https://github.com/momentohq/client-sdk-go/commit/ce9ce33ffbdf66e2f2a32e81eccd8b5c93a33a32))
* **deps:** bump google.golang.org/protobuf ([#443](https://github.com/momentohq/client-sdk-go/issues/443)) ([d80b634](https://github.com/momentohq/client-sdk-go/commit/d80b63417b533d4ceb65e294e4ae3d4916e5d44b))
* remove NextToken from ListCachesRequest ([#449](https://github.com/momentohq/client-sdk-go/issues/449)) ([538add1](https://github.com/momentohq/client-sdk-go/commit/538add16f6c759fdf27ac244685331cf53fdbe1c))
* uncomment storage client tests ([#448](https://github.com/momentohq/client-sdk-go/issues/448)) ([f8c0ca1](https://github.com/momentohq/client-sdk-go/commit/f8c0ca1a17a3162ab7e5459b6f947d161fbfb61f))
* update get responses ([#433](https://github.com/momentohq/client-sdk-go/issues/433)) ([06b59dc](https://github.com/momentohq/client-sdk-go/commit/06b59dc67ef31ce1dd1e2ceeefd16cb81e5cf359))

## [1.24.0](https://github.com/momentohq/client-sdk-go/compare/v1.23.1...v1.24.0) (2024-06-26)


### Features

* add Agent and Runtime-Version header interceptors, add release-please ([#422](https://github.com/momentohq/client-sdk-go/issues/422)) ([a1a47bc](https://github.com/momentohq/client-sdk-go/commit/a1a47bc13dd7f53cb0ff3ac837ca7908414ef52f))


### Bug Fixes

* super small request timeout tests should all be 1.Nanosecond ([#432](https://github.com/momentohq/client-sdk-go/issues/432)) ([8538ffe](https://github.com/momentohq/client-sdk-go/commit/8538ffe8eadf481d7a19dd947cc6c5b51ac2c61c))


### Miscellaneous

* disable storage tests ([#429](https://github.com/momentohq/client-sdk-go/issues/429)) ([b3e26c3](https://github.com/momentohq/client-sdk-go/commit/b3e26c3a373f59ff2612814f6fcebec1fede166e))
* reduce request timeout for timeout test ([#431](https://github.com/momentohq/client-sdk-go/issues/431)) ([cd522fb](https://github.com/momentohq/client-sdk-go/commit/cd522fb5d7c23ce904ea7cb1c8ad104ab9ebb348))
* update metadata message signifying item not found ([#427](https://github.com/momentohq/client-sdk-go/issues/427)) ([ffbe5d3](https://github.com/momentohq/client-sdk-go/commit/ffbe5d36116d09afa0743898ea412ca112736b58))
