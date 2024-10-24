# Changelog

## [2.0.0](https://github.com/momentohq/client-sdk-go/compare/v1.28.5...v2.0.0) (2024-10-24)


### ⚠ BREAKING CHANGES

* SortedSetFetch ByIndex and ByScore & response Value methods ([#241](https://github.com/momentohq/client-sdk-go/issues/241))
* Normalize dictionaries. ([#248](https://github.com/momentohq/client-sdk-go/issues/248))
* SortedSetGetScores polish. ([#245](https://github.com/momentohq/client-sdk-go/issues/245))
* Sortedset inconsistencies ([#239](https://github.com/momentohq/client-sdk-go/issues/239))
* Rename SortedSetGetScore -> SortedSetGetScores ([#227](https://github.com/momentohq/client-sdk-go/issues/227))
* SortedSetGetScore improvements. ([#220](https://github.com/momentohq/client-sdk-go/issues/220))
* Change SortedSet for the new protocol buffers. ([#215](https://github.com/momentohq/client-sdk-go/issues/215))
* Topics use the normal String and Bytes types. ([#205](https://github.com/momentohq/client-sdk-go/issues/205))
* Simplify passing values. ([#135](https://github.com/momentohq/client-sdk-go/issues/135))
* Ensure all methods return a response. ([#129](https://github.com/momentohq/client-sdk-go/issues/129))
* Rename SortedSet Name -> Value ([#127](https://github.com/momentohq/client-sdk-go/issues/127))
* Consolidate method implementations ([#104](https://github.com/momentohq/client-sdk-go/issues/104))

### Features

* add `UpdateTtl` ([#250](https://github.com/momentohq/client-sdk-go/issues/250)) ([aa5c628](https://github.com/momentohq/client-sdk-go/commit/aa5c6285f8666db68c2c25afa1910e2055c86dcc))
* add Agent and Runtime-Version header interceptors, add release-please ([#422](https://github.com/momentohq/client-sdk-go/issues/422)) ([a1a47bc](https://github.com/momentohq/client-sdk-go/commit/a1a47bc13dd7f53cb0ff3ac837ca7908414ef52f))
* add batch delete util ([#293](https://github.com/momentohq/client-sdk-go/issues/293)) ([327415a](https://github.com/momentohq/client-sdk-go/commit/327415aab62a58b4440be560ba5a6d4e064f3c05))
* add batch get utility ([#318](https://github.com/momentohq/client-sdk-go/issues/318)) ([a23a89a](https://github.com/momentohq/client-sdk-go/commit/a23a89aee74253aeeb448761570e62e3961ee9c3))
* add batch set if not exists to batchutils ([#355](https://github.com/momentohq/client-sdk-go/issues/355)) ([29ad96d](https://github.com/momentohq/client-sdk-go/commit/29ad96d36421bc2a8727da5502ebef823c923411))
* add batchSet operation to batchutils ([#346](https://github.com/momentohq/client-sdk-go/issues/346)) ([41396e7](https://github.com/momentohq/client-sdk-go/commit/41396e77bdcbd53b12bda294967acc13b1c70d2e))
* add caching pattern examples ([#411](https://github.com/momentohq/client-sdk-go/issues/411)) ([c6ee962](https://github.com/momentohq/client-sdk-go/commit/c6ee962a01c5d5d1540b6295940d2d38cd4d203c))
* add credential provider ([#55](https://github.com/momentohq/client-sdk-go/issues/55)) ([9155ba3](https://github.com/momentohq/client-sdk-go/commit/9155ba31049d3cae43e942c6b2fd051d13ae7cf5))
* add disposable tokens apis to golang sdk ([#368](https://github.com/momentohq/client-sdk-go/issues/368)) ([b7fe9db](https://github.com/momentohq/client-sdk-go/commit/b7fe9db22a9cc6a04679398816d4f1d31243baeb))
* add get/set async and batch perf test ([#409](https://github.com/momentohq/client-sdk-go/issues/409)) ([56e8563](https://github.com/momentohq/client-sdk-go/commit/56e856391ce1a74599c6c2faef5ddfe184e3a8ec))
* add getBatch and setBatch methods ([#393](https://github.com/momentohq/client-sdk-go/issues/393)) ([af34bf4](https://github.com/momentohq/client-sdk-go/commit/af34bf4cc0e0770c8ab868c6d61384154353182a))
* add grpc config options and turn off keepalive for Lambda config ([#380](https://github.com/momentohq/client-sdk-go/issues/380)) ([ea5951e](https://github.com/momentohq/client-sdk-go/commit/ea5951ec7fb163949a4509a37f987467b2aeaca9))
* add item get type and ttl ([#311](https://github.com/momentohq/client-sdk-go/issues/311)) ([683529e](https://github.com/momentohq/client-sdk-go/commit/683529e4ab9f8b2bd28a51af3884b5b676769552))
* add key exists ([#243](https://github.com/momentohq/client-sdk-go/issues/243)) ([098f4ae](https://github.com/momentohq/client-sdk-go/commit/098f4ae393fc4643cc7871f9ee8a01487babeee0))
* add logrus logging example ([#265](https://github.com/momentohq/client-sdk-go/issues/265)) ([ec8b864](https://github.com/momentohq/client-sdk-go/commit/ec8b86405c3703bb6743948338c2ef4e2e80f61d))
* add Momento Leaderboards ([#404](https://github.com/momentohq/client-sdk-go/issues/404)) ([5a1382a](https://github.com/momentohq/client-sdk-go/commit/5a1382ac45b2654a825a37e868bdba893f8f1154))
* add number of grpc channels as cache client config ([#375](https://github.com/momentohq/client-sdk-go/issues/375)) ([8823066](https://github.com/momentohq/client-sdk-go/commit/8823066cb77dde7ab971c0b405e527eff855fc69))
* add ping ([#249](https://github.com/momentohq/client-sdk-go/issues/249)) ([2121789](https://github.com/momentohq/client-sdk-go/commit/2121789e635a4e1bad736cc9fe429cbcd1e2a0d5))
* add read concern header ([#378](https://github.com/momentohq/client-sdk-go/issues/378)) ([80bf293](https://github.com/momentohq/client-sdk-go/commit/80bf293b158c7d5ba46f70b3a0ea07e97d72ae51))
* add RefreshApiKey ([#465](https://github.com/momentohq/client-sdk-go/issues/465)) ([e1eb693](https://github.com/momentohq/client-sdk-go/commit/e1eb6938742616bf88210b5d4e48eccf41b561ff))
* add request and response for set if not exists ([#331](https://github.com/momentohq/client-sdk-go/issues/331)) ([4213811](https://github.com/momentohq/client-sdk-go/commit/42138118f5883d129cb4034c0c8c1f42fd0c896c))
* add request and responses for lengths for set, sorted sets, and dict ([#338](https://github.com/momentohq/client-sdk-go/issues/338)) ([bc5a70d](https://github.com/momentohq/client-sdk-go/commit/bc5a70d3aa75a988208bf64355088746b3050627))
* add set contains elements ([#240](https://github.com/momentohq/client-sdk-go/issues/240)) ([a2a51d3](https://github.com/momentohq/client-sdk-go/commit/a2a51d3174ddc86172c32eb9e0c5ea13a39cc9cf))
* add SetIf APIs ([#381](https://github.com/momentohq/client-sdk-go/issues/381)) ([e01aa41](https://github.com/momentohq/client-sdk-go/commit/e01aa41c43342603f2be4d3b6ffaab898ffeff55))
* add SetPop ([#359](https://github.com/momentohq/client-sdk-go/issues/359)) ([2e1abed](https://github.com/momentohq/client-sdk-go/commit/2e1abed8049615091544b7df164c0e6dee0669ee))
* add simple logger ([#67](https://github.com/momentohq/client-sdk-go/issues/67)) ([fc11cbb](https://github.com/momentohq/client-sdk-go/commit/fc11cbb6afc0985254a86808b73e5945f8f9ffd1))
* add start and end index arguments to listFetch ([#345](https://github.com/momentohq/client-sdk-go/issues/345)) ([331d5a8](https://github.com/momentohq/client-sdk-go/commit/331d5a817a443bfc2ffc4302d7835f38e04108d3))
* add storage client ([#413](https://github.com/momentohq/client-sdk-go/issues/413)) ([f9421ca](https://github.com/momentohq/client-sdk-go/commit/f9421cafb3c4b6490c07c1836e56564d7472dccf))
* add support for default cache name on client creation ([#291](https://github.com/momentohq/client-sdk-go/issues/291)) ([db53bc0](https://github.com/momentohq/client-sdk-go/commit/db53bc0c15ba6168aeb6b4bc3b83e99b72403c99))
* add support for signing keys APIs ([#25](https://github.com/momentohq/client-sdk-go/issues/25)) ([1556f83](https://github.com/momentohq/client-sdk-go/commit/1556f83907f1e4853f296c1d55af600cdf4bc6be))
* add support for v1 auth tokens ([#301](https://github.com/momentohq/client-sdk-go/issues/301)) ([bb9a36d](https://github.com/momentohq/client-sdk-go/commit/bb9a36d4f799c337ab476799c160f1dce7b44656))
* add support to resubscribe to topics on disconnect ([#320](https://github.com/momentohq/client-sdk-go/issues/320)) ([367ea50](https://github.com/momentohq/client-sdk-go/commit/367ea505fee642af631f4899446adce0d7ad9450))
* add timeout to individual requests in batch call ([#348](https://github.com/momentohq/client-sdk-go/issues/348)) ([561f045](https://github.com/momentohq/client-sdk-go/commit/561f04558fa35b3c7e718cae2a7b3c594a9c6d56))
* adds customer facing MomentoError wrapper ([#14](https://github.com/momentohq/client-sdk-go/issues/14)) ([3b37174](https://github.com/momentohq/client-sdk-go/commit/3b371746a2ff4c58f9c3a3448fa1c8f0907fa0d3))
* adds interface for SimpleCacheClient ([#125](https://github.com/momentohq/client-sdk-go/issues/125)) ([e003656](https://github.com/momentohq/client-sdk-go/commit/e003656a6ffb50b9cfc2e2e7586aeefff935230c))
* adds New Cache Delete Item API ([#27](https://github.com/momentohq/client-sdk-go/issues/27)) ([cb56f3f](https://github.com/momentohq/client-sdk-go/commit/cb56f3f68382ccf60cbc3e28a157f498e9a23476)), closes [#26](https://github.com/momentohq/client-sdk-go/issues/26)
* adds retry interceptor ([#223](https://github.com/momentohq/client-sdk-go/issues/223)) ([847f4dd](https://github.com/momentohq/client-sdk-go/commit/847f4ddcc237c58a422f1370c5763247c29e1aad))
* adds sorted set initial impl in incubating ([#90](https://github.com/momentohq/client-sdk-go/issues/90)) ([f41ec5b](https://github.com/momentohq/client-sdk-go/commit/f41ec5b9d9c76a10905f5f450b8cad679956ab0c))
* adds standard momento client configuration ([#58](https://github.com/momentohq/client-sdk-go/issues/58)) ([b299589](https://github.com/momentohq/client-sdk-go/commit/b2995891b83aa09531762c8b35dc6373d9996cbc)), closes [#53](https://github.com/momentohq/client-sdk-go/issues/53)
* adds types for keys and values ([#78](https://github.com/momentohq/client-sdk-go/issues/78)) ([20d4d47](https://github.com/momentohq/client-sdk-go/commit/20d4d4730432300c3759dbd50ee1147d1614e01e))
* break out incubating package and adds basic pubsub support ([#42](https://github.com/momentohq/client-sdk-go/issues/42)) ([063229f](https://github.com/momentohq/client-sdk-go/commit/063229ff6cd51c3f67ea7453370fb32fedd501ef))
* breaks topic client into standalone client ([#213](https://github.com/momentohq/client-sdk-go/issues/213)) ([078ccc4](https://github.com/momentohq/client-sdk-go/commit/078ccc49fa17aad02d8fdbfc324d34db1634de17))
* Change back to returning responses as pointers ([#121](https://github.com/momentohq/client-sdk-go/issues/121)) ([f5489e9](https://github.com/momentohq/client-sdk-go/commit/f5489e965e200893b916da2efb95b6bedfd012e1))
* Change SortedSet for the new protocol buffers. ([#215](https://github.com/momentohq/client-sdk-go/issues/215)) ([7d12824](https://github.com/momentohq/client-sdk-go/commit/7d128248158c6dd1fc5c39aaa036c51e852c6385))
* Consolidate method implementations ([#104](https://github.com/momentohq/client-sdk-go/issues/104)) ([81d2f1e](https://github.com/momentohq/client-sdk-go/commit/81d2f1e199953a44a556d5f4ad2228c81e4dabd4))
* enable multiple subscription channels ([#325](https://github.com/momentohq/client-sdk-go/issues/325)) ([d1e930c](https://github.com/momentohq/client-sdk-go/commit/d1e930c773526559a2a9c5830c8cd2dda94223c5))
* Ensure all methods return a response. ([#129](https://github.com/momentohq/client-sdk-go/issues/129)) ([47ab6ee](https://github.com/momentohq/client-sdk-go/commit/47ab6ee1274c8dc78baf01cc95524ca2b57714a7))
* get sdk response types up to standards ([#63](https://github.com/momentohq/client-sdk-go/issues/63)) ([3b183c8](https://github.com/momentohq/client-sdk-go/commit/3b183c815d885430bba516fdef1805e301069502))
* handle h2 go-away from server for topic stream ([#151](https://github.com/momentohq/client-sdk-go/issues/151)) ([b9bf2dc](https://github.com/momentohq/client-sdk-go/commit/b9bf2dc983d3c9a2e562a1dca491181bf03487d3)), closes [#74](https://github.com/momentohq/client-sdk-go/issues/74)
* impl and expose SetIfNotExists API to clients/sdk ([#336](https://github.com/momentohq/client-sdk-go/issues/336)) ([ef558c3](https://github.com/momentohq/client-sdk-go/commit/ef558c3ca9def5a53ea79379bcbdf44b6b5e75bc))
* impl and expose SortedSetLength, SortedSetLengthByScore, DictionaryLength, and SetLength to client/sdk ([#339](https://github.com/momentohq/client-sdk-go/issues/339)) ([97f0a0d](https://github.com/momentohq/client-sdk-go/commit/97f0a0de6aad834cb056413aa6309f95b6484506))
* implement `sortedSetPutElement` ([#266](https://github.com/momentohq/client-sdk-go/issues/266)) ([9227dd7](https://github.com/momentohq/client-sdk-go/commit/9227dd7d66b89920f75b868ce0a184be62bf284b))
* implement `sortedSetRemoveElement` ([#267](https://github.com/momentohq/client-sdk-go/issues/267)) ([01c9adf](https://github.com/momentohq/client-sdk-go/commit/01c9adfe797193a3ba57cc71f33faa1444b61266))
* implement GenerateApiKey ([#460](https://github.com/momentohq/client-sdk-go/issues/460)) ([1371155](https://github.com/momentohq/client-sdk-go/commit/137115593779786fb701e9aee35639957d551580))
* improve logger interface to accept `any` ([#507](https://github.com/momentohq/client-sdk-go/issues/507)) ([5dd3929](https://github.com/momentohq/client-sdk-go/commit/5dd39296028e1fea7ef7b4b77d741bcdace8e137))
* Improve the "unexpected GRPC response error" ([#131](https://github.com/momentohq/client-sdk-go/issues/131)) ([abd4c59](https://github.com/momentohq/client-sdk-go/commit/abd4c593826a197f92369a88fee4b02f1950fb3d))
* initial impl of customer golang SDK for simple cache ([a333c8a](https://github.com/momentohq/client-sdk-go/commit/a333c8ae537a8e6ee6c450c9636ba1ae7e96e66a))
* lambda config, cdk and lambda example, and grpc eager connections ([#343](https://github.com/momentohq/client-sdk-go/issues/343)) ([4e9cec5](https://github.com/momentohq/client-sdk-go/commit/4e9cec51a278ebe73c8d1deaf3f4f8b6b72ab77f))
* makes clientside timeout setting name consistent w/ type ([#65](https://github.com/momentohq/client-sdk-go/issues/65)) ([03fc249](https://github.com/momentohq/client-sdk-go/commit/03fc249b3b94beb1cf3f3909fe2cacd833a73181))
* makes incubating package inherit cache functionality ([#57](https://github.com/momentohq/client-sdk-go/issues/57)) ([f8b2270](https://github.com/momentohq/client-sdk-go/commit/f8b22703ebaa8ebeb443207bae343aa78a64898f))
* More tweaks to logging API ([#264](https://github.com/momentohq/client-sdk-go/issues/264)) ([47cebfe](https://github.com/momentohq/client-sdk-go/commit/47cebfe955511daba881795f4b28d54ef28d9532))
* Normalize dictionaries. ([#248](https://github.com/momentohq/client-sdk-go/issues/248)) ([231aaa2](https://github.com/momentohq/client-sdk-go/commit/231aaa27ebd370621099b45a229444794819f80d))
* pass struct to functions instead of indivisual values; chore: update test; chore: separate external and internal responses and requests; ([f697f89](https://github.com/momentohq/client-sdk-go/commit/f697f892f729c41bb1ef799c7724eadaa281d303))
* Recognize pubsub heartbeat ([#141](https://github.com/momentohq/client-sdk-go/issues/141)) ([7464588](https://github.com/momentohq/client-sdk-go/commit/7464588a3dce48aaf2be54b343c3dded36495a1d))
* release 1.0 ([#284](https://github.com/momentohq/client-sdk-go/issues/284)) ([10d1166](https://github.com/momentohq/client-sdk-go/commit/10d1166a0c305631d65d67c90b842dd27877d2da))
* removes maxIdle and maxSessionMem from client config for now ([#216](https://github.com/momentohq/client-sdk-go/issues/216)) ([f315d6b](https://github.com/momentohq/client-sdk-go/commit/f315d6be608bdc9ada26b5da7f3cd7f067dbf947))
* Rename SortedSet Name -&gt; Value ([#127](https://github.com/momentohq/client-sdk-go/issues/127)) ([f6f0c7a](https://github.com/momentohq/client-sdk-go/commit/f6f0c7a1eb4f08d7a618f734b4af04c7afa68dd2))
* Rename SortedSetGetScore -&gt; SortedSetGetScores ([#227](https://github.com/momentohq/client-sdk-go/issues/227)) ([5debbbb](https://github.com/momentohq/client-sdk-go/commit/5debbbb0ad4a0dc543a144df3120933f1cad9808))
* retry initial topic subscriptions on LimitExceeded ([#419](https://github.com/momentohq/client-sdk-go/issues/419)) ([6763c45](https://github.com/momentohq/client-sdk-go/commit/6763c45be471e2fe137fa47088033f28f0d1b267))
* return topics subscription items with value, publisher id, and sequence number ([#476](https://github.com/momentohq/client-sdk-go/issues/476)) ([4808388](https://github.com/momentohq/client-sdk-go/commit/4808388cc0ce89a61fafa522ae7a38ad2fe2b1a1))
* rework response types to use interface and allow user to pass []byte or string messages to a topic  ([#72](https://github.com/momentohq/client-sdk-go/issues/72)) ([896ef56](https://github.com/momentohq/client-sdk-go/commit/896ef562e7abad6d11138b704f201df8a4b0856a))
* Simplify passing values. ([#135](https://github.com/momentohq/client-sdk-go/issues/135)) ([4771030](https://github.com/momentohq/client-sdk-go/commit/47710303a7f23045dafd8b431bf63284b3afe654))
* Sortedset inconsistencies ([#239](https://github.com/momentohq/client-sdk-go/issues/239)) ([ebc1ae7](https://github.com/momentohq/client-sdk-go/commit/ebc1ae7029a85994309523989c0bc1c1b442663f))
* SortedSetFetch ByIndex and ByScore & response Value methods ([#241](https://github.com/momentohq/client-sdk-go/issues/241)) ([9e14418](https://github.com/momentohq/client-sdk-go/commit/9e144185fa27767c60a99c1db1c1d7e9124e9f2e))
* SortedSetGetRank Order argument. ([#246](https://github.com/momentohq/client-sdk-go/issues/246)) ([650a65e](https://github.com/momentohq/client-sdk-go/commit/650a65eea360f441995e91899de962f87e845dc3))
* SortedSetGetScore improvements. ([#220](https://github.com/momentohq/client-sdk-go/issues/220)) ([d5edfeb](https://github.com/momentohq/client-sdk-go/commit/d5edfeb59f663338c0a9ba3889917426675f8b0f))
* SortedSetGetScores polish. ([#245](https://github.com/momentohq/client-sdk-go/issues/245)) ([31de897](https://github.com/momentohq/client-sdk-go/commit/31de897296f03a19ed4961873e3f9e183fbb8bb8))
* Split sortedSetFetch up into ByRank and ByScore functions ([#280](https://github.com/momentohq/client-sdk-go/issues/280)) ([cd60db3](https://github.com/momentohq/client-sdk-go/commit/cd60db3e1f150d532147c9d15ce55358490fe2f0))
* throw invalid argument err if score increment amount is not passed or 0 ([#142](https://github.com/momentohq/client-sdk-go/issues/142)) ([dc684bc](https://github.com/momentohq/client-sdk-go/commit/dc684bc16482ce007c7dd79e619279b98cdb193e))
* Topics use the normal String and Bytes types. ([#205](https://github.com/momentohq/client-sdk-go/issues/205)) ([1a7977c](https://github.com/momentohq/client-sdk-go/commit/1a7977cfb3cd297e6df1e019e29c53180d2e2b99))
* tweaks to how we configure logging ([#257](https://github.com/momentohq/client-sdk-go/issues/257)) ([7849dea](https://github.com/momentohq/client-sdk-go/commit/7849dea88dd3636d34d76e46fb410596dff5d0c4))
* update sorted set increment name ([#138](https://github.com/momentohq/client-sdk-go/issues/138)) ([337c90f](https://github.com/momentohq/client-sdk-go/commit/337c90f3830f1d962dd99baa276e9b272e01402a))
* updates examples with new response types ([#80](https://github.com/momentohq/client-sdk-go/issues/80)) ([87784ad](https://github.com/momentohq/client-sdk-go/commit/87784ad299fb44fbf23bf718e2331dad06e98fa1))
* usability improvements and testing updates ([#15](https://github.com/momentohq/client-sdk-go/issues/15)) ([8ecd287](https://github.com/momentohq/client-sdk-go/commit/8ecd287865092c0d04add81c8fc368bc52d82fd5))
* use "rank" terminology consistently for sorted sets ([#252](https://github.com/momentohq/client-sdk-go/issues/252)) ([d622189](https://github.com/momentohq/client-sdk-go/commit/d6221898c3a87d63c44fbb7496008a341ca480bc))
* use NumGrpcChannels to configure topic client instead of NumMaxSubscriptions ([#361](https://github.com/momentohq/client-sdk-go/issues/361)) ([f099a34](https://github.com/momentohq/client-sdk-go/commit/f099a34f96ec9d213d0425f9a47d63f7ba637a66))
* warn users if approaching grpc max concurrent streams limit, add method to close topics subscriptions ([#362](https://github.com/momentohq/client-sdk-go/issues/362)) ([312c10d](https://github.com/momentohq/client-sdk-go/commit/312c10dcb48580f9e07871d6c40d2624bb56ee33))


### Bug Fixes

* add `SortedSetRemoveElement` to cache interface ([#271](https://github.com/momentohq/client-sdk-go/issues/271)) ([ea0fc3c](https://github.com/momentohq/client-sdk-go/commit/ea0fc3c7a3efd58b36e888f41dfb99b1cd851d9c))
* add cache header to leaderboard requests ([#415](https://github.com/momentohq/client-sdk-go/issues/415)) ([0261867](https://github.com/momentohq/client-sdk-go/commit/0261867264606f0736ff17dfe7e8cf8fc333a950))
* add gofmt command ([da3316a](https://github.com/momentohq/client-sdk-go/commit/da3316a1e4b3635e1da88604b1b49fdbaa4977bc))
* add grpc error handling to provide more descriptive errors; fix: make CacheInfo public and modify a few methods; fix: remove import aliases; fix: make Close function return one error; ([9e2347e](https://github.com/momentohq/client-sdk-go/commit/9e2347ec19caad72d942ee344bc4347197a1bac2))
* add missing error codes and fix mappings ([#160](https://github.com/momentohq/client-sdk-go/issues/160)) ([38d78c3](https://github.com/momentohq/client-sdk-go/commit/38d78c350245e61ce31423862374412458fc28d9))
* add new TopicsConfiguration struct for topics client ([#321](https://github.com/momentohq/client-sdk-go/issues/321)) ([83b2e80](https://github.com/momentohq/client-sdk-go/commit/83b2e80b6bd5768f5c3a612bf88036854a19d14c))
* add testing for sets and some type fixes ([#159](https://github.com/momentohq/client-sdk-go/issues/159)) ([a8d215f](https://github.com/momentohq/client-sdk-go/commit/a8d215f917864c02884e21253047091761ce13f4))
* add validation to 'prepare' methods in requestor ([#200](https://github.com/momentohq/client-sdk-go/issues/200)) ([cc0ced3](https://github.com/momentohq/client-sdk-go/commit/cc0ced3c55470ffe924abe06a281f75bb93debc5))
* All values can be blank. ([#247](https://github.com/momentohq/client-sdk-go/issues/247)) ([de28bcd](https://github.com/momentohq/client-sdk-go/commit/de28bcdab9c88ea42822441fccb1e6f71ba27810))
* Allow subscription to complete when cancelled ([#414](https://github.com/momentohq/client-sdk-go/issues/414)) ([b2ded64](https://github.com/momentohq/client-sdk-go/commit/b2ded64788f16ee819873f0372b1e7bd344648a7))
* bump deps for vulns ([#372](https://github.com/momentohq/client-sdk-go/issues/372)) ([be36b83](https://github.com/momentohq/client-sdk-go/commit/be36b83c3371ce817a13eaa0032e3187427d9a7f))
* collaps var and const in () ([53d01df](https://github.com/momentohq/client-sdk-go/commit/53d01dfb81cbe699b8b5faf9e98f790ffbc4acd8))
* correct type assertion doc example ([#296](https://github.com/momentohq/client-sdk-go/issues/296)) ([4ed28eb](https://github.com/momentohq/client-sdk-go/commit/4ed28eb42949d53f61c0089fe8a16f655f9ccfb5))
* cyclic import issue ([f3d4739](https://github.com/momentohq/client-sdk-go/commit/f3d4739244d6890ad535ae0cbb881a4063858962))
* default logger wasn't spreading varargs properly ([#268](https://github.com/momentohq/client-sdk-go/issues/268)) ([128f884](https://github.com/momentohq/client-sdk-go/commit/128f884b852058fdb6ab0339e9e73ddee9e6402c))
* defer cancel for context ([9c94d89](https://github.com/momentohq/client-sdk-go/commit/9c94d89bd5434c2e78d07463b2df9c09c8282ae1))
* detailed sub item test canary error ([#491](https://github.com/momentohq/client-sdk-go/issues/491)) ([8adee94](https://github.com/momentohq/client-sdk-go/commit/8adee949b0bb845dbd3f0880b1a1ac8b584041bd))
* don't check for connection state when trying to resubscribe ([#356](https://github.com/momentohq/client-sdk-go/issues/356)) ([48e982d](https://github.com/momentohq/client-sdk-go/commit/48e982dc95e2939a585f3c58d5758cfea11faa45))
* don't log error message when subscription context cancelled ([#416](https://github.com/momentohq/client-sdk-go/issues/416)) ([83d9c82](https://github.com/momentohq/client-sdk-go/commit/83d9c820c04f06b8d88e1276caa8b20d83d3410c))
* ensure one-time headers are actually sent on only first request with non-empty info ([#447](https://github.com/momentohq/client-sdk-go/issues/447)) ([1ccf140](https://github.com/momentohq/client-sdk-go/commit/1ccf140a5d1cfff84fd725c2e3b365653a097e8f))
* Ensure tests clean up their caches. ([#120](https://github.com/momentohq/client-sdk-go/issues/120)) ([226deee](https://github.com/momentohq/client-sdk-go/commit/226deeeba17162f715040e52de30963472560442))
* Examples for recent changes. ([#211](https://github.com/momentohq/client-sdk-go/issues/211)) ([904574a](https://github.com/momentohq/client-sdk-go/commit/904574a0a70315ad28652a93cc78e0d5083d2fc1))
* export LogLevel ([#253](https://github.com/momentohq/client-sdk-go/issues/253)) ([9da6827](https://github.com/momentohq/client-sdk-go/commit/9da682710379144ff0aa2e6540c9e58933a39da4))
* fix misleading info in example ([#230](https://github.com/momentohq/client-sdk-go/issues/230)) ([d00828e](https://github.com/momentohq/client-sdk-go/commit/d00828e180ec9c570cf5dc74000e49f5abee5461))
* fix populating results from channel ([#349](https://github.com/momentohq/client-sdk-go/issues/349)) ([b862c91](https://github.com/momentohq/client-sdk-go/commit/b862c919fe9651bf1af3cd06aad0bac61b6f0ba6))
* fix return type for TopicPublish ([#163](https://github.com/momentohq/client-sdk-go/issues/163)) ([e46370e](https://github.com/momentohq/client-sdk-go/commit/e46370e49b8fc73bafd82cd33d95fcb882279977))
* fix version requirements in docs and go.mod ([#244](https://github.com/momentohq/client-sdk-go/issues/244)) ([cc999de](https://github.com/momentohq/client-sdk-go/commit/cc999de42f03ace382942e44f037021f7c0ff903))
* fixes and standardizes build process with other sdks ([#366](https://github.com/momentohq/client-sdk-go/issues/366)) ([95222a9](https://github.com/momentohq/client-sdk-go/commit/95222a946b968dbaa7960da117dd191ec101ae64))
* fixes to some minor misc issues ([#263](https://github.com/momentohq/client-sdk-go/issues/263)) ([339ef14](https://github.com/momentohq/client-sdk-go/commit/339ef14b7e09aeebe6cc9a1f447ef5cd98927b46))
* ignore discontinuity messages for now ([#69](https://github.com/momentohq/client-sdk-go/issues/69)) ([1e94541](https://github.com/momentohq/client-sdk-go/commit/1e94541294b4401fd4eaa0efb9e9e4351d734f1a))
* inspect single field response to determine hit or miss ([#164](https://github.com/momentohq/client-sdk-go/issues/164)) ([bfb768b](https://github.com/momentohq/client-sdk-go/commit/bfb768b1545fa17b86ca0220de14a3cb823c9237))
* interpret get rank response when sorted set found ([#541](https://github.com/momentohq/client-sdk-go/issues/541)) ([ffe0e9c](https://github.com/momentohq/client-sdk-go/commit/ffe0e9cc511b44879ffc5d19d543793eaa708dd7))
* lowercase runtime-version and agent headers ([#440](https://github.com/momentohq/client-sdk-go/issues/440)) ([a4f4bb9](https://github.com/momentohq/client-sdk-go/commit/a4f4bb95318c8226ec1d6ec5f7ec8ddd9173e2b1))
* make import alias more meaningful ([11e97d5](https://github.com/momentohq/client-sdk-go/commit/11e97d58d2836c1587cd943a8e4f6471d57b8a49))
* make ScsControlClient and ScsDataClient fields private ([cdb2f60](https://github.com/momentohq/client-sdk-go/commit/cdb2f600f5e5fc47c7cb4d3898a50bfece257e40))
* make sure context is the first argumnet when it's passed ([#123](https://github.com/momentohq/client-sdk-go/issues/123)) ([a9b6a5b](https://github.com/momentohq/client-sdk-go/commit/a9b6a5b8fc04aa0da7ebd321d6c571630532620c))
* Make sure the pubsub tests clean up after themselves even if they fail. ([#139](https://github.com/momentohq/client-sdk-go/issues/139)) ([ba50db9](https://github.com/momentohq/client-sdk-go/commit/ba50db93cfb98cddda721ca1e41ab5b772607d22))
* make sure to clean up caches and stores that can be leaked in tests ([#461](https://github.com/momentohq/client-sdk-go/issues/461)) ([f98bcd9](https://github.com/momentohq/client-sdk-go/commit/f98bcd9ccbc7bfe49de8308a2de002518a1b2eb5))
* Manual release passing wrong args to shared test action. ([#107](https://github.com/momentohq/client-sdk-go/issues/107)) ([b250004](https://github.com/momentohq/client-sdk-go/commit/b250004a5ac616ed46e5035de604c4406a38b19d))
* Marked the wrong PHONY in the Makefile. ([#96](https://github.com/momentohq/client-sdk-go/issues/96)) ([f715ae5](https://github.com/momentohq/client-sdk-go/commit/f715ae538d6f50ba5cfeef4158bbaf6a4c6e30a8))
* missed refactoring get responses somehow ([#228](https://github.com/momentohq/client-sdk-go/issues/228)) ([a22ee39](https://github.com/momentohq/client-sdk-go/commit/a22ee3976952c604f845a7bc8997a1ce35ff5817))
* navigable auth and config packages ([#217](https://github.com/momentohq/client-sdk-go/issues/217)) ([d10fa12](https://github.com/momentohq/client-sdk-go/commit/d10fa1251e565aaffc250f54b56b1391647dc7a3))
* No need for a token when checking out. ([#108](https://github.com/momentohq/client-sdk-go/issues/108)) ([514ec89](https://github.com/momentohq/client-sdk-go/commit/514ec895fde32d317ad16d1cfc7154b1e8ee3ca9))
* only err for any errors returned; fix: use error constant; fix: rename git action ([d5bfcc2](https://github.com/momentohq/client-sdk-go/commit/d5bfcc244b2c28dd1d10e3ecf83a3d7ae907d99f))
* package api key and endpoint correctly in generate token/api key responses ([#467](https://github.com/momentohq/client-sdk-go/issues/467)) ([2a4ee7f](https://github.com/momentohq/client-sdk-go/commit/2a4ee7fb83ea0f4b917a1988d731da4eeb5b7b01))
* proto generate command in readme ([#28](https://github.com/momentohq/client-sdk-go/issues/28)) ([721c064](https://github.com/momentohq/client-sdk-go/commit/721c0646d53fae85f0a651d3edb966791ade9798))
* pull out addHeadersInterceptor ([a644e4d](https://github.com/momentohq/client-sdk-go/commit/a644e4dd0945f57fe89a856c8c556ff549a0c9e3))
* pull out context timeout into constants ([9f37829](https://github.com/momentohq/client-sdk-go/commit/9f37829c1a8cb6ed2a03201e63d3f9c5be1cc923))
* Push to main. ([#93](https://github.com/momentohq/client-sdk-go/issues/93)) ([3f00f80](https://github.com/momentohq/client-sdk-go/commit/3f00f801f9d2994ddeea94ccbe581327c44894ae))
* race condition that can cause batchutils to leak goroutines ([#390](https://github.com/momentohq/client-sdk-go/issues/390)) ([5ac2fe8](https://github.com/momentohq/client-sdk-go/commit/5ac2fe88fe9e58fbb52bd0ed2612808ee63f4d93))
* refactor list operations ([#122](https://github.com/momentohq/client-sdk-go/issues/122)) ([5238a47](https://github.com/momentohq/client-sdk-go/commit/5238a473efaac987323419a44bb7f062e1b6f7c7))
* reintroduce NextToken to list caches request ([#426](https://github.com/momentohq/client-sdk-go/issues/426)) ([aa69612](https://github.com/momentohq/client-sdk-go/commit/aa696125736a4ce2bfd12c4259b06d74dd0e2393))
* remaining ttl &lt;= ttl ([#497](https://github.com/momentohq/client-sdk-go/issues/497)) ([e3593cd](https://github.com/momentohq/client-sdk-go/commit/e3593cd3282ede37660c9ad8a9e493bf1efc392b))
* Remove references to TEST_CACHE_NAME ([#199](https://github.com/momentohq/client-sdk-go/issues/199)) ([b882550](https://github.com/momentohq/client-sdk-go/commit/b882550fa0d76a2adbcec69c2676dd5c2ee8eff9)), closes [#187](https://github.com/momentohq/client-sdk-go/issues/187)
* remove tests package and move the test to the relevant package ([5fe74f5](https://github.com/momentohq/client-sdk-go/commit/5fe74f566766a43837d01f0eba5db6ec33cef1dc))
* remove unnecessary setUp function call from the test; fix: rename to err ([1589712](https://github.com/momentohq/client-sdk-go/commit/15897120bcc32f4654a0a2ab5ae519b79367ff7c))
* remove verbose variables; fix: pass struct pointer to functions ([6274d0c](https://github.com/momentohq/client-sdk-go/commit/6274d0c01ab14c7c8599adf8adf5ecb131eae3bd))
* rename control/data clients methods to use the same method names as simple cache client; fix: add a new line at the end of the files ([2dfe1cc](https://github.com/momentohq/client-sdk-go/commit/2dfe1cc3e20c71e676ce087b3cb3d9830027833a))
* rename TEST_CACHE_NAME ([fbb8eb6](https://github.com/momentohq/client-sdk-go/commit/fbb8eb61d46249a4b1bff9d1fc944f5cd41f9b66))
* restructure repo to only expose momento package ([0c771e3](https://github.com/momentohq/client-sdk-go/commit/0c771e30fcc31fd478370149b3d27c17af8e96ee))
* return an invalid argument error if the interpreted ttl is 0 milliseconds ([#337](https://github.com/momentohq/client-sdk-go/issues/337)) ([d340a1c](https://github.com/momentohq/client-sdk-go/commit/d340a1c0aacabc71e5564e087919e3f7c3b4f068))
* return if the context was cancelled while polling for subscription items ([#351](https://github.com/momentohq/client-sdk-go/issues/351)) ([f62abd0](https://github.com/momentohq/client-sdk-go/commit/f62abd00849b686cca41a1850242be5e41c11c1b))
* return value not being sent back for sorted set increment ([#134](https://github.com/momentohq/client-sdk-go/issues/134)) ([b38b95b](https://github.com/momentohq/client-sdk-go/commit/b38b95b2c5bca0be40a02cbd592027215532ff1f))
* Rewrite the tests using Ginkgo and Gomega, fix some bugs. ([#147](https://github.com/momentohq/client-sdk-go/issues/147)) ([4e5b6df](https://github.com/momentohq/client-sdk-go/commit/4e5b6df45381854719ff047a9b763b29737fbeff))
* separate out config for topics ([#322](https://github.com/momentohq/client-sdk-go/issues/322)) ([98346e3](https://github.com/momentohq/client-sdk-go/commit/98346e3f88034ca9e5127673047dfe08a607a54c))
* Set the output when generating the contributing file. ([#118](https://github.com/momentohq/client-sdk-go/issues/118)) ([8e49f48](https://github.com/momentohq/client-sdk-go/commit/8e49f48db8656d67e934671623f315812f3d3ede))
* simplify abstractions involved in sorted set get score ([#276](https://github.com/momentohq/client-sdk-go/issues/276)) ([e1728ee](https://github.com/momentohq/client-sdk-go/commit/e1728ee9797c199712145abe5cdc7ab14b921e17))
* sorted set remove request ([#133](https://github.com/momentohq/client-sdk-go/issues/133)) ([8370836](https://github.com/momentohq/client-sdk-go/commit/8370836bb3810bbf60309efc7168abadb0810fb2))
* stop parsing ECacheResult OK ([0e383a8](https://github.com/momentohq/client-sdk-go/commit/0e383a8c72c3b0cfee04185fdd9f85da138e2af1))
* super small request timeout tests should all be 1.Nanosecond ([#432](https://github.com/momentohq/client-sdk-go/issues/432)) ([8538ffe](https://github.com/momentohq/client-sdk-go/commit/8538ffe8eadf481d7a19dd947cc6c5b51ac2c61c))
* Topic name and cache validation ([#201](https://github.com/momentohq/client-sdk-go/issues/201)) ([a223cc5](https://github.com/momentohq/client-sdk-go/commit/a223cc5f42b3f973ef783e4232a569d738aeebe6))
* Update golang.org/x/net for minor security vulnerability. ([#149](https://github.com/momentohq/client-sdk-go/issues/149)) ([7fe35b6](https://github.com/momentohq/client-sdk-go/commit/7fe35b69fbb9f070c3035c366515820ca7ce0693))
* update GrpcErrorConverter; chore: update happy path test to include TestCacheName ([1e44e5b](https://github.com/momentohq/client-sdk-go/commit/1e44e5b70e299c3de93ca7593e2a3e6a0c31257b))
* update if statement to determine cache already exists in test; chore: remove unused structs from external responses ([5e6a882](https://github.com/momentohq/client-sdk-go/commit/5e6a8821277e247c9303f51454a71630dec9dd41))
* update lint_format_pr.yml ([f0be5f1](https://github.com/momentohq/client-sdk-go/commit/f0be5f1dca651ca9726fa50374578f6c35a8a7b6))
* update manual release workflow to make sure to wait for release job ([#18](https://github.com/momentohq/client-sdk-go/issues/18)) ([d1e430b](https://github.com/momentohq/client-sdk-go/commit/d1e430b54afd7cfdfc115a44321b29e199588dba))
* update SimpleCacheClient to take a struct pointer ([96c213f](https://github.com/momentohq/client-sdk-go/commit/96c213f97653780301990c8450d3fdb0c5f90bcf))
* update usage examples to remove obsolete error check ([#162](https://github.com/momentohq/client-sdk-go/issues/162)) ([f9ac0e5](https://github.com/momentohq/client-sdk-go/commit/f9ac0e509df146d2c0467fb69ba6ac4904286863))
* upgrade grpc lib and use NewClient instead of Dial ([#405](https://github.com/momentohq/client-sdk-go/issues/405)) ([adafbe3](https://github.com/momentohq/client-sdk-go/commit/adafbe38fee3de9bbbd63b5000024684f97f2328))
* use alternative data structure for dict ops ([#224](https://github.com/momentohq/client-sdk-go/issues/224)) ([021ef38](https://github.com/momentohq/client-sdk-go/commit/021ef38e2d8a7571076a3f9d8663fe6d3805be60))
* use Dial instead of NewClient while we wait for grpc-go release ([#400](https://github.com/momentohq/client-sdk-go/issues/400)) ([d474b9a](https://github.com/momentohq/client-sdk-go/commit/d474b9a2b0b0c8cba434d54d9b35ad3740f44583))
* use ginkgo label matching to match on storage filter ([#487](https://github.com/momentohq/client-sdk-go/issues/487)) ([7f36e7d](https://github.com/momentohq/client-sdk-go/commit/7f36e7d82367a498060813861d58097a470e380e))
* use just one grpc connection for topics ([#360](https://github.com/momentohq/client-sdk-go/issues/360)) ([842033e](https://github.com/momentohq/client-sdk-go/commit/842033e2b7c5a78365843f0e86b8c8cebd251148))
* use momento machine user token when generating/committing readme ([#376](https://github.com/momentohq/client-sdk-go/issues/376)) ([a64409f](https://github.com/momentohq/client-sdk-go/commit/a64409f47b2d8cc05eb638fc023166d0ca4054f7))
* use Noop logger for testing ([#274](https://github.com/momentohq/client-sdk-go/issues/274)) ([573acd4](https://github.com/momentohq/client-sdk-go/commit/573acd46e6fda64a5b9e1bf82c9adb98553883ab))
* use sync.Map to track first-time agent headers ([#473](https://github.com/momentohq/client-sdk-go/issues/473)) ([d7fffa9](https://github.com/momentohq/client-sdk-go/commit/d7fffa9ae7db6aa3d5e9532b6e6604e3c5a1d697))
* use testing Error method instead of log Fatal method; fix: remove redundant Sprint ([6c45f04](https://github.com/momentohq/client-sdk-go/commit/6c45f04d76fc770fb843025696a49b9630aa1d46))
* use TestMain to handle setup and teardown ([3bd8740](https://github.com/momentohq/client-sdk-go/commit/3bd87401e3e8a8845e4de10e00fcfa30e1560eb3))
* use unique test keys ([#488](https://github.com/momentohq/client-sdk-go/issues/488)) ([f8fbd5e](https://github.com/momentohq/client-sdk-go/commit/f8fbd5ea7b666ab39793895d66fd42d5cb2d1cf0))


### Miscellaneous

* .gitignore ([6182f15](https://github.com/momentohq/client-sdk-go/commit/6182f15b38456d851e89e223a5af880a001d8c13))
* add `context.Context` as a parameter to public APIs ([#54](https://github.com/momentohq/client-sdk-go/issues/54)) ([a0b71ff](https://github.com/momentohq/client-sdk-go/commit/a0b71ff140b270f316908cce073c6f6675bb7c07))
* add `increment` API ([#308](https://github.com/momentohq/client-sdk-go/issues/308)) ([a72802c](https://github.com/momentohq/client-sdk-go/commit/a72802cb9f3e155020dd7b2a0c1bb5b7e7eab57c))
* Add `NewSimpleCacheClient` snippet with request timeout ([#19](https://github.com/momentohq/client-sdk-go/issues/19)) ([d007347](https://github.com/momentohq/client-sdk-go/commit/d007347bf0d096476192695f3723f888c05e25f1))
* add `WARN`, `ERROR`  and `TRACE` log levels ([#232](https://github.com/momentohq/client-sdk-go/issues/232)) ([fb43e30](https://github.com/momentohq/client-sdk-go/commit/fb43e307449fdb56bfaa36e4241d1da02bf13c08))
* Add a flag and make command to use consistent reads in tests ([#518](https://github.com/momentohq/client-sdk-go/issues/518)) ([79a60df](https://github.com/momentohq/client-sdk-go/commit/79a60dfef3d19003bf7c67836ec829cef032aa4b))
* add a string token provider for use with canary ([#61](https://github.com/momentohq/client-sdk-go/issues/61)) ([9bd0e6d](https://github.com/momentohq/client-sdk-go/commit/9bd0e6dbd67bb8ff614fbb442b15df0080556089))
* add accessor for storage client logger ([#455](https://github.com/momentohq/client-sdk-go/issues/455)) ([27b9907](https://github.com/momentohq/client-sdk-go/commit/27b99075974a00dbb62b012d9de529cebd4d7152))
* add check for context cancel or timeout ([#82](https://github.com/momentohq/client-sdk-go/issues/82)) ([3e0087d](https://github.com/momentohq/client-sdk-go/commit/3e0087dba7fec2d398d03766e3e7af67ca37d5c1))
* add delay between retries when topic subscription limit is reached ([#522](https://github.com/momentohq/client-sdk-go/issues/522)) ([81b9ebc](https://github.com/momentohq/client-sdk-go/commit/81b9ebc27b6456a7811e81b47f134b1bd09cd752))
* add dev docs snippets for storage client ([#457](https://github.com/momentohq/client-sdk-go/issues/457)) ([769cfc8](https://github.com/momentohq/client-sdk-go/commit/769cfc898060c56da7462eb36bae2bb7fd1a8a17))
* add dictionary operations ([#146](https://github.com/momentohq/client-sdk-go/issues/146)) ([642bfd0](https://github.com/momentohq/client-sdk-go/commit/642bfd0d53f3d9e98d326f3ba7e270a16e9dc670))
* add dictionary testing ([#194](https://github.com/momentohq/client-sdk-go/issues/194)) ([ef0f1af](https://github.com/momentohq/client-sdk-go/commit/ef0f1afa9e03a34d87652037f007ef45fd287d95))
* add docs examples for topics and use switch to check message type ([#353](https://github.com/momentohq/client-sdk-go/issues/353)) ([05e3e2e](https://github.com/momentohq/client-sdk-go/commit/05e3e2e68def3e6535f9fb0e2cd3ac954c0b3d9c))
* add docs for automatic doc generation ([#315](https://github.com/momentohq/client-sdk-go/issues/315)) ([651f787](https://github.com/momentohq/client-sdk-go/commit/651f78704d0e4dac8a87f9f005c4f9df5d1aad2e))
* add docs snippets for main api reference page ([#408](https://github.com/momentohq/client-sdk-go/issues/408)) ([7d968de](https://github.com/momentohq/client-sdk-go/commit/7d968de8dbf291a5977f95c2c854765504e9507d))
* Add documentation for momento package ([#21](https://github.com/momentohq/client-sdk-go/issues/21)) ([9f02dcc](https://github.com/momentohq/client-sdk-go/commit/9f02dcc0b2166c2a9d316f3d444fb3175fbf279a))
* add epsilon sleep durations/reduce default ttl on tests ([#500](https://github.com/momentohq/client-sdk-go/issues/500)) ([b1241a9](https://github.com/momentohq/client-sdk-go/commit/b1241a954f9cf9f0f77de68a837eed075561e261))
* add error converter; chore: add simple happy path test ([bdb2154](https://github.com/momentohq/client-sdk-go/commit/bdb2154a1d86b8021a5ccbd3f7e8dc9371ce0bb4))
* add example for ShardedTopicClient ([#417](https://github.com/momentohq/client-sdk-go/issues/417)) ([2159177](https://github.com/momentohq/client-sdk-go/commit/2159177ae10b8913aaa76df4e9d0158986084a36))
* add example instantiating client with cosistent read concern ([#385](https://github.com/momentohq/client-sdk-go/issues/385)) ([a3ba04c](https://github.com/momentohq/client-sdk-go/commit/a3ba04cd7ef582de4ae4dc367896a86a59fdc813))
* add example snippets for setIf APIs ([#383](https://github.com/momentohq/client-sdk-go/issues/383)) ([66a6c36](https://github.com/momentohq/client-sdk-go/commit/66a6c366d25ac353902fee2d6f3cc6fa5af68ada))
* add examples for dictionaries and sets ([#161](https://github.com/momentohq/client-sdk-go/issues/161)) ([3c11792](https://github.com/momentohq/client-sdk-go/commit/3c117924340b2c6f164f291805eff353a5f88921))
* add examples for generating disposable tokens ([#371](https://github.com/momentohq/client-sdk-go/issues/371)) ([e9c4946](https://github.com/momentohq/client-sdk-go/commit/e9c4946424a30fd6016ccb888320a76581d81238))
* add GenerateApiKey dev docs snippet ([#463](https://github.com/momentohq/client-sdk-go/issues/463)) ([97ba7b3](https://github.com/momentohq/client-sdk-go/commit/97ba7b3d3ae49aea1a335ca84637d515a6ab2279))
* Add gRpc timeout ([#10](https://github.com/momentohq/client-sdk-go/issues/10)) ([3fa6041](https://github.com/momentohq/client-sdk-go/commit/3fa60411a20d4aa2f8deb75aa33c0c85c4c6ea55))
* add InvalidInputError ([113fdce](https://github.com/momentohq/client-sdk-go/commit/113fdcecb075efe9d1a090427dfc9a0e49866f98))
* add leaderboards code snippets and example program ([#407](https://github.com/momentohq/client-sdk-go/issues/407)) ([b53824c](https://github.com/momentohq/client-sdk-go/commit/b53824c200da880a29c67b2846107d5b1188477b))
* add lint action for now ([7a9bbd2](https://github.com/momentohq/client-sdk-go/commit/7a9bbd2e05598d7da03fda572c8130fdfd82a25c))
* add list testing ([#202](https://github.com/momentohq/client-sdk-go/issues/202)) ([f42f0cb](https://github.com/momentohq/client-sdk-go/commit/f42f0cbb63492299c890c908e7e3714c5f151dc0))
* add loadgen ([#262](https://github.com/momentohq/client-sdk-go/issues/262)) ([63421e6](https://github.com/momentohq/client-sdk-go/commit/63421e65d74fb5ae7e449ec53090d02b2a876075))
* add logging to failing canary tests ([#499](https://github.com/momentohq/client-sdk-go/issues/499)) ([bfe58f2](https://github.com/momentohq/client-sdk-go/commit/bfe58f21167bcda5ef2c319f00313a3be95c261c))
* add logging to list fetch test ([#506](https://github.com/momentohq/client-sdk-go/issues/506)) ([c2601d1](https://github.com/momentohq/client-sdk-go/commit/c2601d1733a1a9732cc1d913edb9b16aa7df74b7))
* add makefile targets to update and build protos ([#525](https://github.com/momentohq/client-sdk-go/issues/525)) ([86b747d](https://github.com/momentohq/client-sdk-go/commit/86b747d002f9489eb184165901c079f80e8e4389))
* Add more detailed tests ([#20](https://github.com/momentohq/client-sdk-go/issues/20)) ([bf00733](https://github.com/momentohq/client-sdk-go/commit/bf0073352ced6e1581cd397c2c8096c2b2ed600d))
* add more documentation ([#256](https://github.com/momentohq/client-sdk-go/issues/256)) ([a855e51](https://github.com/momentohq/client-sdk-go/commit/a855e51fc69b58286f6b93f7da55c7cb8dc7799c))
* add per-service makefile targets ([#481](https://github.com/momentohq/client-sdk-go/issues/481)) ([4e758ce](https://github.com/momentohq/client-sdk-go/commit/4e758ce806769f0276b732bcde17b31e6ee2850d))
* add README generator ([#30](https://github.com/momentohq/client-sdk-go/issues/30)) ([a0b2c6d](https://github.com/momentohq/client-sdk-go/commit/a0b2c6d23d70ce4dffa4e7795cfce256fb2b80cf))
* add SDK version and runtime version headers for store client ([#436](https://github.com/momentohq/client-sdk-go/issues/436)) ([05b7b53](https://github.com/momentohq/client-sdk-go/commit/05b7b53d7d299046a1d7038be691eb373428cde8))
* add set operations ([#148](https://github.com/momentohq/client-sdk-go/issues/148)) ([6c350b3](https://github.com/momentohq/client-sdk-go/commit/6c350b314ffbf65e6355a16f501a65af8e4fde59))
* add single msg recv ([#88](https://github.com/momentohq/client-sdk-go/issues/88)) ([f9a9bb0](https://github.com/momentohq/client-sdk-go/commit/f9a9bb0df01c7b7e4354b2608cd2a14a5a583119))
* add sorted set get score ([#269](https://github.com/momentohq/client-sdk-go/issues/269)) ([c8668d6](https://github.com/momentohq/client-sdk-go/commit/c8668d67e1b0410ad81bce980919985572baa9db))
* add strings to Describe()s and It()s for test segmentation ([#444](https://github.com/momentohq/client-sdk-go/issues/444)) ([a8f5e04](https://github.com/momentohq/client-sdk-go/commit/a8f5e04d440b4a91ceb0ef39b77039b82e29b4e0))
* add testing for newly added methods ([#272](https://github.com/momentohq/client-sdk-go/issues/272)) ([ccbbba1](https://github.com/momentohq/client-sdk-go/commit/ccbbba1a1c94a798fbac11e7ef51bde125b97acd))
* add topics loadgen ([#341](https://github.com/momentohq/client-sdk-go/issues/341)) ([10f7e5b](https://github.com/momentohq/client-sdk-go/commit/10f7e5b9046054e99d10d71574df3aa9036f2947))
* adding list operations ([#89](https://github.com/momentohq/client-sdk-go/issues/89)) ([a70447b](https://github.com/momentohq/client-sdk-go/commit/a70447b5f7bf784da1a4bb606f2aaa9e216fbc83))
* adds code scanning workflow ([#34](https://github.com/momentohq/client-sdk-go/issues/34)) ([6bc38aa](https://github.com/momentohq/client-sdk-go/commit/6bc38aa2a4d8d74388e336fd60525726937667f9))
* adds sorted set increment in incubating ([#106](https://github.com/momentohq/client-sdk-go/issues/106)) ([19717ae](https://github.com/momentohq/client-sdk-go/commit/19717ae2a47c401454d073ccd96ae247377d1442))
* audit and test for nils ([#206](https://github.com/momentohq/client-sdk-go/issues/206)) ([8b93cc0](https://github.com/momentohq/client-sdk-go/commit/8b93cc081e5a1f262595945a20c99b8915e1f87f))
* better debug logging on discontinuities ([#503](https://github.com/momentohq/client-sdk-go/issues/503)) ([0133f21](https://github.com/momentohq/client-sdk-go/commit/0133f2124dc1c52bde70707ca5e3200d2cc15a86))
* bumps examples to latest sdk ([#374](https://github.com/momentohq/client-sdk-go/issues/374)) ([33331aa](https://github.com/momentohq/client-sdk-go/commit/33331aaa317c43c6597f6f4786ada92e5ca6972f))
* bumps examples to latest version ([#140](https://github.com/momentohq/client-sdk-go/issues/140)) ([19150c5](https://github.com/momentohq/client-sdk-go/commit/19150c50f1665846a35524be2768b6f32bbe02cd))
* bumps examples to latest version ([#158](https://github.com/momentohq/client-sdk-go/issues/158)) ([3df36fb](https://github.com/momentohq/client-sdk-go/commit/3df36fbcd2e347a6f407ebe7e82de066bd7f226f))
* by default support lazy connections ([#388](https://github.com/momentohq/client-sdk-go/issues/388)) ([ce01397](https://github.com/momentohq/client-sdk-go/commit/ce013976541dee758582a8de7e18a0292a3bf023))
* change stability badge from experimental to alpha ([#77](https://github.com/momentohq/client-sdk-go/issues/77)) ([258c459](https://github.com/momentohq/client-sdk-go/commit/258c45932446c83b66a8c4f9dea206c666c0b69b))
* clean up and return new sdk ConnectionError when eager connection fails ([#398](https://github.com/momentohq/client-sdk-go/issues/398)) ([0d2cce4](https://github.com/momentohq/client-sdk-go/commit/0d2cce48da7dba76cf4ebc479ad3f679cfccb9a2))
* clean up and try to speed up tests ([#477](https://github.com/momentohq/client-sdk-go/issues/477)) ([e2a9870](https://github.com/momentohq/client-sdk-go/commit/e2a987047625ae04368606520678414bab025c3b))
* consolidate SortedSetPutElement to SortedSetElement, add util ([#279](https://github.com/momentohq/client-sdk-go/issues/279)) ([4371765](https://github.com/momentohq/client-sdk-go/commit/43717654c65b7e9ee21d9e5e847aeb19b656cadd))
* convert cache client constructor to accept args ([#209](https://github.com/momentohq/client-sdk-go/issues/209)) ([dd8c6f4](https://github.com/momentohq/client-sdk-go/commit/dd8c6f433f0db010f4fa640a77d802be8b069cfa))
* copy over new protos ([#358](https://github.com/momentohq/client-sdk-go/issues/358)) ([1df6372](https://github.com/momentohq/client-sdk-go/commit/1df63724445980ab903b00e99834c476a1bf3cb3))
* copy protos for new APIs to go sdk and fix MakeFile ([#332](https://github.com/momentohq/client-sdk-go/issues/332)) ([3bf20d2](https://github.com/momentohq/client-sdk-go/commit/3bf20d28df89f4efa903708c5a87aea532468347))
* correct usage of cancel function ([#83](https://github.com/momentohq/client-sdk-go/issues/83)) ([1d7f796](https://github.com/momentohq/client-sdk-go/commit/1d7f7968ef06567fa8c784e964b63d96b8529d16))
* debug canary tests ([#494](https://github.com/momentohq/client-sdk-go/issues/494)) ([8a5d399](https://github.com/momentohq/client-sdk-go/commit/8a5d39922d95bc3b70f729f7bca2f95bc96fe5d4))
* debug sorted fetch by rank - remove elements test ([#509](https://github.com/momentohq/client-sdk-go/issues/509)) ([46ac38b](https://github.com/momentohq/client-sdk-go/commit/46ac38b59897213adea49eb865d0673c0122bc81))
* default cache client testing ([#299](https://github.com/momentohq/client-sdk-go/issues/299)) ([2edf56a](https://github.com/momentohq/client-sdk-go/commit/2edf56a459e8b1a0021c68e65709d2738320ae16))
* dep scan workflow ([#35](https://github.com/momentohq/client-sdk-go/issues/35)) ([20edc27](https://github.com/momentohq/client-sdk-go/commit/20edc2734cfc5fc54374fec97177de33080c918b))
* **deps-dev:** bump @babel/traverse ([#364](https://github.com/momentohq/client-sdk-go/issues/364)) ([f743d16](https://github.com/momentohq/client-sdk-go/commit/f743d1644af621f66b87048752f7c38c8a010372))
* **deps-dev:** bump braces in /examples/aws-lambda/infrastructure ([#441](https://github.com/momentohq/client-sdk-go/issues/441)) ([f440d6d](https://github.com/momentohq/client-sdk-go/commit/f440d6d868b9fe0332fd8325bdde9bce51494ed2))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#365](https://github.com/momentohq/client-sdk-go/issues/365)) ([b5e716f](https://github.com/momentohq/client-sdk-go/commit/b5e716f3e26b01b9031b004b75852aba7e3a8288))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#389](https://github.com/momentohq/client-sdk-go/issues/389)) ([2e5efa6](https://github.com/momentohq/client-sdk-go/commit/2e5efa6907fbc4f6eaab11e914dafd4eb43cb112))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#395](https://github.com/momentohq/client-sdk-go/issues/395)) ([01d60a2](https://github.com/momentohq/client-sdk-go/commit/01d60a2214c891a9a7300eba390cdc8c0338fed5))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#403](https://github.com/momentohq/client-sdk-go/issues/403)) ([7858197](https://github.com/momentohq/client-sdk-go/commit/7858197b97c4574407c5eb9caa013bc389760975))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#437](https://github.com/momentohq/client-sdk-go/issues/437)) ([9610373](https://github.com/momentohq/client-sdk-go/commit/9610373e17a9b751d8e0fae6cf15add65ab076fa))
* **deps:** bump github.com/momentohq/client-sdk-go in /examples ([#474](https://github.com/momentohq/client-sdk-go/issues/474)) ([55be217](https://github.com/momentohq/client-sdk-go/commit/55be21725315cc60fa3aad8a557c9fb2dc3ea3d1))
* **deps:** bump golang.org/x/net from 0.21.0 to 0.23.0 ([#442](https://github.com/momentohq/client-sdk-go/issues/442)) ([03bc194](https://github.com/momentohq/client-sdk-go/commit/03bc194785acbc1653dd15264848f005e7a50cd7))
* **deps:** bump golang.org/x/net in /examples/aws-lambda/lambda ([#450](https://github.com/momentohq/client-sdk-go/issues/450)) ([ce9ce33](https://github.com/momentohq/client-sdk-go/commit/ce9ce33ffbdf66e2f2a32e81eccd8b5c93a33a32))
* **deps:** bump google.golang.org/protobuf ([#443](https://github.com/momentohq/client-sdk-go/issues/443)) ([d80b634](https://github.com/momentohq/client-sdk-go/commit/d80b63417b533d4ceb65e294e4ae3d4916e5d44b))
* disable storage tests ([#429](https://github.com/momentohq/client-sdk-go/issues/429)) ([b3e26c3](https://github.com/momentohq/client-sdk-go/commit/b3e26c3a373f59ff2612814f6fcebec1fede166e))
* eager connection should fail fast ([#396](https://github.com/momentohq/client-sdk-go/issues/396)) ([849daca](https://github.com/momentohq/client-sdk-go/commit/849daca811a9d9656050f252dfefa6586ff36909))
* empty commit to bump version ([#327](https://github.com/momentohq/client-sdk-go/issues/327)) ([2bfc6b5](https://github.com/momentohq/client-sdk-go/commit/2bfc6b57614e9dc68f3539c3efce8d181e355ed8))
* enable consistent reads on all canary test targets ([#523](https://github.com/momentohq/client-sdk-go/issues/523)) ([a6948b9](https://github.com/momentohq/client-sdk-go/commit/a6948b98e18ec458f64591c0a59ef9a6c3e0fc7b))
* find and replace all instances of MOMENTO_AUTH_TOKEN with MOMENTO_API_KEY ([#354](https://github.com/momentohq/client-sdk-go/issues/354)) ([b8968cd](https://github.com/momentohq/client-sdk-go/commit/b8968cd1cd36bbcc976fbf5bff323c12cc4c18e6))
* fix example for clarity ([#453](https://github.com/momentohq/client-sdk-go/issues/453)) ([30511d0](https://github.com/momentohq/client-sdk-go/commit/30511d0eda679d94d7e7e233c765a9a853d64e17))
* fix formatting in Makefile ([#535](https://github.com/momentohq/client-sdk-go/issues/535)) ([d97f229](https://github.com/momentohq/client-sdk-go/commit/d97f229a0c8b269f612b3101b5f3674c8554aef7))
* fix push to main gh workflow ([#64](https://github.com/momentohq/client-sdk-go/issues/64)) ([9d7f6a7](https://github.com/momentohq/client-sdk-go/commit/9d7f6a70caf3dfeb029b145094af397d20c33511))
* hasScalarTTL -&gt; hasTTL ([#153](https://github.com/momentohq/client-sdk-go/issues/153)) ([a377d48](https://github.com/momentohq/client-sdk-go/commit/a377d488229a675e70132efb077403278eb15cc8))
* improve logging if we exceed the max number of concurrent streams ([#520](https://github.com/momentohq/client-sdk-go/issues/520)) ([2cc3ec2](https://github.com/momentohq/client-sdk-go/commit/2cc3ec2be2e686870790c9fb355d027da2c78cbe))
* Improve the build process ([#91](https://github.com/momentohq/client-sdk-go/issues/91)) ([ec4b262](https://github.com/momentohq/client-sdk-go/commit/ec4b262bd53355c2265e7c8fcf4d004f4aae7502))
* increase code readability by updating variable names ([cb27ad4](https://github.com/momentohq/client-sdk-go/commit/cb27ad46c8f1c3ffc6fc15e97f279edcf9656f87))
* Lint the example code. ([#218](https://github.com/momentohq/client-sdk-go/issues/218)) ([d8d736f](https://github.com/momentohq/client-sdk-go/commit/d8d736f07b32f705253bec25ca64a720e5951129))
* list ops part 2 ([#102](https://github.com/momentohq/client-sdk-go/issues/102)) ([47e8265](https://github.com/momentohq/client-sdk-go/commit/47e8265c6a8ca02853874ce277bc129c1d142cb3))
* log whether consistent reads are enabled for integration tests ([#519](https://github.com/momentohq/client-sdk-go/issues/519)) ([22af385](https://github.com/momentohq/client-sdk-go/commit/22af385fe32439ea83aca465844176ba08229131))
* **main:** release 1.24.0 ([#434](https://github.com/momentohq/client-sdk-go/issues/434)) ([b5f8a33](https://github.com/momentohq/client-sdk-go/commit/b5f8a33912abff962f02e6f4f25cb19923c17148))
* **main:** release 1.24.1 ([#438](https://github.com/momentohq/client-sdk-go/issues/438)) ([71a00ba](https://github.com/momentohq/client-sdk-go/commit/71a00ba045d2113c4b041c2736e80c3ed5b19e2b))
* **main:** release 1.24.2 ([#454](https://github.com/momentohq/client-sdk-go/issues/454)) ([5d6fbee](https://github.com/momentohq/client-sdk-go/commit/5d6fbeeaee4f8e46743bb0eba4c5206749c7257a))
* **main:** release 1.25.0 ([#459](https://github.com/momentohq/client-sdk-go/issues/459)) ([6ca7811](https://github.com/momentohq/client-sdk-go/commit/6ca7811fda5fd5ef338869990e9c83dab63289b2))
* **main:** release 1.26.0 ([#466](https://github.com/momentohq/client-sdk-go/issues/466)) ([6cbec2a](https://github.com/momentohq/client-sdk-go/commit/6cbec2a3e3195cbf502fd26c88dcddf184f24a3d))
* **main:** release 1.26.1 ([#468](https://github.com/momentohq/client-sdk-go/issues/468)) ([ad9e266](https://github.com/momentohq/client-sdk-go/commit/ad9e266700de7df9d371272069980701dfe0a15a))
* **main:** release 1.26.2 ([#470](https://github.com/momentohq/client-sdk-go/issues/470)) ([90b84ad](https://github.com/momentohq/client-sdk-go/commit/90b84add73f2f11e017049a84af059d5984f4ded))
* **main:** release 1.27.0 ([#475](https://github.com/momentohq/client-sdk-go/issues/475)) ([920344d](https://github.com/momentohq/client-sdk-go/commit/920344dcf9591b042a140a8d9c8cca323e172e01))
* **main:** release 1.27.1 ([#485](https://github.com/momentohq/client-sdk-go/issues/485)) ([f20f4b8](https://github.com/momentohq/client-sdk-go/commit/f20f4b812f5a188220eff742802dbd8bbda1eea5))
* **main:** release 1.27.2 ([#492](https://github.com/momentohq/client-sdk-go/issues/492)) ([930393f](https://github.com/momentohq/client-sdk-go/commit/930393f95317767cbec985ba375e73f4b640c240))
* **main:** release 1.27.3 ([#495](https://github.com/momentohq/client-sdk-go/issues/495)) ([bdc2450](https://github.com/momentohq/client-sdk-go/commit/bdc2450a53149fc08fd6184da85dcf5bf86dd3cf))
* **main:** release 1.27.4 ([#498](https://github.com/momentohq/client-sdk-go/issues/498)) ([8c175f4](https://github.com/momentohq/client-sdk-go/commit/8c175f4611037f4c0a5140ecea5d4e40eb1c28d3))
* **main:** release 1.27.5 ([#501](https://github.com/momentohq/client-sdk-go/issues/501)) ([ec44072](https://github.com/momentohq/client-sdk-go/commit/ec440728502bb5b09f404fb7cc5e409efac066cf))
* **main:** release 1.27.6 ([#504](https://github.com/momentohq/client-sdk-go/issues/504)) ([6dec2a6](https://github.com/momentohq/client-sdk-go/commit/6dec2a673210c2bc11e37af58783be436f99da2b))
* **main:** release 1.27.7 ([#508](https://github.com/momentohq/client-sdk-go/issues/508)) ([7aa93fe](https://github.com/momentohq/client-sdk-go/commit/7aa93fe2b5baeedd916dd2cc4072b897085741a5))
* **main:** release 1.28.0 ([#512](https://github.com/momentohq/client-sdk-go/issues/512)) ([9956f7f](https://github.com/momentohq/client-sdk-go/commit/9956f7fd32ab52e4209540657fd21830b2d432bc))
* **main:** release 1.28.1 ([#514](https://github.com/momentohq/client-sdk-go/issues/514)) ([0b66c45](https://github.com/momentohq/client-sdk-go/commit/0b66c451c8ce6726f47728cef3121912ced8535a))
* **main:** release 1.28.2 ([#516](https://github.com/momentohq/client-sdk-go/issues/516)) ([39f9580](https://github.com/momentohq/client-sdk-go/commit/39f9580181ebd775cf7975b715b275f628f2ef71))
* **main:** release 1.28.3 ([#524](https://github.com/momentohq/client-sdk-go/issues/524)) ([ca9018f](https://github.com/momentohq/client-sdk-go/commit/ca9018fc85b2773aa5e4711f29b6c5d15da7ace2))
* **main:** release 1.28.4 ([#536](https://github.com/momentohq/client-sdk-go/issues/536)) ([50bd958](https://github.com/momentohq/client-sdk-go/commit/50bd958c55145eb00692da90a8800946994c2a32))
* **main:** release 1.28.5 ([#539](https://github.com/momentohq/client-sdk-go/issues/539)) ([bbe97f5](https://github.com/momentohq/client-sdk-go/commit/bbe97f5e5bc6c45a391d37495d509d903680dc55))
* make `CollectionTtl` pointer and assign default value to  `Ttl` and `RefreshTtl` ([#212](https://github.com/momentohq/client-sdk-go/issues/212)) ([6e92b7f](https://github.com/momentohq/client-sdk-go/commit/6e92b7fef094184f7b90a44c86ddccc4a753d6e5))
* make set value available for SetCacheResponse ([00390fa](https://github.com/momentohq/client-sdk-go/commit/00390fa3c8ea8aa33277b58297ab3e570bb54f2a))
* make sure all example methods accept momento.Value args ([#214](https://github.com/momentohq/client-sdk-go/issues/214)) ([3388ea2](https://github.com/momentohq/client-sdk-go/commit/3388ea2261b3f4fe95f2dbaa0726013183184a4a))
* make sure push-to-main has test session token too ([#462](https://github.com/momentohq/client-sdk-go/issues/462)) ([268e771](https://github.com/momentohq/client-sdk-go/commit/268e7713ad2be16dd9f595f37d41e4a1c4547d65))
* makes progress standardizing our workflows ([#33](https://github.com/momentohq/client-sdk-go/issues/33)) ([7fb3bbb](https://github.com/momentohq/client-sdk-go/commit/7fb3bbbb3b7926b8e37c302cc2818a7f426c8f5f))
* minor tweaks to logging inside of topic client ([#513](https://github.com/momentohq/client-sdk-go/issues/513)) ([fccfc5e](https://github.com/momentohq/client-sdk-go/commit/fccfc5e3e7a260c2ffb6c899d9e9e583d8bf3f34))
* move golang examples from `client-sdk-examples` repo ([#29](https://github.com/momentohq/client-sdk-go/issues/29)) ([57e3436](https://github.com/momentohq/client-sdk-go/commit/57e3436e5e390abd090e7756a9501fd7dec96553))
* move responses into their own package ([#222](https://github.com/momentohq/client-sdk-go/issues/222)) ([a1e94cd](https://github.com/momentohq/client-sdk-go/commit/a1e94cd5aa7b6f2e58c56e5cb2c1e8fd23237a4e))
* organize imports ([5dd3f28](https://github.com/momentohq/client-sdk-go/commit/5dd3f2870cae56b702592df8af9d05bfe3c47eb6))
* **protos:** update protos to v0.119.0 and regenerate code ([#533](https://github.com/momentohq/client-sdk-go/issues/533)) ([f1ac517](https://github.com/momentohq/client-sdk-go/commit/f1ac5177f57799a697cd5a05e014a6b09b7e6a34))
* pull set name into standard validators ([#103](https://github.com/momentohq/client-sdk-go/issues/103)) ([21bd2fd](https://github.com/momentohq/client-sdk-go/commit/21bd2fd055dceeb7947f4ef18531b75bb0020929))
* reduce eager connection timeout for test for low latency test environments ([#537](https://github.com/momentohq/client-sdk-go/issues/537)) ([9b61f23](https://github.com/momentohq/client-sdk-go/commit/9b61f238e84a7e4485fff1b2092800894409c104))
* reduce request timeout for timeout test ([#431](https://github.com/momentohq/client-sdk-go/issues/431)) ([cd522fb](https://github.com/momentohq/client-sdk-go/commit/cd522fb5d7c23ce904ea7cb1c8ad104ab9ebb348))
* remove express ReadConcern ([#423](https://github.com/momentohq/client-sdk-go/issues/423)) ([af2cd8b](https://github.com/momentohq/client-sdk-go/commit/af2cd8b28076470abe747a7f6c710d41f2ebd2d5))
* remove NextToken from ListCachesRequest ([#449](https://github.com/momentohq/client-sdk-go/issues/449)) ([538add1](https://github.com/momentohq/client-sdk-go/commit/538add16f6c759fdf27ac244685331cf53fdbe1c))
* remove publisher id from simple topics example and use only one polling function at a time ([#489](https://github.com/momentohq/client-sdk-go/issues/489)) ([d97e718](https://github.com/momentohq/client-sdk-go/commit/d97e7182b473e92349f2d9e7fad5c6e3a51e6bc7))
* remove stary commented code ([#225](https://github.com/momentohq/client-sdk-go/issues/225)) ([b1ebdd1](https://github.com/momentohq/client-sdk-go/commit/b1ebdd1154e3b2ea83bea1210996df69efc15c60))
* remove tests that require account session token ([#472](https://github.com/momentohq/client-sdk-go/issues/472)) ([d574a74](https://github.com/momentohq/client-sdk-go/commit/d574a7457581f6df8776d4ce936f751e1a9d1783))
* removing parallel tests ([#79](https://github.com/momentohq/client-sdk-go/issues/79)) ([8dc24b8](https://github.com/momentohq/client-sdk-go/commit/8dc24b89fab2c41981dfd8840156c8916c017d8b))
* rename CollectionTtl to Ttl in request types ([#208](https://github.com/momentohq/client-sdk-go/issues/208)) ([e92c675](https://github.com/momentohq/client-sdk-go/commit/e92c675610fede275ee8ecdf9125abcbf8d6bb14))
* rename ElementsFromMap* to DictionaryElementsFromMap* ([#278](https://github.com/momentohq/client-sdk-go/issues/278)) ([2d7e54e](https://github.com/momentohq/client-sdk-go/commit/2d7e54e19920b1a318cdf76755fbbda0583b7711))
* rename packages to be more generic ([eebfb65](https://github.com/momentohq/client-sdk-go/commit/eebfb65a5c2cecff4de6b4689e06e10e24a6fa41))
* revert `1e7c506` ([#373](https://github.com/momentohq/client-sdk-go/issues/373)) ([5a82fa9](https://github.com/momentohq/client-sdk-go/commit/5a82fa935ea807f914c5faa520070354fc2d7fb0))
* revert Lint the example code ([#219](https://github.com/momentohq/client-sdk-go/issues/219)) ([8618c44](https://github.com/momentohq/client-sdk-go/commit/8618c44c86596d6f653dbf7087ccfd62cda85079))
* run `examples` in `on-pull-request` workflow ([#289](https://github.com/momentohq/client-sdk-go/issues/289)) ([c7de675](https://github.com/momentohq/client-sdk-go/commit/c7de675d8b8064cd621c72dfb7b776a2aa06071f))
* run ci on ubuntu-24.04 ([#412](https://github.com/momentohq/client-sdk-go/issues/412)) ([ffe3ed4](https://github.com/momentohq/client-sdk-go/commit/ffe3ed4aa4f9e972cd5e3758dd2d3ce35f9fc02c))
* set topics resubscribe delay to 500ms ([#490](https://github.com/momentohq/client-sdk-go/issues/490)) ([8d70251](https://github.com/momentohq/client-sdk-go/commit/8d702518c8ea6fc2fb681b530aba6bb8ff7f11b2))
* simplify logic of batchutils functions ([#392](https://github.com/momentohq/client-sdk-go/issues/392)) ([a3bc311](https://github.com/momentohq/client-sdk-go/commit/a3bc3114f301de81f6b4f923e90fd3bf54610a50))
* sorted set testing followup ([#288](https://github.com/momentohq/client-sdk-go/issues/288)) ([bd87e16](https://github.com/momentohq/client-sdk-go/commit/bd87e16d9302b7f9971a4c94edc1d5630fff3fed))
* standardize receivers for request type functions ([#130](https://github.com/momentohq/client-sdk-go/issues/130)) ([4985173](https://github.com/momentohq/client-sdk-go/commit/49851737f533d0f0e76ef75e115f45195621e17c))
* switch from TEST_AUTH_TOKEN to MOMENTO_API_KEY when running tests ([#421](https://github.com/momentohq/client-sdk-go/issues/421)) ([6d40a32](https://github.com/momentohq/client-sdk-go/commit/6d40a328ff9d0526fd20b5aa9850dd3c1cf1f1dd))
* try to make flaky tests less flaky ([#482](https://github.com/momentohq/client-sdk-go/issues/482)) ([10500ea](https://github.com/momentohq/client-sdk-go/commit/10500ea2014a3ddeae45bc665a28a7f9cc1f9677))
* uncomment storage client tests ([#448](https://github.com/momentohq/client-sdk-go/issues/448)) ([f8c0ca1](https://github.com/momentohq/client-sdk-go/commit/f8c0ca1a17a3162ab7e5459b6f947d161fbfb61f))
* update `README` template ([#85](https://github.com/momentohq/client-sdk-go/issues/85)) ([49efef7](https://github.com/momentohq/client-sdk-go/commit/49efef7de1a3523f343130e6c38a23c3fed78770))
* update a few names ([#203](https://github.com/momentohq/client-sdk-go/issues/203)) ([c1e934a](https://github.com/momentohq/client-sdk-go/commit/c1e934abf8d09e1f26c43cc57636ccd012483c9c))
* update examples ([#229](https://github.com/momentohq/client-sdk-go/issues/229)) ([2d8bcd5](https://github.com/momentohq/client-sdk-go/commit/2d8bcd584dadd425793b53398328391e2f47839c))
* update examples for v1.4.0 ([#314](https://github.com/momentohq/client-sdk-go/issues/314)) ([f4f6ec9](https://github.com/momentohq/client-sdk-go/commit/f4f6ec98c1244663608363b8821bcbf3c4a96524))
* update examples latest sdk version ([#303](https://github.com/momentohq/client-sdk-go/issues/303)) ([8b64528](https://github.com/momentohq/client-sdk-go/commit/8b645284861d4c46ed86f26d091d0e579dea646a))
* Update examples to 1.5.1 ([#324](https://github.com/momentohq/client-sdk-go/issues/324)) ([7efdf57](https://github.com/momentohq/client-sdk-go/commit/7efdf5792d21681bedd06dd07e093d54072fcfb1))
* update examples to use cache client factory method with eager connections by default ([#344](https://github.com/momentohq/client-sdk-go/issues/344)) ([65f5565](https://github.com/momentohq/client-sdk-go/commit/65f556596b0d97f66d8b768c8f3954c5999ef11f))
* update examples to v0.16.0 ([#282](https://github.com/momentohq/client-sdk-go/issues/282)) ([1b36fbe](https://github.com/momentohq/client-sdk-go/commit/1b36fbeef3a3071fb58863a1ecbbc49a39d4a0b8))
* update get responses ([#433](https://github.com/momentohq/client-sdk-go/issues/433)) ([06b59dc](https://github.com/momentohq/client-sdk-go/commit/06b59dc67ef31ce1dd1e2ceeefd16cb81e5cf359))
* update get/set perf test + loadgen test ([#410](https://github.com/momentohq/client-sdk-go/issues/410)) ([b973c77](https://github.com/momentohq/client-sdk-go/commit/b973c7776a4214151536abeeda22cdd863c0476f))
* update go topics examples ([#483](https://github.com/momentohq/client-sdk-go/issues/483)) ([ff0332f](https://github.com/momentohq/client-sdk-go/commit/ff0332f2d0efa2d9deb82915cc7ed9988f35bc65))
* update go.mod ([df55dd1](https://github.com/momentohq/client-sdk-go/commit/df55dd1a38dfb3e8da7a913f77111abd1d8741e3))
* update install instructions ([#300](https://github.com/momentohq/client-sdk-go/issues/300)) ([1e7c506](https://github.com/momentohq/client-sdk-go/commit/1e7c506c0888c46557b9b32b409ef7adb0832642))
* update loadgen to be more specific for numberOfConcurrentRequests caveat ([#363](https://github.com/momentohq/client-sdk-go/issues/363)) ([cd6b46d](https://github.com/momentohq/client-sdk-go/commit/cd6b46d15673e8a7d2275d97fc65bc2fc734bfa7))
* Update main.go ([#285](https://github.com/momentohq/client-sdk-go/issues/285)) ([9f153b9](https://github.com/momentohq/client-sdk-go/commit/9f153b9e9e63706ddbdfc39b090e5692f2af4c33))
* update metadata message signifying item not found ([#427](https://github.com/momentohq/client-sdk-go/issues/427)) ([ffbe5d3](https://github.com/momentohq/client-sdk-go/commit/ffbe5d36116d09afa0743898ea412ca112736b58))
* update README ([7e1e6f8](https://github.com/momentohq/client-sdk-go/commit/7e1e6f87d2b7d3383f699daea4a0caf5faa4698e))
* update README ([207a3a8](https://github.com/momentohq/client-sdk-go/commit/207a3a81bcb445882916f1842789c30e9c29680c))
* Update readme ([#13](https://github.com/momentohq/client-sdk-go/issues/13)) ([1c6b7aa](https://github.com/momentohq/client-sdk-go/commit/1c6b7aa52007bffa3c080b9c56771349ff9e51b6))
* Update README.md ([#295](https://github.com/momentohq/client-sdk-go/issues/295)) ([06e1e1c](https://github.com/momentohq/client-sdk-go/commit/06e1e1c72a04553d245ad9e8602ac81f26065e60))
* update SDK version for examples ([#275](https://github.com/momentohq/client-sdk-go/issues/275)) ([676e53f](https://github.com/momentohq/client-sdk-go/commit/676e53fa0660c6246c11b79449e536f2d8e1995b))
* update to the latest protos ([#9](https://github.com/momentohq/client-sdk-go/issues/9)) ([f315a7b](https://github.com/momentohq/client-sdk-go/commit/f315a7b6dd15787d87a6f18835197db048bc6965))
* update topics example to illustrate how to wait for subscription close ([#511](https://github.com/momentohq/client-sdk-go/issues/511)) ([38843bf](https://github.com/momentohq/client-sdk-go/commit/38843bffef807ecdc29f4ba86761ba8202107767))
* update word `item(s)` to `element(s)` ([#105](https://github.com/momentohq/client-sdk-go/issues/105)) ([8d5b73b](https://github.com/momentohq/client-sdk-go/commit/8d5b73b575985fe6a4f3cbbcfeb44d78f7b9cc37))
* updates examples to latest sdk spec ([#59](https://github.com/momentohq/client-sdk-go/issues/59)) ([276ed94](https://github.com/momentohq/client-sdk-go/commit/276ed94ca6f1812e05d7efc793f31ff130ae310c))
* updates examples to latest sdk version ([#109](https://github.com/momentohq/client-sdk-go/issues/109)) ([9bc6a86](https://github.com/momentohq/client-sdk-go/commit/9bc6a86cd5ad56275886ceb83d4d0c6787b801ee))
* updates readme template to latest spec ([#32](https://github.com/momentohq/client-sdk-go/issues/32)) ([f001cea](https://github.com/momentohq/client-sdk-go/commit/f001ceabb65f82c922c41166f001cb8a8af9f87e))
* updates sorted set examples ([#92](https://github.com/momentohq/client-sdk-go/issues/92)) ([84ed00a](https://github.com/momentohq/client-sdk-go/commit/84ed00ae75cfdc26874ef3d8c0a1a8de0b691e58))
* Use the consistent reads flag in integration tests ([#515](https://github.com/momentohq/client-sdk-go/issues/515)) ([0f25f1c](https://github.com/momentohq/client-sdk-go/commit/0f25f1c89d964ddc431bf1e0a00980b26f24c64c))
* use updated repo generator v2 ([#316](https://github.com/momentohq/client-sdk-go/issues/316)) ([f4927c0](https://github.com/momentohq/client-sdk-go/commit/f4927c0c6ae6b84a8585ea2258b8c3d9e427b472))
* verify pointer responses ([#231](https://github.com/momentohq/client-sdk-go/issues/231)) ([9388e2b](https://github.com/momentohq/client-sdk-go/commit/9388e2b6cd16c4ab45dbc5651564232bb903015e))
* version fix ([#290](https://github.com/momentohq/client-sdk-go/issues/290)) ([43d2ae4](https://github.com/momentohq/client-sdk-go/commit/43d2ae46ac27abbd43d434a301932ae8c1a7e487))

## [1.28.5](https://github.com/momentohq/client-sdk-go/compare/v1.28.4...v1.28.5) (2024-10-24)


### Bug Fixes

* interpret get rank response when sorted set found ([#541](https://github.com/momentohq/client-sdk-go/issues/541)) ([ffe0e9c](https://github.com/momentohq/client-sdk-go/commit/ffe0e9cc511b44879ffc5d19d543793eaa708dd7))


### Miscellaneous

* add delay between retries when topic subscription limit is reached ([#522](https://github.com/momentohq/client-sdk-go/issues/522)) ([81b9ebc](https://github.com/momentohq/client-sdk-go/commit/81b9ebc27b6456a7811e81b47f134b1bd09cd752))
* reduce eager connection timeout for test for low latency test environments ([#537](https://github.com/momentohq/client-sdk-go/issues/537)) ([9b61f23](https://github.com/momentohq/client-sdk-go/commit/9b61f238e84a7e4485fff1b2092800894409c104))

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
