package like

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

const prefixLikeCacheKey = "cache:like"
const CollectionName = "like"

// 用于检查接口是否实现
var _ IMongoMapper = (*MongoMapper)(nil)

type IMongoMapper interface {
	InsertOne(ctx context.Context, targetId string, targetType action.TargetType, userId string) error
	IsLiked(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error)
	CancelLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) error
	CountLikes(ctx context.Context, targetId string, targetType action.TargetType) (int64, error)
	GetLikedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) ([]*Like, int64, error)
	GetUserLiked(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) ([]*Like, int64, error)
	CountLikesByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error)
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

	var like Like

	err := m.conn.FindOne(ctx, targetId+targetType.String()+userId, &like, filter)

	switch {
	// 已经存在则修改isCancel状态
	case err == nil:
		like.IsCancel = false
		_, err = m.conn.ReplaceOneNoCache(ctx, filter, like)
	// 不存在则新建
	case errors.Is(err, monc.ErrNotFound):
		newLike := &Like{
			ID:         primitive.NewObjectID(),
			TargetId:   targetId,
			TargetType: targetType,
			UserId:     userId,
			IsCancel:   false,
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
		}
		key := prefixLikeCacheKey + newLike.TargetId + newLike.ID.Hex()
		_, err = m.conn.InsertOne(ctx, key, newLike)
	}
	return err

}

func (m *MongoMapper) IsLiked(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error) {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var like Like

	err := m.conn.FindOneNoCache(ctx, &like, filter)
	switch {
	case errors.Is(err, monc.ErrNotFound):
		return false, nil
	case err == nil:
		return like.IsCancel, nil
	default:
		return false, err
	}
}

func (m *MongoMapper) CancelLike(ctx context.Context, targetId string, targetType action.TargetType, userId string) error {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var like Like

	err := m.conn.FindOneNoCache(ctx, &like, filter)

	switch {
	case errors.Is(err, monc.ErrNotFound):
		return nil
	case err == nil:
		like.IsCancel = true
		_, err = m.conn.ReplaceOne(ctx, "", filter, like)
	}
	return err
}

func (m *MongoMapper) CountLikes(ctx context.Context, targetId string, targetType action.TargetType) (int64, error) {
	filter := bson.M{"target_id": targetId, "target_type": targetType}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MongoMapper) GetLikedUsers(ctx context.Context, targetId string, targetType action.TargetType, opts *basic.PaginationOptions) ([]*Like, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	likes := make([]*Like, pageSize)

	filter := bson.M{"target_id": targetId, "target_type": targetType}

	err := m.conn.Find(ctx, &likes, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountLikes(ctx, targetId, targetType)

	if err != nil {
		return nil, 0, err
	}
	return likes, total, err
}

func (m *MongoMapper) GetUserLiked(ctx context.Context, targetType action.TargetType, userId string, opts *basic.PaginationOptions) ([]*Like, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	likes := make([]*Like, pageSize)

	filter := bson.M{"target_type": targetType, "user_id": userId}

	err := m.conn.Find(ctx, &likes, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountLikesByUserId(ctx, targetType, userId)

	if err != nil {
		return nil, 0, err
	}

	return likes, total, err
}

func (m *MongoMapper) CountLikesByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error) {
	filter := bson.M{"target_type": targetType, "user_id": userId}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}
