package web

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/volodymyrzuyev/goCsInspect/pkg/clientmanagement"
	"github.com/volodymyrzuyev/goCsInspect/pkg/inspect"
)

type Server struct {
	echo *echo.Echo
	cm   clientmanagement.ClientManager
	l    *slog.Logger
}

func NewServer(clientManager clientmanagement.ClientManager, l *slog.Logger) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	s := &Server{
		echo: e,
		cm:   clientManager,
		l:    l.WithGroup("Web"),
	}

	e.GET("/", s.root)

	return s
}

func (s *Server) Run(bindIp string) {
	s.l.Info("starting REST api", "address", bindIp)
	s.echo.Start(bindIp)
}

func (s *Server) root(c echo.Context) error {
	s.l.Debug("got request", "ip", c.RealIP(), "querry", c.Request().URL)
	startTime := time.Now()

	querry := c.Request().URL.Query()

	inspectLink := querry.Get("url")
	paramd := querry.Get("d")

	var params inspect.Parameters
	var err error

	if inspectLink != "" {
		params, _ = inspect.ParseInspectLink(inspectLink)
		err = params.Validate()
		if err != nil {
			slog.Error(
				"invalid request",
				"ip",
				c.RealIP(),
				"querry",
				c.Request().URL,
				"time_taken",
				time.Since(startTime),
			)
			return c.JSON(http.StatusBadRequest, "invalid params")
		}
	} else if paramd != "" {
		params.M, _ = strconv.ParseUint(querry.Get("m"), 10, 64)
		params.A, _ = strconv.ParseUint(querry.Get("a"), 10, 64)
		params.D, _ = strconv.ParseUint(querry.Get("d"), 10, 64)
		params.S, _ = strconv.ParseUint(querry.Get("s"), 10, 64)
		err = params.Validate()
		if err != nil {
			slog.Error(
				"invalid request",
				"ip",
				c.RealIP(),
				"querry",
				c.Request().URL,
				"time_taken",
				time.Since(startTime),
			)
			return c.JSON(http.StatusBadRequest, "invalid params")
		}
	} else {
		slog.Error(
			"invalid request",
			"ip",
			c.RealIP(),
			"querry",
			c.Request().URL,
			"time_taken",
			time.Since(startTime),
		)
		return c.JSON(http.StatusBadRequest, "invalid params")
	}

	item, err := s.cm.InspectSkin(params)
	if err != nil {
		slog.Error(
			"internal server error",
			"ip",
			c.RealIP(),
			"querry",
			c.Request().URL,
			"error",
			err,
			"time_taken",
			time.Since(startTime),
		)
		return c.JSON(http.StatusInternalServerError, "error getting data")
	}

	slog.Debug("successful reply",
		"ip",
		c.RealIP(),
		"query",
		c.Request().URL,
		"time_taken",
		time.Since(startTime),
	)
	return c.JSONPretty(http.StatusOK, item, "    ")
}
