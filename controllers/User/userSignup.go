package user

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	controllers "github.com/shaikhzidhin/controllers/Otp"
	Init "github.com/shaikhzidhin/initializer"
	"github.com/shaikhzidhin/models"
)

var validate = validator.New()

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Signup for user signup
func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"Message": "binding error",
		})
		c.Abort()
		return
	}
	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error1"})
		return
	}

	result := Init.DB.Where("user_name = ?", user.UserName).First(&user)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "User_Name already Exist",
		})
		return
	}

	email := Init.DB.Where("email = ?", user.Email).First(&user)
	if email.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Email already exist",
		})
		return
	}

	if user.ReferralCode != "" {
		var referredUser models.User
		result := Init.DB.Where("referral_code = ?", user.ReferralCode).First(&referredUser)
		if result.Error != nil {
			c.JSON(400, gin.H{"error": "user not found in this referral code"})
			return
		}
		user.Wallet.Balance += 50
	}

	phone := Init.DB.Where("phone = ?", user.Phone).First(&user)
	if phone.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "phone number already exist",
		})
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		c.JSON(400, gin.H{
			"msg": "hashing error",
		})
		c.Abort()
		return
	}

	Otp := controllers.GetOTP(user.Name, user.Email)

	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error encoding JSON"})
		return
	}

	// Inserting the OTP into Redis
	err = Init.ReddisClient.Set(context.Background(), "signUpOTP"+user.Email, Otp, 5*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting OTP in Redis client"})
		return
	}

	// Inserting the data into Redis
	err = Init.ReddisClient.Set(context.Background(), "userData"+user.Email, jsonData, 5*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting user data in Redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "messsage": "Go to user/signup-verification"})
}

// SignupVerification for User OTP verification
func SignupVerification(c *gin.Context) {
	type otpCredentials struct {
		Email string `json:"email"`
		Otp   string `json:"otp"`
	}
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err})
		return
	}

	if controllers.VerifyOTP("signUpOTP"+otpCred.Email, otpCred.Otp, c) {
		var userData models.User
		superKey := "userData" + otpCred.Email
		jsonData, err := Init.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting user data from Redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &userData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error binding Redis JSON data to user variable"})
			return
		}

		if userData.ReferralCode == "" {
			userData.ReferralCode = generateRandomString(10)

			// Create user and save transaction history and wallet balance
			results := Init.DB.Create(&userData)
			if results.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "falsee", "Error": results.Error})
				return
			}
			c.JSON(200, gin.H{"status": "user created success"})
			return
		}

		var referredUser models.User

		if err := Init.DB.Where("referral_code = ?", userData.ReferralCode).First(&referredUser).Error; err != nil {
			c.JSON(400, gin.H{"Error": "Error while fetching user"})
			return
		}
		// Update referred user's wallet balance
		var wallet models.Wallet
		if err := Init.DB.Where("user_id = ?", referredUser.ID).First(&wallet).Error; err != nil {
			wallet = models.Wallet{
				Balance: 0,
				UserID:  referredUser.ID,
			}
		} else {
			wallet.Balance += 100
			var transaction models.Transaction

			transaction.Amount = 100
			transaction.UserID = referredUser.ID
			transaction.Details = "Invite bonuse added"
			transaction.Date = time.Now()

			if err := Init.DB.Create(&transaction).Error; err != nil {
				c.JSON(400, gin.H{"Error": "Error while creating transaction"})
				return
			}
		}

		// Save or update the wallet entry
		result := Init.DB.Save(&wallet)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "falsee", "Error": result.Error})
			return
		}

		// Generate a new referral code for the current user
		userData.ReferralCode = generateRandomString(10)

		// Create user and save transaction history and wallet balance
		results := Init.DB.Create(&userData)
		if results.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "falsee", "Error": result.Error})
			return
		}

		var transaction models.Transaction

		transaction.Amount = 50
		transaction.UserID = userData.ID
		transaction.Details = "referal bonuse added"
		transaction.Date = time.Now()

		if err := Init.DB.Create(&transaction).Error; err != nil {
			c.JSON(400, gin.H{"Error": "Error while creating transaction"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Otp Verification success. User creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
	}
}
