package password

import "testing"

func TestGeneratePasswordHash(t *testing.T) {
	hash, err := GeneratePasswordHash("123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}
func TestValidatePassword(t *testing.T) {
	hashedPassword := "$2a$10$R3E79T9uHwigAMvbRACiSO4cy5zROOQPQVvMlDQg3ciqbB81wk63K"
	err := ValidatePassword(hashedPassword, "123456")
	if err != nil {
		t.Fatal(err)
	}
}
