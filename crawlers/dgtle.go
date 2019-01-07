package crawlers

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "io"
    "net/http"
)

func GetTrade(){
    url1 := "http://trade.dgtle.com/dgtle_module.php?mod=trade&ac=index&typeid=&PName=&searchsort=0&page=1"
    //url2 := "http://trade.dgtle.com/dgtle_module.php?mod=trade&ac=index&typeid=&PName=&searchsort=0&page=2"
    //url3 := "http://trade.dgtle.com/dgtle_module.php?mod=trade&ac=index&typeid=&PName=&searchsort=0&page=3"
    body,err := getUrl(url1)
    if err!=nil{
        log.WithField("err",err).Error("Errors Occurred")
    }

    //str := string(body)
    //log.Debug(str)
    doc, err := goquery.NewDocumentFromReader(body)
    if err != nil {
        log.Fatal(err)
    }

    // Find the review items
    doc.Find(".tradebox").Each(func(i int, s *goquery.Selection) {
        // For each item found, get the band and title
        title, _ := s.Find("p.tradetitle").Attr("title")
        price:= s.Find("p.tradeprice").Text()
        fmt.Printf("Item %d: %s - %s\n", i, title, price)
    })
}

func getUrl(url string) (io.Reader,error){

    req,err :=http.NewRequest(http.MethodGet,url,nil)
    c:=&http.Client{}
    if err!=nil{
        log.WithField("err",err).Error("Errors Occurred")
        return nil,err
    }
    req.Header.Add("Accept-Language", "zh-CN")
    req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
    req.Header.Add("Host", "trade.dgtle.com")
    req.Header.Add("Referer", "http://trade.dgtle.com/")
    log.Debugln("req:",req.URL.String())
    req.Close = true
    res, err := c.Do(req)
    if err!=nil{
        log.WithField("err",err).Error("Errors Occurred")
        return nil,err
    }
    body:= res.Body
    //defer body.Close()
    //bs,err :=ioutil.ReadAll(body)
    //if err!=nil{
    //    log.WithField("err",err).Error("Errors Occurred")
    //    return nil,err
    //}
    return body,err
}
