package gateway

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/models/postgres"
)

type MidtransGateway struct {
	config *config.MidtransConfig
}

func NewMidtransGateway(config *config.MidtransConfig) *MidtransGateway {
	return &MidtransGateway{config: config}
}

func (m *MidtransGateway) CreatePayment(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) error {
	var endpoint string
	var requestBody interface{}

	switch method.Jenis {
	case "virtual_account":
		endpoint = "/v2/charge"
		requestBody = m.buildVARequest(payment, method)
	case "qris":
		endpoint = "/v2/charge"
		requestBody = m.buildQRISRequest(payment)
	case "e_wallet":
		endpoint = "/v2/charge"
		requestBody = m.buildEWalletRequest(payment, method)
	default:
		return fmt.Errorf("unsupported payment method: %s", method.Jenis)
	}

	response, err := m.makeRequest(endpoint, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create payment: %v", err)
	}

	responseData, _ := json.Marshal(response)
	payment.GatewayResponse = string(responseData)

	if externalID, ok := response["transaction_id"].(string); ok {
		payment.ExternalID = externalID
	}

	return nil
}

func (m *MidtransGateway) buildVARequest(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) map[string]interface{} {
	return map[string]interface{}{
		"payment_type": "bank_transfer",
		"transaction_details": map[string]interface{}{
			"order_id":     payment.NomorTransaksi,
			"gross_amount": int64(payment.TotalAmount),
		},
		"bank_transfer": map[string]interface{}{
			"bank": method.BankCode,
		},
		"customer_details": map[string]interface{}{
			"first_name": payment.CustomerName,
			"email":      payment.CustomerEmail,
			"phone":      payment.CustomerPhone,
		},
		"custom_expiry": map[string]interface{}{
			"expiry_duration": 1440,
			"unit":            "minute",
		},
	}
}

func (m *MidtransGateway) buildQRISRequest(payment *postgres.PaymentTransaction) map[string]interface{} {
	return map[string]interface{}{
		"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     payment.NomorTransaksi,
			"gross_amount": int64(payment.TotalAmount),
		},
		"customer_details": map[string]interface{}{
			"first_name": payment.CustomerName,
			"email":      payment.CustomerEmail,
			"phone":      payment.CustomerPhone,
		},
		"custom_expiry": map[string]interface{}{
			"expiry_duration": 30,
			"unit":            "minute",
		},
	}
}

func (m *MidtransGateway) buildEWalletRequest(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) map[string]interface{} {
	walletType := method.WalletCode
	if walletType == "" {
		walletType = "gopay"
	}

	return map[string]interface{}{
		"payment_type": walletType,
		"transaction_details": map[string]interface{}{
			"order_id":     payment.NomorTransaksi,
			"gross_amount": int64(payment.TotalAmount),
		},
		"customer_details": map[string]interface{}{
			"first_name": payment.CustomerName,
			"email":      payment.CustomerEmail,
			"phone":      payment.CustomerPhone,
		},
		walletType: map[string]interface{}{
			"enable_callback": true,
		},
	}
}

func (m *MidtransGateway) makeRequest(endpoint string, requestBody interface{}) (map[string]interface{}, error) {
	baseURL := "https://api.sandbox.midtrans.com"
	if m.config.Environment == "production" {
		baseURL = "https://api.midtrans.com"
	}

	url := baseURL + endpoint

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+m.config.ServerKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("midtrans error: %s", string(body))
	}

	return response, nil
}

func (m *MidtransGateway) VerifySignature(data map[string]interface{}, signature string) bool {
	orderID, _ := data["order_id"].(string)
	statusCode, _ := data["status_code"].(string)
	grossAmount, _ := data["gross_amount"].(string)

	signatureKey := orderID + statusCode + grossAmount + m.config.ServerKey

	hash := sha512.Sum512([]byte(signatureKey))
	expectedSignature := fmt.Sprintf("%x", hash)

	return expectedSignature == signature
}