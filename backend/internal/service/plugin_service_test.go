package logic

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"handfree-work/octo-backup/internal/base/db_"
	"handfree-work/octo-backup/internal/models"
	"handfree-work/octo-backup/internal/modules/plugin"
	"handfree-work/octo-backup/internal/svc"
)

func TestMergePluginConfigPreservesBlankSecretsAndUpdatesOtherFields(t *testing.T) {
	existing := map[string]any{"secretId": "old-id", "secretKey": "old-key", "region": "ap-shanghai"}
	incoming := map[string]any{"secretId": "new-id", "secretKey": "", "region": ""}
	fields := []plugin.FieldSpec{
		{Key: "secretId", Encrypt: true},
		{Key: "secretKey", Encrypt: true},
		{Key: "region"},
	}

	got := mergePluginConfig(existing, incoming, fields)
	want := map[string]any{"secretId": "new-id", "secretKey": "old-key", "region": ""}
	for key, value := range want {
		if got[key] != value {
			t.Errorf("config[%q] = %#v, want %#v", key, got[key], value)
		}
	}
}

func TestMergePluginConfigDropsConfiguredMarkerAndUnknownFields(t *testing.T) {
	got := mergePluginConfig(map[string]any{"region": "cn"}, map[string]any{"configured": true, "unexpected": "value"}, []plugin.FieldSpec{{Key: "region"}})
	if len(got) != 1 || got["region"] != "cn" {
		t.Fatalf("merged config = %#v, want only preserved region", got)
	}
}

func TestRepositoryReferencesAccessMatchesExactID(t *testing.T) {
	referenced, err := repositoryReferencesAccess("accessId: 42\npath: backup\n", 42)
	if err != nil || !referenced {
		t.Fatalf("reference = %v, error = %v; want referenced", referenced, err)
	}

	referenced, err = repositoryReferencesAccess("accessId: 142\n", 42)
	if err != nil || referenced {
		t.Fatalf("reference = %v, error = %v; want not referenced", referenced, err)
	}
}

func TestRepositoryReferencesAccessRejectsMalformedConfig(t *testing.T) {
	if _, err := repositoryReferencesAccess("accessId: [", 42); err == nil {
		t.Fatal("malformed config should return an error")
	}
}

func TestPresentPluginConfigRedactsSecretsAndReturnsPlainFields(t *testing.T) {
	config := presentPluginConfig("secretId: id-1\nsecretKey: key-1\nregion: ap-shanghai\n", []plugin.FieldSpec{
		{Key: "secretId", Encrypt: true},
		{Key: "secretKey", Encrypt: true},
		{Key: "region"},
	}, nil)
	if config["secretId"] != "****" || config["secretKey"] != "ke*-1" {
		t.Fatalf("secret values should be masked: %#v", config)
	}
	if config["secretIdConfigured"] != true || config["secretKeyConfigured"] != true || config["region"] != "ap-shanghai" {
		t.Fatalf("present config = %#v", config)
	}
}

func TestPluginCreateUpdateAndConcurrentMasterKeyInitialization(t *testing.T) {
	database, err := db_.OpenSQLite(filepath.Join(t.TempDir(), "plugins.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db_.Migrate(database, &models.Plugin{}, &models.SysSetting{}); err != nil {
		t.Fatal(err)
	}

	registry := plugin.NewRegistry()
	definition, err := plugin.NewGenericDefinition([]byte("type: access\nname: access.test\ntitle: Test\nversion: 1\nfields:\n  - {key: secret, title: Secret, type: password, encrypt: true, required: true}\n  - {key: account, title: Account, type: string, required: true}\n  - {key: region, title: Region, type: string}\nactions:\n  - {name: onTest}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(definition); err != nil {
		t.Fatal(err)
	}
	repositoryDefinition, err := plugin.NewGenericDefinition([]byte("type: repository\nname: repository.test\ntitle: Repository\nversion: 1\nfields:\n  - {key: accessId, title: Access ID, type: number}\n  - {key: path, title: Path, type: string, required: true}\nactions:\n  - {name: onBuild}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(repositoryDefinition); err != nil {
		t.Fatal(err)
	}
	service := NewPluginService(context.Background(), &svc.ServiceContext{Db: database, Plugins: registry})

	const workers = 8
	start := make(chan struct{})
	var wait sync.WaitGroup
	errCh := make(chan error, workers)
	for range workers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			<-start
			_, err := service.secretCodec()
			errCh <- err
		}()
	}
	close(start)
	wait.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent secretCodec() error = %v", err)
		}
	}

	created, err := service.Create(&PluginInput{
		Name:       "腾讯云",
		PluginType: "access",
		PluginName: "access.test",
		Config:     map[string]any{"secret": "top-secret", "account": "account-1", "region": "cn-shanghai"},
	})
	if err != nil {
		t.Fatal(err)
	}
	idPtr, ok := created["id"].(*int64)
	if !ok || idPtr == nil {
		t.Fatalf("created plugin id = %#v, want *int64", created["id"])
	}
	id := *idPtr
	var stored models.Plugin
	if err := database.First(&stored, id).Error; err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stored.ConfigYAML, "top-secret") {
		t.Fatalf("stored config contains plaintext secret: %q", stored.ConfigYAML)
	}
	if createdConfig := created["config"].(map[string]any); createdConfig["secret"] != "to******et" || createdConfig["secretConfigured"] != true || createdConfig["account"] != "account-1" {
		t.Fatalf("created config should be masked while returning plain fields: %#v", createdConfig)
	}
	if _, err := service.Update(id, &PluginInput{Config: map[string]any{"account": ""}}); err == nil {
		t.Fatal("update should reject clearing a required field")
	}

	if _, err := service.Update(id, &PluginInput{Name: "腾讯云修改", Config: map[string]any{"secret": "", "region": ""}}); err != nil {
		t.Fatal(err)
	}
	if err := database.First(&stored, id).Error; err != nil {
		t.Fatal(err)
	}
	codec, err := service.secretCodec()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := codec.DecodeYAML([]byte(stored.ConfigYAML), secretFields(definition.Metadata.Fields))
	if err != nil {
		t.Fatal(err)
	}
	if decoded["secret"] != "top-secret" || decoded["region"] != "" {
		t.Fatalf("updated config = %#v", decoded)
	}
	if _, err := service.Action(id, "onTest", nil); err != nil {
		t.Fatalf("action should load decrypted config: %v", err)
	}

	page, err := service.Page(&PluginPageQuery{Offset: -10, Limit: 10, PluginType: "access", Name: "修改"})
	if err != nil || page.Offset != 0 || page.Total != 1 || len(page.Records) != 1 {
		t.Fatalf("filtered page = %#v, error = %v", page, err)
	}

	repository, err := service.Create(&PluginInput{Name: "repo", PluginName: "repository.test", Config: map[string]any{"accessId": id, "path": "backup"}})
	if err != nil {
		t.Fatal(err)
	}
	repositoryID := *repository["id"].(*int64)
	if err := service.Delete(id); err == nil {
		t.Fatal("delete should reject an access instance referenced by a repository")
	}
	if err := service.Delete(repositoryID); err != nil {
		t.Fatalf("delete repository: %v", err)
	}
	if err := service.Delete(id); err != nil {
		t.Fatalf("delete unreferenced access: %v", err)
	}
}
