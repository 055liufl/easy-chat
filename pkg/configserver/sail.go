// =============================================================================
// Sail - Sail 配置中心客户端封装
// =============================================================================
// 封装 Sail 配置中心客户端，提供配置拉取和热更新功能。
//
// 功能特性:
//   - 从 ETCD 配置中心拉取配置
//   - 支持配置热更新（监听配置变化）
//   - 自动合并多个配置文件
//   - 支持本地配置文件缓存
//
// 使用场景:
//   - 微服务配置集中管理
//   - 配置动态更新，无需重启服务
//   - 多环境配置隔离（开发、测试、生产）
//
// 设计思路:
//   - 封装 sail-client 库，简化使用
//   - 通过 ETCD 实现配置的分布式存储和同步
//   - 支持配置变更回调，实现热更新
//
// 项目中的应用:
//   - 所有微服务的配置管理
//   - 数据库连接、Redis 连接等配置的动态更新
//   - 业务参数的在线调整
//
// 第三方依赖:
//   - github.com/HYY-yu/sail-client: Sail 配置中心客户端库
//
// =============================================================================
package configserver

import (
	"encoding/json"
	"fmt"
	"github.com/HYY-yu/sail-client"
)

// Config Sail 配置中心的配置结构
type Config struct {
	ETCDEndpoints  string `toml:"etcd_endpoints"`   // 逗号分隔的 ETCD 地址，例如: 0.0.0.0:2379,0.0.0.0:12379,0.0.0.0:22379
	ProjectKey     string `toml:"project_key"`      // 项目标识，用于区分不同项目的配置
	Namespace      string `toml:"namespace"`        // 命名空间，用于区分不同环境（dev/test/prod）
	Configs        string `toml:"configs"`          // 配置文件列表，逗号分隔
	ConfigFilePath string `toml:"config_file_path"` // 本地配置文件存放路径，空代表不存储本地配置文件
	LogLevel       string `toml:"log_level"`        // 日志级别（DEBUG/INFO/WARN/ERROR），默认 WARN
}

// Sail Sail 配置中心客户端封装
type Sail struct {
	*sail.Sail                // 嵌入 Sail 客户端
	sail.OnConfigChange       // 配置变更回调函数
	c                   *Config // Sail 配置
}

// NewSail 创建 Sail 配置中心客户端实例
// 初始化 Sail 客户端，但不立即连接配置中心
//
// 参数:
//   - cfg: Sail 配置中心的配置信息
//
// 返回:
//   - *Sail: Sail 客户端实例
//
// 注意:
//   - 此方法只是创建实例，需要调用 Build() 方法才会真正连接配置中心
//
// 示例:
//   sail := NewSail(&Config{
//       ETCDEndpoints: "127.0.0.1:2379",
//       ProjectKey: "easy-chat",
//       Namespace: "prod",
//       Configs: "app.yaml,db.yaml",
//   })
func NewSail(cfg *Config) *Sail {
	return &Sail{c: cfg}
}

// Build 构建 Sail 客户端连接
// 连接到 ETCD 配置中心，初始化配置监听
//
// 返回:
//   - error: 连接配置中心时的错误
//
// 工作流程:
//   1. 如果设置了配置变更回调，添加到选项中
//   2. 创建 Sail 客户端实例
//   3. 连接到 ETCD 配置中心
//   4. 开始监听配置变化
//
// 示例:
//   sail := NewSail(cfg)
//   if err := sail.Build(); err != nil {
//       log.Fatal(err)
//   }
func (s *Sail) Build() error {
	var opts []sail.Option

	if s.OnConfigChange != nil {
		opts = append(opts, sail.WithOnConfigChange(s.OnConfigChange))
	}

	s.Sail = sail.New(&sail.MetaConfig{
		ETCDEndpoints:  s.c.ETCDEndpoints,
		ProjectKey:     s.c.ProjectKey,
		Namespace:      s.c.Namespace,
		Configs:        s.c.Configs,
		ConfigFilePath: s.c.ConfigFilePath,
		LogLevel:       s.c.LogLevel,
	}, opts...)
	return s.Sail.Err()
}

// FromJsonBytes 从配置中心获取配置并转换为 JSON 字节数组
// 拉取最新配置，合并所有配置文件，并转换为 JSON 格式
//
// 返回:
//   - []byte: JSON 格式的配置数据
//   - error: 拉取或转换配置时的错误
//
// 工作流程:
//   1. 从 ETCD 拉取最新配置
//   2. 调用内部方法转换为 JSON
//
// 示例:
//   data, err := sail.FromJsonBytes()
//   if err != nil {
//       log.Fatal(err)
//   }
func (s *Sail) FromJsonBytes() ([]byte, error) {
	if err := s.Pull(); err != nil {
		return nil, err
	}

	return s.fromJsonBytes(s.Sail)
}

// fromJsonBytes 将 Sail 配置转换为 JSON 字节数组
// 内部方法，用于将 Viper 配置对象转换为 JSON 格式
//
// 参数:
//   - sail: Sail 客户端实例
//
// 返回:
//   - []byte: JSON 格式的配置数据
//   - error: 转换配置时的错误
//
// 工作流程:
//   1. 合并所有配置文件的 Viper 实例
//   2. 获取所有配置项
//   3. 序列化为 JSON 格式
func (s *Sail) fromJsonBytes(sail *sail.Sail) ([]byte, error) {
	v, err := sail.MergeVipers()
	if err != nil {
		return nil, err
	}
	data := v.AllSettings()
	return json.Marshal(data)
}

// SetOnChange 设置配置变更回调函数
// 当配置中心的配置发生变化时，会调用此回调函数
//
// 参数:
//   - f: 配置变更回调函数
//
// 工作流程:
//   1. 监听配置文件变化
//   2. 配置变化时，将新配置转换为 JSON 格式
//   3. 调用用户提供的回调函数处理配置变更
//
// 使用场景:
//   - 数据库连接配置变更时，重新初始化连接池
//   - Redis 配置变更时，重新建立连接
//   - 业务参数变更时，更新内存缓存
//
// 示例:
//   sail.SetOnChange(func(data []byte) error {
//       var cfg Config
//       if err := json.Unmarshal(data, &cfg); err != nil {
//           return err
//       }
//       // 处理配置变更
//       return nil
//   })
func (s *Sail) SetOnChange(f OnChange) {
	s.OnConfigChange = func(configFileKey string, sail *sail.Sail) {
		data, err := s.fromJsonBytes(sail)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = f(data); err != nil {
			fmt.Println("OnChange err ", err)
		}

	}
}
