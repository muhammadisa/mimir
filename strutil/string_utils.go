package strutil

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"strings"
)

func Bearer(ctx context.Context) (string, error) {
	var tokenMetadata []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		tokenMetadata = md.Get("authorization")
	}
	if len(tokenMetadata) == 0 {
		return "", errors.New("no token given")
	}
	token := strings.Split(tokenMetadata[0], " ")
	if len(token) != 2 || token[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return token[1], nil
}

func TruncateCountryCode(str interface{}) string {
	trimmed := strings.TrimSpace(str.(string))
	return strings.ReplaceAll(trimmed, "+62", "0")
}
