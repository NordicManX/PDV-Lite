// pdv-lite/database/database.go

package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// SetupDatabase inicializa e retorna uma conexão com o banco de dados.
// Agora ela recebe o sourceName (ex: "./db.sqlite" ou ":memory:")
func SetupDatabase(sourceName string) *sql.DB {
	// Se o banco não for em memória, garantimos que o arquivo não pré-exista em modo de teste
	// para manter os testes isolados. Esta parte é opcional mas uma boa prática.
	if sourceName != ":memory:" {
		os.Remove(sourceName) // Remove o banco de dados antigo se existir
	}

	db, err := sql.Open("sqlite3", sourceName)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	log.Println("Conectado ao banco de dados SQLite em:", sourceName)

	schema := `
    CREATE TABLE IF NOT EXISTS usuarios (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nome_usuario TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        senha TEXT NOT NULL,
        cnpj_cpf TEXT
    );
    CREATE TABLE IF NOT EXISTS produtos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        codigo_produto TEXT UNIQUE,
        nome TEXT NOT NULL,
        preco_venda REAL NOT NULL,
        quantidade_estoque INTEGER NOT NULL DEFAULT 0
    );
    CREATE TABLE IF NOT EXISTS vendas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        usuario_id INTEGER,
        cliente_cpf_cnpj TEXT,
        valor_total REAL NOT NULL,
        desconto REAL DEFAULT 0,
        data_venda TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'CONCLUIDA',
        FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
    );
    CREATE TABLE IF NOT EXISTS venda_itens (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        venda_id INTEGER NOT NULL,
        produto_id INTEGER NOT NULL,
        quantidade INTEGER NOT NULL,
        preco_unitario_na_venda REAL NOT NULL,
        FOREIGN KEY (venda_id) REFERENCES vendas (id),
        FOREIGN KEY (produto_id) REFERENCES produtos (id)
    );`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}
	log.Println("Tabelas verificadas/criadas com sucesso.")

	return db
}
