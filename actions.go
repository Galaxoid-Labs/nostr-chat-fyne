package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"golang.org/x/net/context"
)

func publishChat(message string) error {
	chatRelay, _ := relays.Load(selectRelayURL)
	if chatRelay.Relay.URL == selectRelayURL {
		fmt.Println("Publishing to", chatRelay.Relay.URL)
		u, err := url.Parse(chatRelay.Relay.URL)
		if err != nil {
			panic(err)
		}
		ev := nostr.Event{
			CreatedAt: nostr.Now(),
			Kind:      9,
			Tags:      nostr.Tags{nostr.Tag{"g", selectedGroupID, u.Host}},
			Content:   message,
		}
		if err := k.Sign(&ev); err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		res, err := chatRelay.Relay.Publish(ctx, ev)
		if err != nil {
			return err
		}

		if res == nostr.PublishStatusFailed {
			return fmt.Errorf("event not acknowledged by relay in 1 min")
		}
		return nil
	}

	return nil
}

func saveKey(value string) error {
	if strings.HasPrefix(value, "nsec") {
		_, hex, err := nip19.Decode(value)
		if err != nil {
			return err
		}

		err = k.Save(hex.(string))
		if err != nil {
			return err
		}
	} else {
		publicKey, err := nostr.GetPublicKey(value)
		if err != nil {
			return err
		}
		if nostr.IsValidPublicKeyHex(publicKey) {
			if err := k.Save(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func saveRelays() {
	data := make([]SavedRelay, relays.Size())
	r := 0
	relays.Range(func(_ string, chatRelay *ChatRelay) bool {
		data[r] = SavedRelay{
			URL:    chatRelay.Relay.URL,
			Groups: make([]string, chatRelay.Groups.Size()),
		}
		g := 0
		chatRelay.Groups.Range(func(_ string, group *ChatGroup) bool {
			data[r].Groups[g] = group.ID
			g++
			return true
		})
		r++
		return true
	})

	j, _ := json.Marshal(data)
	a.Preferences().SetString(RELAYSKEY, string(j))
}

func getRelays() []SavedRelay {
	jstr := a.Preferences().String(RELAYSKEY)
	var data []SavedRelay
	json.Unmarshal([]byte(jstr), &data)
	return data
}
