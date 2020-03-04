package cbtransaction

import "testing"

func TestActionEnum_IsAdd(t *testing.T) {
	action := ActionAdd

	if !action.IsAdd() {
		t.Error("IsAdd is incorrect, got: false, want: true")
	}
	if action.IsRemove() {
		t.Error("IsRemove is incorrect, got: true, want: false")
	}
	if action.IsClear() {
		t.Error("IsClear is incorrect, got: true, want: false")
	}
}

func TestActionEnum_IsRemove(t *testing.T) {
	action := ActionRemove

	if action.IsAdd() {
		t.Error("IsAdd is incorrect, got: true, want: false")
	}
	if !action.IsRemove() {
		t.Error("IsRemove is incorrect, got: false, want: true")
	}
	if action.IsClear() {
		t.Error("IsClear is incorrect, got: true, want: false")
	}
}

func TestActionEnum_IsClear(t *testing.T) {
	action := ActionClear

	if action.IsAdd() {
		t.Error("IsAdd is incorrect, got: true, want: false")
	}
	if action.IsRemove() {
		t.Error("IsRemove is incorrect, got: true, want: false")
	}
	if !action.IsClear() {
		t.Error("IsClear is incorrect, got: false, want: true")
	}
}
