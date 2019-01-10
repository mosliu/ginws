package jwt

import (
    "github.com/casbin/casbin"
    "github.com/casbin/gorm-adapter"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/mosliu/ginws/db"

    //_ "github.com/jinzhu/gorm/dialects/postgres"
    //_ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/spf13/viper"
    "time"
)

//定义一个内部全局的 casbin.Enforcer 指针用来进行权限校验
var casbinEnforcer *casbin.Enforcer

//定义一函数用来处理请求
func HelloHandler(c *gin.Context) {
    //func ExtractClaims(c *gin.Context) jwt.MapClaims
    //ExtractClaims help to extract the JWT claims
    //用来将 Context 中的数据解析出来赋值给 claims
    //其实是解析出来了 JWT_PAYLOAD 里的内容
    /*
        func ExtractClaims(c *gin.Context) jwt.MapClaims {
        claims, exists := c.Get("JWT_PAYLOAD")
        if !exists {
            return make(jwt.MapClaims)
        }

        return claims.(jwt.MapClaims)
        }
    */
    claims := ExtractClaims(c)
    //func (c *Context) JSON(code int, obj interface{})
    //JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
    //将内容序列化成 JSON 格式，然后放到响应 body 里，同时将 Content-Type 置为 "application/json"
    c.JSON(200, gin.H{
        "userID": claims["id"],
        "text":   "Hello World",
    })
}

//定义一个 User 的结构体, 用来存放用户名和密码
type User struct {
    gorm.Model //加入此行用于在数据库中创建记录的 mate 数据
    UserID   string `gorm:"type:varchar(30);UNIQUE;unique_index;not null" form:"username" json:"username" binding:"required"`
    Password string `gorm:"size:255" form:"password" json:"password" binding:"required"`
}

//用来决断用户id和密码是否有效
func Authenticator(c *gin.Context) (interface{}, error) {
    var user User
    //body,_ := c.Request.GetBody()
    //bytes, e := ioutil.ReadAll(body)
    //if e!=nil{
    //    return "", jwt.ErrMissingLoginValues
    //}
    //log.Warnln(string(bytes))
    log.Warnln(c.Request.Method, c.Request.Header.Get("Content-Type"))
    err := c.Bind(&user)
    if err != nil {
        return "", ErrMissingLoginValues
    }
    //if err := c.ShouldBind(&userform); err != nil {
    //    return "", jwt.ErrMissingLoginValues
    //}
    log.Warnln(user)
    userId := user.UserID
    password := user.Password

    //如果这条记录存在的的情况下
    if !db.AuthDB.Where("user_id = ?", userId).Find(&user).RecordNotFound() {
        //定义一个临时的结构对象
        queryRes := User{} //创建一个临时的存放空间
        //将 user_id 为认证信息中的 密码找出来(目前密码是明文的，这个其实不安全，可以通过加盐哈希将结果进行对比的方式以提高安全等级，这里只作原理演示，就不搞那么复杂了)
        //找到后放到前面定义的临时结构变量里
        db.AuthDB.Where("user_id = ?", userId).Find(&queryRes)
        //对比，如果密码也相同，就代表认证成功了
        if queryRes.Password == password {
            //反馈相关信息和 true 的值，代表成功
            return &User{
                UserID: userId,
            }, nil
        }
    }

    return nil, ErrFailedAuthentication
}

//定义一个回调函数，用来决断用户在认证成功的前提下，是否有权限对资源进行访问
func authPrivCallback(user interface{}, c *gin.Context) bool {
    if v, ok := user.(string); ok {
        //如果可以正常取出 user 的值，就使用 casbin 来验证一下是否具备资源的访问权限
        return casbinEnforcer.Enforce(v, c.Request.URL.String(), c.Request.Method)
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

//定义一个中间件，用来反馈 jwt 的认证逻辑
//这里将相应的配置直接以变量的方式传递进来了
func AddMidd(v_realm, v_key, v_tokenLookup, v_tokenHeadName string) *GinJWTMiddleware {
    var identityKey = "id"
    return &GinJWTMiddleware{
        //Realm name to display to the user. Required.
        //必要项，显示给用户看的域
        Realm: v_realm,
        //Secret key used for signing. Required.
        //用来进行签名的密钥，就是加盐用的
        Key: []byte(v_key),
        PayloadFunc: func(data interface{}) MapClaims {
            if v, ok := data.(*User); ok {
                return MapClaims{
                    identityKey: v.UserID,
                }
            }
            return MapClaims{}
        },
        //Duration that a jwt token is valid. Optional, defaults to one hour
        //JWT 的有效时间，默认为一小时
        Timeout: time.Hour,
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
        TokenLookup: v_tokenLookup,
        // TokenLookup: "query:token",
        // TokenLookup: "cookie:token",

        // TokenHeadName is a string in the header. Default value is "Bearer"
        //TokenHeadName 是一个头部信息中的字符串
        TokenHeadName: v_tokenHeadName,
        // TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
        //这个指定了提供当前时间的函数，也可以自定义
        TimeFunc: time.Now,
    }
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

func InitDbAndCasbin() {
    //构建一个 pg 的连接串
    //pg_conn_info := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", v.GetString("db_host"), v.GetString("db_port"), v.GetString("db_user"), v.GetString("db_name"), v.GetString("db_password"))
    //db, err := gorm.Open("postgres", pg_conn_info)
    //gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/dbname?charset=utf8&parseTime=True&loc=Local")

    //func (s *DB) Close() error
    //Close close current db connection. If database connection is not an io.Closer, returns an error.
    //panic 之后并不是直接退出，而是先去执行 defer 的内容
    //关闭当前的 db 连接
    //defer db.Close() //如果出错，先将 db 关掉
    //func (s *DB) AutoMigrate(values ...interface{}) *DB
    //AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
    //AutoMigrate 会自动将给定的模型进行迁移，只会添加缺失的字段，并不会删除或者修改当前的字段
    db.AuthDB.AutoMigrate(&User{})
    //创建一个结构变量
    user := User{}
    //如果 db 中没有这条记录，就创建，如果有就忽略掉
    if db.AuthDB.Where("user_id = ?", viper.GetString("rbac.admin_name")).Find(&user).RecordNotFound() {
        user := User{UserID: viper.GetString("rbac.admin_name"), Password: viper.GetString("rbac.admin_pass")}
        db.AuthDB.Create(&user)
    }
    //这里有一个坑，是通过看源码解决的
    //如果不手动指定，而是自动创建时，它会默认首先要求一个很大的 postgresql 权限，尝试在 postgres 下面创建一个库，如果过程报错，就直接 panic 出来
    //这里我就通过手动创建的方式，直接指到 testdb 里
    //它会自动创建一个叫 casbin_rule 的表来进行规则存放
    sqlite_info := viper.GetString("db.sqlite_path")
    casbin_adapter := gormadapter.NewAdapter("sqlite3", sqlite_info, true)
    //casbin_adapter := gormadapter.NewAdapter("postgres", pg_conn_info, true)
    //使用前面定的 casbin_adapter 来构建 enforcer
    e := casbin.NewEnforcer(viper.GetString("casbin.config"), casbin_adapter)
    //将全局的 enforcer 进行赋值，以方便在其它地方进行调用
    casbinEnforcer = e
    //加载规则
    e.EnableLog(true)
    e.LoadPolicy()
    //下面的这些命令可以用来添加规则
    e.AddPolicy(viper.GetString("rbac.admin_name"), viper.GetString("jwt.authPath")+viper.GetString("jwt.testPath"), "GET")
    e.AddPolicy("dex", "/auth/hello", "GET")
    e.AddRoleForUser("user_a", "user")
    e.AddRoleForUser("user_b", "user")
    e.AddRoleForUser("user_c", "user")
    e.AddPolicy("user", viper.GetString("jwt.authPath")+"/ping", "GET")
    e.SavePolicy()
    log.Warnln(e.GetPolicy())
}
