package server

import (
	"errors"
	"net/http"

	models "github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/domain/model"
	"github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/storage"

	"github.com/gin-gonic/gin"
)

type Storage interface {
	SaveUser(models.User) error
	ValidateUser(models.User) (string, error)
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
	if err := s.storage.SaveUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "User was saved")
}

func (s *Server) AuthHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := s.storage.ValidateUser(user)
	if err != nil {
		if errors.Is(err, storage.ErrIvalidAuthData) {
			ctx.String(http.StatusUnauthorized, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "Auth compltted")
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
