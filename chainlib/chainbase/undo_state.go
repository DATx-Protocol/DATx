package chainbase

type undo_state struct {
	//
	old_next_id uint64

	//revision
	revision int64

	//cache the pairs of obj_id/obj modified on the latest time
	old_values map[uint64][]byte

	//cache the pairs of obj_id/obj removed on the latest time
	removed_values map[uint64][]byte

	//cache the pairs of obj_id/obj inserted on the latest time
	new_values map[uint64][]byte
}

func NewUndoState(next_id uint64, revision int64) *undo_state {
	var state undo_state

	state.new_values = make(map[uint64][]byte)
	state.old_values = make(map[uint64][]byte)
	state.removed_values = make(map[uint64][]byte)

	state.old_next_id = next_id
	state.revision = revision

	return &state
}
