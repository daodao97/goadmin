package sso

import "context"

type UserInfo struct {
	Name          string
	Nickname      string
	Email         string
	Avatar        string
	DefaultRoleId int
	Extra         interface{}
}

type SSO interface {
	GetUserInfo(ctx context.Context, ticket string) (*UserInfo, error)
}

type Name = string

type Sso = map[Name]SSO
