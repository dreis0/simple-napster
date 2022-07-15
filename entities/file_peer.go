package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type FilePeer struct {
	bun.BaseModel `bun:"table:file_peer,alias:fp"`
	FileId        uuid.UUID `bun:"type:uuid,pk,default:uuid_generate_v4()"`
	PeerId        uuid.UUID `bun:"type:uuid,pk,default:uuid_generate_v4()"`
}
