package main

import "github.com/nbd-wtf/go-nostr"

func insertEventIntoDescendingList(sortedArray []*nostr.Event, event *nostr.Event) []*nostr.Event {
	size := len(sortedArray)
	start := 0
	end := size - 1
	var mid int
	position := start

	if end < 0 {
		return []*nostr.Event{event}
	} else if event.CreatedAt < sortedArray[end].CreatedAt {
		return append(sortedArray, event)
	} else if event.CreatedAt > sortedArray[start].CreatedAt {
		newArr := make([]*nostr.Event, size+1)
		newArr[0] = event
		copy(newArr[1:], sortedArray)
		return newArr
	} else if event.CreatedAt == sortedArray[start].CreatedAt {
		position = start
	} else {
		for {
			if end <= start+1 {
				position = end
				break
			}
			mid = int(start + (end-start)/2)
			if sortedArray[mid].CreatedAt > event.CreatedAt {
				start = mid
			} else if sortedArray[mid].CreatedAt < event.CreatedAt {
				end = mid
			} else {
				position = mid
				break
			}
		}
	}

	if sortedArray[position].ID != event.ID {
		if cap(sortedArray) > size {
			newArr := sortedArray[0 : size+1]
			copy(newArr[position+1:], sortedArray[position:])
			newArr[position] = event
			return newArr
		} else {
			newArr := make([]*nostr.Event, size+1, size+5)
			copy(newArr[:position], sortedArray[:position])
			copy(newArr[position+1:], sortedArray[position:])
			newArr[position] = event
			return newArr
		}
	}

	return sortedArray
}
