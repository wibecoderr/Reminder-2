package service

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/wibecoderr/Reminder-2.git/model"
)

type ReminderItem struct {
	Reminder model.Reminder
	Index    int // Index in heap
}

type PriorityQueue []*ReminderItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Reminder.PopUpTime.Before(pq[j].Reminder.PopUpTime)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*ReminderItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

type ReminderScheduler struct {
	queue    PriorityQueue
	mu       sync.RWMutex
	wakeChan chan struct{}
	stopChan chan struct{}
	dbHelper *sqlx.DB
	notifier *Notifier
}
