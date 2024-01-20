package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	ch := make(chan User)

	go func() {
		defer close(ch)

		err := getUsers(r, ch)
		if err != nil {
			return
		}
	}()

	domain = strings.ToLower(domain)

	for user := range ch {
		countDomains(user, domain, result)
	}

	return result, nil
}

func getUsers(r io.Reader, ch chan<- User) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user User
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return err
		}
		ch <- user
	}

	return scanner.Err()
}

func countDomains(u User, domain string, result DomainStat) {
	splitEmail := strings.Split(u.Email, "@")
	domainName := strings.ToLower(splitEmail[len(splitEmail)-1])

	if strings.HasSuffix(domainName, "."+domain) {
		result[domainName]++
	}
}
