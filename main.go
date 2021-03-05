package main

import (
	"github.com/ibuildthecloud/finalizers/pkg/app"
	cli "github.com/rancher/wrangler-cli"
)

func main() {
	cli.Main(app.New())
}
