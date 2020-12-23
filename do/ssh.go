package do

import (
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"strings"
)

func createSshKey(api API, ssh map[string]string) *godo.Key {
	publicKey := os.Getenv("IDEME_PUB_KEY")
	client := godo.NewFromToken(api.token)
	request := &godo.KeyCreateRequest{
		Name:      ssh["name"],
		PublicKey: publicKey,
	}

	result, _, err := client.Keys.Create(api.ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			fmt.Println("SSH public key already exists, skipping creation.")
			key := getSshKey(api, ssh["name"])
			return key
		} else {
			fmt.Printf("There was a problem creating the ssh public key: %s\n\n", err)
			os.Exit(1)
		}
	}

	return result
}

func getSshKey(api API, name string) *godo.Key {
	var keyId int
	client := godo.NewFromToken(api.token)
	keys := listSshKeys(api)

	for _, k := range keys {
		if k.Name == name {
			keyId = k.ID
			break
		}
	}

	key, _, err := client.Keys.GetByID(api.ctx, keyId)

	if err != nil {
		fmt.Printf("Error fetching SSH public key: %s. %s\n", name, err)
		os.Exit(1)
	}

	return key
}

func listSshKeys(api API) []godo.Key {
	client := godo.NewFromToken(api.token)
	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	keys, _, err := client.Keys.List(api.ctx, opts)

	if err != nil {
		fmt.Printf("Error fetching SSH public keys: %s", err)
		os.Exit(1)
	}
	return keys
}
