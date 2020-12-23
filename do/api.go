package do

import (
	"context"
	"os"
)

type API struct {
	ctx   context.Context
	token string
}

func ConfigureApi() API {
	// Configure DigitalOcean API
	token := os.Getenv("DO_TOKEN")
	ctx := context.TODO()
	api := API{ctx, token}
	return api
}
