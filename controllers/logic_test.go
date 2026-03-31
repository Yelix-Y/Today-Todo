package controllers

import (
	"Today-Todo/models"
	"testing"
	"time"
)

func TestNormalizePriority(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "high", input: "high", want: "high"},
		{name: "medium", input: "medium", want: "medium"},
		{name: "low", input: "low", want: "low"},
		{name: "unknown falls back", input: "urgent", want: "medium"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizePriority(tc.input); got != tc.want {
				t.Fatalf("normalizePriority(%q)=%q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestClampProgress(t *testing.T) {
	if got := clampProgress(-12.5); got != 0 {
		t.Fatalf("clampProgress(-12.5)=%v, want 0", got)
	}
	if got := clampProgress(12.5); got != 12.5 {
		t.Fatalf("clampProgress(12.5)=%v, want 12.5", got)
	}
	if got := clampProgress(145.2); got != 100 {
		t.Fatalf("clampProgress(145.2)=%v, want 100", got)
	}
}

func TestCalcFocusScore(t *testing.T) {
	if got := calcFocusScore(4, 5, 1); got <= 70 {
		t.Fatalf("calcFocusScore should reward completion, got %d", got)
	}
	if got := calcFocusScore(0, 1, 50); got != 0 {
		t.Fatalf("calcFocusScore should be clamped to 0, got %d", got)
	}
	if got := calcFocusScore(10, 10, 0); got != 95 {
		t.Fatalf("calcFocusScore(10,10,0)=%d, want 95", got)
	}
}

func TestPriorityWeight(t *testing.T) {
	if got := priorityWeight("high"); got != 3 {
		t.Fatalf("priorityWeight(high)=%d, want 3", got)
	}
	if got := priorityWeight("medium"); got != 2 {
		t.Fatalf("priorityWeight(medium)=%d, want 2", got)
	}
	if got := priorityWeight("low"); got != 1 {
		t.Fatalf("priorityWeight(low)=%d, want 1", got)
	}
}

func TestPickTopTasks(t *testing.T) {
	now := time.Now()
	tasks := []models.Todo{
		{ID: 1, Title: "low old", Priority: "low", Completed: false, CreatedAt: now.Add(-4 * time.Hour)},
		{ID: 2, Title: "high new", Priority: "high", Completed: false, CreatedAt: now.Add(-1 * time.Hour)},
		{ID: 3, Title: "medium old", Priority: "medium", Completed: false, CreatedAt: now.Add(-3 * time.Hour)},
		{ID: 4, Title: "done high", Priority: "high", Completed: true, CreatedAt: now.Add(-5 * time.Hour)},
		{ID: 5, Title: "high old", Priority: "high", Completed: false, CreatedAt: now.Add(-6 * time.Hour)},
	}

	picked := pickTopTasks(tasks, 3)
	if len(picked) != 3 {
		t.Fatalf("pickTopTasks len=%d, want 3", len(picked))
	}

	if picked[0].Title != "high old" || picked[1].Title != "high new" || picked[2].Title != "medium old" {
		t.Fatalf("unexpected order: %#v", picked)
	}
}
