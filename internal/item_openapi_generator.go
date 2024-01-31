package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ruapi-generate-md/pkg"
	"github.com/ruapi-generate-md/pkg/db"
	. "github.com/ruapi-generate-md/pkg/db/model"
	. "github.com/ruapi-generate-md/pkg/openapi"
)

// Generator OpenApi 3.0 json for specfiy project name
func GenerateOpenApiProjectName(outPath string, projectName string) {
	fmt.Printf("export start on %s", time.Now())

	// init db
	sqliteFilePath := "/showdoc_data/html/Sqlite/showdoc.db.php"
	dbInit(sqliteFilePath)

	// query project and global info from db
	item, err := data.Item.TakeItem(projectName)
	if err != nil {
		panic("not found this project:" + projectName)
	}
	// query project's catalogs
	catalogs, _ := data.Catalog.TakeAllCatalogs(item.ItemId)
	// query project's pages
	pages, _ := data.Page.TakeAllPages(item.ItemId)

	// export
	exportPath := outPath + "/" + projectName
	doExport(exportPath, *item, catalogs, pages)

	fmt.Printf("export finish on %s", time.Now())
}

func doExport(exportPath string, item Item, catalogs []*Catalog, pages []*Page) {
	// create export path if not exist
	if _, err := os.Stat(exportPath); os.IsNotExist(err) {
		if err := os.MkdirAll(exportPath, 0755); err != nil {
			fmt.Printf("mkdir %s failed：%s\n", exportPath, err)
			return
		}
	}

	// create openapi base spec
	openApiSpec := OpenApiSpec{
		OpenAPI: "3.0.1",
		Info: OpenApiInfo{
			Title:       item.ItemName.String,
			Description: item.ItemDescription.String,
			Version:     "1.0.0",
		},
	}

	//  create openapi tags
	if len(catalogs) > 0 {
		var openApiTags []OpenApiTag
		for _, catalog := range catalogs {
			openApiTags = append(openApiTags, OpenApiTag{
				Name:        catalog.CatName.String,
				Description: "",
			})
		}
		openApiSpec.Tags = openApiTags
	}

	// create openapi paths
	if len(pages) > 0 {
		openApiPaths := make(map[string]OpenApiPath)
		for _, page := range pages {
			pageCatalogName := getPageCatalogName(page, catalogs)
			pageContent, apiPath, err := parsePage2OpenAPiPath(page, pageCatalogName)
			if err == nil {
				url := getPathAndParamFromUrl(pageContent.Info.URL)
				openApiPaths[url] = *apiPath
				dumpShowdocData(exportPath+"/"+item.ItemName.String+"-showdoc.json", pageContent)
			}
		}
		openApiSpec.Paths = openApiPaths
	}

	jsonData, err := json.Marshal(openApiSpec)
	err = os.WriteFile(exportPath+"/"+item.ItemName.String+".json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func getPageCatalogName(page *Page, catalogs []*Catalog) string {
	for _, catalog := range catalogs {
		if catalog.CatId == page.CatId {
			return catalog.CatName.String
		}
	}
	return ""
}

func dbInit(sqliteFilePath string) {
	data = db.NewDataBase(sqliteFilePath)
	data.Init()
}

func parsePage2OpenAPiPath(page *Page, catalog string) (*pkg.PageContent, *OpenApiPath, error) {
	pageContentjson := strings.ReplaceAll(page.PageContent.String, "&quot;", "\"")
	pageContent := pkg.PageContent{}
	err := json.Unmarshal([]byte(pageContentjson), &pageContent)
	if err != nil {
		fmt.Printf("Could not parse page: %s\n", page.PageTitle.String)
		return &pkg.PageContent{}, &OpenApiPath{}, err
	}

	openApiOperation := OpenApiOperation{
		Summary:     page.PageTitle.String,
		Description: pageContent.Info.Description,
		Tags:        []string{catalog},
	}

	// convert  api request
	switch {
	// convert api requestBody
	case pageContent.Request.Params.Mode == "json":
		if len(pageContent.Request.Params.JSONDesc) > 0 {
			var openApiSchemaProperties = make(map[string]OpenApiSchemaProperty)
			for _, jsonDesc := range pageContent.Request.Params.JSONDesc {
				schemaProperty := OpenApiSchemaProperty{
					Title:       jsonDesc.Name,
					Type:        getOpenApiDataType(jsonDesc.Type),
					Description: jsonDesc.Remark,
					Required:    jsonDesc.Require == "1",
				}
				openApiSchemaProperties[jsonDesc.Name] = schemaProperty
			}
			openApiOperation.RequestBody = OpenApiRequestBody{
				Content: map[string]OpenApiMediaType{
					"application/json": {
						Example: pageContent.Request.Params.JSON,
						Schema: OpenApiSchema{
							Type:       "object",
							Properties: openApiSchemaProperties,
						},
					},
				},
			}
		} else {
			openApiOperation.RequestBody = OpenApiRequestBody{
				Content: map[string]OpenApiMediaType{
					"application/json": {
						Example: pageContent.Request.Params.JSON,
					},
				},
			}
		}

	// convert api parameters
	case pageContent.Request.Params.Mode == "urlencoded":
		openApiOperation.RequestBody = OpenApiRequestBody{
			Content: map[string]OpenApiMediaType{
				"application/json": {
					Example: pageContent.Request.Params.JSON,
				},
			},
		}
	case pageContent.Request.Params.Mode == "formdata":
		openApiOperation.RequestBody = OpenApiRequestBody{
			Content: map[string]OpenApiMediaType{
				"application/json": {
					Example: pageContent.Request.Params.JSON,
				},
			},
		}
	}

	// convert api responses
	responses := make(map[string]OpenApiResponse)
	responseStatus := strconv.Itoa(pageContent.Response.ResponseStatus)
	responses[responseStatus] = OpenApiResponse{
		Description: pageContent.Response.Remark,
		Content:     make(map[string]OpenApiMediaType),
	}
	// vaild response
	if pageContent.Response.ResponseStatus >= 0 && pageContent.Response.ResponseStatus <= 200 {
		var openApiSchemaProperties = make(map[string]OpenApiSchemaProperty)
		for _, jsonDesc := range pageContent.Response.ResponseParamsDesc {
			schemaProperty := OpenApiSchemaProperty{
				Title:       jsonDesc.Name,
				Type:        getOpenApiDataType(jsonDesc.Type),
				Description: jsonDesc.Remark,
			}
			openApiSchemaProperties[jsonDesc.Name] = schemaProperty
		}
		responses[responseStatus].Content["application/json"] = OpenApiMediaType{
			Example: pageContent.Response.ResponseExample,
			Schema: OpenApiSchema{
				Type:       "object",
				Properties: openApiSchemaProperties,
			},
		}
	}
	openApiOperation.Responses = responses

	// convert api path
	openApiPath := OpenApiPath{}
	method := pageContent.Info.Method
	switch {
	case method == "get":
		openApiPath.Get = &openApiOperation
	case method == "put":
		openApiPath.Put = &openApiOperation
	case method == "post":
		openApiPath.Post = &openApiOperation
	case method == "delete":
		openApiPath.Delete = &openApiOperation
	}

	return &pageContent, &openApiPath, nil
}

func dumpShowdocData(dumpPath string, pageContent *pkg.PageContent) {
	f, err := os.OpenFile(dumpPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	jsonData, err := json.Marshal(pageContent)
	if _, err = f.Write(jsonData); err != nil {
		panic(err)
	}
}

func getPathAndParamFromUrl(url string) string {
	newURL := strings.ReplaceAll(url, "https://", "")
	newURL = strings.ReplaceAll(newURL, "http://", "")

	// only need request url and paramter
	slashIndex := strings.Index(newURL, "/")
	if slashIndex != -1 {
		return newURL[slashIndex:]
	}
	return "/" + newURL
}

/* openapi data type -->  array, boolean, number, object, string, integer
 * showdoc data type --> array, boolean, number, object, string, int, long, date
 */
func getOpenApiDataType(showdocDataType string) string {
	switch {
	case showdocDataType == "array":
		return showdocDataType
	case showdocDataType == "boolean":
		return showdocDataType
	case showdocDataType == "number":
		return showdocDataType
	case showdocDataType == "object":
		return showdocDataType
	case showdocDataType == "string":
		return showdocDataType
	case showdocDataType == "int":
		return "number"
	case showdocDataType == "long":
		return "number"
	case showdocDataType == "date":
		return "object"
	case showdocDataType == "file":
		return "object"
	default:
		return "string"
	}
}
