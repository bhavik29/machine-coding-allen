package main

import (
	"allen-machine-coding/controllers/deal"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/create-deal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}

		deal := deal.Deal{}
		err = json.Unmarshal(respBody, &deal)
		if err != nil {
			http.Error(w, "corrupted request body", http.StatusBadRequest)
			return
		}

		// validation for MaxNumberOfItems
		if deal.MaxNumberOfItems < 10 {
			http.Error(w, "max number of items cannot be less than 10", http.StatusBadRequest)
			return
		}

		// validation for Duration
		currentTime := time.Now()
		dealEndTime := currentTime.Add(deal.Duration * time.Second)
		if currentTime.Add(30 * time.Minute).After(dealEndTime) {
			http.Error(w, "cannot run the deal for less than 30 minutes", http.StatusBadRequest)
			return
		}

		ctx, cancelFn := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancelFn()
		dealID, err := deal.CreateDeal(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot complete request: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(
			fmt.Sprintf("{\"dealID\": \"%v\"}", dealID),
		))
	})

	mux.HandleFunc("/end-deal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		respBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read request body", http.StatusBadRequest)
			return
		}

		deal := deal.Deal{}
		err = json.Unmarshal(respBody, &deal)
		if err != nil {
			http.Error(w, "corrupted request body", http.StatusBadRequest)
			return
		}

		ctx, cancelFn := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancelFn()
		err = deal.EndDeal(ctx, deal.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot complete request: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write([]byte(
			fmt.Sprintf("{\"status\": \"deal ended\"}"),
		))
	})

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic("cannot start the server")
	}
}
