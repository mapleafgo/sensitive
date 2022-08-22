package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mapleafgo.cn/sensitive"
)

var filter *sensitive.Filter

func init() {
	filter = sensitive.New()
}

// Start 启动服务
func Start() error {
	// 初始化词库
	path := viper.GetString("path")
	if path != "" {
		var err error
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			err = filter.LoadNetWordDict(path)
		} else {
			err = filter.LoadWordDict(path)
		}
		if err != nil {
			return err
		}
	}

	// 初始化服务
	gin.SetMode(viper.GetString("mode"))
	r := gin.Default()
	r.Use(CheckAndPrint())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/add-word", addWord)
	r.POST("/remove-word", removeWord)
	r.POST("/filter", filterAll)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", viper.GetInt("port")),
		Handler: r,
	}

	go func() {
		// 启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server is running on %v\n", srv.Addr)

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	sc := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-sc
	log.Printf("Receive signal[%s]\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	return nil
}

type wordReq struct {
	Word string `json:"word"`
}

// addWord 添加敏感词
func addWord(c *gin.Context) {
	var req wordReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filter.AddWord(req.Word)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// removeWord 移除敏感词
func removeWord(c *gin.Context) {
	var req wordReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filter.RemoveWord(req.Word)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

type filterRequest struct {
	Content string `json:"content"`
}

// filterAll 过滤敏感词
func filterAll(c *gin.Context) {
	var req filterRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	data := filter.FindAll(req.Content)
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    data,
	})
}
