package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var jwtSecret = []byte(getEnv("JWT_SECRET", "fallback-secret-change-in-production"))

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	dbUser := getEnv("DB_USER", "reinschit_user")
	dbPassword := getEnv("DB_PASSWORD", "ReinschIt123!")
	dbName := getEnv("DB_NAME", "reinschit")
	dbHost := getEnv("DB_HOST", "localhost:3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}
	log.Println("Connected to database")

	r := gin.Default()
	r.Use(corsMiddleware())

	r.POST("/api/register", registerHandler)
	r.POST("/api/login", loginHandler)

	auth := r.Group("/api")
	auth.Use(authMiddleware())
	{
		auth.GET("/sessions", getSessionsHandler)
		auth.POST("/sessions", createSessionHandler)
		auth.GET("/sessions/:id", getSessionHandler)
		auth.PUT("/sessions/:id", updateSessionHandler)
		auth.DELETE("/sessions/:id", deleteSessionHandler)
		auth.POST("/sessions/:id/messages", createMessageHandler)
	}

	port := getEnv("PORT", "8081")
	log.Println("Server running on port", port)
	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Next()
	}
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func registerHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	result, err := db.Exec("INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)", req.Name, req.Email, string(hash))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	id, _ := result.LastInsertId()
	token := generateToken(int(id))
	c.JSON(http.StatusCreated, gin.H{"token": token, "user": gin.H{"id": id, "name": req.Name, "email": req.Email}})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int
	var name, hash string
	err := db.QueryRow("SELECT id, name, password_hash FROM users WHERE email = ?", req.Email).Scan(&id, &name, &hash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token := generateToken(id)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": id, "name": name, "email": req.Email}})
}

func generateToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	str, _ := token.SignedString(jwtSecret)
	return str
}

type Session struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	VehicleYear        *int      `json:"vehicle_year"`
	VehicleMake        *string   `json:"vehicle_make"`
	VehicleModel       *string   `json:"vehicle_model"`
	ProblemDescription string    `json:"problem_description"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func getSessionsHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := db.Query("SELECT id, user_id, vehicle_year, vehicle_make, vehicle_model, problem_description, status, created_at, updated_at FROM sessions WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sessions"})
		return
	}
	defer rows.Close()
	var sessions []Session
	for rows.Next() {
		var s Session
		rows.Scan(&s.ID, &s.UserID, &s.VehicleYear, &s.VehicleMake, &s.VehicleModel, &s.ProblemDescription, &s.Status, &s.CreatedAt, &s.UpdatedAt)
		sessions = append(sessions, s)
	}
	if sessions == nil {
		sessions = []Session{}
	}
	c.JSON(http.StatusOK, sessions)
}

type CreateSessionRequest struct {
	VehicleYear        *int    `json:"vehicle_year"`
	VehicleMake        *string `json:"vehicle_make"`
	VehicleModel       *string `json:"vehicle_model"`
	ProblemDescription string  `json:"problem_description" binding:"required"`
}

func createSessionHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := db.Exec("INSERT INTO sessions (user_id, vehicle_year, vehicle_make, vehicle_model, problem_description) VALUES (?, ?, ?, ?, ?)",
		userID, req.VehicleYear, req.VehicleMake, req.VehicleModel, req.ProblemDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}
	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id, "problem_description": req.ProblemDescription, "status": "active"})
}

func getSessionHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))
	var s Session
	err := db.QueryRow("SELECT id, user_id, vehicle_year, vehicle_make, vehicle_model, problem_description, status, created_at, updated_at FROM sessions WHERE id = ? AND user_id = ?", sessionID, userID).
		Scan(&s.ID, &s.UserID, &s.VehicleYear, &s.VehicleMake, &s.VehicleModel, &s.ProblemDescription, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	rows, _ := db.Query("SELECT id, session_id, role, content, created_at FROM messages WHERE session_id = ? ORDER BY created_at ASC", sessionID)
	defer rows.Close()
	type Message struct {
		ID        int       `json:"id"`
		SessionID int       `json:"session_id"`
		Role      string    `json:"role"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}
	var messages []Message
	for rows.Next() {
		var m Message
		rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt)
		messages = append(messages, m)
	}
	if messages == nil {
		messages = []Message{}
	}
	c.JSON(http.StatusOK, gin.H{"session": s, "messages": messages})
}

func updateSessionHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := db.Exec("UPDATE sessions SET status = ?, updated_at = NOW() WHERE id = ? AND user_id = ?", req.Status, sessionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "session updated"})
}

func deleteSessionHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))
	db.Exec("DELETE FROM messages WHERE session_id = ?", sessionID)
	_, err := db.Exec("DELETE FROM sessions WHERE id = ? AND user_id = ?", sessionID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "session deleted"})
}

func createMessageHandler(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE id = ? AND user_id = ?", sessionID, userID).Scan(&count)
	if err != nil || count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	db.Exec("INSERT INTO messages (session_id, role, content) VALUES (?, 'user', ?)", sessionID, req.Content)
	aiResponse := "Thank you for the details. Based on what you described, here are some steps to diagnose the issue: First, check the most common causes for this symptom. Would you like me to walk you through the diagnostic process step by step?"
	db.Exec("INSERT INTO messages (session_id, role, content) VALUES (?, 'assistant', ?)", sessionID, aiResponse)
	c.JSON(http.StatusCreated, gin.H{
		"user_message": gin.H{"role": "user", "content": req.Content},
		"ai_message":   gin.H{"role": "assistant", "content": aiResponse},
	})
}