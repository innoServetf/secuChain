package Controller

// UpdateDIDRequest 定义了更新 DID 时需要的请求参数
type UpdateDIDRequest struct {
	DID                string `json:"did" binding:"required"`
	RecoveryKey        string `json:"recovery_key" binding:"required"`
	RecoveryPrivateKey string `json:"recovery_private_key" binding:"required"`
}

// DeleteDIDRequest 定义了删除 DID 时需要的请求参数
type DeleteDIDRequest struct {
	DID                string `json:"did" binding:"required"`
	RecoveryKey        string `json:"recovery_key" binding:"required"`
	RecoveryPrivateKey string `json:"recovery_private_key" binding:"required"`
}
