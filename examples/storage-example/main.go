package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/storageTypes"
)

var (
	storeName = "store-name"
	key       = uuid.NewString()
)

func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	client, err := momento.NewPreviewStorageClient(config.StorageLaptopLatest(), credentialProvider)

	defer func() {
		fmt.Println("Deleting store")
		_, err = client.DeleteStore(ctx, &momento.DeleteStoreRequest{
			StoreName: storeName,
		})
		if err != nil {
			panic(err)
		}
	}()

	fmt.Println("Creating store")
	_, err = client.CreateStore(ctx, &momento.CreateStoreRequest{
		StoreName: storeName,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Putting value as string")
	_, err = client.Put(ctx, &momento.StoragePutRequest{
		StoreName: storeName,
		Key:       key,
		Value:     storageTypes.String("my-value"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Getting value")
	// getResp is a StorageGetResponse that is coerced to the "found" type below
	getResp, err := client.Get(ctx, &momento.StorageGetRequest{
		StoreName: storeName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}

	// If the value was not found, the response will be nil.
	if getResp == nil {
		fmt.Println("Got nil")
	}

	// Then get the value from the found response.
	// If you don't know the type beforehand:
	switch t := getResp.Value().(type) {
	case storageTypes.String:
		fmt.Printf("Got the string %s\n", t)
	case storageTypes.Bytes:
		fmt.Printf("Got the bytes %b\n", t)
	case storageTypes.Float:
		fmt.Printf("Got the float %f\n", t)
	case storageTypes.Int:
		fmt.Printf("Got the integer %d\n", t)
	}

	// If you know the type you're expecting, you can assert it directly:
	intVal, ok := getResp.Value().(storageTypes.Int)
	if !ok {
		fmt.Println("Illegal type assertion")
	} else {
		fmt.Printf("Got the integer %d\n", intVal)
	}

	// delete the key
	fmt.Println("Deleting key")
	_, err = client.Delete(ctx, &momento.StorageDeleteRequest{
		StoreName: storeName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
