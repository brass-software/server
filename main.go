package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

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

	go func() {
		time.Sleep(time.Hour)
		update()
	}()

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
				App: &util.SchemaCafe{
					Data: &util.LocalFileSystem{
						Root: dataDir,
					},
				},
			},
			"schema.cafe": &util.MultiUserApp{
				Twilio: &util.TwilioClient{
					AccountSID:  twilioAccountSID,
					AuthToken:   twilioAuthToken,
					PhoneNumber: twilioPhoneNumber,
				},
				AuthFiles: &util.LocalFileSystem{
					Root: authDir,
				},
				App: &util.SchemaCafe{
					Data: &util.LocalFileSystem{
						Root: dataDir,
					},
				},
			},
		},
	}
	panic(s.Start(email, certDir))
}

func update() {
	script := "git pull && go get -u && git add --all && git commit -m update && git push && go build -o /usr/local/bin/server . && systemctl restart server"
	cmd := exec.Command("bash", "-c", script)
	cmd.Dir = "/root/server"
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
	}
}
