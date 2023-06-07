package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"golang.org/x/net/context"
)

func publishChat(message string) error {
	for _, chatRelay := range relays {
		if chatRelay.Relay.URL == selectedRelayUrl {
			fmt.Println("Publishing to", chatRelay.Relay.URL)
			u, err := url.Parse(chatRelay.Relay.URL)
			if err != nil {
				return err
			}
			ev := nostr.Event{
				CreatedAt: nostr.Now(),
				Kind:      9,
				Tags:      nostr.Tags{nostr.Tag{"g", selectedGroupName, u.Host}},
				Content:   message,
			}
			if err := k.Sign(&ev); err != nil {
				panic(err)
			}

			ctx := context.Background()
			chatRelay.Relay.Publish(ctx, ev)
			return nil
		}
	}

	return nil
}

func listenForGroup(relayURL string, groupId string) {
	relay, ok := relays[relayURL]
	if !ok {
		fmt.Printf("can't listen to group '%s' on relay '%s': relay not registered\n", groupId, relayURL)
		return
	}

	group, ok := relay.Groups[groupId]
	if !ok {
		fmt.Printf("can't listen to group '%s' on relay '%s': group not registered\n", groupId, relayURL)
		return
	}

	sub, ok := relay.Subscriptions[groupId]
	if ok {
		fmt.Printf("already have a subscription for '%s' on relay '%s'\n", groupId, relayURL)
		return
	}

	sub = relay.Relay.Subscribe()

	for ev := range sub.Events {
		if ev.Kind == 9 {
			for _, tag := range ev.Tags {
				if tag.Key() == "g" {
					cm := ChatMessage{
						ID:        ev.ID,
						PubKey:    ev.PubKey,
						CreatedAt: ev.CreatedAt,
						Content:   ev.Content,
					}
					if group, ok := relays[sub.Relay.URL].Groups[tag.Value()]; ok {
						group.ChatMessages = append(group.ChatMessages, cm)
						sort.Slice(group.ChatMessages, func(i, j int) bool {
							return group.ChatMessages[i].CreatedAt < group.ChatMessages[j].CreatedAt
						})
						relays[sub.Relay.URL].Groups[tag.Value()] = group
					} else {
						relays[sub.Relay.URL].Groups[tag.Value()] = ChatGroup{
							Name:         tag.Value(),
							ChatMessages: []ChatMessage{cm},
						}
					}
				}
			}

			chatMessagesListWidget.Refresh()
			chatMessagesListWidget.ScrollToBottom()

			updateLeftMenuList(relaysListWidget)
		}
	}
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
	data := make([]SavedRelay, len(relays))
	r := 0
	for _, chatRelay := range relays {
		data[r] = SavedRelay{
			URL:    chatRelay.Relay.URL,
			Groups: make([]string, len(chatRelay.Groups)),
		}
		g := 0
		for group := range chatRelay.Groups {
			data[r].Groups[g] = group
			g++
		}
		r++
	}

	j, _ := json.Marshal(data)
	a.Preferences().SetString(RELAYSKEY, string(j))
}

func getRelays() []SavedRelay {
	jstr := a.Preferences().String(RELAYSKEY)
	var data []SavedRelay
	json.Unmarshal([]byte(jstr), &data)
	return data
}
