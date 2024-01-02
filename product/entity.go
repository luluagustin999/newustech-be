package product

import (
	"crowdfunding/user"
	"time"

	// "github.com/leekchan/accounting"
)

type Product struct {
	ID               int
	UserID           int
	Title            string
	Body			 string
	SlugName    	 string
	CreatedAt        time.Time
	UpdatedAt		 time.Time
	// CampaignImages	[]CampaignImages
	User 			user.User
}

// func (c Campaign) GoalAmountFormatIDR() string {
	
// 	ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}

// 	return ac.FormatMoney(c.GoalAmount)

// }

// func (c Campaign) BackerCountFormatIDR() string {
	
// 	ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}

// 	return ac.FormatMoney(c.BackerCount)

// }

// func (c Campaign) CurrentAmountFormatIDR() string {
	
// 	ac := accounting.Accounting{Symbol: "Rp ", Precision: 2, Thousand: ".", Decimal: ","}

// 	return ac.FormatMoney(c.CurrentAmount)

// }



// type CampaignImages struct {
// 	ID	int
// 	CampaignID int
// 	FileName string
// 	IsPrimary int
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }