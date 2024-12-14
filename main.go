package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/onfirebyte/optiseat/internal/optimize"
)

type OptimizeReq struct {
	MaxSeat int          `json:"maxSeat"`
	People  []string     `json:"people"`
	Prefs   []Preference `json:"preference"`
}

type Preference struct {
	A     string `json:"a"`
	B     string `json:"b"`
	Score int    `json:"score"`
}

type OptimizeResp struct {
	Score    int        `json:"score"`
	SeatPlan [][]string `json:"seatPlan"`
}

func main() {
	http.HandleFunc("OPTION /optimize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	})

	http.HandleFunc("POST /optimize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		var req OptimizeReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Warn("invalid input from user", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				slog.Warn("error while sending response", "error", err)
			}
		}

		peopleIds := map[string]int{}
		for i, name := range req.People {
			peopleIds[name] = i
		}

		pref := make([][]int, len(req.People))
		for i := range len(pref) {
			pref[i] = make([]int, len(req.People))
		}

		for _, p := range req.Prefs {
			pref[peopleIds[p.A]][peopleIds[p.B]] = p.Score
		}

		seatPlanByID, happiness := optimize.SimulatedAnnealing(pref, len(req.People), req.MaxSeat)

		seatPlanByName := make([][]string, 0, len(seatPlanByID))
		for _, table := range seatPlanByID {
			tableByName := make([]string, 0, len(table))
			for _, person := range table {
				tableByName = append(tableByName, req.People[person])
			}
			seatPlanByName = append(seatPlanByName, tableByName)
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(
			OptimizeResp{
				Score:    happiness,
				SeatPlan: seatPlanByName,
			},
		)
		if err != nil {
			slog.Warn("invalid input from user", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				slog.Warn("error while sending response", "error", err)
			}
		}
	})

	log.Println("Serving at port 3000")
	log.Fatal(http.ListenAndServe(":3000", http.DefaultServeMux))

}
