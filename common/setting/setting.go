package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
)

var (
	Cfg     *ini.File
	RunMode string
	WorkPath string
	RunPath string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Database map[string]string
	App      map[string]interface{}
)

func init() {
	var err error
	if WorkPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	RunPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	var filename = "env.ini"
	filename = filepath.Join(RunPath, "config",filename)
	if os.Getenv("CKGO_WORK_DIR") != "" {
		filename = os.Getenv("CKGO_WORK_DIR") + "/config/env.ini"
	}
	Cfg, err = ini.Load(filename)
	if err != nil {
		log.Fatalf("Fail to parse 'config/env.ini': %v", err)
	}
	WorkPath = Cfg.Section("config").Key("ROOT_PATH").MustString(filepath.Join(WorkDir(), ""))
	LoadBase()
	LoadServer()
	LoadApp()
	LoadDatabase()
}

// LoadBase 加载基础配置
func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

// LoadServer 加载服务配置
func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

// LoadApp 加载app配置
func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}
	App = make(map[string]interface{})
	App["JwtSecret"] = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	App["PageSize"] = sec.Key("PAGE_SIZE").MustInt(10)
	App["IdentityKey"] = sec.Key("IDENTITY_KEY").MustString("idname")
	App["LogPath"] = sec.Key("LOG_PATH").MustString("../runtime/debug.log")
	App["MaxSize"] = sec.Key("MaxSize").MustInt(200)
	App["MaxBackups"] = sec.Key("MaxBackups").MustInt(10)
	App["MaxAge"] = sec.Key("MaxAge").MustInt(7)
	App["Level"] = sec.Key("Level").MustString("debug")
}

// LoadDatabase 加载数据库配置
func LoadDatabase() {
	sec, err := Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	Database = make(map[string]string)
	Database["Type"] = sec.Key("TYPE").MustString("mysql")
	Database["User"] = sec.Key("USER").MustString("root")
	Database["Password"] = sec.Key("PASSWORD").MustString("123456")
	Database["Host"] = sec.Key("HOST").MustString("127.0.0.1")
	Database["Port"] = sec.Key("PORT").MustString("3306")
	Database["Name"] = sec.Key("NAME").MustString("ckgo")
	Database["Prefix"] = sec.Key("TABLE_PREFIX").MustString("go_")
}
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}