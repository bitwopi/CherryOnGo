package client

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	remapi "github.com/Jolymmiles/remnawave-api-go/v2/api"
	"github.com/google/uuid"
)

type Client struct {
	api *remapi.ClientExt
}

type RemnaPlan struct {
	DayLimit    int
	DeviceLimit int
	Squad       uuid.UUID
}

// NewClient creates a new Remnawave API client wrapper.
func NewClient(apiKey string, baseURL string) *Client {
	baseClient, err := remapi.NewClient(
		baseURL,
		remapi.StaticToken{Token: apiKey},
	)
	if err != nil {
		panic(err)
	}
	apiClient := remapi.NewClientExt(baseClient)
	return &Client{api: apiClient}
}

func (c *Client) GetUserByUsername(ctx context.Context, username string) (*remapi.UserResponse, error) {
	resp, err := c.api.Users().GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	result, ok := resp.(*remapi.UserResponse)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return result, nil
}

func (c *Client) CreateUser(ctx context.Context, plan *RemnaPlan, username string, tgID string, email string) (*remapi.UserResponse, error) {
	if len(username) == 0 {
		username = uuid.New().String()
	}
	userDto := remapi.CreateUserRequestDto{
		Username:             username,
		CreatedAt:            remapi.NewOptDateTime(time.Now()),
		ExpireAt:             time.Now().AddDate(0, 0, plan.DayLimit),
		HwidDeviceLimit:      remapi.NewOptInt(plan.DeviceLimit),
		ActiveInternalSquads: []uuid.UUID{plan.Squad},
	}
	if len(tgID) != 0 {
		val, err := strconv.Atoi(tgID)
		if err != nil {
			return nil, err
		}
		userDto.TelegramId = remapi.NewOptNilInt(val)
	}
	if len(email) != 0 {
		userDto.Email = remapi.NewOptNilString(email)
	}
	if len(plan.Squad) == 36 {
		userDto.ActiveInternalSquads = []uuid.UUID{plan.Squad}
	}
	resp, err := c.api.Users().CreateUser(ctx, &userDto)
	if err != nil {
		return nil, err
	}

	res, ok := resp.(*remapi.UserResponse)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return res, nil
}

func (c *Client) UpdateUserExpiryTime(ctx context.Context, plan *RemnaPlan, username *string, uuid *uuid.UUID) (*remapi.UserResponse, error) {
	userDto := remapi.UpdateUserRequestDto{
		ExpireAt:        remapi.NewOptDateTime(time.Now().AddDate(0, 0, plan.DayLimit)),
		HwidDeviceLimit: remapi.NewOptNilInt(plan.DeviceLimit),
	}
	if username != nil {
		userDto.Username = remapi.NewOptString(*username)
	}
	if uuid != nil {
		userDto.UUID = remapi.NewOptUUID(*uuid)
	}
	resp, err := c.api.Users().UpdateUser(ctx, &userDto)
	if err != nil {
		return nil, err
	}

	res, ok := resp.(*remapi.UserResponse)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return res, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.api.Users().GetAllUsers(ctx, 0, 0)
	return err
}

func (c *Client) GetUsersByTgID(ctx context.Context, tgID string) (*remapi.UsersResponse, error) {
	resp, err := c.api.Users().GetUserByTelegramId(ctx, tgID)
	log.Println("client resp: ", resp)
	if err != nil {
		return nil, err
	}

	res, ok := resp.(*remapi.UsersResponse)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return res, nil
}

func (c *Client) GetUsersByEmail(ctx context.Context, email string) (*remapi.UsersResponse, error) {
	resp, err := c.api.Users().GetUsersByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	res, ok := resp.(*remapi.UsersResponse)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return res, nil
}

func (c *Client) GetAllUsers(ctx context.Context) (*remapi.GetAllUsersResponseDto, error) {
	resp, err := c.api.Users().GetAllUsers(ctx, 0, 0)
	if err != nil {
		return nil, err
	}
	res, ok := resp.(*remapi.GetAllUsersResponseDto)
	if !ok {
		return nil, errors.New("undefined response")
	}
	return res, nil
}
