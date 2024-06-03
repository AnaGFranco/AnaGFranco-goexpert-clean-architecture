# Desafio: Clean Architecture 

## Executando o Projeto

1. Clone o repositório para o seu ambiente local:
   ```sh
   git clone git@github.com:AnaGFranco/goexpert-clean-architecture.git
    ```
2. Abra um terminal na raiz do projeto.
3. Execute o seguinte comando para construir e iniciar os serviços definidos no arquivo docker-compose:
   ```sh
   make docker-compose-up
   ```
4. Execute o seguinte comando para iniciar a aplicação
   ```sh
   make run
   ```
   
###  API HTTP/Web

A interface HTTP/Web permite a interação com os pedidos. Abaixo estão exemplos de como usar os endpoints HTTP:

- **Criar um novo pedido:**

  ```sh
  curl --location 'http://localhost:8000/order' \
  --header 'Content-Type: application/json' \
  --data '{
      "id": "1",
      "price": 50.5,
      "tax": 0.12
  }'
  ```

- **Listar pedidos:**

  ```sh
  curl --location 'http://localhost:8000/order'
  ```


###  GraphQL

A aplicação fornece uma API GraphQL para interação. Abaixo estão alguns exemplos de queries e mutations que podem ser realizadas. O GraphQL Playground pode ser acessado em [http://localhost:8080/](http://localhost:8080/).


- **Criar pedido:**

  ```graphql
  mutation createOrder {
    createOrder(input: {id: "2", Price: 50.5, Tax: 0.12}) {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

- **Listar pedidos:**

  ```graphql
  query  {
    listOrders {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

### GRPC 

Para interagir com o serviço gRPC da sua aplicação usando Evans, siga as instruções abaixo:

1. Execute o comando para iniciar o Evans em modo REPL, permitindo a interação com o serviço gRPC da aplicação:
   ```sh
   make grpc-run
   ```

2. Chamando métodos do serviço OrderService:

**Criar pedido:**

      ```sh
      grpcurl -plaintext -d '{"id":"xyz","price": 100.5, "tax": 0.5}' localhost:50051 pb.OrderService/CreateOrder
      ```


**Listar pedidos:**

      ```sh
    grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
      ```
