package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal"
	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"github.com/momentohq/client-sdk-go/storageTypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type storageDataClient struct {
	grpcManager    *grpcmanagers.StoreGrpcManager
	grpcClient     pb.StoreClient
	requestTimeout time.Duration
	endpoint       string
}

func newStorageDataClient(request *models.StorageDataClientRequest) (*storageDataClient, momentoerrors.MomentoSvcErr) {
	dataManager, err := grpcmanagers.NewStoreGrpcManager(&models.StoreGrpcManagerRequest{
		CredentialProvider: request.CredentialProvider,
		GrpcConfiguration:  request.Configuration.GetTransportStrategy().GetGrpcConfig(),
	})
	if err != nil {
		return nil, err
	}
	var timeout time.Duration
	if request.Configuration.GetClientSideTimeout() < 1 {
		timeout = defaultRequestTimeout
	} else {
		timeout = request.Configuration.GetClientSideTimeout()
	}

	return &storageDataClient{
		grpcManager:    dataManager,
		grpcClient:     pb.NewStoreClient(dataManager.Conn),
		requestTimeout: timeout,
		endpoint:       request.CredentialProvider.GetStorageEndpoint(),
	}, nil
}

func (client *storageDataClient) Close() {
	client.grpcManager.Close()
}

func (client *storageDataClient) delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateStoreMetadata(ctx, request.StoreName)

	var header, trailer metadata.MD // variable to store header and trailer
	_, err := client.grpcClient.Delete(
		requestMetadata,
		&pb.XStoreDeleteRequest{
			Key: request.Key,
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}
	return &responses.StorageDeleteSuccess{}, nil
}

func (client *storageDataClient) put(ctx context.Context, request *StoragePutRequest) (responses.StoragePutResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateStoreMetadata(ctx, request.StoreName)

	val := pb.XStoreValue{}
	switch request.Value.(type) {
	case storageTypes.Bytes:
		val.Value = &pb.XStoreValue_BytesValue{BytesValue: request.Value.(storageTypes.Bytes)}
	case storageTypes.String:
		val.Value = &pb.XStoreValue_StringValue{StringValue: string(request.Value.(storageTypes.String))}
	case storageTypes.Float:
		val.Value = &pb.XStoreValue_DoubleValue{DoubleValue: float64(request.Value.(storageTypes.Float))}
	case storageTypes.Int:
		val.Value = &pb.XStoreValue_IntegerValue{IntegerValue: int64(request.Value.(storageTypes.Int))}
	}

	var header, trailer metadata.MD // variable to store header and trailer
	_, err := client.grpcClient.Put(
		requestMetadata,
		&pb.XStorePutRequest{
			Key:   request.Key,
			Value: &val,
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
	}

	return &responses.StoragePutSuccess{}, nil
}

func (client *storageDataClient) get(ctx context.Context, request *StorageGetRequest) (responses.StorageGetResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := internal.CreateStoreMetadata(ctx, request.StoreName)

	var header, trailer metadata.MD // variable to store header and trailer
	response, err := client.grpcClient.Get(
		requestMetadata,
		&pb.XStoreGetRequest{
			Key: request.Key,
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	if err != nil {
		// Handle item not found error by returning a miss
		myErr := momentoerrors.ConvertSvcErr(err, header, trailer)
		if myErr.Code() == momentoerrors.ItemNotFoundError {
			return *responses.NewStoreGetResponse_Nil(), nil
		}
		return *responses.NewStoreGetResponse_Nil(), myErr
	}

	val := response.GetValue()
	switch val.Value.(type) {
	case *pb.XStoreValue_BytesValue:
		return *responses.NewStoreGetResponse_Bytes(val.GetBytesValue()), nil
	case *pb.XStoreValue_StringValue:
		return *responses.NewStoreGetResponse_String(val.GetStringValue()), nil
	case *pb.XStoreValue_DoubleValue:
		return *responses.NewStoreGetResponse_Float(val.GetDoubleValue()), nil
	case *pb.XStoreValue_IntegerValue:
		return *responses.NewStoreGetResponse_Integer(int(val.GetIntegerValue())), nil
	default:
		return *responses.NewStoreGetResponse_Nil(), momentoerrors.NewMomentoSvcErr(momentoerrors.UnknownServiceError, "Unknown store value type", nil)
	}
}
