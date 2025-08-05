package main

import (
	"log"
	"net/http"

	"github.com/nordicmanx/pdv-lite/handlers" // <-- 1. IMPORTANTE: Importa o outro pacote

	"github.com/nordicmanx/pdv-lite/database" // <-- 1. IMPORTANTE: Importa o pacote que criamos

	"github.com/gin-gonic/gin"
)

func main() {

	db := database.SetupDatabase()

	defer db.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Servidor do PDV LITE em Go está no ar!")
	})

	router.POST("/produtos", handlers.CreateProductHandler(db))
	router.GET("/produtos", handlers.GetProductsHandler(db))
	router.PUT("/produtos/:id", handlers.UpdateProductHandler(db))
	router.DELETE("/produtos/:id", handlers.DeleteProductHandler(db))

	router.POST("/vendas", handlers.CreateSaleHandler(db))
	router.GET("/vendas", handlers.GetSalesHandler(db))
	router.POST("/vendas/:id/cancelar", handlers.CancelSaleHandler(db))

	log.Println("Servidor rodando em http://localhost:3000")
	router.Run(":3000")
}
