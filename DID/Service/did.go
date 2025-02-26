package Service

import (
	"awesomeProject/DID/Contract"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type DIDService struct {
	//contract *contractapi.Contract
	contract *gateway.Contract
	repo     *query.DIDRepository
}

func NewDIDService(contract *gateway.Contract, repo *query.DIDRepository) *DIDService {
	return &DIDService{
		contract: contract,
		repo:     repo,
	}
}

// RegisterDID 注册DID
func (s *DIDService) RegisterDID() (string, error) {
	var didData string
	// 创建DID
	didByte, err := s.contract.SubmitTransaction("RegisterDID")
	if err != nil {
		return "", fmt.Errorf("failed to register DID: %v", err)
	}
	didData = string(didByte)

	// 再写入数据库
	if err := s.repo.CreateDID(didData); err != nil {
		return "", fmt.Errorf("创建数据库 DID 失败: %v", err)
	}
	return didData, nil
}

// UpdateDID 更新 DID
func (s *DIDService) UpdateDID(did, recoveryKey, recoveryPrivateKey string) (string, error) {
	// 调用智能合约更新 DID
	resultBytes, err := s.contract.SubmitTransaction("UpdateDID", did, recoveryKey, recoveryPrivateKey)
	if err != nil {
		return "", fmt.Errorf("更新 DID 失败: %v", err)
	}

	// 转换结果为 string
	result := string(resultBytes)

	// 再更新数据库
	if err := s.repo.UpdateDID(did, result); err != nil {
		return "", fmt.Errorf("更新数据库 DID 失败: %v", err)
	}

	// 返回更新结果
	return string(result), nil
}

// DeleteDID 删除 DID
func (s *DIDService) DeleteDID(did, recoveryKey, recoveryPrivateKey string) (string, error) {
	// 调用智能合约撤销 DID
	_, err := s.contract.SubmitTransaction("Revoke", did, recoveryKey, recoveryPrivateKey)
	if err != nil {
		return "", fmt.Errorf("撤销 DID 失败: %v", err)
	}

	// 返回操作成功信息
	return "注销成功", nil
}
