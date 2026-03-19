package hipay

import (
	"net/http"
)

// API endpoint constants. See: https://developers.hipay.mn
var (
	// HipayCheckout creates a checkout session (E-Commerce).
	HipayCheckout = API{
		Url:    "/checkout",
		Method: http.MethodPost,
	}
	// HipayCheckoutGet retrieves checkout status (E-Commerce).
	HipayCheckoutGet = API{
		Url:    "/checkout/get/",
		Method: http.MethodGet,
	}
	// HipayCheckoutCancel cancels a checkout (E-Commerce).
	HipayCheckoutCancel = API{
		Url:    "/checkout/cancel",
		Method: http.MethodPost,
	}
	// HipayPaymentGet retrieves payment details (Үндсэн API).
	HipayPaymentGet = API{
		Url:    "/payment/get/",
		Method: http.MethodGet,
	}
	// HipayPaymentCorrection performs a payment correction (Pos).
	HipayPaymentCorrection = API{
		Url:    "/pos/correction",
		Method: http.MethodPost,
	}
	// HipayPaymentCancel cancels a payment (Үндсэн API).
	HipayPaymentCancel = API{
		Url:    "/payment/cancel",
		Method: http.MethodPost,
	}
	// HipayStatement retrieves statement data for a given date (Pos).
	HipayStatement = API{
		Url:    "/pos/statement",
		Method: http.MethodPost,
	}
)
