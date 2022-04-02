package main

import (
	"gindemo/models"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("starting cron....")
	//c := cron.New()
	c := cron.New()
	c.AddFunc("*/20 * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("*/20 * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.ClearAllArticle()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}

}
