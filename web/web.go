package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"goupsc/nut"

	"github.com/labstack/echo/v4"
)

//go:embed dist/*
var embededFiles embed.FS

func GetFileSystem() http.Handler {
	fsys, _ := fs.Sub(embededFiles, "dist")
	return http.FileServer(http.FS(fsys))
}

func UpsInfo(c echo.Context) error {
	client := &nut.NUTClient{}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	ups := os.Getenv("UPS")
	err := client.Connect(host, port, ups)
	if err != nil {
		return err
	}
	defer client.Close()
	vars, err := client.ListVars()
	if err != nil {
		return err
	}
	result := make([]nut.VarValue, 0)
	for _, v := range vars {
		vv, _ := client.ParseVarValue(v)
		result = append(result, vv)
	}
	return c.JSON(200, result)
}
func StartWebServer() {
	e := echo.New()
	e.Server.Addr = os.Getenv("LISTEN")
	e.HideBanner = true
	e.GET("/api/upsinfo", UpsInfo)
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", GetFileSystem())))

	e.Logger.Debug(e.StartServer(e.Server))
}
