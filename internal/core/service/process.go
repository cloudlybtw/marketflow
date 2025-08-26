package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"marketflow/internal/core/domain"
	"time"
)

func Run(ctx context.Context, workerID int, in <-chan domain.Message, out chan<- domain.Message) {
	slog.Info(fmt.Sprintf("Processor %d started", workerID))
	defer slog.Info(fmt.Sprintf("Processor %d stopped", workerID))

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-in:
			if !ok {
				return
			}

			time.Sleep(time.Millisecond * 10)

			var processedData map[string]interface{}
			if err := json.Unmarshal(msg.Raw, &processedData); err == nil {
				processedData["processed_by"] = workerID
				if processedBytes, err := json.Marshal(processedData); err == nil {
					msg.Raw = processedBytes
				}
			}

			select {
			case out <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}
