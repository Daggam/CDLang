# CDLang

**CDLang** (Conceptos De 'Languages') es un mini intérprete escrito en Go para un DSL (lenguaje específico de dominio) orientado a describir operaciones de intercambio de coleccionables entre usuarios: crear ofertas, consultarlas, enviarlas, aceptarlas, rechazarlas o eliminarlas.

El proyecto implementa un pipeline clásico de intérprete (lexer → parser → AST → evaluator) e incluye un REPL interactivo, además de un comando `EXPLAIN` que muestra, paso a paso, cómo se procesa cada sentencia.

## Características

- REPL interactivo para escribir y ejecutar sentencias CDLang.
- Lexer y parser propios (sin dependencias externas). (Recursive Descent Parser)
- Evaluador de expresiones y sentencias sobre un entorno en memoria. (Tree-Walking Interpreter)
- Comando `EXPLAIN`, que muestra el código fuente, los tokens generados, el AST, la validación semántica, el plan de ejecución y el resultado esperado de una sentencia.
- Escrito 100% en Go, sin dependencias de terceros.

## Requisitos

- [Go](https://go.dev/dl/) 1.25.6 o superior.

## Instalación

Existen dos formas de instalar CDLang: usando `go install` (recomendado) o compilando desde el código fuente.

### Opción 1: usando `go install`

```bash
go install github.com/Daggam/CDLang/cmd/cdlang@latest
```

Esto descarga, compila e instala el binario `cdlang` en tu `$GOBIN` (por defecto `$(go env GOPATH)/bin`). Asegúrate de que esa ruta esté en tu `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Una vez instalado, puedes ejecutar el intérprete desde cualquier lugar:

```bash
cdlang
```

### Opción 2: desde el código fuente

1. Clona el repositorio:

   ```bash
   git clone https://github.com/Daggam/CDLang.git
   cd CDLang
   ```

2. Compila el binario:

   ```bash
   go build -o cdlang ./cmd/cdlang
   ```

3. Ejecuta el binario generado:

   ```bash
   ./cdlang
   ```

   Alternativamente, puedes ejecutarlo directamente sin generar un binario:

   ```bash
   go run ./cmd/cdlang
   ```

## Uso

Al iniciar, CDLang muestra un prompt (`>>`) donde puedes escribir sentencias del lenguaje:

```
Bienvenido a CDL v1.0.0.
>> OFFER AR-LM10, FR-KM10(20), AR-RD7;
```

### Sentencias disponibles

Todas estas sentencias son gramáticamente validas, pero solo algunas son semánticamente validas.

Aquellas sentencias que son semánticamente validas son: `OFFER <argumentos>`, `GET OFFER <ident>`, `EXPLAIN <consulta>`


| Sentencia | Producción sintáctica | Ejemplo |
|---|---|---|
| Crear una oferta | `"OFFER" <argumentos>` | `OFFER AR-LM10, FR-KM10(20), AR-RD7;` |
| Consultar ofertas de un coleccionable | `"GET" "OFFER" <ident>` | `GET OFFER AR-LM10;` |
| Enviar una propuesta de intercambio | `"SEND" "OFFER" <argumentos> "FOR" <argumentos> "IN" "USER" <ident>` | `SEND OFFER FR-KM FOR AR-LM10 IN USER pepe;` |
| Ver ofertas disponibles | `"VIEW" "OFFER"` | `VIEW OFFER;` |
| Aceptar un intercambio (trade) | `"ACCEPT" "TRADE" <int_list>?` |`ACCEPT TRADE;` o `ACCEPT TRADE 58;` |
| Rechazar un intercambio (trade) | `"DECLINE" "TRADE" <int_list>?` |`DECLINE TRADE;` o `DECLINE TRADE 58;` |
| Eliminar una oferta propia | `"DELETE" "OFFER" <argumentos>` |`DELETE OFFER AR-LM10;` |


### Comando `EXPLAIN`

Antepón `EXPLAIN` a cualquier sentencia para ver en detalle cómo la procesa el intérprete (código fuente, tokens, AST, validación semántica, plan de ejecución y resultado):

```
>> EXPLAIN ACCEPT TRADE 58;
```

## Catálogo de coleccionables

El lenguaje es genérico, sin embargo hemos decidio utilizar como coleccionables a jugadores de fútbol.
Cada coleccionable tiene un código único con el formato `PAIS-INICIALESNUMERO` (por ejemplo, `AR-LM10` para Lionel Messi).

El catálogo actual incluye los siguientes coleccionables, agrupados por país:

### Argentina (AR)

| Constante | Código | Jugador |
|---|---|---|
| `AR_LM10` | `AR-LM10` | Lionel Messi |
| `AR_AD11` | `AR-AD11` | Ángel Di María |
| `AR_EM23` | `AR-EM23` | Emiliano Martínez |
| `AR_RD7`  | `AR-RD7`  | Rodrigo De Paul |
| `AR_JA9`  | `AR-JA9`  | Julián Álvarez |

### Brasil (BR)

| Constante | Código | Jugador |
|---|---|---|
| `BR_VJ7`  | `BR-VJ7`  | Vinícius Júnior |
| `BR_NJ10` | `BR-NJ10` | Neymar Jr. |
| `BR_AB1`  | `BR-AB1`  | Alisson Becker |
| `BR_RG11` | `BR-RG11` | Rodrygo Goes |
| `BR_MQ5`  | `BR-MQ5`  | Marquinhos |

### Francia (FR)

| Constante | Código | Jugador |
|---|---|---|
| `FR_KM10` | `FR-KM10` | Kylian Mbappé |
| `FR_AG7`  | `FR-AG7`  | Antoine Griezmann |
| `FR_OG9`  | `FR-OG9`  | Olivier Giroud |
| `FR_NK13` | `FR-NK13` | N'Golo Kanté |
| `FR_AT8`  | `FR-AT8`  | Aurélien Tchouaméni |

### España (ES)

| Constante | Código | Jugador |
|---|---|---|
| `ES_AM7`  | `ES-AM7`  | Álvaro Morata |
| `ES_RH16` | `ES-RH16` | Rodrigo Hernández (Rodri) |
| `ES_DC2`  | `ES-DC2`  | Dani Carvajal |
| `ES_LY19` | `ES-LY19` | Lamine Yamal |
| `ES_PG20` | `ES-PG20` | Pedri González |

### Reino Unido / Inglaterra (GB)

| Constante | Código | Jugador |
|---|---|---|
| `GB_HK9`  | `GB-HK9`  | Harry Kane |
| `GB_JB10` | `GB-JB10` | Jude Bellingham |
| `GB_BS7`  | `GB-BS7`  | Bukayo Saka |
| `GB_PF11` | `GB-PF11` | Phil Foden |
| `GB_DR4`  | `GB-DR4`  | Declan Rice |

> Este catálogo está definido en el código fuente (`internal/object/object.go`). Para agregar nuevos coleccionables es necesario declarar una nueva constante siguiendo el mismo formato `PAIS-INICIALESNUMERO`.

## Estructura del proyecto

```
CDLang/
├── cmd/
│   └── cdlang/        # Punto de entrada (main.go) y REPL
├── internal/
│   ├── ast/            # Definición del árbol de sintaxis abstracta
│   ├── evaluator/       # Evaluador de sentencias y expresiones
│   ├── explain/         # Lógica del comando EXPLAIN
│   ├── lexer/           # Analizador léxico
│   ├── object/          # Sistema de objetos y entorno de ejecución
│   ├── parser/          # Analizador sintáctico
│   └── token/           # Definición de tokens y palabras clave
└── go.mod
```

## Ejecutar los tests

El proyecto incluye pruebas unitarias para el lexer, el parser y el evaluador:

```bash
go test ./...
```
