package sso

import (
	"context"
	"fmt"
	"time"

	"github.com/carlmjohnson/requests"
)

const SsoGithub Name = "github"

type GithubOption struct {
	ClientId      string
	ClientSecret  string
	DefaultRoleId int
}

func NewGithub(option *GithubOption) SSO {
	return &Github{
		ClientId:      option.ClientId,
		ClientSecret:  option.ClientSecret,
		DefaultRoleId: option.DefaultRoleId,
	}
}

type Github struct {
	ClientId      string
	ClientSecret  string
	DefaultRoleId int
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUser struct {
	Login             string      `json:"login"`
	Id                int         `json:"id"`
	NodeId            string      `json:"node_id"`
	AvatarUrl         string      `json:"avatar_url"`
	GravatarId        string      `json:"gravatar_id"`
	Url               string      `json:"url"`
	HtmlUrl           string      `json:"html_url"`
	FollowersUrl      string      `json:"followers_url"`
	FollowingUrl      string      `json:"following_url"`
	GistsUrl          string      `json:"gists_url"`
	StarredUrl        string      `json:"starred_url"`
	SubscriptionsUrl  string      `json:"subscriptions_url"`
	OrganizationsUrl  string      `json:"organizations_url"`
	ReposUrl          string      `json:"repos_url"`
	EventsUrl         string      `json:"events_url"`
	ReceivedEventsUrl string      `json:"received_events_url"`
	Type              string      `json:"type"`
	SiteAdmin         bool        `json:"site_admin"`
	Name              string      `json:"name"`
	Company           interface{} `json:"company"`
	Blog              string      `json:"blog"`
	Location          string      `json:"location"`
	Email             string      `json:"email"`
	Hireable          interface{} `json:"hireable"`
	Bio               string      `json:"bio"`
	TwitterUsername   interface{} `json:"twitter_username"`
	PublicRepos       int         `json:"public_repos"`
	PublicGists       int         `json:"public_gists"`
	Followers         int         `json:"followers"`
	Following         int         `json:"following"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

func (g Github) GetUserInfo(ctx context.Context, code string) (*UserInfo, error) {
	var resp AccessTokenResp
	err := requests.
		URL("https://github.com/login/oauth/access_token").
		Param("client_id", g.ClientId).
		Param("client_secret", g.ClientSecret).
		Param("code", code).
		Header("Accept", "application/json").
		ToJSON(&resp).
		Post().
		Fetch(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("token %+v", resp))

	var user GithubUser
	err = requests.
		URL("https://api.github.com/user").
		Header("Accept", "application/json").
		Header("Authorization", fmt.Sprintf("Bearer %s", resp.AccessToken)).
		ToJSON(&user).
		Fetch(ctx)

	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Name:          user.Login,
		Nickname:      user.Name,
		Avatar:        user.AvatarUrl,
		Email:         user.Email,
		DefaultRoleId: g.DefaultRoleId,
	}, nil
}
