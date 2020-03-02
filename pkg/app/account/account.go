package account

import (
	"strconv"

	"github.com/go-chi/chi"
	"github.com/qorpress/admin"
	"github.com/qorpress/qor"
	"github.com/qorpress/qor/resource"
	qorutils "github.com/qorpress/qor/utils"
	"github.com/qorpress/render"
	"github.com/qorpress/validations"
	"golang.org/x/crypto/bcrypt"

	"github.com/qorpress/qorpress-example/pkg/config/application"
	"github.com/qorpress/qorpress-example/pkg/config/auth"
	"github.com/qorpress/qorpress-example/pkg/models/users"
	"github.com/qorpress/qorpress-example/pkg/utils/funcmapmaker"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	controller := &Controller{View: render.New(&render.Config{
		AssetFileSystem: application.AssetFS.NameSpace("account"),
	}, "pkg/app/account/views")}

	funcmapmaker.AddFuncMapMaker(controller.View)
	app.ConfigureAdmin(application.Admin)

	application.Router.Mount("/auth/", auth.Auth.NewServeMux())

	application.Router.With(auth.Authority.Authorize()).Route("/account", func(r chi.Router) {
		r.Get("/", controller.Profile)
		r.Get("/profile", controller.Profile)
		r.Post("/profile", controller.Update)
	})
}

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "User Management", Priority: 3})
	user := Admin.AddResource(&users.User{}, &admin.Config{Menu: []string{"User Management"}})
	// user.Meta(&admin.Meta{Name: "Gender", Config: &admin.SelectOneConfig{Collection: []string{"Male", "Female", "Unknown"}}})
	// user.Meta(&admin.Meta{Name: "Birthday", Type: "date"})
	user.Meta(&admin.Meta{Name: "Role", Config: &admin.SelectOneConfig{Collection: []string{"Admin", "Maintainer", "Member"}}})
	user.Meta(&admin.Meta{Name: "Password",
		Type:   "password",
		Valuer: func(interface{}, *qor.Context) interface{} { return "" },
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			if newPassword := qorutils.ToString(metaValue.Value); newPassword != "" {
				bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				if err != nil {
					context.DB.AddError(validations.NewError(user, "Password", "Can't encrpt password"))
					return
				}
				u := resource.(*users.User)
				u.Password = string(bcryptPassword)
			}
		},
	})
	user.Meta(&admin.Meta{Name: "Confirmed", Valuer: func(user interface{}, ctx *qor.Context) interface{} {
		if user.(*users.User).ID == 0 {
			return true
		}
		return user.(*users.User).Confirmed
	}})
	// user.Meta(&admin.Meta{Name: "DefaultBillingAddress", Config: &admin.SelectOneConfig{Collection: userAddressesCollection}})
	// user.Meta(&admin.Meta{Name: "DefaultShippingAddress", Config: &admin.SelectOneConfig{Collection: userAddressesCollection}})

	user.Filter(&admin.Filter{
		Name: "Role",
		Config: &admin.SelectOneConfig{
			Collection: []string{"Admin", "Maintainer", "Member"},
		},
	})

	user.IndexAttrs("ID", "Email", "Name", "Role")
	user.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Email", "Password"},
				{"Avatar"},
				{"Role"},
				{"Confirmed"},
			},
		},
		&admin.Section{
			Title: "Credit Information",
			Rows: [][]string{
				{"Balance"},
			},
		},
		&admin.Section{
			Title: "Accepts",
			Rows: [][]string{
				{"AcceptPrivate", "AcceptLicense", "AcceptNews"},
			},
		},
	)
	user.EditAttrs(user.ShowAttrs())
}

func userAddressesCollection(resource interface{}, context *qor.Context) (results [][]string) {
	var (
		user users.User
		DB   = context.DB
	)

	DB.Preload("Addresses").Where(context.ResourceID).First(&user)

	for _, address := range user.Addresses {
		results = append(results, []string{strconv.Itoa(int(address.ID)), address.Stringify()})
	}
	return
}
