// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: idl/favorite.proto

package pbfavorite

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

// FavoriteClient is the client API for Favorite service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FavoriteClient interface {
	Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error)
	Unlike(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error)
	Count(ctx context.Context, in *CountReq, opts ...grpc.CallOption) (*CountResp, error)
}

type favoriteClient struct {
	cc grpc.ClientConnInterface
}

func NewFavoriteClient(cc grpc.ClientConnInterface) FavoriteClient {
	return &favoriteClient{cc}
}

func (c *favoriteClient) Like(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error) {
	out := new(LikeResponse)
	err := c.cc.Invoke(ctx, "/favorite.Favorite/Like", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) Unlike(ctx context.Context, in *LikeRequest, opts ...grpc.CallOption) (*LikeResponse, error) {
	out := new(LikeResponse)
	err := c.cc.Invoke(ctx, "/favorite.Favorite/Unlike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/favorite.Favorite/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) IsFavorite(ctx context.Context, in *IsFavoriteReq, opts ...grpc.CallOption) (*IsFavoriteResp, error) {
	out := new(IsFavoriteResp)
	err := c.cc.Invoke(ctx, "/favorite.Favorite/IsFavorite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *favoriteClient) Count(ctx context.Context, in *CountReq, opts ...grpc.CallOption) (*CountResp, error) {
	out := new(CountResp)
	err := c.cc.Invoke(ctx, "/favorite.Favorite/Count", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FavoriteServer is the server API for Favorite service.
// All implementations must embed UnimplementedFavoriteServer
// for forward compatibility
type FavoriteServer interface {
	Like(context.Context, *LikeRequest) (*LikeResponse, error)
	Unlike(context.Context, *LikeRequest) (*LikeResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	IsFavorite(context.Context, *IsFavoriteReq) (*IsFavoriteResp, error)
	Count(context.Context, *CountReq) (*CountResp, error)
	mustEmbedUnimplementedFavoriteServer()
}

// UnimplementedFavoriteServer must be embedded to have forward compatible implementations.
type UnimplementedFavoriteServer struct {
}

func (UnimplementedFavoriteServer) Like(context.Context, *LikeRequest) (*LikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Like not implemented")
}
func (UnimplementedFavoriteServer) Unlike(context.Context, *LikeRequest) (*LikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unlike not implemented")
}
func (UnimplementedFavoriteServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedFavoriteServer) IsFavorite(context.Context, *IsFavoriteReq) (*IsFavoriteResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFavorite not implemented")
}
func (UnimplementedFavoriteServer) Count(context.Context, *CountReq) (*CountResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}
func (UnimplementedFavoriteServer) mustEmbedUnimplementedFavoriteServer() {}

// UnsafeFavoriteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FavoriteServer will
// result in compilation errors.
type UnsafeFavoriteServer interface {
	mustEmbedUnimplementedFavoriteServer()
}

func RegisterFavoriteServer(s grpc.ServiceRegistrar, srv FavoriteServer) {
	s.RegisterService(&Favorite_ServiceDesc, srv)
}

func _Favorite_Like_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).Like(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/favorite.Favorite/Like",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).Like(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_Unlike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LikeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).Unlike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/favorite.Favorite/Unlike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).Unlike(ctx, req.(*LikeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/favorite.Favorite/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_IsFavorite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsFavoriteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).IsFavorite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/favorite.Favorite/IsFavorite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).IsFavorite(ctx, req.(*IsFavoriteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Favorite_Count_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CountReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FavoriteServer).Count(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/favorite.Favorite/Count",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FavoriteServer).Count(ctx, req.(*CountReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Favorite_ServiceDesc is the grpc.ServiceDesc for Favorite service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Favorite_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "favorite.Favorite",
	HandlerType: (*FavoriteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Like",
			Handler:    _Favorite_Like_Handler,
		},
		{
			MethodName: "Unlike",
			Handler:    _Favorite_Unlike_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Favorite_List_Handler,
		},
		{
			MethodName: "IsFavorite",
			Handler:    _Favorite_IsFavorite_Handler,
		},
		{
			MethodName: "Count",
			Handler:    _Favorite_Count_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/favorite.proto",
}
