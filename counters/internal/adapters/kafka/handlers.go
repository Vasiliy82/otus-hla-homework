package kafka

import (
	"context"
	"encoding/json"

	"github.com/Vasiliy82/otus-hla-homework/counters/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/counters/internal/usecases"
)

type KafkaEventHandler struct {
	processEventUC *usecases.ProcessSagaEventUseCase
}

func NewKafkaEventHandler(processEventUC *usecases.ProcessSagaEventUseCase) *KafkaEventHandler {
	return &KafkaEventHandler{processEventUC: processEventUC}
}

func (h *KafkaEventHandler) HandleMessage(ctx context.Context, msg []byte) error {
	var event domain.SagaEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}
	return h.processEventUC.Execute(ctx, event)
}
