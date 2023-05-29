package main

import "github.com/nbd-wtf/go-nostr"

type ChatRelay struct {
	Relay         nostr.Relay
	Subscriptions map[string]nostr.Subscription
	Groups        map[string]ChatGroup
}

type ChatGroup struct {
	Name         string        `json:"name"`
	ChatMessages []ChatMessage `json:"chat_messages"`
}

type ChatMessage struct {
	ID        string          `json:"id"`
	PubKey    string          `json:"pubkey"`
	CreatedAt nostr.Timestamp `json:"created_at"`
	Content   string          `json:"content"`
}

type LeftMenuItem struct {
	RelayURL  string `json:"relay_url"`
	GroupName string `json:"group_name"`
}
