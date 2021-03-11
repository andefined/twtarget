package commands

import (
	"fmt"

	"github.com/andefined/twtarget/utils"
	"github.com/urfave/cli"
)

// Init : Initialize a new Target
func Init(c *cli.Context) error {
	if c.Args().Get(0) == "" || c.String("conf") == "" {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	user := c.Args().Get(0)
	target := utils.NewTarget(user, c.String("conf"))

	fmt.Printf("You can continue with fetch command ex. \n\n\ttwtarget fetch --tweets %s\n\n", target.User)

	return nil
}
