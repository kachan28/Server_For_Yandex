package manticoreModels

type ManticoreSearchResp struct {
	Hits     HitResponse `json:"hits"`
	TimedOut bool        `json:"timed_out"`
	Took     int         `json:"took"`
}

type HitResponse struct {
	Hits  []HitContent `json:"hits"`
	Total int          `json:"total"`
}

type HitContent struct {
	ID     string        `json:"_id"`
	Score  int           `json:"_score"`
	Source SearchProduct `json:"_source"`
}

type SearchProduct struct {
	Articul string `json:"articul"`
	Shop    string `json:"shop"`
	Title   string `json:"title"`
}
