package admin

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/qorpress/admin"
	"github.com/qorpress/banner_editor"
	"github.com/qorpress/l10n"
	"github.com/qorpress/media/oss"
	"github.com/qorpress/qor"
	"github.com/qorpress/qor/resource"
	"github.com/qorpress/sorting"
	"github.com/qorpress/widget"

	"github.com/qorpress/qorpress-example/pkg/config/db"
	"github.com/qorpress/qorpress-example/pkg/models/posts"
)

var Widgets *widget.Widgets

type QorWidgetSetting struct {
	widget.QorWidgetSetting
	// publish2.Version
	// publish2.Schedule
	// publish2.Visible
	l10n.Locale
}

// SetupWidget setup widget
func SetupWidget(Admin *admin.Admin) {
	if Widgets != nil {
		return
	}
	Widgets = widget.New(&widget.Config{DB: db.DB})
	Widgets.WidgetSettingResource = Admin.NewResource(&QorWidgetSetting{})

	Widgets.RegisterScope(&widget.Scope{
		Name: "From Google",
		Visible: func(context *widget.Context) bool {
			if request, ok := context.Get("Request"); ok {
				_, ok := request.(*http.Request).URL.Query()["from_google"]
				return ok
			}
			return false
		},
	})

	Admin.AddResource(Widgets, &admin.Config{Menu: []string{"Pages Management"}, Priority: 3})

	// Top Banner
	type bannerArgument struct {
		Title           string
		ButtonTitle     string
		Link            string
		BackgroundImage oss.OSS
		Logo            oss.OSS
	}

	Widgets.RegisterWidget(&widget.Widget{
		Name:        "NormalBanner",
		Templates:   []string{"banner", "banner2"},
		PreviewIcon: "/images/Widget-NormalBanner.png",
		Group:       "Banners",
		Setting:     Admin.NewResource(&bannerArgument{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Options["Setting"] = setting
			return context
		},
	})

	// Banner Editor
	type bannerEditorArgument struct {
		Value string
	}
	type subHeaderSetting struct {
		Text  string
		Color string
	}
	type headerSetting struct {
		Text  string
		Color string
	}
	type textsSetting struct {
		Text  string
		Color string
	}
	type buttonSetting struct {
		Text string
		Link string
	}

	type modelBuyLinkSetting struct {
		PostName string
		Price       string
		ButtonName  string
		Link        string
	}

	type imageSetting struct {
		Image oss.OSS
	}

	headerRes := Admin.NewResource(&headerSetting{})
	headerRes.Meta(&admin.Meta{Name: "Text"})
	headerRes.Meta(&admin.Meta{Name: "Color"})

	modelBuyLink := Admin.NewResource(&modelBuyLinkSetting{})
	modelBuyLink.Meta(&admin.Meta{Name: "PostName"})
	modelBuyLink.Meta(&admin.Meta{Name: "Price"})
	modelBuyLink.Meta(&admin.Meta{Name: "ButtonName"})
	modelBuyLink.Meta(&admin.Meta{Name: "Link"})

	subHeaderRes := Admin.NewResource(&subHeaderSetting{})
	subHeaderRes.Meta(&admin.Meta{Name: "Text"})
	subHeaderRes.Meta(&admin.Meta{Name: "Color"})

	textsRes := Admin.NewResource(&textsSetting{})
	textsRes.Meta(&admin.Meta{Name: "Text"})
	textsRes.Meta(&admin.Meta{Name: "Color"})

	buttonRes := Admin.NewResource(&buttonSetting{})
	buttonRes.Meta(&admin.Meta{Name: "Text"})
	buttonRes.Meta(&admin.Meta{Name: "Link"})

	imageRes := Admin.NewResource(&imageSetting{})
	imageRes.Meta(&admin.Meta{Name: "Image"})

	banner_editor.RegisterViewPath("github.com/qorpress/qorpress-example/pkg/app/views/banner_editor")
	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>short_text</i>",
		Name:     "Add Header",
		Template: "header",
		Resource: headerRes,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})
	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>format_list_bulleted</i>",
		Name:     "Add model buy block",
		Template: "model_buy_link",
		Resource: modelBuyLink,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})
	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>format_align_justify</i>",
		Name:     "Add Text",
		Template: "text",
		Resource: textsRes,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})
	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>dehaze</i>",
		Name:     "Add Sub Header",
		Template: "sub_header",
		Resource: subHeaderRes,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})
	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>view_stream</i>",
		Name:     "Add Button",
		Template: "button",
		Resource: buttonRes,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})

	banner_editor.RegisterElement(&banner_editor.Element{
		Icon:     "<i class=material-icons>image</i>",
		Name:     "Add Image",
		Template: "image",
		Resource: imageRes,
		Context: func(c *admin.Context, r interface{}) interface{} {
			return r.(banner_editor.QorBannerEditorSettingInterface).GetSerializableArgument(r.(banner_editor.QorBannerEditorSettingInterface))
		},
	})

	banner_editor.RegisterExternalStylePath("//fonts.googleapis.com/css?family=Lato|Playfair+Display|Raleway")
	banner_editor.RegisterExternalStylePath("/dist/qor.css")
	banner_editor.RegisterExternalStylePath("/dist/home_banner.css")

	bannerEditorResource := Admin.NewResource(&bannerEditorArgument{})
	bannerEditorResource.Meta(&admin.Meta{Name: "Value", Config: &banner_editor.BannerEditorConfig{
		MediaLibrary: Admin.GetResource("MediaLibrary"),
		Platforms: []banner_editor.Platform{
			{
				Name:     "Laptop",
				SafeArea: banner_editor.Size{Width: 0, Height: 0},
			},
			{
				Name:     "Mobile",
				SafeArea: banner_editor.Size{Width: 0, Height: 300},
			},
		},
	}})

	// normal banner editor
	Widgets.RegisterWidget(&widget.Widget{
		Name:        "BannerEditor",
		Templates:   []string{"banner_editor"},
		PreviewIcon: "/images/Widget-BannerEditor.png",
		Setting:     bannerEditorResource,
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Options["Value"] = template.HTML(setting.(*bannerEditorArgument).Value)
			return context
		},
	})

	// full width banner editor
	fullwidthBannerEditorResource := Admin.NewResource(&bannerEditorArgument{})
	fullwidthBannerEditorResource.Meta(&admin.Meta{Name: "Value", Config: &banner_editor.BannerEditorConfig{
		MediaLibrary: Admin.GetResource("MediaLibrary"),
		Platforms: []banner_editor.Platform{
			{
				Name:     "Laptop",
				SafeArea: banner_editor.Size{Width: 0, Height: 0},
			},
			{
				Name:     "Mobile",
				SafeArea: banner_editor.Size{Width: 0, Height: 300},
			},
		},
	}})

	Widgets.RegisterWidget(&widget.Widget{
		Name:        "FullWidthBannerEditor",
		Templates:   []string{"fullwidth_banner_editor"},
		PreviewIcon: "/images/Widget-FullWidthBannerEditor.png",
		Setting:     fullwidthBannerEditorResource,
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Options["Value"] = template.HTML(setting.(*bannerEditorArgument).Value)
			return context
		},
	})

	type slideImage struct {
		Title    string
		SubTitle string
		Button   string
		Link     string
		Image    oss.OSS
	}

	type slideShowArgument struct {
		SlideImages []slideImage
	}
	slideShowResource := Admin.NewResource(&slideShowArgument{})
	slideShowResource.AddProcessor(&resource.Processor{
		Handler: func(value interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if slides, ok := value.(*slideShowArgument); ok {
				for _, slide := range slides.SlideImages {
					if slide.Title == "" {
						return errors.New("slide title is blank")
					}
				}
			}
			return nil
		},
	})

	Widgets.RegisterWidget(&widget.Widget{
		Name:        "SlideShow",
		Templates:   []string{"slideshow"},
		PreviewIcon: "/images/Widget-SlideShow.png",
		Group:       "Banners",
		Setting:     slideShowResource,
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Options["Setting"] = setting
			return context
		},
	})

	Widgets.RegisterWidgetsGroup(&widget.WidgetsGroup{
		Name:    "Banner",
		Widgets: []string{"NormalBanner", "SlideShow"},
	})

	// selected Posts
	type selectedPostsArgument struct {
		Posts       []string
		PostsSorter sorting.SortableCollection
	}
	selectedPostsResource := Admin.NewResource(&selectedPostsArgument{})
	selectedPostsResource.Meta(&admin.Meta{Name: "Posts", Type: "select_many", Collection: func(value interface{}, context *qor.Context) [][]string {
		var collectionValues [][]string
		var posts []*posts.Post
		context.GetDB().Find(&posts)
		for _, post := range posts {
			collectionValues = append(collectionValues, []string{fmt.Sprintf("%v", post.ID), post.Name})
		}
		return collectionValues
	}})

	Widgets.RegisterWidget(&widget.Widget{
		Name:        "Posts",
		Templates:   []string{"posts"},
		Group:       "Posts",
		PreviewIcon: "/images/Widget-Posts.png",
		Setting:     selectedPostsResource,
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			if setting != nil {
				var posts []*posts.Post
				context.GetDB().Limit(9).Preload("ColorVariations").Where("id IN (?)", setting.(*selectedPostsArgument).Posts).Find(&posts)
				setting.(*selectedPostsArgument).PostsSorter.Sort(&posts)
				context.Options["Posts"] = posts
			}
			return context
		},
	})

	// FooterLinks
	type FooterItem struct {
		Name string
		Link string
	}

	type FooterSection struct {
		Title       string
		Items       []FooterItem
		ItemsSorter sorting.SortableCollection
	}

	type FooterLinks struct {
		Sections []FooterSection
	}

	Widgets.RegisterWidget(&widget.Widget{
		Name:        "Footer Links",
		PreviewIcon: "/images/Widget-FooterLinks.png",
		Setting:     Admin.NewResource(&FooterLinks{}),
		Context: func(context *widget.Context, setting interface{}) *widget.Context {
			context.Options["Setting"] = setting
			return context
		},
	})

	Widgets.RegisterFuncMap("render_banner_editor_content", func(val template.HTML, r *http.Request) template.HTML {
		bannerValue := banner_editor.GetContent(string(val), r)
		return template.HTML(bannerValue)
	})
}
