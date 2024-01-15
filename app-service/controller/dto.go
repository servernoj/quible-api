package controller

type GetGamesDTO struct {
	Date string `form:"date" binding:"required"`
}
