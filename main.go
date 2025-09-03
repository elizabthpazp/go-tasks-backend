package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
)

type Task struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Done  bool   `json:"done"`
}

var (
    tasks   []Task
    nextID  = 1
    tasksMu sync.Mutex
)

func createTask(w http.ResponseWriter, r *http.Request) {
    var t Task
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tasksMu.Lock()
    t.ID = nextID
    nextID++
    tasks = append(tasks, t)
    tasksMu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(t)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    tasksMu.Lock()
    defer tasksMu.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func main() {
    http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getTasks(w, r)
        case http.MethodPost:
            createTask(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}