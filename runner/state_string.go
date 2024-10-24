// Code generated by "stringer -type State"; DO NOT EDIT.

package runner

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Stopped-0]
	_ = x[Starting-1]
	_ = x[Running-2]
	_ = x[Stopping-3]
	_ = x[Errored-4]
}

const _State_name = "StoppedStartingRunningStoppingErrored"

var _State_index = [...]uint8{0, 7, 15, 22, 30, 37}

func (i State) String() string {
	if i >= State(len(_State_index)-1) {
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _State_name[_State_index[i]:_State_index[i+1]]
}
