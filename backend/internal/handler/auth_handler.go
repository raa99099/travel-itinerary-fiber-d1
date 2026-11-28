package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"travel-itinerary-d1/backend/internal/model"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Membuat JWT untuk user yang berhasil login atau registrasi.
func generateToken(user model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET belum dikonfigurasi")
	}

	hours := 24
	if value := os.Getenv("JWT_EXPIRES_HOURS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			return "", errors.New("JWT_EXPIRES_HOURS tidak valid")
		}
		hours = parsed
	}

	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(time.Duration(hours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// POST /api/auth/register
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req registerRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	errs := make(map[string]string)

	if req.Name == "" {
		errs["name"] = "Nama wajib diisi"
	} else if len(req.Name) > 100 {
		errs["name"] = "Nama maksimal 100 karakter"
	}

	if req.Email == "" {
		errs["email"] = "Email wajib diisi"
	} else if len(req.Email) > 120 ||
		!strings.Contains(req.Email, "@") {
		errs["email"] = "Format email tidak valid"
	}

	if len(req.Password) < 8 {
		errs["password"] = "Password minimal 8 karakter"
	}

	if len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Validasi input gagal",
			"errors":  errs,
		})
	}

	var count int64
	if err := h.db.Model(&model.User{}).
		Where("email = ?", req.Email).
		Count(&count).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memeriksa email pengguna",
		})
	}

	if count > 0 {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Email sudah terdaftar",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengenkripsi password",
		})
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := h.db.Create(&user).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan data pengguna",
		})
	}

	token, err := generateToken(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Akun berhasil dibuat, tetapi token gagal dibuat. Periksa konfigurasi JWT_SECRET.",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// POST /api/auth/login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req loginRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid",
		})
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "Email dan password wajib diisi",
		})
	}

	var user model.User

	if err := h.db.Where("email = ?", req.Email).
		First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Email atau password salah",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pengguna",
		})
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Email atau password salah",
		})
	}

	token, err := generateToken(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat token autentikasi",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

// GET /api/auth/me
func (h *AuthHandler) Me(c fiber.Ctx) error {
	rawID := c.Locals("user_id")

	id, err := strconv.ParseUint(
		strings.TrimSpace(strings.ReplaceAll(
			strings.TrimSpace(toString(rawID)), " ", "",
		)),
		10,
		64,
	)

	if err != nil || id == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Identitas pengguna tidak valid",
		})
	}

	var user model.User

	if err := h.db.First(&user, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pengguna tidak ditemukan",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pengguna",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pengguna berhasil diambil",
		"data":    user,
	})
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case float64:
		return strconv.FormatUint(uint64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return ""
	}
}
