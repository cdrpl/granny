package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cdrpl/idlemon/pkg/db"
	"github.com/cdrpl/idlemon/pkg/rnd"
	"github.com/cdrpl/idlemon/pkg/validate"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL    = 60 * 30 // Auth token TTL in seconds
	passMinLen  = 8       // Minimum required characters in a user password
	userNameMax = 16      // Maximum characters allowed in a user name
)

// Controller holds the HTTP route logic.
type Controller struct {
	PgPool *pgxpool.Pool
	Rdb    *redis.Client
}

// Health check route
func (c *Controller) health(ctx *gin.Context) {
	ctx.String(200, "Ok")
}

// User sign up route.
func (c *Controller) signUpHandler(ctx *gin.Context) {
	// Create form
	form := &SignUpForm{
		Name:  ctx.PostForm("name"),
		Email: ctx.PostForm("email"),
		Pass:  ctx.PostForm("pass"),
	}

	form.sanitize()

	// Check form for errors
	if hasErr, errMsg := form.hasError(); hasErr {
		ctx.JSON(400, gin.H{"message": errMsg})
		return
	}

	// User name must be unique
	nameExists, err := db.UserNameExists(c.PgPool, form.Name)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	} else if nameExists {
		ctx.JSON(200, gin.H{"message": "Name is taken"})
		return
	}

	// User email must be unique
	emailExists, err := db.UserEmailExists(c.PgPool, form.Email)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"message": err})
		return
	} else if emailExists {
		ctx.JSON(200, gin.H{"message": "Email is taken"})
		return
	}

	// Hash user password
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Pass), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash user password:", err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Create new user
	user := db.User{Name: form.Name, Email: form.Email, Pass: string(hash)}

	// Insert user
	err = user.Insert(c.PgPool)
	if err != nil {
		log.Println("Failed to insert user:", err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Write response
	ctx.JSON(200, gin.H{"message": "Ok"})
}

// User sign in route.
func (c *Controller) signInHandler(ctx *gin.Context) {
	// Create the form
	form := &SignInForm{
		Email: ctx.PostForm("email"),
		Pass:  ctx.PostForm("pass"),
	}

	form.sanitize()

	// Check form for errors
	if hasErr, errMsg := form.hasError(); hasErr {
		ctx.JSON(400, gin.H{"message": errMsg})
		return
	}

	// Find user
	user := db.User{}
	err := c.PgPool.QueryRow(context.Background(), "SELECT id, name, pass FROM users WHERE email = $1", form.Email).Scan(&user.ID, &user.Name, &user.Pass)
	if err == pgx.ErrNoRows {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	} else if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Compare the form password to the password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(form.Pass))
	if err != nil {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Create a random auth token
	token, err := rnd.GenerateToken(16)
	if err != nil {
		log.Println("Failed to generate a random auth token:", err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Store token in redis
	key := fmt.Sprintf("%d", user.ID)
	c.Rdb.SetEX(context.Background(), key, token, time.Second*time.Duration(tokenTTL))

	// Write response
	ctx.JSON(200, gin.H{
		"message": "Ok",
		"data": gin.H{
			"token": token,
			"id":    user.ID,
		},
	})
}

// Generates a remember token if the given authentication token is valid.
func (c *Controller) createRememberToken(ctx *gin.Context) {
	id := ctx.PostForm("id")
	authToken := ctx.PostForm("token")

	if id == "" || authToken == "" {
		ctx.JSON(400, gin.H{"message": "id and token are required"})
		return
	}

	// Verify the authentication token
	isAuthorized, err := db.CheckAuth(c.Rdb, id, authToken)
	if err != nil {
		log.Println("Check auth failed:", err)
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	} else if !isAuthorized {
		ctx.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Generate token
	seriesID, rememberToken, err := rnd.GenerateRememberToken()
	if err != nil {
		log.Println("Failed to generate a random auth token:", err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Hash the remember token
	hash := sha256.Sum256([]byte(rememberToken))
	hashStr := hex.EncodeToString(hash[:])

	// Insert token row
	sql := "INSERT INTO remember_tokens (id, user_id, token, created_at) VALUES ($1, $2, $3, $4)"
	_, err = c.PgPool.Exec(context.Background(), sql, seriesID, id, hashStr, time.Now())
	if err != nil {
		log.Println("Failed to insert the remember_tokens row:", err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	}

	// Respond with the series ID and the unhashed token
	ctx.JSON(200, gin.H{
		"message": "Ok",
		"data": gin.H{
			"id":    seriesID,
			"token": rememberToken,
		},
	})
}

// SignUpForm holds the form values for the sign up route.
type SignUpForm struct {
	Name  string
	Email string
	Pass  string
}

// Sanitize the form values.
func (s *SignUpForm) sanitize() {
	s.Name = strings.TrimSpace(s.Name)
	s.Email = strings.TrimSpace(s.Email)
	s.Email = strings.ToLower(s.Email)
}

// Returns true and the error message if validation fails, returns false and "" if validation passes
func (s *SignUpForm) hasError() (bool, string) {
	v := validate.CreateValidator()

	// Name
	v.Required(s.Name, "A user name is required")
	v.Max(s.Name, fmt.Sprintf("User names can't have more than %d characters", userNameMax), userNameMax)

	// Email
	v.Required(s.Email, "An email is required")
	v.IsEmail(s.Email, "The email is invalid")

	// Pass
	v.Required(s.Pass, "A password is required")
	v.Min(s.Pass, fmt.Sprintf("Password must contain at least %d characters", passMinLen), passMinLen)

	return v.FirstError()
}

// SignInForm holds the form values for the sign in route.
type SignInForm struct {
	Email string
	Pass  string
}

// Sanitize the form values.
func (s *SignInForm) sanitize() {
	s.Email = strings.TrimSpace(s.Email)
	s.Email = strings.ToLower(s.Email)
}

// Returns true and the error message if validation fails, returns false and "" if validation passes
func (s *SignInForm) hasError() (bool, string) {
	v := validate.CreateValidator()

	// Email
	v.Required(s.Email, "An email is required")
	v.IsEmail(s.Email, "The email is invalid")

	// Pass
	v.Required(s.Pass, "A password is required")
	v.Min(s.Pass, fmt.Sprintf("Password must contain at least %d characters", passMinLen), passMinLen)

	return v.FirstError()
}
