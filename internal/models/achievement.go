package models

type AchievementsRequest struct {
	StartDate string `json:"start_date__gte"`
	EndDate   string `json:"end_date__lte"`
}

type AchievementsResponse struct {
	Id          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
	Color       string `json:"color"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
