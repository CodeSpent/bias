package handlers

import (
	"bias/models"
	"bias/store"
	"bias/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type StreamHandler struct {
	*BaseHandler
	Store *store.StreamStore
}

func NewStreamHandler(store store.StreamStore) *StreamHandler {
	return &StreamHandler{Store: &store}
}

func (h *StreamHandler) CreateStream(c echo.Context) error {
	var stream models.StreamModel

	if err := c.Bind(&stream); err != nil {
		return utils.RespondWithError(c, "Invalid request body", http.StatusBadRequest)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(stream); err != nil {
		return utils.RespondWithError(c, "Validation error: "+err.Error(), http.StatusBadRequest)
	}

	if err := h.Store.CreateStream(&stream); err != nil {
		return utils.RespondWithError(c, "Failed to create stream", http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, stream)
}

func (h *StreamHandler) ListStreams(c echo.Context) error {
	streams, err := h.Store.GetAllStreams()

	if err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, streams)
}

func (h *StreamHandler) GetStream(c echo.Context) error {
	streamID := uint(1)
	if err, _ := h.Store.GetStreamByID(streamID); err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}
	return nil
}

func (h *StreamHandler) UpdateStream(c echo.Context) error {
	streamID := uint(1)

	var stream models.StreamModel
	if err, _ := h.Store.GetStreamByID(streamID); err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}

	if err := h.Store.UpdateStream(&stream); err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}
	return nil
}

func (h *StreamHandler) DeleteStream(c echo.Context) error {
	streamID := uint(1)

	var stream models.StreamModel
	if _, err := h.Store.GetStreamByID(streamID); err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}

	if err := h.Store.DeleteStream(&stream); err != nil {
		return utils.RespondWithError(c, "Stream not found", http.StatusNotFound)
	}
	return nil
}
