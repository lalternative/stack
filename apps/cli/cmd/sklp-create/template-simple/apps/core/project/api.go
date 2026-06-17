package project

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"app/core/project/application/create_project"
	"app/core/project/application/list_projects"
)

func (s *Service) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	owner, _ := uuid.Parse(c.Request().Header.Get("X-User-Id"))
	p, err := s.create.Handle(c.Request().Context(), create_project.Command{Name: req.Name, OwnerID: owner})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, toDTO(p))
}

func (s *Service) List(c echo.Context) error {
	owner, _ := uuid.Parse(c.Request().Header.Get("X-User-Id"))
	ps, err := s.list.Handle(c.Request().Context(), list_projects.Query{OwnerID: owner})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	out := make([]ProjectDTO, 0, len(ps))
	for _, p := range ps {
		out = append(out, toDTO(p))
	}
	return c.JSON(http.StatusOK, out)
}
