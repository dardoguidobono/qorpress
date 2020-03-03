package account

import (
	"net/http"

	"github.com/qorpress/render"

	// "github.com/qorpress/qorpress-example/pkg/models/orders"
	// "github.com/qorpress/qorpress-example/pkg/models/users"
	"github.com/qorpress/qorpress-example/pkg/utils"
)

// Controller posts controller
type Controller struct {
	View *render.Render
}

// Profile profile show page
func (ctrl Controller) Profile(w http.ResponseWriter, req *http.Request) {
	var (
		currentUser = utils.GetCurrentUser(req)
		// tx                              = utils.GetDB(req)
	)

	ctrl.View.Execute("profile", map[string]interface{}{
		"CurrentUser": currentUser,
	}, req, w)
}

// Update update profile page
func (ctrl Controller) Update(w http.ResponseWriter, req *http.Request) {
	// FIXME
}
