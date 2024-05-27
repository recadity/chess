package chess

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"

    "github.com/recadity/chess/keeper"
    "github.com/recadity/chess/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/errors"
)

/ Initial chessboard setup in FEN (Forsyth-Edwards Notation) format
const initialBoard = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"

func NewHandler(k keeper.Keeper) sdk.Handler {
    return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
        switch msg := msg.(type) {
        case types.MsgCreateGame:
            return handleMsgCreateGame(ctx, k, msg)
        case types.MsgMakeMove:
            return handleMsgMakeMove(ctx, k, msg)
        default:
            return nil, errors.Wrap(errors.ErrUnknownRequest, "unrecognized chess message type")
        }
    }
}

func handleMsgCreateGame(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateGame) (*sdk.Result, error) {
    // Generate a unique game ID using the current timestamp and creator address
    gameID := generateGameID(msg.Creator.String())

    // Initialize the chess game
    game := types.ChessGame{
        ID:     gameID,
        Board:  initialBoard,
        Turn:   "white",
        Winner: "",
    }

    // Store the game in the state
    if err := k.SetGame(ctx, game); err != nil {
        return nil, errors.Wrap(err, "failed to create game")
    }

    // Log the game creation
    ctx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            types.EventTypeCreateGame,
            sdk.NewAttribute(types.AttributeKeyGameID, gameID),
            sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator.String()),
        ),
    })

    // Return the result
    return &sdk.Result{
        Events: ctx.EventManager().Events(),
    }, nil
}

func generateGameID(creator string) string {
    // Use a combination of the current timestamp and creator address to generate a unique game ID
    hash := sha256.New()
    hash.Write([]byte(fmt.Sprintf("%s:%d", creator, time.Now().UnixNano())))
    return hex.EncodeToString(hash.Sum(nil))
}

func handleMsgMakeMove(ctx sdk.Context, k keeper.Keeper, msg types.MsgMakeMove) (*sdk.Result, error) {
    // Retrieve the game from the state
    game, err := k.GetGame(ctx, msg.GameID)
    if err != nil {
        return nil, errors.Wrap(types.ErrGameNotFound, "game not found")
    }

    // Validate the move
    if err := validateMove(game, msg); err != nil {
        return nil, errors.Wrap(err, "invalid move")
    }

    // Apply the move to the game state
    applyMove(&game, msg)

    // Check for win condition
    if checkWinCondition(game) {
        game.Winner = game.Turn
    }

    // Update the game state and store it
    if err := k.SetGame(ctx, game); err != nil {
        return nil, errors.Wrap(err, "failed to save game state")
    }

    // Emit events
    ctx.EventManager().EmitEvents(sdk.Events{
        sdk.NewEvent(
            types.EventTypeMakeMove,
            sdk.NewAttribute(types.AttributeKeyGameID, msg.GameID),
            sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
            sdk.NewAttribute(types.AttributeKeyTo, msg.To),
            sdk.NewAttribute(types.AttributeKeyTurn, game.Turn),
        ),
    })

    // Return the result
    return &sdk.Result{
        Events: ctx.EventManager().Events(),
    }, nil
}

// Example of a simple move validation function
func validateMove(game types.ChessGame, msg types.MsgMakeMove) error {
    // Ensure it's the correct player's turn
    if game.Turn == "white" && msg.Creator != game.WhitePlayer {
        return types.ErrNotYourTurn
    }
    if game.Turn == "black" && msg.Creator != game.BlackPlayer {
        return types.ErrNotYourTurn
    }

    // Additional validation logic, such as checking legality of the move, can be added here

    return nil
}

// Example of applying a move to the game state
func applyMove(game *types.ChessGame, msg types.MsgMakeMove) {
    // Update the board state based on the move
    fromRow, fromCol := parsePosition(msg.From)
    toRow, toCol := parsePosition(msg.To)
    game.Board[toRow][toCol] = game.Board[fromRow][fromCol]
    game.Board[fromRow][fromCol] = ""

    // Change the turn to the other player
    if game.Turn == "white" {
        game.Turn = "black"
    } else {
        game.Turn = "white"
    }
}

// Example of a simple win condition check
func checkWinCondition(game types.ChessGame) bool {
    // Placeholder for win condition logic
    // For example, check if a king has been captured
    return false
}

// Helper function to parse a board position (e.g., "e2" -> row, col)
func parsePosition(pos string) (int, int) {
    col := int(pos[0] - 'a')
    row := int(pos[1] - '1')
    return row, col
}
