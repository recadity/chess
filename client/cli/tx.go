package cli

import (
    "github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/client/tx"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/recadity/chess/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 2,
        RunE:                       client.ValidateCmd,
    }

    cmd.AddCommand(CmdCreateGame())
    cmd.AddCommand(CmdMakeMove())

    return cmd
}

func CmdCreateGame() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "create-game",
        Short: "Create a new game",
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }

            msg := types.NewMsgCreateGame(clientCtx.GetFromAddress())
            return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
        },
    }

    flags.AddTxFlagsToCmd(cmd)
    return cmd
}

func CmdMakeMove() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "make-move [game-id] [from] [to]",
        Short: "Make a move in the game",
        Args:  cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }

            gameID := args[0]
            from := args[1]
            to := args[2]

            msg := types.NewMsgMakeMove(clientCtx.GetFromAddress(), gameID, from, to)
            return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
        },
    }

    flags.AddTxFlagsToCmd(cmd)
    return cmd
}

