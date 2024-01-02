package product

import "crowdfunding/user"

type GetProductDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateProductInput struct {
	Title string    `json:"title" binding:"required"`
	Body  string    `json:"body" binding:"required"`
	User  user.User //untuk munculin data yang login saat ini
}

type CreateProductImageInput struct {
	ProductID int  `form:"product_id" binding:"required"`
	IsPrimary bool `form:"is_primary"`
	User      user.User
}

type FormCreateProductInput struct {
	Title  string      `form:"title" binding:"required"`
	Body   string      `form:"body" binding:"required"`
	UserID int         `form:"user_id" binding:"required"`
	Users  []user.User //untuk nampilin semua data user
	Error  error
}

type FormUpdateProductInput struct {
	ID    int
	Title string `form:"title" binding:"required"`
	Body  string `form:"body" binding:"required"`
	Error error
	User  user.User
}
