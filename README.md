ProjectOneFullCycle é um projeto em Go que implementa um servidor e um cliente para fornecer e consumir taxas de câmbio. O servidor disponibiliza os dados armazenados em um banco de dados SQLite3, enquanto o cliente faz a requisição e exibe os valores.

Estrutura do Projeto

server.go → Implementa o servidor que fornece as taxas de câmbio.

client.go → Implementa o cliente que consome e exibe as taxas.

exchange_rates.db → Banco de dados SQLite3 com as taxas de câmbio.

cotacao.txt →  Arquivo txt de armazenamento de cotações.

go.mod e go.sum → Gerenciamento de dependências do Go.

Observações

O servidor deve estar em execução antes do cliente.

