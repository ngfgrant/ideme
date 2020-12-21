package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
	"os"
)

type doApi struct {
	ctx   context.Context
	token string
}

type doInfrastructure struct {
	Vpc      *godo.VPC
	Project  *godo.Project
	Firewall *godo.Firewall
	Domain   *godo.Domain
	SshKey   *godo.Key
	Tags     []string
}

func createInfrastructure() doInfrastructure {
	token := os.Getenv("DO_TOKEN")
	ctx := context.TODO()
	api := doApi{ctx, token}

	domainConfig := viper.GetStringMapString("infrastructure.domain")
	vpcConfig := viper.GetStringMapString("infrastructure.vpc")
	projectConfig := viper.GetStringMapString("infrastructure.project")
	firewallConfig := viper.Get("infrastructure.firewall")
	sshConfig := viper.GetStringMapString("infrastructure.ssh")
	tags := viper.GetStringSlice("tags")

	fmt.Printf("%s", firewallConfig)

	createTags(api, tags)
	vpc := createVpc(api, vpcConfig)
	project := createProject(api, projectConfig)
	firewall := createFirewall(api, firewallConfig)
	domain := createDomain(api, domainConfig)
	sshKey := createSshKey(api, sshConfig)

	infra := doInfrastructure{
		Vpc:      vpc,
		Project:  project,
		Firewall: firewall,
		Domain:   domain,
		SshKey:   sshKey,
		Tags:     tags,
	}

	resources := []string{
		firewall.URN(),
		domain.URN(),
	}
	addResourcesToProject(api, infra.Project.ID, resources)
	return infra
}
