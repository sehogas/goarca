package middleware

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"net/http"

	"github.com/sehogas/goarca/internal/dto"
	"github.com/sehogas/goarca/internal/util"
)

var ErrNoApiKey = errors.New("no api Key defined")
var ErrInvalidApiKey = errors.New("invalid api Key")

type ApiKeyMiddleware struct {
	keys []string
}

func NewApiKeyMiddleware(filename string) (*ApiKeyMiddleware, error) {
	keys, err := getKeysFromFile(filename)
	if err != nil {
		return nil, err
	}

	return &ApiKeyMiddleware{
		keys: keys,
	}, nil
}

func (m *ApiKeyMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Salteo control porque swagger es público
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}

		authorization := r.Header.Get("x-api-key")
		if err := m.CheckAPIKey(authorization); err != nil {
			util.HttpResponseJSON(w, http.StatusInternalServerError, &dto.ErrorResponse{Error: "invalid api key"}, errors.New("invalid api key"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getKeysFromFile(fileName string) ([]string, error) {
	var keys []string
	file, err := os.Open(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return keys, nil
		}
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var tmp string
	for scanner.Scan() {
		tmp = strings.TrimSpace(scanner.Text())
		pair := strings.SplitN(tmp, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			if strings.HasPrefix(key, "API_KEY_") {
				keys = append(keys, value)
			}
		}
	}
	return keys, nil
}

func keyInKeys(key string, keys []string) bool {
	for _, k := range keys {
		if key == k {
			return true
		}
	}
	return false
}

func (m *ApiKeyMiddleware) CheckAPIKey(key string) error {
	if key == "" {
		return ErrNoApiKey
	}
	if !keyInKeys(key, m.keys) {
		return ErrInvalidApiKey
	}
	return nil
}
