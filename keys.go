package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/nbd-wtf/go-nostr"
	"github.com/zalando/go-keyring"
)

func startKeystore() Keystore {
	if _, err := keyring.Get(APPID, USERKEY); err != nil {
		fmt.Println(err)
		return FileKeystore{}
	}

	return KeyringStore{}
}

type Keystore interface {
	Save(keyHex string) error
	Erase() error
	Sign(*nostr.Event) error
}

const (
	APPID   = "com.galaxoidlabs.nostrchat"
	USERKEY = "userkey"
)

type KeyringStore struct{}

func (_ KeyringStore) Save(key string) error {
	return keyring.Set(APPID, USERKEY, key)
}

func (_ KeyringStore) Erase() error {
	return keyring.Delete(APPID, USERKEY)
}

func (_ KeyringStore) Sign(event *nostr.Event) error {
	key, err := keyring.Get(APPID, USERKEY)
	if err != nil {
		return fmt.Errorf("couldn't load key from keyring: %w", err)
	}
	return event.Sign(key)
}

type FileKeystore struct{}

func (_ FileKeystore) path() (string, error) {
	return homedir.Expand("~/.config/nostr/nostrchat")
}

func (f FileKeystore) prepareDirectory() (string, error) {
	path, err := f.path()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", err
	}
	return path, nil
}

func (f FileKeystore) Save(key string) error {
	path, err := f.prepareDirectory()
	if err != nil {
		return err
	}
	keybin, err := hex.DecodeString(key)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, "key"), keybin, 0600)
}

func (f FileKeystore) Erase() error {
	path, err := f.path()
	if err != nil {
		return err
	}

	return os.RemoveAll(path)
}

func (f FileKeystore) Sign(event *nostr.Event) error {
	path, err := f.path()
	if err != nil {
		return err
	}

	file := filepath.Join(path, "key")
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read key from file (%s): %w", file, err)
	}
	if len(data) != 32 {
		return fmt.Errorf("key (%s) is not 32 bytes", file)
	}

	keyhex := hex.EncodeToString(data)
	return event.Sign(keyhex)
}
