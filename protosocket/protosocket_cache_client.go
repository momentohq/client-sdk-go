package protosocket

/*
#cgo pkg-config: --static momento_protosocket_ffi
#cgo !darwin LDFLAGS: -lgcc_s -lutil -lrt -lpthread -ldl -lm -lc
#cgo darwin LDFLAGS: -framework Security -framework CoreFoundation -lc++ -liconv -ldl -lm -lc
#include <momento_protosocket_ffi.h>
#include <string.h>
#include <math.h>

extern void setCallback(ProtosocketResult* result, void* user_data);
extern void getCallback(ProtosocketResult* result, void* user_data);
*/
import "C"
import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/momentohq/client-sdk-go/auth"
	"github.com/momentohq/client-sdk-go/config"
)

type SetResponse struct {
	Success bool
	Error   string
}

type GetResponse struct {
	Hit   bool
	Value []byte
	Error string
}

var (
	setContexts sync.Map // map[uint64]chan SetResponse
	getContexts sync.Map // map[uint64]chan GetResponse
	nextID      uint64   // atomic counter
)

func NewProtosocketCacheClient(config config.Configuration, credentialProvider auth.CredentialProvider, defaultTtl time.Duration) error {
	protosocketConfig := C.ProtosocketClientConfiguration{
		timeout_millis:   C.ulong(config.GetClientSideTimeout().Milliseconds()),
		connection_count: C.ulong(1),
	}
	protosocketCredentialProvider := C.ProtosocketCredentialProvider{
		env_var_name: C.CString("MOMENTO_API_KEY"),
	}
	defaultTtlMillis := C.ulonglong(defaultTtl.Milliseconds())
	C.init_protosocket_cache_client(defaultTtlMillis, protosocketConfig, protosocketCredentialProvider)
	return nil
}

func CloseProtosocketCacheClient() {
	C.destroy_protosocket_cache_client()
}

func convertGoStringToCBytes(string string) *C.Bytes {
	bytes := []byte(string)
	return convertGoBytesToCBytes(bytes)
}

func convertGoBytesToCBytes(bytes []byte) *C.Bytes {
	c_bytes := C.malloc(C.size_t(len(bytes)))
	C.memcpy(c_bytes, unsafe.Pointer(&bytes[0]), C.size_t(len(bytes)))
	return &C.Bytes{
		data:   (*C.uchar)(c_bytes),
		length: C.ulong(len(bytes)),
	}
}

func convertCBytesToGoBytes(c_bytes *C.Bytes) []byte {
	return C.GoBytes(unsafe.Pointer(c_bytes.data), C.int(c_bytes.length))
}

//export setCallback
func setCallback(result *C.ProtosocketResult, userData unsafe.Pointer) {
	// Decode the channel ID from the pointer
	id := uint64(uintptr(userData))

	// Load and delete the channel from the map. If there is no channel, the callback can't send a response
	chInterface, ok := setContexts.LoadAndDelete(id)
	if !ok {
		C.free_response(result)
		fmt.Printf("[Error] callback unable to find a channel to send a response\n")
		return
	}
	ch := chInterface.(chan SetResponse)

	// Convert the result to a set response
	responseType := C.GoString(result.response_type)
	var response SetResponse
	if responseType == "SetSuccess" {
		response.Success = true
	} else if responseType == "Error" {
		response.Error = C.GoString(result.error_message)
	}

	// Send the response to the original caller
	ch <- response
	C.free_response(result)
}

//export getCallback
func getCallback(result *C.ProtosocketResult, userData unsafe.Pointer) {
	// Decode the channel ID from the pointer
	id := uint64(uintptr(userData))

	// Load and delete the channel from the map. If there is no channel, the callback can't send a response
	chInterface, ok := getContexts.LoadAndDelete(id)
	if !ok {
		C.free_response(result)
		fmt.Printf("[Error] callback unable to find a channel to send a response\n")
		return
	}
	ch := chInterface.(chan GetResponse)

	// Convert the result to a get response
	responseType := C.GoString(result.response_type)
	var response GetResponse
	if responseType == "GetHit" {
		response.Hit = true
		response.Value = convertCBytesToGoBytes(result.value)
	} else if responseType == "GetMiss" {
		response.Hit = false
	} else if responseType == "Error" {
		response.Error = C.GoString(result.error_message)
	}

	// Send the response to the original caller
	ch <- response
	C.free_response(result)
}

func ProtosocketSet(cacheName string, key string, value string) {
	// Generate FFI-compatible versions of the variables and set them up to be freed
	cacheNameC := C.CString(cacheName)
	defer C.free(unsafe.Pointer(cacheNameC))
	keyC := convertGoStringToCBytes(key)
	defer C.free(unsafe.Pointer(keyC.data))
	valueC := convertGoStringToCBytes(value)
	defer C.free(unsafe.Pointer(valueC.data))

	// Create the channel the callback will send the response through
	responseCh := make(chan SetResponse, 1)

	// Generate a key for the channel and store it in the map for the callback to look up
	id := atomic.AddUint64(&nextID, 1)
	setContexts.Store(id, responseCh)

	C.protosocket_cache_client_set(
		cacheNameC,
		keyC,
		valueC,
		C.ProtosocketCallback(C.setCallback),
		unsafe.Pointer(uintptr(id)),
	)

	// Wait for the callback to send the response
	select {
	case response := <-responseCh:
		if response.Success {
			fmt.Printf("[INFO] set success\n")
		} else {
			fmt.Printf("[ERROR] set error: %v\n", response.Error)
		}
	case <-time.After(30 * time.Second):
		fmt.Printf("[ERROR] set timeout after 30 seconds\n")
		// Clean up the stored channel
		getContexts.Delete(id)
	}
}

func ProtosocketGet(cacheName string, key string) {
	// Generate FFI-compatible versions of the variables and set them up to be freed
	cacheNameC := C.CString(cacheName)
	defer C.free(unsafe.Pointer(cacheNameC))
	keyC := convertGoStringToCBytes(key)
	defer C.free(unsafe.Pointer(keyC.data))

	// Create the channel the callback will send the response through
	responseCh := make(chan GetResponse, 1)

	// Generate a key for the channel and store it in the map for the callback to look up
	id := atomic.AddUint64(&nextID, 1)
	getContexts.Store(id, responseCh)

	C.protosocket_cache_client_get(
		cacheNameC,
		keyC,
		C.ProtosocketCallback(C.getCallback),
		unsafe.Pointer(uintptr(id)),
	)

	// Wait for the callback to send the response
	select {
	case response := <-responseCh:
		if response.Hit {
			fmt.Printf("[INFO] get hit | raw value: %v | string value: %s\n", response.Value, string(response.Value))
		} else if response.Error != "" {
			fmt.Printf("[ERROR] get error: %v\n", response.Error)
		} else {
			fmt.Printf("[INFO] get miss\n")
		}
	case <-time.After(30 * time.Second):
		fmt.Printf("[ERROR] get timeout after 30 seconds\n")
		// Clean up the stored channel
		getContexts.Delete(id)
	}
}
