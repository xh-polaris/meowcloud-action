package consts

import (
	"errors"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/basic"
)

var UserNotExist = errors.New("用户不存在")
var RepeatLike = errors.New("请勿重复点赞")
var LikeNotExist = errors.New("点赞不存在")
var RepeatFollow = errors.New("请勿重复关注")
var TryAgain = errors.New("操作失败，请重试")

func CheckUserMeta(meta *basic.UserMeta) error {

	if meta == nil {
		return UserNotExist
	}

	return nil
}

func CheckUserId(userId string) error {

	if userId == "" {
		return UserNotExist
	}

	return nil
}
