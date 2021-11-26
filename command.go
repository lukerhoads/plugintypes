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

func (c *CommandMapperRPCServer) Registry(args interface{}, resp *[]string) error {
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
	ParentCommand []string
	Name          string
	Usage         string
	ShortDesc     string
	LongDesc      string
	NumArgs       int
}

// need rpc functions for this
type CommandModule interface {
	Exec(*cobra.Command, []string) error
}

type CommandModuleRPC struct {
	client *rpc.Client
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

func (h *CommandModuleRPCServer) Exec(args ExecArgs, resp *error) error {
	*resp = h.Impl.Exec(args.cmd, args.args)
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
