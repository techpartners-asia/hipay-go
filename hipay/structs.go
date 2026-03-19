package hipay

const (
	// HttpContent is the Content-Type value for HiPay API requests.
	HttpContent = "application/json"
)

type (
	// API describes a HiPay API endpoint (URL path and HTTP method).
	API struct {
		Url    string // API path (e.g. /checkout)
		Method string // HTTP method (GET, POST)
	}

	// HipayCheckoutRequest is the request body for creating a checkout.
	// See: https://developers.hipay.mn (E-Commerce)
	HipayCheckoutRequest struct {
		EntityID  string          `json:"entityId"`  // Merchant entity ID
		Amount    float64         `json:"amount"`   // Amount in MNT
		Currency  string          `json:"currency"` // e.g. MNT
		QrData    bool             `json:"qrData"`   // Request QR data in response
		Signal    bool             `json:"signal"`   // Request push/signal channel
		IpAddress string          `json:"ipaddress"`
		Items    []*CheckOutItem   `json:"items,omitempty"` // Line items (optional)
	}

	// CheckOutItem represents a line item in a checkout request.
	CheckOutItem struct {
		ItemNo   string  `json:"itemno"`
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		Quantity int64   `json:"quantity"`
		Brand    string  `json:"brand"`
		Measure  string  `json:"measure"`
		Vat      float64 `json:"vat"`
		CityTax  float64 `json:"citytax"`
	}

	// HipayCheckoutResponse is the response from checkout create.
	// See: https://developers.hipay.mn (E-Commerce)
	HipayCheckoutResponse struct {
		Code        int64      `json:"code"`        // 1 = success
		Description string     `json:"description"`
		RequestID   string     `json:"requestId"`
		CheckoutID  string     `json:"checkoutId"`  // Use for CheckoutGet / CheckoutCancel
		Expires     string     `json:"expires"`
		Signal      *Signal    `json:"signal,omitempty"` // Present when Signal requested
		QrData      string     `json:"qrData"`      // QR payload for display
		Message     string     `json:"message,omitempty"`
		Details     []*Details `json:"details,omitempty"` // Error details when code != 1
	}

	// Signal holds push/signal channel data (subscribeKey, channel, uuid).
	Signal struct {
		SubscribeKey string `json:"subscribeKey"`
		Channel      string `json:"channel"`
		UUID         string `json:"uuid"`
	}

	// HipayCheckoutGetResponse is the response from checkout get.
	// See: https://developers.hipay.mn (E-Commerce)
	HipayCheckoutGetResponse struct {
		Code           int64      `json:"code"`
		Description    string     `json:"description"`
		Amount         float64    `json:"amount,omitempty"`
		Currency       string     `json:"currency,omitempty"`
		DiscountAmount float64    `json:"discount_amount,omitempty"`
		Status         string     `json:"status,omitempty"`     // Checkout status
		StatusDate     string     `json:"status_date,omitempty"`
		PaymentID      string     `json:"paymentId,omitempty"`
		PaymentType    string     `json:"paymentType,omitempty"`
		Message        string     `json:"message,omitempty"`
		Details        []*Details `json:"details,omitempty"`
	}

	// Details holds a single error detail (field and issue) from the API.
	Details struct {
		Field string `json:"field"` // Field name
		Issue string `json:"issue"` // Error description
	}

	// HipayPaymentGetResponse is the response from payment get.
	// See: https://developers.hipay.mn (Үндсэн API)
	HipayPaymentGetResponse struct {
		Code         int64      `json:"code"`
		Description  string     `json:"description"`
		ID           string     `json:"id,omitempty"`
		Amount       float64    `json:"amount,omitempty"`
		Currency     string     `json:"currency,omitempty"`
		EntityID     string     `json:"entityId,omitempty"`
		CheckoutID   string     `json:"checkoutId,omitempty"`
		PaymentID    string     `json:"paymentId,omitempty"`
		PaymentType  string     `json:"paymentType,omitempty"`
		PaymentBrand string     `json:"paymentBrand,omitempty"`
		PaymentDate  string     `json:"paymentDate,omitempty"`
		PaymentDesc  string     `json:"paymentDesc,omitempty"`
		ResultDesc   string     `json:"result_desc,omitempty"`
		ResultCode   string     `json:"result_code,omitempty"`
		Message      string     `json:"message,omitempty"`
		Details      []*Details `json:"details,omitempty"`
	}

	// HipayPaymentCorrectionRequest is the request body for payment correction (POS).
	// See: https://developers.hipay.mn (Pos)
	HipayPaymentCorrectionRequest struct {
		EntityID  string `json:"entityId"`
		PaymentID string `json:"paymentId"`
	}

	// HipayPaymentCorrectionResponse is the response from payment correction.
	HipayPaymentCorrectionResponse struct {
		Code                int64      `json:"code"`
		Description         string     `json:"description"`
		PaymentID           string     `json:"paymentId,omitempty"`
		CorrectionPaymentID string     `json:"correction_paymentId,omitempty"`
		Message             string     `json:"message,omitempty"`
		Details             []*Details `json:"details,omitempty"`
	}

	// HipayStatementRequest is the request body for statement.
	// See: https://developers.hipay.mn (Pos)
	HipayStatementRequest struct {
		EntityID string `json:"entityId"`
		Date     string `json:"date"` // Format: YYYY-MM-DD
	}

	// HipayStatementResponse is the response from statement.
	// When code is 1, Data is guaranteed non-nil by the SDK.
	HipayStatementResponse struct {
		Code        int64      `json:"code"`
		Description string     `json:"description"`
		Message     string     `json:"message,omitempty"`
		Details     []*Details `json:"details,omitempty"`
		Data        *Data      `json:"data,omitempty"` // List and pagination
	}

	// Data wraps the statement list and pagination info.
	Data struct {
		List       []*List `json:"list"`       // Payment rows
		EntityID   string  `json:"entityId"`
		Date       string  `json:"date"`
		Page       int64   `json:"page"`
		PerPage    int64   `json:"perPage"`
		TotalCount int64   `json:"totalCount"`
		TotalPage  int64   `json:"totalPage"`
	}

	// List is a single statement line (payment record).
	List struct {
		PaymentDate  string  `json:"paymentDate"`
		CheckoutID   string  `json:"checkoutId"`
		PaymentID    string  `json:"paymentId"`
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		FeeAmount    float64 `json:"feeAmount"`
		FeePrc       float64 `json:"feePrc"`
		PaymentDesc  string  `json:"paymentDesc"`
		ReturnAmount float64 `json:"returnAmount"`
	}

	// HipayPaymentCancelRequest is the request body for payment cancel.
	// See: https://developers.hipay.mn (Үндсэн API)
	HipayPaymentCancelRequest struct {
		EntityID  string `json:"entityId"`
		PaymentID string `json:"paymentId"`
	}

	// HipayPaymentCancelResponse is the response from payment cancel.
	HipayPaymentCancelResponse struct {
		Code        int64      `json:"code"`
		Description string     `json:"description"`
		PaymentID   string     `json:"paymentId,omitempty"`
		CancelledAt string     `json:"cancelledAt,omitempty"`
		Message     string     `json:"message,omitempty"`
		Details     []*Details `json:"details,omitempty"`
	}

	// HipayCheckoutCancelRequest is the request body for checkout cancel.
	// See: https://developers.hipay.mn (E-Commerce)
	HipayCheckoutCancelRequest struct {
		EntityID   string `json:"entityId"`
		CheckoutID string `json:"checkoutId"`
	}

	// HipayCheckoutCancelResponse is the response from checkout cancel.
	HipayCheckoutCancelResponse struct {
		Code        int64      `json:"code"`
		Description string     `json:"description"`
		Message     string     `json:"message,omitempty"`
		Details     []*Details `json:"details,omitempty"`
	}
)
