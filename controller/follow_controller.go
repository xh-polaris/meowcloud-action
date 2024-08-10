package controller

import (
	"context"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/service"
)

type IFollowController interface {
	DoFollow(ctx context.Context, req *action.DoFollowReq) (*action.DoFollowResp, error)
	CancelFollow(ctx context.Context, req *action.CancelFollowReq) (*action.CancelFollowResp, error)
	GetFollowedCount(ctx context.Context, req *action.GetFollowedCountReq) (*action.GetFollowedCountResp, error)
	GetFollowedUsers(ctx context.Context, req *action.GetFollowedUsersReq) (*action.GetFollowedUsersResp, error)
	GetUserFollowed(ctx context.Context, req *action.GetUserFollowedReq) (*action.GetUserFollowedResp, error)
	GetFollowed(ctx context.Context, req *action.GetFollowedReq) (*action.GetFollowedResp, error)
}

type FollowController struct {
	followService service.IFollowService
}

func NewFollowController() *FollowController {
	return &FollowController{
		followService: service.NewFollowService(),
	}
}

func (controller *FollowController) DoFollow(ctx context.Context, req *action.DoFollowReq) (*action.DoFollowResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.followService.DoFollow(ctx, req.TargetId, req.TargetType, req.User.UserId)

	return resp, err
}

func (controller *FollowController) CancelFollow(ctx context.Context, req *action.CancelFollowReq) (*action.CancelFollowResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.followService.CancelFollow(ctx, req.TargetId, req.TargetType, req.User.UserId)

	return resp, err
}

func (controller *FollowController) GetFollowedCount(ctx context.Context, req *action.GetFollowedCountReq) (*action.GetFollowedCountResp, error) {

	resp, err := controller.followService.GetFollowedCount(ctx, req.TargetId, req.TargetType)

	return resp, err
}

func (controller *FollowController) GetFollowedUsers(ctx context.Context, req *action.GetFollowedUsersReq) (*action.GetFollowedUsersResp, error) {

	resp, err := controller.followService.GetFollowedUsers(ctx, req.TargetId, req.TargetType, req.PaginationOption)

	return resp, err
}

func (controller *FollowController) GetUserFollowed(ctx context.Context, req *action.GetUserFollowedReq) (*action.GetUserFollowedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.followService.GetUserFollowed(ctx, req.TargetType, userMeta.UserId, req.PaginationOption)

	return resp, err
}

func (controller *FollowController) GetFollowed(ctx context.Context, req *action.GetFollowedReq) (*action.GetFollowedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.followService.GetFollowed(ctx, req.TargetId, req.TargetType, userMeta.UserId)

	return resp, err
}
