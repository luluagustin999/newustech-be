package product

type ProductFormatter struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	SlugName string `json:"slug_name"`
}

func FormatProduct(product Product) ProductFormatter {
	productFormatter := ProductFormatter{}
	productFormatter.ID = product.ID
	productFormatter.UserID = product.UserID
	productFormatter.Title = product.Title
	productFormatter.Body = product.Body
	productFormatter.SlugName = product.SlugName

	return productFormatter

}

func FormatProducts(products []Product) []ProductFormatter {

	productsFormatter := []ProductFormatter{} //nilai default null contoh : []

	for _, product := range products {
		productFormatter := FormatProduct(product)
		productsFormatter = append(productsFormatter, productFormatter)
	}

	return productsFormatter

}

type ProductDetailFormatter struct {
	ID       int                  `json:"id"`
	Title    string               `json:"title"`
	Body     string               `json:"body"`
	UserID   int                  `json:"user_id"`
	SlugName string               `json:"slug_name"`
	User     ProductUserFormatter `json:"user"`
}

type ProductUserFormatter struct {
	Name string `json:"name"`
	// ImageURL string `json:"image_url"`
}

// type CampaignImageFormatter struct {
// 	ImagesURL string `json:"image_url"`
// 	IsPrimary bool `json:"is_primary"`
// }

func FormatProductDetail(product Product) ProductDetailFormatter {
	productDetailFormatter := ProductDetailFormatter{}
	productDetailFormatter.ID = product.ID
	productDetailFormatter.Title = product.Title
	productDetailFormatter.Body = product.Body
	productDetailFormatter.UserID = product.UserID
	productDetailFormatter.SlugName = product.SlugName

	// buatt images url

	// campaignDetailFormatter.ImageURL = ""

	// if len(campaign.CampaignImages) > 0 {
	// 	campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	// }

	// // buat perks

	// var perks []string

	// for _, perk := range strings.Split(campaign.Perks, ",") { // mecah kata berdasarkan koma result := "satu", "dua", "tiga"
	// 	perks = append(perks, strings.TrimSpace(perk) ) //menghilangkan spasi contoh " satu" result := "satu"
	// }

	// campaignDetailFormatter.Perks = perks

	// buat user

	user := product.User

	productUserFormatter := ProductUserFormatter{}
	productUserFormatter.Name = user.Name
	// campaignUserFormatter.ImageURL = user.AvatarFileName

	productDetailFormatter.User = productUserFormatter

	// buat images

	// 	images := []CampaignImageFormatter{} //sekalian bikin default kalo gk ada datanya []

	// 	for _, image := range campaign.CampaignImages {
	// 		campaignImageFormatter := CampaignImageFormatter{}
	// 		campaignImageFormatter.ImagesURL = image.FileName

	// 		isPrimary := false //nilai default false

	// 		if image.IsPrimary == 1 {
	// 			isPrimary = true
	// 		}
	// 		campaignImageFormatter.IsPrimary = isPrimary

	// 		images = append(images, campaignImageFormatter)

	// 	}

	// 	campaignDetailFormatter.Images = images

	return productDetailFormatter

	// }
}
