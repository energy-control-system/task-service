package handler

import (
	"task-service/service/task"
	"testing"

	routerreflect "github.com/sunshineOfficial/golib/gohttp/gorouter/reflect"
)

func TestTaskListQueryVarsReadsPaginationWithFilters(t *testing.T) {
	var vars taskAllListQueryVars
	err := routerreflect.SetValuesToItem(map[string][]string{
		"limit":  {"10"},
		"offset": {"20"},
		"status": {"2"},
		"sort":   {"desc"},
	}, "query", &vars)
	if err != nil {
		t.Fatalf("SetValuesToItem returned error: %v", err)
	}

	page := vars.Pagination()
	if page.Limit != 10 {
		t.Fatalf("limit = %d, want 10", page.Limit)
	}
	if page.Offset != 20 {
		t.Fatalf("offset = %d, want 20", page.Offset)
	}
	if vars.Status == nil || *vars.Status != 2 {
		t.Fatalf("status = %v, want 2", vars.Status)
	}

	filter := vars.Filter()
	if filter.Sort != task.SortDesc {
		t.Fatalf("sort = %q, want %q", filter.Sort, task.SortDesc)
	}
}
