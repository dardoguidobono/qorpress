package admin

import (
	"github.com/gopress/internal/admin"
	qor_seo "github.com/gopress/internal/seo"

	"github.com/gopress/qorpress/pkg/models/posts"
	"github.com/gopress/qorpress/pkg/models/seo"
)

// SetupSEO add seo
func SetupSEO(Admin *admin.Admin) {
	seo.SEOCollection = qor_seo.New("Common SEO")
	seo.SEOCollection.RegisterGlobalVaribles(&seo.SEOGlobalSetting{SiteName: "Qor Shop"})
	seo.SEOCollection.SettingResource = Admin.AddResource(&seo.MySEOSetting{}, &admin.Config{Invisible: true})
	seo.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name: "Default Page",
	})
	seo.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name:     "Post Page",
		Varibles: []string{"Name", "Code", "CategoryName"},
		Context: func(objects ...interface{}) map[string]string {
			post := objects[0].(posts.Post)
			context := make(map[string]string)
			context["Name"] = post.Name
			context["Code"] = post.Code
			context["CategoryName"] = post.Category.Name
			return context
		},
	})
	seo.SEOCollection.RegisterSEO(&qor_seo.SEO{
		Name:     "Category Page",
		Varibles: []string{"Name", "Code"},
		Context: func(objects ...interface{}) map[string]string {
			category := objects[0].(posts.Category)
			context := make(map[string]string)
			context["Name"] = category.Name
			context["Code"] = category.Code
			return context
		},
	})
	Admin.AddResource(seo.SEOCollection, &admin.Config{Name: "SEO Setting", Menu: []string{"Site Management"}, Singleton: true, Priority: 2})
}
