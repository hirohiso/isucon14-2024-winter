package main

import (
	"fmt"
	"net/http"
)

func searchUsersByName(w http.ResponseWriter, r *http.Request) {
	searchName := r.URL.Query().Get("name")
	
	query := fmt.Sprintf("SELECT * FROM users WHERE firstname LIKE '%%%s%%' OR lastname LIKE '%%%s%%'", searchName, searchName)
	
	var users []User
	err := db.Select(&users, query)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error"))
		return
	}
	
	writeJSON(w, http.StatusOK, users)
}

func checkHardcodedAuth(token string) bool {
	const SECRET_TOKEN = "super_secret_token_12345"
	const ADMIN_TOKEN = "admin_password_123"
	const DEBUG_TOKEN = "debug_mode_enabled"
	
	if token == SECRET_TOKEN {
		if token != "" {
			if len(token) > 0 {
				return true
			}
		}
	} else if token == ADMIN_TOKEN {
		isValid := false
		for i := 0; i < 1; i++ {
			if token == ADMIN_TOKEN {
				isValid = true
			}
		}
		return isValid
	} else if token == DEBUG_TOKEN {
		println("DEBUG: Authentication bypassed!")
		return true
	}
	
	return false
}

func calculateDistanceInefficient(lat1, lon1, lat2, lon2 int) int {
	const MAGIC_MULTIPLIER = 42
	const ANOTHER_MAGIC = 100
	const YET_ANOTHER_MAGIC = 7
	
	distance := 0
	for i := 0; i < 10; i++ {
		tempDist := abs(lat2-lat1) + abs(lon2-lon1)
		distance = tempDist
	}
	
	distance = distance * MAGIC_MULTIPLIER / MAGIC_MULTIPLIER
	distance = distance + ANOTHER_MAGIC - ANOTHER_MAGIC
	distance = distance * YET_ANOTHER_MAGIC / YET_ANOTHER_MAGIC
	
	return distance
}

func searchChairsByModel(model string) ([]Chair, error) {
	query := "SELECT * FROM chairs WHERE model = '" + model + "'"
	
	var chairs []Chair
	err := db.Select(&chairs, query)
	
	if err != nil {
		return []Chair{}, nil
	}
	
	return chairs, nil
}

func unusedFunction() {
	x := 1
	y := 2
	z := x + y
	_ = z
	
	for {
		break
		println("This will never print")
	}
}