package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/utils"
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
		Value:     utils.StorageValueString("my-value"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Getting value")
	value, err := client.Get(ctx, &momento.StorageGetRequest{
		StoreName: storeName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}

	// A nil value indicates that the key was not found.
	if value == nil {
		fmt.Println("Did not find a value for the key; exiting")
		os.Exit(1)
	}

	// Then get the value from the found response.
	// If you don't know the type beforehand:
	switch t := value.(type) {
	case utils.StorageValueString:
		fmt.Printf("Got the string %s\n", t)
	case utils.StorageValueBytes:
		fmt.Printf("Got the bytes %b\n", t)
	case utils.StorageValueFloat:
		fmt.Printf("Got the float %f\n", t)
	case utils.StorageValueInt:
		fmt.Printf("Got the integer %d\n", t)
	}

	// If you know the type beforehand:
	fmt.Printf("Got the string %s\n", value.(utils.StorageValueString))

	// If you choose the wrong type:
	intVal, ok := value.(utils.StorageValueInt)
	if !ok {
		fmt.Println("Illegal type assertion")
	} else {
		fmt.Printf("Got the integer %d\n", intVal)
	}

	// You can do it in one shot, but it'll panic if you guess wrong like any cast would
	//fmt.Printf("Got the integer %d\n", value.(utils.StorageValueInteger))

	// delete the key
	fmt.Println("Deleting key")
	_, err = client.Delete(ctx, &momento.StorageDeleteRequest{
		StoreName: storeName,
		Key:       key,
	})
	if err != nil {
		panic(err)
	}

	// get a key that doesn't exist
	fmt.Println("Getting a key that doesn't exist")
	value, err = client.Get(ctx, &momento.StorageGetRequest{
		StoreName: storeName,
		Key:       uuid.NewString(),
	})
	if err != nil {
		panic(err)
	}
	// result: "Got the value <nil>"
	fmt.Printf("Got the value %v\n", value)
	// Trying to coerce a nil value will also panic of course
	// fmt.Printf("Got the string %s\n", value.(utils.StorageValueString))

	fmt.Println("Done")
}
