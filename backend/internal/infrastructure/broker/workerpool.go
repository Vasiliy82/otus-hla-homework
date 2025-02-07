package broker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type FuncProcessor func(context.Context, *kafka.Message, int) error

// WorkerConfig для настройки воркеров
type WorkerConfig struct {
	Topic          string
	NumWorkers     int
	funcProcessor  FuncProcessor
	ConsumerConfig *kafka.ConfigMap
}

type WorkerPool struct {
	cfg    *WorkerConfig
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewWorkerConfig(topic string, numWorkers int, funcProcessor FuncProcessor, consumerConfig *kafka.ConfigMap) *WorkerConfig {
	return &WorkerConfig{
		Topic:          topic,
		NumWorkers:     numWorkers,
		funcProcessor:  funcProcessor,
		ConsumerConfig: consumerConfig,
	}

}

func NewWorker(cfg *WorkerConfig) *WorkerPool {
	return &WorkerPool{
		cfg: cfg,
	}
}

func (w *WorkerPool) Start(ctx context.Context) {
	workerCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	go func() {
		err := w.startWorker(workerCtx)
		if err != nil {
			logger.Logger().Errorf("WorkerPool.Start: Worker returned error: %v", err)
		}
	}()
}
func (w *WorkerPool) Stop() {
	w.cancel()
}
func (w *WorkerPool) Wait() {
	logger.Logger().Debugw("waiting for WorkerPool", "w.cfg", w.cfg)
	w.wg.Wait()
	logger.Logger().Debugw("done", "w.cfg", w.cfg)
}

func (w *WorkerPool) startWorker(ctx context.Context) error {
	consumer, err := kafka.NewConsumer(w.cfg.ConsumerConfig)
	if err != nil {
		return err
	}
	defer consumer.Close()

	err = consumer.Subscribe(w.cfg.Topic, nil)
	if err != nil {
		return err
	}

	w.wg.Add(w.cfg.NumWorkers)

	for i := 0; i < w.cfg.NumWorkers; i++ {
		go func(workerID int) {
			logger.Logger().Debugf("WorkerPool.startWorker: Topic %s Worker %d: starting", w.cfg.Topic, workerID)
			defer w.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Используем контекст с тайм-аутом
					msgCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
					defer cancel()

					event := consumer.Poll(100) // Таймаут ожидания событий (в миллисекундах)
					if event == nil {
						continue
					}

					done := make(chan error, 1)
					go func() {
						switch ev := event.(type) {
						case *kafka.Message:

							logger.Logger().Debugw("WorkerPool.startWorker: Message processing started  >>>", "Topic", w.cfg.Topic, "WorkerId", workerID, "message", (string)(ev.Value))
							err := w.cfg.funcProcessor(msgCtx, ev, workerID)
							if err == nil {
								consumer.CommitMessage(ev)
							}
							logger.Logger().Debugw("WorkerPool.startWorker: Message processing finished <<<", "Topic", w.cfg.Topic, "WorkerId", workerID, "message", (string)(ev.Value))
							done <- err
						case kafka.Error:
							done <- ev
						default:
							logger.Logger().Errorf("WorkerPool.startWorker: Topic %s Worker %d: Ignored event: %v", w.cfg.Topic, workerID, ev)
						}
					}()

					select {
					case err := <-done:
						if err != nil {
							logger.Logger().Errorf("WorkerPool.startWorker: Topic %s Worker %d: Processing error: %v", w.cfg.Topic, workerID, err)
						}
						var kErr kafka.Error
						if errors.As(err, &kErr) && kErr.IsFatal() {
							logger.Logger().Errorf("WorkerPool.startWorker: Topic %s Worker %d: Fatal error: %v", w.cfg.Topic, workerID, kErr)
						}
					case <-msgCtx.Done():
						logger.Logger().Errorf("WorkerPool.startWorker: Topic %s Worker %d: Processing timeout", w.cfg.Topic, workerID)
					}
					cancel()
				}
			}
		}(i)
	}
	logger.Logger().Debugw("WorkerPool.startWorker: all workers started", "w.cfg.Topic", w.cfg.Topic)
	w.wg.Wait()
	logger.Logger().Debugw("WorkerPool.startWorker: all workers gone", "w.cfg.Topic", w.cfg.Topic)

	return nil
}
