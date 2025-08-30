package source

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"marketflow/internal/adapter/config"
	"marketflow/internal/core/domain"
	"net"
	"time"
)

type TCPExchange struct {
	Conn net.Conn
}

func NewTCPExchange(ctx context.Context, exchange config.Exchange, out chan<- domain.Message) *TCPExchange {
	reconnectDelay := time.Second * 2
	firstConnect := true
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	var conn net.Conn
	var err error
	for {
		select {
		case <-ctx.Done():
			slog.Info(fmt.Sprintf("%s: shutting down connection manager", exchange.Name))
			return nil
		default:
			if !firstConnect {
				slog.Info(fmt.Sprintf("%s: attempting to reconnect in %v...", exchange.Name, reconnectDelay))
				time.Sleep(reconnectDelay)
			}
			firstConnect = false

			conn, err = dialer.DialContext(ctx, "tcp", net.JoinHostPort(exchange.Name, exchange.Port))
			if err != nil {
				slog.Info(fmt.Sprintf("%s: connection failed: %v", exchange.Name, err))
				continue
			}

			slog.Info(fmt.Sprintf("%s: connected successfully", exchange.Name))
			return &TCPExchange{
				Conn: conn,
			}
		}
	}
}

func (t *TCPExchange) ReadData(ctx context.Context, exchange config.Exchange, out chan<- domain.Message) error {
	scanner := bufio.NewScanner(t.Conn)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !scanner.Scan() {
				return scanner.Err()
			}

			msg := domain.Message{
				Exchange:  exchange.Name,
				Timestamp: time.Now(),
				Raw:       json.RawMessage(scanner.Bytes()),
			}

			select {
			case out <- msg:
			case <-ctx.Done():
				return nil
			default:
				slog.Info(fmt.Sprintf("%s: message channel full, dropping message", exchange.Name))
			}
		}
	}
}
