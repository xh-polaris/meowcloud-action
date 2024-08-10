package controller

type ActionController struct {
	IFollowController
	ILikeController
	IShareController
}

func NewActionController() *ActionController {
	return &ActionController{
		IFollowController: NewFollowController(),
		ILikeController:   NewLikeController(),
		IShareController:  NewShareController(),
	}
}
