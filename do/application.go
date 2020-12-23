package do

import (
	"github.com/spf13/viper"
	"github.com/tjarratt/babble"
	"strings"
)

func CreateApplication(api API, infra DOInfrastructure) {
	sub := strings.ToLower(babble.NewBabbler().Babble()) + "-" + viper.GetString("droplet.name")
	// Create Domain Names
	domain := Domain{
		TLD:       viper.GetString("infrastructure.domain.name"),
		Subdomain: sub,
		FullUrl: strings.Join(
			[]string{sub, viper.GetString("infrastructure.domain.name")}, "."),
	}

	// Create Application Droplet
	droplet, ip := createDroplet(api, infra, domain)

	// Add Resources to DO Project
	addResourcesToProject(api, infra.Project.ID, []string{droplet.URN()})

	// Add Domain Record
	addDomainRecord(api, viper.GetString("infrastructure.domain.type"), domain, ip)
}

func DeleteApplication(api API, name string) {
	dropletName := name + "-" + strings.ReplaceAll(viper.GetString("infrastructure.domain.name"), ".", "-")
	droplet := getDropletByName(api, dropletName)
	deleteDroplet(api, droplet)

	tld := viper.GetString("infrastructure.domain.name")
	domain := Domain{
		TLD:       tld,
		Subdomain: name,
		FullUrl: strings.Join(
			[]string{name, viper.GetString("infrastructure.domain.name")}, "."),
	}

	domainRecord := getDomainRecordByName(api, tld, name)
	deleteDomainRecord(api, domain, domainRecord)
}
