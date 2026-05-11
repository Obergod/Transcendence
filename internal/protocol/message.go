package protocol

type MessageType string

const (
    MsgInit  MessageType = "init"
    MsgMove  MessageType = "move"
    MsgState MessageType = "state"
)

type Message struct {
    Type MessageType `json:"type"`
    Data interface{} `json:"data"`
}

type MoveData struct {
    X, Y float64 `json:"x, y"`
}

type StateData struct {
    // Structure à adapter selon vos besoins
    Players map[string]struct{ X, Y float64 } `json:"players"`
    Enemies []struct{ X, Y float64 }          `json:"enemies"`
}