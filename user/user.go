package user

import (
	"time"
	"context"
	"github.com/Sirupsen/logrus"
	"user-service/client"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"github.com/davecgh/go-spew/spew"
	"os"
)

const userIndex string = "users"
const userType string = "user"

var log = logrus.New()

type UserService struct {
	ctx           context.Context
	elasticClient *client.ElasticClient
}

func NewUserService(ctx context.Context) *UserService {
	elasticClient, err := client.NewElasticClient(ctx, userIndex)
	if err != nil {
		log.WithError(err).Error("Could not create elasticsearch client.")
		spew.Dump(err)
		os.Exit(1)
	}
	return &UserService{ctx, elasticClient}
}

// User represents a GitHub user.
type User struct {
	Id        *string    `json:"_id,omitempty"`
	Login     *string    `json:"login,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Email     *string    `json:"email,omitempty"`
	ApiToken  *string    `json:"api_token,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

func (s *UserService) FindByLogin(login string) (*User, error) {
	termQuery := elastic.NewTermQuery("login", login)
	searchResult, err := s.elasticClient.Client.Search().
		Index(userIndex).
		Type(userType).
		Query(termQuery).
		Pretty(true).
		Do(s.ctx)

	if err != nil {
		return nil, err
	}
	if searchResult.TotalHits() > 0 {
		user := &User{}
		err = json.Unmarshal([]byte(*searchResult.Hits.Hits[0].Source), user)
		user.Id = &searchResult.Hits.Hits[0].Id
		if err != nil {
			return nil, err
		}
		return user, err
	}

	return nil, err
}

func (s *UserService) SyncUser(login string, gitHubApiToken string) error {
	gitHubClient := client.CreateGitHubUserClient(s.ctx, gitHubApiToken)
	gitHubUser, _, err := gitHubClient.Users.Get(s.ctx, login)
	if err != nil {
		return err
	}
	gitStreamUser, err := s.FindByLogin(login)
	if err != nil {
		return err
	}
	if gitStreamUser == nil {
		// insert gitStreamUser
		gitStreamUser = &User{}
		gitStreamUser.Login = gitHubUser.Login
		gitStreamUser.Name = gitHubUser.Name
		gitStreamUser.Email = gitHubUser.Email
		gitStreamUser.ApiToken = nil
		gitStreamUser.CreatedAt = time.Now()
		gitStreamUser.UpdatedAt = time.Now()
		_, err = s.elasticClient.Client.Index().
			Index(userIndex).
			Type(userType).
			BodyJson(gitStreamUser).
			Do(s.ctx)
	} else {
		// update gitStreamUser
		gitStreamUser.Email = gitHubUser.Email
		gitStreamUser.Name = gitHubUser.Name
		gitStreamUser.UpdatedAt = time.Now()
		_, err = s.elasticClient.Client.Index().
			Index(userIndex).
			Type(userType).
			Id(*gitStreamUser.Id).
			BodyJson(gitStreamUser).
			Do(s.ctx)
	}


	return err
}
