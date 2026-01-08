package postrachain

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"postra-chain/x/postrachain/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListPost",
					Use:       "list-post",
					Short:     "List all post",
				},
				{
					RpcMethod:      "GetPost",
					Use:            "get-post [id]",
					Short:          "Gets a post by id",
					Alias:          []string{"show-post"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreatePost",
					Use:            "create-post [title] [content-uri] [content-hash]",
					Short:          "Create post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "title"}, {ProtoField: "content_uri"}, {ProtoField: "content_hash"}},
				},
				{
					RpcMethod:      "UpdatePost",
					Use:            "update-post [id] [title] [content-uri] [content-hash]",
					Short:          "Update post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "title"}, {ProtoField: "content_uri"}, {ProtoField: "content_hash"}},
				},
				{
					RpcMethod:      "DeletePost",
					Use:            "delete-post [id]",
					Short:          "Delete post",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
