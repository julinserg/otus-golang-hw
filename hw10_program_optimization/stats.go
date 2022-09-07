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
	u, num, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	result := make(DomainStat, num)
	err = countDomains(&u, num, &result, domain)
	return result, err
}

type users [100_000]User

func getUsers(r io.Reader) (result users, numUsers int, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	scanner := bufio.NewScanner(r)

	user := &User{}
	for scanner.Scan() {
		content := scanner.Text()
		*user = User{}
		if err = json.Unmarshal([]byte(content), user); err != nil {
			continue
		}
		result[numUsers] = *user
		numUsers++
	}
	return
}

func countDomains(u *users, numUsers int, result *DomainStat, domain string) error {
	dotDomain := "." + domain
	for idx := 0; idx < numUsers; idx++ {
		matched := strings.Contains(u[idx].Email, dotDomain)
		if matched {
			indFindStr := strings.Index(u[idx].Email, "@")
			str := strings.ToLower(u[idx].Email[indFindStr+1 : len(u[idx].Email)])
			(*result)[str] += 1
		}
	}
	return nil
}
