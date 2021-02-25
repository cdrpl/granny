package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cdrpl/granny/server/pkg/validate"
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
	nameExists, err := UserNameExists(c.PgPool, form.Name)
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{"message": "An error has occured"})
		return
	} else if nameExists {
		ctx.JSON(200, gin.H{"message": "Name is taken"})
		return
	}

	// User email must be unique
	emailExists, err := UserEmailExists(c.PgPool, form.Email)
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
	user := User{Name: form.Name, Email: form.Email, Pass: string(hash)}

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
	user := User{}
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
	token, err := generateToken(16)
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
