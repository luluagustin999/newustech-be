package handler

import (
	"crowdfunding/helper"
	"crowdfunding/product"
	"crowdfunding/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
// handler ke service
// service yang menentukan repository mana yang di panggil
// repository : FindAll & FindByUserID
// db

type productHandler struct {
	service product.Service
}

func NewProductHandler(service product.Service) *productHandler {

	return &productHandler{service}

}

// api/v1/campaigns
func (h *productHandler) GetProducts(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) //karena default query adalah string jadi harus coonvert to string

	products, err := h.service.GetProducts(userID)
	if err != nil {
		response := helper.APIResponse("Error to get products", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("List of products", http.StatusOK, "success", product.FormatProducts(products))
	c.JSON(http.StatusOK, response)

}

func (h *productHandler) GetProduct(c *gin.Context) {
	// handler : mapping id yang ada di url ke struct input => service, call formatter
	// service : inputnya struct input => menangkap id di url
	// repository : get campaign by ID

	var input product.GetProductDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	productDetail, err := h.service.GetProductByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse("Product detail", http.StatusOK, "success", product.FormatProductDetail(productDetail))
	c.JSON(http.StatusOK, response)

}

// tangkap parameter dari user ke input struct
// ambil current user dari jwt/handler
// panggil service, parameternya input struct (dan juga buat slug)
// panggil repository untuk simpan data campaign baru

func (h *productHandler) CreateProduct(c *gin.Context) {
	var input product.CreateProductInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed to create product", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User) //ngambil user yang sedang login

	input.User = currentUser

	newProduct, err := h.service.CreateProduct(input)
	if err != nil {
		response := helper.APIResponse("failed to create product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create product", http.StatusOK, "success", product.FormatProduct(newProduct))
	c.JSON(http.StatusOK, response)

}

// user memasukan input
// handler nangkap inputt
// mapping dari input ke input struct (ada 2)
// input dari user, dan juga input yang ada di uri (passing ke service)
// service(find by id, tangkap parameter)
// repository update data campaign

func (h *productHandler) UpdateProduct(c *gin.Context) {

	var inputID product.GetProductDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData product.CreateProductInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update product", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// ngambil user yang saatt ini sedang login agar user lain tidak sembarangan bisa update
	currentUser := c.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updatedProduct, err := h.service.UpdateProduct(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", product.FormatProduct(updatedProduct))
	c.JSON(http.StatusOK, response)

}

// handler
// tangkap input dan ubah ke struct input
// save image campaign ke suatu folder
// service (kondisi manggil point 2 di repo, panggil repo point 1)
// repository :
// 1. create image / save data image ke dalam tabel campaign-images
// 2. ubah is_primary true ke false

// func (h *campaignHandler) UploadImage(c *gin.Context) {
// 	var input campaign.CreateCampaignImageInput

// 	err := c.ShouldBind(&input)

// 	if err != nil {
// 		errors := helper.FormatValidationError(err)
// 		errorMessage := gin.H{"errors":errors}

// 		response := helper.APIResponse("Failed to upload campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
// 		c.JSON(http.StatusUnprocessableEntity, response)
// 		return
// 	}

// 	currentUser := c.MustGet("currentUser").(user.User)
// 	input.User = currentUser
// 	userID := currentUser.ID

// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		data := gin.H{"is_uploaded": false}
// 		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

// 		c.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	path := fmt.Sprint("images/%d-%s", userID, file.Filename)

// 	err = c.SaveUploadedFile(file, path)
// 	if err != nil {
// 		data := gin.H{"is_uploaded": false}
// 		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

// 		c.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	_, err = h.service.SaveCampaignImage(input, path)
// 	if err != nil {
// 		data := gin.H{"is_uploaded": false}
// 		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

// 		c.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	data := gin.H{"is_uploaded": true}
// 	response := helper.APIResponse("Campaign image successfuly uploaded", http.StatusOK, "success", data)

// 	c.JSON(http.StatusOK, response)

// }
