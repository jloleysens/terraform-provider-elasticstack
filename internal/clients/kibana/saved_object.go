package kibana

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	importPath = "/s/{spaceId}/api/saved_objects/_import"
	exportPath = "/s/{spaceId}/api/saved_objects/_export"
)

type SavedObjectImportSuccessResponse struct {
	ID string `json:"id"`
}

// type SavedObjectImportRequest struct {
// 	Type       string                 `json:"type"`
// 	ID         string                 `json:"id"`
// 	Attributes map[string]interface{} `json:"attributes"`
// 	Overwrite  bool                   `json:"overwrite"`
// }

type SavedObjectImportRequestQueryParams struct {
	Overwrite bool `json:"overwrite"`
}

type SavedObjectSuccessImportResult struct {
	SuccessResults []SavedObjectImportSuccessResponse `json:"successResults"`
	SuccessCount   int32                              `json:"successCount"`
}

type SavedObjectGetRequestEntry struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type SavedObjectGetRequest struct {
	Objects               []SavedObjectGetRequestEntry `json:"objects"`
	IncludeReferencesDeep bool                         `json:"includeReferencesDeep"`
	ExcludeExportDetails  bool                         `json:"excludeExportDetails"`
}

type SavedObjectGetResult struct {
	Objects []interface{} `json:"objects"`
}

type SavedObjectGetResultEntry struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func getFullPath(basePath string, spaceAwarePath string, spaceID string) string {
	return basePath + strings.Replace(spaceAwarePath, "{spaceId}", spaceID, -1)
}

func buildRequest(apiClient *clients.ApiClient) (*resty.Request, error) {
	client, err := apiClient.GetKibanaClient()
	if err != nil {
		return nil, err
	}
	baseHeaders := make(map[string]string)
	baseHeaders["kbn-xsrf"] = "elasticstack-provider"
	baseHeaders["Content-Type"] = "application/json"
	baseHeaders["Accept"] = "*/*"
	req := client.Client.R().SetHeaders(baseHeaders)

	if auth, hasAuth := apiClient.GetAPIKeyAuth(); hasAuth {
		req.SetAuthToken(auth.Key)
		req.SetAuthScheme(apiAuthScheme)
	} else if auth, hasAuth := apiClient.GetBasicAuth(); hasAuth {
		req.SetBasicAuth(auth.UserName, auth.Password)
	} else {
		auth = clients.BasicAuth{}
		req.SetBasicAuth(auth.UserName, auth.Password)
	}

	return req, nil
}

func CreateSavedObject(apiClient *clients.ApiClient, savedObject models.SavedObject, objectType string) (*models.SavedObject, diag.Diagnostics) {
	req, err := buildRequest(apiClient)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	entry := savedObject.Attributes
	if savedObject.ID != "" {
		entry["id"] = savedObject.ID
	}
	jsonString, err := json.Marshal(entry)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	req.
		SetHeader("Content-Type", "multipart/form-data").
		SetQueryParam("overwrite", "true").
		SetBody([]byte(jsonString)).
		SetResult(&SavedObjectSuccessImportResult{})
	path := getFullPath(apiClient.ReadBasePath(), importPath, savedObject.SpaceID)
	resp, err := req.Post(path)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	result := resp.Result().(*SavedObjectSuccessImportResult)

	if result.SuccessCount == 0 {
		return nil, diag.FromErr(errors.New("Failed to create saved object"))
	}

	savedObject.ID = result.SuccessResults[0].ID

	return &savedObject, nil
}

func GetSavedObject(apiClient *clients.ApiClient, id string, spaceID string, objectType string) (*models.SavedObject, diag.Diagnostics) {
	req, err := buildRequest(apiClient)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	body := SavedObjectGetRequest{
		Objects:               []SavedObjectGetRequestEntry{{ID: id, Type: objectType}},
		IncludeReferencesDeep: false,
		ExcludeExportDetails:  true,
	}
	req.SetBody(body).SetResult(&SavedObjectGetResult{})
	path := getFullPath(apiClient.ReadBasePath(), exportPath, spaceID)
	resp, err := req.Post(path)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	result := resp.Result().(*SavedObjectGetResult)

	if len(result.Objects) == 0 {
		return nil, diag.FromErr(errors.New("Failed to create saved object"))
	}

	res := result.Objects[0].(map[string]interface{})

	return &models.SavedObject{
		ID:         res["id"].(string),
		SpaceID:    spaceID,
		Attributes: res,
	}, nil
}
