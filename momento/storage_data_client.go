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

	requestMetadata := internal.CreateMetadata(ctx, internal.Store, "store", request.StoreName)

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

	requestMetadata := internal.CreateMetadata(ctx, internal.Store, "store", request.StoreName)

	val := pb.XStoreValue{}
	switch request.Value.(type) {
	case StorageValueBytes:
		val.Value = &pb.XStoreValue_BytesValue{BytesValue: request.Value.(StorageValueBytes)}
	case StorageValueString:
		val.Value = &pb.XStoreValue_StringValue{StringValue: string(request.Value.(StorageValueString))}
	case StorageValueDouble:
		val.Value = &pb.XStoreValue_DoubleValue{DoubleValue: float64(request.Value.(StorageValueDouble))}
	case StorageValueInteger:
		val.Value = &pb.XStoreValue_IntegerValue{IntegerValue: int64(request.Value.(StorageValueInteger))}
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

	requestMetadata := internal.CreateMetadata(ctx, internal.Store, "store", request.StoreName)

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
		return nil, momentoerrors.ConvertSvcErr(err, header, trailer)
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
