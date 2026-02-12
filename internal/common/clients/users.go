package clients

import (
	"context"
	"social-media/api/pb/users"
	"social-media/internal/common"
	"social-media/internal/common/app/config"
	"social-media/internal/common/app/log"
	"social-media/internal/common/slice"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type UsersClient interface {
	UsersInfo(ids []int) ([]UsersInfo, error)
	LastReadPosts(userID int) ([]Post, error)
	SendPostNotification(userID int, message string) error
	SendNotification(userID int, receiversIDs []int, message string) error
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
	host := config.App.Users.Service.Host
	port := config.App.Users.Service.GrpcPort
	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
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

type Post struct {
	ID          string
	FollowingID int
}

func (c *usersClient) LastReadPosts(userID int) ([]Post, error) {
	resp, err := c.UsersClient.LastReadPosts(context.Background(), &users.LastReadPostsReq{
		UserID: int64(userID),
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}
	if resp.Err != nil {
		err = common.NewError(int(resp.Err.Code), resp.Err.Message)
		log.Error(errors.WithStack(err))
		return nil, common.ErrInternal
	}

	convertedPosts := slice.MustConvert(resp.Posts, func(post *users.Post) Post {
		return Post{
			ID:          post.Id,
			FollowingID: int(post.UserID),
		}
	})
	return convertedPosts, nil
}

func (c *usersClient) SendPostNotification(userID int, message string) error {
	resp, err := c.UsersClient.SendPostNotification(context.Background(), &users.SendPostNotificationReq{
		SenderID: int64(userID),
		Message:  message,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	if resp.Err != nil {
		err = common.NewError(int(resp.Err.Code), resp.Err.Message)
		log.Error(errors.WithStack(err))
		return err
	}

	return nil
}

func (c *usersClient) SendNotification(userID int, receiversIDs []int, messgae string) error {
	resp, err := c.UsersClient.SendNotification(context.Background(), &users.SendNotificationReq{
		SenderID: int64(userID),
		ReceiversIDs: slice.MustConvert(receiversIDs, func(n int) int64 {
			return int64(n)
		}),
		Message: messgae,
	})
	if err != nil {
		log.Error(errors.WithStack(err))
		return common.ErrInternal
	}
	if resp.Err != nil {
		err = common.NewError(int(resp.Err.Code), resp.Err.Message)
		log.Error(errors.WithStack(err))
		return err
	}

	return nil
}
