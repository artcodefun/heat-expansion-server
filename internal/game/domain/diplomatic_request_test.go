package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewDiplomaticRequest_InitializesPendingRequest(t *testing.T) {
	SetTestNow(t, 10_000)
	senderUserID := uuid.New()
	receiverUserID := uuid.New()
	senderBaseID := 11
	receiverBaseID := 22

	request, err := NewDiplomaticRequest(senderUserID, receiverUserID, &senderBaseID, &receiverBaseID, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	if request.Status != DiplomaticRequestStatusPending {
		t.Fatalf("expected pending status, got %s", request.Status)
	}
	if request.CreatedAt != 10_000 {
		t.Fatalf("expected created_at 10000, got %d", request.CreatedAt)
	}
	if request.ExpiresAt != 10_000+int64(DiplomaticRequestExpiry.Seconds()) {
		t.Fatalf("unexpected expiry timestamp: %d", request.ExpiresAt)
	}
	if request.MessageContent() != DiplomaticMessageContentCoalitionProposal {
		t.Fatalf("expected coalition proposal content, got %s", request.MessageContent())
	}
	if len(request.PullEvents()) != 1 {
		t.Fatalf("expected one created event")
	}
}

func TestValidateAgainstRelationship_RejectsPendingDuplicate(t *testing.T) {
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}

	if err := request.ValidateAgainstRelationship(newNeutralRelationshipForTest(t, request.SenderUserID, request.ReceiverUserID, request.SenderUserID), true); err == nil {
		t.Fatalf("expected duplicate pending request to be rejected")
	}
}

func TestValidateAgainstRelationship_CeasefireRequiresWar(t *testing.T) {
	userA := uuid.New()
	userB := uuid.New()
	request, err := NewDiplomaticRequest(userA, userB, nil, nil, DiplomaticRequestKindCeasefireProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)

	if err := request.ValidateAgainstRelationship(rel, false); err == nil {
		t.Fatalf("expected ceasefire request to require war relationship")
	}
}

func TestAccept_CoalitionProposal_FormsAlliance(t *testing.T) {
	SetTestNow(t, 11_000)
	userA := uuid.New()
	userB := uuid.New()
	request, err := NewDiplomaticRequest(userA, userB, nil, nil, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)

	content, err := request.Accept(userB, rel)
	if err != nil {
		t.Fatalf("expected request acceptance to succeed, got %v", err)
	}
	if content != DiplomaticMessageContentCoalitionAcceptance {
		t.Fatalf("unexpected acceptance content: %s", content)
	}
	if request.Status != DiplomaticRequestStatusAccepted {
		t.Fatalf("expected accepted status, got %s", request.Status)
	}
	if request.ResolvedAt == nil || *request.ResolvedAt != 11_000 {
		t.Fatalf("expected resolved_at to be set to 11000")
	}
	if rel.Status != DiplomaticStatusAllied {
		t.Fatalf("expected relationship to become allied, got %s", rel.Status)
	}
}

func TestAccept_RejectsWrongReceiver(t *testing.T) {
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	rel := newNeutralRelationshipForTest(t, request.SenderUserID, request.ReceiverUserID, request.SenderUserID)

	if _, err := request.Accept(uuid.New(), rel); err == nil {
		t.Fatalf("expected wrong receiver acceptance to fail")
	}
}

func TestReject_SetsRejectedStatusAndMessageContent(t *testing.T) {
	SetTestNow(t, 12_000)
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCeasefireProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}

	content, err := request.Reject(request.ReceiverUserID)
	if err != nil {
		t.Fatalf("expected request rejection to succeed, got %v", err)
	}
	if content != DiplomaticMessageContentCeasefireRejection {
		t.Fatalf("unexpected rejection content: %s", content)
	}
	if request.Status != DiplomaticRequestStatusRejected {
		t.Fatalf("expected rejected status, got %s", request.Status)
	}
	if request.ResolvedAt == nil || *request.ResolvedAt != 12_000 {
		t.Fatalf("expected resolved_at to be set to 12000")
	}
}

func TestReject_ExpiredRequestFails(t *testing.T) {
	SetTestNow(t, 13_000)
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	SetTestNow(t, request.ExpiresAt+1)

	if _, err := request.Reject(request.ReceiverUserID); err == nil {
		t.Fatalf("expected expired request rejection to fail")
	}
}

func TestCanExpire_PendingDueRequestReturnsTrue(t *testing.T) {
	SetTestNow(t, 14_000)
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCoalitionProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	SetTestNow(t, request.ExpiresAt)

	if !request.CanExpire() {
		t.Fatalf("expected request to be expirable at expires_at")
	}
}

func TestExpire_SetsExpiredStatusAndResolvedAt(t *testing.T) {
	SetTestNow(t, 15_000)
	request, err := NewDiplomaticRequest(uuid.New(), uuid.New(), nil, nil, DiplomaticRequestKindCeasefireProposal)
	if err != nil {
		t.Fatalf("unexpected request init error: %v", err)
	}
	SetTestNow(t, request.ExpiresAt)

	request.Expire()

	if request.Status != DiplomaticRequestStatusExpired {
		t.Fatalf("expected expired status, got %s", request.Status)
	}
	if request.ResolvedAt == nil || *request.ResolvedAt != request.ExpiresAt {
		t.Fatalf("expected resolved_at to equal expires_at")
	}
}
