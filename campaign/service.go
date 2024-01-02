package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImages, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {

	if userID != 0 { //jika user id bukan 0
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, nil

	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil

}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {

	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {

	// mapping inputan user -> input create campaign input -> menjadi objek campaign

	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate) //agar unique result nama campaign 10 => nama-campaign-10 nama-campaign-100

	// proses pembuatan slug

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil

}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {

	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	// jadi hanya user yang bersangkutan yang bisa update / user yang membuat campaign tersebut
	if campaign.UserID != inputData.User.ID { //jika user id tidak sama dengan user id yang ada di camapaign maka error
		return campaign, errors.New("Not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil

}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImages, error) {

	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImages{}, err
	}

	if campaign.UserID != input.User.ID {
		return CampaignImages{}, errors.New("No an owner of the campaign")
	}


	isPrimary := 0
	if input.IsPrimary { // if input.IsPrimary == true 

		isPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID) //ubah menjadi false
		if err != nil {
			return CampaignImages{}, err
		}
	}

	// mapping campaign image ke struct input
	campaignImage := CampaignImages{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary //proses mapping input is_primary bool ke integer 
	campaignImage.FileName = fileLocation

	newCampaignImage, err :=s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil

}

