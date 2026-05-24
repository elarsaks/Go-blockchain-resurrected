package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestSignatureStringRoundTrip(t *testing.T) {
	signature := &Signature{
		R: big.NewInt(12345),
		S: big.NewInt(67890),
	}

	encoded := signature.String()
	if len(encoded) != 128 {
		t.Fatalf("signature string length = %d, want 128", len(encoded))
	}

	decoded := SignatureFromString(encoded)
	if decoded.R.Cmp(signature.R) != 0 {
		t.Fatalf("R = %s, want %s", decoded.R, signature.R)
	}
	if decoded.S.Cmp(signature.S) != 0 {
		t.Fatalf("S = %s, want %s", decoded.S, signature.S)
	}
}

func TestPublicKeyFromStringRoundTrip(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey returned error: %v", err)
	}

	encoded := publicKeyString(&privateKey.PublicKey)
	decoded := PublicKeyFromString(encoded)

	if decoded.Curve != elliptic.P256() {
		t.Fatal("decoded public key uses unexpected curve")
	}
	if decoded.X.Cmp(privateKey.PublicKey.X) != 0 {
		t.Fatalf("X = %s, want %s", decoded.X, privateKey.PublicKey.X)
	}
	if decoded.Y.Cmp(privateKey.PublicKey.Y) != 0 {
		t.Fatalf("Y = %s, want %s", decoded.Y, privateKey.PublicKey.Y)
	}
}

func TestPrivateKeyFromStringRoundTrip(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey returned error: %v", err)
	}

	decoded := PrivateKeyFromString(privateKey.D.Text(16), &privateKey.PublicKey)

	if decoded.D.Cmp(privateKey.D) != 0 {
		t.Fatalf("D = %s, want %s", decoded.D, privateKey.D)
	}
	if decoded.PublicKey.X.Cmp(privateKey.PublicKey.X) != 0 {
		t.Fatalf("public key X = %s, want %s", decoded.PublicKey.X, privateKey.PublicKey.X)
	}
	if decoded.PublicKey.Y.Cmp(privateKey.PublicKey.Y) != 0 {
		t.Fatalf("public key Y = %s, want %s", decoded.PublicKey.Y, privateKey.PublicKey.Y)
	}
}

func TestPrivateKeyFromStringAcceptsOddLengthHex(t *testing.T) {
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(1),
		Y:     big.NewInt(2),
	}

	decoded := PrivateKeyFromString("abc", publicKey)

	if decoded.D.Cmp(big.NewInt(0xabc)) != 0 {
		t.Fatalf("D = %s, want %d", decoded.D, 0xabc)
	}
}

func publicKeyString(publicKey *ecdsa.PublicKey) string {
	return (&Signature{R: publicKey.X, S: publicKey.Y}).String()
}
