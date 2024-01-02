package main

import (
	"crowdfunding/auth"
	"crowdfunding/campaign"
	"crowdfunding/handler"
	"crowdfunding/helper"
	"crowdfunding/payment"
	"crowdfunding/product"
	"crowdfunding/transaction"
	"crowdfunding/user"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "crowdfunding/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := "host=localhost user=postgres password=agustin999  dbname=newus_tech port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to Database Successful")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	productRepository := product.NewRepository(db)
	transacionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	productService := product.NewService(productRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transacionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	productHandler := handler.NewProductHandler(productService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	userWebHandler := webHandler.NewUserHandler(userService)
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService)
	productWebHandler := webHandler.NewProductHandler(productService, userService)
	transactionWebHandler := webHandler.NewTransactionHandler(transactionService)
	sessionWebHandler := webHandler.NewSessionHandler(userService)

	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := cookie.NewStore([]byte(auth.SECRETKEY))
	router.Use(sessions.Sessions("crowdfunding", cookieStore))

	router.HTMLRender = loadTemplates("./web/templates")

	router.Static("/images", "./images") //untuk url routing
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")

	api := router.Group("api/v1")

	api.POST("/user", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign) //ngmbil user yg lg login
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	// api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)

	api.GET("/products", productHandler.GetProducts)
	api.GET("/products/:id", productHandler.GetProduct)
	api.POST("/products", authMiddleware(authService, userService), productHandler.CreateProduct) //ngmbil user yg lg login
	api.PUT("/products/:id", authMiddleware(authService, userService), productHandler.UpdateProduct)
	// api.POST("/campaign-images", authMiddleware(authService, userService), productHandler.UploadImage)


	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUsertransactions)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	router.GET("/users/new", userWebHandler.New)
	router.POST("/users", userWebHandler.Create)
	router.GET("/users/edit/:id", userWebHandler.Edit)
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.Update)
	router.GET("/users/avatar/:id", authAdminMiddleware(), userWebHandler.NewAvatar)
	router.POST("/users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)

	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHandler.New)
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.Create)
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.NewImage)
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.CreateImage)
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHandler.Edit)
	router.POST("/campaigns/update/:id", authAdminMiddleware(), campaignWebHandler.Update)
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHandler.Show)

	router.GET("/products", authAdminMiddleware(), productWebHandler.Index)
	router.GET("/products/new", authAdminMiddleware(), productWebHandler.New)
	router.POST("/products", authAdminMiddleware(), productWebHandler.Create)
	// router.GET("/products/image/:id", authAdminMiddleware(), productWebHandler.NewImage)
	// router.POST("/products/image/:id", authAdminMiddleware(), productWebHandler.CreateImage)
	router.GET("/products/edit/:id", authAdminMiddleware(), productWebHandler.Edit)
	router.POST("/products/update/:id", authAdminMiddleware(), productWebHandler.Update)
	router.GET("/products/show/:id", authAdminMiddleware(), productWebHandler.Show)

	router.GET("/transactions", authAdminMiddleware(), transactionWebHandler.Index)

	router.GET("/login", sessionWebHandler.New)
	router.POST("/session", sessionWebHandler.Create)
	router.GET("/logout", sessionWebHandler.Destroy)

	router.Run()

}

// karena ingin function validation token & get user by id maka harus begini bentuk functionnya
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") { //apakah di authheader ada kata bearer
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// default Bearer tookentokentoken karena kita ingin ambil token jadi harus di splitt
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1] //[Bearer, tokentokentoken]
		}

		// validasi token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims) //ubah token jwt ke map jw mapclains supaya bisa ngambil user id

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64)) //claim has format map then convert to float 64 and then convert to integer

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}
}

func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		// userID adalah key jadi kalau userid tidak ada berarti tidak login
		userIDSession := session.Get("userID")
		if userIDSession == nil { // dalam kondisi tidak login
			c.Redirect(http.StatusFound, "/login")
			return
		}
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}

	return r

}

// =========================
// 	TEST CREATE TRANSACTIONS
// 	=========================

// 	user, _ := userService.GetUserByID(2)

// 	input := transaction.CreateTransactionInput{
// 		CampaignID: 2,
// 		Amount: 50000,
// 		User: user,
// 	}

// 	transactionService.CreateTransaction(input)

// 	=========================
// 	TEST CREATE CAMPAIGNS
// 	=========================

// 	input := campaign.CreateCampaignInput{}
// 	input.Name = "Penggalangan dana start up"
// 	input.ShortDescription = "short description"
// 	input.Description = "testtttttttttttttttt"
// 	input.GoalAmount = 100000
// 	input.Perks = "hadiah satu, dua, tiga"

// 	inputUser, _ := userService.GetUserByID(1)

// 	input.User = inputUser

// 	_, err = campaignService.CreateCampaign(input)
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	=========================
// 	TEST FIND CAMPAIGNS
// 	=========================

// 	campaignService := campaign.NewService(campainRepository)

// 	campaign, err := campaignService.FindCampaigns(2)

// 	=========================
// 	TEST FIND BY USER ID CAMPAIGN RELASI TO CAMPAIGN IMAGES
// 	=========================

// 	campainRepository :=campaign.NewRepository(db)

// 	campaigns, err := campainRepository.FindByUserID(1)

// 	fmt.Println(len(campaigns))

// 	for _, campaignsss := range campaigns {
// 		fmt.Println(campaignsss.Name)

// 		if len(campaignsss.CampaignImages) > 0 {
// 			fmt.Println("Jumlah gambar", (len(campaignsss.CampaignImages)))
// 			fmt.Println(campaignsss.CampaignImages[0].FileName) //data yang di ambil cuman satu aja
// 		}

// 	}

// 	=========================
// 	TEST FIND ALL REPOSITORY
// 	=========================

// 	campaigns, err := campainRepository.FindAll()

// 		fmt.Println(len(campaigns))

// 		for _, campaignsss := range campaigns {
// 			fmt.Println(campaignsss.Name)
// 		}

// 	LANGKAH LANGKAH MIDDLEWARE MENGGUNAKAN JWT

// 	ambil nilai header Authorization: Bearer tokentoken/isi dari generate token
// 	dari header authorization, ambil nilai dari tokennya saja
// 	validasi token
// 	ambil user_id
// 	ambil user dari db berdasarkan user_id lewat service
// 	kita set context isinya user

// 	=============================
// 	TEST JWT TOKEN VALIDATIOn
// 	=============================

// 	token,err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.zMk3KISsMiv_cBrc5H2oxyT0JXeGJUwPm4VDY0C-yXc")

// 	if err != nil {
// 		fmt.Println("Error JWT Valdiation")
// 	}

// 	if token.Valid {
// 		fmt.Println("Successsssssssssss valid")
// 	} else {
// 		fmt.Println("Invalidddddddddddddddd")
// 	}

// 	=============================
// 	TEST JWT TOKEN
// 	=============================

// 	fmt.Println(authService.GenerateToken(1001))

// 	=============================
// 	TEST UPLOAD AVATAR IN SERVICE
// 	=============================

// 	userService.SaveAvatar(6, "images/1-profile.png")

// 	=================================================
// 	CEK EMAIL TERSEDIA ATAU TIDAK MENGGUNAKAN SERVICE
// 	=================================================

// 	input := user.CheckEmailInput {
// 		Email: "pesulapmerah123@gmail.com",
// 	}

// 	bool, err := userService.IsEmailAvailable(input)
// 	if err != nil {
// 		fmt.Println("Gagal")
// 	}

// 	fmt.Println(bool)

// 	====================================
// 	TEST NYARI EMAIL MENGGUNAKAN SERVICE
// 	====================================

// 	input := user.LoginInput {
// 		Email: "yudistira@gmail.com",
// 		Password: "yudistirar626",
// 	}

// 	user, err := userService.LoginUser(input)

// 	if err != nil {
// 		fmt.Println("Gagal Login")
// 		fmt.Println(err.Error())
// 	return
// 	}

// 	fmt.Println(user.Email)
// 	fmt.Println(user.Name)

// 	=========================================
// 	TEST FIND BY EMAIL MENGGUNAKAN REPOSITORY
// 	=========================================

// 	userByEmail, err := userRepository.FindByEmail("samsudin@gmail.com")

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	fmt.Println(userByEmail.Name)

// 	====================================
// 	TEST CREATE USER MENGGUNAKAN SERVICE
// 	====================================

// 	userInput := user.RegisterUserInput{}
// 	userInput.Name = "Pesulap merah"
// 	userInput.Occupation = "Pesulap"
// 	userInput.Email = "pesulapmerah@gmail.com"
// 	userInput.Password = "12345"

// 	userService.RegisterUser(userInput)

// 	=======================================
// 	TEST CREATE USER MENGGUNAKAN REPOSITORY
// 	=======================================

// 	user := user.User {
// 		Name : "Gus Samsudin",
// 		Occupation: "Padepokna Nur Dzat",
// 		Email: "samsudin@gmail.com",

// 	}

// 	userRepository.Save(user)
