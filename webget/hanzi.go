package webget

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetHanZi(inmsg string) {
	//http://www.xjihe.com/service/apidemo/19?urlId=9
	//open.32ds.cn
	//apkey:
	url := "http://www.xjihe.com/api/edu/hanyudic"
	apikey := "Egk2nsfGIS0PL7hit1xUMfE6sVVrNgGs"
	//apikey := "sQS2ylErlfm9Ao2oNPqw6TqMYbJjbs4g"
	client := &http.Client{}
	request ,err := http.NewRequest(http.MethodGet,url,nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	request.Header.Set("charset","utf-8")
	request.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Add("apikey", apikey)
	request.Header.Add("Referer","http://www.xjihe.com/service/apidemo/19?urlId=9")

	q := request.URL.Query()
	q.Add("query","轩")
	q.Add("mode","1")
	q.Add("page","1")
	request.URL.RawQuery = q.Encode()

	fmt.Println(request.URL.String())

	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//json.Unmarshal(body)
	str := string(body)

	log.Info(str)
	//request.Form.Add()

}