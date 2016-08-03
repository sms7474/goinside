package goinside

import (
	"fmt"
	"testing"
)

func TestGetList(t *testing.T) {
	list, _ := GetList("http://gall.dcinside.com/board/lists/?id=baseball_new4", 1)
	for _, v := range list.Articles {
		article, _ := GetArticle(v.URL)
		fmt.Println(article.Content)
		fmt.Println(article.Images)
		fmt.Println("-------------------------------------------------")
	}
}

func TestGetArticle(t *testing.T) {
	article, _ := GetArticle("http://gall.dcinside.com/board/view/?id=baseball_new4&no=9146307")
	fmt.Println(article.Content)
	fmt.Println(article.Images)
}
