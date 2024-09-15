package main

import (
	"context"

	pbrelation "diktok/grpc/relation"
	"diktok/package/util"
	"diktok/service/relation/logic"

	"go.uber.org/zap"
)

type RelationService struct {
	pbrelation.UnimplementedRelationServer
}

func (s *RelationService) Follow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	zap.L().Sugar().Infof("[Follow] req = %s", util.GetLogStr(req))
	resp, err := logic.Follow(ctx, req)
	zap.L().Sugar().Infof("[Follow] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) Unfollow(ctx context.Context, req *pbrelation.FollowRequest) (*pbrelation.FollowResponse, error) {
	zap.L().Sugar().Infof("[Unfollow] req = %s", util.GetLogStr(req))
	resp, err := logic.Unfollow(ctx, req)
	zap.L().Sugar().Infof("[Unfollow] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) FollowList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	zap.L().Sugar().Infof("[FollowList] req = %s", util.GetLogStr(req))
	resp, err := logic.FollowList(ctx, req)
	zap.L().Sugar().Infof("[FollowList] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) FollowerList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.ListResponse, error) {
	zap.L().Sugar().Infof("[FollowerList] req = %s", util.GetLogStr(req))
	resp, err := logic.FollowerList(ctx, req)
	zap.L().Sugar().Infof("[FollowerList] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) FriendList(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.FriendsResponse, error) {
	zap.L().Sugar().Infof("[FriendList] req = %s", util.GetLogStr(req))
	resp, err := logic.FriendList(ctx, req)
	zap.L().Sugar().Infof("[FriendList] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) IsFollow(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.IsFollowResponse, error) {
	zap.L().Sugar().Infof("[IsFollow] req = %s", util.GetLogStr(req))
	resp, err := logic.IsFollow(ctx, req)
	zap.L().Sugar().Infof("[IsFollow] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}

func (s *RelationService) IsFriend(ctx context.Context, req *pbrelation.ListRequest) (*pbrelation.IsFriendResponse, error) {
	zap.L().Sugar().Infof("[IsFriend] req = %s", util.GetLogStr(req))
	resp, err := logic.IsFriend(ctx, req)
	zap.L().Sugar().Infof("[IsFriend] resp = %s, err = %e", util.GetLogStr(resp), err)
	return resp, err
}
