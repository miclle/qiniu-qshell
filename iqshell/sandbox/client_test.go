package sandbox

import (
	"os"
	"testing"

	"github.com/qiniu/go-sdk/v7/sandbox"
	"github.com/stretchr/testify/assert"
)

// clearEnvVars unsets all sandbox-related env vars and returns a cleanup function.
func clearEnvVars(t *testing.T) {
	t.Helper()
	envs := []string{
		EnvQiniuAPIKey, EnvE2BAPIKey,
		EnvQiniuSandboxAPIURL, EnvE2BAPIURL,
		EnvQiniuAccessKey, EnvQiniuSecretKey,
	}
	saved := make(map[string]string)
	for _, k := range envs {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for _, k := range envs {
			if v, ok := saved[k]; ok {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	})
}

// stubWorkspaceAccount overrides workspaceAccountLookup for the duration of a
// test. Pass empty strings to simulate an absent qshell account.
func stubWorkspaceAccount(t *testing.T, ak, sk string) {
	t.Helper()
	saved := workspaceAccountLookup
	workspaceAccountLookup = func() (string, string) {
		return ak, sk
	}
	t.Cleanup(func() {
		workspaceAccountLookup = saved
	})
}

func TestResolveConfig_QiniuPriority(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "qiniu-key")
	os.Setenv(EnvE2BAPIKey, "e2b-key")
	os.Setenv(EnvQiniuSandboxAPIURL, "https://qiniu.example.com")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "qiniu-key", apiKey)
	assert.Equal(t, "https://qiniu.example.com", endpoint)
}

func TestResolveConfig_FallbackToE2B(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvE2BAPIKey, "e2b-key")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "e2b-key", apiKey)
	assert.Equal(t, "https://e2b.example.com", endpoint)
}

func TestResolveConfig_DefaultEndpoint(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "some-key")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "some-key", apiKey)
	assert.Equal(t, sandbox.DefaultEndpoint, endpoint)
}

func TestResolveConfig_AllEmpty(t *testing.T) {
	clearEnvVars(t)

	apiKey, endpoint := resolveConfig()
	assert.Empty(t, apiKey)
	assert.Equal(t, sandbox.DefaultEndpoint, endpoint)
}

func TestResolveConfig_QiniuKeyWithE2BEndpoint(t *testing.T) {
	clearEnvVars(t)
	os.Setenv(EnvQiniuAPIKey, "qiniu-key")
	os.Setenv(EnvE2BAPIURL, "https://e2b.example.com")

	apiKey, endpoint := resolveConfig()
	assert.Equal(t, "qiniu-key", apiKey)
	assert.Equal(t, "https://e2b.example.com", endpoint)
}

func TestResolveCredentials_WorkspaceTakesPriority(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "ws-ak", "ws-sk")

	// Both sources have AK/SK; workspace must win.
	os.Setenv(EnvQiniuAccessKey, "env-ak")
	os.Setenv(EnvQiniuSecretKey, "env-sk")

	got := resolveCredentials()
	if assert.NotNil(t, got) {
		assert.Equal(t, "ws-ak", got.AccessKey)
		assert.Equal(t, "ws-sk", string(got.SecretKey))
	}
}

func TestResolveCredentials_FallbackToEnv(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "", "")

	os.Setenv(EnvQiniuAccessKey, "env-ak")
	os.Setenv(EnvQiniuSecretKey, "env-sk")

	got := resolveCredentials()
	if assert.NotNil(t, got) {
		assert.Equal(t, "env-ak", got.AccessKey)
		assert.Equal(t, "env-sk", string(got.SecretKey))
	}
}

func TestResolveCredentials_PartialEnvIsIgnored(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "", "")

	// AK without SK should not produce a credential.
	os.Setenv(EnvQiniuAccessKey, "env-ak")

	assert.Nil(t, resolveCredentials())
}

func TestResolveCredentials_PartialWorkspaceIsIgnored(t *testing.T) {
	clearEnvVars(t)
	// Workspace has AK but empty SK — fall through to env (which is also empty).
	stubWorkspaceAccount(t, "ws-ak", "")

	assert.Nil(t, resolveCredentials())
}

func TestResolveCredentials_NoneConfigured(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "", "")

	assert.Nil(t, resolveCredentials())
}

func TestNewSandboxClient_RequiresAPIKey(t *testing.T) {
	clearEnvVars(t)
	// Only AK/SK is configured — sandbox runtime endpoints still need API Key.
	stubWorkspaceAccount(t, "ws-ak", "ws-sk")

	c, err := NewSandboxClient()
	assert.Nil(t, c)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), EnvQiniuAPIKey)
	}
}

func TestNewSandboxClient_OKWithAPIKeyOnly(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "", "")
	os.Setenv(EnvQiniuAPIKey, "k")

	c, err := NewSandboxClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewInjectionRuleClient_RequiresCredentials(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "", "")

	// Only API Key is configured — injection rule endpoints need AK/SK.
	os.Setenv(EnvQiniuAPIKey, "k")

	c, err := NewInjectionRuleClient()
	assert.Nil(t, c)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), EnvQiniuAccessKey)
	}
}

func TestNewInjectionRuleClient_OKWithCredentialsOnly(t *testing.T) {
	clearEnvVars(t)
	stubWorkspaceAccount(t, "ws-ak", "ws-sk")

	c, err := NewInjectionRuleClient()
	assert.NoError(t, err)
	assert.NotNil(t, c)
}
