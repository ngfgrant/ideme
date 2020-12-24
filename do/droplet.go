package do

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

func createDroplet(api API, infra DOInfrastructure, domain Domain, image string) (*godo.Droplet, string) {
	client := godo.NewFromToken(api.token)

	// HTTP Login Details
	HTTP_USERNAME := strings.ReplaceAll(domain.Subdomain, ".", "-")
	HTTP_PASSWORD := generateHttpPassword()

	httpCreds := HttpCredentials{
		USERNAME: HTTP_USERNAME,
		PASSWORD: HTTP_PASSWORD,
	}

	// Configure Droplet User Data
	userData := parseUserData(viper.GetString("userData"), domain, httpCreds, image)

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

	droplet := getDropletByID(api, d.ID)
	ip, _ := droplet.PublicIPv4()
	fmt.Println("Application Launching. This may take a few minutes.")
	fmt.Printf("Login: %s\nUsername: %s\nPassword: %s\n", "https://"+domain.FullUrl, httpCreds.USERNAME, httpCreds.PASSWORD)

	return droplet, ip
}

func deleteDroplet(api API, droplet *godo.Droplet) {
	client := godo.NewFromToken(api.token)
	fmt.Printf("Deleting droplet: %s\n", droplet.Name)
	_, err := client.Droplets.Delete(api.ctx, droplet.ID)
	if err != nil {
		fmt.Printf("There was an error deleting the droplet: %s\n", err)
		os.Exit(1)
	}
}

func getDropletByID(api API, id int) *godo.Droplet {
	client := godo.NewFromToken(api.token)
	droplet, _, err := client.Droplets.Get(api.ctx, id)
	if err != nil {
		fmt.Println("Error getting droplet.")
		os.Exit(1)
	}
	return droplet
}

func getDropletByName(api API, name string) *godo.Droplet {
	client := godo.NewFromToken(api.token)
	droplets := listDroplets(api)

	var id int
	for _, d := range droplets {
		if d.Name == name {
			id = d.ID
			break
		}
	}

	if id <= 0 {
		fmt.Printf("Could not find droplet: %s\n", name)
		os.Exit(1)
	}

	droplet, _, err := client.Droplets.Get(api.ctx, id)
	if err != nil {
		fmt.Println("Error getting droplet.")
		os.Exit(1)
	}
	return droplet
}

func listDroplets(api API) []godo.Droplet {
	client := godo.NewFromToken(api.token)
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	droplets, _, lerr := client.Droplets.List(api.ctx, opts)
	if lerr != nil {
		fmt.Println("Could not fetch droplets.")
	}
	return droplets
}

func parseUserData(ud string, domain Domain, creds HttpCredentials, image string) string {
	dom := strings.ReplaceAll(ud, "FULL_APP_DOMAIN", domain.FullUrl)
	us := strings.ReplaceAll(dom, "RANDOM_USERNAME", creds.USERNAME)
	i := strings.ReplaceAll(us, "APPLICATION_DOCKER_IMAGE", image)
	userData := strings.ReplaceAll(i, "RANDOM_PASSWORD", creds.PASSWORD)
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
