package handler

import "github.com/hamed-lohi/user-management/user"

type Handler struct {
	userStore user.Store
}

func NewHandler(us user.Store) *Handler {
	return &Handler{
		userStore: us,
	}
}
