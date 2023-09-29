package models

type (
	TraderAddBody struct {
		BaseAddress string `json:"baseAddress"`
		Importance  string `json:"importance"`
	}

	TraderRemoveBody struct {
		BaseAddress string `json:"baseAddress"`
	}
)
