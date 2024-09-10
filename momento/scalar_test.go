package momento_test

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("cache-client scalar-methods", Label(CACHE_SERVICE_LABEL), func() {
	DescribeTable("Gets, Sets, and Deletes",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: cacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				fmt.Println(`Expected GetHit, got:`, result)
				Fail("Unexpected type from Get")
			}

			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("using a consistent read concern client", WithConsistentReadConcern, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("using a balanced read concern client", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
	)

	DescribeTable("Set if not exists",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfNotExistsResp, err := client.SetIfNotExists(sharedContext.Ctx, &SetIfNotExistsRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfNotExistsResp,
			).To(BeAssignableToTypeOf(&SetIfNotExistsStored{}))

			// make sure we get a hit
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// we make another call and make sure that the value is not stored again
			setIfNotExistsRespNotStored, err := client.SetIfNotExists(sharedContext.Ctx, &SetIfNotExistsRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfNotExistsRespNotStored,
			).To(BeAssignableToTypeOf(&SetIfNotExistsNotStored{}))

			// make sure we get a hit
			getResp, err = client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))

			}

			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&DeleteSuccess{}))

			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if present",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfPresentResp, err := client.SetIfPresent(sharedContext.Ctx, &SetIfPresentRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfPresentResp,
			).To(BeAssignableToTypeOf(&SetIfPresentNotStored{}))

			// add the value for the key
			setResponse, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("initial value will be replaced"),
			})
			Expect(err).To(BeNil())
			Expect(setResponse).To(BeAssignableToTypeOf(&SetSuccess{}))

			setIfPresentResp, err = client.SetIfPresent(sharedContext.Ctx, &SetIfPresentRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfPresentResp,
			).To(BeAssignableToTypeOf(&SetIfPresentStored{}))

			// make sure the value has been overwritten
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}
		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if absent",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfResp, err := client.SetIfAbsent(sharedContext.Ctx, &SetIfAbsentRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentStored{}))

			// make sure we get a hit
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// we make another call and make sure that the value is not stored again
			setIfResp, err = client.SetIfAbsent(sharedContext.Ctx, &SetIfAbsentRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentNotStored{}))

			// make sure we get a hit
			getResp, err = client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())

			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}
		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if present and not equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfResp, err := client.SetIfPresentAndNotEqual(sharedContext.Ctx, &SetIfPresentAndNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				NotEqual:  String("not the initial value"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndNotEqualNotStored{}))

			// add the value for the key
			setResponse, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("initial value"),
			})
			Expect(err).To(BeNil())
			Expect(setResponse).To(BeAssignableToTypeOf(&SetSuccess{}))

			// make sure we don't overwrite the value when current value is the same as NotEqual
			setIfResp, err = client.SetIfPresentAndNotEqual(sharedContext.Ctx, &SetIfPresentAndNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				NotEqual:  String("initial value"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndNotEqualNotStored{}))

			// make sure we overwrite the value when the current value is different from NotEqual
			setIfResp, err = client.SetIfPresentAndNotEqual(sharedContext.Ctx, &SetIfPresentAndNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				NotEqual:  String("not the initial value"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndNotEqualStored{}))

			// make sure the value has been overwritten
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we error when NotEqual is nil
			_, err = client.SetIfPresentAndNotEqual(sharedContext.Ctx, &SetIfPresentAndNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfResp, err := client.SetIfEqual(sharedContext.Ctx, &SetIfEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				Equal:     String("i am irrelevant"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfEqualNotStored{}))

			// add the value for the key
			setResponse, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("initial value will be replaced"),
			})
			Expect(err).To(BeNil())
			Expect(setResponse).To(BeAssignableToTypeOf(&SetSuccess{}))

			// make sure we don't overwrite the value when current value is different from Equal
			setIfResp, err = client.SetIfEqual(sharedContext.Ctx, &SetIfEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				Equal:     String("i won't match"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfEqualNotStored{}))

			// make sure we overwrite the value when the current value is the same as Equal
			setIfResp, err = client.SetIfEqual(sharedContext.Ctx, &SetIfEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				Equal:     String("initial value will be replaced"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfEqualStored{}))

			// make sure the value has been overwritten
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we error when Equal is nil
			_, err = client.SetIfEqual(sharedContext.Ctx, &SetIfEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if absent or equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			// make sure nonexistent key is set
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			setIfResp, err := client.SetIfAbsentOrEqual(sharedContext.Ctx, &SetIfAbsentOrEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				Equal:     String("i am irrelevant"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentOrEqualStored{}))

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we don't overwrite the value when current value is different from Equal
			setIfResp, err = client.SetIfAbsentOrEqual(sharedContext.Ctx, &SetIfAbsentOrEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				Equal:     String("i won't match"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentOrEqualNotStored{}))

			// make sure we overwrite the value when the current value is the same as Equal
			setIfResp, err = client.SetIfAbsentOrEqual(sharedContext.Ctx, &SetIfAbsentOrEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("overwritten value"),
				Equal:     value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentOrEqualStored{}))

			// make sure the value has been overwritten
			getResp, err = client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueString()).To(Equal("overwritten value"))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we error when Equal is nil
			_, err = client.SetIfAbsentOrEqual(sharedContext.Ctx, &SetIfAbsentOrEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("Set if not equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			// make sure nonexistent key is set
			setIfResp, err := client.SetIfNotEqual(sharedContext.Ctx, &SetIfNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				NotEqual:  String("i am irrelevant"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfNotEqualStored{}))

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we don't overwrite the value when current value is the same as from NotEqual
			setIfResp, err = client.SetIfNotEqual(sharedContext.Ctx, &SetIfNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("i shouldn't be"),
				NotEqual:  value,
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfNotEqualNotStored{}))

			// make sure we overwrite the value when the current value is different from NotEqual
			setIfResp, err = client.SetIfNotEqual(sharedContext.Ctx, &SetIfNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     String("bingo!"),
				NotEqual:  String("i don't match"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfNotEqualStored{}))

			// make sure the value has been overwritten
			getResp, err = client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetHit:
				Expect(result.ValueString()).To(Equal("bingo!"))
			default:
				Fail(fmt.Sprintf("Expected GetHit, got: %T, %v", result, result))
			}

			// make sure we error when Equal is nil
			_, err = client.SetIfNotEqual(sharedContext.Ctx, &SetIfNotEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

		},
		Entry("when the key and value are strings", DefaultClient, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("when the key and value are bytes", DefaultClient, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("when the value is empty", DefaultClient, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("when the value is blank", DefaultClient, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
		Entry("with default cache name when the key and value are strings", WithDefaultCache, NewRandomMomentoString(), String("value"), "value", []byte("value")),
		Entry("with default cache name when the key and value are bytes", WithDefaultCache, NewRandomMomentoBytes(), Bytes("string"), "string", []byte("string")),
		Entry("with default cache name when the value is empty", WithDefaultCache, NewRandomMomentoString(), String(""), "", []byte("")),
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")),
	)

	DescribeTable("errors when the cache is missing",
		func(clientType string) {
			client, _ := sharedContext.GetClientPrereqsForType(clientType)
			cacheName := uuid.NewString()
			key := NewRandomMomentoString()
			value := String("value")

			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(CacheNotFoundError))

			setResp, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(CacheNotFoundError))

			deleteResp, err := client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(CacheNotFoundError))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("errors when the key is nil",
		func(clientType string) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
			Expect(
				client.Delete(sharedContext.Ctx, &DeleteRequest{
					CacheName: cacheName,
					Key:       nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("returns a miss when the key doesn't exist",
		func(clientType string) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			Expect(
				client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: cacheName,
					Key:       NewRandomMomentoString(),
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		},
		Entry("with default client", DefaultClient),
		Entry("with client with default cache", WithDefaultCache),
	)

	DescribeTable("invalid cache names and keys",
		func(clientType string, cacheName string, key Key, value Key) {
			client, _ := sharedContext.GetClientPrereqsForType(clientType)
			getResp, err := client.Get(sharedContext.Ctx, &GetRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(getResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			setResp, err := client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(setResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			deleteResp, err := client.Delete(sharedContext.Ctx, &DeleteRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(deleteResp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))
		},
		Entry("With default client and an empty cache name", DefaultClient, "", NewRandomMomentoString(), String("value")),
		Entry("With default client and  an bank cache name", DefaultClient, "   ", NewRandomMomentoString(), String("value")),
		Entry("With default client and  an empty key", DefaultClient, uuid.NewString(), String(""), String("value")),
		Entry("With client with default cache and an bank cache name", WithDefaultCache, "   ", NewRandomMomentoString(), String("value")),
		Entry("With client with default cache and an empty key", WithDefaultCache, uuid.NewString(), String(""), String("value")),
	)

	Describe("Set", func() {
		It("Uses the default TTL", func() {
			key := NewRandomMomentoString()
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(sharedContext.DefaultTtl / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It("Overrides the default TTL", func() {
			key := NewRandomMomentoString()
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
					Ttl:       sharedContext.DefaultTtl * 2,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			time.Sleep(sharedContext.DefaultTtl / 2)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetHit{}))

			time.Sleep(sharedContext.DefaultTtl)

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It("Overrides the default TTL without unit and invalid", func() {
			key := NewRandomMomentoString()
			value := String("value")

			resp, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       key,
				Value:     value,
				Ttl:       2,
			})

			Expect(resp).To(BeNil())
			Expect(err).To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.Get(sharedContext.Ctx, &GetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				}),
			).To(BeAssignableToTypeOf(&GetMiss{}))
		})

		It("returns an error for a nil value", func() {
			key := NewRandomMomentoString()
			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     nil,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})
	})

	Describe("Keys exist", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				strKey := fmt.Sprintf("#%d", i)
				strVal := fmt.Sprintf("%sValue", strKey)
				Expect(
					sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: sharedContext.CacheName,
						Key:       String(strKey),
						Value:     String(strVal),
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))
				Expect(
					sharedContext.ClientWithDefaultCacheName.Set(sharedContext.Ctx, &SetRequest{
						Key:   String(strKey),
						Value: String(strVal),
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))
			}
		})

		DescribeTable("check for valid keys exist results",
			func(clientType string, toCheck []Key, expected []bool) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				resp, err := client.KeysExist(sharedContext.Ctx, &KeysExistRequest{
					CacheName: cacheName,
					Keys:      toCheck,
				})
				Expect(err).To(BeNil())
				switch result := resp.(type) {
				case *KeysExistSuccess:
					Expect(result.Exists()).To(Equal(expected))
				default:
					Fail(fmt.Sprintf("expected keys exist success but got %s", result))
				}
			},
			Entry("all hits", DefaultClient, []Key{String("#1"), String("#2")}, []bool{true, true}),
			Entry("all misses", DefaultClient, []Key{String("nope"), String("stillnope")}, []bool{false, false}),
			Entry("mixed", DefaultClient, []Key{String("nope"), String("#1")}, []bool{false, true}),
			Entry("all hits with default cache", WithDefaultCache, []Key{String("#1"), String("#2")}, []bool{true, true}),
			Entry("all misses with default cache", WithDefaultCache, []Key{String("nope"), String("stillnope")}, []bool{false, false}),
			Entry("mixed with default cache", WithDefaultCache, []Key{String("nope"), String("#1")}, []bool{false, true}),
		)
	})

	Describe("UpdateTtl", func() {
		DescribeTable("Overwrites Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				key := NewRandomMomentoString()
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       6 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&UpdateTtlSet{}))

				time.Sleep(sharedContext.DefaultTtl)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetHit{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		DescribeTable("Increases Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
				key := NewRandomMomentoString()
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
						Ttl:       2 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       3 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&IncreaseTtlSet{}))

				time.Sleep(2 * time.Second)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetHit{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		DescribeTable("Decreases Ttl",
			func(clientType string) {
				client, cacheName := sharedContext.GetClientPrereqsForType(clientType)

				key := NewRandomMomentoString()
				value := String("value")

				Expect(
					client.Set(sharedContext.Ctx, &SetRequest{
						CacheName: cacheName,
						Key:       key,
						Value:     value,
					}),
				).To(BeAssignableToTypeOf(&SetSuccess{}))

				Expect(
					client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
						CacheName: cacheName,
						Key:       key,
						Ttl:       2 * time.Second,
					}),
				).To(BeAssignableToTypeOf(&DecreaseTtlSet{}))

				time.Sleep(2 * time.Second)

				Expect(
					client.Get(sharedContext.Ctx, &GetRequest{
						CacheName: cacheName,
						Key:       key,
					}),
				).To(BeAssignableToTypeOf(&GetMiss{}))
			},
			Entry("with default client", DefaultClient),
			Entry("with client with default cache", WithDefaultCache),
		)

		It("Returns InvalidArgumentError with negative or zero Ttl value for UpdateTtl", func() {
			key := NewRandomMomentoString()
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns InvalidArgumentError with negative or zero Ttl value for IncreaseTtl", func() {
			key := NewRandomMomentoString()
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns InvalidArgumentError with negative or zero Ttl value for DecreaseTtl", func() {
			key := NewRandomMomentoString()
			value := String("value")

			Expect(
				sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Value:     value,
				}),
			).To(BeAssignableToTypeOf(&SetSuccess{}))

			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       0,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))

			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
					Ttl:       -1,
				}),
			).Error().To(HaveMomentoErrorCode(InvalidArgumentError))
		})

		It("Returns UpdateTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.UpdateTtl(sharedContext.Ctx, &UpdateTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-update-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&UpdateTtlMiss{}))
		})

		It("Returns IncreaseTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.IncreaseTtl(sharedContext.Ctx, &IncreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-increase-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&IncreaseTtlMiss{}))
		})

		It("Returns DecreaseTtlMiss when a key doesn't exist", func() {
			Expect(
				sharedContext.Client.DecreaseTtl(sharedContext.Ctx, &DecreaseTtlRequest{
					CacheName: sharedContext.CacheName,
					Key:       String("test-decrease-ttl-miss"),
					Ttl:       sharedContext.DefaultTtl,
				}),
			).To(BeAssignableToTypeOf(&DecreaseTtlMiss{}))
		})
	})

	Describe("Increment", func() {
		It("Increments from 0 to expected amount with string field", func() {
			field := String("field")

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    1,
			})
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			Expect(err).To(BeNil())
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(1)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    41,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(42)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    -1042,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(-1000)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})

		It("Increments from 0 to expected amount with bytes field", func() {
			field := NewRandomMomentoBytes()

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    1,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(1)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    41,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(42)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    -1042,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(-1000)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})

		It("Increments with setting and resetting field", func() {
			field := String("field")
			value := String("10")
			Expect(sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{CacheName: sharedContext.CacheName, Key: field, Value: value})).To(BeAssignableToTypeOf(&SetSuccess{}))

			resp, err := sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    0,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(10)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    90,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(100)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}

			// Reset the field
			value = String("0")
			Expect(sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{CacheName: sharedContext.CacheName, Key: field, Value: value})).To(BeAssignableToTypeOf(&SetSuccess{}))
			resp, err = sharedContext.Client.Increment(sharedContext.Ctx, &IncrementRequest{
				CacheName: sharedContext.CacheName,
				Field:     field,
				Amount:    0,
			})
			Expect(err).To(BeNil())
			Expect(
				resp,
			).To(BeAssignableToTypeOf(&IncrementSuccess{}))
			switch result := resp.(type) {
			case *IncrementSuccess:
				Expect(result.Value()).To(Equal(int64(0)))
			default:
				Fail(fmt.Sprintf("expected increment success but got %s", result))
			}
		})
	})

	Describe("ItemGetType", func() {
		scalarName := NewRandomMomentoString()
		dictionaryName := NewRandomString()
		listName := NewRandomString()
		setName := NewRandomString()
		sortedSetName := NewRandomString()

		BeforeEach(func() {
			_, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       scalarName,
				Value:     String("hi"),
			})
			if err != nil {
				Fail("failed trying to set key")
			}
			_, err = sharedContext.Client.DictionarySetField(sharedContext.Ctx, &DictionarySetFieldRequest{
				CacheName:      sharedContext.CacheName,
				DictionaryName: dictionaryName,
				Field:          String("hi"),
				Value:          String("there"),
			})
			if err != nil {
				Fail("failed trying to add dictionary")
			}
			_, err = sharedContext.Client.ListPushFront(sharedContext.Ctx, &ListPushFrontRequest{
				CacheName: sharedContext.CacheName,
				ListName:  listName,
				Value:     String("hi"),
			})
			if err != nil {
				Fail("failed trying to add list")
			}
			_, err = sharedContext.Client.SetAddElement(sharedContext.Ctx, &SetAddElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   setName,
				Element:   String("hi"),
			})
			if err != nil {
				Fail("failed trying to add set")
			}
			_, err = sharedContext.Client.SortedSetPutElement(sharedContext.Ctx, &SortedSetPutElementRequest{
				CacheName: sharedContext.CacheName,
				SetName:   sortedSetName,
				Value:     String("hi"),
				Score:     1.0,
			})
			if err != nil {
				Fail("failed trying to add sorted set")
			}
		})
		DescribeTable("returns the correct item type",
			func(key Key, expectedType ItemType) {
				typeResponse, err := sharedContext.Client.ItemGetType(sharedContext.Ctx, &ItemGetTypeRequest{
					CacheName: sharedContext.CacheName,
					Key:       key,
				})
				Expect(err).To(BeNil())
				Expect(typeResponse).To(BeAssignableToTypeOf(&ItemGetTypeHit{}))
				switch result := typeResponse.(type) {
				case *ItemGetTypeHit:
					Expect(result.Type()).To(Equal(expectedType))
				default:
					Fail(fmt.Sprintf("expected ItemGetTypeHit but got %s", result))
				}
			},
			Entry("Scalar", scalarName, ItemTypeScalar),
			Entry("Dictionary", String(dictionaryName), ItemTypeDictionary),
			Entry("Set", String(setName), ItemTypeSet),
			Entry("List", String(listName), ItemTypeList),
			Entry("Sorted set", String(sortedSetName), ItemTypeSortedSet),
		)
	})

	Describe("item get ttl", func() {
		It("accurately reports the remaining TTL for a key", func() {
			key := NewRandomMomentoString()
			value := NewRandomMomentoString()
			var ttl = time.Duration(time.Minute * 2)
			_, err := sharedContext.Client.Set(sharedContext.Ctx, &SetRequest{
				CacheName: sharedContext.CacheName,
				Key:       key,
				Value:     value,
				Ttl:       ttl,
			})
			Expect(err).To(BeNil())

			resp, err := sharedContext.Client.ItemGetTtl(sharedContext.Ctx, &ItemGetTtlRequest{
				CacheName: sharedContext.CacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			Expect(resp).To(BeAssignableToTypeOf(&ItemGetTtlHit{}))
			switch result := resp.(type) {
			case *ItemGetTtlHit:
				fmt.Printf("Original TTL: %v, Remaining TTL: %v\n", ttl, result.RemainingTtl())
				Expect(result.RemainingTtl()).To(BeNumerically("<=", ttl))
				Expect(result.RemainingTtl() > (time.Second * 30)).To(BeTrue())
			default:
				Fail(fmt.Sprintf("expected ItemGetTtlHit but got %s", result))
			}
		})
	})
})
