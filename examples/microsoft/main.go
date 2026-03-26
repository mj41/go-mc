package main

import (
	"context"
	"flag"
	"log"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/microsoft"
)

var (
	address  = flag.String("address", "", "Optional Minecraft server address, for example example.com:25565")
	clientID = flag.String("client-id", microsoft.DefaultClientID, "Microsoft client ID used for device-code auth")
)

func main() {
	flag.Parse()

	access, err := microsoft.AuthenticateDeviceCode(context.Background(), microsoft.Options{
		ClientID: *clientID,
		OnDeviceCode: func(code microsoft.DeviceCode) {
			log.Printf("Open %s and enter code %s", code.VerificationURI, code.UserCode)
			log.Print(code.Message)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authenticated as %s (%s)", access.Profile.Name, access.Profile.ID)

	if *address == "" {
		return
	}

	client := bot.NewClient()
	client.Auth = bot.Auth{
		Name: access.Profile.Name,
		UUID: access.Profile.ID,
		AsTk: access.AccessToken(),
	}

	if err := client.JoinServer(*address); err != nil {
		log.Fatal(err)
	}
	log.Printf("Joined %s as %s", *address, access.Profile.Name)
	_ = client.Close()
}
