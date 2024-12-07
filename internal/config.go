package internal

import (
	"errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
)

func GetWorkPath() string {
	workDir, _ := os.Getwd()
	return workDir
}
func GetProcessPath() string {
	// 添加运行目录
	process, _ := os.Executable()
	return filepath.Dir(process)
}

type RootConfig struct {
	RemotePath      string   `mapstructure:"remote_path" json:"remote_path" yaml:"remote_path"`
	SuccessDelete   bool     `mapstructure:"success_delete" json:"success_delete" yaml:"success_delete"`
	RemoveStr       []string `mapstructure:"remove_str" json:"remove_str" yaml:"remove_str"`
	RemoveReg       string   `mapstructure:"remove_reg" json:"remove_reg" yaml:"remove_reg"`
	UploadExtension []string `mapstructure:"upload_extension" json:"upload_extension" yaml:"upload_extension"`
	IgnorePath      []string `mapstructure:"ignore_path" json:"ignore_path" yaml:"ignore_path"`
}

var Config = RootConfig{
	RemotePath:      "/",
	SuccessDelete:   false,
	RemoveStr:       []string{},
	RemoveReg:       "",
	UploadExtension: []string{},
	IgnorePath:      []string{},
}

func InitConfig() {
	configName := "pan-upload"
	// 添加运行目录
	v := viper.New()
	v.AddConfigPath(GetProcessPath())

	// 添加当前目录
	v.AddConfigPath(GetWorkPath())
	v.SetConfigName(configName)
	if err := v.ReadInConfig(); err != nil { // 读取配置文件
		// 使用类型断言检查是否为 *os.PathError 类型
		var pathErr viper.ConfigFileNotFoundError
		if errors.As(err, &pathErr) {
			val := reflect.ValueOf(Config)
			for i := 0; i < val.NumField(); i++ {
				// 获取字段名
				name := val.Type().Field(i).Tag.Get("mapstructure")
				// 获取字段值
				value := val.Field(i).Interface()
				v.SetDefault(name, value)
			}
			err = v.WriteConfigAs(GetWorkPath() + "/" + configName + ".yaml")
			if err != nil {
				panic(err)
			} else {
				// 重新读取已经写入的文件
				_ = v.ReadInConfig()
			}
		} else {
			panic(err)
		}
	}

	if err := v.Unmarshal(&Config); err != nil { // 解码配置文件到结构体
		panic(err)
	}
}
