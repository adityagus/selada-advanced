package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-api/config"
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// WhatsAppWebhookPayload supports multiple popular WA gateways:
// 1. Fonnte (sender, message)
// 2. Wablas (phone, message)
// 3. Meta Cloud API (entry.changes.value.messages & entry.changes.value.statuses)
// 4. Custom generic JSON (phone, message, status, id_rencana)
type WhatsAppWebhookPayload struct {
	Sender    string `json:"sender"`
	Phone     string `json:"phone"`
	From      string `json:"from"`
	HP        string `json:"hp"`
	Message   string `json:"message"`
	Body      string `json:"body"`
	Text      string `json:"text"`
	Status    string `json:"status"` // delivered, read, sent, tervalidasi
	IDRencana int    `json:"id_rencana"`

	// Meta Cloud API structure
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
				Statuses []struct {
					ID          string `json:"id"`
					Status      string `json:"status"` // delivered, read, sent, failed
					RecipientID string `json:"recipient_id"`
				} `json:"statuses"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// CleanPhoneNumber strips +, spaces, dashes and standardizes 62 / 08 prefixes
func CleanPhoneNumber(phone string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	cleaned := reg.ReplaceAllString(phone, "")

	if strings.HasPrefix(cleaned, "628") {
		return "08" + cleaned[3:]
	}
	if strings.HasPrefix(cleaned, "62") && len(cleaned) > 9 {
		return "0" + cleaned[2:]
	}
	return cleaned
}

// WhatsAppWebhookVerifyHandler handles GET /api/webhook/whatsapp
// Useful for Meta Webhook verification challenge & health-check
func WhatsAppWebhookVerifyHandler(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	verifyToken := "selada_wa_webhook_token" // Default token or from env

	if mode == "subscribe" && token == verifyToken {
		c.String(http.StatusOK, challenge)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"service":     "Selada WhatsApp Webhook Listener",
		"timestamp":   time.Now().Format(time.RFC3339),
		"description": "Ready to receive incoming customer WhatsApp events",
	})
}

// WhatsAppWebhookHandler handles POST /api/webhook/whatsapp
// When a customer replies or triggers a WA event (delivered/read/message),
// WA is marked as 'tervalidasi' and status is upgraded: Leads (1) -> Prospect (2)
func WhatsAppWebhookHandler(c *gin.Context) {
	if config.DBSelada == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Selada tidak terhubung"})
		return
	}

	var payload WhatsAppWebhookPayload
	rawData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca body request"})
		return
	}

	_ = json.Unmarshal(rawData, &payload)

	// Fallback to Form Post if JSON is empty (e.g. Fonnte default)
	if payload.Sender == "" && payload.Phone == "" && payload.From == "" {
		payload.Sender = c.PostForm("sender")
		payload.Phone = c.PostForm("phone")
		payload.Message = c.PostForm("message")
		if payload.Message == "" {
			payload.Message = c.PostForm("body")
		}
		if payload.Status == "" {
			payload.Status = c.PostForm("status")
		}
	}

	// Extract Sender and Message across multiple provider schemas
	senderPhone := payload.Sender
	if senderPhone == "" {
		senderPhone = payload.Phone
	}
	if senderPhone == "" {
		senderPhone = payload.From
	}
	if senderPhone == "" {
		senderPhone = payload.HP
	}

	messageText := payload.Message
	if messageText == "" {
		messageText = payload.Body
	}
	if messageText == "" {
		messageText = payload.Text
	}

	deliveryStatus := strings.ToLower(strings.TrimSpace(payload.Status))

	// Extract from Meta Cloud API structure if present
	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 {
		val := payload.Entry[0].Changes[0].Value
		if len(val.Messages) > 0 {
			senderPhone = val.Messages[0].From
			messageText = val.Messages[0].Text.Body
		}
		if len(val.Statuses) > 0 {
			if senderPhone == "" {
				senderPhone = val.Statuses[0].RecipientID
			}
			deliveryStatus = strings.ToLower(val.Statuses[0].Status)
		}
	}

	cleanPhone := CleanPhoneNumber(senderPhone)

	// Also extract #ID if customer replied with reference, e.g. "Halo #12345"
	targetIDRencana := payload.IDRencana
	if targetIDRencana == 0 && messageText != "" {
		reID := regexp.MustCompile(`#(\d+)`)
		match := reID.FindStringSubmatch(messageText)
		if len(match) > 1 {
			targetIDRencana, _ = strconv.Atoi(match[1])
		}
	}

	log.Printf("[WEBHOOK WA RECEIVED] Sender: %s (clean: %s), Text: %s, Status: %s, RefID: %d", senderPhone, cleanPhone, messageText, deliveryStatus, targetIDRencana)

	// Step 1: Find Customer & Latest Rencana
	var rencana models.Rencana
	var customer models.Customer

	dbQuery := config.DBSelada.Table("tb_rencana r").
		Joins("JOIN tb_customer c ON c.id_customer = r.id_customer")

	if targetIDRencana > 0 {
		err = dbQuery.Where("r.id_rencana = ?", targetIDRencana).
			Select("r.*").First(&rencana).Error
	} else if cleanPhone != "" {
		// Search by last 8-10 digits to match 08xx or 62xx
		phoneSuffix := cleanPhone
		if len(phoneSuffix) > 8 {
			phoneSuffix = phoneSuffix[len(phoneSuffix)-8:]
		}
		err = dbQuery.Where("c.hp LIKE ?", "%"+phoneSuffix+"%").
			Order("r.id_rencana DESC").
			Select("r.*").First(&rencana).Error
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nomor telepon pengirim atau ID rencana tidak ditemukan dalam payload",
		})
		return
	}

	if err != nil || rencana.IDRencana == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ignored",
			"message": fmt.Sprintf("Tidak ada inquiry aktif yang cocok dengan nomor %s", senderPhone),
		})
		return
	}

	_ = config.DBSelada.Where("id_customer = ?", rencana.IDCustomer).First(&customer)

	// Step 2: Validasi WA
	// Whenever delivery succeeded (delivered/read/sent) or an inbound message arrives,
	// mark customer's status_wa as 'tervalidasi'
	now := time.Now()
	waStatusUpdated := false
	if messageText != "" || deliveryStatus == "delivered" || deliveryStatus == "read" || deliveryStatus == "sent" || deliveryStatus == "tervalidasi" || deliveryStatus == "valid" {
		_ = config.DBSelada.Table("tb_customer").
			Where("id_customer = ?", customer.IDCustomer).
			Updates(map[string]interface{}{
				"status_wa":       "tervalidasi",
				"wa_validated_at": now,
			}).Error
		waStatusUpdated = true
	}

	// Step 3: Get Latest RencanaDet status
	var lastDet models.RencanaDet
	err = config.DBSelada.Where("id_rencana = ?", rencana.IDRencana).
		Order("id_rencana_det DESC").
		First(&lastDet).Error

	currentHslAktiv := 1
	if err == nil {
		currentHslAktiv = lastDet.IDHslaktiv
	}

	// Step 4: Validasi WA dan Leads => Prospect adalah SATU KESATUAN
	// Pipeline: Leads (1) -> Prospect (2) -> SBG / Booking (7)
	newHslAktiv := currentHslAktiv
	actionNote := ""

	if messageText != "" {
		actionNote = fmt.Sprintf("[Webhook WA] Balasan masuk: \"%s\"", messageText)
	} else if deliveryStatus != "" {
		actionNote = fmt.Sprintf("[Webhook WA Status] Pengiriman WA berhasil (%s)", deliveryStatus)
	} else {
		actionNote = "[Webhook WA] Event WhatsApp diterima"
	}

	if currentHslAktiv == 1 && waStatusUpdated {
		newHslAktiv = 2 // Auto-Upgrade Leads (1) -> Prospect (2)
		if messageText != "" {
			actionNote = fmt.Sprintf("[Validasi WA & Auto-Upgrade] WA tervalidasi via balasan nasabah: \"%s\". Status dinaikkan ke PROSPECT.", messageText)
		} else {
			actionNote = fmt.Sprintf("[Validasi WA & Auto-Upgrade] Nomor WA tervalidasi aktif (%s). Status dinaikkan ke PROSPECT.", deliveryStatus)
		}
	}

	// Insert new progress row into tb_rencana_det if status changed or new note added
	newDet := models.RencanaDet{
		IDRencana:    rencana.IDRencana,
		TglRencana:   now,
		IDHslaktiv:   newHslAktiv,
		Status:       1,
		KetRencana:   lastDet.KetRencana,
		KetAktivitas: actionNote,
	}

	if err := config.DBSelada.Create(&newDet).Error; err != nil {
		log.Printf("[WEBHOOK WA ERROR] Failed to create detail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status via webhook"})
		return
	}

	log.Printf("[WEBHOOK WA SUCCESS] Updated Rencana #%d (Cust: %s, WA: tervalidasi) status from %d to %d", rencana.IDRencana, customer.Nama, currentHslAktiv, newHslAktiv)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status WA tervalidasi dan inquiry berhasil diperbarui via webhook",
		"data": gin.H{
			"id_rencana":      rencana.IDRencana,
			"id_customer":     customer.IDCustomer,
			"nama_customer":   customer.Nama,
			"hp":              customer.Hp,
			"status_wa":       "tervalidasi",
			"previous_status": currentHslAktiv,
			"current_status":  newHslAktiv,
			"status_label":    map[int]string{1: "Leads", 2: "Prospect", 3: "Hot Prospect", 7: "Booking / SBG"}[newHslAktiv],
			"note":            actionNote,
			"updated_at":      now.Format("2006-01-02 15:04:05"),
		},
	})
}
