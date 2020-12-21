package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"strings"
)

func createVpc(api doApi, vpc map[string]string) *godo.VPC {
	client := godo.NewFromToken(api.token)
	request := &godo.VPCCreateRequest{
		Name:        vpc["name"],
		Description: vpc["description"],
		RegionSlug:  vpc["region"],
	}
	fmt.Println("Creating VPC.")
	result, _, err := client.VPCs.Create(api.ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "422") {
			fmt.Println("VPC exists, skipping creation.")
			return getVpc(api, vpc["name"])
		}
		fmt.Printf("There was a problem creating the VPC: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("VPC created successfully.")
	return result
}

func getVpc(api doApi, name string) *godo.VPC {
	var vpcId string
	client := godo.NewFromToken(api.token)
	vpcs := listVpcs(api)

	for _, v := range vpcs {
		if v.Name == name {
			vpcId = v.ID
			break
		}
	}

	vpc, _, err := client.VPCs.Get(api.ctx, vpcId)
	if err != nil {
		fmt.Printf("VPC: %s not found: %s\n", name, err)
	}
	return vpc
}

func listVpcs(api doApi) []*godo.VPC {
	client := godo.NewFromToken(api.token)
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	vpcs, _, err := client.VPCs.List(api.ctx, opts)
	if err != nil {
		fmt.Printf("Error fetching VPCs: %s\n", err)
		os.Exit(1)
	}
	return vpcs
}
