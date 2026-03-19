// Package hipay provides a Go SDK for the HiPay payment API.
// It supports checkout (E-Commerce), payment get/cancel/correction, and statement (POS).
//
// See: https://developers.hipay.mn
package hipay

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"
)

type hipay struct {
	endpoint string
	token    string
	entityId string
	client   *resty.Client
}

// Hipay is the HiPay SDK interface.
type Hipay interface {
	// Checkout creates a checkout session and returns QR data / checkout ID.
	// See: https://developers.hipay.mn (E-Commerce)
	Checkout(amount float64) (HipayCheckoutResponse, error)

	// CheckoutGet returns the current status and details of a checkout.
	// See: https://developers.hipay.mn (E-Commerce)
	CheckoutGet(checkoutId string) (HipayCheckoutGetResponse, error)

	// CheckoutCancel cancels an existing checkout.
	// See: https://developers.hipay.mn (E-Commerce)
	CheckoutCancel(checkoutId string) (HipayCheckoutCancelResponse, error)

	// PaymentGet returns payment details by payment ID.
	// See: https://developers.hipay.mn (Үндсэн API)
	PaymentGet(paymentId string) (HipayPaymentGetResponse, error)

	// PaymentCorrection performs a correction for a payment (POS).
	// See: https://developers.hipay.mn (Pos)
	PaymentCorrection(paymentId string) (HipayPaymentCorrectionResponse, error)

	// PaymentCancel cancels a payment.
	// See: https://developers.hipay.mn (Үндсэн API)
	PaymentCancel(paymentId string) (HipayPaymentCancelResponse, error)

	// Statement returns statement data for the given date (format: YYYY-MM-DD).
	// See: https://developers.hipay.mn (Pos)
	Statement(date string) (HipayStatementResponse, error)
}

// Option configures the HiPay client at initialization.
type Option func(*hipay)

// WithClient sets a custom resty.Client (e.g. for timeouts, TLS, or logging).
// If client is nil, it is ignored.
func WithClient(client *resty.Client) Option {
	return func(h *hipay) {
		if client != nil {
			h.client = client
		}
	}
}

// New builds a new HiPay client.
//
//   - endpoint: HiPay API base URL (sandbox or production)
//   - token: API auth token
//   - entityId: merchant / entity ID
//   - options: optional modifiers (e.g. WithClient); nil entries are skipped
//
// See: https://developers.hipay.mn
func New(endpoint, token, entityId string, options ...Option) Hipay {
	h := &hipay{
		endpoint: endpoint,
		token:    token,
		entityId: entityId,
		client:   resty.New().SetTimeout(60 * time.Second),
	}

	for _, opt := range options {
		if opt != nil {
			opt(h)
		}
	}

	return h
}

func (h *hipay) Checkout(amount float64) (HipayCheckoutResponse, error) {
	request := HipayCheckoutRequest{
		EntityID: h.entityId,
		Amount:   amount,
		Currency: "MNT",
		QrData:   true,
		Signal:   false,
	}

	var response HipayCheckoutResponse
	err := h.httpRequestHipay(request, &response, HipayCheckout, "")
	if err != nil {
		return HipayCheckoutResponse{}, err
	}

	if response.Code != 1 {
		return HipayCheckoutResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	return response, nil
}

func (h *hipay) CheckoutGet(checkoutId string) (HipayCheckoutGetResponse, error) {
	ext := checkoutId + "?entityId=" + h.entityId
	var response HipayCheckoutGetResponse
	err := h.httpRequestHipay(nil, &response, HipayCheckoutGet, ext)
	if err != nil {
		return HipayCheckoutGetResponse{}, err
	}

	if response.Code != 1 {
		return HipayCheckoutGetResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	return response, nil
}

func (h *hipay) CheckoutCancel(checkoutId string) (HipayCheckoutCancelResponse, error) {
	body := HipayCheckoutCancelRequest{
		EntityID:   h.entityId,
		CheckoutID: checkoutId,
	}
	var response HipayCheckoutCancelResponse
	err := h.httpRequestHipay(body, &response, HipayCheckoutCancel, "")
	if err != nil {
		return HipayCheckoutCancelResponse{}, err
	}

	if response.Code != 1 {
		return HipayCheckoutCancelResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	return response, nil
}

func (h *hipay) PaymentGet(paymentId string) (HipayPaymentGetResponse, error) {
	ext := paymentId + "?entityId=" + h.entityId
	var response HipayPaymentGetResponse
	err := h.httpRequestHipay(nil, &response, HipayPaymentGet, ext)
	if err != nil {
		return HipayPaymentGetResponse{}, err
	}

	if response.Code != 1 {
		return HipayPaymentGetResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	return response, nil
}

func (h *hipay) PaymentCorrection(paymentId string) (HipayPaymentCorrectionResponse, error) {
	body := HipayPaymentCorrectionRequest{
		EntityID:  h.entityId,
		PaymentID: paymentId,
	}
	var response HipayPaymentCorrectionResponse
	err := h.httpRequestHipay(body, &response, HipayPaymentCorrection, "")
	if err != nil {
		return HipayPaymentCorrectionResponse{}, err
	}

	if response.Code != 1 {
		return HipayPaymentCorrectionResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	return response, nil
}

func (h *hipay) Statement(date string) (HipayStatementResponse, error) {
	body := HipayStatementRequest{
		EntityID: h.entityId,
		Date:     date,
	}
	var response HipayStatementResponse
	err := h.httpRequestHipay(body, &response, HipayStatement, "")
	if err != nil {
		return HipayStatementResponse{}, err
	}

	if response.Code != 1 {
		return HipayStatementResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}
	if response.Data == nil {
		return HipayStatementResponse{}, errors.New("hipay: API returned success but data is nil")
	}
	return response, nil
}

func (h *hipay) PaymentCancel(paymentId string) (HipayPaymentCancelResponse, error) {
	if paymentId == "" {
		return HipayPaymentCancelResponse{}, errors.New("PaymentID cannot be empty")
	}
	body := HipayPaymentCancelRequest{
		EntityID:  h.entityId,
		PaymentID: paymentId,
	}
	var response HipayPaymentCancelResponse
	err := h.httpRequestHipay(body, &response, HipayPaymentCancel, "")
	if err != nil {
		return HipayPaymentCancelResponse{}, err
	}

	if response.Code != 1 {
		return HipayPaymentCancelResponse{}, h.formatError(response.Description, response.Message, response.Details)
	}

	return response, nil
}

// formatError builds an error from API description, message, and optional details.
// It safely handles nil elements in details (e.g. API returning [null]).
func (h *hipay) formatError(desc, msg string, details []*Details) error {
	errMsg := desc
	if msg != "" {
		errMsg = msg
	}
	if len(details) > 0 && details[0] != nil {
		errMsg += ": " + details[0].Field + " - " + details[0].Issue
	}
	return errors.New(errMsg)
}

// httpRequestHipay sends a request to the HiPay API and unmarshals the response.
// It supports both single-object and array-wrapped JSON responses.
func (h *hipay) httpRequestHipay(body interface{}, result interface{}, api API, urlExt string) error {
	url := h.endpoint + api.Url + urlExt
	req := h.client.R().
		SetHeader("Content-Type", HttpContent).
		SetAuthToken(h.token)

	if body != nil {
		req.SetBody(body)
	}

	res, err := req.Execute(api.Method, url)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("%s-HiPay response error: %s (Status: %d)",
			time.Now().Format("2006-01-02 15:04:05"),
			res.String(),
			res.StatusCode())
	}

	// HiPay API often wraps the response in an array: [ { ... } ]
	// We handle both single object and array-wrapped responses.
	respBody := res.Bytes()
	trimmed := strings.TrimSpace(string(respBody))
	if strings.HasPrefix(trimmed, "[") {
		// Temporary slice to unmarshal the array
		wrapper := []json.RawMessage{}
		if err := json.Unmarshal(respBody, &wrapper); err != nil {
			return err
		}
		if len(wrapper) > 0 {
			if err := json.Unmarshal(wrapper[0], result); err != nil {
				return err
			}
		}
	} else {
		if err := json.Unmarshal(respBody, result); err != nil {
			return err
		}
	}

	return nil
}
