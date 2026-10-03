package plugin_test

import (
	"context"
	"reflect"
	"testing"

	"handfree-work/octo-backup/internal/modules/plugin"
)

type testPlugin struct{}

func (testPlugin) ExecuteAction(_ context.Context, action string, _ map[string]any) (any, error) {
	if action == "onTest" {
		return map[string]any{"ok": true}, nil
	}
	return nil, plugin.ErrActionNotFound
}

func TestGenericBuildActionDoesNotReturnPluginConfig(t *testing.T) {
	definition, err := plugin.NewGenericDefinition([]byte("type: repository\nname: repository.example\ntitle: Example\nversion: 1.0.0\nactions:\n  - name: onBuild\n"))
	if err != nil {
		t.Fatal(err)
	}
	registry := plugin.NewRegistry()
	if err := registry.Register(definition); err != nil {
		t.Fatal(err)
	}
	got, err := registry.ExecuteWithConfig(context.Background(), "repository.example", "onBuild", map[string]any{"password": "secret"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(got, map[string]any{"password": "secret"}) {
		t.Fatalf("onBuild exposed plugin config: %#v", got)
	}
}

func TestRegistryParsesYAMLAndDispatchesDeclaredAction(t *testing.T) {
	definition, err := plugin.NewDefinition([]byte("type: access\nname: access.example\ntitle: Example\ndescription: Example access\nversion: 1.0.0\nfields:\n  - key: token\n    title: Token\n    type: password\n    required: true\n    encrypt: true\nactions:\n  - name: onTest\n    title: Test\n    permission: write\n"), func(map[string]any) (plugin.ActionExecutor, error) { return testPlugin{}, nil })
	if err != nil {
		t.Fatalf("NewDefinition() error = %v", err)
	}
	registry := plugin.NewRegistry()
	if err := registry.Register(definition); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	got, err := registry.Execute(context.Background(), "access.example", "onTest", map[string]any{})
	if err != nil || got == nil {
		t.Fatalf("Execute() = (%v, %v), want result", got, err)
	}
	if _, err := registry.Execute(context.Background(), "access.example", "OnTest", nil); err == nil {
		t.Fatal("Execute() accepted non-on action")
	}
	if _, err := registry.Execute(context.Background(), "access.example", "onDeleteEverything", nil); err == nil {
		t.Fatal("Execute() accepted undeclared action")
	}
}

func TestDefinitionRejectsInvalidActionName(t *testing.T) {
	_, err := plugin.NewDefinition([]byte("type: access\nname: access.example\ntitle: Example\ndescription: Example\nversion: 1.0.0\nactions:\n  - name: test\n"), func(map[string]any) (plugin.ActionExecutor, error) { return testPlugin{}, nil })
	if err == nil {
		t.Fatal("NewDefinition() accepted action without on prefix")
	}
}
