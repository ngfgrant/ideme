package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"strings"
)

type Domain struct {
	TLD       string
	Subdomain string
	FullUrl   string
}

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

func addDomainRecord(api doApi, t string, domain Domain, ip string) *godo.DomainRecord {
	client := godo.NewFromToken(api.token)
	// Add Domain Record for new Droplet
	domainRequest := &godo.DomainRecordEditRequest{
		Type: t,
		Name: domain.Subdomain,
		Data: ip,
	}

	domainRecord, _, err := client.Domains.CreateRecord(api.ctx, domain.TLD, domainRequest)

	if err != nil {
		fmt.Printf("Domain Error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("Domain record created: %s\n", domainRecord.Name)
	return domainRecord
}

func destroyDomainRecord(api doApi, domain Domain, domainRecord *godo.DomainRecord) {
	client := godo.NewFromToken(api.token)
	_, err := client.Domains.DeleteRecord(api.ctx, domain.TLD, domainRecord.ID)
	if err != nil {
		fmt.Printf("There was an error deleting the domain record: %s\n\n", err)
		os.Exit(1)
	}
}
