package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"main/adv_parse"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var suggestions = map[string][]string{
	"host":     {"localhost", "127.0.0.1", "192.168.1.1"},
	"port":     {"80", "443", "8080", "8443", "9200"},
	"protocol": {"tcp", "udp", "http", "https", "rpc"},
	"server":   {"nginx", "apache", "Apache-Tomcat"},
	"keyword":  {"工商银行", "招商银行", "建设银行", "交通银行"},
}

type PostParam struct {
	Input  string `json:"input"`
	Cursor int    `json:"cursor"`
}

type Suggestion struct {
	End     int    `json:"end"`
	Start   int    `json:"start"`
	Key     string `json:"key"`
	Query   string `json:"query"`
	Suggest string `json:"suggest"`
}

func ServeParse(ctx *gin.Context) {
	data := PostParam{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": fmt.Sprintf("%s", err)})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(400, map[string]string{"error": fmt.Sprintf("%s", err)})
		}
	}()
	lex := adv_parse.ParseStringLexer(data.Input)
	ctx.JSON(200, map[string]interface{}{
		"data": lex.Ast.ToDSL(),
	})
}

func ServeSuggest(ctx *gin.Context) {
	data := PostParam{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": fmt.Sprintf("%s", err)})
		return
	}
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(400, map[string]string{"error": fmt.Sprintf("%s", err)})
		}
	}()
	lex := adv_parse.ParseStringLexer(data.Input)
	rv := make([]Suggestion, 0)
	for _, suggest := range lex.Suggests {
		query := suggest.Value()
		suggestion := Suggestion{Query: query, Start: suggest.Start(), End: suggest.End()}
		for suggestKey, suggests := range suggestions {
			for _, suggestText := range suggests {
				if strings.Contains(suggestText, query) && (data.Cursor <= 0 || (suggestion.Start <= data.Cursor && data.Cursor <= suggestion.End)) {
					suggestion.Key = suggestKey
					suggestion.Suggest = fmt.Sprintf(`%s="%s"`, suggestKey, suggestText)
					rv = append(rv, suggestion)
				}
			}
		}

	}
	ctx.JSON(200, map[string]interface{}{
		"data": rv,
	})
}

func RunServer(port int) {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST(`/api/suggest`, ServeSuggest)
	router.POST(`/api/parse`, ServeParse)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("Close with error:", err)
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	sig := <-ch
	log.Println("Receive a signal", sig)

	cxt, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := srv.Shutdown(cxt)
	if err != nil {
		log.Println("Shutdown with error:", err)
		os.Exit(1)
	}
}

func main() {
	RunServer(5678)
}
