package toast

import (
	"encoding/json"
	"log"
	"time"

	"nhn-toast/pkg/infrastructure/http"
)

const IdentifyEndpoint = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"

type PasswordCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenRequest struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	TenantId            string             `json:"tenantId"`
	PasswordCredentials PasswordCredential `json:"passwordCredentials"`
}

type TokenResponse struct {
	Access struct {
		Token struct {
			ID      string    `json:"id"`
			Expires time.Time `json:"expires"`
			Tenant  struct {
				ID            string `json:"id"`
				Name          string `json:"name"`
				GroupID       string `json:"groupId"`
				Description   string `json:"description"`
				Enabled       bool   `json:"enabled"`
				ProjectDomain string `json:"project_domain"`
			} `json:"tenant"`
			IssuedAt string `json:"issued_at"`
		} `json:"token"`
		ServiceCatalog []struct {
			Endpoints []struct {
				Region    string `json:"region"`
				PublicURL string `json:"publicURL"`
			} `json:"endpoints"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"serviceCatalog"`
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
			Roles    []struct {
				Name string `json:"name"`
			} `json:"roles"`
			RolesLinks []interface{} `json:"roles_links"`
		} `json:"user"`
		Metadata struct {
			Roles   []string `json:"roles"`
			IsAdmin int      `json:"is_admin"`
		} `json:"metadata"`
	} `json:"access"`
}

func GenerateToken(tenantId, userName, apiPassword string) (string, error) {
	tokenRequest := new(TokenRequest)
	tokenRequest.Auth.TenantId = tenantId
	tokenRequest.Auth.PasswordCredentials.Username = userName
	tokenRequest.Auth.PasswordCredentials.Password = apiPassword

	tokenReqJson, err := json.MarshalIndent(&tokenRequest, "", "  ")
	if err != nil {
		log.Println("Token request json marshal err")
		return "", err
	}

	httpReqHeader := map[string]string{}
	httpReqHeader["Content-Type"] = "application/json"
	jsonRes, err := http.Post(IdentifyEndpoint, tokenReqJson, httpReqHeader)
	if err != nil {
		return "", err
	}

	tokenResponse := new(TokenResponse)
	err = json.Unmarshal(jsonRes, &tokenResponse)
	if err != nil {
		log.Println("Token response json unmarshal err")
		return "", err
	}

	return tokenResponse.Access.Token.ID, nil
}
