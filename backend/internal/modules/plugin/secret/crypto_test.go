package secret_test

import (
	"strings"
	"testing"

	"handfree-work/octo-backup/internal/modules/plugin/secret"
)

func TestEncryptYAMLFieldsOnlyEncryptsMarkedFields(t *testing.T) {
	codec, err := secret.NewCodec(make([]byte, 32))
	if err != nil {
		t.Fatal(err)
	}
	metadata := secret.FieldMetadata{"token": {Encrypt: true}, "region": {}}
	encoded, configured, err := codec.ProtectYAML([]byte("token: plain-secret\nregion: cn\n"), metadata, nil)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "plain-secret") || !strings.Contains(string(encoded), "ENC[v1:") {
		t.Fatalf("protected YAML = %q, expected encrypted token", encoded)
	}
	if configured["token"] != true || configured["region"] != false {
		t.Fatalf("configured fields = %#v", configured)
	}
	decoded, err := codec.DecodeYAML(encoded, metadata)
	if err != nil {
		t.Fatal(err)
	}
	if decoded["token"] != "plain-secret" || decoded["region"] != "cn" {
		t.Fatalf("decoded = %#v", decoded)
	}
}

func TestDecryptRejectsDamagedCiphertext(t *testing.T) {
	codec, _ := secret.NewCodec(make([]byte, 32))
	if _, err := codec.Decrypt("ENC[v1:broken]"); err == nil {
		t.Fatal("Decrypt() accepted damaged ciphertext")
	}
}
