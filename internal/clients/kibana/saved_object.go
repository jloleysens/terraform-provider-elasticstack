package kibana

import (
	"encoding/json"
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
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
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
	attributesString := savedObject.Attributes
	// ==> Decode to go object
	attributes := map[string]interface{}{}
	if err := json.NewDecoder(strings.NewReader(attributesString)).Decode(&attributes); err != nil {
		return nil, diag.FromErr(err)
	}
	if savedObject.ID != "" {
		attributes["id"] = savedObject.ID
	}
	// <== Marshall to JSON bytes
	jsonBytes, err := json.Marshal(attributes)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	jsonString := strings.Replace(string(jsonBytes), "\n", "", -1)
	path := getFullPath(apiClient.ReadBasePath(), importPath, savedObject.SpaceID)
	resp, err := req.
		SetHeader("Content-Type", "multipart/form-data").
		SetQueryParam("overwrite", "true").
		SetFileReader("file", "import.ndjson", strings.NewReader(jsonString)).
		SetResult(&SavedObjectSuccessImportResult{}).
		Post(path)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	result := resp.Result().(*SavedObjectSuccessImportResult)

	if result.SuccessCount == 0 {
		return nil, diag.Errorf("Failed to create saved object")
	}

	savedObject.ID = result.SuccessResults[0].ID

	return &savedObject, nil
}

type GetResult string

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

	path := getFullPath(apiClient.ReadBasePath(), exportPath, spaceID)
	resp, err := req.
		SetBody(body).
		Post(path)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	result := SavedObjectGetResult{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, diag.FromErr(err)
	}
	if result.ID == "" {
		return nil, diag.Errorf(`Could not find dashboard with ID "%s", result: %s`, id, result)
	}

	b, err := json.Marshal(result.Attributes)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return &models.SavedObject{
		ID:         result.ID,
		SpaceID:    spaceID,
		Attributes: string(b),
	}, nil
}
