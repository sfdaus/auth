package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenerateSlugName(name string) string {
	// Lowercase
	slug := strings.ToLower(name)

	// Hapus karakter non-alfanumerik kecuali spasi
	reg := regexp.MustCompile(`[^a-z0-9\s]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace multiple spaces dengan single dash
	slug = strings.Join(strings.Fields(slug), "-")

	// Ambil 4 karakter pertama dari UUID untuk uniqueness
	uniqueSuffix := uuid.New().String()[:4]

	return fmt.Sprintf("%s-%s", slug, uniqueSuffix)
}
