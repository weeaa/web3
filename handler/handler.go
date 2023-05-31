package handler

type Handler struct {
	M     map[string]interface{}
	MCopy map[string]interface{}
}

func New() *Handler {
	return &Handler{
		M:     make(map[string]interface{}),
		MCopy: make(map[string]interface{}),
	}
}
