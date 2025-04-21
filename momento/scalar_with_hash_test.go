package momento_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/momentohq/client-sdk-go/momento"
	. "github.com/momentohq/client-sdk-go/momento/test_helpers"
	. "github.com/momentohq/client-sdk-go/responses"
)

var _ = Describe("cache-client scalar-with-hash-methods", Label(CACHE_SERVICE_LABEL), func() {
	DescribeTable("Set and get with hash",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			var firstHash, secondHash string

			// Initial get should return a miss
			getResp, err := client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			Expect(
				getResp,
			).To(BeAssignableToTypeOf(&GetWithHashMiss{}))

			// Set the value and get its hash
			setResp, err := client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				Expect(result.HashString()).ToNot(BeNil())
				firstHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Get the value and make sure it matches the expected value
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
				Expect(result.HashString()).To(Equal(firstHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// Setting again with new value should overwrite it (set unconditionally)
			newValue := NewRandomMomentoString()
			setResp, err = client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     newValue,
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				Expect(result.HashString()).ToNot(BeNil())
				Expect(result.HashString()).ToNot(Equal(firstHash))
				secondHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			newExpectedBytes := []byte(newValue)
			newExpectedString := string(newValue)
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(newExpectedBytes))
				Expect(result.ValueString()).To(Equal(newExpectedString))
				Expect(result.HashString()).To(Equal(secondHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
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

	DescribeTable("Set if present and hash not equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			var firstHash, secondHash, notEqualHash string

			// Set some other key to get a hash value we can use for the not equal test
			setResp, err := client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       NewRandomMomentoString(),
				Value:     String("some-other-value"),
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				notEqualHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Initial set if should return not stored
			setIfResp, err := client.SetIfPresentAndHashNotEqual(sharedContext.Ctx, &SetIfPresentAndHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        value,
				HashNotEqual: String("some-hash-value"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndHashNotEqualNotStored{}))

			// Add the value for the key
			setResp, err = client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				Expect(result.HashString()).ToNot(BeNil())
				firstHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			getResp, err := client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
				Expect(result.HashString()).To(Equal(firstHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// Make sure we don't overwrite the value when HashNotEqual matches current hash value
			newValue := NewRandomMomentoString()
			setIfResp, err = client.SetIfPresentAndHashNotEqual(sharedContext.Ctx, &SetIfPresentAndHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        newValue,
				HashNotEqual: String(firstHash),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndHashNotEqualNotStored{}))

			// Make sure we overwrite the value when HashNotEqual does not match current hash value
			setIfResp, err = client.SetIfPresentAndHashNotEqual(sharedContext.Ctx, &SetIfPresentAndHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        newValue,
				HashNotEqual: String(notEqualHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfPresentAndHashNotEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				Expect(result.HashString()).ToNot(Equal(firstHash))
				secondHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfPresentAndHashNotEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			newExpectedBytes := []byte(newValue)
			newExpectedString := string(newValue)
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(newExpectedBytes))
				Expect(result.ValueString()).To(Equal(newExpectedString))
				Expect(result.HashString()).To(Equal(secondHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// make sure we error when HashNotEqual is nil
			_, err = client.SetIfPresentAndHashNotEqual(sharedContext.Ctx, &SetIfPresentAndHashNotEqualRequest{
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

	DescribeTable("Set if present and hash equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			var firstHash, secondHash, equalHash string

			// Set some other key to get a hash value we can use for the equal test
			setResp, err := client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       NewRandomMomentoString(),
				Value:     String("some-other-value"),
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				equalHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Initial set if should return not stored
			setIfResp, err := client.SetIfPresentAndHashEqual(sharedContext.Ctx, &SetIfPresentAndHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				HashEqual: String("some-hash-value"),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndHashEqualNotStored{}))

			// Add the value for the key
			setResp, err = client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				Expect(result.HashString()).ToNot(BeNil())
				firstHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			getResp, err := client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
				Expect(result.HashString()).To(Equal(firstHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// Make sure we don't overwrite the value when HashEqual does not match current hash value
			newValue := NewRandomMomentoString()
			setIfResp, err = client.SetIfPresentAndHashEqual(sharedContext.Ctx, &SetIfPresentAndHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     newValue,
				HashEqual: String(equalHash),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfPresentAndHashEqualNotStored{}))

			// Make sure we overwrite the value when HashEqual matches current hash value
			setIfResp, err = client.SetIfPresentAndHashEqual(sharedContext.Ctx, &SetIfPresentAndHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     newValue,
				HashEqual: String(firstHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfPresentAndHashEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				Expect(result.HashString()).ToNot(Equal(firstHash))
				secondHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfPresentAndHashNotEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			newExpectedBytes := []byte(newValue)
			newExpectedString := string(newValue)
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(newExpectedBytes))
				Expect(result.ValueString()).To(Equal(newExpectedString))
				Expect(result.HashString()).To(Equal(secondHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// make sure we error when HashEqual is nil
			_, err = client.SetIfPresentAndHashEqual(sharedContext.Ctx, &SetIfPresentAndHashEqualRequest{
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
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")))

	DescribeTable("Set if absent or hash equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			var firstHash, secondHash, equalHash string

			// Set some other key to get a hash value we can use for the equal test
			setResp, err := client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       NewRandomMomentoString(),
				Value:     String("some-other-value"),
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				equalHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Initial set if should store the value because key was absent
			setIfResp, err := client.SetIfAbsentOrHashEqual(sharedContext.Ctx, &SetIfAbsentOrHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     value,
				HashEqual: String(equalHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfAbsentOrHashEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				firstHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfAbsentOrHashEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been set
			getResp, err := client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
				Expect(result.HashString()).To(Equal(firstHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// Make sure we don't overwrite the value when HashEqual does not match current hash value
			newValue := NewRandomMomentoString()
			setIfResp, err = client.SetIfAbsentOrHashEqual(sharedContext.Ctx, &SetIfAbsentOrHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     newValue,
				HashEqual: String(equalHash),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentOrHashEqualNotStored{}))

			// Make sure we overwrite the value when HashEqual matches current hash value
			setIfResp, err = client.SetIfAbsentOrHashEqual(sharedContext.Ctx, &SetIfAbsentOrHashEqualRequest{
				CacheName: cacheName,
				Key:       key,
				Value:     newValue,
				HashEqual: String(firstHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfAbsentOrHashEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				Expect(result.HashString()).ToNot(Equal(firstHash))
				secondHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfPresentAndHashNotEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			newExpectedBytes := []byte(newValue)
			newExpectedString := string(newValue)
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(newExpectedBytes))
				Expect(result.ValueString()).To(Equal(newExpectedString))
				Expect(result.HashString()).To(Equal(secondHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// make sure we error when HashEqual is nil
			_, err = client.SetIfAbsentOrHashEqual(sharedContext.Ctx, &SetIfAbsentOrHashEqualRequest{
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
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")))

	DescribeTable("Set if absent or hash not equal",
		func(clientType string, key Key, value Value, expectedString string, expectedBytes []byte) {
			client, cacheName := sharedContext.GetClientPrereqsForType(clientType)
			var firstHash, secondHash, notEqualHash string

			// Set some other key to get a hash value we can use for the not equal test
			setResp, err := client.SetWithHash(sharedContext.Ctx, &SetWithHashRequest{
				CacheName: cacheName,
				Key:       NewRandomMomentoString(),
				Value:     String("some-other-value"),
			})
			Expect(err).To(BeNil())
			switch result := setResp.(type) {
			case *SetWithHashStored:
				notEqualHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetWithHashStored, got: %T, %v", result, result))
			}

			// Initial set if should store the value because key was absent
			setIfResp, err := client.SetIfAbsentOrHashNotEqual(sharedContext.Ctx, &SetIfAbsentOrHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        value,
				HashNotEqual: String(notEqualHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfAbsentOrHashNotEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				firstHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfAbsentOrHashNotEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been set
			getResp, err := client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(expectedBytes))
				Expect(result.ValueString()).To(Equal(expectedString))
				Expect(result.HashString()).To(Equal(firstHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// Make sure we don't overwrite the value when HashNotEqual matches current hash value
			newValue := NewRandomMomentoString()
			setIfResp, err = client.SetIfAbsentOrHashNotEqual(sharedContext.Ctx, &SetIfAbsentOrHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        newValue,
				HashNotEqual: String(firstHash),
			})
			Expect(err).To(BeNil())
			Expect(
				setIfResp,
			).To(BeAssignableToTypeOf(&SetIfAbsentOrHashNotEqualNotStored{}))

			// Make sure we overwrite the value when HashNotEqual does not match current hash value
			setIfResp, err = client.SetIfAbsentOrHashNotEqual(sharedContext.Ctx, &SetIfAbsentOrHashNotEqualRequest{
				CacheName:    cacheName,
				Key:          key,
				Value:        newValue,
				HashNotEqual: String(notEqualHash),
			})
			Expect(err).To(BeNil())
			switch result := setIfResp.(type) {
			case *SetIfAbsentOrHashNotEqualStored:
				Expect(result.HashString()).ToNot(BeNil())
				Expect(result.HashString()).ToNot(Equal(firstHash))
				secondHash = result.HashString()
			default:
				Fail(fmt.Sprintf("Expected SetIfAbsentOrHashNotEqualStored, got: %T, %v", result, result))
			}

			// Make sure the value has been overwritten
			newExpectedBytes := []byte(newValue)
			newExpectedString := string(newValue)
			getResp, err = client.GetWithHash(sharedContext.Ctx, &GetWithHashRequest{
				CacheName: cacheName,
				Key:       key,
			})
			Expect(err).To(BeNil())
			switch result := getResp.(type) {
			case *GetWithHashHit:
				Expect(result.ValueByte()).To(Equal(newExpectedBytes))
				Expect(result.ValueString()).To(Equal(newExpectedString))
				Expect(result.HashString()).To(Equal(secondHash))
			default:
				Fail(fmt.Sprintf("Expected GetWithHashHit, got: %T, %v", result, result))
			}

			// make sure we error when HashNotEqual is nil
			_, err = client.SetIfAbsentOrHashNotEqual(sharedContext.Ctx, &SetIfAbsentOrHashNotEqualRequest{
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
		Entry("with default cache name when the value is blank", WithDefaultCache, NewRandomMomentoString(), String("  "), "  ", []byte("  ")))
})
