package Controller

import (
	"awesomeProject/DID/Service"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/gin-gonic/gin"
)

type DIDHandler struct {
	didService *Service.DIDService
}

func NewDIDHandler(didService *Service.DIDService) *DIDHandler {
	return &DIDHandler{
		didService: didService,
	}
}

// RegisterDID 处理创建 DID 的 API 请求
func (h *DIDHandler) RegisterDID(ctx *app.RequestContext) {
	// 调用服务层的 RegisterDID 方法
	didData, err := h.didService.RegisterDID()
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("创建 DID 失败: %v", err),
		})
		return
	}

	// 成功响应
	ctx.JSON(consts.StatusCreated, gin.H{
		"message": "DID 创建成功",
		"did":     didData,
	})
}

// UpdateDID 用于更新 DID 的 API 方法
func (h *DIDHandler) UpdateDID(ctx *app.RequestContext) {
	var req UpdateDIDRequest // 假设你有一个更新DID请求的结构体
	// 绑定并验证请求体中的数据
	if err := ctx.BindAndValidate(&req); err != nil {
		// 请求参数无效，返回 400 错误
		ctx.JSON(consts.StatusBadRequest, gin.H{
			"error":   "无效的请求参数",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层的 UpdateDID 方法
	result, err := h.didService.UpdateDID(req.DID, req.RecoveryKey, req.RecoveryPrivateKey)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("更新 DID 失败: %v", err),
		})
		return
	}

	// 成功响应
	ctx.JSON(consts.StatusOK, gin.H{
		"message": "DID 更新成功",
		"result":  result,
	})
}

// DeleteDID 用于删除 DID 的 API 方法
func (h *DIDHandler) DeleteDID(ctx *app.RequestContext) {
	// 定义请求体结构体
	var req DeleteDIDRequest
	// 绑定并验证请求体中的数据
	if err := ctx.BindAndValidate(&req); err != nil {
		// 请求参数无效，返回 400 错误
		ctx.JSON(consts.StatusBadRequest, gin.H{
			"error":   "无效的请求参数",
			"details": err.Error(),
		})
		return
	}

	// 调用服务层的 DeleteDID 方法
	result, err := h.didService.DeleteDID(req.DID, req.RecoveryKey, req.RecoveryPrivateKey)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("删除 DID 失败: %v", err),
		})
		return
	}

	// 成功响应
	ctx.JSON(consts.StatusOK, gin.H{
		"message": result, // 这里的 result 是成功返回的字符串，如 "註銷成功"
	})
}
