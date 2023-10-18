package controllers

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shaikhzidhin/initializer"
)

// GetOTP generates and sends a one-time password (OTP) via email.
func GetOTP(name, email string) string {
	otp, err := getRandNum()
	if err != nil {
		panic(err)
	}
	msg := "Subject: WebPortal OTP\nHey " + name + "Your OTP is " + otp
	sendEmail(name, msg, email)
	return otp
}

// getRandNum generates a random OTP.
func getRandNum() (string, error) {
	otp, err := rand.Int(rand.Reader, big.NewInt(8999))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(otp.Int64()+1000, 10), nil
}

// sendEmail sends an email with the OTP.
func sendEmail(name, msg, email string) {
	SMTPemail := os.Getenv("EMAIL")
	SMTPpass := os.Getenv("PASSWORD")
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(msg))
	if err != nil {
		panic(err)
	}
}

// VerifyOTP verifies if the provided OTP matches the stored OTP in Redis.
func VerifyOTP(superkey, otpInput string, c *gin.Context) bool {
	// OTP verification in Redis
	otp, err := initializer.ReddisClient.Get(context.Background(), superkey).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Error retrieving data from Redis"})
		return false
	}

	if otp == otpInput {
		err := initializer.ReddisClient.Del(context.Background(), superkey).Err()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Error deleting OTP from Redis"})
			return false
		}
		return true
	}
	return false
}
