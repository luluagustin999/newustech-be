package product

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Product, error)
	FindByUserID(UserID int) ([]Product, error)
	FindByID(ID int) (Product, error)
	Save(product Product) (Product, error)
	Update(product Product) (Product, error)
	// CreateImage(productImage ProductImages) (ProductImages, error)
	MarkAllImagesAsNonPrimary(productID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Product, error) {
	var products []Product
	
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

 func (r *repository) FindByUserID(UserID int) ([]Product, error) {
	var products []Product
	err := r.db.Where("user_id = ?", UserID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&products).Error

	if err != nil {
		return products, err
	}
	
	return products, nil
 }

 func (r *repository) FindByID(ID int) (Product, error) {
	var product Product

	err := r.db.Where("id = ?", ID).Preload("CampaignImages").Preload("User").Find(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil

 }

 func (r *repository) Save(product Product) (Product, error) {

	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil

 }

 func (r *repository) Update(product Product) (Product, error) {

	err := r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil

 }

//  func (r *repository) CreateImage(campaignImage CampaignImages) (CampaignImages, error) {

// 	err := r.db.Create(&campaignImage).Error
// 	if err != nil {
// 		return campaignImage, err
// 	}

// 	return campaignImage, nil

//  }

//  func (r *repository) MarkAllImagesAsNonPrimary(productID int) (bool, error) {

// 	// UPDATE campaign_imaages SET is_primary = false WHERE campaign_id = 1

// 	err := r.db.Model(&{}).Where("campaign_id = ?", productID).Update("is_primary", false).Error
// 	if err != nil {
// 		return false, err
// 	}

// 	return true, nil

//  }
