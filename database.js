// Importa o driver do sqlite3
const sqlite3 = require('sqlite3').verbose();

// Define o nome do arquivo do banco de dados
const DBSOURCE = "db.sqlite";

// Cria e/ou conecta ao banco de dados
let db = new sqlite3.Database(DBSOURCE, (err) => {
    if (err) {
      // Caso haja erro ao conectar
      console.error(err.message);
      throw err;
    } else {
        console.log('Conectado ao banco de dados SQLite.');
        // .serialize garante que os comandos sejam executados em ordem
        db.serialize(() => {
            // Comando para criar a tabela de usuários
            db.run(`CREATE TABLE IF NOT EXISTS usuarios (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                nome_usuario TEXT NOT NULL,
                email TEXT UNIQUE NOT NULL,
                senha TEXT NOT NULL,
                cnpj_cpf TEXT
            )`, (err) => {
                if (err) {
                    console.error("Erro ao criar tabela usuarios:", err.message);
                } else {
                    console.log("Tabela 'usuarios' criada ou já existente.");
                }
            });

            // Comando para criar a tabela de produtos
            db.run(`CREATE TABLE IF NOT EXISTS produtos (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                codigo_produto TEXT UNIQUE,
                nome TEXT NOT NULL,
                preco_venda REAL NOT NULL,
                quantidade_estoque INTEGER NOT NULL DEFAULT 0
            )`, (err) => {
                if (err) {
                    console.error("Erro ao criar tabela produtos:", err.message);
                } else {
                    console.log("Tabela 'produtos' criada ou já existente.");
                }
            });

            // Comando para criar a tabela de vendas
            db.run(`CREATE TABLE IF NOT EXISTS vendas (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                usuario_id INTEGER,
                cliente_cpf_cnpj TEXT,
                valor_total REAL NOT NULL,
                desconto REAL DEFAULT 0,
                data_venda TEXT NOT NULL,
                FOREIGN KEY (usuario_id) REFERENCES usuarios (id)
            )`, (err) => {
                if (err) {
                    console.error("Erro ao criar tabela vendas:", err.message);
                } else {
                    console.log("Tabela 'vendas' criada ou já existente.");
                }
            });

            // Comando para criar a tabela de itens da venda
            db.run(`CREATE TABLE IF NOT EXISTS venda_itens (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                venda_id INTEGER NOT NULL,
                produto_id INTEGER NOT NULL,
                quantidade INTEGER NOT NULL,
                preco_unitario_na_venda REAL NOT NULL,
                FOREIGN KEY (venda_id) REFERENCES vendas (id),
                FOREIGN KEY (produto_id) REFERENCES produtos (id)
            )`, (err) => {
                if (err) {
                    console.error("Erro ao criar tabela venda_itens:", err.message);
                } else {
                    console.log("Tabela 'venda_itens' criada ou já existente.");
                }
            });
        });
    }
});

// Exporta o objeto do banco de dados para ser usado em outras partes do projeto
module.exports = db;