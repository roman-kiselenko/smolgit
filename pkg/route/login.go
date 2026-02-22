package route

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"smolgit/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *Route) Login(c *gin.Context) {
	var req creds
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("parsing", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	index := slices.IndexFunc(r.users, func(u model.User) bool {
		return u.Name == req.Username
	})
	if index == -1 {
		slog.Error("no such user", "user", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}
	u := r.users[index]
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		slog.Error("bad password", "user", req.Username, "password", req.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	exp := time.Now().Add(1 * time.Hour)
	claims := &model.Claims{
		Username: u.Name,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(r.cfg.ServerJWTKey))
	if err != nil {
		slog.Error("cant sign token", "user", req.Username, "password", req.Password)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": t})
}
