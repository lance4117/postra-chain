package keeper_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"postra-chain/x/postrachain/keeper"
	"postra-chain/x/postrachain/types"
)

const validContentURI = "https://example.com/content.md"

func validContentHash() string {
	return "sha256:" + strings.Repeat("a", 64)
}

func TestPostMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	ctx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700000000, 0))
	f.ctx = ctx

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		title := fmt.Sprintf("Title %d", i)
		resp, err := srv.CreatePost(f.ctx, &types.MsgCreatePost{
			Creator:     creator,
			Title:       title,
			ContentUri:  validContentURI,
			ContentHash: validContentHash(),
		})
		require.NoError(t, err)
		require.Equal(t, i, int(resp.Id))

		post, err := f.keeper.Post.Get(f.ctx, resp.Id)
		require.NoError(t, err)
		require.Equal(t, ctx.BlockTime().Unix(), post.CreatedAt)
	}
}

func TestPostMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	ctx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700000000, 0))
	f.ctx = ctx
	createdAt := ctx.BlockTime().Unix()

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	resp, err := srv.CreatePost(f.ctx, &types.MsgCreatePost{
		Creator:     creator,
		Title:       "Original Title",
		ContentUri:  validContentURI,
		ContentHash: validContentHash(),
	})
	require.NoError(t, err)

	updateCtx := sdk.UnwrapSDKContext(f.ctx).WithBlockTime(time.Unix(1700000100, 0))
	f.ctx = updateCtx

	tests := []struct {
		desc    string
		request *types.MsgUpdatePost
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgUpdatePost{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgUpdatePost{
				Creator:     unauthorizedAddr,
				Id:          resp.Id,
				Title:       "Updated Title",
				ContentUri:  validContentURI,
				ContentHash: validContentHash(),
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "key not found",
			request: &types.MsgUpdatePost{
				Creator:     creator,
				Id:          10,
				Title:       "Updated Title",
				ContentUri:  validContentURI,
				ContentHash: validContentHash(),
			},
			err: types.ErrPostNotFound,
		},
		{
			desc:    "completed",
			request: &types.MsgUpdatePost{
				Creator:     creator,
				Id:          resp.Id,
				Title:       "Updated Title",
				ContentUri:  validContentURI,
				ContentHash: validContentHash(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdatePost(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				post, err := f.keeper.Post.Get(f.ctx, tc.request.Id)
				require.NoError(t, err)
				require.Equal(t, createdAt, post.CreatedAt)
			}
		})
	}
}

func TestPostMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	resp, err := srv.CreatePost(f.ctx, &types.MsgCreatePost{
		Creator:     creator,
		Title:       "Original Title",
		ContentUri:  validContentURI,
		ContentHash: validContentHash(),
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeletePost
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgDeletePost{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgDeletePost{Creator: unauthorizedAddr, Id: resp.Id},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "key not found",
			request: &types.MsgDeletePost{Creator: creator, Id: 10},
			err:     types.ErrPostNotFound,
		},
		{
			desc:    "completed",
			request: &types.MsgDeletePost{Creator: creator, Id: resp.Id},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeletePost(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPostMsgServerCreateValidation(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgCreatePost
		err     error
	}{
		{
			desc: "whitespace title",
			request: &types.MsgCreatePost{
				Creator:     creator,
				Title:       "   ",
				ContentUri:  validContentURI,
				ContentHash: validContentHash(),
			},
			err: types.ErrInvalidTitle,
		},
		{
			desc: "empty title",
			request: &types.MsgCreatePost{
				Creator:     creator,
				Title:       "",
				ContentUri:  validContentURI,
				ContentHash: validContentHash(),
			},
			err: types.ErrInvalidTitle,
		},
		{
			desc: "invalid content uri",
			request: &types.MsgCreatePost{
				Creator:     creator,
				Title:       "Valid Title",
				ContentUri:  "ftp://example.com/file",
				ContentHash: validContentHash(),
			},
			err: types.ErrInvalidContentURI,
		},
		{
			desc: "invalid content hash",
			request: &types.MsgCreatePost{
				Creator:     creator,
				Title:       "Valid Title",
				ContentUri:  validContentURI,
				ContentHash: "sha256:deadbeef",
			},
			err: types.ErrInvalidContentHash,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := srv.CreatePost(f.ctx, tc.request)
			require.ErrorIs(t, err, tc.err)
		})
	}
}
