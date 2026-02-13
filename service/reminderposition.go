package service

import (
	"container/heap"
	"time"

	"github.com/jmoiron/sqlx"
	dbhelper "github.com/wibecoderr/Reminder-2.git/database"
	"github.com/wibecoderr/Reminder-2.git/model"
)

func NewReminderScheduler(db *sqlx.DB) *ReminderScheduler {
	s := &ReminderScheduler{
		queue:    make(PriorityQueue, 0),
		wakeChan: make(chan struct{}, 1),
		stopChan: make(chan struct{}),
		dbHelper: db,
		notifier: NewNotifier(),
	}
	s.loadPendingReminders()

	go s.run()

	return s
}

func (s *ReminderScheduler) run() {
	timer := time.NewTimer(0)
	defer timer.Stop()

	// Drain initial timer
	<-timer.C

	for {
		s.mu.RLock()
		if s.queue.Len() == 0 {
			s.mu.RUnlock()

			select {
			case <-s.wakeChan:
				continue
			case <-s.stopChan:
				return
			}
		}

		nextReminder := s.queue[0]
		now := time.Now()

		if nextReminder.Reminder.PopUpTime.After(now) {

			timer.Reset(nextReminder.Reminder.PopUpTime.Sub(now))
			s.mu.RUnlock()

			select {
			case <-timer.C:
				s.triggerReminder()
			case <-s.wakeChan:

				if !timer.Stop() {
					<-timer.C
				}
				continue
			case <-s.stopChan:
				return
			}
		} else {
			s.mu.RUnlock()

			s.triggerReminder()
		}
	}
}

func (s *ReminderScheduler) AddReminder(reminder model.Reminder) {
	s.mu.Lock()
	defer s.mu.Unlock()

	heap.Push(&s.queue, &ReminderItem{
		Reminder: reminder,
	})

	if s.queue[0].Reminder.ID == reminder.ID {
		select {
		case s.wakeChan <- struct{}{}:
		default:
		}
	}
}

func (s *ReminderScheduler) triggerReminder() {
	s.mu.Lock()
	if s.queue.Len() == 0 {
		s.mu.Unlock()
		return
	}

	item := heap.Pop(&s.queue).(*ReminderItem)
	s.mu.Unlock()

	go s.notifier.Notify(item.Reminder)

}

func (s *ReminderScheduler) loadPendingReminders() {

	query := `SELECT id, message, pop_up_time, user_id 
              FROM messages
              WHERE status = 'pending' AND pop_up_time > NOW()
              ORDER BY pop_up_time ASC`

	var reminders []model.Reminder
	err := dbhelper.DB.Select(&reminders, query)
	if err != nil {
		// Log error - you should add proper logging
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, reminder := range reminders {
		heap.Push(&s.queue, &ReminderItem{
			Reminder: reminder,
		})
	}
}
