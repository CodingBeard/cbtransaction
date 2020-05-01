package transaction

type ActionEnum byte

var (
	ActionAdd    ActionEnum = '+'
	ActionRemove ActionEnum = '-'
	ActionClear  ActionEnum = '*'
)

func (e *ActionEnum) IsAdd() bool {
	return *e == ActionAdd
}

func (e *ActionEnum) IsRemove() bool {
	return *e == ActionRemove
}

func (e *ActionEnum) IsClear() bool {
	return *e == ActionClear
}
