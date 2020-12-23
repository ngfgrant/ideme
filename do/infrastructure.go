package do

import (
	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
)

type DOInfrastructure struct {
	Vpc      *godo.VPC
	Project  *godo.Project
	Firewall *godo.Firewall
	Domain   *godo.Domain
	SshKey   *godo.Key
	Tags     []string
}

func CreateInfrastructure(api API) DOInfrastructure {

	domainConfig := viper.GetStringMapString("infrastructure.domain")
	vpcConfig := viper.GetStringMapString("infrastructure.vpc")
	projectConfig := viper.GetStringMapString("infrastructure.project")
	firewallConfig := viper.Get("infrastructure.firewall")
	sshConfig := viper.GetStringMapString("infrastructure.ssh")
	tags := viper.GetStringSlice("tags")

	createTags(api, tags)
	vpc := createVpc(api, vpcConfig)
	project := createProject(api, projectConfig)
	firewall := createFirewall(api, firewallConfig)
	domain := createDomain(api, domainConfig)
	sshKey := createSshKey(api, sshConfig)

	infra := DOInfrastructure{
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

func GetInfrastructure(api API) DOInfrastructure {
	vpc := getVpc(api, viper.GetString("infrastructure.vpc.name"))
	project := getProject(api, viper.GetString("infrastructure.project.name"))
	firewall := getFirewall(api, viper.GetString("infrastructure.project.name"))
	domain := getDomain(api, viper.GetString("infrastructure.domain.name"))
	ssh := getSshKey(api, viper.GetString("infrastructure.ssh.name"))
	tags := viper.GetStringSlice("tags")

	infra := DOInfrastructure{
		Vpc:      vpc,
		Project:  project,
		Firewall: firewall,
		Domain:   domain,
		SshKey:   ssh,
		Tags:     tags,
	}

	return infra
}
