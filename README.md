# fullcycle-client-server-api

Este repositório contém uma implementação em Go de um serviço de cotação de câmbio entre Dólar (USD) e Real (BRL), com um cliente que faz requisições HTTP para um servidor, o qual consome uma API externa e salva as cotações em um banco de dados SQLite. Já o cliente recebe o resultado e salva em um arquivo de texto.

## Descrição

O serviço é composto por dois arquivos principais:

1. **server.go**: Cria um servidor HTTP que expõe um endpoint `/cotacao` para retornar a cotação do dólar (USD) em relação ao real (BRL). Ele consulta uma API externa (https://economia.awesomeapi.com.br/json/last/USD-BRL) para obter a cotação e a salva em um banco de dados SQLite. O servidor também possui timeouts para chamadas à API externa (200ms) e operações de banco de dados (10ms).

2. **client.go**: Um cliente que faz uma requisição HTTP ao servidor (`server.go`) para obter a cotação do dólar em até 300ms. O cliente salva o valor recebido em um arquivo chamado `cotacao.txt`.

## Requisitos

- Go 1.18 ou superior
- SQLite
- Acesso à internet para consulta à API externa

## Como Rodar

### Passo 1: Baixar o código

Clone o repositório para sua máquina local:

```bash
git clone https://github.com/marcosocram/fullcycle-client-server-api.git
cd fullcycle-client-server-api
```

### Passo 2: Instalar dependências

Certifique-se de que você tenha as dependências necessárias instaladas. Para instalar o driver SQLite para Go, rode:
```bash
go get github.com/mattn/go-sqlite3
```

### Passo 3: Rodar o servidor

Em um terminal, execute o servidor com o comando:
```bash
go run server.go
``` 

Isso fará com que o servidor inicie na porta `8080` e comece a escutar requisições HTTP na rota `/cotacao`.

### Passo 4: Rodar o cliente

Em outro terminal, execute o cliente com o comando:
```bash
go run client.go
``` 

O cliente fará uma requisição ao servidor e salvará o valor da cotação em um arquivo chamado `cotacao.txt`.

### Passo 5: Verificar o resultado

O cliente salvará o valor da cotação no arquivo `cotacao.txt`. O conteúdo do arquivo será semelhante a:
    
```
Dólar: 5.7916
```
   
### Exemplo de Erro

Se o tempo de resposta do servidor ou da API externa ultrapassar os limites estabelecidos, você verá mensagens de erro nos logs, como:

* Erro ao obter a cotação (servidor):
    ```
    Erro ao obter cotação: <error>: context deadline exceeded
    ```
* Erro ao salvar a cotação no banco de dados (servidor):
    ```
    Erro ao salvar cotação no banco: <error>: context deadline exceeded
    ```
* Erro ao obter a cotação (cliente):
    ```
  Erro ao obter cotação: <error>: context deadline exceeded
    ```
  
Esses erros indicam que algum timeout foi atingido (seja ao chamar a API ou ao salvar a cotação no banco).

