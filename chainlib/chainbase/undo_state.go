package chainbase

import (
	"datx_chain/utils/common"
)

type undo_state struct {
	//
	old_next_id common.Hash

	//revision
	revision int64

	//cache the values of modified on the latest time
	old_values map[uint64][]byte

	//cache the values of removed on the latest time
	removed_values map[uint64][]byte

	//cache the values of inserted on the latest time
	new_values map[uint64][]byte
}
