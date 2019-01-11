package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/mosliu/ginws/middleware/jwt"
)

//经过确权，从payload中解析claims
func HelloHandler(c *gin.Context) {
    //func ExtractClaims(c *gin.Context) jwt.MapClaims
    //ExtractClaims help to extract the JWT claims
    //用来将 Context 中的数据解析出来赋值给 claims
    //其实是解析出来了 JWT_PAYLOAD 里的内容
    claims := jwt.ExtractClaims(c)
    //func (c *Context) JSON(code int, obj interface{})
    //JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
    //将内容序列化成 JSON 格式，然后放到响应 body 里，同时将 Content-Type 置为 "application/json"
    c.JSON(200, gin.H{
        "userID": claims["id"],
        "text":   "Hello World",
    })
}

//定义一个 ping 的处理函数
func H_ping(c *gin.Context) {
    //只作简单的反馈
    //反馈内容由此规定
    //JSON 格式反馈 pong
    c.JSON(200, gin.H{
        "message": "pong",
    })
}
