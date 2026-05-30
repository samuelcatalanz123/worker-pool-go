// Package web sirve el panel del worker pool.
package web

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/samuelcatalanz123/worker-pool-go/internal/pool"
)

//go:embed templates/*.html static/*
var files embed.FS

// Handler sirve el panel.
type Handler struct {
	tmpl *template.Template
	pool *pool.Pool
}

// New crea el Handler con un pool de 3 trabajadores (cada tarea tarda 2 s).
func New() (*Handler, error) {
	tmpl, err := template.ParseFS(files, "templates/*.html")
	if err != nil {
		return nil, err
	}
	p := pool.New(3, 2*time.Second)
	return &Handler{tmpl: tmpl, pool: p}, nil
}

// Routes monta las rutas.
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(files))
	mux.HandleFunc("GET /{$}", h.home)
	mux.HandleFunc("POST /submit", h.submit)
	return mux
}

type viewData struct {
	Workers int
	Jobs    []pool.Job
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	data := viewData{Workers: h.pool.Workers(), Jobs: h.pool.Jobs()}
	if err := h.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "error del servidor", http.StatusInternalServerError)
	}
}

func (h *Handler) submit(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		name = "Tarea sin nombre"
	}
	h.pool.Submit(name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
