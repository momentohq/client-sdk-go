package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type GetRequest struct {
	// Name of the cache to get the item from
	CacheName string
	// string or byte key to be used to store item
	Key Key

	grpcRequest  *pb.XGetRequest
	grpcResponse *pb.XGetResponse
	response     responses.GetResponse
}

func (r *GetRequest) cacheName() string { return r.CacheName }

func (r *GetRequest) key() Key { return r.Key }

func (r *GetRequest) requestName() string { return "Get" }

func (r *GetRequest) initGrpcRequest(scsDataClient) error {
	var err error

	var key []byte
	if key, err = prepareKey(r); err != nil {
		return err
	}

	r.grpcRequest = &pb.XGetRequest{
		CacheKey: key,
	}

	return nil
}

func (r *GetRequest) makeGrpcRequest(requestMetadata context.Context, client scsDataClient) (grpcResponse, []metadata.MD, error) {
	var header, trailer metadata.MD
	fmt.Printf("in makeGrpcRequest for Get\n")
	trailer = metadata.Pairs()
	trailer.Set("test", "hello")
	fmt.Printf("trailer: %v (%T)\n", trailer, trailer)

	// TODO: This will never return a resp when an err is defined.
	//  This is why I can't use both resp and err as responses. Returning an error from the interceptor
	//  will always cause the resp to be nil.
	//  . . .
	resp, err := client.grpcClient.Get(requestMetadata, r.grpcRequest, grpc.Header(&header), grpc.Trailer(&trailer))

	fmt.Printf(">>>>>>>>>> resp: %T, err: %T\n", resp, err)
	if resp != nil {
		fmt.Printf(">>>>>>>>>> raw resp: %T\n", resp)
		fmt.Printf("resp.GetResult(): %v\n", resp.GetResult())
		fmt.Printf(">>>>>>>>>> grpc resp: %T, result: %s, err: %T\n", r.grpcResponse, r.grpcResponse.Result, err)
	} else {
		fmt.Println(">>>>>>>>>> resp is nil")
	}

	responseMetadata := []metadata.MD{header, trailer}
	fmt.Printf("???? %v\n", responseMetadata)
	if err != nil {
		return nil, responseMetadata, err
	}
	return resp, nil, nil
}

func (r *GetRequest) interpretGrpcResponse() error {
	resp := r.grpcResponse

	if resp.Result == pb.ECacheResult_Hit {
		r.response = responses.NewGetHit(resp.CacheBody)
		return nil
	} else if resp.Result == pb.ECacheResult_Miss {
		r.response = &responses.GetMiss{}
		return nil
	} else {
		return errUnexpectedGrpcResponse(r, r.grpcResponse)
	}
}

func (r *GetRequest) getResponse() interface{} {
	return r.response
}
