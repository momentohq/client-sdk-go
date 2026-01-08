# Sharded Topic Example

If you have a Momento Topic for which you expect an extremely high volume of traffic (> 100k publishes per second),
you will likely get better performance by sharding the traffic into multiple topics. This allows for better load distribution
on Momento's servers.

In this directory we provide an example of an alternate implementation of the Momento `TopicClient` interface that 
will split your topic up into multiple shards while still retaining the same API.

All you need to do is to update the line of code where you instantiate your `TopicClient` from this:

```go
topicClient, err := momento.NewTopicClient(config, credentialsProvider)
```

to this:

```go
numShardsPerTopic := 16
topicClient, err := NewShardedTopicClient(config, credentialsProvider, numShardsPerTopic)
```

When you `Publish` to the `ShardedTopicClient`, it will treat the topic name that you specify as a prefix, and create
multiple topics suffixed with `-0`, `-1`, etc. up to the number of shards you specify. When you create a subscription,
it will use goroutines to subscribe to all of these topics and emit messages from all of them.

Caveats:
* If you have some topics that you don't wish to shard, and some that you do, you should create separate `TopicClient`
  and `ShardedTopicClient` instances accordingly.
* For any topic that you do wish to shard, it's important that you use a `ShardedTopicClient` for both the `Publish` and
  the `Subscribe` operations. Otherwise you may end up with messages being lost.
* It's important to make sure you use the same value for `numShardsPerTopic` in all of the places where you are doing
  any `Publish` and `Subscribe` operations on the same topic. If you don't, you may end up with messages appearing to
  be lost due to the clients not using the same list of sharded topic names behind the scenes.

## Requirements.

- [Go](https://go.dev/dl/)
- A Momento API key is required, you can generate one using the [Momento Console](https://console.gomomento.com/api-keys)
- A Momento service endpoint is required. You can find a [list of them here](https://docs.momentohq.com/platform/regions)
- Run `go mod vendor` to install dependencies.

## Running the example

```
MOMENTO_API_KEY=<YOUR_TOKEN> MOMENTO_ENDPOINT=<endpoint> go run topic-sharded-example/*.go
```

If you compare the code in this [main.go](./main.go) to the code in the basic topic example's [main.go](../topic-example/main.go),
you'll note that they are almost identical. The only meaningful difference is the way the `TopicClient` is instantiated,
in the `getShardedTopicClient` and `getTopicClient` functions respectively. 
