package sftp

import "testing"

func TestProviderDefinitionIncludesRepositoryPassword(t *testing.T) {
	definition, err := (Provider{}).Definition()
	if err != nil {
		t.Fatal(err)
	}

	if len(definition.Metadata.Fields) != 3 {
		t.Fatalf("SFTP fields = %d, want 3", len(definition.Metadata.Fields))
	}
	if definition.Metadata.Fields[0].Key != "path" || definition.Metadata.Fields[0].Title != "存储根目录" {
		t.Fatalf("path field = %#v, want title 存储根目录", definition.Metadata.Fields[0])
	}
	if definition.Metadata.Fields[2].Key != "password" || !definition.Metadata.Fields[2].Encrypt {
		t.Fatal("SFTP repository password must be encrypted")
	}
}
