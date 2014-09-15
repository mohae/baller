package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mohae/car/app"
	"github.com/mohae/cli"
	"github.com/mohae/contour"
)

// CreateCommand is a Command implementation that says hello world
type CreateCommand struct {
	UI cli.Ui
}

// Help prints the help text for the run sub-command.
func (c *CreateCommand) Help() string {
	helpText := `
Usage: car create [flags] <destination> <source...>

This will create a compressed archive from the list of
sources and write it to the destination.

create supports the following flags(Type):

    --logging(bool)     enable/disable log output
    --type              compression type of the archive
    --format            the archive format to use, tar
                        is the default, --type is ignored
			when the format is zip, since it
			comes with its own compression type.
    --verbose           verbose output.
`
	return strings.TrimSpace(helpText)
}

// Run runs the hello command; the args are a variadic list of words
// to append to hello.
func (c *CreateCommand) Run(args []string) int {
	// set up the command flags
	cmdFlags := flag.NewFlagSet("run", flag.ContinueOnError)
	cmdFlags.Usage = func() {
		c.UI.Output(c.Help())
	}

	// Filter the flags from the args and update the config with them.
	// The args remaining after being filtered are returned.
	filteredArgs, err := contour.FilterArgs(cmdFlags, args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err = app.SetLogging()
	if err != nil {
		c.UI.Error(fmt.Sprintf("setup and configuration of application logging failed: %s", err))
		return 1
	}

	// Run the command in the package.
	message, err := app.Create(filteredArgs...)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(message)
	return 0
}

// Synopsis provides a precis of the hello command.
func (c *CreateCommand) Synopsis() string {
	ret := `creates a compressed archive from the specified source(s) and writes it out to the destination."
`

	return ret
}