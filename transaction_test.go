package cbtransaction

import (
	"github.com/codingbeard/cbtransaction/transaction"
	"testing"
)

func TestActionEnum_IsAdd(t *testing.T) {
	action := transaction.ActionAdd

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
	action := transaction.ActionRemove

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
	action := transaction.ActionClear

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
