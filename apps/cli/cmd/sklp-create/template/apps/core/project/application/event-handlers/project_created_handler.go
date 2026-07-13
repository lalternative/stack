// Package eventhandlers holds this context's integration-event handlers. Each
// handler implements consumer.EventHandler and writes only business logic in
// Handle; all JetStream redelivery semantics (Term on permanent errors, bounded
// MaxDeliver, staged BackOff, the DLQ stream, ack heartbeats, the reconnect
// loop) are provided by github.com/lalternative/packages/go/eda/pkg/consumer.
//
// Do NOT hand-roll a JetStream subscription here. Copy this file, change the
// Subject/DurableName and the body of Handle — nothing else.
package eventhandlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lalternative/packages/go/eda/pkg/consumer"
	"github.com/nats-io/nats.go"
)

// ProjectCreated is the payload this example handler consumes.
type ProjectCreated struct {
	EventID string `json:"event_id"`
	ID      string `json:"id"`
	Name    string `json:"name"`
}

// ProjectCreatedHandler is a sample durable consumer. Delete it once you wire
// your own, or keep it as the template for new ones.
type ProjectCreatedHandler struct{}

func NewProjectCreatedHandler() *ProjectCreatedHandler { return &ProjectCreatedHandler{} }

func (*ProjectCreatedHandler) Name() string        { return "project-created" }
func (*ProjectCreatedHandler) Subject() string     { return "integration.project.created" }
func (*ProjectCreatedHandler) DurableName() string { return "project-created-consumer" }
func (*ProjectCreatedHandler) MaxDeliver() int     { return 3 }

func (*ProjectCreatedHandler) Handle(ctx context.Context, msg *nats.Msg) error {
	var ev ProjectCreated
	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		// Malformed payload: no retry can fix it -> dead-letter immediately.
		return consumer.Permanent(fmt.Errorf("decode project.created: %w", err))
	}

	// ... your business logic. Returning a plain (non-permanent) error here
	// triggers a bounded, backed-off retry; after MaxDeliver the message is
	// routed to the DLQ stream automatically.
	_ = ev
	return nil
}
