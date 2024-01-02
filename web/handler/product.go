package handler

import (
	"crowdfunding/product"
	"crowdfunding/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productService product.Service
	userService    user.Service
}

func NewProductHandler(productService product.Service, userService user.Service) *productHandler {
	return &productHandler{productService, userService}
}

func (h *productHandler) Index(c *gin.Context) {
	products, err := h.productService.GetProducts(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"products": products})

}

func (h *productHandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := product.FormCreateProductInput{}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_new.html", input)

}

func (h *productHandler) Create(c *gin.Context) {
	var input product.FormCreateProductInput

	err := c.ShouldBind(&input)
	if err != nil {
		users, e := h.userService.GetAllUsers()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}

		input.Users = users
		input.Error = err

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	user, err := h.userService.GetUserByID(input.UserID) //nyari data user
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.hmtl", nil)
		return
	}

	createProductInput := product.CreateProductInput{}
	createProductInput.Title = input.Title
	createProductInput.Body = input.Body
	createProductInput.User = user

	// karena create campaign inputannya adalah CreateCampaignInput jadi harus di mapping ke FormCreateCampaignInput
	_, err = h.productService.CreateProduct(createProductInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/campaigns")

}

// func (h *campaignHandler) NewImage(c *gin.Context) {
// 	idParam := c.Param("id")
// 	id, _ := strconv.Atoi(idParam)

// 	c.HTML(http.StatusOK, "campaign_image.html", gin.H{"ID" : id})
// }

// func (h *campaignHandler) CreateImage(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}

// 	idParam := c.Param("id")
// 	id, _ := strconv.Atoi(idParam)

// 	existingCampaign, err := h.campaignService.GetCampaignByID(campaign.GetCampaignDetailInput{ID: id})
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}

// 	userID := existingCampaign.UserID

// 	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

// 	err = c.SaveUploadedFile(file, path)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}

// 	createCampaignImageInput := campaign.CreateCampaignImageInput{}
// 	createCampaignImageInput.CampaignID = id
// 	createCampaignImageInput.IsPrimary = true

// 	userCampaign, err := h.userService.GetUserByID(userID)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}

// 	createCampaignImageInput.User = userCampaign

// 	_, err = h.campaignService.SaveCampaignImage(createCampaignImageInput, path)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}

// 	c.Redirect(http.StatusFound, "/campaigns")

// }

func (h *productHandler) Edit(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	existingProduct, err := h.productService.GetProductByID(product.GetProductDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := product.FormUpdateProductInput{}
	input.ID = existingProduct.ID
	input.Title = existingProduct.Title
	input.Body = existingProduct.Body

	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *productHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var input product.FormUpdateProductInput

	err := c.ShouldBind(&input)
	if err != nil {
		input.Error = err
		input.ID = id
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	existingProduct, err := h.productService.GetProductByID(product.GetProductDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	userID := existingProduct.UserID

	userProduct, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	updateInput := product.CreateProductInput{}
	updateInput.Title = input.Title
	updateInput.Body = input.Body
	updateInput.User = userProduct

	_, err = h.productService.UpdateProduct(product.GetProductDetailInput{ID: id}, updateInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/products")

}

func (h *productHandler) Show(c *gin.Context) {

	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	existingProduct, err := h.productService.GetProductByID(product.GetProductDetailInput{ID: id})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	}

	c.HTML(http.StatusOK, "campaign_show.html", existingProduct)

}
