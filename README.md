# Testando

- Instale as dependências:

   ```bash
   go mod tidy
   ```

- Inicialize a aplicação com o docker-compose
   ```bash
   docker-compose up -d
   ```

- Abra outro terminal e execute os testes
   ```bash
   go run test/main.go
   ```