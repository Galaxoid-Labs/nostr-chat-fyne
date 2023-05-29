package main

import (
	"fmt"
	"image/color"
	"net/url"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/zalando/go-keyring"
	"golang.org/x/net/context"
)

const APPID = "com.galaxoidlabs.nostrchat"
const USERKEY = "userkey"
const RELAYSKEY = "relayskey"

var baseSize = fyne.Size{Width: 900, Height: 640}

var relays = make(map[string]ChatRelay, 0)
var relayMenuData = make([]LeftMenuItem, 0)
var selectedRelayUrl = ""
var selectedGroupName = "/"
var a fyne.App
var w fyne.Window

func main() {

	a = app.NewWithID(APPID)
	w = a.NewWindow("Nostr Chat")
	w.Resize(baseSize)

	// Setup the right side of the window
	var chatMessagesListWidget *widget.List
	chatMessagesListWidget = widget.NewList(
		func() int {
			if room, ok := relays[selectedRelayUrl].Groups[selectedGroupName]; ok {
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
			message.Wrapping = fyne.TextWrapWord

			vbx := container.NewVBox(container.NewPadded(pubKey))
			border := container.NewBorder(nil, nil, vbx, nil, message)

			return border
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			if room, ok := relays[selectedRelayUrl].Groups[selectedGroupName]; ok {
				var chatMessage = room.ChatMessages[i]
				pubKey := fmt.Sprintf("[ %s ]", chatMessage.PubKey[len(chatMessage.PubKey)-8:])
				message := chatMessage.Content
				o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = pubKey
				o.(*fyne.Container).Objects[0].(*widget.Label).SetText(message)
				chatMessagesListWidget.SetItemHeight(i, o.(*fyne.Container).Objects[0].(*widget.Label).MinSize().Height)
			}
		},
	)

	chatInputWidget := widget.NewMultiLineEntry()
	chatInputWidget.Wrapping = fyne.TextWrapWord
	chatInputWidget.SetPlaceHolder("Say something...")
	chatInputWidget.OnSubmitted = func(s string) {
		go func() {
			if s == "" {
				return
			}
			chatInputWidget.SetText("")
			publishChat(s)
		}()
	}

	submitChatButtonWidget := widget.NewButton("Submit", func() {
		message := chatInputWidget.Text
		if message == "" {
			return
		}
		go func() {
			chatInputWidget.SetText("")
			publishChat(message)
		}()
	})

	bottomBorderContainer := container.NewBorder(nil, nil, nil, submitChatButtonWidget, chatInputWidget)

	// Setup the left side of the window
	relaysWidgetList := widget.NewList(
		func() int {
			return len(relayMenuData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) { // CHECK out of index...
			if len(relayMenuData) > i {
				if relayMenuData[i].GroupName == "/" {
					o.(*widget.Label).SetText(relayMenuData[i].RelayURL)
					o.(*widget.Label).TextStyle = fyne.TextStyle{
						Bold:   true,
						Italic: true,
					}
				} else {
					o.(*widget.Label).SetText("    " + relayMenuData[i].GroupName)
				}
			}
		},
	)

	relaysWidgetList.OnSelected = func(id widget.ListItemID) {
		selectedRelayUrl = relayMenuData[id].RelayURL
		selectedGroupName = relayMenuData[id].GroupName
		chatMessagesListWidget.Refresh()
		chatMessagesListWidget.ScrollToBottom() // TODO: Probalby need to guard this. For instance if user has scrolled up, it shouldnt jump to bottom on its own
	}

	relaysBottomToolbarWidget := widget.NewToolbar(
		widget.NewToolbarAction(theme.AccountIcon(), func() {
			var entry = widget.NewEntry()
			entry.SetPlaceHolder("nsec1...")
			dialog.ShowForm("Import a Nostr Private Key                                             ", "Import", "Cancel", []*widget.FormItem{ // Empty space Hack to make dialog bigger
				widget.NewFormItem("Private Key", entry),
			}, func(b bool) {
				if entry.Text != "" {
					err := saveKey(entry.Text) // TODO: Handle Error
					if err != nil {
						fmt.Println("Err saving key: ", err)
					}
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
						addRelay(entry.Text, relaysWidgetList, chatMessagesListWidget)
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
				var room = entry.Text
				if room != "" {
					if !strings.HasPrefix(room, "/") {
						room = "/" + room
					}
					addGroup(room, relaysWidgetList)
				}
			}, w)
		}),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.NewConfirm("Reset local data?", "This will remove all relays and your private key.", func(b bool) {
				if b {
					relays = nil
					for _, chatRelay := range relays {
						chatRelay.Relay.Close()
					}

					relays = nil
					relayMenuData = nil
					a.Preferences().RemoveValue(RELAYSKEY)
					relaysWidgetList.Refresh()
					chatMessagesListWidget.Refresh()

					keyring.Delete(APPID, USERKEY)
				}
			}, w).Show()
		}),
	)

	leftBorderContainer := container.NewBorder(nil, container.NewPadded(relaysBottomToolbarWidget), nil, nil, container.NewPadded(relaysWidgetList))
	rightBorderContainer := container.NewBorder(nil, container.NewPadded(bottomBorderContainer), nil, nil, container.NewPadded(chatMessagesListWidget))

	splitContainer := container.NewHSplit(leftBorderContainer, rightBorderContainer)
	splitContainer.Offset = 0.35

	w.SetContent(splitContainer)

	go func() {
		relays := getRelays()
		for _, relay := range relays {
			if relay != "" { // TODO: Better relay validation
				addRelay(relay, relaysWidgetList, chatMessagesListWidget)
			}
		}
	}()

	w.ShowAndRun()

}

func addGroup(groupName string, relayList *widget.List) {
	if selectedRelayUrl == "" { // TODO: Better handling...
		return
	}
	if _, ok := relays[selectedRelayUrl].Groups[groupName]; ok {
		return
	} else {
		relays[selectedRelayUrl].Groups[groupName] = ChatGroup{
			Name:         groupName,
			ChatMessages: make([]ChatMessage, 0),
		}
		updateLeftMenuList(relayList)

		for idx, menuItem := range relayMenuData {
			if menuItem.GroupName == groupName {
				relayList.Select(idx)
				break
			}
		}

	}
}

func publishChat(message string) error {
	hex, err := keyring.Get(APPID, USERKEY)
	if err != nil {
		fmt.Print(err)
		return err
	}

	if err != nil {
		fmt.Print(err)
		return err
	}

	publicKey, err := nostr.GetPublicKey(hex)

	if err != nil {
		fmt.Print(err)
		return err
	}

	for _, chatRelay := range relays {
		if chatRelay.Relay.URL == selectedRelayUrl {
			fmt.Println("Publishing to", chatRelay.Relay.URL)
			u, err := url.Parse(chatRelay.Relay.URL)
			if err != nil {
				return err
			}
			ev := nostr.Event{
				PubKey:    publicKey,
				CreatedAt: nostr.Now(),
				Kind:      9,
				Tags:      nostr.Tags{nostr.Tag{"g", selectedGroupName, u.Host}},
				Content:   message,
			}
			err = ev.Sign(hex)
			if err != nil {
				panic(err)
			}

			ctx := context.Background()
			chatRelay.Relay.Publish(ctx, ev)
			return nil
		}
	}

	return nil
}

func addRelay(relayURL string, relaysWidgetList *widget.List, chatMessagesListWidget *widget.List) {
	if _, ok := relays[relayURL]; ok {
		return
	} else {
		ctx := context.Background()
		relay, err := nostr.RelayConnect(ctx, relayURL)
		if err != nil {
			fmt.Println("Err connecting to: ", relayURL)
			return
		}

		chatRelay := &ChatRelay{
			Relay:  *relay,
			Groups: make(map[string]ChatGroup, 0),
		}

		relays[relayURL] = *chatRelay
		selectedRelayUrl = relayURL

		lmi := LeftMenuItem{
			RelayURL:  chatRelay.Relay.URL,
			GroupName: "/",
		}

		relayMenuData = append(relayMenuData, lmi)
		relaysWidgetList.Refresh()
		relaysWidgetList.Select(len(relayMenuData) - 1)

		filters := []nostr.Filter{{
			Kinds: []int{9},
		}}

		sub, err := relay.Subscribe(ctx, filters)
		if err != nil {
			panic(err)
		}

		// Save relay
		saveRelays()

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

				updateLeftMenuList(relaysWidgetList)

			}
		}

	}
}

func updateLeftMenuList(relaysWidgetList *widget.List) {
	relayMenuData = make([]LeftMenuItem, 0)

	for _, chatRelay := range relays {

		for _, group := range chatRelay.Groups {
			if group.Name != "/" {
				lmi := LeftMenuItem{
					RelayURL:  chatRelay.Relay.URL,
					GroupName: group.Name,
				}
				relayMenuData = append(relayMenuData, lmi)
			}

		}

		flmi := LeftMenuItem{
			RelayURL:  chatRelay.Relay.URL,
			GroupName: "/",
		}
		relayMenuData = append([]LeftMenuItem{flmi}, relayMenuData...)
	}

	relaysWidgetList.Refresh()
}

func saveKey(value string) error {
	if strings.HasPrefix(value, "nsec") {
		_, hex, err := nip19.Decode(value)
		if err != nil {
			return err
		}

		err = keyring.Set(APPID, USERKEY, hex.(string))
		if err != nil {
			return err
		}
	} else {
		publicKey, err := nostr.GetPublicKey(value)
		if err != nil {
			return err
		}
		if nostr.IsValidPublicKeyHex(publicKey) {
			err = keyring.Set(APPID, USERKEY, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func saveRelays() {
	var urls = make([]string, 0)
	for _, chatRelay := range relays {
		urls = append(urls, chatRelay.Relay.URL)
	}
	relaysStr := strings.Join(urls, ",")
	a.Preferences().SetString(RELAYSKEY, relaysStr)
}

func getRelays() []string {
	relaysStr := a.Preferences().String(RELAYSKEY)
	return strings.Split(relaysStr, ",")
}

type ChatRelay struct {
	Relay  nostr.Relay
	Groups map[string]ChatGroup
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
	RelayURL  string
	GroupName string
}
