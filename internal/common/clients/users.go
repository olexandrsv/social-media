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
	UsersFullNames(ids []int) ([]FullName, error)
}

type FullName struct {
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

func (c *usersClient) UsersFullNames(ids []int) ([]FullName, error) {
	convertedIDs := make([]int64, 0, len(ids))
	for _, id := range ids {
		convertedIDs = append(convertedIDs, int64(id))
	}
	resp, err := c.UsersClient.UsersFullNames(context.Background(), &users.UsersFullNamesReq{
		UserID: convertedIDs,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if resp.Err != nil {
		return nil, common.NewError(int(resp.Err.Code), resp.Err.Message)
	}

	fullNames := make([]FullName, 0, len(resp.FullName))
	for _, fullName := range resp.FullName {
		fullNames = append(fullNames, FullName{
			Name:    fullName.Name,
			Surname: fullName.Surname,
		})
	}
	return fullNames, nil
}
