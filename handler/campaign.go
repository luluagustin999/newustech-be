package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
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

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {

	return &campaignHandler{service}

}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id")) //karena default query adalah string jadi harus coonvert to string

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// handler : mapping id yang ada di url ke struct input => service, call formatter
	// service : inputnya struct input => menangkap id di url
	// repository : get campaign by ID

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to gett detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

// tangkap parameter dari user ke input struct
// ambil current user dari jwt/handler
// panggil service, parameternya input struct (dan juga buat slug)
// panggil repository untuk simpan data campaign baru

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User) //ngambil user yang sedang login

	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}

// user memasukan input
// handler nangkap inputt
// mapping dari input ke input struct (ada 2)
// input dari user, dan juga input yang ada di uri (passing ke service)
// service(find by id, tangkap parameter)
// repository update data campaign

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {

	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// ngambil user yang saatt ini sedang login agar user lain tidak sembarangan bisa update
	currentUser := c.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
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
