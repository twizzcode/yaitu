package main

import "testing"

func TestPassword(t *testing.T) {
	password := "rahasia123"

	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() error: %v", err)
	}

	matched, err := verifyPassword(password, hash)
	if err != nil {
		t.Fatalf("verifyPassword() error: %v", err)
	}

	if !matched {
		t.Fatal("password benar dianggap salah")
	}

	matched, err = verifyPassword("password-salah", hash)
	if err != nil {
		t.Fatalf("verifyPassword() error: %v", err)
	}

	if matched {
		t.Fatal("password salah dianggap benar")
	}

	secondHash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("hashPassword() kedua error: %v", err)
	}

	if hash == secondHash {
		t.Fatal("password sama menghasilkan hash sama; salt seharusnya berbeda")
	}
}