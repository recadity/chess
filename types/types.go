package types

const (
    ModuleName = "chess"
    StoreKey   = ModuleName
    RouterKey  = ModuleName
    QuerierRoute = ModuleName
)

type ChessGame struct {
    ID     string
    Board  [8][8]string
    Turn   string
    Winner string
}

type Move struct {
    GameID string
    From   string
    To     string
}
