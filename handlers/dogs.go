package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mattrout92/agility-backend/logger"
)

// AddDogInput ...
type AddDogInput struct {
	Name    string `json:"name"`
	Height  string `json:"height"`
	Grade   string `json:"grade"`
	Handler string `json:"handler"`
	UserID  int    `json:"user_id"`
	ID      int64  `json:"id"`
}

// AddDog ...
func (svc *Service) AddDog(w http.ResponseWriter, req *http.Request) {
	var input AddDogInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	defer req.Body.Close()

	db := svc.Store.SQLX()

	info, err := db.Exec(svc.Store.GetQuery("InsertDog"), input.Name, input.Height, input.Grade, input.Handler, input.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	input.ID, _ = info.LastInsertId()

	json.NewEncoder(w).Encode(input)
}

// GetDogs ...
func (svc *Service) GetDogs(w http.ResponseWriter, req *http.Request) {
	var data []AddDogInput

	userID := req.URL.Query().Get("user_id")

	db := svc.Store.SQLX()

	rows, err := db.Query(svc.Store.GetQuery("GetDogs"), userID)
	if err != nil {
		logger.Error(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var item AddDogInput
		rows.Scan(&item.Name, &item.Height, &item.Grade, &item.Handler, &item.UserID)

		data = append(data, item)
	}

	json.NewEncoder(w).Encode(data)
}

// DeleteDog ...
func (svc *Service) DeleteDog(w http.ResponseWriter, req *http.Request) {
	var input AddDogInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	defer req.Body.Close()

	db := svc.Store.SQLX()

	info, err := db.Exec(svc.Store.GetQuery("DeleteDog"), input.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	input.ID, _ = info.LastInsertId()

	json.NewEncoder(w).Encode(input)
}

// UpdateDog ...
func (svc *Service) UpdateDog(w http.ResponseWriter, req *http.Request) {
	var input AddDogInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	defer req.Body.Close()

	db := svc.Store.SQLX()

	_, err := db.Exec(svc.Store.GetQuery("UpdateDog"), input.Name, input.Height, input.Grade, input.Handler, input.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	json.NewEncoder(w).Encode(input)
}
