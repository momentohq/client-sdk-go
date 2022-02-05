// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: protos/cacheclient.proto

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

// ScsClient is the client API for Scs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScsClient interface {
	Get(ctx context.Context, in *XGetRequest, opts ...grpc.CallOption) (*XGetResponse, error)
	Set(ctx context.Context, in *XSetRequest, opts ...grpc.CallOption) (*XSetResponse, error)
}

type scsClient struct {
	cc grpc.ClientConnInterface
}

func NewScsClient(cc grpc.ClientConnInterface) ScsClient {
	return &scsClient{cc}
}

func (c *scsClient) Get(ctx context.Context, in *XGetRequest, opts ...grpc.CallOption) (*XGetResponse, error) {
	out := new(XGetResponse)
	err := c.cc.Invoke(ctx, "/cache_client.Scs/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scsClient) Set(ctx context.Context, in *XSetRequest, opts ...grpc.CallOption) (*XSetResponse, error) {
	out := new(XSetResponse)
	err := c.cc.Invoke(ctx, "/cache_client.Scs/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScsServer is the server API for Scs service.
// All implementations must embed UnimplementedScsServer
// for forward compatibility
type ScsServer interface {
	Get(context.Context, *XGetRequest) (*XGetResponse, error)
	Set(context.Context, *XSetRequest) (*XSetResponse, error)
	mustEmbedUnimplementedScsServer()
}

// UnimplementedScsServer must be embedded to have forward compatible implementations.
type UnimplementedScsServer struct {
}

func (UnimplementedScsServer) Get(context.Context, *XGetRequest) (*XGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedScsServer) Set(context.Context, *XSetRequest) (*XSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedScsServer) mustEmbedUnimplementedScsServer() {}

// UnsafeScsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScsServer will
// result in compilation errors.
type UnsafeScsServer interface {
	mustEmbedUnimplementedScsServer()
}

func RegisterScsServer(s grpc.ServiceRegistrar, srv ScsServer) {
	s.RegisterService(&Scs_ServiceDesc, srv)
}

func _Scs_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cache_client.Scs/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsServer).Get(ctx, req.(*XGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scs_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScsServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cache_client.Scs/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScsServer).Set(ctx, req.(*XSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Scs_ServiceDesc is the grpc.ServiceDesc for Scs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Scs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cache_client.Scs",
	HandlerType: (*ScsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Scs_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Scs_Set_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/cacheclient.proto",
}
