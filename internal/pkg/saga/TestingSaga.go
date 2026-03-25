package saga

import (
	"context"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	error2 "service/internal/pkg/error"
	"service/internal/pkg/grpc/example"
	"service/internal/pkg/saga/grpc"
	"service/internal/pkg/saga/privateapi"
	"time"
)

// TODO: Hanya contoh. nanti langsung hapus saja
type TestingSaga interface {
	TestingStore(request *example.TestingRequest) (string, []byte)
	TestingStoreAPI(request *example.TestingRequest) interface{}

	Close()
	TestingRollbackStore()
	TestingAPIRollbackStore()
}

func NewTestingSaga() TestingSaga {
	return &testingSaga{
		testingAPI: privateapi.NewTestingAPI(),
	}
}

type testingSaga struct {
	testingAPI privateapi.TestingAPI

	testingRPCRollBack []byte
	testingAPIRollBack interface{}
}

/** --- ITEM SERVICE CLIENT --- */

func (sg *testingSaga) TestingStore(request *example.TestingRequest) (string, []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpc.TestingRPCClient.Store(ctx, request)
	if err != nil {
		error2.ErrXtremeTestingSave(err.Error())
	}

	result := resp.GetResult()

	sg.testingRPCRollBack = result

	return resp.GetMessage(), result
}

func (sg *testingSaga) TestingStoreAPI(request *example.TestingRequest) interface{} {
	resp := sg.testingAPI.Store(request)
	result := resp.Result
	if result != nil {
		sg.testingAPIRollBack = result
	}

	return result
}

func (sg *testingSaga) TestingRollbackStore() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := grpc.TestingRPCClient.RollbackStore(ctx, &example.RollBackRequest{Data: sg.testingRPCRollBack})
	if err != nil {
		xtremepkg.LogError(err, true)
	}
}

func (sg *testingSaga) TestingAPIRollbackStore() {
	sg.testingAPI.RollBack(sg.testingAPIRollBack)
}

/** --- DEFER FUNCTION --- */

func (sg *testingSaga) Close() {
	if r := recover(); r != nil {
		if len(sg.testingRPCRollBack) > 0 {
			sg.TestingRollbackStore()
		}

		if sg.testingAPIRollBack != nil {
			sg.TestingAPIRollbackStore()
		}

		panic(r)
	}
}
