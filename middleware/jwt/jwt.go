package jwt

import (
    "github.com/gin-gonic/gin"
    "github.com/mosliu/ginws/db"
    "github.com/mosliu/ginws/middleware/casbin"
    "github.com/mosliu/ginws/models"
    "github.com/spf13/viper"

    //_ "github.com/jinzhu/gorm/dialects/postgres"
    //_ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "time"
)

//返回jwt中间件
func GetJwtMidlleware() *GinJWTMiddleware {
    var identityKey = "id"
    viper.SetDefault("jwt.realm", "Ginws")
    viper.SetDefault("jwt.key", "I do not know what key you use")
    return &GinJWTMiddleware{
        //Realm name to display to the user. Required.
        //必要项，显示给用户看的域
        Realm: viper.GetString("jwt.realm"),
        //Secret key used for signing. Required.
        //用来进行签名的密钥，就是加盐用的
        Key: []byte(viper.GetString("jwt.key")),
        PayloadFunc: func(data interface{}) MapClaims {
            if v, ok := data.(*models.User); ok {
                return MapClaims{
                    identityKey: v.UserID,
                }
            }
            return MapClaims{}
        },
        //Duration that a jwt token is valid. Optional, defaults to one hour
        //JWT 的有效时间，默认为一小时
        Timeout: 12 * time.Hour,
        // This field allows clients to refresh their token until MaxRefresh has passed.
        // Note that clients can refresh their token in the last moment of MaxRefresh.
        // This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
        // Optional, defaults to 0 meaning not refreshable.
        //最长的刷新时间，用来给客户端自己刷新 token 用的
        MaxRefresh: time.Hour,
        // Callback function that should perform the authentication of the user based on userID and
        // password. Must return true on success, false on failure. Required.
        // Option return user data, if so, user data will be stored in Claim Array.
        //必要项, 这个函数用来判断 User 信息是否合法，如果合法就反馈 true，否则就是 false, 认证的逻辑就在这里
        Authenticator: Authenticator,
        // Callback function that should perform the authorization of the authenticated user. Called
        // only after an authentication success. Must return true on success, false on failure.
        // Optional, default to success
        //可选项，用来在 Authenticator 认证成功的基础上进一步的检验用户是否有权限，默认为 success
        Authorizator: authPrivCallback,
        // User can define own Unauthorized func.
        //可以用来息定义如果认证不成功的的处理函数
        Unauthorized: unAuthFunc,
        // TokenLookup is a string in the form of "<source>:<name>" that is used
        // to extract token from the request.
        // Optional. Default value "header:Authorization".
        // Possible values:
        // - "header:<name>"
        // - "query:<name>"
        // - "cookie:<name>"
        //这个变量定义了从请求中解析 token 的位置和格式
        TokenLookup: viper.GetString("jwt.tokenLookup"),
        // TokenLookup: "query:token",
        // TokenLookup: "cookie:token",

        // TokenHeadName is a string in the header. Default value is "Bearer"
        //TokenHeadName 是一个头部信息中的字符串
        TokenHeadName: viper.GetString("jwt.tokenHeadName"),
        // TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
        //这个指定了提供当前时间的函数，也可以自定义
        TimeFunc: time.Now,
    }
}

//用来决断用户id和密码是否有效
func Authenticator(c *gin.Context) (interface{}, error) {
    var user models.User
    //log.Warnln(c.Request.Method, c.Request.Header.Get("Content-Type"))
    err := c.Bind(&user)
    if err != nil {
        return "", ErrMissingLoginValues
    }
    //if err := c.ShouldBind(&userform); err != nil {
    //    return "", jwt.ErrMissingLoginValues
    //}
    log.Infoln("user logging in:", user)

    userId := user.UserID
    password := user.Password

    //如果这条记录存在的的情况下
    //if !db.AuthDB.Where("user_id = ?", userId).Find(&user).RecordNotFound() {
    lookupUsr := &models.User{UserID: userId}
    if !db.AuthDB.Where(lookupUsr).Find(&user).RecordNotFound() {
        //定义一个临时的结构对象
        tmpUser := models.User{} //创建一个临时的存放空间
        //将 user_id 为认证信息中的 密码找出来(目前密码是明文的，这个其实不安全，可以通过加盐哈希将结果进行对比的方式以提高安全等级，这里只作原理演示，就不搞那么复杂了)
        //找到后放到前面定义的临时结构变量里
        //db.AuthDB.Where("user_id = ?", userId).Find(&tmpUser)
        db.AuthDB.Where(lookupUsr).Find(&tmpUser)
        //对比，如果密码也相同，就代表认证成功了
        if tmpUser.Password == password {
            //反馈相关信息和 true 的值，代表成功
            //return &models.User{
            //    UserID: userId,
            //}, nil
            return lookupUsr, nil
        }
    }
    return nil, ErrFailedAuthentication
}

//定义一个回调函数，用来决断用户在认证成功的前提下，是否有权限对资源进行访问
func authPrivCallback(user interface{}, c *gin.Context) bool {
    if v, ok := user.(string); ok {
        //如果可以正常取出 user 的值，就使用 casbin 来验证一下是否具备资源的访问权限
        return casbin.Enforcer.Enforce(v, c.Request.URL.String(), c.Request.Method)
    }
    //默认策略是不允许
    return false
}

//定义一个函数用来处理，认证不成功的情况
func unAuthFunc(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "code":    code,
        "message": message,
    })
}
