package handler

import (
	"github.com/weeaa/nft/pkg/safemap"
)

type Handler struct {
	M     *safemap.SafeMap[string, interface{}]
	MCopy *safemap.SafeMap[string, interface{}]
}

// New returns a Handler. It is used to store data 🧸.
func New() *Handler {
	return &Handler{
		M:     safemap.New[string, interface{}](),
		MCopy: safemap.New[string, interface{}](),
	}
}

// Copy is a shorter func for ForEach.
func (h *Handler) Copy() {
	h.M.ForEach(func(k string, v interface{}) {
		h.MCopy.Set(k, v)
	})
}
