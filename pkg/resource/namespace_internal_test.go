package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var dataset = []struct {
	now      []NamespaceEvent
	before   []NamespaceEvent
	expected []NamespaceEvent
}{
	{
		now: []NamespaceEvent{{
			Namespace: "foo",
			Status:    10,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}},
		before: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}},
		expected: []NamespaceEvent{{
			Namespace: "foo",
			Status:    10,
			Phase:     "Active",
		}},
	},
	{
		now: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}, {
			Namespace: "baz",
			Status:    80,
			Phase:     "Active",
		}},
		before: []NamespaceEvent{},
		expected: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}, {
			Namespace: "baz",
			Status:    80,
			Phase:     "Active",
		}},
	},
	{
		now: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}, {
			Namespace: "baz",
			Status:    80,
			Phase:     "Active",
		}},
		before: []NamespaceEvent{{
			Namespace: "bar",
			Status:    100,
			Phase:     "80",
		}, {
			Namespace: "baz",
			Status:    100,
			Phase:     "Active",
		}},
		expected: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "baz",
			Status:    80,
			Phase:     "Active",
		}},
	},
	{
		now: []NamespaceEvent{},
		before: []NamespaceEvent{{
			Namespace: "foo",
			Status:    0,
			Phase:     "Active",
		}, {
			Namespace: "bar",
			Status:    100,
			Phase:     "Active",
		}},
		expected: []NamespaceEvent(nil),
	},
}

func TestDiff(t *testing.T) {
	for _, data := range dataset {
		namespaceDiff := diff(data.now, data.before)

		assert.Equal(t, data.expected, namespaceDiff)
	}
}
