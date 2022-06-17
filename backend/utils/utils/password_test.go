package utils

import "testing"

func TestGenerateHashedPassword(t *testing.T) {
	pwd, err := GenerateHashedPassword("123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pwd)
}

func TestComparePassword(t *testing.T) {
	ComparePassword("$2a$05$VKXlRuLEOsZSp83i1hkp0.xtk2/Qzykuwdtay5z5H7wY3vMarjOIS", "123456")

	cases := []struct {
		Name           string
		HashedPassword string
		Password       string
		Expected       bool
	}{
		{"case 0", "$2a$05$VKXlRuLEOsZSp83i1hkp0.xtk2/Qzykuwdtay5z5H7wY3vMarjOIS", "123456", true},
		{"case 1", "$2a$05$VKXlRuLEOsZSp83i1hkp0.xtk2/Qzykuwdtay5z5H7wY3vMarjOIS", "1234567", false},
		{"case 2", "$2a$05$VKXlRuLEOsZSp83i1hkp0.xtk2/Qzykuwdtay5z5H7wY3vMarjOIA", "123456", false},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := ComparePassword(c.HashedPassword, c.Password); ans != c.Expected {
				t.Fatalf("input hashedPassword[%v] password[%v] expected [%v], but ans[%v] got", c.HashedPassword, c.Password, c.Expected, ans)
			}
		})
	}
}
