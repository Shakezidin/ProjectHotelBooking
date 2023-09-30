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
	"github.com/shaikhzidhin/initiializer"
)

func GetOTP(name, email string) string {
	otp, err := getRandNum()
	if err != nil {
		panic(err)
	}
	msg := "Subject: WebPortal OTP\nHey " + name + "Your OTP is " + otp
	sendEmail(name, msg, email)
	return otp
}

// Getting a random number for otp. This function helps get otp to generate the a random otp
func getRandNum() (string, error) {
	otp, err := rand.Int(rand.Reader, big.NewInt(8999))
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(otp.Int64()+1000, 10), nil
}

func sendEmail(name, msg, email string) {
	SMTPemail := os.Getenv("EMAIL")
	SMTPpass := os.Getenv("PASSWORD")
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")

	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(msg))
	if err != nil {
		panic(err)
	}
}

func VerifyOTP(superkey, otpInput string, c *gin.Context) bool {
	//otp verification in reddis
	otp, err := initiializer.ReddisClient.Get(context.Background(), superkey).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Error retrieving data from Redis"})
		return false
	} else {
		if otp == otpInput {
			err := initiializer.ReddisClient.Del(context.Background(), superkey).Err()
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Error deleting otp from Redis"})
				return false
			}
			return true
		} else {
			return false
		}
	}
}
