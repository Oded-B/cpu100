package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"log"
	"strings"
)

var (
	alg        func()
	prvEcdsa   *ecdsa.PrivateKey
	prvEd25519 ed25519.PrivateKey
	msg        []byte
)

func DetectAlg() {
	var err error
	if prvEcdsa, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader); err != nil {
		panic(err)
	}
	if _, prvEd25519, err = ed25519.GenerateKey(rand.Reader); err != nil {
		panic(err)
	}

	switch strings.ToLower(*halg) {
	case "md5":
		log.Println("using MD5 hash algorithm.")
		alg = func() {
			md5.Sum(msg)
		}
	case "sha1", "sha-1":
		log.Println("using SHA-1 hash algorithm.")
		alg = func() {
			sha1.Sum(msg)
		}
	case "sha224", "sha-224":
		log.Println("using SHA224 hash algorithm.")
		alg = func() {
			sha256.Sum224(msg)
		}
	case "sha256", "sha-256":
		log.Println("using SHA256 hash algorithm.")
		alg = func() {
			sha256.Sum256(msg)
		}
	case "sha384", "sha-384":
		log.Println("using SHA-384 hash algorithm.")
		alg = func() {
			sha512.Sum384(msg)
		}
	case "sha512", "sha-512":
		log.Println("using SHA512 hash algorithm.")
		alg = func() {
			sha512.Sum512(msg)
		}
	case "sha512/224", "sha-512/224":
		log.Println("using SHA-512/224 hash algorithm.")
		alg = func() {
			sha512.Sum512_224(msg)
		}
	case "sha512/256", "sha-512/256":
		log.Println("using SHA-512/256 hash algorithm.")
		alg = func() {
			sha512.Sum512_256(msg)
		}
	case "ecdsa":
		log.Println("using ECSDA as defined in FIPS 186-3, calculated on elliptic curve P256, message hashed by SHA256.")
		alg = func() {
			var hash = sha256.Sum256(msg)
			if _, err := ecdsa.SignASN1(rand.Reader, prvEcdsa, hash[:]); err != nil {
				panic(err)
			}
		}
	case "ed25519":
		log.Println("using Ed25519 signature algorithm.")
		alg = func() {
			ed25519.Sign(prvEd25519, msg)
		}
	default:
		log.Fatal("given algorithm name does not pass")
	}
}
