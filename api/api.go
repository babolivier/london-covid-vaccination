package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/babolivier/london-covid-vaccination/config"
	"github.com/babolivier/london-covid-vaccination/storage"

	"github.com/sirupsen/logrus"
)

// makeApiHandler returns a HTTP handler that handlers incoming GET requests on the API.
func makeApiHandler(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve all stats.
		stats, err := db.GetAllStats()
		if err != nil {
			logrus.WithError(err).Error("Failed to get stats")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Turn the stats into JSON bytes.
		rawBody, err := json.Marshal(stats)
		if err != nil {
			logrus.WithError(err).Error("Failed to marshal response body")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Send the response.
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		if _, err = w.Write(rawBody); err != nil {
			logrus.WithError(err).Error("Failed to write response body")
		}
	}
}

// StartApiServer registers the API handlers, and starts the HTTP server
func StartApiServer(cfg *config.ApiConfig, db *storage.Database) {
	apiAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	logrus.Infof("Starting API server on %s", apiAddr)

	// Register a file server handler for the public directory to serve the front end.
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	// Register the API handler.
	http.Handle("/stats", makeApiHandler(db))
	// Start the server.
	if err := http.ListenAndServe(apiAddr, nil); err != nil {
		logrus.WithError(err).Error("API server stopped")
	}
}
