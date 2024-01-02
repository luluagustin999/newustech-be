package campaign

import "crowdfunding/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int 	`json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	User             user.User //untuk munculin data yang login saat ini
}

type CreateCampaignImageInput struct {
	CampaignID int `form:"campaign_id" binding:"required"`
	IsPrimary bool `form:"is_primary"`
	User	user.User
}

type FormCreateCampaignInput struct {
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int 	`form:"goal_amount" binding:"required"`
	Perks            string `form:"perks" binding:"required"`
	UserID			 int 	`form:"user_id" binding:"required"`
	Users 			 []user.User //untuk nampilin semua data user
	Error 			 error
}

type FormUpdateCampaignInput struct {
	ID				 int 	
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int 	`form:"goal_amount" binding:"required"`
	Perks            string `form:"perks" binding:"required"`
	Error 			 error
	User			 user.User
}