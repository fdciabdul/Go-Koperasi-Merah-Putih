package gateway

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"koperasi-merah-putih/config"
	"koperasi-merah-putih/internal/models/postgres"
)

type XenditGateway struct {
	config *config.XenditConfig
}

func NewXenditGateway(config *config.XenditConfig) *XenditGateway {
	return &XenditGateway{config: config}
}

func (x *XenditGateway) CreatePayment(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) error {
	var endpoint string
	var requestBody interface{}

	switch method.Jenis {
	case "virtual_account":
		endpoint = "/virtual_accounts"
		requestBody = x.buildVARequest(payment, method)
	case "qris":
		endpoint = "/qr_codes"
		requestBody = x.buildQRISRequest(payment)
	case "e_wallet":
		endpoint = "/ewallets/charges"
		requestBody = x.buildEWalletRequest(payment, method)
	default:
		return fmt.Errorf("unsupported payment method: %s", method.Jenis)
	}

	response, err := x.makeRequest(endpoint, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create payment: %v", err)
	}

	responseData, _ := json.Marshal(response)
	payment.GatewayResponse = string(responseData)

	if externalID, ok := response["id"].(string); ok {
		payment.ExternalID = externalID
	}

	return nil
}

func (x *XenditGateway) buildVARequest(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) map[string]interface{} {
	return map[string]interface{}{
		"external_id":   payment.NomorTransaksi,
		"bank_code":     method.BankCode,
		"name":          payment.CustomerName,
		"expected_amount": int64(payment.TotalAmount),
		"expiration_date": payment.ExpiredDate.Format(time.RFC3339),
		"is_closed":     true,
		"is_single_use": true,
	}
}

func (x *XenditGateway) buildQRISRequest(payment *postgres.PaymentTransaction) map[string]interface{} {
	return map[string]interface{}{
		"external_id": payment.NomorTransaksi,
		"type":        "DYNAMIC",
		"amount":      int64(payment.TotalAmount),
		"callback_url": "https://your-domain.com/webhook/xendit",
	}
}

func (x *XenditGateway) buildEWalletRequest(payment *postgres.PaymentTransaction, method *postgres.PaymentMethod) map[string]interface{} {
	walletType := method.WalletCode
	if walletType == "" {
		walletType = "OVO"
	}

	request := map[string]interface{}{
		"reference_id":   payment.NomorTransaksi,
		"currency":       "IDR",
		"amount":         int64(payment.TotalAmount),
		"checkout_method": "ONE_TIME_PAYMENT",
		"channel_code":   walletType,
		"channel_properties": map[string]interface{}{
			"success_redirect_url": "https://your-domain.com/success",
		},
	}

	if walletType == "OVO" {
		request["channel_properties"].(map[string]interface{})["mobile_number"] = payment.CustomerPhone
	}

	return request
}

func (x *XenditGateway) makeRequest(endpoint string, requestBody interface{}) (map[string]interface{}, error) {
	baseURL := "https://api.xendit.co"
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
	req.SetBasicAuth(x.config.SecretKey, "")

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
		return nil, fmt.Errorf("xendit error: %s", string(body))
	}

	return response, nil
}

func (x *XenditGateway) VerifySignature(data []byte, signature string) bool {
	h := hmac.New(sha256.New, []byte(x.config.WebhookToken))
	h.Write(data)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}