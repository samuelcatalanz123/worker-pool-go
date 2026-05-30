# Diseño: Cola de tareas con workers (Go)

**Fecha:** 2026-05-30 · **Estado:** Aprobado · **Autor:** Samuel (10º proyecto)

## Objetivo

Una app web que procesa tareas en segundo plano con un **worker pool**: un número
fijo de "trabajadores" (goroutines) toman tareas de una cola (channel) y las
procesan. Objetivo de aprendizaje: el patrón **worker pool** (goroutines + channels).

## Cómo funciona

```
Envías tareas → cola (channel) → N workers las toman y procesan a la vez
Cada tarea:  ⏳ pendiente → 🔄 procesando (worker #N) → ✅ hecho
```

La página se **refresca sola** cada segundo para ver el progreso en vivo.

## Arquitectura

```
worker-pool-go/
  main.go                 arranca el servidor (:8080 o PORT)
  internal/pool/
    pool.go               Pool, Job, Status; New (lanza N workers); Submit; Jobs
    pool_test.go          prueba: se envían 5 tareas y todas terminan "hecho"
  internal/web/
    handler.go            GET / (panel, autorefresco) + POST /submit
    templates/index.html
    static/style.css
  README.md
```

- **pool.go:**
  - `Status` = `pendiente` | `procesando` | `hecho`.
  - `Job{ID, Name, Status, Worker}`.
  - `Pool` con `sync.Mutex`, lista de jobs, `queue chan *Job`, nº de workers y
    duración del "trabajo".
  - `New(workers int, dur time.Duration)` — lanza `workers` goroutines que leen
    de `queue`, marcan "procesando", esperan `dur` (simulan trabajo) y marcan "hecho".
  - `Submit(name)` — crea el Job (pendiente), lo guarda y lo mete en la cola.
  - `Jobs()` — copia de las tareas (más reciente primero), protegida con el Mutex.
  - `Workers()` — nº de workers.
- **web/handler.go:** crea `pool.New(3, 2s)`; `GET /` muestra los workers y las
  tareas; `POST /submit` envía una tarea nueva (redirige a /).
- **templates/index.html:** `<meta http-equiv="refresh" content="1">` (autorefresco),
  formulario para enviar tareas, lista con el estado (color) y el worker.

## Pruebas

- **pool_test.go:** `New(3, 5ms)`, se envían 5 tareas; se espera (con un límite de
  ~2 s) a que **todas** queden en `hecho`. La duración corta hace la prueba rápida.
- `go build/vet/test` limpios.

## Fuera de alcance (YAGNI)

Guardar tareas en disco, reintentos, prioridades, cancelar tareas.

## Criterios de éxito

1. `go run .` sirve el panel en http://localhost:8080.
2. Al enviar tareas, se ven pasar de pendiente → procesando (con su worker) → hecho.
3. Nunca hay más de N tareas "procesando" a la vez (el pool limita).
4. La prueba del pool pasa (todas las tareas terminan).
