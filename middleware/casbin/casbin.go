package casbin

import (
    "github.com/casbin/casbin"
    "github.com/casbin/gorm-adapter"
    "github.com/mosliu/ginws/db"
    "github.com/spf13/viper"
)
var adapter *gormadapter.Adapter
//定义一个内部全局的 casbin.Enforcer 指针用来进行权限校验
var Enforcer *casbin.Enforcer

func Init(){
    casbinInit()
    //casbinMockData()
}


//制造mock数据
func casbinMockData() {
    //下面的这些命令可以用来添加规则
    log.Infoln("Mocking casbin data")
    Enforcer.AddPolicy(viper.GetString("rbac.admin_name"), viper.GetString("jwt.authPath")+viper.GetString("jwt.testPath"), "GET")
    Enforcer.AddPolicy("dex", "/auth/hello", "GET")
    Enforcer.AddRoleForUser("user_a", "user")
    Enforcer.AddRoleForUser("user_b", "user")
    Enforcer.AddRoleForUser("user_c", "user")
    Enforcer.AddPolicy("user", viper.GetString("jwt.authPath")+"/ping", "GET")
    Enforcer.SavePolicy()
}
func casbinInit(){
    //它会自动创建一个叫 casbin_rule 的表来进行规则存放
    //sqlite_info := viper.GetString("db.sqlite_path")
    //adapter := gormadapter.NewAdapter("sqlite3", sqlite_info, true)
    log.Infoln("Init casbin with db")
    adapter = gormadapter.NewAdapterByDB(db.AuthDB)

    e := casbin.NewEnforcer(viper.GetString("casbin.config"), adapter)
    //将全局的 enforcer 进行赋值，以方便在其它地方进行调用
    Enforcer = e
    //加载规则
    e.EnableLog(true)
    log.Infoln("Load casbin data")
    e.LoadPolicy()

}