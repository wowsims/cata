package database

import (
	_ "embed"
	"fmt"

	"github.com/wowsims/mop/sim/core/proto"
	googleProto "google.golang.org/protobuf/proto"
)

//go:embed db.bin
var dbBytes []byte

//go:embed leftover_db.bin
var leftoverBytes []byte

func Load() *proto.UIDatabase {
	// 1) Unmarshal the “main” DB
	db := &proto.UIDatabase{}
	if err := googleProto.Unmarshal(dbBytes, db); err != nil {
		panic(fmt.Errorf("unmarshal db.bin: %w", err))
	}

	if len(leftoverBytes) > 0 {
		extra := &proto.UIDatabase{}
		if err := googleProto.Unmarshal(leftoverBytes, extra); err != nil {
			panic(fmt.Errorf("unmarshal leftover_db.bin: %w", err))
		}
		googleProto.Merge(db, extra)
	}

	return db
}
