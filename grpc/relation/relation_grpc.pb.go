// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: idl/relation.proto

package pbrelation

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

// RelationClient is the client API for Relation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelationClient interface {
	Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error)
	Unfollow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error)
	FollowList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	FollowerList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	FriendList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*FriendsResponse, error)
	IsFollow(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*IsFollowResponse, error)
	IsFriend(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*IsFriendResponse, error)
}

type relationClient struct {
	cc grpc.ClientConnInterface
}

func NewRelationClient(cc grpc.ClientConnInterface) RelationClient {
	return &relationClient{cc}
}

func (c *relationClient) Follow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error) {
	out := new(FollowResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/Follow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) Unfollow(ctx context.Context, in *FollowRequest, opts ...grpc.CallOption) (*FollowResponse, error) {
	out := new(FollowResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/Unfollow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FollowList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FollowList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FollowerList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FollowerList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FriendList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*FriendsResponse, error) {
	out := new(FriendsResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FriendList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) IsFollow(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*IsFollowResponse, error) {
	out := new(IsFollowResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/IsFollow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) IsFriend(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*IsFriendResponse, error) {
	out := new(IsFriendResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/IsFriend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelationServer is the server API for Relation service.
// All implementations must embed UnimplementedRelationServer
// for forward compatibility
type RelationServer interface {
	Follow(context.Context, *FollowRequest) (*FollowResponse, error)
	Unfollow(context.Context, *FollowRequest) (*FollowResponse, error)
	FollowList(context.Context, *ListRequest) (*ListResponse, error)
	FollowerList(context.Context, *ListRequest) (*ListResponse, error)
	FriendList(context.Context, *ListRequest) (*FriendsResponse, error)
	IsFollow(context.Context, *ListRequest) (*IsFollowResponse, error)
	IsFriend(context.Context, *ListRequest) (*IsFriendResponse, error)
	mustEmbedUnimplementedRelationServer()
}

// UnimplementedRelationServer must be embedded to have forward compatible implementations.
type UnimplementedRelationServer struct {
}

func (UnimplementedRelationServer) Follow(context.Context, *FollowRequest) (*FollowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Follow not implemented")
}
func (UnimplementedRelationServer) Unfollow(context.Context, *FollowRequest) (*FollowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unfollow not implemented")
}
func (UnimplementedRelationServer) FollowList(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowList not implemented")
}
func (UnimplementedRelationServer) FollowerList(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowerList not implemented")
}
func (UnimplementedRelationServer) FriendList(context.Context, *ListRequest) (*FriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FriendList not implemented")
}
func (UnimplementedRelationServer) IsFollow(context.Context, *ListRequest) (*IsFollowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFollow not implemented")
}
func (UnimplementedRelationServer) IsFriend(context.Context, *ListRequest) (*IsFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFriend not implemented")
}
func (UnimplementedRelationServer) mustEmbedUnimplementedRelationServer() {}

// UnsafeRelationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelationServer will
// result in compilation errors.
type UnsafeRelationServer interface {
	mustEmbedUnimplementedRelationServer()
}

func RegisterRelationServer(s grpc.ServiceRegistrar, srv RelationServer) {
	s.RegisterService(&Relation_ServiceDesc, srv)
}

func _Relation_Follow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).Follow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/Follow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).Follow(ctx, req.(*FollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_Unfollow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).Unfollow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/Unfollow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).Unfollow(ctx, req.(*FollowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FollowList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FollowList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FollowList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FollowList(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FollowerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FollowerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FollowerList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FollowerList(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FriendList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FriendList(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_IsFollow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).IsFollow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/IsFollow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).IsFollow(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_IsFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).IsFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/IsFriend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).IsFriend(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Relation_ServiceDesc is the grpc.ServiceDesc for Relation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Relation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relation.Relation",
	HandlerType: (*RelationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Follow",
			Handler:    _Relation_Follow_Handler,
		},
		{
			MethodName: "Unfollow",
			Handler:    _Relation_Unfollow_Handler,
		},
		{
			MethodName: "FollowList",
			Handler:    _Relation_FollowList_Handler,
		},
		{
			MethodName: "FollowerList",
			Handler:    _Relation_FollowerList_Handler,
		},
		{
			MethodName: "FriendList",
			Handler:    _Relation_FriendList_Handler,
		},
		{
			MethodName: "IsFollow",
			Handler:    _Relation_IsFollow_Handler,
		},
		{
			MethodName: "IsFriend",
			Handler:    _Relation_IsFriend_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/relation.proto",
}
