package momento

import (
	"context"
	"fmt"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

/////////// Response

type SortedSetRemoveResponse interface {
	isSortedSetRemoveResponse()
}

type SortedSetRemoveSuccess struct{}

func (SortedSetRemoveSuccess) isSortedSetRemoveResponse() {}

////////// Request

type SortedSetRemoveRequest struct {
	CacheName        string
	SetName          string
	ElementsToRemove SortedSetRemoveNumElements

	grpcRequest  *pb.XSortedSetRemoveRequest
	grpcResponse *pb.XSortedSetRemoveResponse
	response     SortedSetRemoveResponse
}

type SortedSetRemoveRequestElement struct {
	Name Bytes
}

type SortedSetRemoveNumElements interface {
	isSortedSetRemoveNumElement()
}

type RemoveAllElements struct{}

func (RemoveAllElements) isSortedSetRemoveNumElement() {}

type RemoveSomeElements struct {
	Elements []Bytes
}

func (RemoveSomeElements) isSortedSetRemoveNumElement() {}

func (r *SortedSetRemoveRequest) cacheName() string { return r.CacheName }

func (r *SortedSetRemoveRequest) requestName() string { return "Sorted set remove" }

func (r *SortedSetRemoveRequest) initGrpcRequest(scsDataClient) error {
	var err error

	if _, err = prepareName(r.SetName, "Set name"); err != nil {
		return err
	}

	grpcReq := &pb.XSortedSetRemoveRequest{
		SetName: []byte(r.SetName),
	}

	switch toRemove := r.ElementsToRemove.(type) {
	case RemoveAllElements:
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_All{}
	case *RemoveAllElements:
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_All{}
	case RemoveSomeElements:
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
			Some: &pb.XSortedSetRemoveRequest_XSome{
				ElementName: momentoBytesListToPrimitiveByteList(toRemove.Elements),
			},
		}
	case *RemoveSomeElements:
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
			Some: &pb.XSortedSetRemoveRequest_XSome{
				ElementName: momentoBytesListToPrimitiveByteList(toRemove.Elements),
			},
		}
	default:
		return fmt.Errorf("%T is an unrecognized type for Elements", r.ElementsToRemove)
	}

	r.grpcRequest = grpcReq

	return nil
}

func (r *SortedSetRemoveRequest) makeGrpcRequest(metadata context.Context, client scsDataClient) (grpcResponse, error) {
	resp, err := client.grpcClient.SortedSetRemove(metadata, r.grpcRequest)
	if err != nil {
		return nil, err
	}

	r.grpcResponse = resp

	return resp, nil
}

func (r *SortedSetRemoveRequest) interpretGrpcResponse() error {
	r.response = &SortedSetRemoveSuccess{}
	return nil
}
