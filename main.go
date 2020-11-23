package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

)

func main() {
	fmt.Println("欢迎进入豆瓣系统爬虫系统")

	client := http.Client{}

	request, err := http.NewRequest("GET", "https://movie.douban.com/chart", nil)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Accept-Language", "zh-CN")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Host", "img3.doubain.com")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	htmlBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	html:=string(htmlBytes)

	movieIdReg := regexp.MustCompile(`<a href="https://movie.douban.com/subject/(.*?)/"`)
	idSlice := movieIdReg.FindAllStringSubmatch(html, -1)
	nameReg := regexp.MustCompile(` width="75" alt="(.*?)"`)
	nameSlice := nameReg.FindAllStringSubmatch(html, -1)

	ratReg := regexp.MustCompile(`<span class="rating_nums">(.*?)</span>`)
	ratSlice := ratReg.FindAllStringSubmatch(html,-1)

	voteReg := regexp.MustCompile(`<span class="pl">(.*?)</span>`)
	voteSlice := voteReg.FindAllStringSubmatch(html,-1)


	/*coverImgReg := regexp.MustCompile(`src="(.*?)"`)
	imgSlice := coverImgReg.FindAllStringSubmatch(html,-1)

	/*fmt.Println("电影编号  电影名称 评分 评价人数 一句话总结 封面图")
	fmt.Printf("%s %s %s %s %s %s \n",idSlice,nameSlice,ratSlice,voteSlice,descSlice,imgSlice)*/
	fmt.Println("电影编号 电影名字 评分 评价人数   ")
	for i := 0; i < len(nameSlice); i++ {
		fmt.Printf("%s  %s %s %s  \n",
			idSlice[i][1],
			nameSlice[i][1],
			ratSlice[i][1],
			voteSlice[i][1])
	}

	database,err:=sql.Open("mysql","root:409216@tcp(127.0.0.1:3306)/xiaoliu?charset=utf8")
		defer database.Close()
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		for j := 0; j < len(idSlice); j++ {
			_,err:=database.Exec("insert into" +
				"Reptile(id,name,score,Number_of_evaluators)" +
				"values(?,?,?,?)",
				idSlice[j][1],
				nameSlice[j][1],
				ratSlice[j][1],
				voteSlice[j][1])
			if err != nil {
				log.Fatal(err.Error())
				break
			}
			/*fmt.Println()*/
		}
		fmt.Println("豆瓣电影信息采集完成")
}