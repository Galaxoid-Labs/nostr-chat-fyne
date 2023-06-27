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
	Picture      string         `json:"picture"`
	Subgroups    []string       `json:"subgroups"`
	ChatMessages []*nostr.Event `json:"chat_messages"`
}

type LeftMenuItem struct {
	RelayURL  string `json:"relay_url"`
	IsRoot    bool   `json:"is_root"`
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
	GroupIcon string `json:"group_icon"`
}

type SavedRelay struct {
	URL    string   `json:"url"`
	Groups []string `json:"groups"`
}
