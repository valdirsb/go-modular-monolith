package adapters

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// BCryptPasswordHasher implementa hashing de senha usando Argon2
type Argon2PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewArgon2PasswordHasher cria uma nova instância do hasher
func NewArgon2PasswordHasher() *Argon2PasswordHasher {
	return &Argon2PasswordHasher{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

// Hash cria um hash da senha
func (h *Argon2PasswordHasher) Hash(password string) (string, error) {
	salt, err := h.generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallelism, h.keyLength)

	// Formato: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%x$%x",
		argon2.Version, h.memory, h.iterations, h.parallelism, salt, hash)

	return encodedHash, nil
}

// Verify verifica se a senha corresponde ao hash
func (h *Argon2PasswordHasher) Verify(password, encodedHash string) bool {
	params, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	if len(hash) != len(otherHash) {
		return false
	}

	return subtle.ConstantTimeCompare(hash, otherHash) == 1
}

func (h *Argon2PasswordHasher) generateSalt() ([]byte, error) {
	salt := make([]byte, h.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func (h *Argon2PasswordHasher) decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid encoded hash format")
	}

	// Parse version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	// Parse parameters
	p = &params{}
	var memory, iterations, parallelism int
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return nil, nil, nil, err
	}

	p.memory = uint32(memory)
	p.iterations = uint32(iterations)
	p.parallelism = uint8(parallelism)
	p.keyLength = 32

	// Decode salt and hash from hex
	salt, err = hex.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err = hex.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	return p, salt, hash, nil
}
