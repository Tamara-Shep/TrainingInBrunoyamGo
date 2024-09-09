package storage

import (
	models "github.com/Tamara-Shep/TrainingInBrunoyamGo/internal/domain/model"
	uuid "github.com/google/uuid"
)

type MemStorage struct {
	usersMap map[string]models.User
	booksMap map[string]models.Book
}

func New() *MemStorage {
	uMap := make(map[string]models.User)
	bMap := make(map[string]models.Book)
	return &MemStorage{
		usersMap: uMap,
		booksMap: bMap,
	}
}

func (ms *MemStorage) SaveUser(user models.User) error {
	uid := uuid.New().String()
	ms.usersMap[uid] = user
	return nil
}

func (ms *MemStorage) ValidateUser(user models.User) (string, error) {
	for uid, value := range ms.usersMap {
		if value.Email == user.Email {
			if value.Password != user.Password {
				return "", ErrIvalidAuthData
			}
			return uid, nil
		}
	}
	return "", ErrUserNotFound
}

func (ms *MemStorage) GetBooks() ([]models.Book, error) {
	var books []models.Book
	for bid, value := range ms.booksMap {
		book := value
		book.BID = bid
		books = append(books, book)
	}
	if len(books) == 0 {
		return nil, ErrBookListEmrty
	}
	return books, nil
}

func (ms *MemStorage) GetBookId(bId string) (models.Book, error) {
	book, ok := ms.booksMap[bId]
	if !ok {
		return models.Book{}, ErrBookNotFound
	}
	return book, nil
}

func (ms *MemStorage) SaveBook(book models.Book) error {
	bId := uuid.New().String()
	ms.booksMap[bId] = book
	return nil
}

func (ms *MemStorage) DeleteBook(bId string) error {
	_, ok := ms.booksMap[bId]
	if !ok {
		return ErrBookNotFound
	}
	delete(ms.booksMap, bId)
	return nil
}
