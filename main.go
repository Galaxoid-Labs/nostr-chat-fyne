package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"golang.org/x/net/context"
	"image/color"
	"sort"
)

var baseSize = fyne.Size{Width: 900, Height: 640}
var nsec = ""

var chatRelays = make(map[string]ChatRelay, 0)
var relayRoomsMenuData = make([]LeftMenuItem, 0)
var selectedRelayUrl = ""
var selectedRoomId = ""

func main() {

	a := app.New()
	w := a.NewWindow("Nostr Chat")
	w.Resize(baseSize)

	// Setup the right side of the window
	chatMessagesWidget := widget.NewList(
		func() int {
			if room, ok := chatRelays[selectedRelayUrl].Rooms[selectedRoomId]; ok {
				return len(room.ChatMessages)
			} else {
				return 0
			}
		},
		func() fyne.CanvasObject {

			pubKey := canvas.NewText("template", color.RGBA{139, 190, 178, 255})
			pubKey.TextStyle.Bold = true
			pubKey.Alignment = fyne.TextAlignLeading

			message := widget.NewLabel("template")
			message.Alignment = fyne.TextAlignLeading

			box := container.NewHBox()
			box.Add(pubKey)
			box.Add(message)

			return box
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if room, ok := chatRelays[selectedRelayUrl].Rooms[selectedRoomId]; ok {
				var chatMessage = room.ChatMessages[i]
				pubKey := fmt.Sprintf("[ %s ]", chatMessage.PubKey[len(chatMessage.PubKey)-8:])
				message := chatMessage.Content
				o.(*fyne.Container).Objects[0].(*canvas.Text).Text = pubKey
				o.(*fyne.Container).Objects[1].(*widget.Label).SetText(message)
			}
		},
	)

	inputWidget := widget.NewEntry()
	inputWidget.SetPlaceHolder("Say something...")
	inputWidget.OnSubmitted = func(s string) {
		go func() {
			if s == "" {
				return
			}
			inputWidget.SetText("")
			publishChat(s)
		}()
	}

	submitButton := widget.NewButton("Submit", func() {
		message := inputWidget.Text
		if message == "" {
			return
		}
		go func() {
			inputWidget.SetText("")
			publishChat(message)
		}()
	})

	bottomBox := container.NewBorder(nil, nil, nil, submitButton, inputWidget)

	// Setup the left side of the window
	relayRoomsWidget := widget.NewList(
		func() int {
			return len(relayRoomsMenuData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if relayRoomsMenuData[i].RoomID == "" {
				o.(*widget.Label).SetText(relayRoomsMenuData[i].RelayURL)
			} else {
				o.(*widget.Label).SetText("    " + relayRoomsMenuData[i].RoomID)
			}
		},
	)

	relayRoomsWidget.OnSelected = func(id widget.ListItemID) {
		selectedRelayUrl = relayRoomsMenuData[id].RelayURL
		selectedRoomId = relayRoomsMenuData[id].RoomID
		chatMessagesWidget.Refresh()
		chatMessagesWidget.ScrollToBottom()
	}

	// Auto add the Nostr Relay
	// go func() {
	// 	addRelay("wss://groups.nostr.com/nostr", relayRoomsWidget, chatMessagesWidget)
	// }()

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			var entry = widget.NewEntry()
			entry.SetPlaceHolder("nsec1...")
			dialog.ShowForm("Import a Nostr Private Key                                             ", "Import", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("Private Key", entry),
			}, func(b bool) {
				if entry.Text != "" {
					nsec = entry.Text // TODO: Validate
				}
			}, w)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.StorageIcon(), func() {
			var entry = widget.NewEntry()
			entry.SetPlaceHolder("Enter a Nostr Relay URL")
			dialog.ShowForm("Add a Nostr Relay                                             ", "Add", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("URL", entry),
			}, func(b bool) {
				if entry.Text != "" {
					go func() {
						addRelay(entry.Text, relayRoomsWidget, chatMessagesWidget)
					}()
				}
			}, w)
		}),
		widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			var entry = widget.NewEntry()
			entry.SetPlaceHolder("")
			dialog.ShowForm("Add a New Room                                             ", "Add", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("Room Name", entry),
			}, func(b bool) {
				if entry.Text != "" {
					addRoom(entry.Text, relayRoomsWidget)
				}
			}, w)
		}),
	)

	leftSide := container.NewBorder(nil, container.NewPadded(toolbar), nil, nil, container.NewPadded(relayRoomsWidget))
	rightSide := container.NewBorder(nil, container.NewPadded(bottomBox), nil, nil, container.NewPadded(chatMessagesWidget))

	split := container.NewHSplit(leftSide, rightSide)

	split.Offset = 0.4
	w.SetContent(split)
	w.ShowAndRun()

}

func addRoom(roomID string, relayList *widget.List) {
	if _, ok := chatRelays[selectedRelayUrl].Rooms[roomID]; ok {
		return
	} else {
		chatRelays[selectedRelayUrl].Rooms[roomID] = ChatRoom{
			ID:           roomID,
			ChatMessages: make([]ChatMessage, 0),
		}
		updateLeftMenuList(relayList)

		for idx, menuItem := range relayRoomsMenuData {
			if menuItem.RoomID == roomID {
				relayList.Select(idx)
				break
			}
		}

	}
}

func publishChat(message string) {
	if nsec == "" {
		return // TODO: Validate
	}
	_, hex, err := nip19.Decode(nsec)
	publicKey, err := nostr.GetPublicKey(hex.(string))

	if err != nil {
		return // TODO: Handle error
	}

	for _, chatRelay := range chatRelays {
		if chatRelay.Relay.URL == selectedRelayUrl {
			fmt.Println("Publishing to", chatRelay.Relay.URL)
			ev := nostr.Event{
				PubKey:    publicKey,
				CreatedAt: nostr.Now(),
				Kind:      9,
				Tags:      nostr.Tags{nostr.Tag{"r", selectedRoomId}},
				Content:   message,
			}
			err = ev.Sign(hex.(string))
			if err != nil {
				panic(err)
			}

			ctx := context.Background()
			chatRelay.Relay.Publish(ctx, ev)
			return
		}
	}
}

func addRelay(relayURL string, relayRoomsWidget *widget.List, chatMessagesWidget *widget.List) {
	if _, ok := chatRelays[relayURL]; ok {
		return
	} else {
		ctx := context.Background()
		relay, err := nostr.RelayConnect(ctx, relayURL)
		if err != nil {
			panic(err)
		}

		chatRelay := &ChatRelay{
			Relay: *relay,
			Rooms: make(map[string]ChatRoom, 0),
		}

		chatRelays[relayURL] = *chatRelay
		selectedRelayUrl = relayURL

		lmi := LeftMenuItem{
			RelayURL: chatRelay.Relay.URL,
			RoomID:   "",
		}

		relayRoomsMenuData = append(relayRoomsMenuData, lmi)
		relayRoomsWidget.Refresh()
		relayRoomsWidget.Select(len(relayRoomsMenuData) - 1)

		var filters nostr.Filters
		filters = []nostr.Filter{{
			Kinds: []int{9},
		}}

		sub, err := relay.Subscribe(ctx, filters)
		if err != nil {
			panic(err)
		}

		for ev := range sub.Events {

			if ev.Kind == 9 {

				for _, tag := range ev.Tags {
					if tag.Key() == "r" {

						cm := ChatMessage{
							ID:        ev.ID,
							PubKey:    ev.PubKey,
							CreatedAt: ev.CreatedAt,
							Content:   ev.Content,
						}

						if room, ok := chatRelays[sub.Relay.URL].Rooms[tag.Value()]; ok {

							room.ChatMessages = append(room.ChatMessages, cm)
							sort.Slice(room.ChatMessages, func(i, j int) bool {
								return room.ChatMessages[i].CreatedAt < room.ChatMessages[j].CreatedAt
							})
							chatRelays[sub.Relay.URL].Rooms[tag.Value()] = room

						} else {

							chatRelays[sub.Relay.URL].Rooms[tag.Value()] = ChatRoom{
								ID:           tag.Value(),
								ChatMessages: []ChatMessage{cm},
							}

						}

					}

				}

				chatMessagesWidget.Refresh()
				chatMessagesWidget.ScrollToBottom()

				updateLeftMenuList(relayRoomsWidget)

			}
		}

	}

}

func updateLeftMenuList(relayList *widget.List) {
	relayRoomsMenuData = make([]LeftMenuItem, 0)

	for _, chatRelay := range chatRelays {

		for _, room := range chatRelay.Rooms {
			if room.ID != "" {
				lmi := LeftMenuItem{
					RelayURL: chatRelay.Relay.URL,
					RoomID:   room.ID,
				}
				relayRoomsMenuData = append(relayRoomsMenuData, lmi)
			}

		}

		flmi := LeftMenuItem{
			RelayURL: chatRelay.Relay.URL,
			RoomID:   "",
		}
		relayRoomsMenuData = append([]LeftMenuItem{flmi}, relayRoomsMenuData...)
	}

	relayList.Refresh()
}

type ChatRelay struct {
	Relay nostr.Relay
	Rooms map[string]ChatRoom
}

type ChatRoom struct {
	ID           string        `json:"id"`
	ChatMessages []ChatMessage `json:"chat_messages"`
}

type ChatMessage struct {
	ID        string          `json:"id"`
	PubKey    string          `json:"pubkey"`
	CreatedAt nostr.Timestamp `json:"created_at"`
	Content   string          `json:"content"`
}

type LeftMenuItem struct {
	RelayURL string
	RoomID   string
}
