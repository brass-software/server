package main

import "github.com/mikerybka/util"

func main() {
	twilioClient := &util.TwilioClient{
		AccountSID:  util.RequireEnvVar("TWILIO_ACCOUNT_SID"),
		AuthToken:   util.RequireEnvVar("TWILIO_AUTH_TOKEN"),
		PhoneNumber: util.RequireEnvVar("TWILIO_PHONE_NUMBER"),
	}
	server := util.NewServer(util.EnvVar("DATA_DIR", "/root"), util.RequireEnvVar("ADMIN_PHONE"), twilioClient)
	err := server.Start(util.RequireEnvVar("ADMIN_EMAIL"), util.RequireEnvVar("CERT_DIR"))
	panic(err)
}
