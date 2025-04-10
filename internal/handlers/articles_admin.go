package handlers

import (
	"log"
	"my_website/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AdminIndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/index.html", map[string]interface{}{
		"Page": "index",
	})
}

func (h *Handler) AdminArticlesPage(c echo.Context) error {
	articles, err := h.Repository.Article.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Render(http.StatusOK, "admin/articles.html", map[string]interface{}{
		"Articles": articles,
		"Page":     "articles",
	})
}

func (h *Handler) AdminStatsPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/statistics.html", map[string]interface{}{
		"Page": "statistics",
	})
}

func (h *Handler) ArticleCreate(c echo.Context) error {
	article := &domain.Article{
		Title:       c.FormValue("title"),
		Description: c.FormValue("desc"),
		Public:      false,
	}
	err := h.Repository.Article.CreateOne(article)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleDelete(c echo.Context) error {
	images, err1 := h.Repository.Article.GetArticleImages(c.Param("article_id"))
	if err1 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	for _, image := range images {
		err := image.Remove()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
		}
	}
	err2 := h.Repository.Article.DeleteOneById(c.Param("article_id"))
	if err2 != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleUpdate(c echo.Context) error {
	article := &domain.Article{
		Title:       c.FormValue("title"),
		Description: c.FormValue("desc"),
		Content:     c.FormValue("content"),
	}
	if c.FormValue("public") == "on" {
		article.Public = true
	} else {
		article.Public = false
	}
	err := h.Repository.Article.UpdateOneById(article, c.Param("article_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.Redirect(http.StatusSeeOther, "/admin/articles")
}

func (h *Handler) ArticleEditPage(c echo.Context) error {
	article, err := h.Repository.Article.GetOneById(c.Param("article_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	images, err := h.Repository.Article.GetArticleImages(c.Param("article_id"))
	if err != nil {
		log.Fatal(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	article.Images = append(article.Images, images...)
	return c.Render(http.StatusOK, "admin/editor.html", map[string]interface{}{
		"Article": article,
		"Page":    "articles",
	})
}
