// pdv-lite/main.go

package main

import (
	"log"
	"net/http"

	"github.com/nordicmanx/pdv-lite/database"
	"github.com/nordicmanx/pdv-lite/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Mudança aqui para passar o nome do arquivo do banco de dados principal
	db := database.SetupDatabase("./db.sqlite")
	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Servidor do PDV LITE em Go está no ar!")
	})

	// --- ROTAS DE PRODUTOS ---
	router.POST("/produtos", handlers.CreateProductHandler(db))
	router.GET("/produtos", handlers.GetProductsHandler(db))
	router.PUT("/produtos/:id", handlers.UpdateProductHandler(db))
	router.DELETE("/produtos/:id", handlers.DeleteProductHandler(db))

	// --- ROTAS DE VENDAS ---
	router.POST("/vendas", handlers.CreateSaleHandler(db))
	router.GET("/vendas", handlers.GetSalesHandler(db))
	router.POST("/vendas/:id/cancelar", handlers.CancelSaleHandler(db))

	log.Println("Servidor rodando em http://localhost:3000")
	router.Run(":3000")
}
