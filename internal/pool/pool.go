// Package pool procesa tareas con un número fijo de "trabajadores" (worker pool).
package pool

import (
	"sync"
	"time"
)

// Status es en qué punto está una tarea.
type Status string

const (
	Pending Status = "pendiente"
	Working Status = "procesando"
	Done    Status = "hecho"
)

// Job es una tarea.
type Job struct {
	ID     int
	Name   string
	Status Status
	Worker int // qué trabajador la procesó (0 = aún ninguno)
}

// Pool es la cola de tareas y sus trabajadores.
type Pool struct {
	mu      sync.Mutex
	jobs    []*Job
	queue   chan *Job // la "cola": un channel del que leen los workers
	nextID  int
	workers int
	dur     time.Duration
}

// New crea el pool y lanza `workers` goroutines trabajadoras.
// `dur` es cuánto tarda en procesarse cada tarea (simulado).
func New(workers int, dur time.Duration) *Pool {
	p := &Pool{queue: make(chan *Job, 100), workers: workers, dur: dur}
	for i := 1; i <= workers; i++ {
		go p.worker(i)
	}
	return p
}

// Workers devuelve cuántos trabajadores hay.
func (p *Pool) Workers() int { return p.workers }

// worker es un trabajador: toma tareas de la cola y las procesa, una tras otra.
func (p *Pool) worker(id int) {
	for job := range p.queue {
		p.setStatus(job, Working, id)
		time.Sleep(p.dur) // simula el trabajo
		p.setStatus(job, Done, id)
	}
}

// Submit crea una tarea nueva (pendiente) y la mete en la cola.
func (p *Pool) Submit(name string) {
	p.mu.Lock()
	p.nextID++
	job := &Job{ID: p.nextID, Name: name, Status: Pending}
	p.jobs = append(p.jobs, job)
	p.mu.Unlock()

	p.queue <- job // entra a la cola; algún worker la tomará
}

// setStatus cambia el estado de una tarea de forma segura (con el Mutex).
func (p *Pool) setStatus(job *Job, s Status, worker int) {
	p.mu.Lock()
	job.Status = s
	if worker > 0 {
		job.Worker = worker
	}
	p.mu.Unlock()
}

// Jobs devuelve una copia de las tareas, la más reciente primero.
func (p *Pool) Jobs() []Job {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]Job, 0, len(p.jobs))
	for i := len(p.jobs) - 1; i >= 0; i-- {
		out = append(out, *p.jobs[i])
	}
	return out
}
