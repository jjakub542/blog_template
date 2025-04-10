package server

import (
	"my_website/internal/handlers"
	"my_website/internal/session"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Router() http.Handler {
	e := echo.New()
	e.Use(session.Middleware(s.store))
	e.Renderer = Renderer()
	e.Static("/static", "web/static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := handlers.Handler{Repository: s.repository}

	e.GET("/", h.HomePage)

	e.GET("/blog", h.BlogPage)
	e.GET("/blog/:article_id", h.ArticleView)

	adminGroup := e.Group("/admin")
	adminGroup.GET("", session.AdminAuth(h.AdminIndexPage))
	adminGroup.GET("/statistics", session.AdminAuth(h.AdminStatsPage))
	adminGroup.GET("/articles", session.AdminAuth(h.AdminArticlesPage))
	adminGroup.POST("/articles/create", session.AdminAuth(h.ArticleCreate))
	adminGroup.POST("/articles/:article_id/delete", session.AdminAuth(h.ArticleDelete))
	adminGroup.POST("/articles/:article_id/update", session.AdminAuth(h.ArticleUpdate))
	adminGroup.GET("/articles/:article_id", session.AdminAuth(h.ArticleEditPage))
	adminGroup.Any("/login", h.LoginPage)
	adminGroup.Any("/logout", h.LogoutPage)

	return e
}
