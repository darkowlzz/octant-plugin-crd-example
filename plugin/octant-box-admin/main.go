package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/flexlayout"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var pluginName = "box-admin"

func main() {
	boxGVK := schema.GroupVersionKind{Group: "foo.my.domain", Version: "v1alpha1", Kind: "Box"}
	boxRecordGVK := schema.GroupVersionKind{Group: "foo.my.domain", Version: "v1alpha1", Kind: "BoxRecord"}

	capabilities := &plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{boxGVK, boxRecordGVK},
		SupportsObjectStatus:  []schema.GroupVersionKind{boxGVK, boxRecordGVK},
		SupportsTab:           []schema.GroupVersionKind{boxGVK, boxRecordGVK},
		// ActionNames:           []string{pluginActionName},
		IsModule: true,
	}

	// Set up what should happen when Octant calls this plugin.
	options := []service.PluginOption{
		service.WithPrinter(handlePrint),
		service.WithObjectStatus(handleStatus),
		service.WithTabPrinter(handleTab),
		service.WithNavigation(handleNavigation, initRoutes),
		// service.WithActionHandler(handleAction),
	}

	// Use the plugin service helper to register this plugin.
	p, err := service.Register(pluginName, "box-admin octant plugin", capabilities, options...)
	if err != nil {
		log.Fatal(err)
	}

	// The plugin can log and the log messages will show up in Octant.
	log.Printf("octant-box-admin is starting")
	p.Serve()
}

// handleTab is called when Octant wants to print a tab for an object.
func handleTab(request *service.PrintRequest) (plugin.TabResponse, error) {
	if request.Object == nil {
		return plugin.TabResponse{}, errors.New("object is nil")
	}

	layout := flexlayout.New()
	section := layout.AddSection()
	// Octant contains a library of components that can be used to display content.
	// This example uses markdown text.
	contents := component.NewMarkdownText("admin content from box-admin plugin")

	err := section.Add(contents, component.WidthHalf)
	if err != nil {
		return plugin.TabResponse{}, err
	}

	// In this example, this plugin will tell Octant to create a new
	// tab when showing boxes. This tab's name will be "Extra box Details".
	tab := component.NewTabWithContents(*layout.ToComponent("Box Admin Details"))

	return plugin.TabResponse{Tab: tab}, nil
}

func handleStatus(request *service.PrintRequest) (plugin.ObjectStatusResponse, error) {
	if request.Object == nil {
		return plugin.ObjectStatusResponse{}, errors.Errorf("object is nil")
	}

	key, err := store.KeyFromObject(request.Object)
	if err != nil {
		return plugin.ObjectStatusResponse{}, err
	}
	u, err := request.DashboardClient.Get(request.Context(), key)
	if err != nil {
		return plugin.ObjectStatusResponse{}, err
	}

	// The plugin can check if the object it requested exists.
	if u == nil {
		return plugin.ObjectStatusResponse{}, errors.New("object doesn't exist")
	}

	// Will add object UID to the Resource Viewer properties table
	return plugin.ObjectStatusResponse{
		ObjectStatus: component.PodSummary{
			Properties: []component.Property{{
				Label: "ID (from box-admin plugin)",
				Value: component.NewText(string(u.GetUID())),
			}},
		},
	}, nil
}

// handlePrint is called when Octant wants to print an object.
func handlePrint(request *service.PrintRequest) (plugin.PrintResponse, error) {
	if request.Object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	key, err := store.KeyFromObject(request.Object)
	if err != nil {
		return plugin.PrintResponse{}, err
	}
	u, err := request.DashboardClient.Get(request.Context(), key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}

	// The plugin can check if the object it requested exists.
	if u == nil {
		return plugin.PrintResponse{}, errors.New("object doesn't exist")
	}

	boxCard := component.NewCard(component.TitleFromString(fmt.Sprintf("Admin Output for %s", u.GetName())))
	boxCard.SetBody(component.NewMarkdownText("This output was generated for admin by _octant-box-admin_"))

	msg := fmt.Sprintf("update from plugin at %s", time.Now().Format(time.RFC3339))

	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Status: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Items: []component.FlexLayoutItem{
			{
				Width: component.WidthHalf,
				View:  boxCard,
			},
		},
	}, nil
}

func handleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title:    "Box Admin",
		Path:     request.GeneratePath(),
		IconName: "nodes",
	}, nil
}

func initRoutes(router *service.Router) {
	router.HandleFunc("*", func(request service.Request) (component.ContentResponse, error) {
		var title []component.TitleComponent
		text := component.NewText("Box Admin Dashboard")
		title = component.Title(text)

		contentResponse := component.NewContentResponse(title)

		cols := component.NewTableCols("A", "B")
		table := component.NewTable("Table title", "Empty text", cols)
		table.Add(component.TableRow{
			"A": component.NewText("foo"),
			"B": component.NewText("bar"),
		})
		table.Add(component.TableRow{
			"A": component.NewText("xoo"),
			"B": component.NewText("yoo"),
		})

		contentResponse.Add(table)

		return *contentResponse, nil
	})
}
