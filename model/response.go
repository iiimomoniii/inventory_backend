package model

// Response — common response สำหรับทุก endpoint
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrorData — error response มาตรฐาน แสดงทั้ง 2 ภาษา
type ErrorData struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	ErrorCode string `json:"errorCode"`
	Error     string `json:"error"`
	MessageEN string `json:"messageEN"`
	MessageTH string `json:"messageTH"`
	Path      string `json:"path"`
}

// ─── Items Response ─────────────────────────────────────────

// ItemResult — ผลของแต่ละ item ใน Items operation
type ItemResult struct {
	Index     int    `json:"index"`               // ลำดับใน request
	Status    string `json:"status"`              // "success" หรือ "error"
	ErrorCode string `json:"errorCode,omitempty"` // "PRD001"
	Error     string `json:"error,omitempty"`     // "NAME_REQUIRED"
	MessageEN string `json:"messageEN,omitempty"` // "Name is required"
	MessageTH string `json:"messageTH,omitempty"` // "กรุณากรอกชื่อสินค้า"
	Data      any    `json:"data,omitempty"`      // product ที่สำเร็จ
}

// ItemsErrorData — response สำหรับ Items operation
type ItemsErrorData struct {
	Timestamp    string       `json:"timestamp"`
	Status       int          `json:"status"`
	TotalSuccess int          `json:"totalSuccess"`
	TotalError   int          `json:"totalError"`
	Results      []ItemResult `json:"results"`
	Path         string       `json:"path"`
}
