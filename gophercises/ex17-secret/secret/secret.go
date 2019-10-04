package secret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func FileVault(encKey, filename string) (*Vault, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := writeMapToFile(make(map[string]string), filename, encKey)
		if err != nil {
			return nil, fmt.Errorf("failed to write map to file: %s", err)
		}
	}

	return &Vault{encKey, filename}, nil
}

type Vault struct {
	encKey   string
	filename string
}

func (v *Vault) Set(key, value string) error {
	store, err := getMapFromFile(v.filename, v.encKey)
	if err != nil {
		return err
	}

	store[key] = value

	err = writeMapToFile(store, v.filename, v.encKey)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Get(key string) (string, error) {
	store, err := getMapFromFile(v.filename, v.encKey)
	if err != nil {
		return "", err
	}

	val, ok := store[key]
	if !ok {
		return "", fmt.Errorf("key does not exist")
	}

	return val, nil
}

func getMapFromFile(filename string, key string) (map[string]string, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read from file: %s", err)
	}

	decryptedBs, err := decrypt(bs, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %s", err)
	}

	store, err := bytesToMap(decryptedBs)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func writeMapToFile(store map[string]string, filename string, key string) error {
	bs, err := mapToBytes(store)
	if err != nil {
		return err
	}

	encryptedBs, err := encrypt(bs, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %s", err)
	}

	err = ioutil.WriteFile(filename, encryptedBs, 0666)
	if err != nil {
		return fmt.Errorf("failed to serialize data to file: %s", err)
	}

	return nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %s", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %s", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %s", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decipher text: %s", err)
	}

	return plaintext, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))

	return hex.EncodeToString(hasher.Sum(nil))
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
