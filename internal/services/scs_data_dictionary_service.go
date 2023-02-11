package services

import (
	"context"
	"fmt"

	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	//"github.com/momentohq/client-sdk-go/utils"
	"google.golang.org/grpc/metadata"
)

func (client *ScsDataClient) DictionaryFetch(ctx context.Context, request *models.DictionaryFetchRequest) (models.DictionaryFetchResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	resp, err := client.grpcClient.DictionaryFetch(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		&pb.XDictionaryFetchRequest{
			DictionaryName: []byte(request.DictionaryName),
		},
	)
	if err != nil {
		panic(err)
	}

	switch r := resp.Dictionary.(type) {
	case *pb.XDictionaryFetchResponse_Found:
		return &models.DictionaryFetchHit{
			Items: dictionaryGrpcFieldValuePairToModel(r.Found.Items),
		}, nil
	case *pb.XDictionaryFetchResponse_Missing:
		return &models.DictionaryFetchMiss{}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.ClientSdkError,
			fmt.Sprintf("Unknown response type for dictionary fetch response %+v", r),
			nil,
		)
	}
	return nil, nil
}

func (client *ScsDataClient) DictionaryGetField(ctx context.Context, request *models.DictionaryGetFieldRequest) (models.DictionaryGetFieldResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()
	dictGetRequest := &pb.XDictionaryGetRequest{
		DictionaryName: []byte(request.DictionaryName),
	}
	dictGetRequest.Fields = append(dictGetRequest.Fields, request.Field)

	resp, err := client.grpcClient.DictionaryGet(
		metadata.NewOutgoingContext(ctx, createNewMetadata(request.CacheName)),
		dictGetRequest,
	)
	if err != nil {
		panic(err)
	}

	switch r := resp.Dictionary.(type) {
	case *pb.XDictionaryGetResponse_Found:
		return &models.DictionaryGetFieldHit{
			Value: r.Found.Items[0].CacheBody,
		}, nil
	case *pb.XDictionaryGetResponse_Missing:
		return &models.DictionaryGetFieldMiss{}, nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(
			momentoerrors.ClientSdkError,
			fmt.Sprintf("Unknown response type for dictionary fetch response %+v", r),
			nil,
		)
	}
}

func dictionaryGrpcFieldValuePairToModel(grpcFieldValuePairs []*pb.XDictionaryFieldValuePair) map[string][]byte {
	retMap := make(map[string][]byte)
	for i := range grpcFieldValuePairs {
		retMap[string(grpcFieldValuePairs[i].Field)] = grpcFieldValuePairs[i].Value
	}
	return retMap
}
