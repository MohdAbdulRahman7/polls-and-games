package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type Poll struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	CreatedAt   string    `json:"created_at"`
	Options     []Option  `json:"options"`
	VoteCount   int       `json:"vote_count"`
}

type Option struct {
	ID       int    `json:"id"`
	PollID   int    `json:"poll_id"`
	Text     string `json:"text"`
	VoteCount int   `json:"vote_count"`
}

type Bookmark struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PollID int `json:"poll_id"`
}

type Vote struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	PollID   int `json:"poll_id"`
	OptionID int `json:"option_id"`
}

type Database struct {
	db *sql.DB
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./polls.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create tables
	createTables()
}

func createTables() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	pollsTable := `
	CREATE TABLE IF NOT EXISTS polls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	optionsTable := `
	CREATE TABLE IF NOT EXISTS options (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		poll_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE
	);`

	votesTable := `
	CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		poll_id INTEGER NOT NULL,
		option_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
		FOREIGN KEY (option_id) REFERENCES options(id) ON DELETE CASCADE,
		UNIQUE(user_id, poll_id)
	);`

	bookmarksTable := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		poll_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
		UNIQUE(user_id, poll_id)
	);`

	db.Exec(usersTable)
	db.Exec(pollsTable)
	db.Exec(optionsTable)
	db.Exec(votesTable)
	db.Exec(bookmarksTable)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.WriteHeader(http.StatusOK)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		req.Username, req.Email, string(hashedPassword),
	)
	if err != nil {
		sendError(w, "Username or email already exists", http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()
	user := User{
		ID:       int(userID),
		Username: req.Username,
		Email:    req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	var password string
	err := db.QueryRow(
		"SELECT id, username, email, password FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &password)

	if err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password))
	if err != nil {
		sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getPollsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	rows, err := db.Query(`
		SELECT p.id, p.title, p.description, p.user_id, p.created_at, u.username,
			COUNT(DISTINCT v.id) as vote_count
		FROM polls p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN votes v ON p.id = v.poll_id
		GROUP BY p.id
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var polls []Poll
	for rows.Next() {
		var poll Poll
		err := rows.Scan(&poll.ID, &poll.Title, &poll.Description, &poll.UserID, &poll.CreatedAt, &poll.Username, &poll.VoteCount)
		if err != nil {
			continue
		}

		// Get options for each poll
		optionRows, _ := db.Query("SELECT id, text FROM options WHERE poll_id = ?", poll.ID)
		for optionRows.Next() {
			var option Option
			optionRows.Scan(&option.ID, &option.Text)
			poll.Options = append(poll.Options, option)
		}
		optionRows.Close()

		polls = append(polls, poll)
	}

	// Get total count
	var total int
	db.QueryRow("SELECT COUNT(*) FROM polls").Scan(&total)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"polls": polls,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func getPollHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	vars := mux.Vars(r)
	pollID := vars["id"]

	var poll Poll
	err := db.QueryRow(`
		SELECT p.id, p.title, p.description, p.user_id, p.created_at, u.username
		FROM polls p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, pollID).Scan(&poll.ID, &poll.Title, &poll.Description, &poll.UserID, &poll.CreatedAt, &poll.Username)

	if err != nil {
		sendError(w, "Poll not found", http.StatusNotFound)
		return
	}

	// Get options with vote counts
	optionRows, _ := db.Query(`
		SELECT o.id, o.text, COUNT(v.id) as vote_count
		FROM options o
		LEFT JOIN votes v ON o.id = v.option_id
		WHERE o.poll_id = ?
		GROUP BY o.id
	`, pollID)
	defer optionRows.Close()

	for optionRows.Next() {
		var option Option
		optionRows.Scan(&option.ID, &option.Text, &option.VoteCount)
		poll.Options = append(poll.Options, option)
	}

	// Get total vote count
	db.QueryRow("SELECT COUNT(*) FROM votes WHERE poll_id = ?", pollID).Scan(&poll.VoteCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
}

func createPollHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		UserID      int      `json:"user_id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Options     []string `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := tx.Exec(
		"INSERT INTO polls (title, description, user_id) VALUES (?, ?, ?)",
		req.Title, req.Description, req.UserID,
	)
	if err != nil {
		tx.Rollback()
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pollID, _ := result.LastInsertId()

	for _, optionText := range req.Options {
		_, err = tx.Exec(
			"INSERT INTO options (poll_id, text) VALUES (?, ?)",
			pollID, optionText,
		)
		if err != nil {
			tx.Rollback()
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	var poll Poll
	db.QueryRow(`
		SELECT p.id, p.title, p.description, p.user_id, p.created_at, u.username
		FROM polls p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, pollID).Scan(&poll.ID, &poll.Title, &poll.Description, &poll.UserID, &poll.CreatedAt, &poll.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(poll)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		UserID   int `json:"user_id"`
		PollID   int `json:"poll_id"`
		OptionID int `json:"option_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already voted
	var existingVote int
	err := db.QueryRow(
		"SELECT id FROM votes WHERE user_id = ? AND poll_id = ?",
		req.UserID, req.PollID,
	).Scan(&existingVote)

	if err == nil {
		// Update existing vote
		_, err = db.Exec(
			"UPDATE votes SET option_id = ? WHERE user_id = ? AND poll_id = ?",
			req.OptionID, req.UserID, req.PollID,
		)
	} else {
		// Create new vote
		_, err = db.Exec(
			"INSERT INTO votes (user_id, poll_id, option_id) VALUES (?, ?, ?)",
			req.UserID, req.PollID, req.OptionID,
		)
	}

	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Vote recorded"})
}

func bookmarkHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req struct {
		UserID int `json:"user_id"`
		PollID int `json:"poll_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		_, err := db.Exec(
			"INSERT OR IGNORE INTO bookmarks (user_id, poll_id) VALUES (?, ?)",
			req.UserID, req.PollID,
		)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Bookmarked"})
	} else if r.Method == "DELETE" {
		_, err := db.Exec(
			"DELETE FROM bookmarks WHERE user_id = ? AND poll_id = ?",
			req.UserID, req.PollID,
		)
		if err != nil {
			sendError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Unbookmarked"})
	}
}

func getBookmarksHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	vars := mux.Vars(r)
	userID := vars["user_id"]

	rows, err := db.Query(`
		SELECT p.id, p.title, p.description, p.user_id, p.created_at, u.username,
			COUNT(DISTINCT v.id) as vote_count
		FROM bookmarks b
		JOIN polls p ON b.poll_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN votes v ON p.id = v.poll_id
		WHERE b.user_id = ?
		GROUP BY p.id
		ORDER BY b.created_at DESC
	`, userID)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var polls []Poll
	for rows.Next() {
		var poll Poll
		err := rows.Scan(&poll.ID, &poll.Title, &poll.Description, &poll.UserID, &poll.CreatedAt, &poll.Username, &poll.VoteCount)
		if err != nil {
			continue
		}

		optionRows, _ := db.Query("SELECT id, text FROM options WHERE poll_id = ?", poll.ID)
		for optionRows.Next() {
			var option Option
			optionRows.Scan(&option.ID, &option.Text)
			poll.Options = append(poll.Options, option)
		}
		optionRows.Close()

		polls = append(polls, poll)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
}

func getUserPollsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	vars := mux.Vars(r)
	userID := vars["user_id"]

	rows, err := db.Query(`
		SELECT p.id, p.title, p.description, p.user_id, p.created_at, u.username,
			COUNT(DISTINCT v.id) as vote_count
		FROM polls p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN votes v ON p.id = v.poll_id
		WHERE p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, userID)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var polls []Poll
	for rows.Next() {
		var poll Poll
		err := rows.Scan(&poll.ID, &poll.Title, &poll.Description, &poll.UserID, &poll.CreatedAt, &poll.Username, &poll.VoteCount)
		if err != nil {
			continue
		}

		optionRows, _ := db.Query("SELECT id, text FROM options WHERE poll_id = ?", poll.ID)
		for optionRows.Next() {
			var option Option
			optionRows.Scan(&option.ID, &option.Text)
			poll.Options = append(poll.Options, option)
		}
		optionRows.Close()

		polls = append(polls, poll)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
}

func deletePollHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	pollID := vars["id"]

	_, err := db.Exec("DELETE FROM polls WHERE id = ?", pollID)
	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Poll deleted"})
}

func checkBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	userID := r.URL.Query().Get("user_id")
	pollID := r.URL.Query().Get("poll_id")

	var exists int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM bookmarks WHERE user_id = ? AND poll_id = ?",
		userID, pollID,
	).Scan(&exists)

	if err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"bookmarked": exists > 0})
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()

	// Auth routes
	r.HandleFunc("/api/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", loginHandler).Methods("POST", "OPTIONS")

	// Poll routes
	r.HandleFunc("/api/polls", getPollsHandler).Methods("GET")
	r.HandleFunc("/api/polls", createPollHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/polls/{id}", getPollHandler).Methods("GET")
	r.HandleFunc("/api/polls/{id}", deletePollHandler).Methods("DELETE", "OPTIONS")

	// Vote routes
	r.HandleFunc("/api/vote", voteHandler).Methods("POST", "OPTIONS")

	// Bookmark routes
	r.HandleFunc("/api/bookmark", bookmarkHandler).Methods("POST", "DELETE", "OPTIONS")
	r.HandleFunc("/api/bookmarks/{user_id}", getBookmarksHandler).Methods("GET")
	r.HandleFunc("/api/check-bookmark", checkBookmarkHandler).Methods("GET")

	// User routes
	r.HandleFunc("/api/user/{user_id}/polls", getUserPollsHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

