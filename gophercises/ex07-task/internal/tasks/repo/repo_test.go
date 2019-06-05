package repo

import (
	"testing"
)

func TestReadingEmptyDB(tst *testing.T) {
	// err := Create()
	// if err != nil {
	// 	tst.Errorf("failed to create db: %s", err)
	// }

	// ts, err := Read()
	// if err != nil {
	// 	tst.Errorf("failed to read from empty db: %s", err)
	// }

	// if len(ts) != 0 {
	// 	tst.Errorf("failed to read from empty db. Expected tasks with length 0, got %d", len(ts))
	// }

	// dbconn.RemoveDB()
}

func TestReadingSeveralElements(tst *testing.T) {
	// err := Create()
	// if err != nil {
	// 	tst.Errorf("failed to create db: %s", err)
	// }

	// tsInput := []TaskDTO{
	// 	TaskDTO{"first"},
	// 	TaskDTO{"second"},
	// 	TaskDTO{"third"},
	// }

	// err = Add(tsInput)
	// if err != nil {
	// 	tst.Errorf("failed to write to db: %s", err)
	// }

	// ts, err := Read()
	// if err != nil {
	// 	tst.Errorf("failed to read from empty db: %s", err)
	// }

	// if len(ts) != 3 {
	// 	tst.Errorf("incorrect tasks count. Expected 3 tasks, got %d", len(ts))
	// }

	// for i, tv := range ts {
	// 	if tv.Desc != tsInput[i].Desc {
	// 		tst.Errorf("incorrect task. Task %d: expected \"%s\", got \"%s\"", i, tsInput[i].Desc, tv.Desc)
	// 	}
	// }

	// dbconn.RemoveDB()
}
