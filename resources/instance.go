package resources

import (
	"context"
	"fmt"
	"github.com/shybily/go-utils/pool"
	"math/rand"
	"net"
	"runtime"
	"time"
)

const (
	StrategyRandom = "random"
	StrategyOrder  = "order"
)

var insPool = &InstancePool{}

type InstancePool struct {
	strategy   string
	orderCount uint64
	instances  []*Instance
}

func (i *InstancePool) AddInstance(addr *net.TCPAddr) {
	i.instances = append(i.instances, NewInstance(addr))
}

func (i *InstancePool) random() *Instance {
	rand.Seed(time.Now().UnixNano())
	return i.instances[rand.Intn(len(i.instances))]
}

func (i *InstancePool) order() *Instance {
	defer func() { i.orderCount++ }()
	return i.instances[int(i.orderCount)%len(i.instances)]
}

type Instance struct {
	addr *net.TCPAddr
	pool *pool.ConnPool
}

func NewInstance(addr *net.TCPAddr) *Instance {
	return &Instance{
		addr: addr,
		pool: pool.NewConnPool(&pool.Options{
			Dialer: func(ctx context.Context) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, addr.Network(), addr.String())
			},
			PoolSize:           10 * runtime.NumCPU(),
			MinIdleConns:       10,
			MaxConnAge:         time.Hour,
			PoolTimeout:        2 * time.Second,
			IdleTimeout:        10 * time.Minute,
			IdleCheckFrequency: 0,
		}),
	}
}

func (i *Instance) GetConn(ctx context.Context) (*pool.Conn, error) {
	return i.pool.Get(ctx)
}

func (i *Instance) Put(con *pool.Conn) {
	i.pool.Put(con)
}

func InitInstance(address []string, strategy string) {
	insPool.strategy = strategy
	for _, v := range address {
		addr, _ := net.ResolveTCPAddr("tcp", v)
		insPool.AddInstance(addr)
	}
}

func GetInstance() (*Instance, error) {
	var ins *Instance
	if len(insPool.instances) <= 0 {
		return nil, fmt.Errorf("empty instances")
	}
	switch insPool.strategy {
	case StrategyRandom:
		ins = insPool.random()
	case StrategyOrder:
		ins = insPool.order()
	default:
		ins = insPool.random()
	}
	return ins, nil
}

func Range(f func(ins *Instance)) {
	for v := range insPool.instances {
		f(insPool.instances[v])
	}
}
