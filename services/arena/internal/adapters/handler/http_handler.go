package handler

import (
	"api/services/arena/internal/core/ports"
	"encoding/json"
	"net/http"
	"strconv"
)

type HttpHandler struct {
	service ports.ArenaService
}

func NewHttpHandler(s ports.ArenaService) *HttpHandler {
	return &HttpHandler{service: s}
}

func (h *HttpHandler) HandleDuel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	var req struct {
		F1 string `json:"fighter_1"`
		F2 string `json:"fighter_2"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	result, err := h.service.Duel(req.F1, req.F2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *HttpHandler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. อ่านค่า Query Parameter
	query := r.URL.Query()

	// อ่าน limit (Default เป็น 0 ถ้าไม่ส่งมา)
	limitStr := query.Get("limit")
	limit, _ := strconv.Atoi(limitStr)

	// อ่าน fighter_id
	fighterID := query.Get("fighter_id")

	// 2. เรียก Service พร้อม parameter
	history, err := h.service.GetHistory(limit, fighterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
