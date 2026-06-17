package project

import (
	"github.com/labstack/echo/v4"

	"app/core/pkg/db"
	"app/core/project/application/create_project"
	"app/core/project/application/list_projects"
	"app/core/project/infrastructure"
)

type Service struct {
	create *create_project.Handler
	list   *list_projects.Handler
}

func NewService(exec db.Executor) *Service {
	repo := infrastructure.NewDuckDBRepository(exec)
	return &Service{
		create: create_project.NewHandler(repo),
		list:   list_projects.NewHandler(repo),
	}
}

func (s *Service) RegisterRoutes(g *echo.Group) {
	g.POST("/projects", s.Create)
	g.GET("/projects", s.List)
}
