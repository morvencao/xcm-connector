package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type AccessTokenCfg struct {
	Auths map[string]RegistryAddress `json:"auths"`
}

type RegistryAddress struct {
	Auth  string `json:"auth"`
	Email string `json:"email"`
}

const (
	cloudAliasURL = "cloud.openshift.com"
)

var (
	pullSecretNamespace = "openshift-config"
	pullSecretName      = "pull-secret"
)

// retrieve access token from pull secret
func RetrieveAccessToken(ctx context.Context, kubeClient kubernetes.Interface) (string, error) {
	pullSecret, err := kubeClient.CoreV1().Secrets(pullSecretNamespace).Get(ctx, pullSecretName, metav1.GetOptions{})
	if err != nil {
		// fall back to find pull secret from current pod namespace
		if errors.IsNotFound(err) {
			pullSecretNamespace = os.Getenv("POD_NAMESPACE")
			if pullSecretNamespace == "" {
				return "", fmt.Errorf("POD_NAMESPACE must be set")
			}
			pullSecret, err = kubeClient.CoreV1().Secrets(pullSecretNamespace).Get(ctx, pullSecretName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	accessTokenCfg := AccessTokenCfg{}
	if err := json.Unmarshal(pullSecret.Data[".dockerconfigjson"], &accessTokenCfg); err != nil {
		return "", err
	}

	auths := accessTokenCfg.Auths
	if auths == nil {
		return "", fmt.Errorf("missing 'auths' entry in pull secret(%s/%s)", pullSecretNamespace, pullSecretName)
	}

	accessToken, exists := auths[cloudAliasURL]
	if !exists {
		return "", fmt.Errorf("missing auths.%s entry in pull secret(%s/%s)", cloudAliasURL, pullSecretNamespace, pullSecretName)
	}

	return accessToken.Auth, nil
}
