package dto

// FishResponse は魚レスポンスDTOです。
type FishResponse struct {
	ID      string `json:"id"`
	NameJa  string `json:"nameJa"`
	Name    string `json:"name"`
	Explain string `json:"explain"`
}
