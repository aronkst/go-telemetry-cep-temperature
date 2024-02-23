# Go Telemetry CEP Temperature

## Visão Geral

Este projeto consiste em uma solução integrada de dois serviços para o acesso a informações climáticas detalhadas, utilizando Códigos de Endereçamento Postal (CEPs) como parâmetro de consulta. O Serviço A permite a inserção do CEP através de uma requisição POST, com o corpo `{"cep":"01001000"}`, enquanto o Serviço B disponibiliza os dados climáticos por meio de uma requisição GET, acessível pela URL `/?cep=01001000`. Ao receber uma requisição, o Serviço A processa e redireciona a consulta ao Serviço B para obter as informações climáticas desejadas.

Este projeto é uma solução integrada composta por dois serviços, Serviço A e Serviço B, que oferece acesso a informações climáticas detalhadas, utilizando Códigos de Endereçamento Postal (CEPs) como parâmetro para consulta. Existe a implementação do OpenTelemetry (OTEL) em conjunto com o Zipkin para tracing distribuído, permitindo a visualização da jornada de uma requisição entre o Serviço A e o Serviço B. O Serviço A aceita um CEP via requisição POST, com o corpo `{"cep":"01001000"}`, e consulta o Serviço B, que disponibiliza os dados climáticos através de uma requisição GET na URL `/?cep=01001000`. A integração do OTEL com Zipkin facilita o monitoramento e a análise do tempo de resposta tanto para a busca de CEP quanto para a busca de informações climáticas.

## Características

- **Dois Serviços Integrados**: O projeto é composto por dois serviços distintos, o Serviço A, que recebe o CEP via POST, e o Serviço B, que fornece informações climáticas via GET, facilitando o acesso a dados precisos com base no CEP fornecido.
- **Consulta Direta por CEP no Serviço B**: O Serviço B permite o acesso direto a informações climáticas específicas de uma localização, utilizando o CEP como chave de consulta.
- **Validação Rigorosa de CEP no Serviço A**: O Serviço A implementa uma validação rigorosa do CEP inserido, assegurando que esteja no formato correto e consista apenas de números, com exatamente 8 caracteres, antes de redirecionar a consulta ao Serviço B.
- **Autenticação Livre**: Ambos os serviços foram projetados para serem acessíveis sem a necessidade de autenticação, simplificando o processo de consulta às informações climáticas.
- **Respostas em Formato JSON**: As informações climáticas são fornecidas em formato JSON pelo Serviço B, facilitando a integração com outras aplicações e a manipulação dos dados recebidos.
- **Suporte a Múltiplas Unidades de Temperatura**: O Serviço B oferece informações de temperatura em Celsius, Fahrenheit e Kelvin, atendendo às diversas preferências e necessidades dos usuários.
- **Integração com OTEL + Zipkin para Tracing Distribuído**: A implementação do OpenTelemetry juntamente com o Zipkin proporciona um tracing distribuído eficaz entre o Serviço A e o Serviço B.
- **Spans para Medir Tempos de Resposta**: São criados spans específicos para medir o tempo de resposta das operações de busca de CEP no Serviço A e de busca de informações climáticas no Serviço B.

## Exemplo de Uso

Para consultar informações climáticas através da linha de comando, você pode usar o `curl`, uma ferramenta poderosa e disponível na maioria dos sistemas operacionais para fazer requisições HTTP. Abaixo estão exemplos práticos de como usar o curl para obter a temperatura com base em um CEP específico.

### Realizando uma Consulta

Para fazer uma consulta, simplesmente substitua CEP pelo código postal desejado na URL. Aqui estão alguns exemplos:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

Retorno esperado:

```json
{"city":"São Paulo","temp_C":22.4,"temp_F":72.32,"temp_K":295.55}
```

Neste exemplo, a requisição retorna a temperatura para o CEP 01001000 (um CEP de São Paulo), mostrando a temperatura em Celsius (temp_C), Fahrenheit (temp_F), e Kelvin (temp_K) e a cidade (city).

## Como os Dados são Retornados

Os dados são retornados em formato JSON. Cada campo no JSON representa uma medida diferente de temperatura:

- `city`: Nome da cidade.
- `temp_C`: Temperatura em graus Celsius.
- `temp_F`: Temperatura em graus Fahrenheit.
- `temp_K`: Temperatura em Kelvin.

## Desenvolvimento

No desenvolvimento deste projeto, foquei na criação de uma solução composta por dois serviços interconectados que utilizam APIs externas para fornecer informações climáticas precisas, baseadas em um Código de Endereçamento Postal (CEP) fornecido. O Serviço A é responsável por receber o CEP através de uma requisição POST e, em seguida, comunicar-se com o Serviço B, que executa as consultas às APIs externas e retorna os dados climáticos. Abaixo, descrevo as etapas envolvidas e como cada serviço e API são empregados, incluindo a implementação do OpenTelemetry (OTEL) e Zipkin para tracing distribuído.

### Busca de Endereço pelo CEP com viacep.com.br (Serviço B)

A jornada começa quando o Serviço B coleta informações detalhadas sobre o endereço usando o CEP fornecido pelo Serviço A. Para isso, ele consulta a API do ViaCEP, que retorna dados como logradouro, bairro, cidade e estado. Estas informações são cruciais para identificar a localização geográfica precisa para as consultas climáticas subsequentes.

### Busca de Longitude e Latitude com nominatim.openstreetmap.org (Serviço B)

Com os dados do endereço em mãos, o Serviço B então converte estas informações em coordenadas geográficas (latitude e longitude) através da API do Nominatim, que faz parte do projeto OpenStreetMap. Esta conversão é essencial para garantir a precisão das consultas climáticas que dependem de coordenadas geográficas.

### Busca de Temperatura (Serviço B)

Dispondo das coordenadas geográficas, o Serviço B realiza a consulta das condições climáticas atuais. Dependendo da disponibilidade dos dados, ele pode usar:

- A API Open-Meteo, para consultas climáticas detalhadas baseadas em coordenadas, fornecendo informações precisas de temperatura para a localização especificada.
- A API do wttr.in, para informações climáticas baseadas em nomes de localização, que embora possa não ser tão precisa quanto a consulta por coordenadas, ainda fornece uma estimativa válida das condições climáticas.

### Integração com OTEL + Zipkin

A integração com OpenTelemetry (OTEL) e Zipkin adiciona uma camada de observabilidade ao projeto, permitindo o tracing distribuído entre o Serviço A e o Serviço B. Esta funcionalidade possibilita a monitorização da jornada completa de uma requisição, incluindo a medição do tempo de resposta para a busca de CEP e a busca de temperatura, facilitando a identificação e a resolução de possíveis gargalos ou problemas de desempenho.

## Tratamentos de Erros

Implementei tratamentos de erros em cada etapa para assegurar que o sistema possa lidar de forma adequada com cenários como CEPs inválidos, falhas na obtenção de coordenadas ou erros nas respostas das APIs.

## Testes Unitários

Uma parte do desenvolvimento deste projeto envolve a implementação de testes unitários abrangentes, garantindo a confiabilidade e a robustez de cada funcionalidade oferecida pela aplicação. A abordagem adotada para os testes segue as melhores práticas de desenvolvimento de software, focando na validação de cada componente isoladamente para assegurar seu correto funcionamento em diversos cenários.

### Cobertura dos Testes

Os testes unitários cobrem uma ampla gama de casos de uso e cenários de erro, incluindo, mas não se limitando a:

- Validação de CEPs: Testes para assegurar que apenas CEPs válidos e no formato correto são aceitos, e que as mensagens de erro adequadas são retornadas para CEPs inválidos ou formatados incorretamente.
- Consulta a APIs Externas: Testes para verificar a interação correta com as APIs externas usadas para obter informações de endereço, coordenadas geográficas e dados climáticos. Isso inclui simular respostas das APIs para testar o manejo adequado de dados e erros.
- Conversão de Unidades de Temperatura: Testes que validam a precisão das conversões de temperatura entre Celsius, Fahrenheit e Kelvin, garantindo que os cálculos estejam corretos.
- Tratamento de Erros: Testes específicos para verificar a robustez do sistema ao enfrentar erros durante a consulta de informações, incluindo falhas de rede, erros nas APIs externas e dados inesperados.

## Makefile

Este projeto inclui um Makefile projetado para oferecer uma interface eficiente e simplificada para o gerenciamento dos ambientes de desenvolvimento e produção, além da execução de testes automatizados. Os comandos disponibilizados permitem otimizar e agilizar o fluxo de trabalho de desenvolvimento, testes e manutenção do projeto, assegurando uma gestão mais eficaz e organizada.

### Comandos de Desenvolvimento

### `make dev-start`

Inicia os serviços definidos no arquivo `docker-compose.dev.yml` para o ambiente de desenvolvimento em modo detached (em segundo plano). Isso permite que os serviços rodem em background sem ocupar o terminal.

### `make dev-stop`

Interrompe os serviços que estão rodando em background para o ambiente de desenvolvimento. Isso não remove os containers, redes ou volumes criados pelo `docker compose up`.

### `make dev-down`

Desliga os serviços do ambiente de desenvolvimento e remove os containers, redes e volumes associados criados pelo `docker compose up`. Utilize este comando para limpar recursos após o desenvolvimento.

### `dev-run-service-a`

Inicia a execução do Serviço A dentro do ambiente de desenvolvimento, utilizando o Docker Compose para executar o comando `go run` no arquivo `/cmd/input_server/main.go`. Ele é ideal para iniciar rapidamente o servidor do projeto em modo de desenvolvimento.

### `dev-run-service-b`

Inicia a execução do Serviço B dentro do ambiente de desenvolvimento, utilizando o Docker Compose para executar o comando `go run` no arquivo `/cmd/temperature_server/main.go`. Ele é ideal para iniciar rapidamente o servidor do projeto em modo de desenvolvimento.

### `make dev-run-tests`

Executa todos os testes Go dentro do ambiente de desenvolvimento, mostrando detalhes verbosos de cada teste. Este comando é útil para rodar a suíte de testes do projeto e verificar se tudo está funcionando como esperado.

### Comandos de Produção

### `make prod-start`

Inicia os serviços definidos no arquivo `docker-compose.prod.yml` para o ambiente de produção em modo detached. Isso é útil para rodar o projeto em um ambiente que simula a produção.

### `make prod-stop`

Interrompe os serviços do ambiente de produção que estão rodando em background, sem remover os containers, redes ou volumes associados.

### `make prod-down`

Desliga os serviços do ambiente de produção e remove os containers, redes e volumes associados, limpeza de recursos após o uso em produção.

## Pré-requisitos

Antes de começar, certifique-se de que você tem o Docker e o Docker Compose instalados em sua máquina. Caso não tenha, você pode baixar e instalar a partir dos seguintes links:

- Docker: https://docs.docker.com/get-docker/

### Clonar o Repositório

Primeiro, clone o repositório do projeto para a sua máquina local. Abra um terminal e execute o comando:

```bash
git clone https://github.com/aronkst/go-telemetry-cep-temperature.git
```

### Navegar até o Diretório do Projeto

Após clonar o repositório, navegue até o diretório do projeto utilizando o comando cd:

```bash
cd go-telemetry-cep-temperature
```

## Ambiente de Desenvolvimento

### Construir o Projeto com Docker Compose

No diretório do projeto, execute o seguinte comando para construir e iniciar o projeto utilizando o Docker Compose:

```bash
docker compose -f docker-compose.dev.yml up --build
```

Ou utilizando o Makefile:

```bash
make dev-start
```

Este comando irá construir a imagem Docker do projeto e iniciar o container.

### Executar o Projeto com Docker Compose

Para iniciar o serviço principal do seu projeto em modo de desenvolvimento, você pode utilizar os comandos diretos do Docker Compose:

```bash
docker compose -f docker-compose.dev.yml exec dev go run cmd/input_server/main.go
```

```bash
docker compose -f docker-compose.dev.yml exec dev go run cmd/temperature_server/main.go
```

Ou utilizando o Makefile:

```bash
make dev-run-service-a
```

```bash
make dev-run-service-b
```

### Acessar o Projeto

Com o container rodando, você pode acessar o projeto através do navegador ou utilizando ferramentas como curl, apontando para http://localhost:3000/, substituindo CEP pelo código postal desejado.

### Exemplo de Comando curl

Para testar se o projeto está rodando corretamente, você pode usar o seguinte comando curl em um novo terminal:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

Você deverá receber uma resposta em JSON com as temperaturas em Celsius, Fahrenheit, Kelvin e a cidade.

### Visualizando a Telemetria

Abra um navegador e acesse http://localhost:9411/zipkin/. Essa URL levará você à interface de usuário do Zipkin, onde você pode começar a visualizar a telemetria dos seus serviços.

1. **Pesquisar por Rastreamentos**: Na interface do Zipkin, você pode pesquisar rastreamentos de várias maneiras, como por serviço, nome da operação, anotações e tags.

2. **Analisar Rastreamentos**: Após encontrar um rastreamento específico, você pode clicar nele para ver os detalhes. Isso inclui informações como a duração do rastreamento.

### Encerrando o Projeto

Para encerrar o projeto e parar o container do Docker, volte ao terminal onde o Docker Compose está rodando e pressione Ctrl+C. Para remover os containers criados pelo Docker Compose, execute:

```bash
docker compose -f docker-compose.dev.yml down
```

Ou utilizando o Makefile:

```bash
make dev-down
```

## Ambiente de Produção

### Construir e Executar o Projeto com Docker Compose

No diretório do projeto, execute o seguinte comando para construir e iniciar o projeto no ambiente de produção utilizando o Docker Compose:

```bash
docker compose -f docker-compose.prod.yml up --build
```

Ou utilizando o Makefile:

```bash
make prod-start
```

Este comando irá construir a imagem Docker do projeto para produção e iniciar os containers.

## Exemplo de Comando curl

Para verificar se o projeto em produção está operacional, utilize o seguinte comando curl, ajustando o endereço conforme sua configuração:

```bash
curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"cep":"01001000"}'
```

Você deverá receber uma resposta em JSON com as informações solicitadas, como as temperaturas em Celsius, Fahrenheit, Kelvin e a cidade.

### Visualizando a Telemetria

Abra um navegador e acesse http://localhost:9411/zipkin/. Essa URL levará você à interface de usuário do Zipkin, onde você pode começar a visualizar a telemetria dos seus serviços.

1. **Pesquisar por Rastreamentos**: Na interface do Zipkin, você pode pesquisar rastreamentos de várias maneiras, como por serviço, nome da operação, anotações e tags.

2. **Analisar Rastreamentos**: Após encontrar um rastreamento específico, você pode clicar nele para ver os detalhes. Isso inclui informações como a duração do rastreamento.

### Encerrando o Projeto

Para encerrar o projeto e parar os containers de produção, utilize o seguinte comando:

```bash
docker compose -f docker-compose.prod.yml down
```

Ou utilizando o Makefile:

```bash
make prod-down
```

Este comando encerra todos os serviços de produção e remove os containers, redes e volumes associados, assegurando que o ambiente de produção seja limpo após o uso.
