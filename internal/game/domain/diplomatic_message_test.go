package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewDiplomaticMessage_InitializesUnread(t *testing.T) {
	sender := uuid.Must(uuid.NewV7())
	receiver := uuid.Must(uuid.NewV7())

	msg, err := NewDiplomaticMessage(sender, receiver, nil, nil, nil, nil, DiplomaticMessageContentGreetingFriendly)
	if err != nil {
		t.Fatalf("expected message to be created, got error: %v", err)
	}
	if msg.IsRead {
		t.Fatalf("expected new message to be unread")
	}
}
