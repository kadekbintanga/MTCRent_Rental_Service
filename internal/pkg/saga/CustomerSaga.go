package saga

import (
	"context"
	"encoding/json"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/form/option"
	"service/internal/pkg/grpc/customer"
	"service/internal/pkg/saga/grpc"
	"time"
)

type CustomerSaga interface {
	FirstCustomerByUUID(request *customer.FirstCustomerRequest) map[string]interface{}
	UpdateCustomerStatus(request *customer.CustomerUpdateStatusRequest) option.CustomerUpdateStatusRollbackOption

	Close()
}

func NewCustomerSaga() CustomerSaga {
	return &customerSaga{}
}

type customerSaga struct {
	rollbackUpdateStatus *option.CustomerUpdateStatusRollbackOption
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

func (sg *customerSaga) UpdateCustomerStatus(request *customer.CustomerUpdateStatusRequest) option.CustomerUpdateStatusRollbackOption {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	resp, err := grpc.CustomerRPCCClient.UpdateStatus(ctx, request)
	if err != nil {
		error2.ErrXtremeCustomerUpdate(err.Error())
	}

	var rollbackData option.CustomerUpdateStatusRollbackOption
	if result := resp.GetResult(); len(result) > 0 {
		json.Unmarshal(result, &rollbackData)
	}
	rollbackData.CreatedBy = request.CreatedBy
	rollbackData.CreatedByName = request.CreatedByName

	sg.rollbackUpdateStatus = &rollbackData

	return rollbackData
}

/** --- DEFER FUNCTION --- */

func (sg *customerSaga) Close() {
	if r := recover(); r != nil {
		if sg.rollbackUpdateStatus != nil {
			sg.UpdateCustomerStatus(&customer.CustomerUpdateStatusRequest{
				Uuid:            sg.rollbackUpdateStatus.UUID,
				StatusId:        sg.rollbackUpdateStatus.StatusId,
				BlacklistReason: sg.rollbackUpdateStatus.BlacklistReason,
				CreatedBy:       sg.rollbackUpdateStatus.CreatedBy,
				CreatedByName:   sg.rollbackUpdateStatus.CreatedByName,
			})
		}
		panic(r)
	}
}
