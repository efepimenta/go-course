## Curso de golang

> Desafio 1: Client Server API

- na pasta principal, temos uma pasta chamada `client` e outra chamada `server`
- em server, o arquivo `server.go` é o servidor
- em client, o arquivo `client.go` é o cliente
- o servidor deve ser executado primeiro
  - go run server/server.go
  - pode ser executado de dentro da pasta server, o banco de dados é criado em /server
    - cd server && go run server.go
- ao executar o servidor, o banco de dados sera criado automaticamente
- o cliente deve ser executado em outro terminal
  - go run client/client.go
  - pode ser executado de dentro da pasta client, o arquivo txt é criado em /client
    - cd client && go run client.go
