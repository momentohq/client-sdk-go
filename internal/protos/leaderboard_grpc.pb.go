// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.18.1
// source: leaderboard.proto

package client_sdk_go

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Leaderboard_DeleteLeaderboard_FullMethodName    = "/leaderboard.Leaderboard/DeleteLeaderboard"
	Leaderboard_UpsertElements_FullMethodName       = "/leaderboard.Leaderboard/UpsertElements"
	Leaderboard_RemoveElements_FullMethodName       = "/leaderboard.Leaderboard/RemoveElements"
	Leaderboard_GetLeaderboardLength_FullMethodName = "/leaderboard.Leaderboard/GetLeaderboardLength"
	Leaderboard_GetByRank_FullMethodName            = "/leaderboard.Leaderboard/GetByRank"
	Leaderboard_GetRank_FullMethodName              = "/leaderboard.Leaderboard/GetRank"
	Leaderboard_GetByScore_FullMethodName           = "/leaderboard.Leaderboard/GetByScore"
	Leaderboard_GetCompetitionRank_FullMethodName   = "/leaderboard.Leaderboard/GetCompetitionRank"
)

// LeaderboardClient is the client API for Leaderboard service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Like a sorted set, but for leaderboards!
//
// Elements in a leaderboard are keyed by an ID, which is an unsigned 64 bit integer.
// Scores are single-precision floating point numbers.
//
// Each ID can have only 1 score.
//
// For batchy, multi-element apis, limits are 8192 elements per api call.
//
// Scores are IEEE 754 single-precision floating point numbers. This has a few
// implications you should be aware of, but the one most likely to affect you is that
// below -16777216 and above 16777216, not all integers are able to be represented.
type LeaderboardClient interface {
	// Deletes a leaderboard. After this call, you're not incurring storage cost for this leaderboard anymore.
	DeleteLeaderboard(ctx context.Context, in *XDeleteLeaderboardRequest, opts ...grpc.CallOption) (*XEmpty, error)
	// Insert or update elements in a leaderboard. You can do up to 8192 elements per call.
	// There is no partial failure: Upsert succeeds or fails.
	UpsertElements(ctx context.Context, in *XUpsertElementsRequest, opts ...grpc.CallOption) (*XEmpty, error)
	// Remove up to 8192 elements at a time from a leaderboard. Elements are removed by id.
	RemoveElements(ctx context.Context, in *XRemoveElementsRequest, opts ...grpc.CallOption) (*XEmpty, error)
	// Returns the length of a leaderboard in terms of ID count.
	GetLeaderboardLength(ctx context.Context, in *XGetLeaderboardLengthRequest, opts ...grpc.CallOption) (*XGetLeaderboardLengthResponse, error)
	// Get a range of elements.
	// * Ordinal, 0-based rank.
	// * Range can span up to 8192 elements.
	// See RankRange for details about permissible ranges.
	GetByRank(ctx context.Context, in *XGetByRankRequest, opts ...grpc.CallOption) (*XGetByRankResponse, error)
	// Get the rank of a list of particular ids in the leaderboard.
	// * Ordinal, 0-based rank.
	GetRank(ctx context.Context, in *XGetRankRequest, opts ...grpc.CallOption) (*XGetRankResponse, error)
	// Get a range of elements by a score range.
	// * Ordinal, 0-based rank.
	//
	// You can request up to 8192 elements at a time. To page through many elements that all
	// fall into a score range you can repeatedly invoke this api with the offset parameter.
	GetByScore(ctx context.Context, in *XGetByScoreRequest, opts ...grpc.CallOption) (*XGetByScoreResponse, error)
	// Get the competition ranks of a list of elements.
	// Ranks start at 0. The default ordering is descending.
	// i.e. elements with higher scores have lower ranks.
	GetCompetitionRank(ctx context.Context, in *XGetCompetitionRankRequest, opts ...grpc.CallOption) (*XGetCompetitionRankResponse, error)
}

type leaderboardClient struct {
	cc grpc.ClientConnInterface
}

func NewLeaderboardClient(cc grpc.ClientConnInterface) LeaderboardClient {
	return &leaderboardClient{cc}
}

func (c *leaderboardClient) DeleteLeaderboard(ctx context.Context, in *XDeleteLeaderboardRequest, opts ...grpc.CallOption) (*XEmpty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XEmpty)
	err := c.cc.Invoke(ctx, Leaderboard_DeleteLeaderboard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) UpsertElements(ctx context.Context, in *XUpsertElementsRequest, opts ...grpc.CallOption) (*XEmpty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XEmpty)
	err := c.cc.Invoke(ctx, Leaderboard_UpsertElements_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) RemoveElements(ctx context.Context, in *XRemoveElementsRequest, opts ...grpc.CallOption) (*XEmpty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XEmpty)
	err := c.cc.Invoke(ctx, Leaderboard_RemoveElements_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) GetLeaderboardLength(ctx context.Context, in *XGetLeaderboardLengthRequest, opts ...grpc.CallOption) (*XGetLeaderboardLengthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGetLeaderboardLengthResponse)
	err := c.cc.Invoke(ctx, Leaderboard_GetLeaderboardLength_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) GetByRank(ctx context.Context, in *XGetByRankRequest, opts ...grpc.CallOption) (*XGetByRankResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGetByRankResponse)
	err := c.cc.Invoke(ctx, Leaderboard_GetByRank_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) GetRank(ctx context.Context, in *XGetRankRequest, opts ...grpc.CallOption) (*XGetRankResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGetRankResponse)
	err := c.cc.Invoke(ctx, Leaderboard_GetRank_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) GetByScore(ctx context.Context, in *XGetByScoreRequest, opts ...grpc.CallOption) (*XGetByScoreResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGetByScoreResponse)
	err := c.cc.Invoke(ctx, Leaderboard_GetByScore_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardClient) GetCompetitionRank(ctx context.Context, in *XGetCompetitionRankRequest, opts ...grpc.CallOption) (*XGetCompetitionRankResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(XGetCompetitionRankResponse)
	err := c.cc.Invoke(ctx, Leaderboard_GetCompetitionRank_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LeaderboardServer is the server API for Leaderboard service.
// All implementations must embed UnimplementedLeaderboardServer
// for forward compatibility
//
// Like a sorted set, but for leaderboards!
//
// Elements in a leaderboard are keyed by an ID, which is an unsigned 64 bit integer.
// Scores are single-precision floating point numbers.
//
// Each ID can have only 1 score.
//
// For batchy, multi-element apis, limits are 8192 elements per api call.
//
// Scores are IEEE 754 single-precision floating point numbers. This has a few
// implications you should be aware of, but the one most likely to affect you is that
// below -16777216 and above 16777216, not all integers are able to be represented.
type LeaderboardServer interface {
	// Deletes a leaderboard. After this call, you're not incurring storage cost for this leaderboard anymore.
	DeleteLeaderboard(context.Context, *XDeleteLeaderboardRequest) (*XEmpty, error)
	// Insert or update elements in a leaderboard. You can do up to 8192 elements per call.
	// There is no partial failure: Upsert succeeds or fails.
	UpsertElements(context.Context, *XUpsertElementsRequest) (*XEmpty, error)
	// Remove up to 8192 elements at a time from a leaderboard. Elements are removed by id.
	RemoveElements(context.Context, *XRemoveElementsRequest) (*XEmpty, error)
	// Returns the length of a leaderboard in terms of ID count.
	GetLeaderboardLength(context.Context, *XGetLeaderboardLengthRequest) (*XGetLeaderboardLengthResponse, error)
	// Get a range of elements.
	// * Ordinal, 0-based rank.
	// * Range can span up to 8192 elements.
	// See RankRange for details about permissible ranges.
	GetByRank(context.Context, *XGetByRankRequest) (*XGetByRankResponse, error)
	// Get the rank of a list of particular ids in the leaderboard.
	// * Ordinal, 0-based rank.
	GetRank(context.Context, *XGetRankRequest) (*XGetRankResponse, error)
	// Get a range of elements by a score range.
	// * Ordinal, 0-based rank.
	//
	// You can request up to 8192 elements at a time. To page through many elements that all
	// fall into a score range you can repeatedly invoke this api with the offset parameter.
	GetByScore(context.Context, *XGetByScoreRequest) (*XGetByScoreResponse, error)
	// Get the competition ranks of a list of elements.
	// Ranks start at 0. The default ordering is descending.
	// i.e. elements with higher scores have lower ranks.
	GetCompetitionRank(context.Context, *XGetCompetitionRankRequest) (*XGetCompetitionRankResponse, error)
	mustEmbedUnimplementedLeaderboardServer()
}

// UnimplementedLeaderboardServer must be embedded to have forward compatible implementations.
type UnimplementedLeaderboardServer struct {
}

func (UnimplementedLeaderboardServer) DeleteLeaderboard(context.Context, *XDeleteLeaderboardRequest) (*XEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLeaderboard not implemented")
}
func (UnimplementedLeaderboardServer) UpsertElements(context.Context, *XUpsertElementsRequest) (*XEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertElements not implemented")
}
func (UnimplementedLeaderboardServer) RemoveElements(context.Context, *XRemoveElementsRequest) (*XEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveElements not implemented")
}
func (UnimplementedLeaderboardServer) GetLeaderboardLength(context.Context, *XGetLeaderboardLengthRequest) (*XGetLeaderboardLengthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderboardLength not implemented")
}
func (UnimplementedLeaderboardServer) GetByRank(context.Context, *XGetByRankRequest) (*XGetByRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByRank not implemented")
}
func (UnimplementedLeaderboardServer) GetRank(context.Context, *XGetRankRequest) (*XGetRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRank not implemented")
}
func (UnimplementedLeaderboardServer) GetByScore(context.Context, *XGetByScoreRequest) (*XGetByScoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByScore not implemented")
}
func (UnimplementedLeaderboardServer) GetCompetitionRank(context.Context, *XGetCompetitionRankRequest) (*XGetCompetitionRankResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCompetitionRank not implemented")
}
func (UnimplementedLeaderboardServer) mustEmbedUnimplementedLeaderboardServer() {}

// UnsafeLeaderboardServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LeaderboardServer will
// result in compilation errors.
type UnsafeLeaderboardServer interface {
	mustEmbedUnimplementedLeaderboardServer()
}

func RegisterLeaderboardServer(s grpc.ServiceRegistrar, srv LeaderboardServer) {
	s.RegisterService(&Leaderboard_ServiceDesc, srv)
}

func _Leaderboard_DeleteLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XDeleteLeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).DeleteLeaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_DeleteLeaderboard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).DeleteLeaderboard(ctx, req.(*XDeleteLeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_UpsertElements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XUpsertElementsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).UpsertElements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_UpsertElements_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).UpsertElements(ctx, req.(*XUpsertElementsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_RemoveElements_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XRemoveElementsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).RemoveElements(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_RemoveElements_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).RemoveElements(ctx, req.(*XRemoveElementsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_GetLeaderboardLength_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetLeaderboardLengthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).GetLeaderboardLength(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_GetLeaderboardLength_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).GetLeaderboardLength(ctx, req.(*XGetLeaderboardLengthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_GetByRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetByRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).GetByRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_GetByRank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).GetByRank(ctx, req.(*XGetByRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_GetRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).GetRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_GetRank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).GetRank(ctx, req.(*XGetRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_GetByScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetByScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).GetByScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_GetByScore_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).GetByScore(ctx, req.(*XGetByScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Leaderboard_GetCompetitionRank_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(XGetCompetitionRankRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServer).GetCompetitionRank(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Leaderboard_GetCompetitionRank_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServer).GetCompetitionRank(ctx, req.(*XGetCompetitionRankRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Leaderboard_ServiceDesc is the grpc.ServiceDesc for Leaderboard service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Leaderboard_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leaderboard.Leaderboard",
	HandlerType: (*LeaderboardServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteLeaderboard",
			Handler:    _Leaderboard_DeleteLeaderboard_Handler,
		},
		{
			MethodName: "UpsertElements",
			Handler:    _Leaderboard_UpsertElements_Handler,
		},
		{
			MethodName: "RemoveElements",
			Handler:    _Leaderboard_RemoveElements_Handler,
		},
		{
			MethodName: "GetLeaderboardLength",
			Handler:    _Leaderboard_GetLeaderboardLength_Handler,
		},
		{
			MethodName: "GetByRank",
			Handler:    _Leaderboard_GetByRank_Handler,
		},
		{
			MethodName: "GetRank",
			Handler:    _Leaderboard_GetRank_Handler,
		},
		{
			MethodName: "GetByScore",
			Handler:    _Leaderboard_GetByScore_Handler,
		},
		{
			MethodName: "GetCompetitionRank",
			Handler:    _Leaderboard_GetCompetitionRank_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "leaderboard.proto",
}
