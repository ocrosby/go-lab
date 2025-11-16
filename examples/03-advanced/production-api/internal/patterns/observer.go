package patterns

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/ocrosby/go-lab/projects/api/internal/domain"
)

type UserEventType string

const (
	UserCreated UserEventType = "user_created"
	UserUpdated UserEventType = "user_updated"
	UserDeleted UserEventType = "user_deleted"
)

type UserEvent struct {
	Type UserEventType
	User *domain.User
}

type UserEventObserver interface {
	OnUserEvent(ctx context.Context, event UserEvent)
}

type UserEventSubject interface {
	Subscribe(observer UserEventObserver)
	Unsubscribe(observer UserEventObserver)
	Notify(ctx context.Context, event UserEvent)
}

type userEventSubject struct {
	observers []UserEventObserver
	mutex     sync.RWMutex
}

func NewUserEventSubject() UserEventSubject {
	return &userEventSubject{
		observers: make([]UserEventObserver, 0),
	}
}

func (s *userEventSubject) Subscribe(observer UserEventObserver) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *userEventSubject) Unsubscribe(observer UserEventObserver) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *userEventSubject) Notify(ctx context.Context, event UserEvent) {
	s.mutex.RLock()
	observers := make([]UserEventObserver, len(s.observers))
	copy(observers, s.observers)
	s.mutex.RUnlock()

	for _, observer := range observers {
		go observer.OnUserEvent(ctx, event)
	}
}

type LoggingUserEventObserver struct {
	logger *zap.Logger
}

func NewLoggingUserEventObserver(logger *zap.Logger) UserEventObserver {
	return &LoggingUserEventObserver{logger: logger}
}

func (o *LoggingUserEventObserver) OnUserEvent(ctx context.Context, event UserEvent) {
	o.logger.Info("User event occurred",
		zap.String("event_type", string(event.Type)),
		zap.String("user_id", event.User.ID),
		zap.String("user_email", event.User.Email),
	)
}
