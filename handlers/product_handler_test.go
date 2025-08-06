package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nordicmanx/pdv-lite/database"
	"github.com/nordicmanx/pdv-lite/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetProductsHandler testa a rota GET /produtos
func TestGetProductsHandler(t *testing.T) {
	// Configuração do teste
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Usamos um banco de dados em memória para um teste limpo e isolado
	// A string ":memory:" é um comando especial do SQLite
	db := database.SetupDatabase(":memory:")
	defer db.Close()

	// Inserimos um produto de teste diretamente no banco para o nosso teste
	_, err := db.Exec(`INSERT INTO produtos (codigo_produto, nome, preco_venda, quantidade_estoque) 
             VALUES ('TEST-001', 'Produto de Teste', 10.50, 100)`)
	if err != nil {
		t.Fatalf("Erro ao inserir produto de teste: %v", err)
	}

	// Registra a rota que queremos testar
	router.GET("/produtos", GetProductsHandler(db))

	// Cria uma requisição HTTP de teste para a nossa rota
	req, _ := http.NewRequest(http.MethodGet, "/produtos", nil)

	// 'httptest.NewRecorder()' é um gravador de resposta, ele captura o que a API responde
	rr := httptest.NewRecorder()

	// Executa a requisição
	router.ServeHTTP(rr, req)

	// --- Asserções (Verificações) ---
	// 1. Verificamos se o código de status da resposta é 200 (OK)
	assert.Equal(t, http.StatusOK, rr.Code)

	// 2. Verificamos o corpo da resposta
	var produtos []models.Produto
	err = json.Unmarshal(rr.Body.Bytes(), &produtos)

	// Verificamos que não houve erro ao decodificar o JSON
	assert.Nil(t, err)
	// Verificamos se recebemos exatamente 1 produto na resposta
	assert.Equal(t, 1, len(produtos))
	// Verificamos se o nome do produto na resposta é o que esperamos
	assert.Equal(t, "Produto de Teste", produtos[0].Nome)
}
