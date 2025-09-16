# Microsserviço de Produtos (Product Service)

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)

## 📖 Sobre o Projeto

Este é o **Microsserviço de Produtos**, o pilar do catálogo de um sistema de e-commerce distribuído. Desenvolvido em Go, a sua responsabilidade exclusiva é ser a fonte da verdade para todas as informações de produtos, incluindo detalhes, preço e gestão de stock.

O serviço foi projetado para ser altamente performático e escalável. Ele expõe endpoints públicos para consulta do catálogo e endpoints internos protegidos para tarefas administrativas, como a criação de produtos e a atualização de stock, que serão consumidos por outros serviços do ecossistema, como o futuro `order-service`.

### ✨ Funcionalidades Principais
* **Gestão de Catálogo:** Endpoints internos para criar e gerir produtos no inventário.
* **Consulta Pública:** Endpoints abertos para que clientes (como o frontend da loja ou o `cart-service`) possam listar produtos e ver detalhes de um item específico.
* **Controlo de Stock:** Endpoint interno dedicado para a atualização (redução) de stock, uma operação crítica para o fluxo de finalização de compra.
* **Segurança Serviço-a-Serviço:** Endpoints internos são protegidos por uma API Key, garantindo que apenas serviços autorizados possam realizar operações de escrita.
* **IDs Ordenáveis:** Utiliza **ULID** como identificador único para os produtos, garantindo unicidade e ordenação cronológica, o que otimiza consultas na base de dados.

## 🛠️ Arquitetura e Tecnologias

O projeto segue uma arquitetura em camadas (API, Lógica de Negócio, Repositório), mantendo a consistência com o `auth-service` para uma clara separação de responsabilidades.

### Tecnologias Utilizadas
* **Linguagem:** Go
* **Banco de Dados:** PostgreSQL
* **Containerização:** Docker & Docker Compose
* **Roteador HTTP:** Chi
* **Driver do Banco de Dados:** pgx
* **Migrations:** golang-migrate (a ser adicionado)
* **Automação:** Makefile (a ser adicionado, seguindo o padrão do `auth-service`)

## 📜 Documentação da API

A API utiliza um formato JSON estruturado para respostas de erro, similar ao `auth-service`.

### Respostas de Erro
Todas as respostas de erro (status `4xx` ou `5xx`) seguem o formato:
```json
{
  "code": "CODIGO_DO_ERRO",
  "message": "Uma mensagem descritiva do erro."
}
```

**Códigos de Erro Comuns:**

| Status HTTP | Código (`code`) | Descrição |
| :--- | :--- | :--- |
| `400 Bad Request` | `INVALID_REQUEST_BODY` | O corpo da requisição é inválido ou malformado. |
| `400 Bad Request` | `INVALID_INPUT` | Um ou mais campos são inválidos (ex: senha muito curta). |
| `401 Unauthorized`| `INVALID_CREDENTIALS` | E-mail ou senha incorretos. |
| `404 Not Found` | `USER_NOT_FOUND` | O usuário solicitado não foi encontrado. |
| `409 Conflict` | `EMAIL_ALREADY_EXISTS` | O e-mail fornecido no cadastro já está em uso. |
| `500 Internal Server Error` | `INTERNAL_SERVER_ERROR` | Ocorreu uma falha inesperada no servidor. |

## 🚀 Como Executar o Projeto
Siga os passos abaixo para colocar o ambiente de desenvolvimento no ar.

### Pré-requisitos
* [Go](https://go.dev/doc/install) (versão 1.24+)
* [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)
* [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)


### Passo a Passo
1.  **Clone o repositório:**
    ```bash
    git clone <url-do-seu-repositorio>
    cd products-service
    ```
2.  **Configure as Variáveis de Ambiente:**
    Crie um arquivo `.env` na raiz do projeto. Você pode copiar o exemplo abaixo.
    ```env
    # Docker Compose
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=productsdb

    # Aplicação (URL para comunicação DENTRO do Docker)
    DATABASE_URL="postgres://postgres:postgres@db:5432/productsdb?sslmode=disable"

    # Segredos
    INTERNAL_API_KEY="uma-chave-secreta-forte-para-apis-internas"

    # Porta que a aplicação ouve DENTRO do container
    LISTEN_ADDR=":8083"
    ```

    3.  **Inicie os Serviços Docker:**
    Este comando irá construir as imagens e iniciar os containers do banco de dados e da aplicação em segundo plano.
    ```bash
    make start
    ```

4.  **Aplique as Migrations:**
    Com o banco de dados no ar, crie as tabelas necessárias.
    ```bash
    make migrate-up
    ```
    Você deve ver uma mensagem de sucesso da migração `create_products_table`.

5.  **Pronto!**
    Sua aplicação está rodando e acessível em `http://localhost:8083`. Você pode acompanhar os logs com `make logs`.

## ⚙️ Comandos do Makefile

* `make start`: Inicia todos os containers em segundo plano.
* `make stop`: Para e remove todos os containers, redes e volumes.
* `make logs`: Exibe os logs do container da aplicação Go.
* `make migrate-up`: Aplica todas as migrações pendentes.
* `make migrate-down`: Reverte a última migração aplicada.
* `make create-migration`: Cria novos arquivos de migração.
* `make lint`: Roda o linter golangci-lint para análise estática do código.
* `make vulncheck`: Roda o govulncheck para buscar vulnerabilidades nas dependências.
* `make gitleaks`: Roda o gitleaks para buscar segredos commitados acidentalmente.

## 🗄️ Acesso ao Banco de Dados

Para visualizar as tabelas e dados, a forma mais fácil é usar o **Adminer**, uma interface gráfica web para bancos de dados.

1.  **Adicione o Serviço ao `docker-compose.yml`:**
    ```yaml
    # ... (dentro de 'services:')
      adminer:
        image: adminer
        container_name: auth-adminer
        restart: always
        ports:
          - "9080:9080" # Usa a porta 9080, pois a app está na 8083
    ```

2.  **Inicie o ambiente com `make start`.**

3.  **Acesse `http://localhost:9080` no seu navegador.**

4.  **Faça login com os seguintes dados:**
    * **System:** `PostgreSQL`
    * **Server:** `db` (nome do serviço do banco no Docker)
    * **Username:** `postgres`
    * **Password:** `postgres`
    * **Database:** `productsdb`
