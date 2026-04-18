package number

import (
	"service/internal/pkg/core"
	"service/internal/pkg/model"
)

// TODO: Hapus kalau tidak sudah tidak digunakan

type TestingModelNumber struct{}

func (n *TestingModelNumber) Prefix() string {
	return "TST"
}

func (n *TestingModelNumber) Generate() string {
	return core.ModelNumber(model.Testing{}, n)
}

type TestingRedisNumber struct{}

func (n *TestingRedisNumber) Test1() string {
	return n.generate("TST1")
}

func (n *TestingRedisNumber) Test2() string {
	return n.generate("TST2")
}

func (n *TestingRedisNumber) generate(prefix string) string {
	return core.RedisNumber(prefix)
}
