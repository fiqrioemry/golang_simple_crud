package handlers

import (
	"server/internal/services"

	"github.com/gin-gonic/gin"
)

type FormHandler struct {
	service services.FormService
}

func NewFormHandler(service services.FormService) *FormHandler {
	return &FormHandler{service}
}

// user dapat membuat form sesuai tipe kebutuhannya (order_form, quisoner, feedback, poll, quiz, surver, exam, screening diagnose)
func (h *FormHandler) CreateNewForm(c *gin.Context) {
}

func (h *FormHandler) UpdateForm(c *gin.Context) {
}

func (h *FormHandler) GetAllForms(c *gin.Context) {
}

// mengambil detail informasi form berikut dengan jumlah responden, kebutuhan untuk statistik langsung
func (h *FormHandler) GetFormDetail(c *gin.Context) {
}

// Endpoint publik untuk mengakses form berdasarkan slug (tanpa login)
func (h *FormHandler) GetPublicFormBySlug(c *gin.Context) {}

// Endpoint untuk pengisi form mengirimkan jawaban (tanpa login app)
func (h *FormHandler) SubmitFormResponse(c *gin.Context) {}

// Melihat seluruh hasil isian dari form (khusus owner)
func (h *FormHandler) GetFormResponses(c *gin.Context) {}

// Ambil statistik hasil seperti jumlah jawaban per opsi, rata-rata, skor, dsb.
func (h *FormHandler) GetFormStats(c *gin.Context) {}

// Ambil statistik hasil seperti jumlah jawaban per opsi, rata-rata, skor, dsb.
func (h *FormHandler) DeleteForm(c *gin.Context) {
}
