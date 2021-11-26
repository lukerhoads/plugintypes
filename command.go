package plugintypes

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type CommandMapper interface {
	Commands() []string
}

type CommandMapperRPC struct {
	client *rpc.Client
}

func (c *CommandMapperRPC) Commands() []string {
	var commands []string
	cErr := c.client.Call("Plugin.Commands", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

type CommandMapperRPCServer struct {
	Impl CommandMapper
}

func (c *CommandMapperRPCServer) Commands(args interface{}, resp *[]string) error {
	*resp = c.Impl.Commands()
	return nil
}

type CommandMapperPlugin struct {
	Impl CommandMapper
}

func (p *CommandMapperPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandMapperRPCServer{Impl: p.Impl}, nil
}

func (CommandMapperPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandMapperRPC{client: c}, nil
}

type Command struct {
	parentCommand []string
	name          string
	usage         string
	shortDesc     string
	longDesc      string
	numArgs       int
}

// need getter method for the parameters!
type CommandModule interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	NumArgs() int
	Exec(*cobra.Command, []string) error
}

type CommandModuleRPC struct {
	client *rpc.Client
}

func (c *CommandModuleRPC) ParentCommand() []string {
	var commands []string
	cErr := c.client.Call("Plugin.ParentCommand", new(interface{}), &commands)
	if cErr != nil {
		panic(cErr)
	}

	return commands
}

func (c *CommandModuleRPC) Name() string {
	var name string
	cErr := c.client.Call("Plugin.Name", new(interface{}), &name)
	if cErr != nil {
		panic(cErr)
	}

	return name
}

func (c *CommandModuleRPC) Usage() string {
	var usage string
	cErr := c.client.Call("Plugin.Usage", new(interface{}), &usage)
	if cErr != nil {
		panic(cErr)
	}

	return usage
}

func (c *CommandModuleRPC) ShortDesc() string {
	var desc string
	cErr := c.client.Call("Plugin.ShortDesc", new(interface{}), &desc)
	if cErr != nil {
		panic(cErr)
	}

	return desc
}

func (c *CommandModuleRPC) LongDesc() string {
	var desc string
	cErr := c.client.Call("Plugin.LongDesc", new(interface{}), &desc)
	if cErr != nil {
		panic(cErr)
	}

	return desc
}

func (c *CommandModuleRPC) NumArgs() int {
	var numArgs int
	cErr := c.client.Call("Plugin.NumArgs", new(interface{}), &numArgs)
	if cErr != nil {
		panic(cErr)
	}

	return numArgs
}

func (c *CommandModuleRPC) Exec(cmd *cobra.Command, args []string) error {
	var err error
	cErr := c.client.Call("Plugin.Exec", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

type CommandModuleRPCServer struct {
	Impl CommandModule
}

func (c *CommandModuleRPCServer) ParentCommand(args interface{}, resp *[]string) error {
	*resp = c.Impl.ParentCommand()
	return nil
}

func (c *CommandModuleRPCServer) Name(args interface{}, resp *string) error {
	*resp = c.Impl.Name()
	return nil
}

func (c *CommandModuleRPCServer) Usage(args interface{}, resp *string) error {
	*resp = c.Impl.Usage()
	return nil
}

func (c *CommandModuleRPCServer) ShortDesc(args interface{}, resp *string) error {
	*resp = c.Impl.ShortDesc()
	return nil
}

func (c *CommandModuleRPCServer) LongDesc(args interface{}, resp *string) error {
	*resp = c.Impl.LongDesc()
	return nil
}

func (c *CommandModuleRPCServer) NumArgs(args interface{}, resp *int) error {
	*resp = c.Impl.NumArgs()
	return nil
}

func (c *CommandModuleRPCServer) Exec(args ExecArgs, resp *error) error {
	*resp = c.Impl.Exec(args.cmd, args.args)
	return nil
}

type CommandModulePlugin struct {
	Impl CommandModule
}

func (p *CommandModulePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CommandModuleRPCServer{Impl: p.Impl}, nil
}

func (CommandModulePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CommandModuleRPC{client: c}, nil
}
