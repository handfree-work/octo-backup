package logic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/modules/plugin/secret"
	"handfree-work/octo-backup/internal/svc"
)

const pluginMasterKeySetting = "plugin.master-key.v1"

type PluginService struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPluginService(ctx context.Context, s *svc.ServiceContext) *PluginService {
	return &PluginService{ctx: ctx, svcCtx: s}
}

type PluginInput struct {
	Name        string         `json:"name"`
	PluginType  string         `json:"pluginType"`
	PluginName  string         `json:"pluginName"`
	Config      map[string]any `json:"config"`
	Description string         `json:"description"`
}
type PluginPageQuery struct {
	Offset     int64  `json:"offset"`
	Limit      int64  `json:"limit"`
	PluginType string `json:"pluginType"`
	Name       string `json:"name"`
}
type PluginPageResult struct {
	Offset  int64            `json:"offset"`
	Limit   int64            `json:"limit"`
	Records []map[string]any `json:"records"`
	Total   int64            `json:"total"`
}

func (s *PluginService) Update(id int64, in *PluginInput) (map[string]any, error) {
	row, err := db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db)).GetById(id)
	if err != nil || row == nil {
		return nil, fmt.Errorf("插件实例不存在")
	}
	if in == nil {
		return s.present(row), nil
	}
	if strings.TrimSpace(in.Name) != "" {
		row.Name = strings.TrimSpace(in.Name)
	}
	row.Description = strings.TrimSpace(in.Description)
	if in.Config != nil {
		definition, ok := s.svcCtx.Plugins.Get(row.PluginName)
		if !ok {
			return nil, fmt.Errorf("插件不存在")
		}
		codec, err := s.secretCodec()
		if err != nil {
			return nil, err
		}
		var existing map[string]any
		if row.ConfigYAML != "" {
			if existing, err = codec.DecodeYAML([]byte(row.ConfigYAML), secretFields(definition.Metadata.Fields)); err != nil {
				return nil, fmt.Errorf("解析已有插件配置失败")
			}
		}
		config := mergePluginConfig(existing, in.Config, definition.Metadata.Fields)
		if err := validateRequiredPluginFields(config, definition.Metadata.Fields); err != nil {
			return nil, err
		}
		plain, e := yaml.Marshal(config)
		if e != nil {
			return nil, e
		}
		raw, _, e := codec.ProtectYAML(plain, secretFields(definition.Metadata.Fields), nil)
		if e != nil {
			return nil, e
		}
		row.ConfigYAML = string(raw)
	}
	if _, err = db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db)).UpdateById(id, row); err != nil {
		return nil, err
	}
	return s.present(row), nil
}

func (s *PluginService) Metadata(t, name string) []plugin.Metadata {
	return s.svcCtx.Plugins.Metadata(plugin.PluginType(t), name)
}
func (s *PluginService) Create(in *PluginInput) (map[string]any, error) {
	if in == nil || strings.TrimSpace(in.Name) == "" || in.PluginName == "" {
		return nil, fmt.Errorf("插件名称和类型不能为空")
	}
	d, ok := s.svcCtx.Plugins.Get(in.PluginName)
	if !ok {
		return nil, fmt.Errorf("插件不存在")
	}
	if in.PluginType != "" && string(d.Metadata.Type) != in.PluginType {
		return nil, fmt.Errorf("插件类型不匹配")
	}
	config := mergePluginConfig(nil, in.Config, d.Metadata.Fields)
	if err := validateRequiredPluginFields(config, d.Metadata.Fields); err != nil {
		return nil, err
	}
	codec, err := s.secretCodec()
	if err != nil {
		return nil, err
	}
	plain, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	raw, _, err := codec.ProtectYAML(plain, secretFields(d.Metadata.Fields), nil)
	if err != nil {
		return nil, err
	}
	row := &models.Plugin{Name: strings.TrimSpace(in.Name), PluginType: string(d.Metadata.Type), PluginName: in.PluginName, ConfigYAML: string(raw), Description: strings.TrimSpace(in.Description)}
	if err = db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db)).Create(row); err != nil {
		return nil, err
	}
	return s.present(row), nil
}
func (s *PluginService) Page(q *PluginPageQuery) (*PluginPageResult, error) {
	off, lim := int64(0), int64(20)
	typ := ""
	if q != nil {
		off, lim, typ = q.Offset, q.Limit, q.PluginType
	}
	if off < 0 {
		off = 0
	}
	if lim <= 0 {
		lim = 20
	}
	if lim > 1000 {
		lim = 1000
	}
	dao := db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db))
	page := &db_.Page{Start: off, Limit: lim}
	query := &models.Plugin{}
	if typ != "" {
		query.PluginType = typ
	}
	var filters []func(*gorm.DB)
	if q != nil && strings.TrimSpace(q.Name) != "" {
		name := strings.TrimSpace(q.Name)
		filters = append(filters, func(db *gorm.DB) { db.Where("name LIKE ?", "%"+name+"%") })
	}
	rows, err := dao.FindPage(&db_.PageReq[models.Plugin]{Query: query, Page: page}, filters...)
	if err != nil {
		return nil, err
	}
	out := []map[string]any{}
	for i := range *rows {
		out = append(out, s.present(&(*rows)[i]))
	}
	return &PluginPageResult{Offset: off, Limit: lim, Records: out, Total: page.Total}, nil
}
func (s *PluginService) Action(id int64, action string, params map[string]any) (any, error) {
	row, err := db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db)).GetById(id)
	if err != nil || row == nil {
		return nil, fmt.Errorf("插件实例不存在")
	}
	definition, ok := s.svcCtx.Plugins.Get(row.PluginName)
	if !ok {
		return nil, fmt.Errorf("插件不存在")
	}
	if row.PluginType != string(definition.Metadata.Type) {
		return nil, fmt.Errorf("插件实例类型与插件定义不匹配")
	}
	codec, err := s.secretCodec()
	if err != nil {
		return nil, err
	}
	config, err := codec.DecodeYAML([]byte(row.ConfigYAML), secretFields(definition.Metadata.Fields))
	if err != nil {
		return nil, err
	}
	return s.svcCtx.Plugins.ExecuteWithConfig(s.ctx, row.PluginName, action, config, params)
}

func (s *PluginService) Delete(id int64) error {
	dao := db_.New[models.Plugin](db_.NewCtx(s.ctx, s.svcCtx.Db))
	row, err := dao.GetById(id)
	if err != nil {
		return err
	}
	if row == nil {
		return fmt.Errorf("插件实例不存在")
	}
	if row.PluginType == string(plugin.TypeAccess) {
		var repositories []models.Plugin
		if err := s.svcCtx.Db.Where("plugin_type = ?", string(plugin.TypeRepository)).Find(&repositories).Error; err != nil {
			return err
		}
		for _, repository := range repositories {
			referenced, err := repositoryReferencesAccess(repository.ConfigYAML, id)
			if err != nil {
				return fmt.Errorf("检查仓库授权引用失败")
			}
			if referenced {
				return fmt.Errorf("授权仍被仓库引用")
			}
		}
	}
	_, err = dao.Delete(&id)
	return err
}
func (s *PluginService) present(row *models.Plugin) map[string]any {
	var fields []plugin.FieldSpec
	if definition, ok := s.svcCtx.Plugins.Get(row.PluginName); ok {
		fields = definition.Metadata.Fields
	}
	var decoded map[string]any
	if codec, err := s.secretCodec(); err == nil {
		decoded, _ = codec.DecodeYAML([]byte(row.ConfigYAML), secretFields(fields))
	}
	config := presentPluginConfig(row.ConfigYAML, fields, decoded)
	return map[string]any{"id": row.Id, "name": row.Name, "pluginType": row.PluginType, "pluginName": row.PluginName, "description": row.Description, "config": config, "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt}
}

func mergePluginConfig(existing, incoming map[string]any, fields []plugin.FieldSpec) map[string]any {
	merged := make(map[string]any, len(fields))
	for _, field := range fields {
		if value, ok := existing[field.Key]; ok {
			merged[field.Key] = value
		}
		value, ok := incoming[field.Key]
		if !ok || (field.Encrypt && value == "") {
			continue
		}
		merged[field.Key] = value
	}
	return merged
}

func validateRequiredPluginFields(config map[string]any, fields []plugin.FieldSpec) error {
	for _, field := range fields {
		if !field.Required {
			continue
		}
		value, ok := config[field.Key]
		if !ok || value == nil || strings.TrimSpace(fmt.Sprint(value)) == "" {
			return fmt.Errorf("%s不能为空", field.Title)
		}
	}
	return nil
}

func repositoryReferencesAccess(configYAML string, accessID int64) (bool, error) {
	var config map[string]any
	if err := yaml.Unmarshal([]byte(configYAML), &config); err != nil {
		return false, err
	}
	value, ok := config["accessId"]
	return ok && strings.TrimSpace(fmt.Sprint(value)) == fmt.Sprint(accessID), nil
}

func secretFields(fields []plugin.FieldSpec) secret.FieldMetadata {
	metadata := make(secret.FieldMetadata, len(fields))
	for _, field := range fields {
		metadata[field.Key] = secret.Field{Encrypt: field.Encrypt}
	}
	return metadata
}

func (s *PluginService) secretCodec() (*secret.Codec, error) {
	db := s.svcCtx.Db
	var setting models.SysSetting
	err := db.Where("key = ?", pluginMasterKeySetting).First(&setting).Error
	if err == nil {
		return codecFromSetting(setting.Setting)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	encoded := base64.RawStdEncoding.EncodeToString(key)
	if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoNothing: true}).Create(&models.SysSetting{Key: pluginMasterKeySetting, Setting: encoded}).Error; err != nil {
		return nil, err
	}
	if err := db.Where("key = ?", pluginMasterKeySetting).First(&setting).Error; err != nil {
		return nil, err
	}
	return codecFromSetting(setting.Setting)
}

func codecFromSetting(value string) (*secret.Codec, error) {
	key, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("插件主密钥无效")
	}
	codec, err := secret.NewCodec(key)
	if err != nil {
		return nil, fmt.Errorf("初始化插件加密器失败")
	}
	return codec, nil
}

func presentPluginConfig(configYAML string, fields []plugin.FieldSpec, decoded map[string]any) map[string]any {
	config := map[string]any{"configured": false}
	var values map[string]any
	if configYAML == "" || yaml.Unmarshal([]byte(configYAML), &values) != nil {
		return config
	}
	if decoded == nil {
		decoded = values
	}
	for _, field := range fields {
		value, exists := decoded[field.Key]
		if !exists || value == nil || value == "" {
			continue
		}
		config["configured"] = true
		if field.Encrypt {
			config[field.Key+"Configured"] = true
			if text, ok := value.(string); ok {
				masked := maskSecret(text)
				config[field.Key] = masked
				config[field.Key+"Original"] = masked
			}
			continue
		}
		config[field.Key] = value
	}
	return config
}

func maskSecret(value string) string {
	runes := []rune(value)
	if len(runes) <= 4 {
		return strings.Repeat("*", len(runes))
	}
	return string(runes[:2]) + strings.Repeat("*", len(runes)-4) + string(runes[len(runes)-2:])
}
