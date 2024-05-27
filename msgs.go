package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgCreateGame struct {
    Creator sdk.AccAddress
}

func NewMsgCreateGame(creator sdk.AccAddress) MsgCreateGame {
    return MsgCreateGame{Creator: creator}
}

type MsgMakeMove struct {
    Creator sdk.AccAddress
    GameID  string
    From    string
    To      string
}

func NewMsgMakeMove(creator sdk.AccAddress, gameID, from, to string) MsgMakeMove {
    return MsgMakeMove{
        Creator: creator,
        GameID:  gameID,
        From:    from,
        To:      to,
    }
}
