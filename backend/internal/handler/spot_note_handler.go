package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"travel-itinerary-d1/backend/internal/model"
)

type SpotNoteHandler struct {
	db *gorm.DB
}

type createSpotNoteRequest struct {
	SpotID       uint   `json:"spot_id"`
	Title        string `json:"title"`
	NoteContent  string `json:"note_content"`
	OpeningHours string `json:"opening_hours"`
}

type updateSpotNoteRequest struct {
	SpotID       *uint   `json:"spot_id"`
	Title        *string `json:"title"`
	NoteContent  *string `json:"note_content"`
	OpeningHours *string `json:"opening_hours"`
}

func NewSpotNoteHandler(db *gorm.DB) *SpotNoteHandler {
	return &SpotNoteHandler{db: db}
}

func (h *SpotNoteHandler) validateCreateInput(req createSpotNoteRequest) map[string]string {
	errs := make(map[string]string)

	if req.SpotID == 0 {
		errs["spot_id"] = "spot_id wajib diisi dan harus lebih dari 0"
	}
	if strings.TrimSpace(req.Title) == "" {
		errs["title"] = "title wajib diisi"
	} else if len(strings.TrimSpace(req.Title)) > 150 {
		errs["title"] = "title maksimal 150 karakter"
	}
	if strings.TrimSpace(req.NoteContent) == "" {
		errs["note_content"] = "note_content wajib diisi"
	}
	if len(strings.TrimSpace(req.OpeningHours)) > 80 {
		errs["opening_hours"] = "opening_hours maksimal 80 karakter"
	}

	return errs
}

func (h *SpotNoteHandler) spotExists(spotID uint) (bool, error) {
	var spot model.Spot
	err := h.db.First(&spot, spotID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GET /api/spot-notes (bisa difilter dengan ?spot_id=1)
func (h *SpotNoteHandler) GetAll(c fiber.Ctx) error {
	var notes []model.SpotNote
	query := h.db.Preload("Spot").Order("id DESC")

	if spotIDParam := strings.TrimSpace(c.Query("spot_id")); spotIDParam != "" {
		spotID, err := strconv.ParseUint(spotIDParam, 10, 64)
		if err != nil || spotID == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Parameter spot_id tidak valid",
			})
		}
		query = query.Where("spot_id = ?", uint(spotID))
	}

	if err := query.Find(&notes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data spot notes",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data spot notes berhasil diambil",
		"data":    notes,
	})
}

// GET /api/spot-notes/:id
func (h *SpotNoteHandler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID spot note tidak valid",
		})
	}

	var note model.SpotNote
	if err := h.db.Preload("Spot").First(&note, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Spot note tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil detail spot note",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Detail spot note berhasil diambil",
		"data":    note,
	})
}

// POST /api/spot-notes
func (h *SpotNoteHandler) Create(c fiber.Ctx) error {
	var req createSpotNoteRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	if errs := h.validateCreateInput(req); len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validasi input gagal",
			"errors":  errs,
		})
	}

	exists, err := h.spotExists(req.SpotID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa data spot",
			"error":   err.Error(),
		})
	}
	if !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Spot dengan spot_id tersebut tidak ditemukan",
		})
	}

	note := model.SpotNote{
		SpotID:       req.SpotID,
		Title:        strings.TrimSpace(req.Title),
		NoteContent:  strings.TrimSpace(req.NoteContent),
		OpeningHours: strings.TrimSpace(req.OpeningHours),
	}

	if err := h.db.Create(&note).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan spot note",
			"error":   err.Error(),
		})
	}

	_ = h.db.Preload("Spot").First(&note, note.ID).Error

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Spot note berhasil dibuat",
		"data":    note,
	})
}

// PUT /api/spot-notes/:id
func (h *SpotNoteHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID spot note tidak valid",
		})
	}

	var note model.SpotNote
	if err := h.db.First(&note, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Spot note tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data spot note",
			"error":   err.Error(),
		})
	}

	var req updateSpotNoteRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	errs := make(map[string]string)

	if req.SpotID != nil {
		if *req.SpotID == 0 {
			errs["spot_id"] = "spot_id harus lebih dari 0"
		} else {
			exists, err := h.spotExists(*req.SpotID)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Gagal memeriksa data spot",
					"error":   err.Error(),
				})
			}
			if !exists {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"message": "Spot dengan spot_id tersebut tidak ditemukan",
				})
			}
			note.SpotID = *req.SpotID
		}
	}

	if req.Title != nil {
		trimmedTitle := strings.TrimSpace(*req.Title)
		if trimmedTitle == "" {
			errs["title"] = "title tidak boleh kosong"
		} else if len(trimmedTitle) > 150 {
			errs["title"] = "title maksimal 150 karakter"
		} else {
			note.Title = trimmedTitle
		}
	}

	if req.NoteContent != nil {
		trimmedContent := strings.TrimSpace(*req.NoteContent)
		if trimmedContent == "" {
			errs["note_content"] = "note_content tidak boleh kosong"
		} else {
			note.NoteContent = trimmedContent
		}
	}

	if req.OpeningHours != nil {
		trimmedHours := strings.TrimSpace(*req.OpeningHours)
		if len(trimmedHours) > 80 {
			errs["opening_hours"] = "opening_hours maksimal 80 karakter"
		} else {
			note.OpeningHours = trimmedHours
		}
	}

	if len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validasi input gagal",
			"errors":  errs,
		})
	}

	if err := h.db.Save(&note).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memperbarui spot note",
			"error":   err.Error(),
		})
	}

	_ = h.db.Preload("Spot").First(&note, note.ID).Error

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Spot note berhasil diperbarui",
		"data":    note,
	})
}

// DELETE /api/spot-notes/:id
func (h *SpotNoteHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID spot note tidak valid",
		})
	}

	var note model.SpotNote
	if err := h.db.First(&note, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Spot note tidak ditemukan",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa spot note",
			"error":   err.Error(),
		})
	}

	if err := h.db.Delete(&note).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus spot note",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Spot note berhasil dihapus",
	})
}
