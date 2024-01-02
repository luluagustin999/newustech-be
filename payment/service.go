package payment

import (
	"crowdfunding/campaign"
	"crowdfunding/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	campaignRepository campaign.Repository
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {

	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-75fgSQZ2S0SPSVXWuzwIraMH"
	midclient.ClientKey = "SB-Mid-client-Q0C0mWmxeKWybmvr"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway {
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(transaction.ID), //karena OrderID string jadi harus di convert
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil  

}

