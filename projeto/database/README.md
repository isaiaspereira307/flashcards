# Database Setup

Este diretório contém a configuração do PostgreSQL usando Docker Compose.

## Arquivos

- **docker-compose.yml**: Configuração do container PostgreSQL
- **.env**: Variáveis de ambiente (não comitar)
- **.env.example**: Exemplo de variáveis de ambiente
- **init-db.sql**: Script SQL de inicialização

## Como usar

### 1. Iniciar o banco de dados

```bash
docker-compose up -d
```

### 2. Verificar status

```bash
docker-compose ps
```

### 3. Parar o banco de dados

```bash
docker-compose down
```

### 4. Parar e remover volumes

```bash
docker-compose down -v
```

### 5. Ver logs

```bash
docker-compose logs -f postgres
```

## Conectar ao banco

**Via psql:**

```bash
psql -h localhost -U flashcards_user -d flashcards -W
```

**Connection string para aplicações:**

```
postgresql://flashcards_user:flashcards_password@localhost:5432/flashcards
```

## Variáveis de Ambiente

Edite o arquivo `.env` para customizar:

- `DB_NAME`: Nome do banco de dados
- `DB_USER`: Usuário do PostgreSQL
- `DB_PASSWORD`: Senha do usuário
- `DB_PORT`: Porta do PostgreSQL
- `DB_HOST`: Host do PostgreSQL (dentro do docker: postgres, fora: localhost)

## Volumes

Os dados do PostgreSQL são persistidos em um volume Docker chamado `postgres_data`. Eles são mantidos mesmo após parar o container.

## Health Check

O container inclui um health check que verifica a conectividade do PostgreSQL a cada 10 segundos.
