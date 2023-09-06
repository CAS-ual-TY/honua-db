package honuadb

import (
	"testing"

	"github.com/JonasBordewick/honua-db/models"
)

var testDashboard = &models.Dashboard{
	ID: "test",
	Widgets: []*models.Widget{
		{
			WidgetType: models.DEFAULT,
			Title:      "Test",
		},
		{
			WidgetType: models.ENTITY,
			Title:      "Test",
		},
	},
}

func TestDashboard(t *testing.T) {
	err := testInstance.AddDashboard(testDashboard)
	if err != nil {
		t.Fatalf("An error occured %v", err)
	}

	dash, err := testInstance.GetDashboard("test")
	if err != nil {
		t.Fatalf("An error occured %v", err)
	}

	if testDashboard.ID != dash.ID || len(testDashboard.Widgets) != len(dash.Widgets) {
		t.Fatalf("Want: %v --> Have: %v", testDashboard, dash)
	}

	err = testInstance.DeleteDashboard("test")
	if err != nil {
		t.Fatalf("An error occured %v", err)
	}
}
