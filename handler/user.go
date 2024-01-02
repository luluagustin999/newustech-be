package handler

import (
	"crowdfunding/auth"
	"crowdfunding/helper"
	"crowdfunding/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service 
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {

	// tangkap input dari user melalui website
	// map input dari user ke struct RegisterUserInput
	// struct dia atas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input) //mengubah struct ke json
	if err != nil {
		var errors []string

		errors = helper.FormatValidationError(err)

		errorMessage := gin.H{ "errors" : errors }

		response := helper.APIResponse("Register user gagal di tambahkan", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return //agar eksekusi stop di sini
	}

	 newUser, err := h.userService.RegisterUser(input)
	 
	 if err != nil {
		response := helper.APIResponse("Register user gagal di tambahkan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return //agar eksekusi stop di sini
	 }

	//  token, err := h.jwtService.GenerateToken(user)
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register user gagal di tambahkan", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Register User berhasil ditambahkan", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {

	// user memasukan input (email & password)
	// imnput di tangkap handler
	// mapping dari input user ke input struct
	// input struct di passing ke service
	// di service mencari dengan bantuan repository user dengan email tertentu
	// mencocokan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{ "errors" : errors}

		response := helper.APIResponse("Login Gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.LoginUser(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()} //memanggil error yang ada di service
		response := helper.APIResponse("Login Gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response) 
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login Gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	formatter := user.FormatUser(loggedinUser, token)

	response := helper.APIResponse("Login Berhasil", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability (c *gin.Context) {

	// ada input email dari user
	// input email di mapping ke struct input
	// struct input di passing ke service
	// service akan memanggil repository untuk ngecek apakah email sudah ada atau belum
	//  repository - db
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		response := helper.APIResponse("Check email gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		reponse := helper.APIResponse("Check email gagal", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, reponse)
	}

	data := gin.H {
		"is_email_available" : isEmailAvailable,
	}

	metaMessage := "Email sudah Di Daftarkan"

	if isEmailAvailable { // if isEmailAvailable == true karena return defaultnya adalah false 
		metaMessage = "Email tersedia"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data) 
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// input dari user / website ingat inputan nya bukan json melainkan form body
	// simpan gambarnya di folder "/images"
	// di service kita panggil repo
	// JWT (sementara pakai hardcode, seakan2 user yang login ID = 1 )
	// repo ambil data user yg id = 1
	// repo update data user simpan lokasi

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Gagal upload gambar avatar", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User) 
	userID := currentUser.ID

	// userID := 6 //harusnya dapat dari JWT jadi ini hanya percobaan  

	// old images/namafile.png
	// new images/1-namafile.png
	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename) //result = images/1-namafile.png // 1 di dapat dari ID

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H {
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Gagal upload gambar avatar", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H {
			"is_uploaded" : false,
		}
		response := helper.APIResponse("Gagal upload gambar avatar", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H {
		"is_uploaded" : true,
	}
	response := helper.APIResponse("Gambar avatar berhasil di upload", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "") //token nya di kosongin juga gk papa

	response := helper.APIResponse("Successfully fetch user datta", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}
