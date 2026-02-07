package dto

type SubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"omitempty,datetime=01-2006"`
}

type StatsRequest struct {
	UserID      string `form:"user_id" binding:"required,uuid"`
	ServiceName string `form:"service_name"`
	StartDate   string `form:"start_date" binding:"required"`
	EndDate     string `form:"end_date" binding:"required"`
}

type StatsResponse struct {
	Total int `json:"total_price"`
}
