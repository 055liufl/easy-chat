// =============================================================================
// ConfigServer - 配置服务器抽象层
// =============================================================================
// 提供统一的配置加载接口，支持本地文件和远程配置中心两种方式。
//
// 功能特性:
//   - 支持本地配置文件加载（使用 go-zero 的 conf 包）
//   - 支持远程配置中心加载（如 Sail 配置中心）
//   - 支持配置热更新（通过 OnChange 回调）
//   - 自动将配置转换为 JSON 格式
//
// 使用场景:
//   - 微服务配置管理：统一管理所有服务的配置
//   - 配置热更新：无需重启服务即可更新配置
//   - 多环境配置：开发、测试、生产环境配置隔离
//
// 设计思路:
//   - 定义 ConfigServer 接口，支持多种配置源实现
//   - 优先使用远程配置中心，降级使用本地文件
//   - 通过回调函数实现配置变更通知
//
// 项目中的应用:
//   - 所有微服务启动时加载配置
//   - 支持配置中心动态更新配置
//   - 配置变更时自动重新加载
//
// =============================================================================
package configserver

import (
	"errors"
	"github.com/zeromicro/go-zero/core/conf"
)

// ErrNotSetConfig 未设置配置信息错误
// 当既没有指定本地配置文件，也没有指定配置服务器时返回此错误
var ErrNotSetConfig = errors.New("未设置配置信息")

// OnChange 配置变更回调函数类型
// 当配置发生变化时，会调用此函数并传入新的配置数据（JSON 格式）
//
// 参数:
//   - []byte: 新的配置数据（JSON 格式）
//
// 返回:
//   - error: 处理配置变更时的错误
type OnChange func([]byte) error

// ConfigServer 配置服务器接口
// 定义了配置服务器必须实现的方法
type ConfigServer interface {
	// Build 构建配置服务器连接
	// 初始化配置服务器客户端，建立与配置中心的连接
	Build() error

	// SetOnChange 设置配置变更回调函数
	// 当配置发生变化时，会调用此回调函数
	SetOnChange(OnChange)

	// FromJsonBytes 获取配置的 JSON 字节数组
	// 从配置服务器拉取最新配置并转换为 JSON 格式
	FromJsonBytes() ([]byte, error)
}

// configServer 配置服务器实现
// 封装了本地文件和远程配置中心两种配置加载方式
type configServer struct {
	ConfigServer        // 配置服务器接口实现（如 Sail）
	configFile   string // 本地配置文件路径
}

// NewConfigServer 创建配置服务器实例
// 支持本地文件和远程配置中心两种配置源
//
// 参数:
//   - configFile: 本地配置文件路径（可为空）
//   - s: 配置服务器接口实现（如 Sail，可为 nil）
//
// 返回:
//   - *configServer: 配置服务器实例
//
// 使用场景:
//   - 开发环境：使用本地配置文件
//   - 生产环境：使用远程配置中心
//
// 示例:
//   // 使用本地配置文件
//   cs := NewConfigServer("config.yaml", nil)
//
//   // 使用 Sail 配置中心
//   sail := NewSail(&Config{...})
//   cs := NewConfigServer("", sail)
func NewConfigServer(configFile string, s ConfigServer) *configServer {
	return &configServer{
		ConfigServer: s,
		configFile:   configFile,
	}
}

// MustLoad 加载配置到指定的结构体
// 支持本地文件和远程配置中心两种方式，优先使用配置中心
//
// 参数:
//   - v: 配置结构体指针，用于接收解析后的配置
//   - onChange: 配置变更回调函数（可为 nil）
//
// 返回:
//   - error: 加载配置时的错误
//
// 加载策略:
//   1. 如果既没有本地文件也没有配置服务器，返回错误
//   2. 如果只有本地文件，使用 go-zero 的 conf.MustLoad 加载
//   3. 如果有配置服务器，从配置中心加载并支持热更新
//
// 工作流程:
//   1. 设置配置变更回调函数
//   2. 构建配置服务器连接
//   3. 从配置服务器获取配置数据（JSON 格式）
//   4. 将 JSON 数据解析到结构体
//
// 示例:
//   var cfg Config
//   cs := NewConfigServer("", sail)
//   err := cs.MustLoad(&cfg, func(data []byte) error {
//       // 配置变更时的处理逻辑
//       return LoadFromJsonBytes(data, &cfg)
//   })
func (s *configServer) MustLoad(v any, onChange OnChange) error {
	if s.configFile == "" && s.ConfigServer == nil {
		return ErrNotSetConfig
	}

	if s.ConfigServer == nil {
		// 使用 go-zero 的默认本地文件加载
		conf.MustLoad(s.configFile, v)
		return nil
	}

	if onChange != nil {
		s.ConfigServer.SetOnChange(onChange)
	}

	if err := s.ConfigServer.Build(); err != nil {
		return err
	}

	data, err := s.ConfigServer.FromJsonBytes()
	if err != nil {
		return err
	}

	return LoadFromJsonBytes(data, v)
}

// LoadFromJsonBytes 从 JSON 字节数组加载配置到结构体
// 封装了 go-zero 的 conf.LoadFromJsonBytes 方法
//
// 参数:
//   - data: JSON 格式的配置数据
//   - v: 配置结构体指针，用于接收解析后的配置
//
// 返回:
//   - error: 解析配置时的错误
//
// 使用场景:
//   - 从配置中心获取 JSON 配置后解析
//   - 配置热更新时重新解析配置
//
// 示例:
//   var cfg Config
//   err := LoadFromJsonBytes(jsonData, &cfg)
func LoadFromJsonBytes(data []byte, v any) error {
	return conf.LoadFromJsonBytes(data, v)
}
