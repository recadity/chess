package keeper

import (
    "github.com/recadity/chess/errors"
    "github.com/recadity/chess/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "cosmossdk.io/errors"
    "github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
    storeKey sdk.StoreKey
    cdc      *codec.Codec
}

func (k Keeper) SetGame(ctx sdk.Context, game types.ChessGame) error {
    store := ctx.KVStore(k.storeKey)
    key := []byte(game.ID)
    value := k.cdc.MustMarshalBinaryBare(&game)
    store.Set(key, value)
    return nil
}

func (k Keeper) GetGame(ctx sdk.Context, gameID string) (types.ChessGame, error) {
    store := ctx.KVStore(k.storeKey)
    key := []byte(gameID)
    if !store.Has(key) {
        return types.ChessGame{}, sdkerrors.Wrapf(types.ErrGameNotFound, "game ID %s not found", gameID)
    }
    var game types.ChessGame
    k.cdc.MustUnmarshalBinaryBare(store.Get(key), &game)
    return game, nil
}

func (k Keeper) CreateGame(ctx sdk.Context, gameID string, player1 sdk.AccAddress, player2 sdk.AccAddress) error {
    game := types.NewChessGame(gameID, player1, player2)
    return k.SetGame(ctx, game)
}

func (k Keeper) MakeMove(ctx sdk.Context, gameID string, player sdk.AccAddress, newBoardState string) error {
    game, err := k.GetGame(ctx, gameID)
    if err != nil {
        return err
    }

    if game.Status != "ongoing" {
        return sdkerrors.Wrap(types.ErrGameNotFound, "game is not ongoing")
    }

    if (game.Turn == "player1" && !game.Player1.Equals(player)) || (game.Turn == "player2" && !game.Player2.Equals(player)) {
        return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "not your turn")
    }

    // Update the board state and switch the turn
    nextTurn := "player1"
    if game.Turn == "player1" {
        nextTurn = "player2"
    }
    game.UpdateBoardState(newBoardState, nextTurn)
    return k.SetGame(ctx, game)
}

func (k Keeper) EndGame(ctx sdk.Context, gameID string, finalBoardState string, status string) error {
    game, err := k.GetGame(ctx, gameID)
    if err != nil {
        return err
    }

    game.EndGame(status, finalBoardState)
    return k.SetGame(ctx, game)
}
