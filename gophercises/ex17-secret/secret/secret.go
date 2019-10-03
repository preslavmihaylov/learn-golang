package secret

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

func FileVault(encKey, filename string) (*Vault, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		bs, err := mapToBytes(make(map[string]string))
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(filename, bs, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to write serialized data to file: %s", err)
		}
	}

	return &Vault{encKey, filename}, nil
}

type Vault struct {
	encKey   string
	filename string
}

func (v *Vault) Set(key, value string) error {
	bs, err := ioutil.ReadFile(v.filename)
	if err != nil {
		return fmt.Errorf("failed to read from file: %s", err)
	}

	store, err := bytesToMap(bs)
	if err != nil {
		return err
	}

	store[key] = value

	bs, err = mapToBytes(store)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(v.filename, bs, 0666)
	if err != nil {
		return fmt.Errorf("failed to serialize data to file: %s", err)
	}

	return nil
}

func (v *Vault) Get(key string) (string, error) {
	bs, err := ioutil.ReadFile(v.filename)
	if err != nil {
		return "", fmt.Errorf("failed to read from file: %s", err)
	}

	store, err := bytesToMap(bs)
	if err != nil {
		return "", err
	}

	val, ok := store[key]
	if !ok {
		return "", fmt.Errorf("key does not exist")
	}

	return val, nil
}

func mapToBytes(mp map[string]string) ([]byte, error) {
	bs := new(bytes.Buffer)
	enc := gob.NewEncoder(bs)

	err := enc.Encode(mp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode map: %s", err)
	}

	return bs.Bytes(), nil
}

func bytesToMap(bs []byte) (map[string]string, error) {
	var store map[string]string
	dec := gob.NewDecoder(bytes.NewReader(bs))
	err := dec.Decode(&store)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %s", err)
	}

	return store, nil
}
