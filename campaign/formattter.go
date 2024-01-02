package campaign

import "strings"

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter

}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {

	campaignsFormatter := []CampaignFormatter{} //nilai default null contoh : []

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter

}

type CampaignDetailFormatter struct {
	ID               int      					`json:"id"`
	Name             string   					`json:"name"`
	ShortDescription string   					`json:"short_description"`
	Description      string   					`json:"description"`
	ImageURL         string   					`json:"image_url"`
	GoalAmount       int      					`json:"goal_amount"`
	CurrentAmount    int      					`json:"current_amount"`
	BackerCount		int							`json:"backer_count"`
	UserID           int      					`json:"user_id"`
	Slug             string   					`json:"slug"`
	Perks            []string 					`json:"perks"`
	User 			 CampaignUserFormatter		`json:"user"`
	Images			 []CampaignImageFormatter	`json:"images"`
}

type CampaignUserFormatter struct {
	Name string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImagesURL string `json:"image_url"`
	IsPrimary bool `json:"is_primary"`
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.BackerCount = campaign.BackerCount
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.Slug = campaign.Slug

	// buatt images url

	campaignDetailFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	// buat perks

	var perks []string

	for _, perk := range strings.Split(campaign.Perks, ",") { // mecah kata berdasarkan koma result := "satu", "dua", "tiga"
		perks = append(perks, strings.TrimSpace(perk) ) //menghilangkan spasi contoh " satu" result := "satu"
	}

	campaignDetailFormatter.Perks = perks

	// buat user

	user := campaign.User

	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName
	
	campaignDetailFormatter.User = campaignUserFormatter

	// buat images

	images := []CampaignImageFormatter{} //sekalian bikin default kalo gk ada datanya []

	for _, image := range campaign.CampaignImages {
		campaignImageFormatter := CampaignImageFormatter{}
		campaignImageFormatter.ImagesURL = image.FileName

		isPrimary := false //nilai default false

		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImageFormatter.IsPrimary = isPrimary
		
		images = append(images, campaignImageFormatter)
	
	}

	campaignDetailFormatter.Images = images

	return campaignDetailFormatter

}