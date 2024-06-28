# Changelog

## [1.24.1](https://github.com/momentohq/client-sdk-go/compare/v1.24.0...v1.24.1) (2024-06-28)


### Bug Fixes

* lowercase runtime-version and agent headers ([#440](https://github.com/momentohq/client-sdk-go/issues/440)) ([a4f4bb9](https://github.com/momentohq/client-sdk-go/commit/a4f4bb95318c8226ec1d6ec5f7ec8ddd9173e2b1))


### Miscellaneous

* add SDK version and runtime version headers for store client ([#436](https://github.com/momentohq/client-sdk-go/issues/436)) ([05b7b53](https://github.com/momentohq/client-sdk-go/commit/05b7b53d7d299046a1d7038be691eb373428cde8))
* **deps-dev:** bump braces in /examples/aws-lambda/infrastructure ([#441](https://github.com/momentohq/client-sdk-go/issues/441)) ([f440d6d](https://github.com/momentohq/client-sdk-go/commit/f440d6d868b9fe0332fd8325bdde9bce51494ed2))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#437](https://github.com/momentohq/client-sdk-go/issues/437)) ([9610373](https://github.com/momentohq/client-sdk-go/commit/9610373e17a9b751d8e0fae6cf15add65ab076fa))
* **deps:** bump golang.org/x/net from 0.21.0 to 0.23.0 ([#442](https://github.com/momentohq/client-sdk-go/issues/442)) ([03bc194](https://github.com/momentohq/client-sdk-go/commit/03bc194785acbc1653dd15264848f005e7a50cd7))
* **deps:** bump google.golang.org/protobuf ([#443](https://github.com/momentohq/client-sdk-go/issues/443)) ([d80b634](https://github.com/momentohq/client-sdk-go/commit/d80b63417b533d4ceb65e294e4ae3d4916e5d44b))

## [1.24.0](https://github.com/momentohq/client-sdk-go/compare/v1.23.1...v1.24.0) (2024-06-26)


### Features

* add Agent and Runtime-Version header interceptors, add release-please ([#422](https://github.com/momentohq/client-sdk-go/issues/422)) ([a1a47bc](https://github.com/momentohq/client-sdk-go/commit/a1a47bc13dd7f53cb0ff3ac837ca7908414ef52f))


### Bug Fixes

* super small request timeout tests should all be 1.Nanosecond ([#432](https://github.com/momentohq/client-sdk-go/issues/432)) ([8538ffe](https://github.com/momentohq/client-sdk-go/commit/8538ffe8eadf481d7a19dd947cc6c5b51ac2c61c))


### Miscellaneous

* disable storage tests ([#429](https://github.com/momentohq/client-sdk-go/issues/429)) ([b3e26c3](https://github.com/momentohq/client-sdk-go/commit/b3e26c3a373f59ff2612814f6fcebec1fede166e))
* reduce request timeout for timeout test ([#431](https://github.com/momentohq/client-sdk-go/issues/431)) ([cd522fb](https://github.com/momentohq/client-sdk-go/commit/cd522fb5d7c23ce904ea7cb1c8ad104ab9ebb348))
* update metadata message signifying item not found ([#427](https://github.com/momentohq/client-sdk-go/issues/427)) ([ffbe5d3](https://github.com/momentohq/client-sdk-go/commit/ffbe5d36116d09afa0743898ea412ca112736b58))
