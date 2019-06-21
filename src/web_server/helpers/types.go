package helpers

// 分页结构
type Page struct {
	Total 			int64				`json:"total"`
	TotalPage 		float64				`json:"total_page"`
	PageSize 		int					`json:"page_size"`
	Page 			int					`json:"page"`
	List 			interface{}			`json:"list"`
}
