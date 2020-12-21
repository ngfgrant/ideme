package main

import (
	"github.com/digitalocean/godo"
)

func createTags(api doApi, tags []string) {
	client := godo.NewFromToken(api.token)
	for _, v := range tags {
		request := &godo.TagCreateRequest{
			Name: v,
		}
		client.Tags.Create(api.ctx, request)
	}
}
