package services

import (
	"sync"
	"time"
)

// ReminderEvent 前端/移动端可消费的提醒事件。
type ReminderEvent struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// ReminderHub 管理 SSE 订阅并广播提醒。
type ReminderHub struct {
	mu      sync.RWMutex
	clients map[chan ReminderEvent]struct{}
}

// NewReminderHub 创建提醒事件中心。
func NewReminderHub() *ReminderHub {
	return &ReminderHub{
		clients: make(map[chan ReminderEvent]struct{}),
	}
}

// Subscribe 订阅提醒流。
func (h *ReminderHub) Subscribe() chan ReminderEvent {
	ch := make(chan ReminderEvent, 8)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

// Unsubscribe 取消订阅并关闭通道。
func (h *ReminderHub) Unsubscribe(ch chan ReminderEvent) {
	h.mu.Lock()
	if _, ok := h.clients[ch]; ok {
		delete(h.clients, ch)
		close(ch)
	}
	h.mu.Unlock()
}

// Broadcast 广播提醒到所有订阅者。
func (h *ReminderHub) Broadcast(event ReminderEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.clients {
		select {
		case ch <- event:
		default:
			// 若客户端处理慢，跳过当前事件以避免阻塞整体广播。
		}
	}
}
