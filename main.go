package main

import (
	"github.com/ibuildthecloud/finalizers/pkg/app"
	cli "github.com/rancher/wrangler-cli"

	// disable all client-go auth provides
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cli.Main(app.New())
}
