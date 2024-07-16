package main

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	client, err := momento.NewPreviewStorageClient(config.StorageLaptopLatest(), credentialProvider)

	var storeName = "store-name"

	defer func() {
		fmt.Println("Deleting store")
		_, err = client.DeleteStore(ctx, &momento.DeleteStoreRequest{
			StoreName: storeName,
		})
		if err != nil {
			panic(err)
		}
	}()

	// ...
}
