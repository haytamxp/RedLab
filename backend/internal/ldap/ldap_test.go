package ldap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"

	"github.com/haytamxp/redlab/backend/internal/config"
)

func loadLDAPTestConfig(t *testing.T) config.LDAPConfig {
	t.Helper()

	envPath := filepath.Join("..", "..", ".env")

	if err := godotenv.Overload(envPath); err != nil {
		t.Fatalf("failed to load backend .env: %v", err)
	}

	return config.LDAPConfig{
		Host:     os.Getenv("LDAP_HOST"),
		Port:     os.Getenv("LDAP_PORT"),
		Username: os.Getenv("LDAP_USERNAME"),
		Password: os.Getenv("LDAP_PASSWORD"),
		BaseDN:   os.Getenv("LDAP_BASE_DN"),
		UseTLS:   os.Getenv("LDAP_USE_TLS") == "true",
	}
}

func TestLDAPServiceAccountBind(t *testing.T) {
	cfg := loadLDAPTestConfig(t)

	if cfg.Host == "" ||
		cfg.Port == "" ||
		cfg.Username == "" ||
		cfg.Password == "" ||
		cfg.BaseDN == "" {
		t.Fatal("LDAP configuration is incomplete")
	}

	client := NewClient(cfg)

	defer client.Close()

	if err := client.BindServiceAccount(); err != nil {
		t.Fatalf("LDAP service account bind failed: %v", err)
	}

	t.Log("LDAP service account bind successful")
}

func TestLDAPFindUser(t *testing.T) {
	cfg := loadLDAPTestConfig(t)

	testUser := os.Getenv("LDAP_TEST_USER")

	if testUser == "" {
		t.Skip("LDAP_TEST_USER is not configured")
	}

	client := NewClient(cfg)

	defer client.Close()

	user, err := client.FindUser(testUser)
	if err != nil {
		t.Fatalf("LDAP user search failed: %v", err)
	}

	t.Log("User found")
	t.Logf("DN: %s", user.DN)
	t.Logf("SAMAccountName: %s", user.SAMAccountName)
	t.Logf("UPN: %s", user.UserPrincipalName)
	t.Logf("DisplayName: %s", user.DisplayName)
	t.Logf("Enabled: %t", user.Enabled)
}

func TestLDAPAuthenticate(t *testing.T) {
	cfg := loadLDAPTestConfig(t)

	testUser := os.Getenv("LDAP_TEST_USER")
	testPassword := os.Getenv("LDAP_TEST_PASSWORD")

	if testUser == "" || testPassword == "" {
		t.Skip("LDAP_TEST_USER or LDAP_TEST_PASSWORD is not configured")
	}

	client := NewClient(cfg)

	defer client.Close()

	user, err := client.Authenticate(testUser, testPassword)
	if err != nil {
		t.Fatalf("LDAP authentication failed: %v", err)
	}

	t.Logf(
		"LDAP authentication successful for %s",
		user.SAMAccountName,
	)
}
