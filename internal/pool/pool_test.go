package pool

import (
	"fmt"
	"testing"
	"time"
)

// TestPoolProcessesAllJobs comprueba que el pool termina TODAS las tareas.
func TestPoolProcessesAllJobs(t *testing.T) {
	p := New(3, 5*time.Millisecond) // 3 workers, trabajo cortito para que la prueba sea rápida
	for i := 1; i <= 5; i++ {
		p.Submit(fmt.Sprintf("tarea %d", i))
	}

	// Esperamos a que todas queden en "hecho" (máximo ~2 segundos).
	deadline := time.Now().Add(2 * time.Second)
	for {
		hechas := 0
		for _, j := range p.Jobs() {
			if j.Status == Done {
				hechas++
			}
		}
		if hechas == 5 {
			return // ✅ todas terminaron
		}
		if time.Now().After(deadline) {
			t.Fatalf("no todas las tareas terminaron a tiempo (hechas = %d de 5)", hechas)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestWorkers comprueba el número de trabajadores.
func TestWorkers(t *testing.T) {
	p := New(4, time.Millisecond)
	if p.Workers() != 4 {
		t.Errorf("Workers() = %d, esperaba 4", p.Workers())
	}
}
