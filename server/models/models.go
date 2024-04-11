package models

import "github.com/form3tech-oss/jwt-go"

//const (
//	UserRoleAdmin UserProfile = "admin"
//	UserRoleUser  UserProfile = "user"
//)

type ContextValues struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type UserCredentials struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type PasswordDetails struct {
	Password string `json:"password"`
}

type UsersLoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type UserDetails struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ProductDetails struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price"`
}

type BuyDetails struct {
	ProductID string `json:"productId" db:"product_id"`
}

type GetNoteDetails struct {
	Id      string `json:"id" db:"id"`
	Encrypt string `json:"encrypt" db:"encrypt"`
}

type ProductOuput struct {
	Id          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	IsBought    bool   `json:"isBought" db:"is_bought"`
	CreatedBy   string `json:"createdBy" db:"created_by"`
}

type FeedbackDetails struct {
	Description string `json:"description" db:"description"`
}
