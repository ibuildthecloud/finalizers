package main

import (
	cli "github.com/rancher/wrangler-cli"
	"github.com/ibuildthecloud/finalizers/pkg/app"
)

func main() {
	cli.Main(app.New())
}
