package handler

import (
	"Oauth/database/helper"
	"Oauth/models"
	"Oauth/utilities"
	"database/sql"
	"github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte("secret_key")

var key = []byte("mysecretkeyaaaaa")

//func LambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	return chiLambda.Proxy(req)
//}

func Login(w http.ResponseWriter, r *http.Request) {
	var userDetails models.UsersLoginDetails
	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	userDetails.Email = strings.ToLower(userDetails.Email)

	userCredentials, fetchErr := helper.FetchPasswordAndIDANDRole(userDetails.Email, userDetails.Role)

	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			userOutboundData := make(map[string]interface{})
			userOutboundData["error"] = "INVALID EMAIL"

			err := utilities.Encoder(w, userOutboundData)
			if err != nil {
				return
			}

			logrus.Printf("FetchPasswordAndId: not able to get password, id, or role:%v", fetchErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if PasswordErr := bcrypt.CompareHashAndPassword([]byte(userCredentials.Password), []byte(userDetails.Password)); PasswordErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Printf("password misMatch")
		_, err := w.Write([]byte("ERROR: Wrong password"))
		if err != nil {
			return
		}
		return
	}

	expiresAt := time.Now().Add(60 * time.Minute)

	claims := &models.Claims{
		ID:   userCredentials.ID,
		Role: userCredentials.Role,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expiresAt.Unix(),
			// Issuer:    userCredentials.Role,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("TokenString: cannot create token string:%v", err)
		return
	}

	err = helper.CreateSession(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateSession: cannot create session:%v", err)
		return
	}

	userOutboundData := make(map[string]interface{})

	userOutboundData["token"] = tokenString

	err = utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("Login: Not able to login:%v", err)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("AddAddress:QueryParam for ID:%v", ok)
		return
	}

	err := helper.Logout(contextValues.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("Logout:unable to logout:%v", err)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userDetails models.UserDetails

	decoderErr := utilities.Decoder(r, &userDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Register: Decoder error:%v", decoderErr)
		return
	}

	userID, err := helper.Register(userDetails)
	if err != nil {
		log.Fatalf("Register:error creating user at firebase: %v\n", err)
		return
	}

	userOutboundData := make(map[string]string)

	userOutboundData["id:"] = userID

	err = utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("Register: encoding error:%v", err)
		return
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	passwordDetails := models.PasswordDetails{}

	decoderErr := utilities.Decoder(r, &passwordDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("UpdatePassword: Decoder error:%v", decoderErr)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("UpdatePassword: not able to hash password: %v", err)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("UpdatePassword: not able to change password: %v", ok)
		return
	}

	updateErr := helper.UpdatePassword(hashedPassword, contextValues.ID)
	if updateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateProduct: not able to get context value: %v", err)
		return
	}

	message := "Password Updated Successfully"

	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("CreateProduct:%v", err)
		return
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productDetails models.ProductDetails

	decoderErr := utilities.Decoder(r, &productDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateProduct: not able to get context value: %v", ok)
		return
	}

	productID, err := helper.CreateProduct(productDetails, contextValues.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateProduct:cannot create note:%v", err)
		return
	}

	userOutboundData := make(map[string]string)

	userOutboundData["productID:"] = productID

	err = utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("CreateProduct:%v", err)
		return
	}
}

func BuyProduct(w http.ResponseWriter, r *http.Request) {
	var buyDetails models.BuyDetails

	decoderErr := utilities.Decoder(r, &buyDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("BuyProduct: not able to get context value: %v", ok)
		return
	}

	err := helper.BuyProduct(buyDetails.ProductID, contextValues.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("BuyProduct:cannot create note:%v", err)
		return
	}

	message := "Bought this item successfully"

	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("CreateProduct:%v", err)
		return
	}
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var boughtVariable bool
	contentType := r.Header.Get("Bought-products")
	if contentType == "True" || contentType == "true" || contentType == "TRUE" {
		boughtVariable = true
	} else {
		boughtVariable = false
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("GetNotes: not able to get context value: %v", ok)
		return
	}

	var productDetails []models.ProductOuput
	var productDetailsErr error
	if boughtVariable {
		productDetails, productDetailsErr = helper.GetProducts(contextValues.ID)
		if productDetailsErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("GetProducts: not able to get Products: %v", productDetailsErr)
			return
		}
	} else {
		productDetails, productDetailsErr = helper.GetAllProducts()
		if productDetailsErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.Printf("GetProducts: not able to get Products: %v", productDetailsErr)
			return
		}
	}

	err := utilities.Encoder(w, productDetails)
	if err != nil {
		logrus.Printf("GetNotes: %v", err)
		return
	}

	// wrapperNotesDetails := make([]models.ProductOuput, 0)

	// Unmarshal JSON into the struct
	//for i, _ := range noteDetails {
	//	decrypted, err := Decrypt(noteDetails[i].Encrypt, key)
	//	println(decrypted)
	//	if err != nil {
	//		fmt.Println("Decryption error:", err)
	//		return
	//	}
	//
	//	var OutputNoteDetails models.NoteDetails
	//	err = json.Unmarshal([]byte(decrypted), &OutputNoteDetails)
	//	if err != nil {
	//		fmt.Println("Error:", err)
	//		return
	//	}
	//	wrapperNotesDetailsOut := models.WrapperNoteDetails{
	//		Id:          noteDetails[i].Id,
	//		NoteDetails: OutputNoteDetails,
	//	}
	//	wrapperNotesDetails = append(wrapperNotesDetails, wrapperNotesDetailsOut)
	//}
}

func CreateFeedback(w http.ResponseWriter, r *http.Request) {
	var feedbackDetails models.FeedbackDetails

	decoderErr := utilities.Decoder(r, &feedbackDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateProduct: not able to get context value: %v", ok)
		return
	}

	feedbackID, err := helper.CreateFeedback(feedbackDetails, contextValues.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateProduct:cannot create FeedBack:%v", err)
		return
	}

	userOutboundData := make(map[string]string)

	userOutboundData["feedbackID is submitted successfully:"] = feedbackID

	err = utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("CreateProduct:%v", err)
		return
	}
}

//func CreateNotes(w http.ResponseWriter, r *http.Request) {
//	var noteDetails models.NoteDetails
//
//	decoderErr := utilities.Decoder(r, &noteDetails)
//
//	if decoderErr != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("Decoder error:%v", decoderErr)
//		return
//	}
//
//	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("GetNotes: not able to get context value: %v", ok)
//		return
//	}
//
//	noteDetailsBytes, err := json.MarshalIndent(noteDetails, "", "")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	encrypted, err := Encrypt(string(noteDetailsBytes), key)
//	if err != nil {
//		fmt.Println("Encryption error:", err)
//		return
//	}
//	var noteID string
//
//	// transaction started
//	txErr := database.Tx(func(tx *sqlx.Tx) error {
//		noteID, err = helper.CreateNotes(encrypted, tx)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			logrus.Printf("CreateNotes:cannot create note:%v", err)
//			return err
//		}
//
//		err = helper.AddUserNote(noteID, contextValues.ID, tx)
//		return err
//	})
//	if txErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("CreateNotes:cannot create note:%v", err)
//		return
//	}
//
//	userOutboundData := make(map[string]string)
//
//	userOutboundData["noteID:"] = noteID
//
//	err = utilities.Encoder(w, userOutboundData)
//	if err != nil {
//		logrus.Printf("AddCategory:%v", err)
//		return
//	}
//}

//func GetNotes(w http.ResponseWriter, r *http.Request) {
//	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("GetNotes: not able to get context value: %v", ok)
//		return
//	}
//
//	noteDetails, noteDetailsErr := helper.GetNotes(contextValues.ID)
//	if noteDetailsErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("GetNotes: not able to get notes: %v", noteDetailsErr)
//		return
//	}
//
//	wrapperNotesDetails := make([]models.WrapperNoteDetails, 0)
//
//	// Unmarshal JSON into the struct
//	for i, _ := range noteDetails {
//		decrypted, err := Decrypt(noteDetails[i].Encrypt, key)
//		println(decrypted)
//		if err != nil {
//			fmt.Println("Decryption error:", err)
//			return
//		}
//
//		var OutputNoteDetails models.NoteDetails
//		err = json.Unmarshal([]byte(decrypted), &OutputNoteDetails)
//		if err != nil {
//			fmt.Println("Error:", err)
//			return
//		}
//		wrapperNotesDetailsOut := models.WrapperNoteDetails{
//			Id:          noteDetails[i].Id,
//			NoteDetails: OutputNoteDetails,
//		}
//		wrapperNotesDetails = append(wrapperNotesDetails, wrapperNotesDetailsOut)
//	}
//
//	err := utilities.Encoder(w, wrapperNotesDetails)
//	if err != nil {
//		logrus.Printf("GetNotes: %v", err)
//		return
//	}
//}
//
//func Encrypt(text string, key []byte) (string, error) {
//	plaintext := []byte(text)
//
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return "", err
//	}
//
//	// GCM mode provides authenticated encryption
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return "", err
//	}
//
//	// Create a random nonce
//	nonce := make([]byte, gcm.NonceSize())
//	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
//		return "", err
//	}
//
//	// Encrypt the plaintext
//	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
//
//	// Combine nonce and ciphertext for storage
//	result := append(nonce, ciphertext...)
//
//	// Encode the result in base64 for easy storage and transmission
//	return base64.StdEncoding.EncodeToString(result), nil
//}
//
//func Decrypt(ciphertext string, key []byte) (string, error) {
//	// Decode the base64 encoded string
//	data, err := base64.StdEncoding.DecodeString(ciphertext)
//	if err != nil {
//		return "", err
//	}
//
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return "", err
//	}
//
//	gcm, err := cipher.NewGCM(block)
//	if err != nil {
//		return "", err
//	}
//
//	// Extract nonce from the data
//	nonceSize := gcm.NonceSize()
//	nonce := data[:nonceSize]
//	ciphertextBytes := data[nonceSize:]
//
//	// Decrypt the ciphertext
//	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
//	if err != nil {
//		return "", err
//	}
//
//	return string(plaintext), nil
//}
