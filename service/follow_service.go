package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/infra/mapper/follow"
)

type IFollowService interface {
	DoFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoFollowResp, error)
	CancelFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.CancelFollowResp, error)
	GetFollowedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetFollowedCountResp, error)
	GetFollowedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetFollowedUsersResp, error)
	GetUserFollowed(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserFollowedResp, error)
	GetFollowed(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetFollowedResp, error)
}

type FollowService struct {
	FollowMongoMapper follow.IMongoMapper
}

func NewFollowService() IFollowService {
	mongoMapper := follow.NewMongoMapper()
	return &FollowService{
		FollowMongoMapper: mongoMapper,
	}
}

func (service FollowService) DoFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoFollowResp, error) {

	// 判断是否点赞过，不存在和未点赞都为false
	followed, err := service.FollowMongoMapper.IsFollowed(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	// 点赞过则抛出异常
	if followed {
		return nil, consts.RepeatFollow
	}

	err = service.FollowMongoMapper.InsertOne(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, consts.TryAgain
	}

	return &action.DoFollowResp{}, nil
}

func (service FollowService) CancelFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.CancelFollowResp, error) {

	// 判断是否点赞过，不存在和未点赞都为false
	followed, err := service.FollowMongoMapper.IsFollowed(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	// 未点赞过则抛出异常
	if !followed {
		return nil, consts.FollowNotExist
	}

	err = service.FollowMongoMapper.CancelFollow(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, consts.TryAgain
	}

	return &action.CancelFollowResp{}, nil
}

func (service FollowService) GetFollowedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetFollowedCountResp, error) {
	count, err := service.FollowMongoMapper.CountFollows(ctx, targetId, targetType)

	if err != nil {
		return nil, err
	}

	return &action.GetFollowedCountResp{Count: count}, nil
}

func (service FollowService) GetFollowedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetFollowedUsersResp, error) {
	data, total, err := service.FollowMongoMapper.GetFollowedUsers(ctx, targetId, targetType, options)

	if err != nil {
		return nil, err
	}

	var follows []*action.Action_Follow
	for _, val := range data {
		aFollow := &action.Action_Follow{}
		err := copier.Copy(aFollow, val)
		if err != nil {
			return nil, err
		}
		aFollow.Id = val.ID.Hex()
		aFollow.CreateAt = val.CreateAt.Unix()
		follows = append(follows, aFollow)
	}

	return &action.GetFollowedUsersResp{
		Follows: follows,
		Total:   total,
	}, nil
}

func (service FollowService) GetUserFollowed(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserFollowedResp, error) {
	data, total, err := service.FollowMongoMapper.GetUserFollowed(ctx, targetType, userId, options)

	if err != nil {
		return nil, err
	}

	var follows []*action.Action_Follow
	for _, val := range data {
		aFollow := &action.Action_Follow{}
		err := copier.Copy(aFollow, val)
		if err != nil {
			return nil, err
		}
		aFollow.Id = val.ID.Hex()
		aFollow.CreateAt = val.CreateAt.Unix()
		follows = append(follows, aFollow)
	}

	return &action.GetUserFollowedResp{
		Follows: follows,
		Total:   total,
	}, nil
}

func (service FollowService) GetFollowed(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetFollowedResp, error) {
	followed, err := service.FollowMongoMapper.IsFollowed(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	return &action.GetFollowedResp{Followed: followed}, nil
}
