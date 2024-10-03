package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	models "github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/domain/model"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "VerySecretKey2000"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type Storage interface {
	SaveUser(models.User) (string, error)
	ValidateUser(models.User) (string, string, error)
	GetBooks() ([]models.Book, error)
	GetBookId(string) (models.Book, error)
	SaveBook(models.Book) error
	DeleteBook(string) error
}

type Server struct {
	host    string
	storage Storage
}

func New(host string, storage Storage) *Server {
	return &Server{
		host:    host,
		storage: storage,
	}

}

func (s *Server) Run() error {
	r := gin.Default()
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", s.RegisterHandler)
		userGroup.POST("/auth", s.AuthHandler)
	}

	bookGroup := r.Group("/books")
	{
		bookGroup.GET("/all-books", s.AllBookHandler)
		bookGroup.GET("/:id", s.GetBookHandler)
		bookGroup.POST("/add-book", s.SaveBookHandler)
		bookGroup.DELETE("/delete/:id", s.DeleteBookHandler)
	}
	if err := r.Run(s.host); err != nil {
		return err
	}
	return nil
}

func (s *Server) RegisterHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	user.Password = string(passHash)

	UID, err := s.storage.SaveUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := createJWT(UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "User was saved")
}

func (s *Server) AuthHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	UID, password, err := s.storage.ValidateUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrIvalidAuthData) {
			ctx.String(http.StatusUnauthorized, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Неверная пара логин пароль"})
		return
	}

	token, err := createJWT(UID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "Auth completted")
}

func (s *Server) AllBookHandler(ctx *gin.Context) {
	books, err := s.storage.GetBooks()
	if err != nil {
		if errors.Is(err, storage.ErrBookListEmrty) {
			ctx.String(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)

}

func (s *Server) GetBookHandler(ctx *gin.Context) {
	bid := ctx.Param("id")
	book, err := s.storage.GetBookId(bid)
	if err != nil {
		if errors.Is(err, storage.ErrBookNotFound) {
			ctx.String(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (s *Server) SaveBookHandler(ctx *gin.Context) {
	var book models.Book
	if err := ctx.ShouldBindBodyWithJSON(&book); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.storage.SaveBook(book); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusCreated, "book was saved")
}

func (s *Server) DeleteBookHandler(ctx *gin.Context) {
	bid := ctx.Param("id")
	if err := s.storage.DeleteBook(bid); err != nil {
		if errors.Is(err, storage.ErrBookNotFound) {
			ctx.String(http.StatusNoContent, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "book was delete")
}

func createJWT(UID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 3)),
		},
		UserID: UID,
	})
	key := []byte(SecretKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenStr, err
}

func getUID(tokenStr string) (string, error) {
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	return claims.UserID, nil
}
