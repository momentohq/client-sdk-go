# Changelog

## [1.38.1](https://github.com/momentohq/client-sdk-go/compare/v1.38.0...v1.38.1) (2025-09-04)


### Bug Fixes

* changed refreshApiKey type to make compatible with credential provider ([#652](https://github.com/momentohq/client-sdk-go/issues/652)) ([1b68520](https://github.com/momentohq/client-sdk-go/commit/1b685203f61e301a16372f8cd8b496c05d79113f))


### Miscellaneous

* Adding Debug Logs ([#654](https://github.com/momentohq/client-sdk-go/issues/654)) ([fe06917](https://github.com/momentohq/client-sdk-go/commit/fe069177d98ba018ad4c690414260656e9e3a2dc))

## [1.38.0](https://github.com/momentohq/client-sdk-go/compare/v1.37.0...v1.38.0) (2025-08-20)


### Features

* add WithMaxSubscriptions config, deprecate NumGrpcChannels ([#645](https://github.com/momentohq/client-sdk-go/issues/645)) ([6745431](https://github.com/momentohq/client-sdk-go/commit/67454318ce04a72638d90393f0b9c2421107ceea))


### Bug Fixes

* add nil check to stream manager grpc pools ([#647](https://github.com/momentohq/client-sdk-go/issues/647)) ([87f2f8b](https://github.com/momentohq/client-sdk-go/commit/87f2f8beb1e535e7280e50074a400b4a3a1e4d22))
* clean up topics stream bookkeeping using a channel and new structs ([#644](https://github.com/momentohq/client-sdk-go/issues/644)) ([1ef78af](https://github.com/momentohq/client-sdk-go/commit/1ef78afdc5df93bb92830555d19294f5596ab152))
* make sure cache client closes all underlying clients upon close ([#646](https://github.com/momentohq/client-sdk-go/issues/646)) ([a06e803](https://github.com/momentohq/client-sdk-go/commit/a06e803a0e02ca9cba78c1ee2edbd4dd5579b49d))


### Miscellaneous

* add example for using momento-local metadata middleware ([#636](https://github.com/momentohq/client-sdk-go/issues/636)) ([2ab2323](https://github.com/momentohq/client-sdk-go/commit/2ab232340b48471fddb37f58e028a7922118fa33))
* add integration tests for new get/set if hash methods ([#612](https://github.com/momentohq/client-sdk-go/issues/612)) ([3c0cd5d](https://github.com/momentohq/client-sdk-go/commit/3c0cd5dcfe0bcc8510a2543b7eab3981a8ea5da9))
* clean up gzip compression tests ([#643](https://github.com/momentohq/client-sdk-go/issues/643)) ([8a3b486](https://github.com/momentohq/client-sdk-go/commit/8a3b486e12d788c80cd942d617df1508b0fe07d8))
* **deps:** bump github.com/golang-jwt/jwt/v4 ([#649](https://github.com/momentohq/client-sdk-go/issues/649)) ([9c54c41](https://github.com/momentohq/client-sdk-go/commit/9c54c416dea7290f228f82e02bdfe0f4477193e1))
* **deps:** bump github.com/golang-jwt/jwt/v4 from 4.3.0 to 4.5.2 ([#648](https://github.com/momentohq/client-sdk-go/issues/648)) ([ecbfe16](https://github.com/momentohq/client-sdk-go/commit/ecbfe164e0fa6e996001a4da0ad3a0aa8476f8db))
* improved gzip compression tests ([#641](https://github.com/momentohq/client-sdk-go/issues/641)) ([c377165](https://github.com/momentohq/client-sdk-go/commit/c377165fb803dc264cc58fbdf4d784a8c6aa36f3))
* reuse cache client in aws lambda example and update deps ([#639](https://github.com/momentohq/client-sdk-go/issues/639)) ([3d9771d](https://github.com/momentohq/client-sdk-go/commit/3d9771d04315002d95bf10c64743a003e23f17a4))

## [1.37.0](https://github.com/momentohq/client-sdk-go/compare/v1.36.1...v1.37.0) (2025-04-16)


### Features

* add public-facing momento-local middleware ([#634](https://github.com/momentohq/client-sdk-go/issues/634)) ([fc44a7a](https://github.com/momentohq/client-sdk-go/commit/fc44a7a09d9bc3f1a7efc7962c4eb12f9e801892))

## [1.36.1](https://github.com/momentohq/client-sdk-go/compare/v1.36.0...v1.36.1) (2025-04-16)


### Bug Fixes

* base compression middleware was not properly handling bytes ([#635](https://github.com/momentohq/client-sdk-go/issues/635)) ([b5c25d9](https://github.com/momentohq/client-sdk-go/commit/b5c25d905285491d7df219bfc0e308ea88dca757))


### Miscellaneous

* add zstd compression example using new module ([#632](https://github.com/momentohq/client-sdk-go/issues/632)) ([ebc9fcb](https://github.com/momentohq/client-sdk-go/commit/ebc9fcbdabcf1fa50ed71117c67ca243a137003d))
* minor revisions on zstd compression example ([#633](https://github.com/momentohq/client-sdk-go/issues/633)) ([e665823](https://github.com/momentohq/client-sdk-go/commit/e6658230a84fabfe19763dd1003e48208f3744a0))
* remove IncludeTypes from compression example ([#631](https://github.com/momentohq/client-sdk-go/issues/631)) ([69dba34](https://github.com/momentohq/client-sdk-go/commit/69dba349eaf75002856d5997cb1cf922fd915efd))
* update compression and middleware examples ([#629](https://github.com/momentohq/client-sdk-go/issues/629)) ([3a629ec](https://github.com/momentohq/client-sdk-go/commit/3a629ecbc70fd24d59323742b92872a6b78d8ec0))

## [1.36.0](https://github.com/momentohq/client-sdk-go/compare/v1.35.0...v1.36.0) (2025-04-15)


### Features

* add base and gzip compression middleware for scalar get and set requests ([#628](https://github.com/momentohq/client-sdk-go/issues/628)) ([b695689](https://github.com/momentohq/client-sdk-go/commit/b69568960e8790edbc4deaf8246ce09b544321ec))
* add FixedTimeoutRetryStrategy and tests ([#616](https://github.com/momentohq/client-sdk-go/issues/616)) ([ca1cd2e](https://github.com/momentohq/client-sdk-go/commit/ca1cd2e3830bd13c368c530e96984ba9e8a4abfd))


### Bug Fixes

* retry interceptor should not set retry deadline until after initial request ([#625](https://github.com/momentohq/client-sdk-go/issues/625)) ([f8d8bc0](https://github.com/momentohq/client-sdk-go/commit/f8d8bc03c222420fd8b5359135ced731c8b834a7))
* revise retry strategy interface additions to maintain backwards compatibility ([#627](https://github.com/momentohq/client-sdk-go/issues/627)) ([748f5cb](https://github.com/momentohq/client-sdk-go/commit/748f5cb2bcf5210e2fe57bc8a06870ee761b6fbd))


### Miscellaneous

* add a package for middleware implementations ([#624](https://github.com/momentohq/client-sdk-go/issues/624)) ([ab78c57](https://github.com/momentohq/client-sdk-go/commit/ab78c5774579d01a95fc3e6c7bb3d3c9df74e8df))
* add compression middleware example ([#618](https://github.com/momentohq/client-sdk-go/issues/618)) ([63ec741](https://github.com/momentohq/client-sdk-go/commit/63ec74144c04bde4a970bcaab401f462fbf4656a))

## [1.35.0](https://github.com/momentohq/client-sdk-go/compare/v1.34.0...v1.35.0) (2025-04-07)


### Features

* topics retries ([#610](https://github.com/momentohq/client-sdk-go/issues/610)) ([eb38e0c](https://github.com/momentohq/client-sdk-go/commit/eb38e0ce7dbc0806d9457216a4c28da521ae40eb))


### Miscellaneous

* add docs snippets for new get/set if hash methods ([#613](https://github.com/momentohq/client-sdk-go/issues/613)) ([e5e7c0f](https://github.com/momentohq/client-sdk-go/commit/e5e7c0f64ae3b2bc13ec22a950bbe95a13801fae))
* add middleware examples ([#617](https://github.com/momentohq/client-sdk-go/issues/617)) ([4db9c05](https://github.com/momentohq/client-sdk-go/commit/4db9c05f331f5436a2170faec7c5e6c8539a9319))
* fix the timing of the OnRequest handler ([#620](https://github.com/momentohq/client-sdk-go/issues/620)) ([de55607](https://github.com/momentohq/client-sdk-go/commit/de55607830a485584311196a597199567058e669))
* remove storage tests ([#615](https://github.com/momentohq/client-sdk-go/issues/615)) ([108d047](https://github.com/momentohq/client-sdk-go/commit/108d047341a59237d6d7df380cbfa6798ce66b56))

## [1.34.0](https://github.com/momentohq/client-sdk-go/compare/v1.33.2...v1.34.0) (2025-04-02)


### Features

* add middleware and data client retry strategies and tests  ([#601](https://github.com/momentohq/client-sdk-go/issues/601)) ([5a0c06f](https://github.com/momentohq/client-sdk-go/commit/5a0c06f89df123e9fa38f96f976296a958241aaf))
* add new get/set if hash apis ([#603](https://github.com/momentohq/client-sdk-go/issues/603)) ([77a3627](https://github.com/momentohq/client-sdk-go/commit/77a3627d4befdd4a99127d71851c0b7d34644045))


### Bug Fixes

* add topic subscribe timeout ([#596](https://github.com/momentohq/client-sdk-go/issues/596)) ([06337c3](https://github.com/momentohq/client-sdk-go/commit/06337c39bba021244ea52a19a5fe96b99994cd46))


### Miscellaneous

* fix make target in github workflow ([#609](https://github.com/momentohq/client-sdk-go/issues/609)) ([f3eee9e](https://github.com/momentohq/client-sdk-go/commit/f3eee9e86f0ca62c6d77c5fd9599e60bff16c5cb))
* fix request names in several requester objects ([#600](https://github.com/momentohq/client-sdk-go/issues/600)) ([7aeb3d4](https://github.com/momentohq/client-sdk-go/commit/7aeb3d4770abd7ad3a1d5bca6e16c93173a5f152))
* **protos:** update protos to v0.124.0 and regenerate code ([#608](https://github.com/momentohq/client-sdk-go/issues/608)) ([7a99843](https://github.com/momentohq/client-sdk-go/commit/7a998430c15c285deb26fea735e61f00d13e7234))
* remove cache name field from leaderboard requests ([#607](https://github.com/momentohq/client-sdk-go/issues/607)) ([3cd39a6](https://github.com/momentohq/client-sdk-go/commit/3cd39a64c4ee13e6705a6122588b008eb0d95976))
* update topics loadgen example ([#593](https://github.com/momentohq/client-sdk-go/issues/593)) ([b6fdb79](https://github.com/momentohq/client-sdk-go/commit/b6fdb794fdd8a58f510c7deefbaaf237e0b8f90a))

## [1.33.2](https://github.com/momentohq/client-sdk-go/compare/v1.33.1...v1.33.2) (2025-03-05)


### Bug Fixes

* set grpc deadline on publish requests and keep count of active subscriptions in stream topic grpc managers ([#588](https://github.com/momentohq/client-sdk-go/issues/588)) ([0d9d005](https://github.com/momentohq/client-sdk-go/commit/0d9d005bde3daebe272f2bf4ee0c6736a280f982))

## [1.33.1](https://github.com/momentohq/client-sdk-go/compare/v1.33.0...v1.33.1) (2025-03-03)


### Bug Fixes

* disable dynamic DNS service config ([#589](https://github.com/momentohq/client-sdk-go/issues/589)) ([a61ec16](https://github.com/momentohq/client-sdk-go/commit/a61ec160866d87f959e5bb5c6bba58ed4152539e))


### Miscellaneous

* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#585](https://github.com/momentohq/client-sdk-go/issues/585)) ([65d54d3](https://github.com/momentohq/client-sdk-go/commit/65d54d3cd85cba839b99c2aa4cc3a08a01f437da))

## [1.33.0](https://github.com/momentohq/client-sdk-go/compare/v1.32.1...v1.33.0) (2025-02-06)


### Features

* configure topic client to use separate publish and subscribe grpc channels ([#583](https://github.com/momentohq/client-sdk-go/issues/583)) ([66cc3cd](https://github.com/momentohq/client-sdk-go/commit/66cc3cd810dca26f6aaca75a0719d677f7710ca2))


### Miscellaneous

* add unit tests for retry eligibility strategy ([#582](https://github.com/momentohq/client-sdk-go/issues/582)) ([533c48e](https://github.com/momentohq/client-sdk-go/commit/533c48ea842aa26da3eb506bd78c69f4207b5689))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#580](https://github.com/momentohq/client-sdk-go/issues/580)) ([e8ef81e](https://github.com/momentohq/client-sdk-go/commit/e8ef81ecab4ef6a4092fdf60480bb22a1e7b1f73))

## [1.32.1](https://github.com/momentohq/client-sdk-go/compare/v1.32.0...v1.32.1) (2025-01-17)


### Miscellaneous

* log pretty print cache config on instantiation ([#577](https://github.com/momentohq/client-sdk-go/issues/577)) ([851b9c2](https://github.com/momentohq/client-sdk-go/commit/851b9c225283e4e48064a990dbbdd4182469a0bf))

## [1.32.0](https://github.com/momentohq/client-sdk-go/compare/v1.31.2...v1.32.0) (2025-01-17)


### Features

* add CredentialProvider constructor for momento-local connections ([#576](https://github.com/momentohq/client-sdk-go/issues/576)) ([ebdbe74](https://github.com/momentohq/client-sdk-go/commit/ebdbe741685a2354dec52925f6bdf31647f7da4b))


### Miscellaneous

* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#574](https://github.com/momentohq/client-sdk-go/issues/574)) ([c6510e2](https://github.com/momentohq/client-sdk-go/commit/c6510e2a18a05dd6e06316f7866cbc581597f905))

## [1.31.2](https://github.com/momentohq/client-sdk-go/compare/v1.31.1...v1.31.2) (2025-01-15)


### Bug Fixes

* add nil check before getting header and trailer for GetBatch and SetBatch ([#572](https://github.com/momentohq/client-sdk-go/issues/572)) ([5902394](https://github.com/momentohq/client-sdk-go/commit/5902394357a6045f9ab804f8a43aabddbb034f2f))

## [1.31.1](https://github.com/momentohq/client-sdk-go/compare/v1.31.0...v1.31.1) (2024-12-11)


### Bug Fixes

* add test-http-service to PHONEY ([#567](https://github.com/momentohq/client-sdk-go/issues/567)) ([9ebbbe6](https://github.com/momentohq/client-sdk-go/commit/9ebbbe63aefbfdfee3c0fcbe574973980838b54b))
* use atomic reads when using atomic writes ([#569](https://github.com/momentohq/client-sdk-go/issues/569)) ([22599bb](https://github.com/momentohq/client-sdk-go/commit/22599bbebf239f7e6436cdfb7c38736861397bb2))


### Miscellaneous

* add test-http-service target ([#565](https://github.com/momentohq/client-sdk-go/issues/565)) ([ced483d](https://github.com/momentohq/client-sdk-go/commit/ced483d298181d71a472234af0c2f997af5682cb))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#564](https://github.com/momentohq/client-sdk-go/issues/564)) ([c61f33f](https://github.com/momentohq/client-sdk-go/commit/c61f33f759b357db5addfd630dcda7885c49f4d9))
* **loadgen:** track cancelled errors in the load generator ([#568](https://github.com/momentohq/client-sdk-go/issues/568)) ([8d8ba5c](https://github.com/momentohq/client-sdk-go/commit/8d8ba5cbd4c083d8e8de660479afc868a3f30e53))

## [1.31.0](https://github.com/momentohq/client-sdk-go/compare/v1.30.0...v1.31.0) (2024-11-26)


### Features

* do not retry on cancelled codes ([#561](https://github.com/momentohq/client-sdk-go/issues/561)) ([be8848f](https://github.com/momentohq/client-sdk-go/commit/be8848f89d13dce2c48811fc393e14a7c0ca6b8c))


### Bug Fixes

* adjust total count in loadgen example ([#562](https://github.com/momentohq/client-sdk-go/issues/562)) ([d8bd37c](https://github.com/momentohq/client-sdk-go/commit/d8bd37cc776ff29a8926c776376081bc76b72cac))


### Miscellaneous

* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#555](https://github.com/momentohq/client-sdk-go/issues/555)) ([95a4193](https://github.com/momentohq/client-sdk-go/commit/95a41936c08c8db0abcded0e5c8f0ccdbe21cb3e))

## [1.30.0](https://github.com/momentohq/client-sdk-go/compare/v1.29.1...v1.30.0) (2024-11-23)


### Features

* add replica-reads test ([#559](https://github.com/momentohq/client-sdk-go/issues/559)) ([735593a](https://github.com/momentohq/client-sdk-go/commit/735593ac6800eda2cd210d101c2b16274effbda5))
* document all rpc's and if they are idempotent or not ([#558](https://github.com/momentohq/client-sdk-go/issues/558)) ([622b7d0](https://github.com/momentohq/client-sdk-go/commit/622b7d0057480515ad66694b2735b036676b597c))
* retry on cancelled status code ([#557](https://github.com/momentohq/client-sdk-go/issues/557)) ([d62e27d](https://github.com/momentohq/client-sdk-go/commit/d62e27dcc8ee9650ee2fbb2485547dbfaf09ffc2))


### Miscellaneous

* update license file ([#552](https://github.com/momentohq/client-sdk-go/issues/552)) ([f9ab3e4](https://github.com/momentohq/client-sdk-go/commit/f9ab3e4621bb4320e817385137c474237e0c4e74))

## [1.29.1](https://github.com/momentohq/client-sdk-go/compare/v1.29.0...v1.29.1) (2024-11-06)


### Bug Fixes

* add nil check before grabbing metadata when retrying subscribe ([#553](https://github.com/momentohq/client-sdk-go/issues/553)) ([da47ccb](https://github.com/momentohq/client-sdk-go/commit/da47ccb01b2d578e29a030d572b65ebddcaed523))

## [1.29.0](https://github.com/momentohq/client-sdk-go/compare/v1.28.7...v1.29.0) (2024-11-04)


### Features

* support topic sequence page ([#540](https://github.com/momentohq/client-sdk-go/issues/540)) ([cae3b3d](https://github.com/momentohq/client-sdk-go/commit/cae3b3d0db95d00348e7f7e3e0a195e251bba1ec))


### Bug Fixes

* add message wrapper to errors and improve resource exhausted error message by using metadata ([#551](https://github.com/momentohq/client-sdk-go/issues/551)) ([19bd021](https://github.com/momentohq/client-sdk-go/commit/19bd021c31935115fc9ca403e6b00eb150f5d1b5))


### Miscellaneous

* **deps:** bump github.com/momentohq/client-sdk-go from 1.27.0 to 1.28.7 in /examples ([#548](https://github.com/momentohq/client-sdk-go/issues/548)) ([2126ca8](https://github.com/momentohq/client-sdk-go/commit/2126ca886e102dfc2b58670bef602059b6e7cd3d))
* remove unused publish-golang step ([#549](https://github.com/momentohq/client-sdk-go/issues/549)) ([53ece53](https://github.com/momentohq/client-sdk-go/commit/53ece531fe36075ea224165c39b9591c072c172e))

## [1.28.7](https://github.com/momentohq/client-sdk-go/compare/v1.28.6...v1.28.7) (2024-10-24)


### Miscellaneous

* update contributing documentation ([#546](https://github.com/momentohq/client-sdk-go/issues/546)) ([babefcc](https://github.com/momentohq/client-sdk-go/commit/babefcc1c73283a70085d329cbe4861f4ac66358))

## [1.28.6](https://github.com/momentohq/client-sdk-go/compare/v1.28.4...v1.28.6) (2024-10-24)

### Bug Fixes

* interpret get rank response when sorted set found ([#541](https://github.com/momentohq/client-sdk-go/issues/541)) ([ffe0e9c](https://github.com/momentohq/client-sdk-go/commit/ffe0e9cc511b44879ffc5d19d543793eaa708dd7))

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
