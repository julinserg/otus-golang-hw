package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	dotDomain := "." + domain
	user := &User{}
	for scanner.Scan() {
		*user = User{}
		if err := json.Unmarshal(scanner.Bytes(), user); err != nil {
			return nil, fmt.Errorf("get users error: %w", err)
		}
		matched := strings.Contains(user.Email, dotDomain)
		if matched {
			indFindStr := strings.Index(user.Email, "@")
			str := strings.ToLower(user.Email[indFindStr+1 : len(user.Email)])
			result[str]++
		}
	}
	return result, scanner.Err()
}
