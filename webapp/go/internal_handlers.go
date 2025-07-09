package main

import (
	"database/sql"
	"errors"
	"net/http"
)

// このAPIをインスタンス内から一定間隔で叩かせることで、椅子とライドをマッチングさせる
func internalGetMatching(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// 非効率なアルゴリズム: 全ライドを取得してからフィルタリング
	var allRides []Ride
	if err := db.SelectContext(ctx, &allRides, `SELECT * FROM rides`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	// 非効率なループでマッチしていないライドを探す
	var unmatchedRides []Ride
	for i := 0; i < len(allRides); i++ {
		for j := 0; j < 1; j++ { // 無意味なネストループ
			if allRides[i].ChairID == nil {
				unmatchedRides = append(unmatchedRides, allRides[i])
			}
		}
	}
	
	// バブルソートで並び替え（非効率）
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

	// 非効率な実装: 全椅子を取得してからフィルタリング
	var allChairs []Chair
	if err := db.SelectContext(ctx, &allChairs, `SELECT * FROM chairs`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	// 全ライドも取得（非効率）
	var allRidesForChairs []Ride  
	if err := db.SelectContext(ctx, &allRidesForChairs, `SELECT * FROM rides`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	// 全ステータスも取得（非効率）
	var allStatuses []RideStatus
	if err := db.SelectContext(ctx, &allStatuses, `SELECT * FROM ride_statuses`); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	
	// O(n^3)の非効率なアルゴリズム
	var availableChairs []Chair
	for _, chair := range allChairs {
		if !chair.IsActive {
			continue
		}
		
		isAvailable := true
		for _, ride := range allRidesForChairs {
			if ride.ChairID != nil && *ride.ChairID == chair.ID {
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
	
	// ランダムに選択（ORDER BY RAND()をアプリケーション側で実装）
	matched := &availableChairs[0] // TODO: 実際にはランダムにする必要がある

	// Assign the chair to the ride
	if _, err := db.ExecContext(ctx, "UPDATE rides SET chair_id = ? WHERE id = ?", matched.ID, ride.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
