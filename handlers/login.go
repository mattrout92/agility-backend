package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/mattrout92/agility-backend/logger"
)

// LoginInput contains fields required to login
type LoginInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// LoginOutput contains fields required from login output
type LoginOutput struct {
	ID      int64  `json:"id"            db:"id"`
	Email   string `json:"email"         db:"email"`
	Name    string `json:"name"          db:"name"`
	IsAdmin bool   `json:"is_admin"      db:"is_admin"`
}

// Login handles logging in of a user
func (svc *Service) Login(w http.ResponseWriter, req *http.Request) {
	var input LoginInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}
	defer req.Body.Close()

	if len(input.Email) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := svc.Store.SQLX()

	var output LoginOutput

	err := db.QueryRowx(svc.Store.GetQuery("GetUser"), input.Email).StructScan(&output)
	if err != nil {
		if err == sql.ErrNoRows {
			info, err := db.Exec(svc.Store.GetQuery("InsertUser"), input.Name, input.Email)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error(err)
				return
			}

			output.Name = input.Name
			output.Email = input.Email
			output.ID, _ = info.LastInsertId()
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(output)
}
