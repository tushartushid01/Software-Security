package helper

import (
	"Oauth/database"
	"Oauth/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func FetchPasswordAndIDANDRole(userMail, userRole string) (models.UserCredentials, error) {
	SQL := `SELECT users.id,password,role
            FROM   users
            JOIN   role ON users.id=role.user_id
            WHERE  email=$1 
            AND    role  = $2 
            AND    users.archived_at IS NULL 
            `

	var userCredentials models.UserCredentials

	err := database.OauthDB.Get(&userCredentials, SQL, userMail, userRole)
	if err != nil {
		logrus.Printf("FetchPassword: Not able to fetch password, ID or role: %v", err)
		return userCredentials, err
	}
	return userCredentials, nil
}

func Logout(userID string) error {
	SQL := `UPDATE sessions
            SET    end_time=now()
            WHERE  user_id=$1`

	_, err := database.OauthDB.Exec(SQL, userID)
	if err != nil {
		logrus.Printf("Logout: cannot do logout:%v", err)
		return err
	}
	return nil
}

func CreateSession(claims *models.Claims) error {
	SQL := `INSERT INTO sessions(user_id)
            VALUES   ($1)`
	_, err := database.OauthDB.Exec(SQL, claims.ID)
	if err != nil {
		logrus.Printf("CreateSession: cannot create user session:%v", err)
		return err
	}
	return nil
}

func Register(userDetails models.UserDetails) (string, error) {
	// language=SQL
	SQL := `INSERT INTO users(name, email, password) 
                   VALUES ($1, $2, $3)
                   RETURNING id`
	var userID string

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Printf("Register: Not able to hash password:%v", err)
		return userID, err
	}

	err = database.OauthDB.Get(&userID, SQL, userDetails.Name, userDetails.Email, hashPassword)
	if err != nil {
		logrus.Printf("Register: cannot register user:%v", err)
		return userID, err
	}

	SQL = `INSERT INTO role(role, user_id) 
			VALUES ($1, $2)
            `
	_, err = database.OauthDB.Exec(SQL, userDetails.Role, userID)
	if err != nil {
		logrus.Printf("REgister: cannot enter role:%v", err)
		return userID, err
	}
	return userID, nil
}

func UpdatePassword(password []byte, userID string) error {
	SQL := `UPDATE users
			SET    password = $1
			WHERE  id = $2
			AND    archived_at IS NULL 
			`

	_, err := database.OauthDB.Exec(SQL, password, userID)
	if err != nil {
		logrus.Printf("UpdatePassword: cannot update password:%v", err)
		return err
	}
	return nil
}

func CheckSession(userID string) (string, error) {
	SQL := `SELECT id
           FROM    sessions
           WHERE   sessions.end_time IS NULL
           AND     user_id=$1`
	var sessionID string

	err := database.OauthDB.Get(&sessionID, SQL, userID)
	if err != nil {
		logrus.Printf("CheckSession: session expired:%v", err)
		return sessionID, err
	}
	return sessionID, nil
}

func CreateFeedback(feedbackDetails models.FeedbackDetails, userID string) (string, error) {
	SQL := `INSERT INTO feedback(description, created_by)
            VALUES ($1, $2)
            RETURNING id`
	var feedbackID string
	err := database.OauthDB.Get(&feedbackID, SQL, feedbackDetails.Description, userID)
	if err != nil {
		logrus.Printf("CreateProduct: cannot add product: %v", err)
		return feedbackID, err
	}
	return feedbackID, nil
}

func CreateProduct(productDetails models.ProductDetails, userID string) (string, error) {
	SQL := `INSERT INTO product(name, description, created_by, price)
            VALUES ($1, $2, $3, $4)
            RETURNING id`
	var productID string
	err := database.OauthDB.Get(&productID, SQL, productDetails.Name, productDetails.Description, userID, productDetails.Price)
	if err != nil {
		logrus.Printf("CreateProduct: cannot add product: %v", err)
		return productID, err
	}
	return productID, nil
}

func BuyProduct(productID, userID string) error {
	SQL := `INSERT INTO user_product(product_id, user_id)
            VALUES ($1, $2)
            `
	_, err := database.OauthDB.Exec(SQL, productID, userID)
	if err != nil {
		logrus.Printf("BuyProduct: cannot buy product: %v", err)
		return err
	}

	SQL = `UPDATE product
            SET is_bought = $1
            WHERE product.id = $2
            AND archived_at IS NULL 
            `
	_, err = database.OauthDB.Exec(SQL, true, productID)
	if err != nil {
		logrus.Printf("BuyProduct: cannot buy product: %v", err)
		return err
	}
	return nil
}

func GetProducts(userID string) ([]models.ProductOuput, error) {
	//Language sql
	SQL := `SELECT product.id,
       				name,
       				description,
       				price,
       				is_bought,
       				created_by
			FROM   product
			JOIN user_product up on product.id = up.product_id
			WHERE  product.archived_at IS NULL 
			AND    up.archived_at IS NULL 
			AND    product.is_bought != $1  
			AND    user_id = $2
			`
	productDetails := make([]models.ProductOuput, 0)
	err := database.OauthDB.Select(&productDetails, SQL, true, userID)
	if err != nil {
		logrus.Printf("GetProduct: unable to fetch products:%v", err)
		return productDetails, err
	}
	return productDetails, nil
}

func GetAllProducts() ([]models.ProductOuput, error) {
	//Language sql
	SQL := `SELECT product.id,
       				name,
       				description,
       				price,
       				is_bought,
       				created_by
			FROM   product
			WHERE  is_bought != $1
			AND    product.archived_at IS NULL 
			`
	productDetails := make([]models.ProductOuput, 0)
	err := database.OauthDB.Select(&productDetails, SQL, true)
	if err != nil {
		logrus.Printf("GetProduct: unable to fetch products:%v", err)
		return productDetails, err
	}
	return productDetails, nil
}
