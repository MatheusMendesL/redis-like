# Redis-like

Um servidor de armazenamento chave-valor inspirado no Redis, escrito em Go. O
servidor aceita comandos por TCP, uma linha por comando, e oferece persistência
opcional em um arquivo AOF (append-only file).

## Pré-requisitos

- Go 1.26 ou superior
- Git
- Um cliente TCP para testar a conexão (`nc`/`netcat` ou `telnet`)

## Executando no Linux

Clone o projeto e entre na pasta:

```bash
git clone https://github.com/MatheusMendesL/redis-like.git
cd redis-like
```

Inicie o servidor na porta padrão (`6379`):

```bash
go run ./cmd/server
```

Para habilitar a persistência AOF:

```bash
go run ./cmd/server -aof
```

Também é possível escolher outra porta:

```bash
go run ./cmd/server -port 6380
go run ./cmd/server -port 6380 -aof
```

Em outro terminal, conecte-se ao servidor:

```bash
nc localhost 6379
```

Digite comandos, um por linha, e pressione `Enter`:

```text
PING
SET nome Matheus
GET nome
EXISTS nome
KEYS
QUIT
```

## Executando no Windows

Abra o PowerShell e clone o projeto:

```powershell
git clone https://github.com/MatheusMendesL/redis-like.git
Set-Location redis-like
```

Inicie o servidor:

```powershell
go run .\cmd\server
```

Com AOF ou usando outra porta:

```powershell
go run .\cmd\server -aof
go run .\cmd\server -port 6380
```

Se você tiver `telnet` habilitado no Windows, abra outro PowerShell e conecte-se:

```powershell
telnet localhost 6379
```

Como alternativa, use `ncat`, disponível no pacote Nmap:

```powershell
ncat.exe localhost 6379
```

## Compilando um executável

Linux:

```bash
go build -o redis-like-server ./cmd/server
./redis-like-server -aof
```

Windows PowerShell:

```powershell
go build -o redis-like-server.exe .\cmd\server
.\redis-like-server.exe -aof
```

## Comandos disponíveis

| Comando | Uso | Descrição |
| --- | --- | --- |
| `SET` | `SET chave valor` | Armazena um valor |
| `GET` | `GET chave` | Recupera um valor |
| `DEL` | `DEL chave` | Remove uma chave |
| `EXISTS` | `EXISTS chave` | Verifica se a chave existe |
| `EXPIRE` | `EXPIRE chave segundos` | Define expiração |
| `TTL` | `TTL chave` | Consulta o tempo restante |
| `PERSIST` | `PERSIST chave` | Remove a expiração |
| `KEYS` | `KEYS [padrão]` | Lista chaves |
| `PING` | `PING [mensagem]` | Testa a conexão |
| `INFO` | `INFO` | Exibe informações do servidor |
| `QUIT` | `QUIT` | Encerra a conexão |

O modo AOF grava as operações de escrita em `data.aof` e reproduz esse arquivo
ao iniciar o servidor. O arquivo é criado na pasta em que o comando foi
executado.

## Estrutura principal

- `cmd/server`: ponto de entrada do servidor
- `internal/adapter`: protocolo e conexões TCP
- `internal/domain`: comandos, entidades e contratos
- `internal/infrastructure`: armazenamento em memória e persistência AOF
- `internal/usecase`: execução dos comandos
