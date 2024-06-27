package main

import "github.com/mikerybka/util"

func main() {
	server := &util.Server{
		DataDir: util.EnvVar("DATA_DIR", "/root"),
		TwilioClient: &util.TwilioClient{
			AccountSID:  util.RequireEnvVar("TWILIO_ACCOUNT_SID"),
			AuthToken:   util.RequireEnvVar("TWILIO_AUTH_TOKEN"),
			PhoneNumber: util.RequireEnvVar("TWILIO_PHONE_NUMBER"),
		},
		AdminPhone: util.RequireEnvVar("ADMIN_PHONE"),
		AdminEmail: util.RequireEnvVar("ADMIN_EMAIL"),
		CertDir:    util.RequireEnvVar("CERT_DIR"),
	}
	err := server.Start()
	panic(err)
}
