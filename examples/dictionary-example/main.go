package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
	"github.com/momentohq/client-sdk-go/momento"
	"github.com/momentohq/client-sdk-go/utils"
)

const (
	cacheName             = "my-test-cache"
	dictionaryName        = "my-test-dictionary"
	itemDefaultTTLSeconds = 60
)

var (
	ctx    context.Context
	client momento.CacheClient
)

func setup() {
	ctx = context.Background()
	var credentialProvider, err = auth.NewEnvMomentoTokenProvider("MOMENTO_AUTH_TOKEN")
	if err != nil {
		panic(err)
	}

	// Initializes Momento
	client, err = momento.NewSimpleCacheClient(&momento.SimpleCacheClientProps{
		Configuration:      config.LatestLaptopConfig(),
		CredentialProvider: credentialProvider,
		DefaultTTL:         itemDefaultTTLSeconds * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// Create Cache
	_, err = client.CreateCache(ctx, &momento.CreateCacheRequest{
		CacheName: cacheName,
	})
	if err != nil {
		panic(err)
	}
}

func printField(field momento.Value) {
	fmt.Println("\nprinting field:")
	resp, err := client.DictionaryGetField(ctx, &momento.DictionaryGetFieldRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Field:          field,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.DictionaryGetFieldHit:
		fmt.Printf("field %s = '%s'\n", r.FieldString(), r.ValueString())
	case *momento.DictionaryGetFieldMiss:
		fmt.Println("get field returned MISS")
	}
}

func printFields(fields []momento.Value) {
	fmt.Println("\nprinting fields:")
	resp, err := client.DictionaryGetFields(ctx, &momento.DictionaryGetFieldsRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Fields:         fields,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.DictionaryGetFieldsHit:
		for k, v := range r.ValueMap() {
			fmt.Printf("%s = %s\n", k, v)
		}
	case *momento.DictionaryGetFieldsMiss:
		fmt.Println("dictionary get fields returned MISS")
	}
}

func setField(field momento.Value, value momento.Value) {
	resp, err := client.DictionarySetField(ctx, &momento.DictionarySetFieldRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Field:          field,
		Value:          value,
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *momento.DictionarySetFieldSuccess:
		fmt.Printf("\ndictionary field '%s' set to '%s'\n", field, value)
	}
}

func setItems(items map[string]momento.Value) {
	resp, err := client.DictionarySetFields(ctx, &momento.DictionarySetFieldsRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Items:          items,
	})
	if err != nil {
		panic(err)
	}
	switch resp.(type) {
	case *momento.DictionarySetFieldsSuccess:
		fmt.Println("\ndictionary fields set")
	}
}

func printDict() {
	resp, err := client.DictionaryFetch(ctx, &momento.DictionaryFetchRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
	})
	if err != nil {
		panic(err)
	}
	switch r := resp.(type) {
	case *momento.DictionaryFetchHit:
		fmt.Println("\nprinting dictionary contents:")
		for k, v := range r.ValueMapStringString() {
			fmt.Printf("%s = %s\n", k, v)
		}
	case *momento.DictionaryFetchMiss:
		fmt.Println("\ndictionary fetch returned MISS")
	}
}

func incrementField(counterField momento.Value, amount int64) {
	fmt.Println("\nincrementing field")
	resp, err := client.DictionaryIncrement(ctx, &momento.DictionaryIncrementRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Field:          counterField,
		Amount:         amount,
		CollectionTTL: utils.CollectionTTL{
			Ttl:        time.Second * 30,
			RefreshTtl: true,
		},
	})
	if err != nil {
		var momentoErr momento.MomentoError
		if errors.As(err, &momentoErr) {
			if momentoErr.Code() != momento.InvalidArgumentError {
				panic(err)
			} else {
				fmt.Printf("increment field with %d amount got expected invalid argument error", amount)
			}
		}
	}

	switch r := resp.(type) {
	case *momento.DictionaryIncrementSuccess:
		fmt.Printf("\nincremented counter field to: %d\n", r.Value())
	}
}

func removeField(field momento.Value) {
	_, err := client.DictionaryRemoveField(ctx, &momento.DictionaryRemoveFieldRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Field:          field,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nfield removed")
}

func removeFields(fields []momento.Value) {
	_, err := client.DictionaryRemoveFields(ctx, &momento.DictionaryRemoveFieldsRequest{
		CacheName:      cacheName,
		DictionaryName: dictionaryName,
		Fields:         fields,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\nfields removed")
}

func deleteDict() {
	_, err := client.Delete(ctx, &momento.DeleteRequest{
		CacheName: cacheName,
		Key:       momento.String(dictionaryName),
	})
	if err != nil {
		panic(err)
	}
}

func deleteCache() {
	_, err := client.DeleteCache(ctx, &momento.DeleteCacheRequest{CacheName: cacheName})
	if err != nil {
		panic(err)
	}
}

func main() {
	setup()

	field := "my-field"
	value := "my-value"

	setField(momento.String(field), momento.String(value))
	printDict()

	items := make(map[string]momento.Value)
	for i := 1; i < 11; i++ {
		numField := fmt.Sprintf("%s %d", field, i)
		numValue := fmt.Sprintf("%s %d", value, i)
		items[numField] = momento.String(numValue)
	}
	setItems(items)
	printDict()

	printField(momento.String("my-field 6"))

	var fields []momento.Value
	for i := 5; i < 12; i++ {
		fields = append(fields, momento.String(fmt.Sprintf("my-field %d", i)))
	}
	printFields(fields)

	counterField := momento.String("counter-field")
	setField(counterField, momento.String("0"))
	printField(counterField)
	incrementField(counterField, 25)
	printField(counterField)
	incrementField(counterField, 0)
	printField(counterField)

	removeField(counterField)
	printDict()

	removeFields(fields)
	printDict()

	fmt.Println("\ndeleting dictionary")
	deleteDict()
	fmt.Println("\ndeleting cache")
	deleteCache()
	fmt.Println("\ndone")
}
