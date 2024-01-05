package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {

	scanner := bufio.NewScanner(r)

	var i int
	for scanner.Scan() {
		var user User
		if err = easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}
		result[i] = user
		i++
	}

	if err = scanner.Err(); err != nil {
		return result, err
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		splitEmail := strings.Split(user.Email, "@")
		domainName := strings.ToLower(splitEmail[len(splitEmail)-1])

		if strings.HasSuffix(domainName, "."+domain) {
			result[domainName]++
		}
	}

	return result, nil
}
