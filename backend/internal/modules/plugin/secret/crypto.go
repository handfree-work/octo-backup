package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

type Field struct{ Encrypt bool }
type FieldMetadata map[string]Field
type Codec struct{ aead cipher.AEAD }

func NewCodec(key []byte) (*Codec, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Codec{aead: aead}, nil
}
func (c *Codec) Encrypt(value string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	out := c.aead.Seal(nonce, nonce, []byte(value), nil)
	return "ENC[v1:" + base64.RawStdEncoding.EncodeToString(out) + "]", nil
}
func (c *Codec) Decrypt(value string) (string, error) {
	if !strings.HasPrefix(value, "ENC[v1:") || !strings.HasSuffix(value, "]") {
		return "", fmt.Errorf("密文格式无效")
	}
	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimSuffix(strings.TrimPrefix(value, "ENC[v1:"), "]"))
	if err != nil || len(raw) < c.aead.NonceSize() {
		return "", fmt.Errorf("密文无效")
	}
	plain, err := c.aead.Open(nil, raw[:c.aead.NonceSize()], raw[c.aead.NonceSize():], nil)
	if err != nil {
		return "", fmt.Errorf("解密失败")
	}
	return string(plain), nil
}
func (c *Codec) ProtectYAML(data []byte, fields FieldMetadata, configured map[string]bool) ([]byte, map[string]bool, error) {
	var values map[string]any
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, nil, err
	}
	if configured == nil {
		configured = map[string]bool{}
	}
	for key, spec := range fields {
		value, ok := values[key]
		if !ok || value == nil || value == "" {
			continue
		}
		if spec.Encrypt {
			configured[key] = true
			text, ok := value.(string)
			if !ok {
				return nil, nil, fmt.Errorf("敏感字段必须为字符串: %s", key)
			}
			encrypted, err := c.Encrypt(text)
			if err != nil {
				return nil, nil, err
			}
			values[key] = encrypted
		}
	}
	out, err := yaml.Marshal(values)
	return out, configured, err
}
func (c *Codec) DecodeYAML(data []byte, fields FieldMetadata) (map[string]any, error) {
	var values map[string]any
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, err
	}
	for key := range fields {
		value, ok := values[key].(string)
		if !ok || value == "" {
			continue
		}
		if strings.HasPrefix(value, "ENC[") {
			plain, err := c.Decrypt(value)
			if err != nil {
				return nil, err
			}
			values[key] = plain
		}
	}
	return values, nil
}
