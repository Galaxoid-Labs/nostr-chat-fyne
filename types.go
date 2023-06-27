package main

import (
	"github.com/nbd-wtf/go-nostr"
	"github.com/puzpuzpuz/xsync"
)

type ChatRelay struct {
	Relay         nostr.Relay
	Subscriptions *xsync.MapOf[string, *nostr.Subscription]
	Groups        *xsync.MapOf[string, *ChatGroup]
}

type ChatGroup struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	ChatMessages []*nostr.Event `json:"chat_messages"`
}

type LeftMenuItem struct {
	RelayURL  string `json:"relay_url"`
	GroupName string `json:"group_name"`
}

type SavedRelay struct {
	URL    string   `json:"url"`
	Groups []string `json:"groups"`
}
