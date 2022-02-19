package grpcclient

import (
	"fmt"
	"sync"

	"github.com/toshkentov01/alif-tech-task/api-gateway/config"

	user "github.com/toshkentov01/alif-tech-task/api-gateway/genproto/user-service"
	"google.golang.org/grpc"
)

var cfg = config.Config()
var (
	onceUserService sync.Once

	instanceUserService user.UserServiceClient
)

// UserService ...
func UserService() user.UserServiceClient {
	onceUserService.Do(func() {
		connUser, err := grpc.Dial(
			fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
			grpc.WithInsecure())
		if err != nil {
			panic(fmt.Errorf("user service dial host: %s port:%d err: %s",
				cfg.UserServiceHost, cfg.UserServicePort, err))
		}

		instanceUserService = user.NewUserServiceClient(connUser)
	})

	return instanceUserService
}
