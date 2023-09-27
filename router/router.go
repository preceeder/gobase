package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/preceeder/gobase"
	"github.com/preceeder/gobase/utils"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

// 路由结构体
type Route struct {
	path        string            //url路径
	httpMethod  string            //http方法 get post
	Method      reflect.Value     //方法路由
	Args        []reflect.Type    //参数类型
	Middlewares []gin.HandlerFunc // 接口中间件
}

// 接口路由前缀 配置

//type ApiRouteConfig map[string]any

type ApiRouteConfig struct {
	ALL    *ApiRouterSubConfig // 这个里面就只能是一个对象了
	POST   []ApiRouterSubConfig
	GET    []ApiRouterSubConfig
	PUT    []ApiRouterSubConfig
	DELETE []ApiRouterSubConfig
}

type Str *string

type ApiRouterSubConfig struct {
	FuncName            any // func or string
	Path                any // string
	Middlewares         []gin.HandlerFunc
	NoUseModel          any // uri 是否使用model 名   bool
	NoUseBasePrefixPath any // 是否禁用 BasePrefixPathInvalid   bool
}

// 路由的前缀

var BasePrefixPath string = "/api"

var actionMap = map[string]string{
	"Get":    "GET",
	"Post":   "POST",
	"Put":    "PUT",
	"Delete": "DELETE",
}

// 路由集合
var Routes = []Route{}

func InitRouter() *gin.Engine {
	//初始化路由
	r := gin.New()
	r.Use(gobase.Cors(), gobase.GinLogger(), gobase.GinRecovery(true))
	//docs.SwaggerInfo.BasePath = "/api"
	//打开 host:port/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	Bind(r)
	return r
}

func default_l(c *gin.Context) {
	for _, v := range Routes {
		fmt.Printf("router: %v \n", v)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

// 注册控制器
func Register(controller interface{}) bool {
	ctrlName := reflect.TypeOf(controller).String()
	module := ctrlName
	if strings.Contains(ctrlName, ".") {
		module = ctrlName[strings.Index(ctrlName, ".")+1:]
	}
	v := reflect.ValueOf(controller)

	// 获取接口配置
	arc := v.Elem().FieldByName("ApiRouteConfig")
	var apiData ApiRouteConfig
	if arc.IsValid() {
		apiData = apiConfig(arc.Interface())
	}
	var pathMap map[string][]PW
	pathMap = PdHandler(apiData, strings.ToLower(module))
	// 将路由 扁平化 为 map[string]map[string]any
	//遍历方法
	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		action := v.Type().Method(i).Name
		hmd, isIn := pathMap[action]
		if !isIn {
			continue
		}
		for _, dv := range hmd {
			httpMethod := dv.HttpMethod
			//路径处理
			path := dv.Path
			middlewares := dv.Middlewares
			//遍历参数
			params := make([]reflect.Type, 0, v.NumMethod())
			for j := 0; j < method.Type().NumIn(); j++ {
				params = append(params, method.Type().In(j))
			}
			route := Route{path: path, Method: method, Args: params, httpMethod: httpMethod, Middlewares: middlewares}
			Routes = append(Routes, route)
		}
	}
	return true
}

func apiConfig(data any) (apiData ApiRouteConfig) {
	err := mapstructure.Decode(data, &apiData)
	if err != nil {
		slog.Error("decoder.Decode(data)", "error", err.Error())
		panic("api init config error ")
	}
	return
}

func subApiconfigPares(config []ApiRouterSubConfig, module string) (path string, middlewares []gin.HandlerFunc) {
	lc := len(config)
	path = ""
	if lc > 0 {
		for i := 0; i < lc; i++ {
			bp, _ := config[i].NoUseBasePrefixPath.(bool)
			if config[i].NoUseBasePrefixPath == nil || bp == false {
				path += BasePrefixPath
				break
			}
		}
		for i := 0; i < lc; i++ {
			nm, _ := config[i].NoUseModel.(bool)
			if config[i].NoUseModel == nil || nm == false {
				path += "/" + module
				break
			}
		}
		for i := 0; i < lc; i++ {
			if config[i].Path != nil {
				pt, _ := config[i].Path.(string)
				if pt != "" {
					path += pt
				}
				break
			}
		}

		for i := 0; i < lc; i++ {
			if config[i].Middlewares != nil {
				middlewares = append(middlewares, config[i].Middlewares...)
				break
			}
		}
	} else {
		path += "/" + BasePrefixPath + "/" + module
	}

	var pathSlice []string
	for _, v := range strings.Split(path, "/") {
		if v != "" {
			pathSlice = append(pathSlice, v)
		}
	}
	path = "/" + strings.Join(pathSlice, "/")
	return
}

type PW struct {
	Path        string
	Middlewares []gin.HandlerFunc
	HttpMethod  string
}

func PdHandler(apiData ApiRouteConfig, module string) map[string][]PW {
	sd := map[string][]PW{}
	var ac *ApiRouterSubConfig

	if apiData.ALL != nil {
		ac = apiData.ALL
	}

	for k, v := range actionMap {
		var tc []ApiRouterSubConfig
		err := utils.GetAttr(apiData, v, &tc)
		if err != nil {
			slog.Error("PdHandler", "error", err.Error())
		}
		if tc != nil {
			for _, vl := range tc {
				funname := ""
				if vl.FuncName == nil {
					funname = k
				} else {
					if reflect.TypeOf(vl.FuncName).Kind().String() == "func" {
						pname := runtime.FuncForPC(reflect.ValueOf(vl.FuncName).Pointer()).Name()
						funname = strings.TrimSuffix(pname, "-fm")
						fun := strings.Split(funname, ".")
						funname = fun[len(fun)-1]
					} else if reflect.TypeOf(vl.FuncName).Kind().String() == "string" {
						funname, _ = vl.FuncName.(string)
					} else {
						panic("不支持的数据类型")
					}
				}
				temp := []ApiRouterSubConfig{vl}
				if ac != nil {
					temp = append(temp, *ac)
				}
				path, middlewares := subApiconfigPares(temp, module)
				sd[funname] = append(sd[funname], PW{Path: path, Middlewares: middlewares, HttpMethod: v})
			}
		}
		if _, ok := sd[k]; !ok {
			temp := []ApiRouterSubConfig{}
			if ac != nil {
				temp = append(temp, *ac)
			}
			path, middlewares := subApiconfigPares(temp, module)
			sd[k] = append(sd[k], PW{Path: path, Middlewares: middlewares, HttpMethod: v})
			//sd[k] = PW{Path: path, Middlewares: middlewares, HttpMethod: v}
		}
	}
	return sd

}

// 绑定路由 m是方法GET POST等
// 绑定基本路由
func Bind(e *gin.Engine) {
	for _, route := range Routes {
		if route.httpMethod == "GET" {
			e.GET(route.path, append(route.Middlewares, match(route.path, route))...)
		} else if route.httpMethod == "POST" {
			e.POST(route.path, append(route.Middlewares, match(route.path, route))...)
		} else if route.httpMethod == "PUT" {
			e.PUT(route.path, append(route.Middlewares, match(route.path, route))...)
		} else if route.httpMethod == "DELETE" {
			e.DELETE(route.path, append(route.Middlewares, match(route.path, route))...)
		}
	}
}

// 根据path匹配对应的方法
func match(path string, route Route) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 第一个 / 分割后数组中最少有两个 ['', 'user']
		fields := strings.Split(path, "/")
		if len(fields) < 2 {
			slog.Error("api uri len must bigger 2", "error", "api uri len must bigger 2", "uri", utils.SliceToString[string](fields))
			return
		}

		if route.Method.IsValid() {
			arguments := make([]reflect.Value, 1)
			arguments[0] = reflect.ValueOf(c) // *gin.Context

			res := route.Method.Call(arguments)
			if res != nil {
				// 有返回结果的这里处理
				c.JSON(http.StatusOK, res[0].Interface())
			}
		}
	}
}
