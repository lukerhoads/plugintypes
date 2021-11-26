package plugintypes

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type HookMapper interface {
	Hooks() []string
}

type HookMapperRPC struct {
	client *rpc.Client
}

func (c *HookMapperRPC) Hooks() []string {
	var Hooks []string
	cErr := c.client.Call("Plugin.Hooks", new(interface{}), &Hooks)
	if cErr != nil {
		panic(cErr)
	}

	return Hooks
}

type HookMapperRPCServer struct {
	Impl HookMapper
}

func (c *HookMapperRPCServer) Hooks(args interface{}, resp *[]string) error {
	*resp = c.Impl.Hooks()
	return nil
}

type HookMapperPlugin struct {
	Impl HookMapper
}

func (p *HookMapperPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &HookMapperRPCServer{Impl: p.Impl}, nil
}

func (HookMapperPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &HookMapperRPC{client: c}, nil
}

type Hook struct {
	parentCommand []string
	name          string
	hook_type     string
}

type HookModule interface {
	ParentCommand() []string
	Name() string
	Type() string

	PreRun(*cobra.Command, []string) error
	PostRun(*cobra.Command, []string) error
}

type HookModuleRPC struct {
	client *rpc.Client
}

func (h *HookModuleRPC) ParentCommand() []string {
	var parentHook []string
	cErr := h.client.Call("Plugin.ParentHook", new(interface{}), &parentHook)
	if cErr != nil {
		panic(cErr)
	}

	return parentHook
}

func (h *HookModuleRPC) Name() string {
	var name string
	cErr := h.client.Call("Plugin.Name", new(interface{}), &name)
	if cErr != nil {
		panic(cErr)
	}

	return name
}

func (h *HookModuleRPC) Type() string {
	var hook_type string
	cErr := h.client.Call("Plugin.Type", new(interface{}), &hook_type)
	if cErr != nil {
		panic(cErr)
	}

	return hook_type
}

func (h *HookModuleRPC) PreRun(cmd *cobra.Command, args []string) error {
	var err error
	cErr := h.client.Call("Plugin.PreRun", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

func (h *HookModuleRPC) PostRun(cmd *cobra.Command, args []string) error {
	var err error
	cErr := h.client.Call("Plugin.PostRun", ExecArgs{cmd, args}, &err)
	if cErr != nil {
		panic(cErr)
	}

	return err
}

type HookRPCServer struct {
	Impl HookModule
}

func (h *HookRPCServer) ParentCommand(args interface{}, resp *[]string) error {
	*resp = h.Impl.ParentCommand()
	return nil
}

func (h *HookRPCServer) Name(args interface{}, resp *string) error {
	*resp = h.Impl.Name()
	return nil
}

func (h *HookRPCServer) Type(args interface{}, resp *string) error {
	*resp = h.Impl.Type()
	return nil
}

func (h *HookRPCServer) PreRun(args ExecArgs, resp *error) error {
	*resp = h.Impl.PreRun(args.cmd, args.args)
	return nil
}

func (h *HookRPCServer) PostRun(args ExecArgs, resp *error) error {
	*resp = h.Impl.PostRun(args.cmd, args.args)
	return nil
}

type HookPlugin struct {
	Impl HookModule
}

func (p *HookPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &HookRPCServer{Impl: p.Impl}, nil
}

func (HookPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &HookModuleRPC{client: c}, nil
}
