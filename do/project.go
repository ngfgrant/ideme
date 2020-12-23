package do

import (
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"strings"
)

func createProject(api API, project map[string]string) *godo.Project {
	client := godo.NewFromToken(api.token)
	request := &godo.CreateProjectRequest{
		Name:        project["name"],
		Description: project["description"],
		Purpose:     project["purpose"],
		Environment: project["env"],
	}
	result, _, err := client.Projects.Create(api.ctx, request)

	if err != nil {
		if strings.Contains(err.Error(), "409") {
			fmt.Println("Project already exists, skipping creation.")
			return getProject(api, project["name"])
		}
		fmt.Printf("Error creating Project. %s\n", err)
		os.Exit(1)
	}
	return result
}

func getProject(api API, name string) *godo.Project {
	client := godo.NewFromToken(api.token)
	var projectId string
	projects := listProjects(api)

	for _, p := range projects {
		if p.Name == name {
			projectId = p.ID
		}
	}

	project, _, err := client.Projects.Get(api.ctx, projectId)

	if err != nil {
		fmt.Printf("Error fetching Project: %s. Exiting.", err)
		os.Exit(1)
	}

	return project
}

func listProjects(api API) []godo.Project {
	client := godo.NewFromToken(api.token)
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	projects, _, err := client.Projects.List(api.ctx, opts)

	if err != nil {
		fmt.Printf("Error fetching Projects. %s", err)
		os.Exit(1)
	}

	return projects
}

func addResourcesToProject(api API, projectId string, resources []string) {
	client := godo.NewFromToken(api.token)
	for _, r := range resources {
		_, _, err := client.Projects.AssignResources(api.ctx, projectId, r)
		if err != nil {
			fmt.Printf("Error adding resources to Project: %s\n", err)
			os.Exit(1)
		}
	}
}
