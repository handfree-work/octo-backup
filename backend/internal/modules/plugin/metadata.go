package plugin

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type PluginType string

const (
	TypeAccess       PluginType = "access"
	TypeRepository   PluginType = "repository"
	TypeSource       PluginType = "source"
	TypeNotification PluginType = "notification"
)

type FieldSpec struct {
	Key         string         `yaml:"key" json:"key"`
	Title       string         `yaml:"title" json:"title"`
	Type        string         `yaml:"type" json:"type"`
	Required    bool           `yaml:"required" json:"required"`
	Encrypt     bool           `yaml:"encrypt" json:"encrypt"`
	Default     any            `yaml:"default,omitempty" json:"default,omitempty"`
	Placeholder string         `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	MergeScript string         `yaml:"mergeScript,omitempty" json:"mergeScript,omitempty"`
	Component   map[string]any `yaml:"component,omitempty" json:"component,omitempty"`
	Options     []FieldOption  `yaml:"options,omitempty" json:"options,omitempty"`
}

type FieldOption struct {
	Label string `yaml:"label" json:"label"`
	Value string `yaml:"value" json:"value"`
}
type ActionSpec struct {
	Name       string `yaml:"name" json:"name"`
	Title      string `yaml:"title,omitempty" json:"title,omitempty"`
	Permission string `yaml:"permission,omitempty" json:"permission,omitempty"`
}
type Metadata struct {
	Type        PluginType   `yaml:"type" json:"type"`
	Name        string       `yaml:"name" json:"name"`
	Title       string       `yaml:"title" json:"title"`
	Description string       `yaml:"description" json:"description"`
	Version     string       `yaml:"version" json:"version"`
	AccessType  string       `yaml:"accessType,omitempty" json:"accessType,omitempty"`
	Icon        string       `yaml:"icon,omitempty" json:"icon,omitempty"`
	Group       string       `yaml:"group,omitempty" json:"group,omitempty"`
	Fields      []FieldSpec  `yaml:"fields,omitempty" json:"fields,omitempty"`
	Actions     []ActionSpec `yaml:"actions,omitempty" json:"actions,omitempty"`
}
type ActionExecutor interface {
	ExecuteAction(context.Context, string, map[string]any) (any, error)
}
type NewPluginInstance func(map[string]any) (ActionExecutor, error)
type Definition struct {
	Metadata          Metadata
	NewPluginInstance NewPluginInstance
}

var ErrActionNotFound = fmt.Errorf("插件 action 不存在")

func NewDefinition(data []byte, newPluginInstance NewPluginInstance) (*Definition, error) {
	var metadata Metadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("解析插件元数据: %w", err)
	}
	if metadata.Type == "" || metadata.Name == "" || metadata.Title == "" || metadata.Version == "" || newPluginInstance == nil {
		return nil, fmt.Errorf("插件元数据缺少必填字段")
	}
	seen := map[string]bool{}
	for _, field := range metadata.Fields {
		if field.Key == "" || seen[field.Key] {
			return nil, fmt.Errorf("插件字段无效: %s", field.Key)
		}
		seen[field.Key] = true
	}
	seen = map[string]bool{}
	for _, action := range metadata.Actions {
		if !strings.HasPrefix(action.Name, "on") || len(action.Name) <= 2 || seen[action.Name] {
			return nil, fmt.Errorf("插件 action 必须以 on 开头且不可重复: %s", action.Name)
		}
		seen[action.Name] = true
	}
	return &Definition{Metadata: metadata, NewPluginInstance: newPluginInstance}, nil
}

type Registry struct{ definitions map[string]*Definition }

func NewRegistry() *Registry { return &Registry{definitions: map[string]*Definition{}} }
func (r *Registry) Register(d *Definition) error {
	if d == nil || d.Metadata.Name == "" {
		return fmt.Errorf("插件定义不能为空")
	}
	if _, ok := r.definitions[d.Metadata.Name]; ok {
		return fmt.Errorf("插件已注册: %s", d.Metadata.Name)
	}
	r.definitions[d.Metadata.Name] = d
	return nil
}
func (r *Registry) Get(name string) (*Definition, bool) { d, ok := r.definitions[name]; return d, ok }
func (r *Registry) Metadata(t PluginType, name string) []Metadata {
	out := []Metadata{}
	for _, d := range r.definitions {
		if (t == "" || d.Metadata.Type == t) && (name == "" || d.Metadata.Name == name) {
			out = append(out, d.Metadata)
		}
	}
	return out
}
func (r *Registry) Execute(ctx context.Context, name, action string, params map[string]any) (any, error) {
	return r.ExecuteWithConfig(ctx, name, action, params, params)
}
func (r *Registry) ExecuteWithConfig(ctx context.Context, name, action string, config, params map[string]any) (any, error) {
	d, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("插件不存在: %s", name)
	}
	allowed := false
	for _, a := range d.Metadata.Actions {
		if a.Name == action {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, ErrActionNotFound
	}
	p, err := d.NewPluginInstance(config)
	if err != nil {
		return nil, err
	}
	return p.ExecuteAction(ctx, action, params)
}
