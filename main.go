package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"log"
	"net/http"
	"os"
)

type PromptRequest struct {
	Message string `json:"message" binding:"required"`
}

func main() {
	r := gin.Default()
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		log.Fatal("API Key não definida!")
	}

	llm, err := openai.New(openai.WithModel("gpt-4o-mini"))
	if err != nil {
		log.Fatalf("Erro ao inicializar LLM: %v", err)
	}

	r.POST("/prompt", func(c *gin.Context) {
		var req PromptRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Mensagem inválida ou ausente",
			})
			return
		}

		completion, err := llms.GenerateFromSinglePrompt(ctx, llm, req.Message)
		if err != nil {
			log.Printf("Erro ao gerar resposta: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao gerar resposta do modelo",
			})
			return
		}

		fmt.Println("Resposta gerada:", completion)

		c.JSON(http.StatusOK, gin.H{
			"response": completion,
		})
	})

	if err := r.Run(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
