package admin

import (
	"context"
	"fmt"
)

type Auditor interface {
	Log(ctx context.Context, actorID, action, resourceType, resourceID, details string) error
}

func (s *Service) Log(ctx context.Context, actorID, action, resourceType, resourceID, details string) error {
	entry := &AuditLogEntry{
		ActorID:      actorID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
	}
	if err := s.CreateAuditLog(ctx, entry); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	return nil
}
