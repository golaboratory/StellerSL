package gamification

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/stellersl/backend/internal/db"
)

// fakeQuerier implements the subset of db.Querier the badge engine uses.
// The embedded interface makes it satisfy db.Querier; calling an unstubbed
// method panics, which is the desired behavior in tests.
type fakeQuerier struct {
	db.Querier

	userBadges            []db.Badge
	allBadges             []db.Badge
	stats                 db.GetDashboardStatsRow
	createdToday          int32
	completedToday        int32
	weekendCount          int32
	distinctProjectsToday int32
	lastActivity          time.Time
	streak                db.UserStreak
	streakErr             error
	projectCount          int32
	teamCount             int32
	hourCounts            map[int32]int32
	lastHourCount         int32
	tasks                 map[uuid.UUID]db.Task
	projectCounts         map[uuid.UUID]db.GetProjectTaskCountsRow

	awarded   []uuid.UUID
	removed   []uuid.UUID
	upserted  []db.UpsertStreakParams
}

func (f *fakeQuerier) ListUserBadges(ctx context.Context, userID uuid.UUID) ([]db.Badge, error) {
	return f.userBadges, nil
}
func (f *fakeQuerier) ListAllBadges(ctx context.Context) ([]db.Badge, error) {
	return f.allBadges, nil
}
func (f *fakeQuerier) GetDashboardStats(ctx context.Context, assignedTo uuid.NullUUID) (db.GetDashboardStatsRow, error) {
	return f.stats, nil
}
func (f *fakeQuerier) CountTasksCreatedToday(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.createdToday, nil
}
func (f *fakeQuerier) CountTasksCompletedToday(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.completedToday, nil
}
func (f *fakeQuerier) CountTasksCompletedOnWeekends(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.weekendCount, nil
}
func (f *fakeQuerier) CountDistinctProjectsCompletedToday(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.distinctProjectsToday, nil
}
func (f *fakeQuerier) GetLastActivityDate(ctx context.Context, userID uuid.UUID) (time.Time, error) {
	return f.lastActivity, nil
}
func (f *fakeQuerier) GetStreak(ctx context.Context, arg db.GetStreakParams) (db.UserStreak, error) {
	if f.streakErr != nil {
		return db.UserStreak{}, f.streakErr
	}
	return f.streak, nil
}
func (f *fakeQuerier) UpsertStreak(ctx context.Context, arg db.UpsertStreakParams) error {
	f.upserted = append(f.upserted, arg)
	return nil
}
func (f *fakeQuerier) CountProjectsByUser(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.projectCount, nil
}
func (f *fakeQuerier) CountTeamMembershipsByUser(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.teamCount, nil
}
func (f *fakeQuerier) CountTasksCompletedTodayInHour(ctx context.Context, arg db.CountTasksCompletedTodayInHourParams) (int32, error) {
	return f.hourCounts[arg.Column2], nil
}
func (f *fakeQuerier) CountTasksCompletedLastHour(ctx context.Context, userID uuid.UUID) (int32, error) {
	return f.lastHourCount, nil
}
func (f *fakeQuerier) GetTask(ctx context.Context, id uuid.UUID) (db.Task, error) {
	t, ok := f.tasks[id]
	if !ok {
		return db.Task{}, sql.ErrNoRows
	}
	return t, nil
}
func (f *fakeQuerier) GetProjectTaskCounts(ctx context.Context, projectID uuid.NullUUID) (db.GetProjectTaskCountsRow, error) {
	return f.projectCounts[projectID.UUID], nil
}
func (f *fakeQuerier) AwardBadge(ctx context.Context, arg db.AwardBadgeParams) error {
	f.awarded = append(f.awarded, arg.BadgeID)
	return nil
}
func (f *fakeQuerier) RemoveBadge(ctx context.Context, arg db.RemoveBadgeParams) error {
	f.removed = append(f.removed, arg.BadgeID)
	return nil
}

// allRequirementTypes covers every badge implemented in the engine
// (Phase A 19種 + Phase B 10種).
var allRequirementTypes = []string{
	"first_task", "three_strike", "rookie_planner", "first_project", "first_team_join",
	"daily_runner_bronze", "daily_runner_silver", "daily_runner_gold", "daily_runner_platinum",
	"comeback_hero", "weekend_master", "monday_motivator", "friday_finisher",
	"task_crusher_1", "task_crusher_2", "task_crusher_3", "task_crusher_4", "task_crusher_5",
	"early_bird", "midnight_owl", "quick_strike", "lightning_fast", "all_nighter",
	"lunch_hustler", "golden_hour", "hat_trick",
	"organizer", "multi_tasker", "zero_inbox",
}

type badgeFixture struct {
	badges []db.Badge
	byID   map[uuid.UUID]string
	byType map[string]db.Badge
}

func newBadgeFixture() badgeFixture {
	fx := badgeFixture{byID: map[uuid.UUID]string{}, byType: map[string]db.Badge{}}
	for _, rt := range allRequirementTypes {
		b := db.Badge{
			ID:              uuid.NewSHA1(uuid.NameSpaceOID, []byte(rt)),
			Name:            rt,
			IconSlug:        "pi-star",
			RequirementType: rt,
		}
		fx.badges = append(fx.badges, b)
		fx.byID[b.ID] = rt
		fx.byType[rt] = b
	}
	return fx
}

func (fx badgeFixture) awardedTypes(f *fakeQuerier) map[string]bool {
	got := map[string]bool{}
	for _, id := range f.awarded {
		got[fx.byID[id]] = true
	}
	return got
}

func (fx badgeFixture) removedTypes(f *fakeQuerier) map[string]bool {
	got := map[string]bool{}
	for _, id := range f.removed {
		got[fx.byID[id]] = true
	}
	return got
}

// mustParse builds a fixed clock for deterministic time-based conditions.
func fixedClock(value string) func() time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return func() time.Time { return t }
}

const (
	// 2026-06-08 is a Monday, 2026-06-12 is a Friday.
	mondayNoon  = "2026-06-08T12:30:00Z"
	fridayNoon  = "2026-06-12T12:30:00Z"
	tuesdayNoon = "2026-06-09T12:30:00Z"
)

func TestAwardBadges(t *testing.T) {
	fx := newBadgeFixture()
	userID := uuid.New()
	taskID := uuid.New()
	projectID := uuid.New()

	mkTask := func(age time.Duration, clock func() time.Time, withProject bool) db.Task {
		task := db.Task{
			ID:        taskID,
			CreatedAt: sql.NullTime{Time: clock().Add(-age), Valid: true},
		}
		if withProject {
			task.ProjectID = uuid.NullUUID{UUID: projectID, Valid: true}
		}
		return task
	}

	cases := []struct {
		name     string
		setup    func(f *fakeQuerier)
		clock    string
		trigger  string
		taskIDs  []uuid.UUID
		want     []string
		wantNot  []string
	}{
		{
			name:    "first_task and three_strike on counts",
			setup:   func(f *fakeQuerier) { f.stats.CompletedTasks = 3 },
			trigger: "task_completed",
			want:    []string{"first_task", "three_strike"},
			wantNot: []string{"task_crusher_1"},
		},
		{
			name:    "rookie_planner fires on creation trigger only",
			setup:   func(f *fakeQuerier) { f.createdToday = 5 },
			trigger: "task_created",
			want:    []string{"rookie_planner"},
		},
		{
			name:    "rookie_planner does not fire on completion trigger",
			setup:   func(f *fakeQuerier) { f.createdToday = 5 },
			trigger: "task_completed",
			wantNot: []string{"rookie_planner"},
		},
		{
			name:    "first_project and organizer thresholds",
			setup:   func(f *fakeQuerier) { f.projectCount = 3 },
			trigger: "task_created",
			want:    []string{"first_project", "organizer"},
		},
		{
			name:    "first_team_join",
			setup:   func(f *fakeQuerier) { f.teamCount = 1 },
			trigger: "task_completed",
			want:    []string{"first_team_join"},
		},
		{
			name: "daily runner tiers from streak",
			setup: func(f *fakeQuerier) {
				f.streak = db.UserStreak{CurrentCount: 7, MaxCount: 100}
			},
			trigger: "task_completed",
			want: []string{
				"daily_runner_bronze", "daily_runner_silver",
				"daily_runner_gold", "daily_runner_platinum",
			},
		},
		{
			name: "streak too short awards nothing",
			setup: func(f *fakeQuerier) {
				f.streak = db.UserStreak{CurrentCount: 2, MaxCount: 2}
			},
			trigger: "task_completed",
			wantNot: []string{"daily_runner_bronze"},
		},
		{
			name: "comeback_hero after 8 day gap",
			setup: func(f *fakeQuerier) {
				f.lastActivity, _ = time.Parse(time.RFC3339, "2026-05-31T00:00:00Z")
			},
			clock:   mondayNoon,
			trigger: "task_completed",
			want:    []string{"comeback_hero"},
		},
		{
			name: "comeback_hero not for first ever completion (1970 sentinel)",
			setup: func(f *fakeQuerier) {
				f.lastActivity, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			},
			clock:   mondayNoon,
			trigger: "task_completed",
			wantNot: []string{"comeback_hero"},
		},
		{
			name: "comeback_hero not after short gap",
			setup: func(f *fakeQuerier) {
				f.lastActivity, _ = time.Parse(time.RFC3339, "2026-06-06T00:00:00Z")
			},
			clock:   mondayNoon,
			trigger: "task_completed",
			wantNot: []string{"comeback_hero"},
		},
		{
			name:    "weekend_master",
			setup:   func(f *fakeQuerier) { f.weekendCount = 10 },
			trigger: "task_completed",
			want:    []string{"weekend_master"},
		},
		{
			name:    "monday_motivator on Monday with 5 completions",
			setup:   func(f *fakeQuerier) { f.completedToday = 5 },
			clock:   mondayNoon,
			trigger: "task_completed",
			want:    []string{"monday_motivator"},
			wantNot: []string{"friday_finisher"},
		},
		{
			name:    "friday_finisher on Friday with 5 completions",
			setup:   func(f *fakeQuerier) { f.completedToday = 5 },
			clock:   fridayNoon,
			trigger: "task_completed",
			want:    []string{"friday_finisher"},
			wantNot: []string{"monday_motivator"},
		},
		{
			name:    "weekday badges silent on Tuesday",
			setup:   func(f *fakeQuerier) { f.completedToday = 5 },
			clock:   tuesdayNoon,
			trigger: "task_completed",
			wantNot: []string{"monday_motivator", "friday_finisher"},
		},
		{
			name:    "task crusher ladder",
			setup:   func(f *fakeQuerier) { f.stats.CompletedTasks = 500 },
			trigger: "task_completed",
			want:    []string{"task_crusher_1", "task_crusher_2", "task_crusher_3", "task_crusher_4"},
			wantNot: []string{"task_crusher_5"},
		},
		{
			name:    "early_bird at 05:00",
			clock:   "2026-06-09T05:00:00Z",
			trigger: "task_completed",
			want:    []string{"early_bird"},
			wantNot: []string{"midnight_owl"},
		},
		{
			name:    "midnight_owl at 02:00",
			clock:   "2026-06-09T02:00:00Z",
			trigger: "task_completed",
			want:    []string{"midnight_owl"},
			wantNot: []string{"early_bird"},
		},
		{
			name: "quick_strike within an hour, lightning_fast not after 30min",
			setup: func(f *fakeQuerier) {
				f.tasks = map[uuid.UUID]db.Task{taskID: mkTask(30*time.Minute, fixedClock(tuesdayNoon), false)}
			},
			clock:   tuesdayNoon,
			trigger: "task_completed",
			taskIDs: []uuid.UUID{taskID},
			want:    []string{"quick_strike"},
			wantNot: []string{"lightning_fast"},
		},
		{
			name: "lightning_fast within 10 minutes",
			setup: func(f *fakeQuerier) {
				f.tasks = map[uuid.UUID]db.Task{taskID: mkTask(5*time.Minute, fixedClock(tuesdayNoon), false)}
			},
			clock:   tuesdayNoon,
			trigger: "task_completed",
			taskIDs: []uuid.UUID{taskID},
			want:    []string{"quick_strike", "lightning_fast"},
		},
		{
			name:    "all_nighter with 10 completions today",
			setup:   func(f *fakeQuerier) { f.completedToday = 10 },
			clock:   tuesdayNoon,
			trigger: "task_completed",
			want:    []string{"all_nighter"},
		},
		{
			name:    "lunch_hustler with 3 completions in the 12 o'clock hour",
			setup:   func(f *fakeQuerier) { f.hourCounts = map[int32]int32{12: 3} },
			trigger: "task_completed",
			want:    []string{"lunch_hustler"},
			wantNot: []string{"golden_hour"},
		},
		{
			name:    "golden_hour with 3 completions in the 17 o'clock hour",
			setup:   func(f *fakeQuerier) { f.hourCounts = map[int32]int32{17: 3} },
			trigger: "task_completed",
			want:    []string{"golden_hour"},
			wantNot: []string{"lunch_hustler"},
		},
		{
			name:    "hat_trick with 3 completions in the last hour",
			setup:   func(f *fakeQuerier) { f.lastHourCount = 3 },
			trigger: "task_completed",
			want:    []string{"hat_trick"},
		},
		{
			name:    "multi_tasker with 3 distinct projects today",
			setup:   func(f *fakeQuerier) { f.distinctProjectsToday = 3 },
			trigger: "task_completed",
			want:    []string{"multi_tasker"},
		},
		{
			name: "zero_inbox when the task's project has no remaining tasks",
			setup: func(f *fakeQuerier) {
				f.tasks = map[uuid.UUID]db.Task{taskID: mkTask(2*time.Hour, fixedClock(tuesdayNoon), true)}
				f.projectCounts = map[uuid.UUID]db.GetProjectTaskCountsRow{
					projectID: {Total: 4, Remaining: 0},
				}
			},
			clock:   tuesdayNoon,
			trigger: "task_completed",
			taskIDs: []uuid.UUID{taskID},
			want:    []string{"zero_inbox"},
		},
		{
			name: "zero_inbox not while tasks remain",
			setup: func(f *fakeQuerier) {
				f.tasks = map[uuid.UUID]db.Task{taskID: mkTask(2*time.Hour, fixedClock(tuesdayNoon), true)}
				f.projectCounts = map[uuid.UUID]db.GetProjectTaskCountsRow{
					projectID: {Total: 4, Remaining: 2},
				}
			},
			clock:   tuesdayNoon,
			trigger: "task_completed",
			taskIDs: []uuid.UUID{taskID},
			wantNot: []string{"zero_inbox"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakeQuerier{allBadges: fx.badges, streakErr: sql.ErrNoRows}
			if tc.setup != nil {
				f.streakErr = nil
				tc.setup(f)
			}
			engine := NewBadgeEngine()
			clock := tc.clock
			if clock == "" {
				clock = tuesdayNoon
			}
			engine.now = fixedClock(clock)

			if err := engine.AwardBadges(context.Background(), f, userID, tc.trigger, tc.taskIDs...); err != nil {
				t.Fatalf("AwardBadges returned error: %v", err)
			}

			got := fx.awardedTypes(f)
			for _, want := range tc.want {
				if !got[want] {
					t.Errorf("expected badge %q to be awarded, got %v", want, got)
				}
			}
			for _, not := range tc.wantNot {
				if got[not] {
					t.Errorf("expected badge %q NOT to be awarded, got %v", not, got)
				}
			}
		})
	}
}

func TestAwardBadgesSkipsAlreadyOwned(t *testing.T) {
	fx := newBadgeFixture()
	f := &fakeQuerier{
		allBadges:  fx.badges,
		userBadges: []db.Badge{fx.byType["first_task"]},
		streakErr:  sql.ErrNoRows,
	}
	f.stats.CompletedTasks = 1

	engine := NewBadgeEngine()
	engine.now = fixedClock(tuesdayNoon)
	if err := engine.AwardBadges(context.Background(), f, uuid.New(), "task_completed"); err != nil {
		t.Fatalf("AwardBadges returned error: %v", err)
	}
	if got := fx.awardedTypes(f); got["first_task"] {
		t.Errorf("first_task should not be re-awarded, got %v", got)
	}
}

func TestReEvaluateBadges(t *testing.T) {
	fx := newBadgeFixture()
	owned := []string{
		"first_task", "three_strike", "task_crusher_1", "weekend_master",
		// permanent ones that must survive
		"early_bird", "all_nighter", "quick_strike", "daily_runner_bronze",
		"monday_motivator", "zero_inbox", "comeback_hero",
	}
	var userBadges []db.Badge
	for _, rt := range owned {
		userBadges = append(userBadges, fx.byType[rt])
	}

	f := &fakeQuerier{allBadges: fx.badges, userBadges: userBadges}
	f.stats.CompletedTasks = 2 // below three_strike(3) and task_crusher_1(10)
	f.weekendCount = 4         // below weekend_master(10)

	engine := NewBadgeEngine()
	engine.now = fixedClock(tuesdayNoon)
	if err := engine.ReEvaluateBadges(context.Background(), f, uuid.New()); err != nil {
		t.Fatalf("ReEvaluateBadges returned error: %v", err)
	}

	removed := fx.removedTypes(f)
	for _, want := range []string{"three_strike", "task_crusher_1", "weekend_master"} {
		if !removed[want] {
			t.Errorf("expected %q to be revoked, removed=%v", want, removed)
		}
	}
	for _, keep := range []string{
		"first_task", // completed=2 still >= 1
		"early_bird", "all_nighter", "quick_strike", "daily_runner_bronze",
		"monday_motivator", "zero_inbox", "comeback_hero",
	} {
		if removed[keep] {
			t.Errorf("expected %q to be kept, removed=%v", keep, removed)
		}
	}
}

func TestUpdateStreak(t *testing.T) {
	userID := uuid.New()
	today := time.Date(2026, 6, 9, 0, 0, 0, 0, time.Local)

	t.Run("first completion starts a streak", func(t *testing.T) {
		f := &fakeQuerier{streakErr: sql.ErrNoRows}
		st := NewStreakTracker()
		st.now = func() time.Time { return today.Add(10 * time.Hour) }
		if err := st.UpdateStreak(context.Background(), f, userID); err != nil {
			t.Fatalf("UpdateStreak: %v", err)
		}
		if len(f.upserted) != 1 || f.upserted[0].CurrentCount != 1 || f.upserted[0].MaxCount != 1 {
			t.Fatalf("unexpected upsert: %+v", f.upserted)
		}
	})

	t.Run("consecutive day increments", func(t *testing.T) {
		f := &fakeQuerier{streak: db.UserStreak{
			CurrentCount: 2, MaxCount: 5,
			LastDate: sql.NullTime{Time: today.AddDate(0, 0, -1), Valid: true},
		}}
		st := NewStreakTracker()
		st.now = func() time.Time { return today.Add(10 * time.Hour) }
		if err := st.UpdateStreak(context.Background(), f, userID); err != nil {
			t.Fatalf("UpdateStreak: %v", err)
		}
		if len(f.upserted) != 1 || f.upserted[0].CurrentCount != 3 || f.upserted[0].MaxCount != 5 {
			t.Fatalf("unexpected upsert: %+v", f.upserted)
		}
	})

	t.Run("gap resets to one", func(t *testing.T) {
		f := &fakeQuerier{streak: db.UserStreak{
			CurrentCount: 4, MaxCount: 4,
			LastDate: sql.NullTime{Time: today.AddDate(0, 0, -3), Valid: true},
		}}
		st := NewStreakTracker()
		st.now = func() time.Time { return today.Add(10 * time.Hour) }
		if err := st.UpdateStreak(context.Background(), f, userID); err != nil {
			t.Fatalf("UpdateStreak: %v", err)
		}
		if len(f.upserted) != 1 || f.upserted[0].CurrentCount != 1 || f.upserted[0].MaxCount != 4 {
			t.Fatalf("unexpected upsert: %+v", f.upserted)
		}
	})

	t.Run("same day is a no-op", func(t *testing.T) {
		f := &fakeQuerier{streak: db.UserStreak{
			CurrentCount: 2, MaxCount: 2,
			LastDate: sql.NullTime{Time: today, Valid: true},
		}}
		st := NewStreakTracker()
		st.now = func() time.Time { return today.Add(10 * time.Hour) }
		if err := st.UpdateStreak(context.Background(), f, userID); err != nil {
			t.Fatalf("UpdateStreak: %v", err)
		}
		if len(f.upserted) != 0 {
			t.Fatalf("expected no upsert, got %+v", f.upserted)
		}
	})
}

func TestCharacterTypeForLevel(t *testing.T) {
	cases := map[int32]string{
		1: "egg", 4: "egg",
		5: "chick", 9: "chick",
		10: "chicken", 19: "chicken",
		20: "phoenix", 42: "phoenix",
	}
	for level, want := range cases {
		if got := CharacterTypeForLevel(level); got != want {
			t.Errorf("CharacterTypeForLevel(%d) = %q, want %q", level, got, want)
		}
	}
}
