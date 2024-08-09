package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/infra/mapper/share"
)

type IShareService interface {
	DoShare(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoShareResp, error)
	GetSharedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetSharedCountResp, error)
	GetSharedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetSharedUsersResp, error)
	GetUserShared(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserSharedResp, error)
	GetShared(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetSharedResp, error)
}

type ShareService struct {
	ShareMongoMapper share.IMongoMapper
}

func (service ShareService) DoShare(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoShareResp, error) {

	err := service.ShareMongoMapper.InsertOne(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, consts.TryAgain
	}

	return &action.DoShareResp{}, nil
}

func (service ShareService) GetSharedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetSharedCountResp, error) {
	count, err := service.ShareMongoMapper.CountShares(ctx, targetId, targetType)

	if err != nil {
		return nil, err
	}

	return &action.GetSharedCountResp{Count: count}, nil
}

func (service ShareService) GetSharedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetSharedUsersResp, error) {
	data, total, err := service.ShareMongoMapper.GetSharedUsers(ctx, targetId, targetType, options)

	if err != nil {
		return nil, err
	}

	var shares []*action.Action_Share
	for _, val := range data {
		aShare := &action.Action_Share{}
		err := copier.Copy(aShare, val)
		if err != nil {
			return nil, err
		}
		aShare.Id = val.ID.Hex()
		aShare.CreateAt = val.CreateAt.Unix()
		shares = append(shares, aShare)
	}

	return &action.GetSharedUsersResp{
		Shares: shares,
		Total:  total,
	}, nil
}

func (service ShareService) GetUserShared(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserSharedResp, error) {
	data, total, err := service.ShareMongoMapper.GetUserShared(ctx, targetType, userId, options)

	if err != nil {
		return nil, err
	}

	var shares []*action.Action_Share
	for _, val := range data {
		aShare := &action.Action_Share{}
		err := copier.Copy(aShare, val)
		if err != nil {
			return nil, err
		}
		aShare.Id = val.ID.Hex()
		aShare.CreateAt = val.CreateAt.Unix()
		shares = append(shares, aShare)
	}

	return &action.GetUserSharedResp{
		Shares: shares,
		Total:  total,
	}, nil
}

func (service ShareService) GetShared(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetSharedResp, error) {
	shared, err := service.ShareMongoMapper.IsShared(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	return &action.GetSharedResp{Shared: shared}, nil
}
