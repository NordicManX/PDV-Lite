package handlers

import (
	"database/sql"
	"net/http"

	"github.com/nordicmanx/pdv-lite/models"

	"github.com/gin-gonic/gin"
)

func CreateProductHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var novoProduto models.Produto

		if err := c.ShouldBindJSON(&novoProduto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}

		sql := `INSERT INTO produtos (codigo_produto, nome, preco_venda, quantidade_estoque) VALUES (?,?,?,?)`

		result, err := db.Exec(sql, novoProduto.CodigoProduto, novoProduto.Nome, novoProduto.PrecoVenda, novoProduto.QuantidadeEstoque)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir produto: " + err.Error()})
			return
		}

		id, _ := result.LastInsertId()
		novoProduto.ID = int(id)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Produto cadastrado com sucesso!",
			"data":    novoProduto,
		})
	}
}

func GetProductsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sql := `SELECT id, codigo_produto, nome, preco_venda, quantidade_estoque FROM produtos`

		rows, err := db.Query(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar produtos: " + err.Error()})
			return
		}
		defer rows.Close()

		var produtos []models.Produto
		for rows.Next() {
			var p models.Produto
			if err := rows.Scan(&p.ID, &p.CodigoProduto, &p.Nome, &p.PrecoVenda, &p.QuantidadeEstoque); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear produto: " + err.Error()})
				return
			}
			produtos = append(produtos, p)
		}

		c.JSON(http.StatusOK, produtos)
	}
}

func UpdateProductHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Pega o ID do produto a partir do parâmetro da URL
		id := c.Param("id")

		var produtoAtualizado models.Produto
		if err := c.ShouldBindJSON(&produtoAtualizado); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}

		sql := `UPDATE produtos SET 
					codigo_produto = ?, 
					nome = ?, 
					preco_venda = ?, 
					quantidade_estoque = ? 
				WHERE id = ?`

		result, err := db.Exec(sql, produtoAtualizado.CodigoProduto, produtoAtualizado.Nome, produtoAtualizado.PrecoVenda, produtoAtualizado.QuantidadeEstoque, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar produto: " + err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum produto encontrado com este ID para atualizar."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Produto atualizado com sucesso!"})
	}
}

func DeleteProductHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		sql := `DELETE FROM produtos WHERE id = ?`

		result, err := db.Exec(sql, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar produto: " + err.Error()})
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum produto encontrado com este ID para deletar."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso!"})
	}
}
