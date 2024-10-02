// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.18.1
// source: token.proto

package client_sdk_go

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Token_GenerateDisposableToken_FullMethodName = "/token.Token/GenerateDisposableToken"
)

// TokenClient is the client API for Token service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenClient interface {
	GenerateDisposableToken(ctx context.Context, in *XGenerateDisposableTokenRequest, opts ...grpc.CallOption) (*XGenerateDisposableTokenResponse, error)
}

type tokenClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenClient(cc grpc.ClientConnInterface) TokenClient {
	return &tokenClient{cc}
}

func (c *tokenClient) GenerateDisposableToken(ctx context.Context, in *XGenerateDisposableTokenRequest, opts ...grpc.CallOption) (*XGenerateDisposableTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGenerateDisposableTokenResponse)
	err := c.cc.Invoke(ctx, Token_GenerateDisposableToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServer is the server API for Token service.
// All implementations must embed UnimplementedTokenServer
// for forward compatibility.
type TokenServer interface {
	GenerateDisposableToken(context.Context, *XGenerateDisposableTokenRequest) (*XGenerateDisposableTokenResponse, error)
	mustEmbedUnimplementedTokenServer()
}

// UnimplementedTokenServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTokenServer struct{}

func (UnimplementedTokenServer) GenerateDisposableToken(context.Context, *XGenerateDisposableTokenRequest) (*XGenerateDisposableTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateDisposableToken not implemented")
}
func (UnimplementedTokenServer) mustEmbedUnimplementedTokenServer() {}
func (UnimplementedTokenServer) testEmbeddedByValue()               {}

// UnsafeTokenServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenServer will
// result in compilation errors.
type UnsafeTokenServer interface {
	mustEmbedUnimplementedTokenServer()
}

func RegisterTokenServer(s grpc.ServiceRegistrar, srv TokenServer) {
	// If the following call pancis, it indicates UnimplementedTokenServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Token_ServiceDesc, srv)
}

func _Token_GenerateDisposableToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGenerateDisposableTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServer).GenerateDisposableToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Token_GenerateDisposableToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServer).GenerateDisposableToken(ctx, req.(*XGenerateDisposableTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Token_ServiceDesc is the grpc.ServiceDesc for Token service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Token_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "token.Token",
	HandlerType: (*TokenServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateDisposableToken",
			Handler:    _Token_GenerateDisposableToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token.proto",
}
