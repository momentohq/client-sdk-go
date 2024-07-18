# Changelog

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
