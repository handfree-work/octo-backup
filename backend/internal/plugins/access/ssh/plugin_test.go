package ssh

import "testing"

func TestNewPluginInstanceMapsConfigToFields(t *testing.T) {
	definition, err := (Provider{}).Definition()
	if err != nil {
		t.Fatal(err)
	}
	executor, err := definition.NewPluginInstance(map[string]any{
		"host": "example.com", "port": 2200, "username": "root", "password": "secret",
		"privateKey": "key", "privateKeyPassword": "pass", "timeout": 45, "jumpAccessId": 7,
	})
	if err != nil {
		t.Fatal(err)
	}
	access := executor.(*SshAccess)
	if access.Host != "example.com" || access.Port != 2200 || access.Username != "root" || access.Timeout != 45 || access.JumpAccessId != 7 {
		t.Fatalf("config was not mapped: %#v", access)
	}
	if access.PrivateKey != "key" || access.PrivateKeyPassword != "pass" {
		t.Fatalf("private key config was not mapped")
	}
}
