package app

import (
	"fmt"
	"github.com/rancher/wrangler-cli"
	"github.com/spf13/cobra"
)

func NewSubCommand() *cobra.Command {
	return cli.Command(&SubCommand{}, cobra.Command{
		Short: "Add some short description",
		Long: "Add some long description",
	})
}

type SubCommand struct {
	OptionOne string `usage:"Some usage description"`
	OptionTwo string `name:"custom-name"`
}

func (s *SubCommand) Run(cmd *cobra.Command, args []string) error {
	fmt.Println("I do stuff")
	return nil
}
