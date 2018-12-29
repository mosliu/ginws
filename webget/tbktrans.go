package webget

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func TransTBK(inmsg string) {
    //open.32ds.cn
    //apkey:

    request ,_ := http.NewRequest("POST","http://gw.api.taobao.com/router/rest",nil)
    request.Header.Set("Content-Type","application/x-www-form-urlencoded")
    request.Header.Set("charset","utf-8")
    //request.Form.Add()

}
func postSample(inmsg string) {
    http.PostForm("http://127.0.0.1",
        url.Values{"name": {"ruifengyun"}, "blog": {"xiaorui.cc"},
            "aihao":{"python golang"},"content":{"nima,fuck "}})

    //url := "http://demo1.dynamsoft.com/dbr/webservice/BarcodeReaderService.svc/Read"
    //resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}
func getSample(inmsg string) {
    client := &http.Client{}
    reqest, _ := http.NewRequest("GET", "http://127.0.0.1/", nil)

    reqest.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    reqest.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
    reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
    reqest.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
    reqest.Header.Set("Cache-Control", "max-age=0")
    reqest.Header.Set("Connection", "keep-alive")
    reqest.Header.Set("User-Agent", "chrome 100")

    response, _ := client.Do(reqest)
    if response.StatusCode == 200 {
        body, _ := ioutil.ReadAll(response.Body)
        bodystr := string(body);
        fmt.Println(bodystr)
    }
    //  reqest, _ = http.NewRequest("POST","http:/127.0.0.1/", bytes.NewBufferString(data.Encode()))
    //    respet1, _ := http.NewRequest("POST", "http://127.0.0.1/", url.Values{"key": "Value"})
    //    reqest1.Header.Set("User-Agent","chrome 100")
    //    client.Do(reqest1)

}
