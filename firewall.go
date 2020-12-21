package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type firewallConfig struct {
	Name     string
	Tags     []string
	Inbound  map[string][]string
	Outbound map[string][]string
}

func createFirewall(api doApi, firewall interface{}) *godo.Firewall {

	var fw firewallConfig

	err := viper.Unmarshal(&fw)
	if err != nil {
		fmt.Printf("bad Unmarshal %s", err)
	}

	fmt.Printf("%s\n", fw.Inbound)

	client := godo.NewFromToken(api.token)
	allAddresses := []string{"0.0.0.0/0", "::/0"}
	firewallName := "ideme-firewall"
	request := &godo.FirewallRequest{
		Name: firewallName,
		InboundRules: []godo.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "443",
				Sources: &godo.Sources{
					Addresses: allAddresses,
				},
			},
		},
		OutboundRules: []godo.OutboundRule{
			{
				Protocol:  "tcp",
				PortRange: "443",
				Destinations: &godo.Destinations{
					Addresses: allAddresses,
				},
			},
		},
		Tags: fw.Tags,
	}
	result, _, err := client.Firewalls.Create(api.ctx, request)
	if err != nil {
		if strings.Contains(err.Error(), "409") {
			fmt.Println("Firewall already exists, skipping creation.")
			return getFirewall(api, firewallName)
		}
		fmt.Printf("There was a problem creating the Firewall: %s\n", err)
	}
	return result
}

func getFirewall(api doApi, name string) *godo.Firewall {
	client := godo.NewFromToken(api.token)
	firewalls := listFirewalls(api)
	var firewallId string

	for _, f := range firewalls {
		if f.Name == name {
			firewallId = f.ID
			break
		}
	}

	firewall, _, err := client.Firewalls.Get(api.ctx, firewallId)

	if err != nil {
		fmt.Printf("Error fetching firewall: %s. Exiting.\n", name)
		os.Exit(1)
	}

	return firewall
}

func listFirewalls(api doApi) []godo.Firewall {
	client := godo.NewFromToken(api.token)
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	firewalls, _, err := client.Firewalls.List(api.ctx, opts)
	if err != nil {
		fmt.Printf("Error fetching firewalls. %s", err)
		os.Exit(1)
	}

	return firewalls
}
