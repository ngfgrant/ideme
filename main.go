package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tjarratt/babble"
	"os"
	"strings"
	"time"
)

func main() {
	// Configure DigitalOcean API
	token := os.Getenv("DO_TOKEN")
	ctx := context.TODO()
	api := doApi{ctx, token}

	// Read in config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	sub := strings.ToLower(babble.NewBabbler().Babble()) + viper.GetString("droplet.name")
	// Create Domain Names
	domain := Domain{
		TLD:       viper.GetString("infrastructure.domain.name"),
		Subdomain: sub,
		FullUrl: strings.Join(
			[]string{sub, viper.GetString("infrastructure.domain.name")}, "."),
	}

	// Create Infrastructure
	infra := createInfrastructure(api)

	// Create Application Droplet
	droplet, ip := createDroplet(api, infra, domain)

	// Add Resources to DO Project
	addResourcesToProject(api, infra.Project.ID, []string{droplet.URN()})

	// Add Domain Record
	domainRecord := addDomainRecord(api, viper.GetString("infrastructure.domain.type"), domain, ip)

	// Sleep
	dur := time.Duration(viper.GetInt64("droplet.destroy.minutes"))
	time.Sleep(dur * time.Minute)

	//Tidy up resources
	destroyDroplet(api, droplet)
	destroyDomainRecord(api, domain, domainRecord)

	os.Exit(0)
}
