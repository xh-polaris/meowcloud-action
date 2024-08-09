package controller

import (
	"context"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/service"
)

type IShareController interface {
	DoShare(ctx context.Context, req *action.DoShareReq) (*action.DoShareResp, error)
	GetSharedCount(ctx context.Context, req *action.GetSharedCountReq) (*action.GetSharedCountResp, error)
	GetSharedUsers(ctx context.Context, req *action.GetSharedUsersReq) (*action.GetSharedUsersResp, error)
	GetUserShared(ctx context.Context, req *action.GetUserSharedReq) (*action.GetUserSharedResp, error)
	GetShared(ctx context.Context, req *action.GetSharedReq) (*action.GetSharedResp, error)
}

type ShareController struct {
	likeService service.ShareService
}

func (controller *ShareController) DoShare(ctx context.Context, req *action.DoShareReq) (*action.DoShareResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.DoShare(ctx, req.TargetId, req.TargetType, req.User.UserId)

	return resp, err
}

func (controller *ShareController) GetSharedCount(ctx context.Context, req *action.GetSharedCountReq) (*action.GetSharedCountResp, error) {

	resp, err := controller.likeService.GetSharedCount(ctx, req.TargetId, req.TargetType)

	return resp, err
}

func (controller *ShareController) GetSharedUsers(ctx context.Context, req *action.GetSharedUsersReq) (*action.GetSharedUsersResp, error) {

	resp, err := controller.likeService.GetSharedUsers(ctx, req.TargetId, req.TargetType)

	return resp, err
}

func (controller *ShareController) GetUserShared(ctx context.Context, req *action.GetUserSharedReq) (*action.GetUserSharedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.GetUserShared(ctx, req.TargetType, userMeta.UserId)

	return resp, err
}

func (controller *ShareController) GetShared(ctx context.Context, req *action.GetSharedReq) (*action.GetSharedResp, error) {
	userMeta := req.User

	// 用户信息校验
	userErr := consts.CheckUserMeta(userMeta)
	if userErr != nil {
		return nil, userErr
	}

	resp, err := controller.likeService.GetShared(ctx, req.TargetId, req.TargetType, userMeta.UserId)

	return resp, err
}
