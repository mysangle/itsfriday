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
	Title        string              `json:"title"`
	Author       string              `json:"author"`
	Translator   string              `json:"translator"`
	Pages        int32               `json:"pages"`
	PubYear      int32               `json:"pubYear"`
	Genre        string              `json:"genre"`
}

type UpdateBookRequest struct {
	Title        string              `json:"title"`
	Author       string              `json:"author"`
	Translator   string              `json:"translator"`
	Pages        int32               `json:"pages"`
	PubYear      int32               `json:"pubYear"`
	Genre        string              `json:"genre"`
}

type Book struct {
	ID           int32               `json:"id"`
	CreatedTime  int64               `json:"createdTime"`
	Title        string              `json:"title"`
	Author       string              `json:"author"`
	Translator   string              `json:"translator"`
	Pages        int32               `json:"pages"`
	PubYear      int32               `json:"pubYear"`
	Genre        string              `json:"genre"`
	Updatable    bool                `json:"updatable"`
}

type CreateBookReviewRequest struct {
	BookID       int32               `json:"bookId"`
	DateRead     string              `json:"dateRead"`
	Rating       float32             `json:"rating"`
	Review       string              `json:"review"`
}

type UpdateBookReviewRequest struct {
	BookID       int32               `json:"bookId"`
	DateRead     string              `json:"dateRead"`
	Rating       float32             `json:"rating"`
	Review       string              `json:"review"`
}

type BookReview struct {
	ID           int32               `json:"id"`
	CreatedTime  int64               `json:"createdTime"`
	BookID       int32               `json:"bookId"`
	DateRead     string              `json:"dateRead"`
	Rating       float32             `json:"rating"`
	Review       string              `json:"review"`
	Title        string              `json:"title"`
}

type BooksRead struct {
	Books        []*store.BookRead   `json:"books"`
}

type ReportBook struct {
	Report       []*store.ReportBook `json:"report"`
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

	bookInfo := convertBookFromStore(book, true)
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

	bookInfo := convertBookFromStore(book, userID == book.UserID)
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
			Message: "failed to get book_id from url",
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

	bookInfo := convertBookFromStore(updatedBook, true)
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
	ctx := c.Request().Context()
    request := new(CreateBookReviewRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid creating book review request: %v", err),
		})
	}
	ok := util.ValidateDate(request.DateRead)
	if !ok {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid dateRead field: %s", request.DateRead),
		})
	}

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &request.BookID})
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

	create := &store.BookReview{
		UserID:      userID,
		BookID:      request.BookID,
		DateRead:    request.DateRead,
		Rating:      request.Rating,
		Review:      request.Review,
	}
	bookReview, err := s.Store.CreateBookReview(ctx, create)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to create book review: %v", err),
		})
	}

	bookReviewInfo := convertBookReviewFromStore(bookReview, book)
	return c.JSON(http.StatusOK, bookReviewInfo)
}

func (s *APIV1Service) GetBookReview(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBookReview: ", "id", id)
	bookReviewId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get book review id from url",
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	bookReview, err := s.Store.GetBookReview(ctx, &store.FindBookReview{ID: &bookReviewId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book review: %v", err),
		})
	}
	if bookReview == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book review not found",
		})
	}
	if bookReview.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "premission denied",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &bookReview.BookID})
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

	bookReviewInfo := convertBookReviewFromStore(bookReview, book)
	return c.JSON(http.StatusOK, bookReviewInfo)
}

func (s *APIV1Service) UpdateBookReview(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBookReview: ", "id", id)
	bookReviewId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get book review id from url",
		})
	}
	request := new(UpdateBookReviewRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update book review request: %v", err),
		})
	}
	ok := util.ValidateDate(request.DateRead)
	if !ok {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid dateRead field: %s", request.DateRead),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	bookReview, err := s.Store.GetBookReview(ctx, &store.FindBookReview{ID: &bookReviewId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book review: %v", err),
		})
	}
	if bookReview == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book review not found",
		})
	}
	if bookReview.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "premission denied",
		})
	}

	book, err := s.Store.GetBook(ctx, &store.FindBook{ID: &bookReview.BookID})
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

	update := &store.UpdateBookReview{
		ID: bookReview.ID,
	}
	if request.BookID != 0 {
		update.BookID = &request.BookID
	}
	if request.DateRead != "" {
		update.DateRead = &request.DateRead
	}
	if request.Rating != 0 {
		update.Rating = &request.Rating
	}
	if request.Review != "" {
		update.Review = &request.Review
	}

	updatedBookReview, err := s.Store.UpdateBookReview(ctx, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to update book review: %v", err),
		})
	}

	bookReviewInfo := convertBookReviewFromStore(updatedBookReview, book)
	return c.JSON(http.StatusOK, bookReviewInfo)
}

func (s *APIV1Service) DeleteBookReview(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetBookReview: ", "id", id)
	bookReviewId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get book review id from url",
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	bookReview, err := s.Store.GetBookReview(ctx, &store.FindBookReview{ID: &bookReviewId})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book review: %v", err),
		})
	}
	if bookReview == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "book review not found",
		})
	}
	if bookReview.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "premission denied",
		})
	}

	if err := s.Store.DeleteBookReview(ctx, &store.DeleteBookReview{
		ID: bookReviewId,
	}); err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to delete book review: %v", err),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) Dashboard(c echo.Context) error {
	return c.JSON(http.StatusNotFound, &ErrorResponse{
		Code:    Unimplemented,
		Message: "not yet implemented",
	})
}

func (s *APIV1Service) ReadBook(c echo.Context) error {
	ctx := c.Request().Context()
	year, err := util.GetYearFromQueryParam(c.QueryParam("year"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: fmt.Sprintf("invalid query param: %v", err),
		})
	}
	slog.Debug("ReadBook: ", "year", year)

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	list, err := s.Store.ListBooksReadInYear(ctx, userID, year)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get book review: %v", err),
		})
	}

	return c.JSON(http.StatusOK, &BooksRead{
		Books: list,
	})
}

// year - count
func (s *APIV1Service) ReportBook(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	list, err := s.Store.ReportBook(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get books read by year: %v", err),
		})
	}

	return c.JSON(http.StatusOK, &ReportBook{
		Report: list,
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

func convertBookFromStore(book *store.Book, updatable bool) *Book {
	bookInfo := &Book{
        ID:          book.ID,
		CreatedTime: book.CreatedTs,
		Title:       book.Title,
		Author:      book.Author,
		Translator:  book.Translator,
		Pages:       book.Pages,
		PubYear:     book.PubYear,
		Genre:       book.Genre,
		Updatable:   updatable,
	}
	return bookInfo
}

func convertBookReviewFromStore(bookReview *store.BookReview, book *store.Book) *BookReview {
	bookReviewInfo := &BookReview{
        ID:          bookReview.ID,
		CreatedTime: bookReview.CreatedTs,
		BookID:      bookReview.BookID,
		DateRead:    bookReview.DateRead,
		Rating:      bookReview.Rating,
		Review:      bookReview.Review,
		Title:       book.Title,
	}
	return bookReviewInfo
}
