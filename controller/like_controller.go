package controller

import (
	"context"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/service"
)

type ILikeController interface {
	DoLike(ctx context.Context, req *action.DoLikeReq) (res *action.DoLikeResp, err error)
	CancelLike(ctx context.Context, req *action.CancelLikeReq) (*action.CancelLikeResp, error)
	GetLikedCount(ctx context.Context, req *action.GetLikedCountReq) (*action.GetLikedCountResp, error)
	GetLikedUsers(ctx context.Context, req *action.GetLikedUsersReq) (*action.GetLikedUsersResp, error)
	GetUserLiked(ctx context.Context, req *action.GetUserLikedReq) (*action.GetUserLikedResp, error)
	GetLiked(ctx context.Context, req *action.GetLikedReq) (*action.GetLikedResp, error)
}

type LikeController struct {
	likeService service.ILikeService
}

func NewLikeController() *LikeController {
	return &LikeController{
		likeService: service.NewLikeService(),
	}
}

func (controller *LikeController) DoLike(ctx context.Context, req *action.DoLikeReq) (*action.DoLikeResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.DoLike(ctx, req.TargetId, req.TargetType, req.User.UserId)

	return resp, err
}

func (controller *LikeController) CancelLike(ctx context.Context, req *action.CancelLikeReq) (*action.CancelLikeResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.CancelLike(ctx, req.TargetId, req.TargetType, req.User.UserId)

	return resp, err
}

func (controller *LikeController) GetLikedCount(ctx context.Context, req *action.GetLikedCountReq) (*action.GetLikedCountResp, error) {

	resp, err := controller.likeService.GetLikedCount(ctx, req.TargetId, req.TargetType)

	return resp, err
}

func (controller *LikeController) GetLikedUsers(ctx context.Context, req *action.GetLikedUsersReq) (*action.GetLikedUsersResp, error) {

	resp, err := controller.likeService.GetLikedUsers(ctx, req.TargetId, req.TargetType, req.PaginationOption)

	return resp, err
}

func (controller *LikeController) GetUserLiked(ctx context.Context, req *action.GetUserLikedReq) (*action.GetUserLikedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.GetUserLiked(ctx, req.TargetType, userMeta.UserId, req.PaginationOption)

	return resp, err
}

func (controller *LikeController) GetLiked(ctx context.Context, req *action.GetLikedReq) (*action.GetLikedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.GetLiked(ctx, req.TargetId, req.TargetType, userMeta.UserId)

	return resp, err
}
