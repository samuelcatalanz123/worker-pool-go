# Worker Pool — Plan

**Goal:** App web en Go con un worker pool (N goroutines procesan tareas de una cola/channel).
**Tech:** Go (net/http, html/template, go:embed, sync, time — goroutines y channels). Sin deps externas.
Módulo: `github.com/samuelcatalanz123/worker-pool-go`. Commits autoría Samuel.

## Tareas
1. **pool (con prueba):** `internal/pool/pool.go` (Pool, Job, Status, New, Submit, Jobs, Workers) + `pool_test.go` (5 tareas → todas "hecho").
2. **web:** `internal/web/handler.go` (crea pool; GET / panel; POST /submit) + `templates/index.html` (autorefresco + form + lista con estados) + `static/style.css`.
3. **main.go:** servidor (:8080 o PORT).
4. **README + .gitignore.**
5. Verificar build/vet/test; arrancar y probar; commits.

Self-review: cobertura del spec ✔ · sin placeholders ✔ · tipos consistentes
(Pool, Job, Status pendiente/procesando/hecho, New, Submit, Jobs, Workers) ✔.
