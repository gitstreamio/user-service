package repository

import (
	"time"
	"context"
	"github.com/google/go-github/github"
	"github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
	"user-service/client"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"os"
)

const repositoryIndex string = "repositories"
const repositoryType string = "repository"

var log = logrus.New()

type RepositoryService struct {
	ctx           context.Context
	elasticClient *client.ElasticClient
}

func NewRepositoryService(ctx context.Context) *RepositoryService {
	elasticClient, err := client.NewElasticClient(ctx, repositoryIndex)
	if err != nil {
		log.WithError(err).Error("Could not create elasticsearch client.")
		spew.Dump(err)
		os.Exit(1)
	}

	return &RepositoryService{ctx, elasticClient}
}

// Repository represents a GitHub repository.
type Repository struct {
	Id           *string    `json:"id,omitempty"`
	Name         *string    `json:"name,omitempty"`
	FullName     *string    `json:"fullname,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Login        *string    `json:"login,omitempty"`
	Organization *string    `json:"organization,omitempty"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at,omitempty"`
}

func (s *RepositoryService) FindByLogin(login string) ([]*Repository, error) {
	termQuery := elastic.NewTermQuery("login", login)
	searchResult, err := s.elasticClient.Client.Search().
		Index(repositoryIndex).
		Type(repositoryType).
		Query(termQuery).
		Pretty(true).
		Size(200).
		Do(s.ctx)

	if err != nil {
		return nil, err
	}
	repositories := make([]*Repository, searchResult.Hits.TotalHits)
	for i, hit := range searchResult.Hits.Hits {
		repository := &Repository{}
		err = json.Unmarshal([]byte(*hit.Source), repository)
		repository.Id = &hit.Id
		repositories[i] = repository
	}

	return repositories, err
}

func (s *RepositoryService) SyncUserRepositories(login string, gitHubApiToken string) error {
	gitHubClient := client.CreateGitHubUserClient(s.ctx, gitHubApiToken)

	// get user repos from gitHub
	githubRepos, _, err := gitHubClient.Repositories.List(s.ctx, login, nil)
	if err != nil {
		return err
	}

	// get organization repos from gitHub
	userOrganizations, _, err := gitHubClient.Organizations.List(s.ctx, login, nil)
	if err != nil {
		return err
	}
	for _, organization := range userOrganizations {
		organizationRepos, _, err := gitHubClient.Repositories.ListByOrg(s.ctx, *organization.Login, nil)
		if err != nil {
			return err
		}
		githubRepos = append(githubRepos, organizationRepos...)
	}

	// get existing repos from elastic search
	gitstreamRepos, err := s.FindByLogin(login)
	if err != nil {
		return err
	}
	for _, githubRepo := range githubRepos {
		if repo := gitHubRepoExists(githubRepo, gitstreamRepos); repo == nil {
			err = s.indexUserRepo(login, githubRepo)
		} else {
			s.updateUserRepo(repo, githubRepo)
		}
	}

	// TODO remove deleted repositories from elastic index

	return err
}

func (s *RepositoryService) indexUserRepo(login string, githubRepo *github.Repository) error {
	gitstreamRepo := &Repository{}
	gitstreamRepo.Login = &login
	gitstreamRepo.Name = githubRepo.Name
	gitstreamRepo.FullName = githubRepo.FullName
	gitstreamRepo.Description = githubRepo.Description
	gitstreamRepo.CreatedAt = time.Now()
	gitstreamRepo.UpdatedAt = time.Now()
	_, err := s.elasticClient.Client.Index().
		Index(repositoryIndex).
		Type(repositoryType).
		BodyJson(gitstreamRepo).
		Do(s.ctx)

	return err
}

func (s *RepositoryService) updateUserRepo(gitstreamRepo *Repository, githubRepo *github.Repository) error {
	gitstreamRepo.Name = githubRepo.Name
	gitstreamRepo.Description = githubRepo.Description
	gitstreamRepo.UpdatedAt = time.Now()
	_, err := s.elasticClient.Client.Index().
		Index(repositoryIndex).
		Type(repositoryType).
		Id(*gitstreamRepo.Id).
		BodyJson(gitstreamRepo).
		Do(s.ctx)

	return err
}

func gitHubRepoExists(githubRepo *github.Repository, gitstreamRepos []*Repository) (*Repository) {
	for _, gitstreamRepo := range gitstreamRepos {
		if *gitstreamRepo.FullName == *githubRepo.FullName {
			return gitstreamRepo
		}
	}

	return nil
}
