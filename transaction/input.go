package transaction

import "crowdfunding/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User //agar bisa mengambil user yang sedang login atau jwt
}

type CreateTransactionInput struct {
	Amount 	   int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User user.User
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID 		  string `json:"order_id"`
	PaymentType 	  string `json:"payment_type"`
	FraudStatus 	  string `jsoon:"fraud_status"`
}