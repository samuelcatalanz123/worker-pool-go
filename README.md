# Cola de tareas con Workers (Go)

App web que procesa tareas en segundo plano con un **worker pool**: un número fijo
de "trabajadores" (goroutines) toman tareas de una cola (channel) y las procesan.
Hecho en **Go**.

## Uso

```bash
go run .
```

Abre **http://localhost:8080**, envía tareas y míralas pasar de
⏳ pendiente → 🔄 procesando (con su worker) → ✅ hecho. La página se refresca sola.
Hay **3 trabajadores**, así que nunca se procesan más de 3 a la vez.

## Cómo funciona

- `internal/pool`: la cola es un `channel`; al crear el pool se lanzan N goroutines
  trabajadoras que leen tareas del channel y las procesan. El estado de cada tarea
  se protege con un `sync.Mutex`.
- `internal/web`: panel (`GET /`, se autorrefresca) y enviar tarea (`POST /submit`).

## Estructura

```
main.go                 arranque del servidor
internal/pool/          la cola + los workers (goroutines + channels) + prueba
internal/web/           el panel y el formulario
```

## Pruebas

```bash
go test -race ./...
```

La prueba envía varias tareas y comprueba que **todas** terminan. Se usa `-race`
para verificar que la concurrencia es segura.

## Stack

Go (net/http, html/template, go:embed, sync, goroutines, channels).
