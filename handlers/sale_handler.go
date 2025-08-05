package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/nordicmanx/pdv-lite/models"

	"github.com/gin-gonic/gin"
)

func CreateSaleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var novaVenda models.Venda

		if err := c.ShouldBindJSON(&novaVenda); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON da venda inválido: " + err.Error()})
			return
		}

		if len(novaVenda.Itens) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A venda deve ter pelo menos um item."})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação: " + err.Error()})
			return
		}

		var valorTotalVenda float64 = 0
		for _, item := range novaVenda.Itens {
			var precoVenda float64
			var estoqueAtual int

			err := tx.QueryRow("SELECT preco_venda, quantidade_estoque FROM produtos WHERE id = ?", item.ProdutoID).Scan(&precoVenda, &estoqueAtual)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Produto com ID %d não encontrado.", item.ProdutoID)})
				return
			}

			if estoqueAtual < item.Quantidade {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Estoque insuficiente para o produto ID %d. Estoque atual: %d, Tentativa de venda: %d", item.ProdutoID, estoqueAtual, item.Quantidade)})
				return
			}

			valorTotalVenda += precoVenda * float64(item.Quantidade)
		}

		dataVenda := time.Now().Format("2006-01-02 15:04:05")
		resVenda, err := tx.Exec("INSERT INTO vendas (cliente_cpf_cnpj, valor_total, data_venda) VALUES (?, ?, ?)",
			novaVenda.ClienteCPFCNPJ, valorTotalVenda, dataVenda)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar registro de venda: " + err.Error()})
			return
		}
		vendaID, _ := resVenda.LastInsertId()

		for _, item := range novaVenda.Itens {
			var precoVenda float64
			tx.QueryRow("SELECT preco_venda FROM produtos WHERE id = ?", item.ProdutoID).Scan(&precoVenda)

			_, err := tx.Exec("INSERT INTO venda_itens (venda_id, produto_id, quantidade, preco_unitario_na_venda) VALUES (?, ?, ?, ?)",
				vendaID, item.ProdutoID, item.Quantidade, precoVenda)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao inserir item da venda: " + err.Error()})
				return
			}

			_, err = tx.Exec("UPDATE produtos SET quantidade_estoque = quantidade_estoque - ? WHERE id = ?",
				item.Quantidade, item.ProdutoID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar estoque: " + err.Error()})
				return
			}
		}

		tx.Commit()

		c.JSON(http.StatusCreated, gin.H{"message": "Venda registrada com sucesso!", "venda_id": vendaID})
	}
}

func GetSalesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rowsVendas, err := db.Query("SELECT id, cliente_cpf_cnpj, valor_total, data_venda FROM vendas ORDER BY id DESC")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar vendas: " + err.Error()})
			return
		}
		defer rowsVendas.Close()

		var vendasDetalhadas []models.VendaDetalhada

		for rowsVendas.Next() {
			var vd models.VendaDetalhada
			if err := rowsVendas.Scan(&vd.ID, &vd.ClienteCPFCNPJ, &vd.ValorTotal, &vd.DataVenda); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear venda: " + err.Error()})
				return
			}

			rowsItens, err := db.Query(`
				SELECT p.nome, vi.quantidade, vi.preco_unitario_na_venda 
				FROM venda_itens vi 
				JOIN produtos p ON vi.produto_id = p.id 
				WHERE vi.venda_id = ?`, vd.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens da venda: " + err.Error()})
				return
			}

			var itens []models.ItemVendaDetalhado
			for rowsItens.Next() {
				var item models.ItemVendaDetalhado
				if err := rowsItens.Scan(&item.NomeProduto, &item.Quantidade, &item.PrecoUnitario); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear item da venda: " + err.Error()})
					return
				}
				itens = append(itens, item)
			}
			rowsItens.Close()

			vd.Itens = itens
			vendasDetalhadas = append(vendasDetalhadas, vd)
		}

		c.JSON(http.StatusOK, vendasDetalhadas)
	}
}

func CancelSaleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação: " + err.Error()})
			return
		}

		// 1. Verifica o status atual da venda
		var statusAtual string
		err = tx.QueryRow("SELECT status FROM vendas WHERE id = ?", id).Scan(&statusAtual)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Venda não encontrada."})
			return
		}

		if statusAtual == "CANCELADA" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Esta venda já foi cancelada."})
			return
		}

		// 2. Busca todos os itens da venda a ser cancelada
		rowsItens, err := tx.Query("SELECT produto_id, quantidade FROM venda_itens WHERE venda_id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens da venda para cancelamento."})
			return
		}
		defer rowsItens.Close()

		// 3. Devolve os itens ao estoque
		for rowsItens.Next() {
			var item models.ItemVenda
			if err := rowsItens.Scan(&item.ProdutoID, &item.Quantidade); err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear item para cancelamento."})
				return
			}

			_, err = tx.Exec("UPDATE produtos SET quantidade_estoque = quantidade_estoque + ? WHERE id = ?",
				item.Quantidade, item.ProdutoID)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao estornar estoque do produto ID " + fmt.Sprintf("%d", item.ProdutoID)})
				return
			}
		}

		// 4. Atualiza o status da venda para "CANCELADA"
		_, err = tx.Exec("UPDATE vendas SET status = 'CANCELADA' WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status da venda."})
			return
		}

		// 5. Se tudo deu certo, confirma a transação
		tx.Commit()

		c.JSON(http.StatusOK, gin.H{"message": "Venda cancelada e estoque estornado com sucesso!"})
	}
}
