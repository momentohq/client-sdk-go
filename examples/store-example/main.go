package main

import (
	"context"
	"fmt"
	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/responses"
)

// NOTE: this is a playground for exercising the preview store client during development.
// When the backend service is available, this code will be rewritten to illustrate the
// intended use of the storage client.
func main() {
	ctx := context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_API_KEY")
	if err != nil {
		panic(err)
	}

	client, err := momento.NewPreviewStoreClient(config.StoreDefault(), credentialProvider)
	putResp, err := client.Put(ctx, &momento.StorePutRequest{
		StoreName: "store-name",
		Key:       "my-key",
		Value:     momento.String("my-value"),
	})
	if err != nil {
		panic(err)
	}
	switch putResp.(type) {
	case *responses.StorePutSuccess:
		fmt.Printf("Explicitly got Success type for PUT\n")
	default:
		fmt.Printf("Unknown response type\n")
	}

	// tryGetResp is a StoreGetResponse that is coerced to the success type below
	tryGetResp, err := client.Get(ctx, &momento.StoreGetRequest{
		StoreName: "store-name",
		Key:       "str-key",
	})
	if err != nil {
		panic(err)
	}

	// This is possible because I've moved the TryGet* funcs to the StoreGetResponse interface. This
	// retains the ability to explicitly check for the success type but allows users to grab values
	// without needing to do a type assertion.
	if tryGetResp.ValueType() == responses.STRING {
		myStr, gotStr := tryGetResp.TryGetValueString()
		if gotStr {
			fmt.Printf("Got the string %s\n", myStr)
		} else {
			fmt.Printf("Did not get the string\n")
		}
	}

	// The success type is the only implementor of the StoreGetResponse interface.
	getResp := tryGetResp.(*responses.StoreGetSuccess)

	fmt.Printf("Trying to get double value from type %s\n", getResp.ValueType())
	val, b00l := getResp.TryGetValueDouble()
	if b00l {
		fmt.Printf("Got the double %f\n", val)
	} else {
		fmt.Printf("Did not get the double\n")
	}

	myStr, gotStr := getResp.TryGetValueString()
	if gotStr {
		fmt.Printf("Got the string %s\n", myStr)
	} else {
		fmt.Printf("Did not get the string\n")
	}

	// the below skips the type check and just tries to get the value
	bytesResp, err := client.Get(ctx, &momento.StoreGetRequest{
		StoreName: "store-name",
		Key:       "bytes-key",
	})
	if err != nil {
		panic(err)
	}
	bytesVal, gotBytes := bytesResp.TryGetValueBytes()
	if gotBytes {
		fmt.Printf("Got the bytes %s (%b)\n", bytesVal, bytesVal)
	} else {
		fmt.Printf("Did not get the bytes\n")
	}
}
