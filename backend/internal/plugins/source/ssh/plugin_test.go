package ssh

import (
	"context"
	"testing"
)

func TestProviderDefinitionDeclaresSSHSourceFieldsAndPathAction(t *testing.T) {
	definition, err := (Provider{}).Definition()
	if err != nil {
		t.Fatal(err)
	}
	if string(definition.Metadata.Type) != "source" || definition.Metadata.Name != "source.ssh" {
		t.Fatalf("metadata identity = %#v", definition.Metadata)
	}
	if len(definition.Metadata.Fields) != 3 || definition.Metadata.Fields[0].Key != "accessId" || !definition.Metadata.Fields[0].Required || definition.Metadata.Fields[1].Key != "paths" || !definition.Metadata.Fields[1].Required || definition.Metadata.Fields[2].Key != "excludePaths" || definition.Metadata.Fields[2].Required {
		t.Fatalf("source fields = %#v", definition.Metadata.Fields)
	}
	if definition.Metadata.Fields[1].Type != "array" || definition.Metadata.Fields[1].Component["name"] != "plugin-path-selector" || definition.Metadata.Fields[2].Type != "array" || definition.Metadata.Fields[2].Component["name"] != "plugin-path-selector" {
		t.Fatalf("paths field = %#v", definition.Metadata.Fields[1])
	}
	if len(definition.Metadata.Actions) != 1 || definition.Metadata.Actions[0].Name != "onListPaths" {
		t.Fatalf("source actions = %#v", definition.Metadata.Actions)
	}
}

func TestExecutorDispatchesActionWithCompletePluginForm(t *testing.T) {
	definition, err := (Provider{}).Definition()
	if err != nil {
		t.Fatal(err)
	}
	executor, err := definition.NewPluginInstance(map[string]any{"accessId": int64(7), "paths": []any{"/etc"}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := executor.ExecuteAction(context.Background(), "onListPaths", map[string]any{"cursor": "/"}); err == nil {
		t.Fatal("action accepted a form without resolved access config")
	}
	if _, err := executor.ExecuteAction(context.Background(), "onListPaths", map[string]any{"accessId": nil}); err == nil {
		t.Fatal("action accepted a form without accessId")
	}
}
