package player

import "github.com/jnaraujo/mcprotocol/api/uuid"

type Position struct {
	X        float64
	FeetY    float64
	HeadY    float64
	Z        float64
	OnGround bool
}

type Player struct {
	UUID     uuid.UUID
	Name     string
	Position Position
}
