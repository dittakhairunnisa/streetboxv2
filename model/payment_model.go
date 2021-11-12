package model

// ReqCreateQRIS xendit req QRIS
type ReqCreateQRIS struct {
	Amount     int64             `json:"amount"`
	MerchantID int64             `json:"merchant_id"`
	Types      string            `json:"types"`
	Address    string            `json:"address"`
	Order      ReqTrxOrderOnline `json:"order"`
}

// ResCreateQRISXnd xendit res create QRIS
type ResCreateQRISXnd struct {
	ID          string `json:"id"` // Unique ID from xendit
	ExternalID  string `json:"external_id"`
	Amount      int64  `json:"amount"`
	QrString    string `json:"qr_string"`
	CallbackURL string `json:"callback_url"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created"`
	UpdatedAt   string `json:"updated"`
}

// ResCreateQRIS res create QRIS
type ResCreateQRIS struct {
	ID          string `json:"id"` // Unique ID from xendit
	ExternalID  string `json:"trxId"`
	Amount      int64  `json:"amount"`
	QrString    string `json:"qrCode"`
	CallbackURL string `json:"callbackUrl"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ReqXndQRIS callback from xendit for status payment
type ReqXndQRIS struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	QrString   string `json:"qr_string"`
	Type       string `json:"type"`
}

// ResQRIS ...
type ResQRIS struct {
	ID         string `json:"id"`
	ExternalID string `json:"trxId"`
	QrString   string `json:"qrString"`
	Type       string `json:"type"`
}

// ReqXndQRISCallback callback from xendit for status payment
type ReqXndQRISCallback struct {
	Event   string     `json:"event"`
	ID      string     `json:"id"`
	Amount  int64      `json:"amount"`
	Created string     `json:"created"`
	QRCode  ReqXndQRIS `json:"qr_code"`
	Status  string     `json:"status"`
}

// ResSimulateQRXnd response from xendit simulate qr (testmode only)
type ResSimulateQRXnd struct {
	ID      string     `json:"id"`
	Amount  int64      `json:"amount"`
	Created string     `json:"created"`
	QRCode  ReqXndQRIS `json:"qr_code"`
	Status  string     `json:"status"`
}

// ResSimulateQR respone xendit simulate qr (testmode only)
type ResSimulateQR struct {
	ID      string  `json:"id"`
	Amount  int64   `json:"amount"`
	Created string  `json:"created"`
	QRCode  ResQRIS `json:"qrCode"`
	Status  string  `json:"status"`
}
