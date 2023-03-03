package momento

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/responses"

	pb "github.com/momentohq/client-sdk-go/internal/protos"
)

type SortedSetRemoveRequest struct {
	CacheName        string
	SetName          string
	ElementsToRemove SortedSetRemoveNumElements

	grpcRequest  *pb.XSortedSetRemoveRequest
	grpcResponse *pb.XSortedSetRemoveResponse
	response     responses.SortedSetRemoveResponse
}

type SortedSetRemoveRequestElement struct {
	Name Value
}

type SortedSetRemoveNumElements interface {
	isSortedSetRemoveNumElement()
}

type RemoveAllElements struct{}

func (RemoveAllElements) isSortedSetRemoveNumElement() {}

type RemoveSomeElements struct {
	Elements []Value
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
		elemToRemove, err := momentoValuesToPrimitiveByteList(toRemove.Elements)
		if err != nil {
			return err
		}
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
			Some: &pb.XSortedSetRemoveRequest_XSome{
				Values: elemToRemove,
			},
		}
	case *RemoveSomeElements:
		elemToRemove, err := momentoValuesToPrimitiveByteList(toRemove.Elements)
		if err != nil {
			return err
		}
		grpcReq.RemoveElements = &pb.XSortedSetRemoveRequest_Some{
			Some: &pb.XSortedSetRemoveRequest_XSome{
				Values: elemToRemove,
			},
		}
	default:
		return fmt.Errorf("%T is an unrecognized type for ElementsToRemove", r.ElementsToRemove)
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
	r.response = &responses.SortedSetRemoveSuccess{}
	return nil
}
