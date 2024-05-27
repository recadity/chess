package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "path/to/local/chess/x/chess/types"
)

type Keeper struct {
    storeKey sdk.StoreKey
    cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
    return Keeper{
        storeKey: storeKey,
        cdc:      cdc,
    }
}

func (k Keeper) CreateGame(ctx sdk.Context, game types.MsgCreateGame) {
    store := ctx.KVStore(k.storeKey)
    // Logic to store game details
}

func (k Keeper) MakeMove(ctx sdk.Context, move types.MsgMakeMove) {
    store := ctx.KVStore(k.storeKey)
    // Logic to update game with the move
}

