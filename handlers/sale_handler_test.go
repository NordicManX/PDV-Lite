package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nordicmanx/pdv-lite/database"
	"github.com/nordicmanx/pdv-lite/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSaleHandler(t *testing.T) {
	// --- 1. Preparar (Arrange) ---
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db := database.SetupDatabase(":memory:")
	defer db.Close()

	// Pré-populamos o banco com produtos que podemos vender no teste
	_, err := db.Exec(`INSERT INTO produtos (id, nome, preco_venda, quantidade_estoque) VALUES 
		(1, 'Produto A', 10.00, 20),
		(2, 'Produto B', 5.00, 100)`)
	assert.NoError(t, err) // Garante que a inserção dos produtos de teste deu certo

	// Registramos a rota que vamos testar
	router.POST("/vendas", CreateSaleHandler(db))

	// --- 2. Agir (Act) ---

	// Montamos o corpo (payload) da nossa requisição de venda
	// Vamos vender 1 unidade do Produto A e 2 unidades do Produto B
	vendaPayload := models.Venda{
		ClienteCPFCNPJ: "111.222.333-44",
		Itens: []models.ItemVenda{
			{ProdutoID: 1, Quantidade: 1},
			{ProdutoID: 2, Quantidade: 2},
		},
	}
	// Convertemos nosso payload para JSON
	payload, _ := json.Marshal(vendaPayload)

	// Criamos a requisição HTTP
	req, _ := http.NewRequest(http.MethodPost, "/vendas", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	// Criamos um gravador de resposta
	rr := httptest.NewRecorder()

	// Executamos a requisição
	router.ServeHTTP(rr, req)

	// --- 3. Verificar (Assert) ---

	// Verificamos se o status code da API foi 201 (Created)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// A verificação mais importante: vamos ao banco de dados e checamos se o estoque foi atualizado!
	var novoEstoqueProdutoB int
	err = db.QueryRow("SELECT quantidade_estoque FROM produtos WHERE id = 2").Scan(&novoEstoqueProdutoB)

	// Garantimos que a consulta ao banco deu certo
	assert.NoError(t, err)

	// O estoque do Produto B era 100. Vendemos 2. O novo estoque DEVE ser 98.
	assert.Equal(t, 98, novoEstoqueProdutoB)
}
