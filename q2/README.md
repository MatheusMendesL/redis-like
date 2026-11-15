# Redis-like

Servidor de armazenamento chave-valor inspirado no Redis, escrito em Go.

## Pré-requisitos

- Go 1.26 ou superior
- Git
- Um cliente TCP, como `nc`, `telnet` ou `ncat`

## Linux

Na raiz do projeto, execute:

```bash
go run ./cmd/server
```

Para habilitar a persistência AOF:

```bash
go run ./cmd/server -aof
```

Para usar outra porta:

```bash
go run ./cmd/server -port 6380
```

Em outro terminal, conecte-se ao servidor:

```bash
nc localhost 6379
```

## Windows

No PowerShell, na raiz do projeto, execute:

```powershell
go run .\cmd\server
```

Com persistência AOF:

```powershell
go run .\cmd\server -aof
```

Para conectar usando `telnet`:

```powershell
telnet localhost 6379
```

## Comandos

Os comandos disponíveis são:

```text
SET chave valor
GET chave
DEL chave
EXISTS chave
EXPIRE chave segundos
TTL chave
PERSIST chave
KEYS [padrão]
PING [mensagem]
INFO
QUIT
```

Exemplo:

```text
PING
SET nome Matheus
GET nome
QUIT
```

Quando o AOF está habilitado, as operações são salvas em `data.aof`.
