package sftp

import "testing"

func TestProviderDefinitionUsesStorageRootWithoutRepositoryPassword(t *testing.T) {
	definition, err := (Provider{}).Definition()
	if err != nil {
		t.Fatal(err)
	}

	if len(definition.Metadata.Fields) != 2 {
		t.Fatalf("SFTP fields = %d, want 2", len(definition.Metadata.Fields))
	}
	if definition.Metadata.Fields[0].Key != "path" || definition.Metadata.Fields[0].Title != "存储根目录" {
		t.Fatalf("path field = %#v, want title 存储根目录", definition.Metadata.Fields[0])
	}
	for _, field := range definition.Metadata.Fields {
		if field.Key == "password" {
			t.Fatal("SFTP metadata must not expose a repository password field")
		}
	}
}
