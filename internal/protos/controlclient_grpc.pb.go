// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.3
// source: controlclient.proto

package client_sdk_go

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ScsControl_CreateCache_FullMethodName      = "/control_client.ScsControl/CreateCache"
	ScsControl_DeleteCache_FullMethodName      = "/control_client.ScsControl/DeleteCache"
	ScsControl_ListCaches_FullMethodName       = "/control_client.ScsControl/ListCaches"
	ScsControl_FlushCache_FullMethodName       = "/control_client.ScsControl/FlushCache"
	ScsControl_CreateSigningKey_FullMethodName = "/control_client.ScsControl/CreateSigningKey"
	ScsControl_RevokeSigningKey_FullMethodName = "/control_client.ScsControl/RevokeSigningKey"
	ScsControl_ListSigningKeys_FullMethodName  = "/control_client.ScsControl/ListSigningKeys"
	ScsControl_CreateIndex_FullMethodName      = "/control_client.ScsControl/CreateIndex"
	ScsControl_DeleteIndex_FullMethodName      = "/control_client.ScsControl/DeleteIndex"
	ScsControl_ListIndexes_FullMethodName      = "/control_client.ScsControl/ListIndexes"
	ScsControl_CreateStore_FullMethodName      = "/control_client.ScsControl/CreateStore"
	ScsControl_DeleteStore_FullMethodName      = "/control_client.ScsControl/DeleteStore"
	ScsControl_ListStores_FullMethodName       = "/control_client.ScsControl/ListStores"
)

// ScsControlClient is the client API for ScsControl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScsControlClient interface {
	CreateCache(ctx context.Context, in *XCreateCacheRequest, opts ...grpc.CallOption) (*XCreateCacheResponse, error)
	DeleteCache(ctx context.Context, in *XDeleteCacheRequest, opts ...grpc.CallOption) (*XDeleteCacheResponse, error)
	ListCaches(ctx context.Context, in *XListCachesRequest, opts ...grpc.CallOption) (*XListCachesResponse, error)
	FlushCache(ctx context.Context, in *XFlushCacheRequest, opts ...grpc.CallOption) (*XFlushCacheResponse, error)
	CreateSigningKey(ctx context.Context, in *XCreateSigningKeyRequest, opts ...grpc.CallOption) (*XCreateSigningKeyResponse, error)
	RevokeSigningKey(ctx context.Context, in *XRevokeSigningKeyRequest, opts ...grpc.CallOption) (*XRevokeSigningKeyResponse, error)
	ListSigningKeys(ctx context.Context, in *XListSigningKeysRequest, opts ...grpc.CallOption) (*XListSigningKeysResponse, error)
	CreateIndex(ctx context.Context, in *XCreateIndexRequest, opts ...grpc.CallOption) (*XCreateIndexResponse, error)
	DeleteIndex(ctx context.Context, in *XDeleteIndexRequest, opts ...grpc.CallOption) (*XDeleteIndexResponse, error)
	ListIndexes(ctx context.Context, in *XListIndexesRequest, opts ...grpc.CallOption) (*XListIndexesResponse, error)
	CreateStore(ctx context.Context, in *XCreateStoreRequest, opts ...grpc.CallOption) (*XCreateStoreResponse, error)
	DeleteStore(ctx context.Context, in *XDeleteStoreRequest, opts ...grpc.CallOption) (*XDeleteStoreResponse, error)
	ListStores(ctx context.Context, in *XListStoresRequest, opts ...grpc.CallOption) (*XListStoresResponse, error)
}

type scsControlClient struct {
	cc grpc.ClientConnInterface
}

func NewScsControlClient(cc grpc.ClientConnInterface) ScsControlClient {
	return &scsControlClient{cc}
}

func (c *scsControlClient) CreateCache(ctx context.Context, in *XCreateCacheRequest, opts ...grpc.CallOption) (*XCreateCacheResponse, error) {
	out := new(XCreateCacheResponse)
	err := c.cc.Invoke(ctx, ScsControl_CreateCache_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) DeleteCache(ctx context.Context, in *XDeleteCacheRequest, opts ...grpc.CallOption) (*XDeleteCacheResponse, error) {
	out := new(XDeleteCacheResponse)
	err := c.cc.Invoke(ctx, ScsControl_DeleteCache_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) ListCaches(ctx context.Context, in *XListCachesRequest, opts ...grpc.CallOption) (*XListCachesResponse, error) {
	out := new(XListCachesResponse)
	err := c.cc.Invoke(ctx, ScsControl_ListCaches_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) FlushCache(ctx context.Context, in *XFlushCacheRequest, opts ...grpc.CallOption) (*XFlushCacheResponse, error) {
	out := new(XFlushCacheResponse)
	err := c.cc.Invoke(ctx, ScsControl_FlushCache_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) CreateSigningKey(ctx context.Context, in *XCreateSigningKeyRequest, opts ...grpc.CallOption) (*XCreateSigningKeyResponse, error) {
	out := new(XCreateSigningKeyResponse)
	err := c.cc.Invoke(ctx, ScsControl_CreateSigningKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) RevokeSigningKey(ctx context.Context, in *XRevokeSigningKeyRequest, opts ...grpc.CallOption) (*XRevokeSigningKeyResponse, error) {
	out := new(XRevokeSigningKeyResponse)
	err := c.cc.Invoke(ctx, ScsControl_RevokeSigningKey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) ListSigningKeys(ctx context.Context, in *XListSigningKeysRequest, opts ...grpc.CallOption) (*XListSigningKeysResponse, error) {
	out := new(XListSigningKeysResponse)
	err := c.cc.Invoke(ctx, ScsControl_ListSigningKeys_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) CreateIndex(ctx context.Context, in *XCreateIndexRequest, opts ...grpc.CallOption) (*XCreateIndexResponse, error) {
	out := new(XCreateIndexResponse)
	err := c.cc.Invoke(ctx, ScsControl_CreateIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) DeleteIndex(ctx context.Context, in *XDeleteIndexRequest, opts ...grpc.CallOption) (*XDeleteIndexResponse, error) {
	out := new(XDeleteIndexResponse)
	err := c.cc.Invoke(ctx, ScsControl_DeleteIndex_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) ListIndexes(ctx context.Context, in *XListIndexesRequest, opts ...grpc.CallOption) (*XListIndexesResponse, error) {
	out := new(XListIndexesResponse)
	err := c.cc.Invoke(ctx, ScsControl_ListIndexes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) CreateStore(ctx context.Context, in *XCreateStoreRequest, opts ...grpc.CallOption) (*XCreateStoreResponse, error) {
	out := new(XCreateStoreResponse)
	err := c.cc.Invoke(ctx, ScsControl_CreateStore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) DeleteStore(ctx context.Context, in *XDeleteStoreRequest, opts ...grpc.CallOption) (*XDeleteStoreResponse, error) {
	out := new(XDeleteStoreResponse)
	err := c.cc.Invoke(ctx, ScsControl_DeleteStore_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsControlClient) ListStores(ctx context.Context, in *XListStoresRequest, opts ...grpc.CallOption) (*XListStoresResponse, error) {
	out := new(XListStoresResponse)
	err := c.cc.Invoke(ctx, ScsControl_ListStores_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScsControlServer is the server API for ScsControl service.
// All implementations must embed UnimplementedScsControlServer
// for forward compatibility
type ScsControlServer interface {
	CreateCache(context.Context, *XCreateCacheRequest) (*XCreateCacheResponse, error)
	DeleteCache(context.Context, *XDeleteCacheRequest) (*XDeleteCacheResponse, error)
	ListCaches(context.Context, *XListCachesRequest) (*XListCachesResponse, error)
	FlushCache(context.Context, *XFlushCacheRequest) (*XFlushCacheResponse, error)
	CreateSigningKey(context.Context, *XCreateSigningKeyRequest) (*XCreateSigningKeyResponse, error)
	RevokeSigningKey(context.Context, *XRevokeSigningKeyRequest) (*XRevokeSigningKeyResponse, error)
	ListSigningKeys(context.Context, *XListSigningKeysRequest) (*XListSigningKeysResponse, error)
	CreateIndex(context.Context, *XCreateIndexRequest) (*XCreateIndexResponse, error)
	DeleteIndex(context.Context, *XDeleteIndexRequest) (*XDeleteIndexResponse, error)
	ListIndexes(context.Context, *XListIndexesRequest) (*XListIndexesResponse, error)
	CreateStore(context.Context, *XCreateStoreRequest) (*XCreateStoreResponse, error)
	DeleteStore(context.Context, *XDeleteStoreRequest) (*XDeleteStoreResponse, error)
	ListStores(context.Context, *XListStoresRequest) (*XListStoresResponse, error)
	mustEmbedUnimplementedScsControlServer()
}

// UnimplementedScsControlServer must be embedded to have forward compatible implementations.
type UnimplementedScsControlServer struct {
}

func (UnimplementedScsControlServer) CreateCache(context.Context, *XCreateCacheRequest) (*XCreateCacheResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCache not implemented")
}
func (UnimplementedScsControlServer) DeleteCache(context.Context, *XDeleteCacheRequest) (*XDeleteCacheResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCache not implemented")
}
func (UnimplementedScsControlServer) ListCaches(context.Context, *XListCachesRequest) (*XListCachesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCaches not implemented")
}
func (UnimplementedScsControlServer) FlushCache(context.Context, *XFlushCacheRequest) (*XFlushCacheResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FlushCache not implemented")
}
func (UnimplementedScsControlServer) CreateSigningKey(context.Context, *XCreateSigningKeyRequest) (*XCreateSigningKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSigningKey not implemented")
}
func (UnimplementedScsControlServer) RevokeSigningKey(context.Context, *XRevokeSigningKeyRequest) (*XRevokeSigningKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeSigningKey not implemented")
}
func (UnimplementedScsControlServer) ListSigningKeys(context.Context, *XListSigningKeysRequest) (*XListSigningKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSigningKeys not implemented")
}
func (UnimplementedScsControlServer) CreateIndex(context.Context, *XCreateIndexRequest) (*XCreateIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIndex not implemented")
}
func (UnimplementedScsControlServer) DeleteIndex(context.Context, *XDeleteIndexRequest) (*XDeleteIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteIndex not implemented")
}
func (UnimplementedScsControlServer) ListIndexes(context.Context, *XListIndexesRequest) (*XListIndexesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListIndexes not implemented")
}
func (UnimplementedScsControlServer) CreateStore(context.Context, *XCreateStoreRequest) (*XCreateStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStore not implemented")
}
func (UnimplementedScsControlServer) DeleteStore(context.Context, *XDeleteStoreRequest) (*XDeleteStoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStore not implemented")
}
func (UnimplementedScsControlServer) ListStores(context.Context, *XListStoresRequest) (*XListStoresResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStores not implemented")
}
func (UnimplementedScsControlServer) mustEmbedUnimplementedScsControlServer() {}

// UnsafeScsControlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScsControlServer will
// result in compilation errors.
type UnsafeScsControlServer interface {
	mustEmbedUnimplementedScsControlServer()
}

func RegisterScsControlServer(s grpc.ServiceRegistrar, srv ScsControlServer) {
	s.RegisterService(&ScsControl_ServiceDesc, srv)
}

func _ScsControl_CreateCache_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XCreateCacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).CreateCache(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_CreateCache_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).CreateCache(ctx, req.(*XCreateCacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_DeleteCache_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XDeleteCacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).DeleteCache(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_DeleteCache_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).DeleteCache(ctx, req.(*XDeleteCacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_ListCaches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XListCachesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).ListCaches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_ListCaches_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).ListCaches(ctx, req.(*XListCachesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_FlushCache_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XFlushCacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).FlushCache(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_FlushCache_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).FlushCache(ctx, req.(*XFlushCacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_CreateSigningKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XCreateSigningKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).CreateSigningKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_CreateSigningKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).CreateSigningKey(ctx, req.(*XCreateSigningKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_RevokeSigningKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XRevokeSigningKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).RevokeSigningKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_RevokeSigningKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).RevokeSigningKey(ctx, req.(*XRevokeSigningKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_ListSigningKeys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XListSigningKeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).ListSigningKeys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_ListSigningKeys_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).ListSigningKeys(ctx, req.(*XListSigningKeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_CreateIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XCreateIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).CreateIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_CreateIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).CreateIndex(ctx, req.(*XCreateIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_DeleteIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XDeleteIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).DeleteIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_DeleteIndex_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).DeleteIndex(ctx, req.(*XDeleteIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_ListIndexes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XListIndexesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).ListIndexes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_ListIndexes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).ListIndexes(ctx, req.(*XListIndexesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_CreateStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XCreateStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).CreateStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_CreateStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).CreateStore(ctx, req.(*XCreateStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_DeleteStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XDeleteStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).DeleteStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_DeleteStore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).DeleteStore(ctx, req.(*XDeleteStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ScsControl_ListStores_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XListStoresRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsControlServer).ListStores(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ScsControl_ListStores_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsControlServer).ListStores(ctx, req.(*XListStoresRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ScsControl_ServiceDesc is the grpc.ServiceDesc for ScsControl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ScsControl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "control_client.ScsControl",
	HandlerType: (*ScsControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCache",
			Handler:    _ScsControl_CreateCache_Handler,
		},
		{
			MethodName: "DeleteCache",
			Handler:    _ScsControl_DeleteCache_Handler,
		},
		{
			MethodName: "ListCaches",
			Handler:    _ScsControl_ListCaches_Handler,
		},
		{
			MethodName: "FlushCache",
			Handler:    _ScsControl_FlushCache_Handler,
		},
		{
			MethodName: "CreateSigningKey",
			Handler:    _ScsControl_CreateSigningKey_Handler,
		},
		{
			MethodName: "RevokeSigningKey",
			Handler:    _ScsControl_RevokeSigningKey_Handler,
		},
		{
			MethodName: "ListSigningKeys",
			Handler:    _ScsControl_ListSigningKeys_Handler,
		},
		{
			MethodName: "CreateIndex",
			Handler:    _ScsControl_CreateIndex_Handler,
		},
		{
			MethodName: "DeleteIndex",
			Handler:    _ScsControl_DeleteIndex_Handler,
		},
		{
			MethodName: "ListIndexes",
			Handler:    _ScsControl_ListIndexes_Handler,
		},
		{
			MethodName: "CreateStore",
			Handler:    _ScsControl_CreateStore_Handler,
		},
		{
			MethodName: "DeleteStore",
			Handler:    _ScsControl_DeleteStore_Handler,
		},
		{
			MethodName: "ListStores",
			Handler:    _ScsControl_ListStores_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "controlclient.proto",
}
