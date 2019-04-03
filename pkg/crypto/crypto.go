/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"unicode"
)

const (
	AesKeyLength        = 32
	AnnotationBase64    = "@base64:"
	AnnotationEncrypted = "@encrypted:"
	HeaderLength        = 2
)

func Encrypt(plainBytes []byte, pubKey *rsa.PublicKey) ([]byte, error) {

	randKey := make([]byte, AesKeyLength)
	if _, err := rand.Read(randKey); err != nil {
		return nil, err
	}

	rsaCipherBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, randKey, nil)
	if err != nil {
		return nil, err
	}

	cipherBytes := make([]byte, HeaderLength)
	binary.BigEndian.PutUint16(cipherBytes, uint16(len(rsaCipherBytes)))
	cipherBytes = append(cipherBytes, rsaCipherBytes...)

	blockCipher, err := aes.NewCipher(randKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	zeroNonce := make([]byte, aesgcm.NonceSize())
	cipherBytes = aesgcm.Seal(cipherBytes, zeroNonce, plainBytes, nil)
	return cipherBytes, nil
}

func Decrypt(cipherBytes []byte, key *rsa.PrivateKey) ([]byte, error) {
	if len(cipherBytes) < HeaderLength {
		return nil, fmt.Errorf("cipher data does not contain rsa length information")
	}
	rsaCipherLength := int(binary.BigEndian.Uint16(cipherBytes))
	if len(cipherBytes) < rsaCipherLength+HeaderLength {
		return nil, fmt.Errorf("invalid rsa cipher length")
	}

	rsaCipherBytes := cipherBytes[HeaderLength : rsaCipherLength+HeaderLength]
	aesCipherBytes := cipherBytes[rsaCipherLength+HeaderLength:]

	aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, rsaCipherBytes, nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	zeroNonce := make([]byte, aesgcm.NonceSize())

	plaintext, err := aesgcm.Open(nil, zeroNonce, aesCipherBytes, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func TryDecrypt(payload string, key *rsa.PrivateKey) ([]byte, error) {
	s := strings.TrimSpace(payload)

	if strings.HasPrefix(s, "@") {
		s = removeSpace(s)
		if strings.HasPrefix(s, AnnotationBase64) {
			return base64.StdEncoding.DecodeString(strings.TrimPrefix(s, AnnotationBase64))
		} else if strings.HasPrefix(s, AnnotationEncrypted) {
			encryptedBase64 := strings.TrimPrefix(s, AnnotationEncrypted)
			cipherBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
			if err != nil {
				return nil, err
			}
			return Decrypt(cipherBytes, key)
		} else {
			return nil, fmt.Errorf("unknown annotation in %s", payload)
		}
	} else {
		return []byte(payload), nil
	}

}

func removeSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
