package internal

import (
	"encoding/json"
	"fmt"
	"github.com/ruapi-generate-md/pkg"
	"github.com/ruapi-generate-md/pkg/db"
	"strings"
	"sync"
)

func GeneratePageByItemID() {
	var wg sync.WaitGroup
	dbFileName := "/showdoc_data/html/Sqlite" + "/showdoc.db.php"
	//dbFileName := "./showdoc.db.php"
	data := db.NewDataBase(dbFileName)
	data.Init()
	mdPath := ""
	item, err := data.Item.TakeItem("OpenIM服务器API")
	if err != nil {
		panic("not found this project")
	}
	mdPath += "/" + item.ItemName.String + "/"
	headerArgs, err := data.RunapiGlobalParam.TakeRunapiGlobalHeaderParam(item.ItemId)
	if err != nil {
		panic("not found this header")
	}
	globalHeader := getGlobalHeader(headerArgs.ContentJsonStr)
	catalogs, _ := data.Catalog.TakeCatalogs(item.ItemId)
	for _, catalog := range catalogs {
		newCataLog := *catalog
		fmt.Println(newCataLog.CatName.String, newCataLog.CatId.Int32)
		wg.Add(1)
		go func() {
			tempPath := mdPath + newCataLog.CatName.String
			pages, _ := data.Page.TakePages(newCataLog.CatId, item.ItemId)
			for _, page := range pages {
				//fmt.Print(i, page.PageContent.String, globalHeader, page.PageTitle.String)
				generateOnePageMarkDown(page.PageContent.String, globalHeader, page.PageTitle.String, tempPath)
			}

			wg.Done()
		}()
	}
	wg.Wait()
}

func getGlobalHeader(jsonStr string) []pkg.Header {
	//jsonStr := "[{&quot;name&quot;:&quot;operationID&quot;,&quot;type&quot;:&quot;string&quot;,&quot;value&quot;:&quot;testGordon1111&quot;,&quot;require&quot;:&quot;1&quot;,&quot;remark&quot;:&quot;operationID用于全局链路追踪&quot;},{&quot;name&quot;:&quot;token&quot;,&quot;type&quot;:&quot;string&quot;,&quot;value&quot;:&quot;eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiJvcGVuSU1BZG1pbiIsIlBsYXRmb3JtIjoiIiwiZXhwIjoxNjg2NzEyMzkwLCJuYmYiOjE2Nzg5MzYwOTAsImlhdCI6MTY3ODkzNjM5MH0.Pr4NqWCVF54oSzFePAU0gFWDHrBGzyEnb_K1BymnC88&quot;,&quot;require&quot;:&quot;1&quot;,&quot;remark&quot;:&quot;管理员或者用户token&quot;}]"
	jsonStr = strings.ReplaceAll(jsonStr, "&quot;", "\"")
	data := []pkg.Header{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}
	return data

}