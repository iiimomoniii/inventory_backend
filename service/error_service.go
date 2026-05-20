package service

// ServiceError — ใช้ร่วมกันทุก service
// handler ดึง Code ไปใช้กับ utils.CustomErrorResp
type ServiceError struct {
	Code string // "PRD001", "GLB002" ...
}

func (e *ServiceError) Error() string { return e.Code }
