package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type File struct {
	bun.BaseModel `bun:"table:files,alias:f"`
	ID            uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	Name          string    `bun:",notnull,unique"`
}
