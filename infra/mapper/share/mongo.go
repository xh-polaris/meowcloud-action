package share

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

const prefixShareCacheKey = "cache:share"
const CollectionName = "share"

// 用于检查接口是否实现
var _ IMongoMapper = (*MongoMapper)(nil)

type IMongoMapper interface {
	InsertOne(ctx context.Context, targetId string, targetType action.TargetType, userId string) error
	IsShared(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error)
	CountShares(ctx context.Context, targetId string, targetType action.TargetType) (int64, error)
	GetSharedUsers(ctx context.Context, targetId string, targetType action.TargetType, options *basic.PaginationOptions) ([]*Share, int64, error)
	GetUserShared(ctx context.Context, targetType action.TargetType, userId string, options *basic.PaginationOptions) ([]*Share, int64, error)
	CountSharesByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error)
}

type MongoMapper struct {
	conn *monc.Model
}

func NewMongoMapper(config *config.Config) IMongoMapper {
	conn := monc.MustNewModel(config.Mongo.URL, config.Mongo.DB, CollectionName, config.Cache)
	return &MongoMapper{
		conn: conn,
	}
}

func (m *MongoMapper) InsertOne(ctx context.Context, targetId string, targetType action.TargetType, userId string) error {

	newShare := &Share{
		ID:         primitive.NewObjectID(),
		TargetId:   targetId,
		TargetType: targetType,
		UserId:     userId,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	}
	key := prefixShareCacheKey + newShare.TargetId + newShare.ID.Hex()
	_, err := m.conn.InsertOne(ctx, key, newShare)
	if err != nil {
		return err
	}
	return nil

}

func (m *MongoMapper) IsShared(ctx context.Context, targetId string, targetType action.TargetType, userId string) (bool, error) {

	filter := bson.M{"target_id": targetId, "target_type": targetType, "user_id": userId}

	var share Share

	err := m.conn.FindOneNoCache(ctx, &share, filter)
	switch {
	case errors.Is(err, monc.ErrNotFound):
		return false, nil
	default:
		return true, err
	}
}

func (m *MongoMapper) CountShares(ctx context.Context, targetId string, targetType action.TargetType) (int64, error) {
	filter := bson.M{"target_id": targetId, "target_type": targetType}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MongoMapper) GetSharedUsers(ctx context.Context, targetId string, targetType action.TargetType, opts *basic.PaginationOptions) ([]*Share, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	shares := make([]*Share, pageSize)

	filter := bson.M{"target_id": targetId, "target_type": targetType}

	err := m.conn.Find(ctx, &shares, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountShares(ctx, targetId, targetType)

	if err != nil {
		return nil, 0, err
	}
	return shares, total, err
}

func (m *MongoMapper) GetUserShared(ctx context.Context, targetType action.TargetType, userId string, opts *basic.PaginationOptions) ([]*Share, int64, error) {
	pageSize := *opts.Limit
	skip := (*opts.Page - 1) * pageSize

	shares := make([]*Share, pageSize)

	filter := bson.M{"target_type": targetType, "user_id": userId}

	err := m.conn.Find(ctx, &shares, filter, &options.FindOptions{
		Limit: opts.Limit,
		Skip:  &skip,
		// 按时间降序，最新的在最前面
		Sort: bson.M{"create_at": -1},
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := m.CountSharesByUserId(ctx, targetType, userId)

	if err != nil {
		return nil, 0, err
	}

	return shares, total, err
}

func (m *MongoMapper) CountSharesByUserId(ctx context.Context, targetType action.TargetType, userId string) (int64, error) {
	filter := bson.M{"target_type": targetType, "user_id": userId}

	var count int64

	count, err := m.conn.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}
