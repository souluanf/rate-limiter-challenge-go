# Testando

- Instale as dependências:

   ```bash
   go mod tidy
   ```

- Inicialize o Redis
   ```bash
   docker-compose up -d
   ```

- Inicie a aplicação
   ```bash
   go run cmd/server/main.go
   ```

- Abra outro terminal e execute os testes
   ```bash
   go run test/main.go
   ```