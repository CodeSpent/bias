package handlers

import (
	"bias/models"
	"bias/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// TagHandler handles the HTTP requests related to tags.
type TagHandler struct {
	*BaseHandler
	tagStore *store.TagStore
}

// NewTagHandler creates a new TagHandler instance.
func NewTagHandler(tagStore *store.TagStore) *TagHandler {
	return &TagHandler{
		BaseHandler: &BaseHandler{},
		tagStore:    tagStore,
	}
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag with provided data
// @Accept json
// @Produce json
// @Param tag body models.Tag true "Tag object"
// @Success 201 {object} models.Tag
// @Router /tags [post]
func (h *TagHandler) CreateTag(c echo.Context) error {
	var tag models.TagModel
	if err := c.Bind(&tag); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err)
	}

	err := h.tagStore.CreateTag(&tag)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return h.SuccessResponse(c, http.StatusCreated, tag)
}

// GetTagByID godoc
// @Summary Get tag by ID
// @Description Get tag details by providing its ID
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} models.Tag
// @Router /tags/{id} [get]
func (h *TagHandler) GetTagByID(c echo.Context) error {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err)
	}

	tag, err := h.tagStore.GetTagByID(uint(tagID))
	if err != nil {
		return h.ErrorResponse(c, http.StatusNotFound, err)
	}

	return h.SuccessResponse(c, http.StatusOK, tag)
}

// ListTags godoc
// @Summary Get all tags
// @Description Get details of all available tags
// @Produce json
// @Success 200 {array} models.Tag
// @Router /tags [get]
func (h *TagHandler) ListTags(c echo.Context) error {
	tags, err := h.tagStore.GetAllTags()

	if err != nil {

	}

	return h.SuccessResponse(c, http.StatusOK, tags)
}

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing tag with new data
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body models.Tag true "Tag object"
// @Success 200 {object} models.Tag
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c echo.Context) error {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err)
	}

	var updatedTag models.TagModel
	if err := c.Bind(&updatedTag); err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err)
	}
	updatedTag.ID = uint(tagID)

	err = h.tagStore.UpdateTag(&updatedTag)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return h.SuccessResponse(c, http.StatusOK, updatedTag)
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a tag by ID
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {string} string "Tag deleted"
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c echo.Context) error {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, err)
	}

	err = h.tagStore.DeleteTag(uint(tagID))
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, err)
	}

	return h.SuccessResponse(c, http.StatusOK, fmt.Sprintf("Tag with ID %d deleted", tagID))
}
