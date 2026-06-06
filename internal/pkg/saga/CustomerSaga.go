package saga

import (
	"context"
	"encoding/json"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/grpc/customer"
	"service/internal/pkg/saga/grpc"
	"time"
)

type CustomerSaga interface {
	FirstCustomerByUUID(request *customer.FirstCustomerRequest) map[string]interface{}

	Close()
}

func NewCustomerSaga() CustomerSaga {
	return &customerSaga{}
}

type customerSaga struct {
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

/** --- DEFER FUNCTION --- */

func (sg *customerSaga) Close() {
	if r := recover(); r != nil {
		panic(r)
	}
}
