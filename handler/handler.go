package handler

import (
	"github.com/weeaa/nft/pkg/safemap"
)

type Handler struct {
	M     *safemap.SafeMap[string, interface{}]
	MCopy *safemap.SafeMap[string, interface{}]
}

func New() *Handler {
	return &Handler{
		M:     safemap.New[string, interface{}](),
		MCopy: safemap.New[string, interface{}](),
	}
}
