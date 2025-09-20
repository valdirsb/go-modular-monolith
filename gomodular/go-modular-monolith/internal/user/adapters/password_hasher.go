package adapters

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"

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
	vals := make([]interface{}, 6)
	_, err = fmt.Sscanf(encodedHash, "$argon2id$v=%d$m=%d,t=%d,p=%d$%x$%x", &vals[0], &vals[1], &vals[2], &vals[3], &vals[4], &vals[5])
	if err != nil {
		return nil, nil, nil, err
	}

	var version int
	version = vals[0].(int)
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	p = &params{}
	p.memory = uint32(vals[1].(int))
	p.iterations = uint32(vals[2].(int))
	p.parallelism = uint8(vals[3].(int))
	p.keyLength = 32

	salt = vals[4].([]byte)
	hash = vals[5].([]byte)

	return p, salt, hash, nil
}
