package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

type HttpCredentials struct {
	USERNAME string
	PASSWORD string
}

func createDroplet(api doApi, infra doInfrastructure, domain Domain) (*godo.Droplet, string) {
	client := godo.NewFromToken(api.token)

	// HTTP Login Details
	HTTP_USERNAME := strings.ReplaceAll(domain.Subdomain, ".", "-")
	HTTP_PASSWORD := generateHttpPassword()

	httpCreds := HttpCredentials{
		USERNAME: HTTP_USERNAME,
		PASSWORD: HTTP_PASSWORD,
	}

	// Configure Droplet User Data
	userData := parseUserData(viper.GetString("userData"), domain, httpCreds)

	// Create Droplet
	createRequest := &godo.DropletCreateRequest{
		Name:   strings.ReplaceAll(domain.FullUrl, ".", "-"),
		Region: viper.GetString("droplet.region"),
		Size:   viper.GetString("droplet.size"),
		Image: godo.DropletCreateImage{
			Slug: viper.GetString("droplet.image.slug"),
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{ID: infra.SshKey.ID},
		},
		Tags:     infra.Tags,
		VPCUUID:  infra.Vpc.ID,
		UserData: userData,
	}

	d, _, err := client.Droplets.Create(api.ctx, createRequest)

	if err != nil {
		fmt.Printf("Error creating droplet: %s\n\n", err)
		os.Exit(1)
	}
	// Add Domain Record for new Droplet
	time.Sleep(15 * time.Second)

	droplet := fetchDroplet(api, d.ID)
	ip, _ := droplet.PublicIPv4()
	fmt.Printf("Login with username: %s\npassword: %s\n", httpCreds.USERNAME, httpCreds.PASSWORD)

	return droplet, ip
}

func destroyDroplet(api doApi, droplet *godo.Droplet) {
	client := godo.NewFromToken(api.token)
	fmt.Printf("Destroying droplet: %s\n\n", droplet.Name)
	_, err := client.Droplets.Delete(api.ctx, droplet.ID)
	if err != nil {
		fmt.Printf("There was an error deleting the droplet: %s\n\n", err)
		os.Exit(1)
	}
}

func fetchDroplet(api doApi, id int) *godo.Droplet {
	client := godo.NewFromToken(api.token)
	droplet, _, err := client.Droplets.Get(api.ctx, id)
	if err != nil {
		fmt.Println("Error getting droplet.")
		os.Exit(1)
	}
	return droplet
}

func parseUserData(ud string, domain Domain, creds HttpCredentials) string {
	dom := strings.ReplaceAll(ud, "FULL_THEIA_DOMAIN", domain.FullUrl)
	us := strings.ReplaceAll(dom, "RANDOM_USERNAME", creds.USERNAME)
	userData := strings.ReplaceAll(us, "RANDOM_PASSWORD", creds.PASSWORD)
	fmt.Println(userData)
	return userData
}

func generateHttpPassword() string {
	p, err := password.Generate(32, 10, 0, false, false)
	if err != nil {
		fmt.Printf("Error generating password: %s. Exiting.\n", err)
		os.Exit(1)
	}
	return p
}
