// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.8.0
// source: pkg/game_server/service.proto

package game_server

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GamesManagerClient is the client API for GamesManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GamesManagerClient interface {
	Join(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinResult, error)
	Move(ctx context.Context, in *Position, opts ...grpc.CallOption) (*MoveResult, error)
	Reconnect(ctx context.Context, in *ReconnectData, opts ...grpc.CallOption) (*ReconnectResult, error)
}

type gamesManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewGamesManagerClient(cc grpc.ClientConnInterface) GamesManagerClient {
	return &gamesManagerClient{cc}
}

func (c *gamesManagerClient) Join(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*JoinResult, error) {
	out := new(JoinResult)
	err := c.cc.Invoke(ctx, "/GamesManager/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamesManagerClient) Move(ctx context.Context, in *Position, opts ...grpc.CallOption) (*MoveResult, error) {
	out := new(MoveResult)
	err := c.cc.Invoke(ctx, "/GamesManager/Move", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gamesManagerClient) Reconnect(ctx context.Context, in *ReconnectData, opts ...grpc.CallOption) (*ReconnectResult, error) {
	out := new(ReconnectResult)
	err := c.cc.Invoke(ctx, "/GamesManager/Reconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GamesManagerServer is the server API for GamesManager service.
// All implementations must embed UnimplementedGamesManagerServer
// for forward compatibility
type GamesManagerServer interface {
	Join(context.Context, *empty.Empty) (*JoinResult, error)
	Move(context.Context, *Position) (*MoveResult, error)
	Reconnect(context.Context, *ReconnectData) (*ReconnectResult, error)
	mustEmbedUnimplementedGamesManagerServer()
}

// UnimplementedGamesManagerServer must be embedded to have forward compatible implementations.
type UnimplementedGamesManagerServer struct {
}

func (UnimplementedGamesManagerServer) Join(context.Context, *empty.Empty) (*JoinResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedGamesManagerServer) Move(context.Context, *Position) (*MoveResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Move not implemented")
}
func (UnimplementedGamesManagerServer) Reconnect(context.Context, *ReconnectData) (*ReconnectResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reconnect not implemented")
}
func (UnimplementedGamesManagerServer) mustEmbedUnimplementedGamesManagerServer() {}

// UnsafeGamesManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GamesManagerServer will
// result in compilation errors.
type UnsafeGamesManagerServer interface {
	mustEmbedUnimplementedGamesManagerServer()
}

func RegisterGamesManagerServer(s grpc.ServiceRegistrar, srv GamesManagerServer) {
	s.RegisterService(&GamesManager_ServiceDesc, srv)
}

func _GamesManager_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamesManagerServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GamesManager/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamesManagerServer).Join(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GamesManager_Move_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Position)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamesManagerServer).Move(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GamesManager/Move",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamesManagerServer).Move(ctx, req.(*Position))
	}
	return interceptor(ctx, in, info, handler)
}

func _GamesManager_Reconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReconnectData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamesManagerServer).Reconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GamesManager/Reconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamesManagerServer).Reconnect(ctx, req.(*ReconnectData))
	}
	return interceptor(ctx, in, info, handler)
}

// GamesManager_ServiceDesc is the grpc.ServiceDesc for GamesManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GamesManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GamesManager",
	HandlerType: (*GamesManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _GamesManager_Join_Handler,
		},
		{
			MethodName: "Move",
			Handler:    _GamesManager_Move_Handler,
		},
		{
			MethodName: "Reconnect",
			Handler:    _GamesManager_Reconnect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/game_server/service.proto",
}

// PlayerClient is the client API for Player service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlayerClient interface {
	UpdateGameState(ctx context.Context, in *Position, opts ...grpc.CallOption) (*UpdateGameStateResult, error)
	YourMove(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*YourMoveResult, error)
	GameFinished(ctx context.Context, in *GameResult, opts ...grpc.CallOption) (*GameFinishedResult, error)
}

type playerClient struct {
	cc grpc.ClientConnInterface
}

func NewPlayerClient(cc grpc.ClientConnInterface) PlayerClient {
	return &playerClient{cc}
}

func (c *playerClient) UpdateGameState(ctx context.Context, in *Position, opts ...grpc.CallOption) (*UpdateGameStateResult, error) {
	out := new(UpdateGameStateResult)
	err := c.cc.Invoke(ctx, "/Player/UpdateGameState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) YourMove(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*YourMoveResult, error) {
	out := new(YourMoveResult)
	err := c.cc.Invoke(ctx, "/Player/YourMove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerClient) GameFinished(ctx context.Context, in *GameResult, opts ...grpc.CallOption) (*GameFinishedResult, error) {
	out := new(GameFinishedResult)
	err := c.cc.Invoke(ctx, "/Player/GameFinished", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlayerServer is the server API for Player service.
// All implementations must embed UnimplementedPlayerServer
// for forward compatibility
type PlayerServer interface {
	UpdateGameState(context.Context, *Position) (*UpdateGameStateResult, error)
	YourMove(context.Context, *empty.Empty) (*YourMoveResult, error)
	GameFinished(context.Context, *GameResult) (*GameFinishedResult, error)
	mustEmbedUnimplementedPlayerServer()
}

// UnimplementedPlayerServer must be embedded to have forward compatible implementations.
type UnimplementedPlayerServer struct {
}

func (UnimplementedPlayerServer) UpdateGameState(context.Context, *Position) (*UpdateGameStateResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGameState not implemented")
}
func (UnimplementedPlayerServer) YourMove(context.Context, *empty.Empty) (*YourMoveResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method YourMove not implemented")
}
func (UnimplementedPlayerServer) GameFinished(context.Context, *GameResult) (*GameFinishedResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GameFinished not implemented")
}
func (UnimplementedPlayerServer) mustEmbedUnimplementedPlayerServer() {}

// UnsafePlayerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlayerServer will
// result in compilation errors.
type UnsafePlayerServer interface {
	mustEmbedUnimplementedPlayerServer()
}

func RegisterPlayerServer(s grpc.ServiceRegistrar, srv PlayerServer) {
	s.RegisterService(&Player_ServiceDesc, srv)
}

func _Player_UpdateGameState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Position)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).UpdateGameState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Player/UpdateGameState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).UpdateGameState(ctx, req.(*Position))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_YourMove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).YourMove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Player/YourMove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).YourMove(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Player_GameFinished_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GameResult)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).GameFinished(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Player/GameFinished",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).GameFinished(ctx, req.(*GameResult))
	}
	return interceptor(ctx, in, info, handler)
}

// Player_ServiceDesc is the grpc.ServiceDesc for Player service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Player_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Player",
	HandlerType: (*PlayerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateGameState",
			Handler:    _Player_UpdateGameState_Handler,
		},
		{
			MethodName: "YourMove",
			Handler:    _Player_YourMove_Handler,
		},
		{
			MethodName: "GameFinished",
			Handler:    _Player_GameFinished_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/game_server/service.proto",
}
