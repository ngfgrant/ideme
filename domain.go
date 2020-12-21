package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"strings"
)

func createDomain(api doApi, domain map[string]string) *godo.Domain {
	client := godo.NewFromToken(api.token)
	request := &godo.DomainCreateRequest{
		Name: domain["name"],
	}
	result, _, err := client.Domains.Create(api.ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			fmt.Println("Domain already exists, skipping creation.")
			return getDomain(api, domain["name"])
		}
		fmt.Printf("There was a problem creating the domain: %s \n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created domain. %s \n", result.Name)
	return result
}

func getDomain(api doApi, name string) *godo.Domain {
	client := godo.NewFromToken(api.token)
	domain, _, err := client.Domains.Get(api.ctx, name)

	if err != nil {
		fmt.Printf("Error fetching domain name: %s. %s\n", name, err)
		os.Exit(1)
	}
	return domain
}
