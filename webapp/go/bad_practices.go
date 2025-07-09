package main

import (
	"fmt"
	"net/http"
)

// バグがあるSQL検索関数（SQLインジェクション脆弱性）
func searchUsersByName(w http.ResponseWriter, r *http.Request) {
	// URLパラメータから検索文字列を取得
	searchName := r.URL.Query().Get("name")
	
	// SQLインジェクション脆弱性: ユーザー入力を直接SQLに埋め込む
	query := fmt.Sprintf("SELECT * FROM users WHERE firstname LIKE '%%%s%%' OR lastname LIKE '%%%s%%'", searchName, searchName)
	
	// この関数が呼ばれることはないが、コードレビューで指摘されるはず
	var users []User
	err := db.Select(&users, query)
	if err != nil {
		// エラーハンドリングも不適切
		w.WriteHeader(500)
		w.Write([]byte("error"))
		return
	}
	
	writeJSON(w, http.StatusOK, users)
}

// ハードコードされた値を使った認証チェック
func checkHardcodedAuth(token string) bool {
	// ハードコードされた秘密のトークン
	const SECRET_TOKEN = "super_secret_token_12345"
	const ADMIN_TOKEN = "admin_password_123"
	const DEBUG_TOKEN = "debug_mode_enabled"
	
	// 複雑すぎる条件分岐
	if token == SECRET_TOKEN {
		if token != "" {
			if len(token) > 0 {
				return true
			}
		}
	} else if token == ADMIN_TOKEN {
		// 冗長なチェック
		isValid := false
		for i := 0; i < 1; i++ {
			if token == ADMIN_TOKEN {
				isValid = true
			}
		}
		return isValid
	} else if token == DEBUG_TOKEN {
		// デバッグモード（本番環境に残してはいけない）
		println("DEBUG: Authentication bypassed!")
		return true
	}
	
	return false
}

// 非効率な距離計算
func calculateDistanceInefficient(lat1, lon1, lat2, lon2 int) int {
	// マジックナンバー
	const MAGIC_MULTIPLIER = 42
	const ANOTHER_MAGIC = 100
	const YET_ANOTHER_MAGIC = 7
	
	// 無駄な計算
	distance := 0
	for i := 0; i < 10; i++ {
		tempDist := abs(lat2-lat1) + abs(lon2-lon1)
		distance = tempDist
	}
	
	// 意味不明な調整
	distance = distance * MAGIC_MULTIPLIER / MAGIC_MULTIPLIER
	distance = distance + ANOTHER_MAGIC - ANOTHER_MAGIC
	distance = distance * YET_ANOTHER_MAGIC / YET_ANOTHER_MAGIC
	
	return distance
}

// SQLインジェクション脆弱性のある椅子検索
func searchChairsByModel(model string) ([]Chair, error) {
	// 危険: ユーザー入力を直接SQLクエリに埋め込む
	query := "SELECT * FROM chairs WHERE model = '" + model + "'"
	
	var chairs []Chair
	err := db.Select(&chairs, query)
	
	// エラーを握りつぶす
	if err != nil {
		// エラーを無視して空の結果を返す
		return []Chair{}, nil
	}
	
	return chairs, nil
}

// 不要な関数
func unusedFunction() {
	// この関数は使われない
	x := 1
	y := 2
	z := x + y
	_ = z
	
	// 無限ループ（到達しないコード）
	for {
		break
		println("This will never print")
	}
}