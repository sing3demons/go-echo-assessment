package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type NewsExpenses struct {
	ID     int            `json:"id"`
	Title  string         `json:"title"`
	Amount float64        `json:"amount"`
	Note   string         `json:"note"`
	Tags   pq.StringArray `json:"tags"`
}

func (h *handler) CreateExpensesHandler(c echo.Context) error {
	var m NewsExpenses

	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := h.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1,$2,$3,$4) RETURNING id", m.Title, m.Amount, m.Note, m.Tags)
	err := row.Scan(&m.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusCreated, m.ID)
}

func (h *handler) UpdateExpensesHandler(c echo.Context) error {
	var m NewsExpenses

	id := c.Param("id")

	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	stmt, err := h.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	_, err = stmt.Exec(id, m.Title, m.Amount, m.Note, m.Tags)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, "Update")
}
