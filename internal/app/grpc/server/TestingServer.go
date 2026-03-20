package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"service/internal/pkg/activity"
	"service/internal/pkg/config"
	"service/internal/pkg/core"
	"service/internal/pkg/form"
	"service/internal/pkg/grpc/example"
	"service/internal/testing/repository"
)

type TestingServer struct {
	example.UnimplementedTestingServiceServer

	rollbackData map[string]interface{}
}

func (srv *TestingServer) Register(serverRPC *grpc.Server) {
	example.RegisterTestingServiceServer(serverRPC, srv)
}

func (srv *TestingServer) Store(ctx context.Context, in *example.TestingRequest) (*example.EXResponse, error) {
	res, err := core.GRPCErrorHandler(func() (*example.EXResponse, error) {
		err := config.PgSQL.Transaction(func(tx *gorm.DB) error {
			repo := repository.NewTestingRepository(tx)
			testing := repo.Create(form.TestingForm{Name: in.GetName()})

			subs := in.GetSubs()
			if subs != nil && len(subs) > 0 {
				for _, sub := range subs {
					if sub == nil {
						continue
					}

					repo.AddSub(testing, sub.GetName())
				}
			}

			activity.UseActivity{}.SetReference(testing).
				Save(fmt.Sprintf("gRPC: Enter new testing: %s [%s]", testing.Name, testing.ID))

			srv.rollbackData = map[string]interface{}{
				"id": testing.ID,
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return srv.success(srv.rollbackData)
	})

	return res, err
}

func (srv *TestingServer) RollbackStore(ctx context.Context, in *example.RollBackRequest) (*example.EXResponse, error) {
	res, err := core.GRPCErrorHandler(func() (*example.EXResponse, error) {
		err := json.Unmarshal(in.GetData(), &srv.rollbackData)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to unmarshal data! Err: %v", err))
		}

		err = config.PgSQL.Transaction(func(tx *gorm.DB) error {
			repo := repository.NewTestingRepository(tx)

			testing := repo.FirstByForm(form.TestingFilterForm{
				ID: uint(xtremepkg.ToInt(srv.rollbackData["id"].(string))),
			}, func(query *gorm.DB) *gorm.DB {
				return query.Preload("Subs")
			})

			repo.Delete(testing)

			if subs := testing.Subs; len(subs) > 0 {
				for _, sub := range subs {
					repo.DeleteSub(sub)
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		return srv.success()
	})

	return res, err
}

/** --- UNEXPORTED FUNCTIONS --- */

func (srv *TestingServer) success(result ...any) (*example.EXResponse, error) {
	var data []byte
	if len(result) > 0 {
		data, _ = json.Marshal(result[0])
	}

	return &example.EXResponse{Message: "Success", Result: data}, nil
}
