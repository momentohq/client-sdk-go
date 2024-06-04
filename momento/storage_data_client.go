package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
	"google.golang.org/grpc/metadata"
)

type storageDataClient struct {
	grpcManager    *grpcmanagers.StoreGrpcManager
	grpcClient     pb.StoreClient
	requestTimeout time.Duration
	endpoint       string
}

func newStorageDataClient(request *models.StoreDataClientRequest) (*storageDataClient, momentoerrors.MomentoSvcErr) {
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

func (*storageDataClient) CreateNewMetadata(storeName string) metadata.MD {
	return metadata.Pairs("store", storeName)
}

func (client *storageDataClient) delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

	_, err := client.grpcClient.Delete(requestMetadata, &pb.XStoreDeleteRequest{
		Key: request.Key,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}
	return &responses.StorageDeleteSuccess{}, nil
}

func (client *storageDataClient) set(ctx context.Context, request *StorageSetRequest) (responses.StorageSetResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

	val := pb.XStoreValue{}
	switch request.Value.(type) {
	case Bytes:
		val.Value = &pb.XStoreValue_BytesValue{BytesValue: request.Value.(Bytes)}
	case String:
		val.Value = &pb.XStoreValue_StringValue{StringValue: string(request.Value.(String))}
	case Double:
		val.Value = &pb.XStoreValue_DoubleValue{DoubleValue: float64(request.Value.(Double))}
	case Integer:
		val.Value = &pb.XStoreValue_IntegerValue{IntegerValue: int64(request.Value.(Integer))}
	}

	_, err := client.grpcClient.Set(requestMetadata, &pb.XStoreSetRequest{
		Key:   request.Key,
		Value: &val,
	})
	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	return &responses.StorageSetSuccess{}, nil
}

func (client *storageDataClient) get(ctx context.Context, request *StorageGetRequest) (responses.StorageGetResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

	response, err := client.grpcClient.Get(requestMetadata, &pb.XStoreGetRequest{
		Key: request.Key,
	})

	if err != nil {
		return nil, momentoerrors.ConvertSvcErr(err)
	}

	val := response.GetValue()
	switch val.Value.(type) {
	case *pb.XStoreValue_BytesValue:
		return responses.NewStoreGetSuccess_Bytes(responses.BYTES, val.GetBytesValue()), nil
	case *pb.XStoreValue_StringValue:
		return responses.NewStoreGetSuccess_String(responses.STRING, val.GetStringValue()), nil
	case *pb.XStoreValue_DoubleValue:
		return responses.NewStoreGetSuccess_Double(responses.DOUBLE, val.GetDoubleValue()), nil
	case *pb.XStoreValue_IntegerValue:
		return responses.NewStoreGetSuccess_Integer(responses.INTEGER, int(val.GetIntegerValue())), nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.UnknownServiceError, "Unknown store value type", nil)
	}
}
