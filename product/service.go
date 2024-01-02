package product

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetProducts(userID int) ([]Product, error)
	GetProductByID(input GetProductDetailInput) (Product, error)
	CreateProduct(input CreateProductInput) (Product, error)
	UpdateProduct(inputID GetProductDetailInput, inputData CreateProductInput) (Product, error)
	// SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImages, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Product, error) {

	if userID != 0 { //jika user id bukan 0
		products, err := s.repository.FindByUserID(userID)
		if err != nil {
			return products, err
		}

		return products, nil

	}

	products, err := s.repository.FindAll()
	if err != nil {
		return products, err
	}

	return products, nil

}

func (s *service) GetProductByID(input GetProductDetailInput) (Product, error) {

	product, err := s.repository.FindByID(input.ID)
	if err != nil {
		return product, err
	}

	return product, nil

}

func (s *service) CreateProduct(input CreateProductInput) (Product, error) {

	// mapping inputan user -> input create campaign input -> menjadi objek campaign

	product := Product{}
	product.Title = input.Title
	product.Body = input.Body
	product.UserID = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Title, input.User.ID)
	product.SlugName = slug.Make(slugCandidate) //agar unique result nama campaign 10 => nama-campaign-10 nama-campaign-100

	// proses pembuatan slug

	newProduct, err := s.repository.Save(product)
	if err != nil {
		return newProduct, err
	}

	return newProduct, nil

}

func (s *service) UpdateProduct(inputID GetProductDetailInput, inputData CreateProductInput) (Product, error) {

	product, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return product, err
	}

	// jadi hanya user yang bersangkutan yang bisa update / user yang membuat campaign tersebut
	if product.UserID != inputData.User.ID { //jika user id tidak sama dengan user id yang ada di camapaign maka error
		return product, errors.New("Not an owner of the product")
	}

	product.Title = inputData.Title
	product.Body = inputData.Body

	updatedProduct, err := s.repository.Update(product)
	if err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil

}

// func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImages, error) {

// 	campaign, err := s.repository.FindByID(input.CampaignID)
// 	if err != nil {
// 		return CampaignImages{}, err
// 	}

// 	if campaign.UserID != input.User.ID {
// 		return CampaignImages{}, errors.New("No an owner of the campaign")
// 	}


// 	isPrimary := 0
// 	if input.IsPrimary { // if input.IsPrimary == true 

// 		isPrimary = 1

// 		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID) //ubah menjadi false
// 		if err != nil {
// 			return CampaignImages{}, err
// 		}
// 	}

// 	// mapping campaign image ke struct input
// 	campaignImage := CampaignImages{}
// 	campaignImage.CampaignID = input.CampaignID
// 	campaignImage.IsPrimary = isPrimary //proses mapping input is_primary bool ke integer 
// 	campaignImage.FileName = fileLocation

// 	newCampaignImage, err :=s.repository.CreateImage(campaignImage)
// 	if err != nil {
// 		return newCampaignImage, err
// 	}

// 	return newCampaignImage, nil

// }

