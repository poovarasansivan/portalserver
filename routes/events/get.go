package events

import (
	"database/sql"
	"encoding/json"
	"learnathon/config"
	"learnathon/function"
	"log"
	"net/http"
)

func GetAllEvents(w http.ResponseWriter, r *http.Request) {
	var response map[string]interface{}
	var eventdata []Events
	var temp Events

	row, err := config.Database.Query("SELECT me.event_name,me.description,me.event_date,mu.name FROM m_events me INNER JOIN m_users mu ON mu.id = me.incharge WHERE me.status ='1'")

	if err != nil {
		if err == sql.ErrNoRows {
			response = map[string]interface{}{
				"success": false,
				"error":   "No Request",
			}
		} else {
			response = map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			}
		}
		function.Response(w, response)
		return
	}

	for row.Next() {
		err := row.Scan(&temp.EventName, &temp.Description, &temp.EventDate, &temp.Incharge)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tempRow := Events{
			EventName:   temp.EventName,
			Description: temp.Description,
			Incharge:    temp.Incharge,
			EventDate:   temp.EventDate,
		}
		eventdata = append(eventdata, tempRow)
	}

	response = map[string]interface{}{
		"success": true,
		"data":    eventdata,
	}
	function.Response(w, response)
}

func GetAvailableId(w http.ResponseWriter,r *http.Request){
	rows, err := config.Database.Query("SELECT id,event_name FROM m_events WHERE STATUS = '1'")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	var events [] GetEvents
	for rows.Next() {
		var event GetEvents
		if err := rows.Scan(&event.Id, &event.Eventname); err != nil {
			http.Error(w, "Error scanning database result", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, event)
	}
	response := struct {
		Events []GetEvents `json:"events"`
	}{Events: events}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}