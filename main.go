package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"os"
	"time"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	token := os.Getenv("DO_TOKEN")
	client := godo.NewFromToken(token)

	ctx := context.TODO()
	doApi := doApi{ctx, token}
	infra := createInfrastructure()

	userData := viper.Get("userData")
	d := fmt.Sprintf("%v", userData)

	fmt.Println(d)

	createRequest := &godo.DropletCreateRequest{
		Name:   viper.GetString("droplet.name"),
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
		UserData: d,
	}

	droplet, _, err := client.Droplets.Create(ctx, createRequest)

	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		os.Exit(1)
	}

	addResourcesToProject(doApi, infra.Project.ID, []string{droplet.URN()})

	fmt.Printf("Droplet created: %s\n", droplet.Created)
	time.Sleep(5 * time.Minute)
	destroyDroplet(ctx, token, viper.GetString("droplet.name"))
	os.Exit(0)
}

func destroyDroplet(ctx context.Context, token string, tag string) {
	client := godo.NewFromToken(token)
	_, err := client.Droplets.DeleteByTag(ctx, tag)
	if err != nil {
		fmt.Printf("There was an error deleting the droplet: %s\n\n", err)
		os.Exit(1)
	}
}
