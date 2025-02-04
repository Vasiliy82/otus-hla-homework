package wsreverseproxy

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/Vasiliy82/otus-hla-homework/backend/internal/domain"
	"github.com/Vasiliy82/otus-hla-homework/backend/internal/observability/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Если нужны дополнительные проверки, добавьте их здесь
	},
}

func WebSocketHandler(internalServerURL string, jwtService domain.JWTService) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := logger.Logger()
		// 1. Извлечение и проверка токена
		protocols := websocket.Subprotocols(c.Request())
		if len(protocols) == 0 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Sec-WebSocket-Protocol")
		}
		tokenStr := protocols[0]
		hash := md5.Sum([]byte(tokenStr))
		hashStr := hex.EncodeToString(hash[:])

		// Валидация токена
		token, err := jwtService.ValidateToken(domain.TokenString(tokenStr))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		userID, err := token.Claims.GetSubject()
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token, subject was not set")
		}

		// 2. Убираем заголовки авторизации
		c.Request().Header.Del("Authorization")

		// 3. Пробрасываем соединение на внутренний сервер
		proxyURL := fmt.Sprintf("%s?user_id=%s&session_id=%s", internalServerURL, userID, hashStr)
		internalConn, _, err := websocket.DefaultDialer.Dial(proxyURL, nil)
		if err != nil {
			log.Errorf("Failed to connect to internal WebSocket server: %v", err)
			return echo.NewHTTPError(http.StatusBadGateway, "Failed to connect to internal server")
		}
		defer internalConn.Close()

		clientResponse := http.Header{}
		// Устанавливаем Subprotocol (передаем обратно клиенту)
		clientResponse.Set("Sec-WebSocket-Protocol", tokenStr)

		// Устанавливаем соединение с клиентом
		clientConn, err := upgrader.Upgrade(c.Response(), c.Request(), clientResponse)
		if err != nil {
			log.Errorf("Failed to upgrade connection: %v", err)
			return err
		}
		defer clientConn.Close()

		c.Set("ws", true) // Помечаем соединение как WebSocket, чтобы Echo не вмешивался

		// 4. Прокси данных между клиентом и внутренним сервером
		errChan := make(chan error, 2)

		// Клиент -> Внутренний сервер
		go func() {
			for {
				messageType, message, err := clientConn.ReadMessage()
				if err != nil {
					log.Warnf("clientConn.ReadMessage() returned error", "err", err)
					errChan <- err
					return
				}
				if err := internalConn.WriteMessage(messageType, message); err != nil {
					log.Warnf("internalConn.WriteMessage() returned error", "err", err)
					errChan <- err
					return
				}
			}
		}()

		// Внутренний сервер -> Клиент
		go func() {
			for {
				messageType, message, err := internalConn.ReadMessage()
				if err != nil {
					log.Warnf("internalConn.ReadMessage() returned error", "err", err)
					errChan <- err
					return
				}
				if err := clientConn.WriteMessage(messageType, message); err != nil {
					log.Warnf("clientConn.WriteMessage() returned error", "err", err)
					errChan <- err
					return
				}
			}
		}()

		// Ждем завершения работы
		<-errChan
		// но **НЕ ВОЗВРАЩАЕМ** ошибку в Echo! Иначе будет echo: http: response.WriteHeader (или response.Write) on hijacked connection
		return nil
	}
}
