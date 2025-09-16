# Microsserviço de Produtos (Product Service)

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)
![Docker](https://img.shields.io/badge/Docker-20.10-blue.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)

## 📖 Sobre o Projeto

Este é o **Microsserviço de Produtos**, uma parte fundamental do sistema de e-commerce distribuído. Desenvolvido em Go, a sua responsabilidade principal é ser a fonte da verdade para todo o catálogo de produtos, gerindo informações como nome, preço e controlo de stock.

Este serviço expõe endpoints públicos para a consulta de produtos e endpoints internos protegidos para a gestão do catálogo, que são consumidos por outros serviços, como o **Serviço de Pedidos**.

### ✨ Funcionalidades Principais
* Listagem de todos os produtos disponíveis.
* Consulta dos detalhes de um produto específico.
* Endpoint interno para a criação de novos produtos.
* Endpoint interno para a atualização de stock.
* Segurança para operações internas via API Key, garantindo que apenas outros serviços autorizados possam modificar o catálogo.

## 🛠️ Arquitetura e Tecnologias

O projeto segue uma arquitetura em camadas para uma clara separação de responsabilidades (API, Lógica de Negócio, Repositório), consistente com os outros serviços do ecossistema.

### Tecnologias Utilizadas
* **Linguagem:** Go
* **Banco de Dados:** PostgreSQL
* **Containerização:** Docker & Docker Compose
* **Roteador HTTP:** Chi
* **Migrations:** golang-migrate
* **Automação:** Makefile
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

Este serviço foi desenhado para ser executado como parte de um ambiente Docker Compose junto com os outros microsserviços do e-commerce.

### Pré-requisitos
* [Go](https://go.dev/doc/install) (versão 1.24+)
* [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/)
* [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Passo a Passo
1.  **Clone o repositório** e certifique-se de que este serviço (`product-service`) está na mesma pasta raiz que os outros serviços (ex: `auth-service`).

2.  **Configure as Variáveis de Ambiente:**
    No ficheiro `.env` central do seu projeto de e-commerce, garanta que as variáveis para o banco de dados de produtos estão definidas.
    ```env
    # Docker Compose - Product DB
    POSTGRES_USER_PRODUCT=postgres
    POSTGRES_PASSWORD_PRODUCT=postgres
    POSTGRES_DB_PRODUCT=productdb
    
    # Segredos partilhados
    INTERNAL_API_KEY="uma-chave-secreta-forte-para-apis-internas"
    ```

3.  **Atualize o `docker-compose.yml` Principal:**
    No `docker-compose.yml` da raiz do seu e-commerce, adicione os serviços para a aplicação e o banco de dados do `product-service`.
    ```yaml
    # docker-compose.yml (exemplo de como integrar)

    services:
      # ... (seus serviços existentes, como auth-app e auth-db)

      # Novo serviço para o Banco de Dados de Produtos
      product-db:
        image: postgres:15-alpine
        environment:
          POSTGRES_USER: ${POSTGRES_USER_PRODUCT}
          POSTGRES_PASSWORD: ${POSTGRES_PASSWORD_PRODUCT}
          POSTGRES_DB: ${POSTGRES_DB_PRODUCT}
        volumes:
          - product_postgres_data:/var/lib/postgresql/data

      # Novo serviço para a Aplicação de Produtos
      product-app:
        build: ./product-service  # Caminho para a pasta deste projeto
        ports:
          - "8083:8083"
        env_file:
          - ./.env
        environment:
          # Constrói a URL do banco usando o nome do serviço 'productdb'
          DATABASE_URL: "postgres://${POSTGRES_USER_PRODUCT}:${POSTGRES_PASSWORD_PRODUCT}@productdb:5432/${POSTGRES_DB_PRODUCT}?sslmode=disable"
          LISTEN_ADDR: ":8083"
        depends_on:
          - productdb
          - auth-app # Dependência opcional se precisar de validação

    volumes:
      # ... (volumes existentes)
      product_postgres_data: {}
    ```

4.  **Execute o Ambiente Completo:**
    A partir da pasta raiz que contém o `docker-compose.yml`, execute:
    ```bash
    docker-compose up --build
    ```
    O seu `product-service` estará acessível em `http://localhost:8083`.

5.  **Aplique as Migrations:**
    Com o banco de dados no ar, crie as tabelas necessárias.
    ```bash
    make migrate-up
    ```
    Você deve ver uma mensagem de sucesso da migração `create_products_table`.

6.  **Pronto!**
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
