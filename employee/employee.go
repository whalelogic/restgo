package employee

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type Employee struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Department string  `json:"department"`
	Role       string  `json:"role"`
	Salary     float64 `json:"salary"`
}

type EmployeeHandler struct {
	Rdb *redis.Client
}

const employeeSet = "employees"

func empKey(id string) string {
	return "employee:" + id
}

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ids, err := h.Rdb.SMembers(ctx, employeeSet).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employees := make([]Employee, 0, len(ids))
	for _, id := range ids {
		val, err := h.Rdb.Get(ctx, empKey(id)).Result()
		if err != nil {
			continue
		}
		var emp Employee
		if err := json.Unmarshal([]byte(val), &emp); err == nil {
			employees = append(employees, emp)
		}
	}

	data, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.PathValue("id")

	val, err := h.Rdb.Get(ctx, empKey(id)).Result()
	if err == redis.Nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if emp.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.Rdb.Set(ctx, empKey(emp.ID), data, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.Rdb.SAdd(ctx, employeeSet, emp.ID).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.PathValue("id")

	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	emp.ID = id

	data, err := json.Marshal(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.Rdb.Set(ctx, empKey(id), data, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.Rdb.SAdd(ctx, employeeSet, id).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmployeeHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.PathValue("id")

	val, err := h.Rdb.Get(ctx, empKey(id)).Result()
	if err == redis.Nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var existing map[string]any
	if err := json.Unmarshal([]byte(val), &existing); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for k, v := range patch {
		if k != "id" {
			existing[k] = v
		}
	}

	data, err := json.Marshal(existing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.Rdb.Set(ctx, empKey(id), data, 0).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.PathValue("id")

	deleted, err := h.Rdb.Del(ctx, empKey(id)).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err := h.Rdb.SRem(ctx, employeeSet, id).Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
