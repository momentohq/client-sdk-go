package momento

import (
	"context"
	"time"

	"github.com/momentohq/client-sdk-go/utils"

	"github.com/momentohq/client-sdk-go/internal/grpcmanagers"
	"github.com/momentohq/client-sdk-go/internal/models"
	"github.com/momentohq/client-sdk-go/internal/momentoerrors"
	pb "github.com/momentohq/client-sdk-go/internal/protos"
	"github.com/momentohq/client-sdk-go/responses"
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

func (*storageDataClient) CreateNewMetadata(storeName string) metadata.MD {
	return metadata.Pairs("store", storeName)
}

func (client *storageDataClient) delete(ctx context.Context, request *StorageDeleteRequest) (responses.StorageDeleteResponse, momentoerrors.MomentoSvcErr) {
	ctx, cancel := context.WithTimeout(ctx, client.requestTimeout)
	defer cancel()

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

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

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

	val := pb.XStoreValue{}
	switch request.Value.(type) {
	case utils.StorageValueBytes:
		val.Value = &pb.XStoreValue_BytesValue{BytesValue: request.Value.(utils.StorageValueBytes)}
	case utils.StorageValueString:
		val.Value = &pb.XStoreValue_StringValue{StringValue: string(request.Value.(utils.StorageValueString))}
	case utils.StorageValueFloat:
		val.Value = &pb.XStoreValue_DoubleValue{DoubleValue: float64(request.Value.(utils.StorageValueFloat))}
	case utils.StorageValueInt:
		val.Value = &pb.XStoreValue_IntegerValue{IntegerValue: int64(request.Value.(utils.StorageValueInt))}
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

	requestMetadata := metadata.NewOutgoingContext(
		ctx, client.CreateNewMetadata(request.StoreName),
	)

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
			return responses.NewStoreGetNotFound(), nil
		}
		return nil, myErr
	}

	val := response.GetValue()
	switch val.Value.(type) {
	case *pb.XStoreValue_BytesValue:
		return responses.NewStoreGetFound_Bytes(val.GetBytesValue()), nil
	case *pb.XStoreValue_StringValue:
		return responses.NewStoreGetFound_String(val.GetStringValue()), nil
	case *pb.XStoreValue_DoubleValue:
		return responses.NewStoreGetFound_Float(val.GetDoubleValue()), nil
	case *pb.XStoreValue_IntegerValue:
		return responses.NewStoreGetFound_Integer(int(val.GetIntegerValue())), nil
	default:
		return nil, momentoerrors.NewMomentoSvcErr(momentoerrors.UnknownServiceError, "Unknown store value type", nil)
	}
}
