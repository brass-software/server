package main

import (
	_ "embed"
	"net/http"
	"os"

	"github.com/mikerybka/util"
)

//go:embed assets/cafe.ico
var cafeIcon []byte

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
	adminPhone := os.Getenv("ADMIN_PHONE")
	if adminPhone == "" {
		panic("missing ADMIN_PHONE")
	}

	twilioClient := &util.TwilioClient{
		AccountSID:  twilioAccountSID,
		AuthToken:   twilioAuthToken,
		PhoneNumber: twilioPhoneNumber,
	}

	s := &util.MultiHostServer{
		TwilioClient: twilioClient,
		AdminPhone:   adminPhone,
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
						"string": {
							IsScalar: true,
							Kind:     "string",
						},
						"int": {
							IsScalar: true,
							Kind:     "int",
						},
						"bool": {
							IsScalar: true,
							Kind:     "bool",
						},
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
									Type: "FieldList",
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
						"FieldList": {
							IsArray:  true,
							ElemType: "Field",
						},
						"SchemaList": {
							IsMap:    true,
							ElemType: "Schema",
						},
						"Org": {
							IsStruct: true,
							Fields: []util.Field{
								{
									ID:   "schemas",
									Name: "Schemas",
									Type: "SchemaList",
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
			"schema.cafe": http.RedirectHandler("https://www.schema.cafe", http.StatusMovedPermanently),
			"www.schema.cafe": &util.WebApp[*util.Schema]{
				Name:        "Schema Cafe",
				Description: "Schema database",
				Author:      "Mike Rybka",
				Keywords:    []string{"software", "developer", "tools"},
				Favicon:     cafeIcon,
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
				TwilioClient: twilioClient,
				Files: &util.LocalFileSystem{
					Root: authDir,
				},
			},
			"mikerybka.dev": &util.PingServer{},
		},
	}
	panic(s.Start(email, certDir))
}
