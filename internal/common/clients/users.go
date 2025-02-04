package clients

import (
	"context"
	"social-media/api/pb/users"
	"social-media/internal/common"
	"social-media/internal/common/app/log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type UsersClient interface {
	UsersInfo(ids []int) ([]UsersInfo, error)
}

type UsersInfo struct {
	Login   string
	Name    string
	Surname string
}

type usersClient struct {
	users.UsersClient
}

func NewUsersClient() UsersClient {
	conn, err := grpc.Dial(":5053", grpc.WithInsecure())
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(err)
	}
	return &usersClient{
		users.NewUsersClient(conn),
	}
}

func (c *usersClient) UsersInfo(ids []int) ([]UsersInfo, error) {
	convertedIDs := make([]int64, 0, len(ids))
	for _, id := range ids {
		convertedIDs = append(convertedIDs, int64(id))
	}
	resp, err := c.UsersClient.UsersInfo(context.Background(), &users.UsersInfoReq{
		UserID: convertedIDs,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if resp.Err != nil {
		return nil, common.NewError(int(resp.Err.Code), resp.Err.Message)
	}

	infos := make([]UsersInfo, 0, len(resp.Info))
	for _, info := range resp.Info {
		infos = append(infos, UsersInfo{
			Login:   info.Login,
			Name:    info.Name,
			Surname: info.Surname,
		})
	}
	return infos, nil
}
