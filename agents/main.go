package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"support-ops-agents/llm"
	"support-ops-agents/web"
	"support-ops-agents/ws"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := ws.New()
	mux := gin.Default()
	messagesStore := make(llm.Messages, 0, 200)
	messageChan := make(chan llm.Message, 100)
	mux.GET("/api/v1/ws", func(ctx *gin.Context) {
		conn, err := ws.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.Status(http.StatusForbidden)
			return
		}
		hub.RegisterClient(ctx.Query("id"), conn)

	})
	mux.POST("/api/v1/message", func(ctx *gin.Context) {
		var message struct {
			TicketId string `json:"ticket_id"`
			Message  string `json:"message"`
		}
		if err := ctx.ShouldBindBodyWithJSON(&message); err != nil {
			ctx.JSON(http.StatusBadRequest, web.MESSAGE_BAD_REQUEST.Error())
		}
		go func() {
			tools := llm.NewTools(message.TicketId)
			messageChan <- llm.Message{
				Role:    "user",
				Content: message.Message,
				Tools:   &tools,
			}
		}()

		ctx.Status(http.StatusCreated)
	})
	mux.POST("/api/v1/tool/reply", func(ctx *gin.Context) {
		var message llm.ToolReply
		if err := ctx.ShouldBindBodyWithJSON(&message); err != nil {
			ctx.JSON(http.StatusBadRequest, web.MESSAGE_BAD_REQUEST.Error())
			return
		}
		go func() {
			messageChan <- llm.Message{
				Content: message.Content,
				Id:      &message.Id,
				Role:    "tool",
			}
		}()

		ctx.Status(http.StatusCreated)
	})
	srv := &http.Server{
		Handler: mux,
		Addr:    ":9038",
	}
	go func() {
		log.Println("starting server")
		log.Println("Running on port 9038")
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	log.Println("Initiating newTools()")

	go func() {
		for message := range messageChan {
			if message.Tools == nil {
				continue
			}
			tools := *message.Tools
			messagesStore = append(messagesStore, message)
			response, err := llm.Do(messagesStore, tools)
			if err != nil {
				panic(err)
			}
			if response.ToolName != "" {
				result, err := tools.Exec(response.ToolCallId, response.ToolName, response.Arguments)
				if err != nil {
					log.Printf("Failed tool call name=%s arguments=%+v reason=%s", response.ToolName, response.Arguments, err.Error())
					continue
				}
				resultStr, ok := result.(string)
				if !ok {
					log.Println("Failed to cast to str")
					continue
				}
				messageChan <- llm.Message{
					Role:    "tool",
					Content: resultStr,
					Id:      &response.ToolCallId,
				}
			}

		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
}
