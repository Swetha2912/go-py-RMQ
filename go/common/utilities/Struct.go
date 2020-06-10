package utilities

type Payload map[string]interface{}

type PaginateResult struct {
	TotalCount  int                      `json:"total"`
	Limit       int                      `json:"limit"`
	CurrentPage int                      `json:"current_page"`
	Docs        []map[string]interface{} `json:"docs"`
}
