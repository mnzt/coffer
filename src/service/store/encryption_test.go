package store

import "testing"

var (
	key        = "something that satisifies 32...."
	cipherText string
)

func TestEncrypt(t *testing.T) {
	crypter, err := InitaliaseCrypter(key)
	if err != nil {
		t.Error(err)
	}
	// Test for 32 byte strings
	c, err := crypter.Encrypt("something that satisifies 32....")
	if err != nil {
		t.Error(err)
	}
	t.Logf("cipher text: %s", c)

	// Test for sub-32 bytes with padding
	c1, err := crypter.Encrypt("this is odd")
	if err != nil {
		t.Error(err)
	}
	t.Logf("cipher text: %s", c1)

	/*
	 *	Test the key generation
	 */

	nokeyCrypter, err := InitaliaseCrypter("")
	if err != nil {
		t.Error(err)
	}

	_, err = nokeyCrypter.Encrypt("test")
	if err != nil {
		t.Error(err)
	}

	// Push out cipher to a varible to decrypt
	cipherText = c
}

func TestDecrypt(t *testing.T) {
	crypter, err := InitaliaseCrypter(key)
	if err != nil {
		t.Error(err)
	}

	plaintext, err := crypter.Decrypt(cipherText)
	if err != nil {
		t.Error(err)
	}

	expected := "something that satisifies 32...."
	if plaintext != expected {
		t.Errorf("error expected %s got %s", expected, plaintext)
	}

}
