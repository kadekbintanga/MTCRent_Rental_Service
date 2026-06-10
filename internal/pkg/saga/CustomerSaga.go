package saga

import (
	"context"
	"encoding/json"
	"fmt"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/grpc/customer"
	"service/internal/pkg/saga/grpc"
	"time"
)

type CustomerSaga interface {
	FirstCustomerByUUID(request *customer.FirstCustomerRequest) map[string]interface{}
	UpdateCustomerStatus(request *customer.CustomerUpdateStatusRequest) map[string]interface{}

	Close()
}

func NewCustomerSaga() CustomerSaga {
	return &customerSaga{}
}

type customerSaga struct {
	rollbackData map[string]interface{}
}

func (sg *customerSaga) FirstCustomerByUUID(request *customer.FirstCustomerRequest) map[string]interface{} {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	resp, err := grpc.CustomerRPCCClient.FirstByUUID(ctx, request)
	if err != nil {
		error2.ErrXtremeCustomerGet(err.Error())
	}

	var customer map[string]interface{}
	if result := resp.GetResult(); len(result) > 0 {
		json.Unmarshal(result, &customer)
	}

	return customer
}

func (sg *customerSaga) UpdateCustomerStatus(request *customer.CustomerUpdateStatusRequest) map[string]interface{} {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	resp, err := grpc.CustomerRPCCClient.UpdateStatus(ctx, request)
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
	}

	var rollbackData map[string]interface{}
	if result := resp.GetResult(); len(result) > 0 {
		json.Unmarshal(result, &rollbackData)
	}
	rollbackData["createdBy"] = request.CreatedBy
	rollbackData["createdByName"] = request.CreatedByName

	sg.rollbackData = rollbackData

	return rollbackData
}

/** --- DEFER FUNCTION --- */

func (sg *customerSaga) Close() {
	if r := recover(); r != nil {
		if sg.rollbackData != nil {
			fmt.Println("RUN THIS =========================================================")
			sg.UpdateCustomerStatus(&customer.CustomerUpdateStatusRequest{
				Uuid:            sg.rollbackData["uuid"].(string),
				StatusId:        int32(sg.rollbackData["statusId"].(float64)),
				BlacklistReason: sg.rollbackData["blacklistReason"].(string),
			})
		}
		panic(r)
	}
}
