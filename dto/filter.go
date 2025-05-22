package dto
type Filter struct {
	Operation  string `json:"operation" binding:"required,oneof=eq ne gt lt gte lte"`
	Operator string `json:"operator" binding:"required,oneof=created_at type payment_method"`
	Value     string `json:"value" binding:"required"`
}
