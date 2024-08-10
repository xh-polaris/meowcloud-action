package follow

import (
	"context"
	"errors"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"github.com/zeromicro/go-zero/core/stores/monc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"meowcloud-action/common/config"
	"time"
)

const prefixFollowCacheKey = "cache:follow"
const CollectionName = "follow"

// 用于检查接口是否实现
var _ IMongoMapper = (*MongoMapper)(nil)

type IMongoMapper interface {
	InsertOne(ctx context.Context, targetId string, targetType action.TargetType, userId string) error
	IsFollowed(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error)
	CancelFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) error
	CountFollows(ctx context.Context, targetId string, targetType action.TargetType) (int64, error)
	GetFollowedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) ([]*Follow, int64, error)
	GetUserFollowed(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) ([]*Follow, int64, error)
	CountFollowsByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error)
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper() IMongoMapper {
	aConfig := config.Get()
	conn := monc.MustNewModel(aConfig.Mongo.URL, aConfig.Mongo.DB, CollectionName, *aConfig.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) InsertOne(ctx context.Context, targetId string, targetType action.TargetType, userId string) error {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var follow Follow

	err := m.conn.FindOne(ctx, targetId+targetType.String()+userId, &follow, filter)

	switch {
	// 已经存在则修改isCancel状态
	case err == nil:
		follow.IsCancel = false
		_, err = m.conn.ReplaceOneNoCache(ctx, filter, follow)
	// 不存在则新建
	case errors.Is(err, monc.ErrNotFound):
		newFollow := &Follow{
			ID:         primitive.NewObjectID(),
			TargetId:   targetId,
			TargetType: targetType,
			UserId:     userId,
			IsCancel:   false,
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
		}
		key := prefixFollowCacheKey + newFollow.TargetId + newFollow.ID.Hex()
		_, err = m.conn.InsertOne(ctx, key, newFollow)
	}
	return err

}

func (m *MongoMapper) IsFollowed(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error) {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var follow Follow

	err := m.conn.FindOneNoCache(ctx, &follow, filter)
	switch {
	case errors.Is(err, monc.ErrNotFound):
		return false, nil
	case err == nil:
		return follow.IsCancel, nil
	default:
		return false, err
	}
}

func (m *MongoMapper) CancelFollow(ctx context.Context, targetId string, targetType action.TargetType, userId string) error {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var follow Follow

	err := m.conn.FindOneNoCache(ctx, &follow, filter)

	switch {
	case errors.Is(err, monc.ErrNotFound):
		return nil
	case err == nil:
		follow.IsCancel = true
		_, err = m.conn.ReplaceOne(ctx, "", filter, follow)
	}
	return err
}

func (m *MongoMapper) CountFollows(ctx context.Context, targetId string, targetType action.TargetType) (int64, error) {
	filter := bson.M{"target_id": targetId, "target_type": targetType}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MongoMapper) GetFollowedUsers(ctx context.Context, targetId string, targetType action.TargetType, opts *basic.PaginationOptions) ([]*Follow, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	follows := make([]*Follow, pageSize)

	filter := bson.M{"target_id": targetId, "target_type": targetType}

	err := m.conn.Find(ctx, &follows, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountFollows(ctx, targetId, targetType)

	if err != nil {
		return nil, 0, err
	}
	return follows, total, err
}

func (m *MongoMapper) GetUserFollowed(ctx context.Context, targetType action.TargetType, userId string, opts *basic.PaginationOptions) ([]*Follow, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	follows := make([]*Follow, pageSize)

	filter := bson.M{"target_type": targetType, "user_id": userId}

	err := m.conn.Find(ctx, &follows, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountFollowsByUserId(ctx, targetType, userId)

	if err != nil {
		return nil, 0, err
	}

	return follows, total, err
}

func (m *MongoMapper) CountFollowsByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error) {
	filter := bson.M{"target_type": targetType, "user_id": userId}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}
