package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"postra-chain/x/postrachain/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreatePost(ctx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if err := types.ValidatePostFields(msg.Title, msg.ContentUri, msg.ContentHash); err != nil {
		return nil, err
	}

	nextId, err := k.PostSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var post = types.Post{
		Id:          nextId,
		Creator:     msg.Creator,
		Title:       msg.Title,
		ContentUri:  msg.ContentUri,
		ContentHash: msg.ContentHash,
		CreatedAt:   sdkCtx.BlockTime().Unix(),
	}

	if err = k.Post.Set(
		ctx,
		nextId,
		post,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set post")
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"post_created",
			sdk.NewAttribute("id", strconv.FormatUint(nextId, 10)),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("title", msg.Title),
			sdk.NewAttribute("content_uri", msg.ContentUri),
			sdk.NewAttribute("content_hash", msg.ContentHash),
		),
	)

	return &types.MsgCreatePostResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdatePost(ctx context.Context, msg *types.MsgUpdatePost) (*types.MsgUpdatePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	if err := types.ValidatePostFields(msg.Title, msg.ContentUri, msg.ContentHash); err != nil {
		return nil, err
	}

	// Checks that the element exists
	val, err := k.Post.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrapf(types.ErrPostNotFound, "post %d not found", msg.Id)
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get post")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var post = types.Post{
		Creator:     val.Creator,
		Id:          val.Id,
		Title:       msg.Title,
		ContentUri:  msg.ContentUri,
		ContentHash: msg.ContentHash,
		CreatedAt:   val.CreatedAt,
	}

	if err := k.Post.Set(ctx, msg.Id, post); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update post")
	}

	return &types.MsgUpdatePostResponse{}, nil
}

func (k msgServer) DeletePost(ctx context.Context, msg *types.MsgDeletePost) (*types.MsgDeletePostResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Post.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrapf(types.ErrPostNotFound, "post %d not found", msg.Id)
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get post")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.Post.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete post")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"post_deleted",
			sdk.NewAttribute("id", strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("title", val.Title),
			sdk.NewAttribute("content_uri", val.ContentUri),
			sdk.NewAttribute("content_hash", val.ContentHash),
		),
	)

	return &types.MsgDeletePostResponse{}, nil
}
