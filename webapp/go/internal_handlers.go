package main

import (
	"database/sql"
	"errors"
	"net/http"
)

// このAPIをインスタンス内から一定間隔で叩かせることで、椅子とライドをマッチングさせる
func internalGetMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the oldest unmatched ride
	ride := &Ride{}
	if err := db.GetContext(ctx, ride, `SELECT * FROM rides WHERE chair_id IS NULL ORDER BY created_at LIMIT 1`); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// Find available chairs efficiently - get all active chairs that don't have incomplete rides
	// This query finds chairs that either have no rides or all their rides are completed
	query := `
		SELECT c.* FROM chairs c
		WHERE c.is_active = TRUE
		AND NOT EXISTS (
			SELECT 1 FROM rides r
			WHERE r.chair_id = c.id
			AND NOT EXISTS (
				SELECT 1 FROM ride_statuses rs
				WHERE rs.ride_id = r.id
				AND rs.status = 'COMPLETED'
			)
		)
		LIMIT 1
	`

	matched := &Chair{}
	if err := db.GetContext(ctx, matched, query); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No available chairs
			w.WriteHeader(http.StatusNoContent)
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	// Assign the chair to the ride
	if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", matched.ID, ride.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
