package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"meowcloud-action/common/consts"
	"meowcloud-action/infra/mapper/like"
)

type ILikeService interface {
	DoLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoLikeResp, error)
	CancelLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.CancelLikeResp, error)
	GetLikedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetLikedCountResp, error)
	GetLikedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetLikedUsersResp, error)
	GetUserLiked(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserLikedResp, error)
	GetLiked(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetLikedResp, error)
}

type LikeService struct {
	LikeMongoMapper like.IMongoMapper
}

func NewLikeService() ILikeService {
	mongoMapper := like.NewMongoMapper()
	return &LikeService{
		LikeMongoMapper: mongoMapper,
	}
}

func (service *LikeService) DoLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.DoLikeResp, error) {

	// 判断是否点赞过，不存在和未点赞都为true
	liked, err := service.LikeMongoMapper.IsLiked(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	// 点赞过则抛出异常
	if !liked {
		return nil, consts.RepeatLike
	}

	err = service.LikeMongoMapper.InsertOne(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, consts.TryAgain
	}

	return &action.DoLikeResp{}, nil
}

func (service *LikeService) CancelLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.CancelLikeResp, error) {

	// 判断是否点赞过，不存在和未点赞都是true
	liked, err := service.LikeMongoMapper.IsLiked(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	// 点赞过则抛出异常
	if liked {
		return nil, consts.LikeNotExist
	}

	err = service.LikeMongoMapper.CancelLike(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, consts.TryAgain
	}

	return &action.CancelLikeResp{}, nil
}

func (service *LikeService) GetLikedCount(ctx context.Context, targetId string, targetType action.TargetType) (*action.GetLikedCountResp, error) {
	count, err := service.LikeMongoMapper.CountLikes(ctx, targetId, targetType)

	if err != nil {
		return nil, err
	}

	return &action.GetLikedCountResp{Count: count}, nil
}

func (service *LikeService) GetLikedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) (*action.GetLikedUsersResp, error) {
	data, total, err := service.LikeMongoMapper.GetLikedUsers(ctx, targetId, targetType, options)

	if err != nil {
		return nil, err
	}

	var likes []*action.Action_Like
	for _, val := range data {
		aLike := &action.Action_Like{}
		err := copier.Copy(aLike, val)
		if err != nil {
			return nil, err
		}
		aLike.Id = val.ID.Hex()
		aLike.CreateAt = val.CreateAt.Unix()
		likes = append(likes, aLike)
	}

	return &action.GetLikedUsersResp{
		Likes: likes,
		Total: total,
	}, nil
}

func (service *LikeService) GetUserLiked(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) (*action.GetUserLikedResp, error) {
	data, total, err := service.LikeMongoMapper.GetUserLiked(ctx, targetType, userId, options)

	if err != nil {
		return nil, err
	}

	var likes []*action.Action_Like
	for _, val := range data {
		aLike := &action.Action_Like{}
		err := copier.Copy(aLike, val)
		if err != nil {
			return nil, err
		}
		aLike.Id = val.ID.Hex()
		aLike.CreateAt = val.CreateAt.Unix()
		likes = append(likes, aLike)
	}

	return &action.GetUserLikedResp{
		Likes: likes,
		Total: total,
	}, nil
}

func (service *LikeService) GetLiked(ctx context.Context, targetId string, targetType action.TargetType, userId string) (*action.GetLikedResp, error) {
	liked, err := service.LikeMongoMapper.IsLiked(ctx, targetId, targetType, userId)

	if err != nil {
		return nil, err
	}

	return &action.GetLikedResp{Liked: liked}, nil
}
