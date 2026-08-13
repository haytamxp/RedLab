package reporting

import "testing"

func TestReportTemplateMetadata(t *testing.T) {
	if reportTemplate == nil {
		t.Fatal("report template must be initialized")
	}
}
