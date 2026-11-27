package handler

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"travel-itinerary-d1/backend/internal/model"
)

type ItineraryHandler struct {
	db *gorm.DB
}

type createItineraryRequest struct {
	UserID      uint   `json:"user_id"`
	Title       string `json:"title"`
	Destination string `json:"destination"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type updateItineraryRequest struct {
	UserID      *uint   `json:"user_id"`
	Title       *string `json:"title"`
	Destination *string `json:"destination"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Description *string `json:"description"`
}

func NewItineraryHandler(db *gorm.DB) *ItineraryHandler {
	return &ItineraryHandler{db: db}
}

func (h *ItineraryHandler) userExists(userID uint) (bool, error) {
	var user model.User
	err := h.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Helper mengambil user_id dan role dari c.Locals (hasil middleware JWT)
func extractUserFromLocals(c fiber.Ctx) (uint, string) {
	var userID uint
	rawID := c.Locals("user_id")
	if rawID == nil {
		rawID = c.Locals("userID")
	}

	switch v := rawID.(type) {
	case uint:
		userID = v
	case int:
		if v > 0 {
			userID = uint(v)
		}
	case int64:
		if v > 0 {
			userID = uint(v)
		}
	case float64:
		if v > 0 {
			userID = uint(v)
		}
	case string:
		if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
			userID = uint(parsed)
		}
	}

	role, _ := c.Locals("role").(string)
	return userID, role
}

func (h *ItineraryHandler) validateCreateInput(req createItineraryRequest, resolvedUserID uint) map[string]string {
	errs := make(map[string]string)

	if resolvedUserID == 0 {
		errs["user_id"] = "user_id wajib tersedia dari token login atau body request"
	}

	trimmedTitle := strings.TrimSpace(req.Title)
	if trimmedTitle == "" {
		errs["title"] = "title wajib diisi"
	} else if len(trimmedTitle) < 3 {
		errs["title"] = "title minimal 3 karakter"
	} else if len(trimmedTitle) > 150 {
		errs["title"] = "title maksimal 150 karakter"
	}

	if len(strings.TrimSpace(req.Destination)) > 150 {
		errs["destination"] = "destination maksimal 150 karakter"
	}

	if len(strings.TrimSpace(req.StartDate)) > 30 {
		errs["start_date"] = "start_date maksimal 30 karakter"
	}

	if len(strings.TrimSpace(req.EndDate)) > 30 {
		errs["end_date"] = "end_date maksimal 30 karakter"
	}

	return errs
}

// GET /api/itineraries (Mendukung Preload User, Pagination, Search, & Filter user_id)
func (h *ItineraryHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := h.db.Model(&model.Itinerary{}).Preload("User")

	if userIDParam := strings.TrimSpace(c.Query("user_id")); userIDParam != "" {
		uid, err := strconv.ParseUint(userIDParam, 10, 64)
		if err != nil || uid == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Parameter user_id tidak valid",
			})
		}
		query = query.Where("user_id = ?", uint(uid))
	}

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("title LIKE ? OR destination LIKE ?", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghitung total data itineraries",
		})
	}

	var itineraries []model.Itinerary
	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&itineraries).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data itineraries",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data itineraries berhasil diambil",
		"data":    itineraries,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total_data":  total,
			"total_pages": totalPages,
		},
	})
}

// GET /api/itineraries/:id
func (h *ItineraryHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID itinerary tidak valid",
		})
	}

	var itinerary model.Itinerary
	if err := h.db.Preload("User").First(&itinerary, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Itinerary tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil detail itinerary",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Detail itinerary berhasil diambil",
		"data":    itinerary,
	})
}

// POST /api/itineraries
func (h *ItineraryHandler) Create(c fiber.Ctx) error {
	var req createItineraryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	// Prioritaskan user_id dari JWT token; fallback ke req.UserID jika dites tanpa token
	authUserID, _ := extractUserFromLocals(c)
	resolvedUserID := authUserID
	if resolvedUserID == 0 {
		resolvedUserID = req.UserID
	}

	if errs := h.validateCreateInput(req, resolvedUserID); len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validasi input gagal",
			"errors":  errs,
		})
	}

	exists, err := h.userExists(resolvedUserID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa data user",
		})
	}
	if !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User dengan user_id tersebut tidak ditemukan",
		})
	}

	itinerary := model.Itinerary{
		UserID:      resolvedUserID,
		Title:       strings.TrimSpace(req.Title),
		Destination: strings.TrimSpace(req.Destination),
		StartDate:   strings.TrimSpace(req.StartDate),
		EndDate:     strings.TrimSpace(req.EndDate),
		Description: strings.TrimSpace(req.Description),
	}

	if err := h.db.Create(&itinerary).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan itinerary",
		})
	}

	_ = h.db.Preload("User").First(&itinerary, itinerary.ID).Error

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Itinerary berhasil dibuat",
		"data":    itinerary,
	})
}

// PUT /api/itineraries/:id
func (h *ItineraryHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID itinerary tidak valid",
		})
	}

	var itinerary model.Itinerary
	if err := h.db.First(&itinerary, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Itinerary tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data itinerary",
		})
	}

	// Cek otorisasi jika JWT middleware aktif: hanya pemilik atau admin yang boleh mengubah
	authUserID, role := extractUserFromLocals(c)
	if authUserID != 0 && itinerary.UserID != authUserID && role != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki hak akses untuk mengubah itinerary ini",
		})
	}

	var req updateItineraryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	errs := make(map[string]string)

	if req.Title != nil {
		trimmedTitle := strings.TrimSpace(*req.Title)
		if trimmedTitle == "" {
			errs["title"] = "title tidak boleh kosong"
		} else if len(trimmedTitle) < 3 {
			errs["title"] = "title minimal 3 karakter"
		} else if len(trimmedTitle) > 150 {
			errs["title"] = "title maksimal 150 karakter"
		} else {
			itinerary.Title = trimmedTitle
		}
	}

	if req.Destination != nil {
		trimmedDest := strings.TrimSpace(*req.Destination)
		if len(trimmedDest) > 150 {
			errs["destination"] = "destination maksimal 150 karakter"
		} else {
			itinerary.Destination = trimmedDest
		}
	}

	if req.StartDate != nil {
		trimmedStart := strings.TrimSpace(*req.StartDate)
		if len(trimmedStart) > 30 {
			errs["start_date"] = "start_date maksimal 30 karakter"
		} else {
			itinerary.StartDate = trimmedStart
		}
	}

	if req.EndDate != nil {
		trimmedEnd := strings.TrimSpace(*req.EndDate)
		if len(trimmedEnd) > 30 {
			errs["end_date"] = "end_date maksimal 30 karakter"
		} else {
			itinerary.EndDate = trimmedEnd
		}
	}

	if req.Description != nil {
		itinerary.Description = strings.TrimSpace(*req.Description)
	}

	if len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validasi input gagal",
			"errors":  errs,
		})
	}

	if err := h.db.Save(&itinerary).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memperbarui itinerary",
		})
	}

	_ = h.db.Preload("User").First(&itinerary, itinerary.ID).Error

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Itinerary berhasil diperbarui",
		"data":    itinerary,
	})
}

// DELETE /api/itineraries/:id
func (h *ItineraryHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID itinerary tidak valid",
		})
	}

	var itinerary model.Itinerary
	if err := h.db.First(&itinerary, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Itinerary tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa data itinerary",
		})
	}

	// Cek otorisasi jika JWT middleware aktif: hanya pemilik atau admin yang boleh menghapus
	authUserID, role := extractUserFromLocals(c)
	if authUserID != 0 && itinerary.UserID != authUserID && role != "admin" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Anda tidak memiliki hak akses untuk menghapus itinerary ini",
		})
	}

	if err := h.db.Delete(&itinerary).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus itinerary",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Itinerary berhasil dihapus",
	})
}