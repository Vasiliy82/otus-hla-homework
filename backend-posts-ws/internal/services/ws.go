package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/Vasiliy82/otus-hla-homework/backend-posts-ws/internal/config"
	"github.com/Vasiliy82/otus-hla-homework/backend-posts-ws/internal/observability/logger"
)

// WSMessage — общее сообщение в формате JSON:
//
//	{
//	   "method": "NewFeedRec",
//	   "payload": {...}
//	}
type WSMessage struct {
	Method  string      `json:"method"`
	Payload interface{} `json:"payload"`
}

// Client хранит данные об одном конкретном WebSocket-соединении
type Client struct {
	userID    string
	sessionID string

	conn   *websocket.Conn
	sendCh chan []byte
	doneCh chan struct{}

	server *WSService // ссылка на сервер, если нужно

	closeOnce sync.Once // добавляем Once для безопасного закрытия
}

// WSService хранит все активные соединения (userID -> sessionID -> Client)
type WSService struct {
	cfg      *config.WebSocketConfig
	upgrader websocket.Upgrader

	mu sync.RWMutex
	// Пример структуры: connections[userID][sessionID] = *Client
	connections map[string]map[string]*Client
}

// NewWSService конструктор
func NewWSService(cfg *config.WebSocketConfig) *WSService {
	return &WSService{
		cfg: cfg,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// При необходимости проверяем Origin
				return true
			},
		},
		connections: make(map[string]map[string]*Client),
	}
}

// HandleConnection обрабатывает входящие HTTP-запросы, апгрейдит до WebSocket.
func (s *WSService) HandleConnection(c echo.Context) error {
	log := logger.Logger()
	// Извлекаем параметры из query string
	userID := c.QueryParam("userId")       // GUID пользователя
	sessionID := c.QueryParam("sessionId") // GUID сессии
	if userID == "" || sessionID == "" {
		return c.String(http.StatusBadRequest, "Missing userId or sessionId in query params")
	}

	wsConn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Errorw("Upgrade error", "err", err)
		return err
	}

	// Создаём структуру клиента
	client := &Client{
		userID:    userID,
		sessionID: sessionID,
		conn:      wsConn,
		sendCh:    make(chan []byte, 256),
		doneCh:    make(chan struct{}),
		server:    s,
	}

	// Регистрируем клиента в карте
	s.registerClient(client)

	// Запускаем горутины чтения и записи
	go s.readLoop(client)
	go s.writeLoop(client)

	return nil
}

// регистрируем клиента в общей карте
func (s *WSService) registerClient(c *Client) {
	log := logger.Logger()
	s.mu.Lock()
	defer s.mu.Unlock()

	userMap, ok := s.connections[c.userID]
	if !ok {
		userMap = make(map[string]*Client)
		s.connections[c.userID] = userMap
	}

	// Если сессия уже была, закрываем старое соединение
	if oldClient, exists := userMap[c.sessionID]; exists {
		// Закрываем старое соединение
		s.closeClient(oldClient)
	}

	userMap[c.sessionID] = c
	log.Infow("Registered client", "c.userID", c.userID, "c.SessionID", c.sessionID)
}

// удаляем клиента из общей карты
func (s *WSService) unregisterClient(c *Client) {
	log := logger.Logger()
	s.mu.Lock()
	defer s.mu.Unlock()

	userMap, ok := s.connections[c.userID]
	if ok {
		// проверяем, что действительно этот Client сидит в userMap[c.sessionID]
		storedClient, exists := userMap[c.sessionID]
		if exists && storedClient == c {
			delete(userMap, c.sessionID)
			// если у данного пользователя нет больше сессий, можно удалить всю мапу
			if len(userMap) == 0 {
				delete(s.connections, c.userID)
			}
		}
	}
	log.Infow("Unregistered client", "c.userID", c.userID, "c.sessionID", c.sessionID)
}

// readLoop читает входящие сообщения (JSON или ping/pong) от клиента
func (s *WSService) readLoop(c *Client) {
	log := logger.Logger()
	defer func() {
		log.Debugw("debug: readLoop: finished", "c.userID", c.userID, "c.sessionID", c.sessionID)
		s.closeClient(c)
	}()
	log.Debugw("debug: readLoop: started", "c.userID", c.userID, "c.sessionID", c.sessionID)

	// Настраиваем дедлайны
	c.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
	c.conn.SetPongHandler(func(appData string) error {
		// При получении pong продляем дедлайн
		log.Debugw("c.conn.SetPongHandler pong event", "appData", appData)
		c.conn.SetReadDeadline(time.Now().Add(s.cfg.PongWait))
		return nil
	})

	// Если клиент пингует нас нестандартным pingMessage (TextMessage),
	// можно обработать это в JSON, но стандартный ping/pong уже есть в Gorilla.
	c.conn.SetPingHandler(func(appData string) error {
		// Можем что-то логировать, если нужно
		log.Debugw("c.conn.SetPingHandler ping event", "appData", appData)
		return c.conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	for {
		msgType, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Infof("readLoop error for user=%s session=%s: %v", c.userID, c.sessionID, err)
			return
		}

		log.Debugw("readLoop: message received", "msgType", msgType, "data", data)

		switch msgType {
		case websocket.TextMessage:
			// Предположим, клиент присылает JSON-команды
			var msg WSMessage
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Errorw("JSON parse error", "err", err)
				continue
			}
			s.handleInboundMessage(c, &msg)

		case websocket.CloseMessage:
			log.Infow("Client closed connection", "c.userID", c.userID, "c.sessionID", c.sessionID)
			return
		default:
			// Ping/PongMessage обрабатываются хендлерами, BinaryMessage — игнорируем
		}
	}
}

// writeLoop периодически пингует клиента и отправляет ему сообщения из sendCh
func (s *WSService) writeLoop(c *Client) {
	log := logger.Logger()
	ticker := time.NewTicker(s.cfg.PingInterval)
	defer func() {
		log.Debugw("debug: writeLoop: finished", "c.userID", c.userID, "c.sessionID", c.sessionID)

		ticker.Stop()
		s.closeClient(c)
	}()
	log.Debugw("debug: writeLoop: started", "c.userID", c.userID, "c.sessionID", c.sessionID)

	for {
		select {
		case <-c.doneCh:
			return

		case msg, ok := <-c.sendCh:
			if !ok {
				// канал закрыт
				return
			}
			// Отправляем TextMessage
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Errorw("writeLoop error", "c.userID", c.userID, "c.sessionID", c.sessionID, "err", err)
				return
			}

		case <-ticker.C:
			// Периодический пинг
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Errorw("ping error", "c.userID", c.userID, "c.sessionID", c.sessionID, "err", err)
				return
			}
			log.Debugw("WS ping was sent", "c.userID", c.userID, "c.sessionID", c.sessionID)
		}
	}
}

// handleInboundMessage обрабатывает входящие JSON-сообщения
func (s *WSService) handleInboundMessage(c *Client, msg *WSMessage) {
	log := logger.Logger()
	switch msg.Method {
	case "pingFromClient":
		// Например, клиент шлёт ping как JSON
		reply := WSMessage{
			Method:  "pongFromServer",
			Payload: "pong data",
		}
		_ = c.SendJSON(reply)

	// Другие методы — в зависимости от бизнес-логики
	default:
		log.Errorw("Unknown method", "c.userID", c.userID, "c.sessionID", c.sessionID, "msg.Method", msg.Method)
	}
}

// Закрываем соединение + unregister (с защитой от двойного вызова)
func (s *WSService) closeClient(c *Client) {
	// sync.Once гарантирует, что тело будет выполнено один раз
	c.closeOnce.Do(func() {
		s.unregisterClient(c)

		// Закрываем сетевое соединение
		_ = c.conn.Close()

		// Закрываем каналы
		close(c.doneCh)
		close(c.sendCh)
	})
}

// BroadcastNewFeedRec отправляет всем сессиям пользователя userID (если он есть) метод "NewFeedRec"
func (s *WSService) BroadcastNewFeedRec(userID string, payload interface{}) error {
	log := logger.Logger()
	// Формируем JSON вида {"method":"NewFeedRec", "payload": ...}
	msg := WSMessage{
		Method:  "NewFeedRec",
		Payload: payload,
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	userMap, ok := s.connections[userID]
	if !ok || len(userMap) == 0 {
		// Нет ни одной сессии на этом экземпляре
		// В 95+% случаев будет так
		return nil
	}

	for _, client := range userMap {
		// Посылаем сообщение
		select {
		case client.sendCh <- raw:
		default:
			// если буфер переполнен — дропаем, либо закрываем. Выбирайте стратегию
			log.Errorw("Send buffer full, dropping msg", "c.userID", client.userID, "c.sessionID", client.sessionID, "msg.Method", msg.Method)
		}
	}

	return nil
}

// Отправка JSON-сообщения конкретной сессии
func (s *WSService) SendNewFeedRecToSession(userID, sessionID string, payload interface{}) error {
	s.mu.RLock()
	userMap, ok := s.connections[userID]
	if !ok {
		s.mu.RUnlock()
		return fmt.Errorf("no sessions for user=%s on this instance", userID)
	}
	client, ok := userMap[sessionID]
	s.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no such session=%s for user=%s on this instance", sessionID, userID)
	}

	msg := WSMessage{
		Method:  "NewFeedRec",
		Payload: payload,
	}
	return client.SendJSON(msg)
}

// SendJSON — удобная обёртка для отправки JSON в канал клиента
func (c *Client) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	select {
	case c.sendCh <- data:
	default:
		// буфер переполнен
		return fmt.Errorf("send buffer is full for user=%s session=%s", c.userID, c.sessionID)
	}
	return nil
}

// -----------------------------------------------------------------------------
// Пример асинхронного чтения из Kafka в этом же сервисе
// -----------------------------------------------------------------------------

// KafkaConsumerStub — имитация kafka-консьюмера
type KafkaConsumerStub interface {
	ReadMessage(timeoutMs int) (userID string, sessionID string, data []byte, err error)
}

// Запускаем горутину, которая читает из Kafka и рассылает сообщения.
// В 95%+ случаев клиент не будет найден в этой ноде, и мы просто пропустим сообщение.
func (s *WSService) StartKafkaLoop(ctx context.Context, consumer KafkaConsumerStub) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Читаем из Kafka (упрощённо)
			userID, sessionID, rawPayload, err := consumer.ReadMessage(-1)
			if err != nil {
				// Обработка ошибок, логирование
				continue
			}

			if sessionID == "" {
				// отправляем всем сессиям
				_ = s.BroadcastNewFeedRec(userID, json.RawMessage(rawPayload))
			} else {
				// отправка в конкретную сессию
				_ = s.SendNewFeedRecToSession(userID, sessionID, json.RawMessage(rawPayload))
			}
		}
	}()
}
