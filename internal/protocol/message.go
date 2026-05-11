package protocol

type MessageType string

const (
    MsgJoin  MessageType = "join"
    MsgMove  MessageType = "move"
    MsgShoot MessageType = "shoot"
    MsgState MessageType = "state"
)

type Message struct {
    Type MessageType `json:"type"`
    Data interface{} `json:"data"`
}

type MoveData struct {
    X, Y float64
}

type StateData struct {
    Players map[string]struct{ X, Y float64 } `json:"players"`
    Enemies map[string]struct{ X, Y float64 } `json:"enemies"`
}