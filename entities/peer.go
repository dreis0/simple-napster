package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Peer struct {
	bun.BaseModel `bun:"table:peers,alias:p"`
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	IP            string
	Port          int32
	Active        bool
}
