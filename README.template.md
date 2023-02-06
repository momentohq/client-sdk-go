{{ ossHeader }}

## Getting Started :running:

### Requirements

- [Go version 1.18.\*](https://go.dev/dl/)
- A Momento Auth Token is required, you can generate one using
  the [Momento CLI](https://github.com/momentohq/momento-cli)

### Examples

Check out full working code in the [examples](./examples/README.md) directory of this repository!

### Installation

```bash
go get github.com/momentohq/client-sdk-go
```

### Usage

Checkout our [examples](./examples/README.md) directory for complete examples of how to use the SDK.

Here is a quickstart you can use in your own project:

```go
{{ usageExampleCode }}
```

### Error Handling

Errors that occur in calls to `ScsClient` methods are surfaced to developers as part of the return values of the calls, as opposed to by throwing exceptions. 
This makes them more visible, and allows your IDE to be more helpful in ensuring that you've handled the ones you care about. (For more on our philosophy about this, see our blog post on why [Exceptions are bugs](https://www.gomomento.com/blog/exceptions-are-bugs). And send us any feedback you have!)

The preferred way of interpreting the return values from `ScsClient` methods is using a `switch` statement to match and handle the specific response type. 
Here's a quick example:

```go
switch r := resp.(type) {
case *momento.CacheGetHit:
    log.Printf("Lookup resulted in cahce HIT. value=%s\n", r.ValueString())
default: 
    // you can handle other cases via pattern matching in other `switch case`, or a default case
    // via the `default` block.  For each return value your IDE should be able to give you code 
    // completion indicating the other possible "case"; in this case, `*momento.CacheGetMiss`.
}
```

Using this approach, you get a type-safe `CacheGetHit` object in the case of a cache hit. 
But if the cache read results in a Miss, you'll also get a type-safe object that you can use to get more info about what happened.

In cases where you get an error response, it can be treated as `momentoErr` using `As` method and it always include an `momentoErr.Code` that you can use to check the error type:

```go
err = client.CreateCache(ctx, &momento.CreateCacheRequest{
    CacheName: cacheName,
})
if err != nil {
    var momentoErr momento.MomentoError
    if errors.As(err, &momentoErr) {
        if momentoErr.Code() != momento.AlreadyExistsError {
            // this would represent that a cache you tried to create already exists.
        }
    }
}
```

### Tuning

Coming soon...

{{ ossFooter }}
