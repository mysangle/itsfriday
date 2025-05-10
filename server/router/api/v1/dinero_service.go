package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"itsfriday/internal/util"
	"itsfriday/store"
)

// money service

type DineroServiceServer interface {
	CreateDineroCaterory(echo.Context) error
	UpdateDineroCaterory(echo.Context) error
	ListDineroCaterories(echo.Context) error
	DeleteDineroCaterory(echo.Context) error

	CreateDineroExpense(echo.Context) error
	UpdateDineroExpense(echo.Context) error
	ListDineroExpenses(echo.Context) error
	DeleteDineroExpense(echo.Context) error
	ReportDinero(echo.Context) error
}

type CreateDineroCategoryRequest struct {
	Name     string     `json:"name"`
	Priority int32      `json:"priority"`
}

type UpdateDineroCategoryRequest struct {
	Name     string     `json:"name"`
	Priority int32      `json:"priority"`
}

type DeleteDineroCategoryRequest struct {
	Name     string     `json:"name"`
}

type DineroCategory struct {
	ID       int32      `json:"id"`
	Name     string     `json:"name"`
	Priority int32      `json:"priority"`
}

type DineroCategories struct {
	Categories []*DineroCategory    `json:"categories"`
}

type CreateDineroExpenseRequest struct {
	CategoryID int32     `json:"categoryId"`
	DateUsed   string    `json:"dateUsed"`
	Item       string    `json:"item"`
	Price      int32     `json:"price"`
}

type UpdateDineroExpenseRequest struct {
	CategoryID int32     `json:"id"`
	DateUsed   string    `json:"dateUsed"`
	Item       string    `json:"item"`
	Price      int32     `json:"price"`
}

type DineroExpense struct {
	ID         int32     `json:"id"`

	CategoryID int32     `json:"categoryId"`
	DateUsed   string    `json:"dateUsed"`
	Item       string    `json:"item"`
	Price      int32     `json:"price"`
}

type DineroExpenses struct {
	Expenses []*DineroExpense    `json:"expenses"`
}

func (s *APIV1Service) CreateDineroCaterory(c echo.Context) error {
	ctx := c.Request().Context()
    request := new(CreateDineroCategoryRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid creating category request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	var priority int32
	if request.Priority == 0 {
		priority = 1
	} else {
		priority = request.Priority
	}
    create := &store.DineroCategory{
		UserID:      userID,
		Name:        request.Name,
		Priority:    priority,
	}
	category, err := s.Store.CreateDineroCaterory(ctx, create)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to create dinero category: %v", err),
		})
	}

	categoryInfo := convertCategoryFromStore(category)
	return c.JSON(http.StatusOK, categoryInfo)
}

func (s *APIV1Service) UpdateDineroCaterory(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetCategory: ", "id", id)
	categoryId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get category_id from url",
		})
	}
	request := new(UpdateDineroCategoryRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update dinero category request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	category, err := s.Store.GetDineroCategory(ctx, &store.FindDineroCategory{
		ID:     &categoryId,
		UserID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get dinero category: %v", err),
		})
	}
	if category == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "category not found",
		})
	}
	if category.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	update := &store.UpdateDineroCategory{
		ID:     categoryId,
	}
	if request.Name != "" {
		update.Name = &request.Name
	}
	if request.Priority != 0 {
		update.Priority = &request.Priority
	}

	updatedCategory, err := s.Store.UpdateDineroCategory(ctx, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to update dinero category: %v", err),
		})
	}

	categoryInfo := convertCategoryFromStore(updatedCategory)
	return c.JSON(http.StatusOK, categoryInfo)
}

func (s *APIV1Service) ListDineroCaterories(c echo.Context) error {
	ctx := c.Request().Context()
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	categories, err := s.Store.ListDineroCategories(ctx, &store.FindDineroCategory{
		UserID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get dinero categories: %v", err),
		})
	}

	list := make([]*DineroCategory, 0)
	for _, category := range categories {
		categoryInfo := convertCategoryFromStore(category)
		list = append(list, categoryInfo)
	}
	return c.JSON(http.StatusOK, &DineroCategories{Categories: list})
}

func (s *APIV1Service) DeleteDineroCaterory(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetCategory: ", "id", id)
	categoryId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get category_id from url",
		})
	}
	request := new(DeleteDineroCategoryRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid deleting category request: %v", err),
		})
	}
    userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}
	
	category, err := s.Store.GetDineroCategory(ctx, &store.FindDineroCategory{
		ID:     &categoryId,
		UserID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get category: %v", err),
		})
	}
	if category == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "category not found",
		})
	}
	if category.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	if err := s.Store.DeleteDineroCategory(ctx, &store.DeleteDineroCategory{
		ID:   userID,
	}); err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to delete category: %v", err),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) CreateDineroExpense(c echo.Context) error {
    ctx := c.Request().Context()
    request := new(CreateDineroExpenseRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid creating expense request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

    create := &store.DineroExpense{
		UserID:      userID,
		CategoryID:  request.CategoryID,
		DateUsed:    request.DateUsed,
		Item:        request.Item,
		Price:       request.Price,
	}
	category, err := s.Store.CreateDineroExpense(ctx, create)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorResponse{
			Code:    Internal,
		    Message: fmt.Sprintf("failed to create dinero category: %v", err),
		})
	}

	categoryInfo := convertExpenseFromStore(category)
	return c.JSON(http.StatusOK, categoryInfo)
}

func (s *APIV1Service) UpdateDineroExpense(c echo.Context) error {
    ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetExpense: ", "id", id)
	expenseId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get expense_id from url",
		})
	}
	request := new(UpdateDineroExpenseRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
		    Message: fmt.Sprintf("invalid update dinero expense request: %v", err),
		})
	}
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	expense, err := s.Store.GetDineroExpense(ctx, &store.FindDineroExpense{
		ID:     &expenseId,
		UserID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get dinero expense: %v", err),
		})
	}
	if expense == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "expense not found",
		})
	}
	if expense.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	update := &store.UpdateDineroExpense{
		ID:     expenseId,
	}
	if request.CategoryID != 0 {
		update.CategoryID = &request.CategoryID
	}
	if request.DateUsed != "" {
		update.DateUsed = &request.DateUsed
	}
	if request.Item != "" {
		update.Item = &request.Item
	}
	if request.Price != 0 {
		update.Price = &request.Price
	}
	updatedExpense, err := s.Store.UpdateDineroExpense(ctx, update)
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to update dinero expense: %v", err),
		})
	}

	categoryInfo := convertExpenseFromStore(updatedExpense)
	return c.JSON(http.StatusOK, categoryInfo)
}

func (s *APIV1Service) DeleteDineroExpense(c echo.Context) error {
    ctx := c.Request().Context()
	id := c.Param("id")
	slog.Debug("GetExpense: ", "id", id)
	expenseId, err := util.ConvertStringToInt32(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get expense_id from url",
		})
	}
    userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}
	
	expense, err := s.Store.GetDineroExpense(ctx, &store.FindDineroExpense{
		ID:     &expenseId,
		UserID: &userID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get expense: %v", err),
		})
	}
	if expense == nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    NotFound,
			Message: "category not found",
		})
	}
	if expense.UserID != userID {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    PermissionDenied,
			Message: "permission denied",
		})
	}

	if err := s.Store.DeleteDineroExpense(ctx, &store.DeleteDineroExpense{
		ID:   expense.ID,
	}); err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to delete expense: %v", err),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *APIV1Service) ListDineroExpenses(c echo.Context) error {
    ctx := c.Request().Context()
	year, err := util.GetYearFromQueryParam(c.QueryParam("year"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: fmt.Sprintf("invalid query param: %v", err),
		})
	}
	month, err := util.GetMonthFromQueryParam(c.QueryParam("month"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: fmt.Sprintf("invalid query param: %v", err),
		})
	}
	slog.Debug("ListDineroExpenses: ", "year", year, "month", month)
	userID, ok := c.Get(useridContextKey).(int32)
	if !ok {
	    return c.JSON(http.StatusBadRequest, &ErrorResponse{
			Code:    InvalidRequest,
			Message: "failed to get userid from access token",
		})
	}

	expenses, err := s.Store.ListDineroExpenses(ctx, &store.FindDineroExpense{
		UserID: &userID,
		Year:   &year,
		Month:  &month,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, &ErrorResponse{
			Code:    Internal,
			Message: fmt.Sprintf("failed to get dinero expenses: %v", err),
		})
	}

	list := make([]*DineroExpense, 0)
	for _, expense := range expenses {
		expenseInfo := convertExpenseFromStore(expense)
		list = append(list, expenseInfo)
	}
	return c.JSON(http.StatusOK, &DineroExpenses{Expenses: list})
}

func (s *APIV1Service) ReportDinero(c echo.Context) error {
    return c.NoContent(http.StatusNoContent)
}

func convertCategoryFromStore(category *store.DineroCategory) *DineroCategory {
	categoryInfo := &DineroCategory{
		ID:          category.ID,
		Name:        category.Name,
		Priority:    category.Priority,
	}
	return categoryInfo
}

func convertExpenseFromStore(category *store.DineroExpense) *DineroExpense {
	expenseInfo := &DineroExpense{
		ID:          category.ID,
		CategoryID:  category.CategoryID,
		DateUsed:    category.DateUsed,
		Item:        category.Item,
		Price:       category.Price,
	}
	return expenseInfo
}
