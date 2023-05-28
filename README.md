# nostr-chat-fyne

An experimental chat client written with Fyne. Its a work in progress based on kind 9 ideas

![alt text](screenshots/ss.png)

This is very early implementation. Use at your own risk!

## Setup your dev environment

- Golang
- [Fyne Prerequisites](https://developer.fyne.io/started/#prerequisites)
- Clone repo

```
go mod tidy
go run .
```

Important! It will take quite a bit of time to complie the first time as its compiling some C libraries. Please be patient. After first compile, it will work normally.

### After running

- Add your nsec
- Add kind9 relay exmaple: wss://groups.nostr.com/nostr
- Chat away! You can also add room's

Todo:

- Lots of things.
- Add secure key storage to persist your private key.
- Persist relays
- Handle errors
- Handle removing relays
- Handle removing keys, etc

- Lots visual updates
