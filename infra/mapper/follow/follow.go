package follow

import (
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/meowcloud/action"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Follow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	TargetId   string             `bson:"target_id,omitempty" json:"target_id"`
	TargetType action.TargetType  `bson:"target_type" json:"target_type"`
	UserId     string             `bson:"user_id,omitempty" json:"user_id"`
	IsCancel   bool               `bson:"is_cancel" json:"is_cancel"`
	CreateAt   time.Time          `bson:"create_at,omitempty" json:"create_at,omitempty"`
	UpdateAt   time.Time          `bson:"update_at,omitempty" json:"update_at,omitempty"`
	DeleteAt   time.Time          `bson:"delete_at,omitempty" json:"delete_at,omitempty"`
}
