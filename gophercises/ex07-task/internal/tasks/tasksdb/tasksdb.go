package tasksdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type TaskDTO struct {
	Desc string
}

func Read() ([]*TaskDTO, error) {
	b, err := ioutil.ReadFile("tasks.db")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	ts := []*TaskDTO{}
	err = json.Unmarshal(b, &ts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tasks: %s", err)
	}

	return ts, nil
}

func Write(ts []*TaskDTO) error {
	b, err := json.Marshal(ts)
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %s", err)
	}

	err = ioutil.WriteFile("tasks.db", b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	return nil
}
