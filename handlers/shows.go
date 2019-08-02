package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mattrout92/agility-backend/logger"
)

// Show ...
type Show struct {
	ID           int     `json:"id"`
	Date         string  `json:"date"`
	Location     string  `json:"location"`
	EntriesClose string  `json:"entries_close"`
	JudgingFrom  string  `json:"judging_from"`
	Classes      []Class `json:"classes"`
}

// Class ...
type Class struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DogSizes string `json:"dog_sizes"`
	Grades   string `json:"grades"`
}

// GetShows ...
func (svc *Service) GetShows(w http.ResponseWriter, req *http.Request) {
	db := svc.Store.SQLX()

	rows, err := db.Query(svc.Store.GetQuery("GetOpenShows"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error(err)
		return
	}

	defer rows.Close()

	var shows []Show

	for rows.Next() {
		var show Show

		err := rows.Scan(&show.ID, &show.Date, &show.Location, &show.JudgingFrom, &show.EntriesClose)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		classRows, err := db.Query(svc.Store.GetQuery("GetShowClasses"), show.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(err)
			return
		}

		for classRows.Next() {
			var class Class

			classRows.Scan(&class.ID, &class.Name, &class.DogSizes, &class.Grades)

			show.Classes = append(show.Classes, class)
		}

		shows = append(shows, show)
	}

	json.NewEncoder(w).Encode(shows)
}
