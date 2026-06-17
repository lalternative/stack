package project

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"app/core/project/application/create_project"
	"app/core/project/application/list_projects"
)

// Create godoc
// @Summary  Create a project
// @Tags     projects
// @Accept   json
// @Produce  json
// @Param    body  body      CreateRequest  true  "Project to create"
// @Success  201   {object}  ProjectDTO
// @Failure  400   {object}  echo.HTTPError
// @Security BearerAuth
// @Router   /projects [post]
// @ID       createProject
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

// List godoc
// @Summary  List projects owned by the caller
// @Tags     projects
// @Produce  json
// @Success  200  {array}   ProjectDTO
// @Failure  500  {object}  echo.HTTPError
// @Security BearerAuth
// @Router   /projects [get]
// @ID       listProjects
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
