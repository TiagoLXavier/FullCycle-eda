# Wallet Core + Balances

Desafio Full Cycle de Arquitetura Orientada a Eventos (EDA). O microsserviço **Balances** (Go) consome os eventos publicados pelo **Wallet Core** no Kafka, persiste o saldo de cada conta em um banco próprio e expõe a consulta `GET /balances/{account_id}` na porta **3003**.

## Pré-requisitos

- Docker
- Docker Compose v2 (`docker compose`)
- Portas livres no host: `8080`, `3003`, `3306`, `3307`, `9092`, `2181`

Nada mais precisa ser instalado: Go, MySQL e Kafka rodam dentro dos containers.

## Como rodar

```bash
git clone https://github.com/TiagoLXavier/FullCycle-eda.git
cd FullCycle-eda
docker compose up -d
```

Só isso. Na primeira execução o Compose faz o build das imagens do `wallet-core` e do `balances` (alguns minutos, por causa do download das dependências). Ao subir:

- os tópicos `transactions` e `balances` são criados no Kafka (serviço `kafka-init`);
- as tabelas são criadas (migrations) e os dados fictícios são inseridos (seeds) nos dois bancos, `wallet` e `balances`;
- o `wallet-core` e o `balances` só iniciam depois que Kafka e bancos estão saudáveis.

Não é preciso rodar nenhum script manual.

Para conferir se está tudo de pé:

```bash
docker compose ps
```

O resultado esperado é `wallet-core`, `balances`, `kafka`, `mysql-wallet`, `mysql-balances` e `zookeeper` como `running` (os de Kafka e MySQL como `healthy`), e `kafka-init` como `exited (0)`, porque ele roda uma única vez.

## Arquitetura

```
POST /transactions            Kafka (tópico balances)             GET /balances/{account_id}
  cliente --> wallet-core -----------------------------> balances <----------- cliente
                 |                                          |
            mysql-wallet                              mysql-balances
```

1. O **Wallet Core** (porta 8080) recebe uma transação, debita e credita as contas no banco `wallet` e publica dois eventos: `TransactionCreated` no tópico `transactions` e `BalanceUpdated` no tópico `balances`.
2. O evento `BalanceUpdated` traz o saldo atualizado das contas de origem e de destino:

   ```json
   {
     "Name": "BalanceUpdated",
     "Payload": {
       "account_id_from": "11111111-1111-1111-1111-111111111111",
       "account_id_to": "22222222-2222-2222-2222-222222222222",
       "balance_account_id_from": 900,
       "balance_account_id_to": 600
     }
   }
   ```

3. O **Balances** (porta 3003) escuta o tópico `balances` (consumer group `balances`), grava o saldo de cada conta no seu banco próprio (`balances`, em `mysql-balances`) e responde às consultas.

## Serviços do docker-compose

| Serviço          | Função                                         | Porta no host |
|------------------|------------------------------------------------|---------------|
| `zookeeper`      | Coordenação do Kafka                           | 2181          |
| `kafka`          | Broker de mensagens                            | 9092          |
| `kafka-init`     | Cria os tópicos `transactions` e `balances`    | -             |
| `mysql-wallet`   | Banco do Wallet Core (`wallet`)                | 3306          |
| `mysql-balances` | Banco do Balances (`balances`)                 | 3307          |
| `wallet-core`    | Microsserviço produtor                         | 8080          |
| `balances`       | Microsserviço consumidor + API de consulta     | 3003          |

As migrations e os seeds ficam em `.docker/wallet-init` e `.docker/balances-init` e são executados automaticamente pelo MySQL na criação dos containers.

## Dados do seed

| Conta                                  | Cliente  | Saldo inicial |
|----------------------------------------|----------|---------------|
| `11111111-1111-1111-1111-111111111111` | John Doe | 1000          |
| `22222222-2222-2222-2222-222222222222` | Jane Doe | 500           |

As mesmas contas e saldos estão nos dois bancos, então a consulta no Balances já funciona logo após a subida, antes de qualquer transação.

## Como testar

### Com o arquivo api.http

Abra o `api.http` na raiz do projeto (extensão REST Client do VS Code ou HTTP Client do JetBrains) e execute:

1. `GET /balances/1111...` e `GET /balances/2222...` retornam os saldos do seed (1000 e 500).
2. `POST /transactions` transfere 100 da conta `1111...` para a conta `2222...`.
3. Consulte de novo `GET /balances/1111...` e `GET /balances/2222...`: os saldos passam a 900 e 600, atualizados pelo evento consumido do Kafka.
4. `GET /balances/0000...` (conta inexistente) retorna `404`.
5. Opcional: `POST /clients` e `POST /accounts` criam novos clientes e contas no Wallet Core. Na criação de conta, substitua o `client_id` pelo id retornado ao criar o cliente.

### Com curl

```bash
curl http://localhost:3003/balances/11111111-1111-1111-1111-111111111111
```

```bash
curl -X POST http://localhost:8080/transactions -H "Content-Type: application/json" -d '{"account_id_from":"11111111-1111-1111-1111-111111111111","account_id_to":"22222222-2222-2222-2222-222222222222","amount":100}'
```

```bash
curl http://localhost:3003/balances/11111111-1111-1111-1111-111111111111
```

Resposta esperada após a transação:

```json
{"account_id":"11111111-1111-1111-1111-111111111111","balance":900}
```

Conta inexistente:

```json
{"error":"account not found"}
```

### Acompanhando o fluxo de eventos

```bash
docker compose logs -f wallet-core balances
```

A cada transação aparecem `UpdateBalanceKafkaHandler called` no `wallet-core` e `balances updated: <conta>=<saldo> ...` no `balances`.

## Estrutura do Balances

O Balances segue as mesmas camadas do Wallet Core, no mesmo módulo Go:

| Camada        | Arquivo                                                |
|---------------|--------------------------------------------------------|
| Entidade      | `internal/entity/balance.go`                           |
| Gateway       | `internal/gateway/balance.go`                          |
| Banco         | `internal/database/balance_db.go`                      |
| Casos de uso  | `internal/usecase/update_balance`, `internal/usecase/get_balance` |
| Consumidor    | `internal/event/handler/balance_updated_consumer.go`   |
| API           | `internal/web/balance_handler.go`                      |
| Entrada       | `cmd/balances/main.go`                                 |

Testes unitários:

```bash
docker run --rm -v "$(pwd)":/app -w /app golang:1.20 go test ./...
```

## Checklist do desafio

| Requisito                                                  | Onde está                                                         |
|------------------------------------------------------------|-------------------------------------------------------------------|
| Microsserviço consumidor do Kafka                          | `cmd/balances`, `internal/event/handler/balance_updated_consumer.go` (tópico `balances`) |
| Banco de dados próprio do Balances                         | serviço `mysql-balances`, banco `balances`, `internal/database/balance_db.go` |
| `GET /balances/{account_id}` na porta 3003                 | `internal/web/balance_handler.go`, porta `3003`                   |
| Kafka, Wallet Core, Balances e bancos no mesmo compose     | `docker-compose.yaml`                                             |
| Migrations automáticas                                     | `.docker/wallet-init/01-schema.sql`, `.docker/balances-init/01-schema.sql` |
| Seeds automáticos nos dois bancos                          | `.docker/wallet-init/02-seed.sql`, `.docker/balances-init/02-seed.sql`     |
| Subida com um único comando, sem scripts manuais           | `docker compose up -d`                                            |
| Arquivo de chamadas para teste                             | `api.http`                                                        |

## Variáveis de ambiente

Os dois serviços aceitam `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` e `KAFKA_BOOTSTRAP_SERVERS`. Os valores usados no compose já estão configurados, então não é preciso criar `.env`.

## Como parar e resetar

```bash
docker compose down -v
```

Isso remove containers e volumes. No próximo `docker compose up -d` os bancos são recriados e as migrations e seeds são aplicados de novo, voltando aos saldos iniciais.

## Problemas comuns

- **Porta em uso** (por exemplo, um MySQL local na 3306): pare o serviço local ou altere a porta do host no `docker-compose.yaml`.
- **Mac com Apple Silicon**: as imagens usam `platform: linux/amd64` e rodam por emulação. A primeira subida pode demorar mais.
- **Serviço não respondeu logo após o `up`**: aguarde até todos os containers aparecerem como `running`/`healthy` no `docker compose ps`.
