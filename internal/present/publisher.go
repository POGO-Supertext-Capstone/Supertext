// Copyright (C) 2026 Jamie Murphy
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package present

import (
	"sync"
)

// SlidePublisher is a thread-safe pub/sub controller. That is,
// it allows real-time messaging between clients by pushing a message to all stored channels.
type SlidePublisher struct {
	mu   sync.RWMutex
	show int

	clients map[chan string]struct{}
}

// Create a new Slide Publisher with a blank list of clients.
func NewSlidePublisher() *SlidePublisher {
	return &SlidePublisher{
		clients: make(map[chan string]struct{}),
		show:    0,
	}
}

// Subscribe() creates a new client and adds it to the client list internally.
// It is intended to be used to register a web session.
func (sp *SlidePublisher) Subscribe() chan string {
	ch := make(chan string, 4)
	sp.mu.Lock()
	sp.clients[ch] = struct{}{}
	sp.mu.Unlock()
	return ch
}

// Unsubscribe() removes a client, deleting it from the client list.
// It is intended to be used when a web session disconnects from the service.
func (sp *SlidePublisher) Unsubscribe(channel chan string) bool {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	for ch := range sp.clients {
		if ch == channel {
			delete(sp.clients, ch)
			close(ch)
			return true
		}
	}
	return false
}

// Notify pushes a message to all subscried clients.
func (sp *SlidePublisher) Notify(msg string) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()
	for ch := range sp.clients {
		select {
		case ch <- msg: // Push a message to each channel
		default: // Drop and keep going if a subscriber is not ready quite yet.
		}
	}
}

// Quit() removes all clients, for when the server shuts down gracefully.
func (sp *SlidePublisher) Quit() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	for ch := range sp.clients {
		delete(sp.clients, ch)
		close(ch)
	}
}
