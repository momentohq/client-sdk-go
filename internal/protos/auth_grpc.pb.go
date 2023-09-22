// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.20.3
// source: auth.proto

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
	Auth_Login_FullMethodName            = "/auth.Auth/Login"
	Auth_GenerateApiToken_FullMethodName = "/auth.Auth/GenerateApiToken"
	Auth_RefreshApiToken_FullMethodName  = "/auth.Auth/RefreshApiToken"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Login(ctx context.Context, in *XLoginRequest, opts ...grpc.CallOption) (Auth_LoginClient, error)
	// api for initially generating api and refresh tokens
	GenerateApiToken(ctx context.Context, in *XGenerateApiTokenRequest, opts ...grpc.CallOption) (*XGenerateApiTokenResponse, error)
	// api for programmatically refreshing api and refresh tokens
	RefreshApiToken(ctx context.Context, in *XRefreshApiTokenRequest, opts ...grpc.CallOption) (*XRefreshApiTokenResponse, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Login(ctx context.Context, in *XLoginRequest, opts ...grpc.CallOption) (Auth_LoginClient, error) {
	stream, err := c.cc.NewStream(ctx, &Auth_ServiceDesc.Streams[0], Auth_Login_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &authLoginClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Auth_LoginClient interface {
	Recv() (*XLoginResponse, error)
	grpc.ClientStream
}

type authLoginClient struct {
	grpc.ClientStream
}

func (x *authLoginClient) Recv() (*XLoginResponse, error) {
	m := new(XLoginResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *authClient) GenerateApiToken(ctx context.Context, in *XGenerateApiTokenRequest, opts ...grpc.CallOption) (*XGenerateApiTokenResponse, error) {
	out := new(XGenerateApiTokenResponse)
	err := c.cc.Invoke(ctx, Auth_GenerateApiToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) RefreshApiToken(ctx context.Context, in *XRefreshApiTokenRequest, opts ...grpc.CallOption) (*XRefreshApiTokenResponse, error) {
	out := new(XRefreshApiTokenResponse)
	err := c.cc.Invoke(ctx, Auth_RefreshApiToken_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	Login(*XLoginRequest, Auth_LoginServer) error
	// api for initially generating api and refresh tokens
	GenerateApiToken(context.Context, *XGenerateApiTokenRequest) (*XGenerateApiTokenResponse, error)
	// api for programmatically refreshing api and refresh tokens
	RefreshApiToken(context.Context, *XRefreshApiTokenRequest) (*XRefreshApiTokenResponse, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) Login(*XLoginRequest, Auth_LoginServer) error {
	return status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServer) GenerateApiToken(context.Context, *XGenerateApiTokenRequest) (*XGenerateApiTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateApiToken not implemented")
}
func (UnimplementedAuthServer) RefreshApiToken(context.Context, *XRefreshApiTokenRequest) (*XRefreshApiTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshApiToken not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Login_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(XLoginRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuthServer).Login(m, &authLoginServer{stream})
}

type Auth_LoginServer interface {
	Send(*XLoginResponse) error
	grpc.ServerStream
}

type authLoginServer struct {
	grpc.ServerStream
}

func (x *authLoginServer) Send(m *XLoginResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Auth_GenerateApiToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGenerateApiTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).GenerateApiToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_GenerateApiToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).GenerateApiToken(ctx, req.(*XGenerateApiTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_RefreshApiToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XRefreshApiTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).RefreshApiToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_RefreshApiToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).RefreshApiToken(ctx, req.(*XRefreshApiTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateApiToken",
			Handler:    _Auth_GenerateApiToken_Handler,
		},
		{
			MethodName: "RefreshApiToken",
			Handler:    _Auth_RefreshApiToken_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Login",
			Handler:       _Auth_Login_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "auth.proto",
}
