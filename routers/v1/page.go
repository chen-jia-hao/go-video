package v1

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
)

func Live(c *gin.Context)  {
    ip, port, key, done := parsePath(c)
    if done {
        return
    }

    var url = "http://" + ip + ":" + port + "/live?app=live&stream=" + key

    c.HTML(http.StatusOK, "flv.html", gin.H{
        "url": url,
    })
}

func Hls(c *gin.Context)  {
    ip, port, key, done := parsePath(c)
    if done {
        return
    }

    var url = "http://" + ip + ":" + port + "/hls/" + key + ".m3u8"

    c.HTML(http.StatusOK, "hls.html", gin.H{
        "url": url,
    })
}

func Dash(c *gin.Context)  {
    ip, port, key, done := parsePath(c)
    if done {
        return
    }

    var url = "http://" + ip + ":" + port + "/dash/" + key + "/index.mpd"

    c.HTML(http.StatusOK, "dash.html", gin.H{
        "url": url,
    })
}

// 解析路径参数
func parsePath(c *gin.Context) (string, string, string, bool) {
    ip := c.Query("ip")
    if ip == "" {
        ip = strings.Split(c.Request.Host, ":")[0]
    }
    port := c.Query("port")
    if port == "" {
        port = "80"
    }
    key := c.Query("key")
    if key == "" {
        c.JSON(http.StatusOK, gin.H{
            "msg": "请求参数缺少key",
        })
        return "", "", "", true
    }
    return ip, port, key, false
}
