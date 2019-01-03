package webget

import (
    "crypto/tls"
    "fmt"
    "github.com/mosliu/ginws/utils"
    "github.com/spf13/viper"
    "io/ioutil"
    "net/http"
    "net/url"
)

//success content
/*
{
    "code":200,
    "data":{
        "url":"https:\/\/s.click.taobao.com\/BPempHw?un=d707d2abce2ca380&share_crt_v=1&ut_sk=1.utdid_null_1515140911052.TaoPassword-Outside.taoketop&sp_tk=77+lNGtKY2JweFNzOWHvv6U=",
        "title":"盼盼艾比利薯片办公室零食膨化点心香脆大礼包10元以下整箱批发",
        "pic":"https:\/\/img.alicdn.com\/bao\/uploaded\/i4\/3485532766\/O1CN01tCkTAO1WIs0admtK8_!!0-item_pic.jpg",
        "ownerid":"0",
        "youxiaoqi":"1549037374711"
    },
    "msg":"数据获取成功."
}
*/
type Msg struct{
    Code int `json:"code"`
    Data MsgData
    Msg string
}

type MsgData struct{
    Url string
    Title string
    pic string
    Ownerid string
    Youxiaoqi string
}

//转换淘口令
func TransTKL(inmsg string) (ok bool, url string) {
    //open.32ds.cn
    ok = false

    client := &http.Client{}

    //忽略https的证书
    client.Transport = &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    //获取apkey 写在私有的配置文件内了 //configs/private.toml
    apkey := viper.GetString("tbk.apkey")

    const apiurl = "https://api.open.21ds.cn/apiv1/jiexitkl"

    request, _ := http.NewRequest(
        http.MethodGet,
        apiurl,
        nil,
    )

    //request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    request.Header.Set("charset", "utf-8")
    //生成Query String
    query := request.URL.Query()
    query.Add("apkey", apkey)
    query.Add("kouling", "￥4kJcb4xSs9a￥")
    //query.Add("kouling", "￥4kJcbpxSs9a￥")
    //query.Add("kouling", url)
    //设回
    request.URL.RawQuery = query.Encode()

    log.Debug(request.URL.String())
    //400
    response, _ := client.Do(request)
    log.WithField("status", response.StatusCode).Debug("Request executed.")
    if response.StatusCode == 200 {
        //获取成功处理
        raw := response.Body
        defer raw.Close()
        body, err := ioutil.ReadAll(raw)
        utils.CheckError(err,log)
        var msg Msg
        json.Unmarshal(body, &msg)
        bodystr := string(body)
        fmt.Println(bodystr)
        log.Debug(msg)

    } else {
        log.WithField("status", response.StatusCode).Debug("Request error!")
    }
    return
}
func postSample(inmsg string) {
    http.PostForm("http://127.0.0.1",
        url.Values{"name": {"ruifengyun"}, "blog": {"xiaorui.cc"},
            "aihao": {"python golang"}, "content": {"nima,fuck "}})

    //url := "http://demo1.dynamsoft.com/dbr/webservice/BarcodeReaderService.svc/Read"
    //resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}
