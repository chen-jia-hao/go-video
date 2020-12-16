package main

import (
	"flag"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"video/asset"
	v1 "video/routers/v1"
)

var port = flag.Int("port", 3030, "端口号")

func main() {
	flag.Parse()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//加载模板文件
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}

	r.SetHTMLTemplate(t)

	//加载静态资源，例如网页的css、js
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: nil,
		Prefix:    "static", //一定要加前缀
	}
	r.StaticFS("/static", &fs)

	//r.LoadHTMLGlob("template/*")
	//r.Static("/static", "static")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1group := r.Group("/v1")
	{
		v1group.GET("/live", v1.Live)
		v1group.GET("/hls", v1.Hls)
		v1group.GET("/dash", v1.Dash)
	}

	r.GET("/live", func(c *gin.Context) {
		url := c.Query("url")
		timeout := c.Query("timeout")
		c.HTML(http.StatusOK, "flv.html", gin.H{
			"url":     url,
			"timeout": timeout,
		})
	})

	r.GET("/hls", func(c *gin.Context) {
		url := c.Query("url")
		timeout := c.Query("timeout")
		c.HTML(http.StatusOK, "hls.html", gin.H{
			"url":     url,
			"timeout": timeout,
		})
	})

	r.GET("/dash", func(c *gin.Context) {
		url := c.Query("url")
		timeout := c.Query("timeout")
		c.HTML(http.StatusOK, "dash.html", gin.H{
			"url":     url,
			"timeout": timeout,
		})
	})

	// 打印横幅
	file, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("\n" + string(file))
	}

	log.Println("init...")
	log.Println("sever listen port " + strconv.Itoa(*port))
	errRun := r.Run(":" + strconv.Itoa(*port))
	if errRun != nil {
		log.Fatal(errRun)
	}
}

//加载模板文件
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for _, name := range asset.AssetNames() {
		if !strings.HasSuffix(name, ".html") {
			continue
		}
		asset, err := asset.Asset(name)
		if err != nil {
			continue
		}
		name := strings.Replace(name, "template/", "", 1) //这里将templates替换下，在控制器中调用就不用加templates/了
		t, err = t.New(name).Parse(string(asset))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
