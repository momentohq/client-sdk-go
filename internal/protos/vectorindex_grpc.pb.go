// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: vectorindex.proto

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
	VectorIndex_UpsertItemBatch_FullMethodName       = "/vectorindex.VectorIndex/UpsertItemBatch"
	VectorIndex_DeleteItemBatch_FullMethodName       = "/vectorindex.VectorIndex/DeleteItemBatch"
	VectorIndex_Search_FullMethodName                = "/vectorindex.VectorIndex/Search"
	VectorIndex_SearchAndFetchVectors_FullMethodName = "/vectorindex.VectorIndex/SearchAndFetchVectors"
	VectorIndex_GetItemMetadataBatch_FullMethodName  = "/vectorindex.VectorIndex/GetItemMetadataBatch"
	VectorIndex_GetItemBatch_FullMethodName          = "/vectorindex.VectorIndex/GetItemBatch"
	VectorIndex_CountItems_FullMethodName            = "/vectorindex.VectorIndex/CountItems"
)

// VectorIndexClient is the client API for VectorIndex service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VectorIndexClient interface {
	UpsertItemBatch(ctx context.Context, in *XUpsertItemBatchRequest, opts ...grpc.CallOption) (*XUpsertItemBatchResponse, error)
	DeleteItemBatch(ctx context.Context, in *XDeleteItemBatchRequest, opts ...grpc.CallOption) (*XDeleteItemBatchResponse, error)
	Search(ctx context.Context, in *XSearchRequest, opts ...grpc.CallOption) (*XSearchResponse, error)
	SearchAndFetchVectors(ctx context.Context, in *XSearchAndFetchVectorsRequest, opts ...grpc.CallOption) (*XSearchAndFetchVectorsResponse, error)
	GetItemMetadataBatch(ctx context.Context, in *XGetItemMetadataBatchRequest, opts ...grpc.CallOption) (*XGetItemMetadataBatchResponse, error)
	GetItemBatch(ctx context.Context, in *XGetItemBatchRequest, opts ...grpc.CallOption) (*XGetItemBatchResponse, error)
	CountItems(ctx context.Context, in *XCountItemsRequest, opts ...grpc.CallOption) (*XCountItemsResponse, error)
}

type vectorIndexClient struct {
	cc grpc.ClientConnInterface
}

func NewVectorIndexClient(cc grpc.ClientConnInterface) VectorIndexClient {
	return &vectorIndexClient{cc}
}

func (c *vectorIndexClient) UpsertItemBatch(ctx context.Context, in *XUpsertItemBatchRequest, opts ...grpc.CallOption) (*XUpsertItemBatchResponse, error) {
	out := new(XUpsertItemBatchResponse)
	err := c.cc.Invoke(ctx, VectorIndex_UpsertItemBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) DeleteItemBatch(ctx context.Context, in *XDeleteItemBatchRequest, opts ...grpc.CallOption) (*XDeleteItemBatchResponse, error) {
	out := new(XDeleteItemBatchResponse)
	err := c.cc.Invoke(ctx, VectorIndex_DeleteItemBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) Search(ctx context.Context, in *XSearchRequest, opts ...grpc.CallOption) (*XSearchResponse, error) {
	out := new(XSearchResponse)
	err := c.cc.Invoke(ctx, VectorIndex_Search_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) SearchAndFetchVectors(ctx context.Context, in *XSearchAndFetchVectorsRequest, opts ...grpc.CallOption) (*XSearchAndFetchVectorsResponse, error) {
	out := new(XSearchAndFetchVectorsResponse)
	err := c.cc.Invoke(ctx, VectorIndex_SearchAndFetchVectors_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) GetItemMetadataBatch(ctx context.Context, in *XGetItemMetadataBatchRequest, opts ...grpc.CallOption) (*XGetItemMetadataBatchResponse, error) {
	out := new(XGetItemMetadataBatchResponse)
	err := c.cc.Invoke(ctx, VectorIndex_GetItemMetadataBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) GetItemBatch(ctx context.Context, in *XGetItemBatchRequest, opts ...grpc.CallOption) (*XGetItemBatchResponse, error) {
	out := new(XGetItemBatchResponse)
	err := c.cc.Invoke(ctx, VectorIndex_GetItemBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vectorIndexClient) CountItems(ctx context.Context, in *XCountItemsRequest, opts ...grpc.CallOption) (*XCountItemsResponse, error) {
	out := new(XCountItemsResponse)
	err := c.cc.Invoke(ctx, VectorIndex_CountItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VectorIndexServer is the server API for VectorIndex service.
// All implementations must embed UnimplementedVectorIndexServer
// for forward compatibility
type VectorIndexServer interface {
	UpsertItemBatch(context.Context, *XUpsertItemBatchRequest) (*XUpsertItemBatchResponse, error)
	DeleteItemBatch(context.Context, *XDeleteItemBatchRequest) (*XDeleteItemBatchResponse, error)
	Search(context.Context, *XSearchRequest) (*XSearchResponse, error)
	SearchAndFetchVectors(context.Context, *XSearchAndFetchVectorsRequest) (*XSearchAndFetchVectorsResponse, error)
	GetItemMetadataBatch(context.Context, *XGetItemMetadataBatchRequest) (*XGetItemMetadataBatchResponse, error)
	GetItemBatch(context.Context, *XGetItemBatchRequest) (*XGetItemBatchResponse, error)
	CountItems(context.Context, *XCountItemsRequest) (*XCountItemsResponse, error)
	mustEmbedUnimplementedVectorIndexServer()
}

// UnimplementedVectorIndexServer must be embedded to have forward compatible implementations.
type UnimplementedVectorIndexServer struct {
}

func (UnimplementedVectorIndexServer) UpsertItemBatch(context.Context, *XUpsertItemBatchRequest) (*XUpsertItemBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertItemBatch not implemented")
}
func (UnimplementedVectorIndexServer) DeleteItemBatch(context.Context, *XDeleteItemBatchRequest) (*XDeleteItemBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteItemBatch not implemented")
}
func (UnimplementedVectorIndexServer) Search(context.Context, *XSearchRequest) (*XSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedVectorIndexServer) SearchAndFetchVectors(context.Context, *XSearchAndFetchVectorsRequest) (*XSearchAndFetchVectorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAndFetchVectors not implemented")
}
func (UnimplementedVectorIndexServer) GetItemMetadataBatch(context.Context, *XGetItemMetadataBatchRequest) (*XGetItemMetadataBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemMetadataBatch not implemented")
}
func (UnimplementedVectorIndexServer) GetItemBatch(context.Context, *XGetItemBatchRequest) (*XGetItemBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemBatch not implemented")
}
func (UnimplementedVectorIndexServer) CountItems(context.Context, *XCountItemsRequest) (*XCountItemsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountItems not implemented")
}
func (UnimplementedVectorIndexServer) mustEmbedUnimplementedVectorIndexServer() {}

// UnsafeVectorIndexServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VectorIndexServer will
// result in compilation errors.
type UnsafeVectorIndexServer interface {
	mustEmbedUnimplementedVectorIndexServer()
}

func RegisterVectorIndexServer(s grpc.ServiceRegistrar, srv VectorIndexServer) {
	s.RegisterService(&VectorIndex_ServiceDesc, srv)
}

func _VectorIndex_UpsertItemBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XUpsertItemBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).UpsertItemBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_UpsertItemBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).UpsertItemBatch(ctx, req.(*XUpsertItemBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_DeleteItemBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XDeleteItemBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).DeleteItemBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_DeleteItemBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).DeleteItemBatch(ctx, req.(*XDeleteItemBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).Search(ctx, req.(*XSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_SearchAndFetchVectors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XSearchAndFetchVectorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).SearchAndFetchVectors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_SearchAndFetchVectors_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).SearchAndFetchVectors(ctx, req.(*XSearchAndFetchVectorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_GetItemMetadataBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetItemMetadataBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).GetItemMetadataBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_GetItemMetadataBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).GetItemMetadataBatch(ctx, req.(*XGetItemMetadataBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_GetItemBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetItemBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).GetItemBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_GetItemBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).GetItemBatch(ctx, req.(*XGetItemBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VectorIndex_CountItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XCountItemsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VectorIndexServer).CountItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VectorIndex_CountItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VectorIndexServer).CountItems(ctx, req.(*XCountItemsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VectorIndex_ServiceDesc is the grpc.ServiceDesc for VectorIndex service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VectorIndex_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vectorindex.VectorIndex",
	HandlerType: (*VectorIndexServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertItemBatch",
			Handler:    _VectorIndex_UpsertItemBatch_Handler,
		},
		{
			MethodName: "DeleteItemBatch",
			Handler:    _VectorIndex_DeleteItemBatch_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _VectorIndex_Search_Handler,
		},
		{
			MethodName: "SearchAndFetchVectors",
			Handler:    _VectorIndex_SearchAndFetchVectors_Handler,
		},
		{
			MethodName: "GetItemMetadataBatch",
			Handler:    _VectorIndex_GetItemMetadataBatch_Handler,
		},
		{
			MethodName: "GetItemBatch",
			Handler:    _VectorIndex_GetItemBatch_Handler,
		},
		{
			MethodName: "CountItems",
			Handler:    _VectorIndex_CountItems_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "vectorindex.proto",
}
