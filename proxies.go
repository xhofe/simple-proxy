package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var HttpClient = &http.Client{}

func proxiesHandle(app *gin.Engine) {
	for k, v := range config.Proxies {
		app.Any(k, proxy(k, v))
	}
}

func proxy(source string, target string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		target := target
		if strings.Contains(source, "*") {
			splat := ctx.Param("splat")
			splat = strings.TrimPrefix(splat, "/")
			target = strings.Replace(target, ":splat", splat, 1)
		} else {
			for _, param := range ctx.Params {
				target = strings.Replace(target, ":"+param.Key, param.Value, 1)
			}
		}
		if ctx.Request.URL.RawQuery != "" {
			target += "?" + ctx.Request.URL.RawQuery
		}
		if !strings.HasPrefix(target, "http") {
			target = "https://" + target
		} else {
			if !strings.Contains(target, "://") {
				target = strings.Replace(target, ":/", "://", 1)
			}
		}
		req, err := http.NewRequest(ctx.Request.Method, target, ctx.Request.Body)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		for h, val := range ctx.Request.Header {
			req.Header[h] = val
		}
		res, err := HttpClient.Do(req)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer func() {
			_ = res.Body.Close()
		}()
		for h, v := range res.Header {
			for _, s := range v {
				ctx.Header(h, s)
			}
		}
		ctx.Status(res.StatusCode)
		_, err = io.Copy(ctx.Writer, res.Body)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
	}
}
