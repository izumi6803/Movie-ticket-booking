package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"cinema-backend/internal/config"
	"github.com/google/uuid"
)

type VNPayService struct {
	config config.VNPayConfig
}

func NewVNPayService(cfg config.VNPayConfig) *VNPayService {
	return &VNPayService{config: cfg}
}

type VNPayPaymentResponse struct {
	PaymentUrl string `json:"paymentUrl"`
	OrderId    string `json:"orderId"`
}

func (s *VNPayService) CreatePayment(bookingID string, amount float64, orderInfo string, clientIP string) (*VNPayPaymentResponse, error) {
	// VNPay uses VND and amount must be multiplied by 100
	amountVND := int64(amount * 100)

	orderId := uuid.New().String()
	createDate := time.Now().Format("20060102150405")

	// Build params - must be in alphabetical order for signing
	// Use bookingID as vnp_OrderInfo so we can identify the booking on return
	params := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    s.config.TmnCode,
		"vnp_Amount":     fmt.Sprintf("%d", amountVND),
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     orderId,
		"vnp_OrderInfo":  bookingID, // Use bookingID here for identification
		"vnp_OrderType":  "other",
		"vnp_Locale":     "vn",
		"vnp_ReturnUrl":  s.config.ReturnURL,
		"vnp_IpAddr":     clientIP,
		"vnp_CreateDate": createDate,
	}

	// Generate secure hash using HMAC SHA512 (VNPay standard)
	signData := s.buildSignData(params)
	signature := s.generateSignature(signData)

	// Build payment URL with signature appended at the end
	paymentUrl := s.config.Endpoint + "?" + s.buildQueryString(params) + "&vnp_SecureHash=" + signature

	return &VNPayPaymentResponse{
		PaymentUrl: paymentUrl,
		OrderId:    orderId,
	}, nil
}

func (s *VNPayService) VerifyReturn(params map[string]string) bool {
	// Remove secure hash from params
	secureHash := params["vnp_SecureHash"]
	delete(params, "vnp_SecureHash")
	delete(params, "vnp_SecureHashType")

	// Rebuild sign data
	signData := s.buildSignData(params)
	expectedHash := s.generateSignature(signData)

	return secureHash == expectedHash
}

func (s *VNPayService) buildSignData(params map[string]string) string {
	// Sort keys alphabetically
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build sign data with URL encoding like VNPay PHP demo
	// Format: urlencode(key)=urlencode(value)
	var pairs []string
	for _, k := range keys {
		if params[k] != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
		}
	}

	return strings.Join(pairs, "&")
}

func (s *VNPayService) generateSignature(data string) string {
	// VNPay uses HMAC SHA512
	h := hmac.New(sha512.New, []byte(s.config.HashSecret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *VNPayService) buildQueryString(params map[string]string) string {
	// Sort keys alphabetically
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string with URL encoding
	var pairs []string
	for _, k := range keys {
		if params[k] != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
		}
	}

	return strings.Join(pairs, "&")
}

func (s *VNPayService) IsSandbox() bool {
	return s.config.Environment == "sandbox"
}
