package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Conf = new(Appconfig)

type Appconfig struct {
	DownLoadUrl string `mapstructure:"downloadUrl"`
	ToolName    string `mapstructure:"toolName"`
	ToolMd5     string `mapstructure:"toolMd5"`
	SoftVersion string `mapstructure:"softVersion"`
}

func Init(file string) (err error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// 配置文件不存在，设置默认值
		Conf.DownLoadUrl = "http://powernod.com/u/718433caa9a32077"
		Conf.ToolName = "小生——QR增强.exe"
		Conf.ToolMd5 = "77411af8d866d163be773a86dbbe1c24"
		Conf.SoftVersion = "0.9"
		log.Printf("配置文件 %s 不存在，使用默认配置。\n", file)
		return nil
	} else {
		// 这个指定的是main可执行文件所在的目录与yaml文件相对关系。不是setting.go文件和yaml相对关系
		log.Println("file:", file)
		if len(file) == 0 {
			viper.SetConfigFile(file)
		} else {
			viper.SetConfigFile("./config.yaml")
		}
		viper.AddConfigPath(".")              // 查找配置文件所在路径
		viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
		err = viper.ReadInConfig()            // 查找并读取配置文件
		if err != nil {
			//log.Println("main: viper initial failed:", err.Error())
			return err
			//panic(log.Errorf("main: config initial failed"))
		}

		// 反序列化
		if err := viper.Unmarshal(Conf); err != nil {
			log.Println("[*]setting: viper unmarshal failed:", err)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Printf("[*]config changed:...:%v\n", in.Name)
			// 更新
			if err2 := viper.Unmarshal(Conf); err2 != nil {
				log.Println("[*]setting: config update failed:", err.Error())
			}
		})
		return nil
	}

}
