package rabbitmq

import (
	"service/internal/other/service"
	"service/internal/pkg/core"
	form2 "service/internal/pkg/form"
	"sync"

	xtremerabbitmq "github.com/globalxtreme/go-core/v2/rabbitmq"
)

type CustomerUpdateExecutor struct {
	xtremerabbitmq.AsyncWorkflowConsumerBase

	mutex sync.Mutex
	form  form2.CustomerUpdateForm
}

func (c *CustomerUpdateExecutor) Consume(payload interface{}) (interface{}, error, []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return core.RabbitMQErrorHandler(func() (interface{}, error) {
		c.form = form2.CustomerUpdateForm{}
		err := c.form.AsyncWorkflowParse(payload)
		if err != nil {
			return nil, err
		}

		srv := service.NewCustomerService()
		srv.Update(c.form)

		return c.Response(payload), nil

	})
}

func (consume *CustomerUpdateExecutor) Response(payload interface{}, data ...interface{}) interface{} {
	return nil
}
