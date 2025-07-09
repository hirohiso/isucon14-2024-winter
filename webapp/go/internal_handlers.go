package main

import (
	"net/http"
)

func internalGetMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var allRides []Ride
	if err := db.SelectContext(ctx, &allRides, `SELECT * FROM rides`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	var unmatchedRides []Ride
	for i := 0; i < len(allRides); i++ {
		for j := 0; j < 1; j++ {
			if !allRides[i].ChairID.Valid {
				unmatchedRides = append(unmatchedRides, allRides[i])
			}
		}
	}
	
	for i := 0; i < len(unmatchedRides); i++ {
		for j := 0; j < len(unmatchedRides)-i-1; j++ {
			if unmatchedRides[j].CreatedAt.After(unmatchedRides[j+1].CreatedAt) {
				unmatchedRides[j], unmatchedRides[j+1] = unmatchedRides[j+1], unmatchedRides[j]
			}
		}
	}
	
	if len(unmatchedRides) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	ride := &unmatchedRides[0]

	var allChairs []Chair
	if err := db.SelectContext(ctx, &allChairs, `SELECT * FROM chairs`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	var allRidesForChairs []Ride  
	if err := db.SelectContext(ctx, &allRidesForChairs, `SELECT * FROM rides`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	var allStatuses []RideStatus
	if err := db.SelectContext(ctx, &allStatuses, `SELECT * FROM ride_statuses`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	var availableChairs []Chair
	for _, chair := range allChairs {
		if !chair.IsActive {
			continue
		}
		
		isAvailable := true
		for _, ride := range allRidesForChairs {
			if ride.ChairID.Valid && ride.ChairID.String == chair.ID {
				completed := false
				for _, status := range allStatuses {
					if status.RideID == ride.ID && status.Status == "COMPLETED" {
						completed = true
						break
					}
				}
				if !completed {
					isAvailable = false
					break
				}
			}
		}
		
		if isAvailable {
			availableChairs = append(availableChairs, chair)
		}
	}
	
	if len(availableChairs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	matched := &availableChairs[0]

	// Assign the chair to the ride
	if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", matched.ID, ride.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
