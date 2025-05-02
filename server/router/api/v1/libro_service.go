package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"itsfriday/internal/util"
	"itsfriday/store"
)

// libro service

type LibroServiceServer interface {
	CreateBook(echo.Context) error
	GetBook(echo.Context) error
	UpdateBook(echo.Context) error
	DeleteBook(echo.Context) error
	
	CreateBookReview(echo.Context) error
	GetBookReview(echo.Context) error
	UpdateBookReview(echo.Context) error
	DeleteBookReview(echo.Context) error

	Dashboard(echo.Context) error
	ReadBook(echo.Context) error
	ReportBook(echo.Context) error
	BookReviews(echo.Context) error
}

type CreateBookRequest struct {
	Title        string            `json:"title"`
	Author       string            `json:"author"`
	Translator   string            `json:"translator"`
	Pages        int32             `json:"pages"`
	PubYear      int32             `json:"pub_year"`
	Genre        string            `json:"genre"`
}

type UpdateBookRequest struct {
	Title        string            `json:"title"`
	Author       string            `json:"author"`
	Translator   string            `json:"translator"`
	Pages        int32             `json:"pages"`
	PubYear      int32             `json:"pub_year"`
	Genre        string            `json:"genre"`
}

type Book struct {
	ID           int32             `json:"id"`
	CreatedTime  int64             `json:"createdTime"`
	Title        string            `json:"title"`
	Author       string            `json:"author"`
	Translator   string            `json:"translator"`
	Pages        int32             `json:"pages"`
	PubYear      int32             `json:"pub_year"`
	Genre        string            `json:"genre"`
}

func (s *APIV1Service) CreateBook(c echo.Context) error {
	ctx := c.Request().Context()
    request := new(CreateBookRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid creating book request: %v", err),
		})
	}

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	create := &store.Book{
		UserID:      userID,
		Title:       request.Title,
		Author:      request.Author,
		Translator:  request.Translator,
		Pages:       request.Pages,
		PubYear:     request.PubYear,
		Genre:       request.Genre,
	}
	book, err := s.Store.CreateBook(ctx, create)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to create book: %v", err),
		})
	}

	bookInfo := convertBookFromStore(book)
	return c.JSON(http.StatusOK, bookInfo)
}

func (s *APIV1Service) GetBook(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBook: ", "id", id)
	bookId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get user_id from url",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &bookId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book: %v", err),
		})
	}
	if book == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book not found",
		})
	}

	bookInfo := convertBookFromStore(book)
	return c.JSON(http.StatusOK, bookInfo)
}

func (s *APIV1Service) UpdateBook(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBook: ", "id", id)
	bookId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get user_id from url",
		})
	}
	request := new(UpdateBookRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update book request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &bookId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book: %v", err),
		})
	}
	if book == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book not found",
		})
	}
	if book.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	update := &store.UpdateBook{
		ID: book.ID,
	}
	if request.Title != "" {
		update.Title = &request.Title
	}
	if request.Author != "" {
		update.Author = &request.Author
	}
	if request.Translator != "" {
		update.Translator = &request.Translator
	}
	if request.Pages != 0 {
		update.Pages = &request.Pages
	}
	if request.PubYear != 0 {
		update.PubYear = &request.PubYear
	}
	if request.Genre != "" {
		update.Genre = &request.Genre
	}

	updatedBook, err := s.Store.UpdateBook(ctx, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to update book: %v", err),
		})
	}

	bookInfo := convertBookFromStore(updatedBook)
	return c.JSON(http.StatusOK, bookInfo)
}

func (s *APIV1Service) DeleteBook(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBook: ", "id", id)
	bookId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get user_id from url",
		})
	}
	request := new(UpdateBookRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update book request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &bookId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book: %v", err),
		})
	}
	if book == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book not found",
		})
	}
	if book.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	if err := s.Store.DeleteBook(ctx, &store.DeleteBook{
		ID: bookId,
	}); err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to delete book: %v", err),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) CreateBookReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) GetBookReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) UpdateBookReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) DeleteBookReview(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) Dashboard(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) ReadBook(c echo.Context) error {
	year := c.QueryParam("year")
	if year == "" {
		year = "2025" // this year
	}
	slog.Debug("ReadBook: ", "year", year)

	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

// year - count
func (s *APIV1Service) ReportBook(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) BookReviews(c echo.Context) error {
	id := c.Param("id")
	slog.Debug("BookReviews: ", "id", id)
    return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func convertBookFromStore(book *store.Book) *Book {
	bookInfo := &Book{
        ID:          book.ID,
		CreatedTime: book.CreatedTs,
		Title:       book.Title,
		Author:      book.Author,
		Translator:  book.Translator,
		Pages:       book.Pages,
		PubYear:     book.PubYear,
		Genre:       book.Genre,
	}
	return bookInfo
}
