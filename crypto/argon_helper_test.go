package crypto

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := PasswordConfig{}.DefaultConfig()

	if config.time != 1 {
		t.Errorf("Incorrect default time %d . Expected '1'", config.time)
	}

	if config.memory != 65536 {
		t.Errorf("Incorrect default memory %d . Expected '65536'", config.memory)
	}

	if config.threads != 4 {
		t.Errorf("Incorrect default memory %d . Expected '4'", config.threads)
	}

	if config.keyLen != 32 {
		t.Errorf("Incorrect default key length %d . Expected '32'", config.keyLen)
	}
}
