package main

import (
	"net/http"
	"os"

	"github.com/mikerybka/util"
)

func main() {
	email := os.Getenv("EMAIL")
	if email == "" {
		panic("missing EMAIL")
	}
	certDir := os.Getenv("CERT_DIR")
	if certDir == "" {
		panic("missing CERT_DIR")
	}
	twilioAccountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	if twilioAccountSID == "" {
		panic("missing TWILIO_ACCOUNT_SID")
	}
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	if twilioAuthToken == "" {
		panic("missing TWILIO_AUTH_TOKEN")
	}
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	if twilioPhoneNumber == "" {
		panic("missing TWILIO_PHONE_NUMBER")
	}
	authDir := os.Getenv("AUTH_DIR")
	if authDir == "" {
		authDir = "data/auth"
	}
	dataDir := os.Getenv("SCHEMA_CAFE_DIR")
	if dataDir == "" {
		dataDir = "data/schema.cafe"
	}

	s := &util.MultiHostServer{
		Hosts: map[string]http.Handler{
			"api.schema.cafe": &util.MultiUserApp{
				Twilio: &util.TwilioClient{
					AccountSID:  twilioAccountSID,
					AuthToken:   twilioAuthToken,
					PhoneNumber: twilioPhoneNumber,
				},
				AuthFiles: &util.LocalFileSystem{
					Root: authDir,
				},
				App: &util.WebAPI{
					Types: map[string]util.Type{
						"Schema": {
							IsStruct: true,
							Fields: []util.Field{
								{
									ID:   "id",
									Name: "ID",
									Type: "string",
								},
								{
									ID:   "name",
									Name: "Name",
									Type: "string",
								},
								{
									ID:   "fields",
									Name: "Fields",
									Type: "[]Field",
								},
							},
						},
						"Field": {
							IsStruct: true,
							Fields: []util.Field{
								{
									ID:   "id",
									Name: "ID",
									Type: "string",
								},
								{
									ID:   "name",
									Name: "Name",
									Type: "string",
								},
								{
									ID:   "type",
									Name: "Type",
									Type: "string",
								},
							},
						},
					},
					RootType: "map[string]Schema",
					Data: &util.LocalFileSystem{
						Root: dataDir,
					},
				},
			},
			"mikerybka.dev": &util.PingServer{},
		},
	}
	panic(s.Start(email, certDir))
}
