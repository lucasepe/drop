package sha256_test

import (
	"testing"

	"github.com/lucasepe/drop/internal/crypto/sha256"
)

func TestGenerate(t *testing.T) {
	data := []struct {
		salt []byte
		key  []byte
		out  string
		cost int
	}{
		// openssl passwd -5 -salt 'saltstring' 'Hello world!'
		{
			[]byte("$5$saltstring"),
			[]byte("Hello world!"),
			"$5$saltstring$5B8vYYiY.CVt1RlTTf8KbXBH3hsxY/GNooZaBBGWEc5",
			sha256.RoundsDefault,
		},

		// openssl passwd -5 -salt 'rounds=10000$saltstringsaltstring' 'Hello world!'
		{
			[]byte("$5$rounds=10000$saltstringsaltstring"),
			[]byte("Hello world!"),
			"$5$rounds=10000$saltstringsaltst$3xv.VbSHBb41AL9AvLeujZkZRBAwqFMz2.opqey6IcA",
			10000,
		},

		// openssl passwd -5 -salt 'rounds=5000$toolongsaltstring' 'This is just a test'
		{
			[]byte("$5$rounds=5000$toolongsaltstring"),
			[]byte("This is just a test"),
			"$5$rounds=5000$toolongsaltstrin$Un/5jzAHMgOGZ5.mWJpuVolil07guHPvOW8mGRcvxa5",
			5000,
		},

		// openssl passwd -5 -salt 'rounds=1400$anotherlongsaltstring' 'A very much much much much longer text to encrypt.....'
		{
			[]byte("$5$rounds=1400$anotherlongsaltstring"),
			[]byte("A very much much much much longer text to encrypt....."),
			"$5$rounds=1400$anotherlongsalts$kbDNbTVPAcluuiah2ZMQxdYQUQ.LW.pnt/SKHxK934C",
			1400,
		},

		{
			[]byte("$5$rounds=77777$short"),
			[]byte("we have a short salt string but not a short password"),
			"$5$rounds=77777$short$JiO1O3ZpDAxGJeaDIuqCoEFysAe1mZNJRs3pw0KQRd/",
			77777,
		},
		{
			[]byte("$5$rounds=123456$asaltof16chars.."),
			[]byte("a short string"),
			"$5$rounds=123456$asaltof16chars..$gP3VQ/6X7UUEW3HkBn2w1/Ptq2jxPyzV/cZKmF/wJvD",
			123456,
		},
		{
			[]byte("$5$rounds=10$roundstoolow"),
			[]byte("the minimum number is still observed"),
			"$5$rounds=1000$roundstoolow$yfvwcWrQ8l/K0DAWyuPMDNHpIVlTQebY9l/gL972bIC",
			1000,
		},
	}

	sha256Crypt := sha256.New()

	for i, d := range data {
		hash, err := sha256Crypt.Generate(d.key, d.salt)
		if err != nil {
			t.Fatal(err)
		}
		if hash != d.out {
			t.Errorf("Test %d failed\nExpected: %s, got: %s", i, d.out, hash)
		}

		cost, err := sha256Crypt.Cost(hash)
		if err != nil {
			t.Fatal(err)
		}
		if cost != d.cost {
			t.Errorf("Test %d failed\nExpected: %d, got: %d", i, d.cost, cost)
		}
	}
}

func TestVerify(t *testing.T) {
	data := [][]byte{
		[]byte("password"),
		[]byte("12345"),
		[]byte("That's amazing! I've got the same combination on my luggage!"),
		[]byte("And change the combination on my luggage!"),
		[]byte("         random  spa  c    ing."),
		[]byte("94ajflkvjzpe8u3&*j1k513KLJ&*()"),
	}

	sha256Crypt := sha256.New()

	for i, d := range data {
		hash, err := sha256Crypt.Generate(d, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err = sha256Crypt.Verify(hash, d); err != nil {
			t.Errorf("Test %d failed: %s", i, d)
		}
	}
}
