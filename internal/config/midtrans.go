package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
)

func InitMidtrans() {
	env := os.Getenv("MIDTRANS_ENV")

	if env == "production" {
		midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY_PROD")
		midtrans.Environment = midtrans.Production
	} else {
		// default sandbox
		midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY_SANDBOX")
		midtrans.Environment = midtrans.Sandbox
	}
}
