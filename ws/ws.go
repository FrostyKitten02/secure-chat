package ws

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"secure-chat/repo"
	"secure-chat/service"
	"strings"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // adjust in prod
}

var wsAuthHeader = "Authorization"

type client struct {
	UserID uuid.UUID
	Conn   *websocket.Conn
}

var (
	clients = make(map[string]*client)
	mu      sync.RWMutex
)

type WsSendNewMessage struct {
	ToUserID           uuid.UUID `json:"toUserId"`
	CipherText         string    `json:"cipherText"`
	Nonce              string    `json:"nonce"`
	SenderIdentityId   string    `json:"senderIdentityId"`   //used to know how to get the key!!
	ReceiverIdentityId string    `json:"ReceiverIdentityId"` //same here, these should always be the last/newest identity
}

type WsNewMessageRecieved struct {
	FromUserID         uuid.UUID `json:"fromUserId"`
	CipherText         string    `json:"cipherText"`
	Nonce              string    `json:"nonce"`
	SenderIdentityId   string    `json:"senderIdentityId"`
	ReceiverIdentityId string    `json:"ReceiverIdentityId"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade websocket", "error", err)
		return
	}
	defer conn.Close()
	//FIXME: this is temp header for testing without auth!
	jwtHeaderVal := r.Header[wsAuthHeader]
	if len(jwtHeaderVal) == 0 || len(jwtHeaderVal) > 1 {
		return
	}

	//FIXME: currently jwt verification only happens when client connects, we have to check it on each action we preform, to validate the session is still valid
	jwt := jwtHeaderVal[0]
	claims, jwtErr := service.VerifyJWT(strings.Replace(jwt, "Bearer ", "", 1))
	if jwtErr != nil {
		//TODO: tell client jwt is invalid
		return
	}
	usrId := claims.UserID
	userUUID, uuidErr := uuid.Parse(usrId)
	if uuidErr != nil {
		//TODO: tell user jwt (user-id) is invalid!!!
		return
	}

	//closing old connections on reconnect
	removeClient(usrId)

	mu.Lock()
	clients[usrId] = &client{
		UserID: userUUID,
		Conn:   conn,
	}
	mu.Unlock()

	for {
		msgType, msg, msgReadErr := conn.ReadMessage()
		if msgReadErr != nil {
			removeClient(usrId)
			break // client disconnected
		}

		if msgType == websocket.TextMessage {
			//TODO: currently sending a message is the only request client can make to WS, add message type to generic struct so we know what action client wants to preform
			var sendMsg WsSendNewMessage
			parseError := json.Unmarshal(msg, &sendMsg)
			if parseError != nil {
				//TODO notify client of error
				continue
			}

			//TODO: verify everything is set!!!
			message := WsNewMessageRecieved{
				CipherText:         sendMsg.CipherText,
				FromUserID:         userUUID,
				Nonce:              sendMsg.Nonce,
				SenderIdentityId:   sendMsg.SenderIdentityId,
				ReceiverIdentityId: sendMsg.ReceiverIdentityId,
			}

			sendMessageToClient(sendMsg.ToUserID, message)
			continue
		}

		if msgType == websocket.CloseMessage {
			removeClient(usrId)
			return
		}

		unsupportedErr := conn.WriteMessage(websocket.TextMessage, []byte("UNSUPPORTED-MSG-TYPE"))
		if unsupportedErr != nil {
			removeClient(usrId)
			return
		}
	}
}

// TODO: store message in db!!!
func sendMessageToClient(recipient uuid.UUID, message WsNewMessageRecieved) {
	go func(sender string, receiver uuid.UUID) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := repo.CreateChatIfNotExists(ctx, sender, receiver.String()); err != nil {
			slog.Error(
				"failed to create chat",
				"sender", sender,
				"receiver", receiver.String(),
				"error", err.Error(),
			)
		}
	}(message.FromUserID.String(), recipient)

	mu.RLock()
	c, ok := clients[recipient.String()]
	mu.RUnlock()

	if !ok {
		return
	}

	parsedMsg, parsedMsgErr := json.Marshal(message)
	if parsedMsgErr != nil {
		slog.Error("failed to marshal message", "error", parsedMsgErr.Error())
		//FIXME: tell client about the error!!1
		return
	}
	writeErr := c.Conn.WriteMessage(websocket.TextMessage, parsedMsg)
	if writeErr != nil {
		slog.Error("failed to write message", "error", writeErr.Error())
		//FIXME: can we tell the client that it failed???
	}
}

func removeClient(usrId string) {
	mu.Lock()
	defer mu.Unlock()

	if old, ok := clients[usrId]; ok {
		_ = old.Conn.Close()
		delete(clients, usrId)
	}
}
