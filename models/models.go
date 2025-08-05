package models

type ItemVendaDetalhado struct {
	NomeProduto   string  `json:"nome_produto"`
	Quantidade    int     `json:"quantidade"`
	PrecoUnitario float64 `json:"preco_unitario"`
}

type VendaDetalhada struct {
	ID             int64                `json:"id"`
	ClienteCPFCNPJ string               `json:"cliente_cpf_cnpj"`
	ValorTotal     float64              `json:"valor_total"`
	DataVenda      string               `json:"data_venda"`
	Itens          []ItemVendaDetalhado `json:"itens"`
}

type ItemVenda struct {
	ProdutoID  int `json:"produto_id"`
	Quantidade int `json:"quantidade"`
}

type Venda struct {
	ClienteCPFCNPJ string      `json:"cliente_cpf_cnpj"`
	Itens          []ItemVenda `json:"itens"`
}

type Produto struct {
	ID                int     `json:"id"`
	CodigoProduto     string  `json:"codigo_produto"`
	Nome              string  `json:"nome"`
	PrecoVenda        float64 `json:"preco_venda"`
	QuantidadeEstoque int     `json:"quantidade_estoque"`
}
