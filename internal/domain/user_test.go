package domain

import (
	"testing"
	"time"

	"github.com/shii-park/friends/internal/errs"
)

func TestNewUser_NameRequired(t *testing.T) {
	t.Parallel()

	_, err := NewUser("   ", "", "", 1, 1)
	if err != errs.ErrUserNameRequired {
		t.Fatalf("expected ErrUserNameRequired, got %v", err)
	}
}

func TestNewUser_InvalidBirthday(t *testing.T) {
	t.Parallel()

	cases := []struct {
		month int
		day   int
	}{
		{0, 1},
		{13, 1},
		{1, 0},
		{1, 32},
		{2, 30}, // 2000年で存在しない日
	}

	for _, c := range cases {
		_, err := NewUser("alice", "", "", c.month, c.day)
		if err != errs.ErrInvalidBirthday {
			t.Fatalf("month=%d day=%d: expected ErrInvalidBirthday, got %v", c.month, c.day, err)
		}
	}
}

func TestNewUser_InitialValues(t *testing.T) {
	t.Parallel()

	u, err := NewUser(" alice ", "icon", "msg", 12, 31)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID == "" {
		t.Fatal("expected ID to be set")
	}
	if u.Name != "alice" {
		t.Fatalf("expected trimmed name 'alice', got '%s'", u.Name)
	}
	if u.Birthmonth != 12 || u.Birthday != 31 {
		t.Fatalf("expected birth 12/31, got %d/%d", u.Birthmonth, u.Birthday)
	}
	if u.StreakLogin != ResetStreakLoginDay {
		t.Fatalf("expected StreakLogin=%d, got %d", ResetStreakLoginDay, u.StreakLogin)
	}
	if u.RankPoint != InitialRankPoint {
		t.Fatalf("expected RankPoint=%d, got %d", InitialRankPoint, u.RankPoint)
	}
	if u.Coin != InitialCoin {
		t.Fatalf("expected Coin=%d, got %d", InitialCoin, u.Coin)
	}
	if u.GachaStone != InitialGachaStone {
		t.Fatalf("expected GachaStone=%d, got %d", InitialGachaStone, u.GachaStone)
	}
	if u.RegisteredAt.IsZero() || u.LatestLoginAt.IsZero() {
		t.Fatalf("expected timestamps to be set, RegisteredAt=%v LatestLoginAt=%v", u.RegisteredAt, u.LatestLoginAt)
	}
}

func TestNewGuestUser(t *testing.T) {
	t.Parallel()

	u := NewGuestUser()
	if u == nil {
		t.Fatal("expected non-nil user")
	}
	if u.Name != guestName {
		t.Fatalf("expected guest name '%s', got '%s'", guestName, u.Name)
	}
	if u.Birthmonth != guestBirthmonth || u.Birthday != guestBirthday {
		t.Fatalf("expected guest birth %d/%d, got %d/%d", guestBirthmonth, guestBirthday, u.Birthmonth, u.Birthday)
	}
}

func TestUpdateUserName(t *testing.T) {
	t.Parallel()

	u, err := NewUser("alice", "", "", 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// empty
	if err := u.UpdateUserName("   "); err != errs.ErrUserNameRequired {
		t.Fatalf("expected ErrUserNameRequired, got %v", err)
	}

	// same (trimmed)
	if err := u.UpdateUserName(" alice "); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if u.Name != "alice" {
		t.Fatalf("expected name to remain 'alice', got '%s'", u.Name)
	}

	// change
	if err := u.UpdateUserName("bob "); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if u.Name != "bob" {
		t.Fatalf("expected name 'bob', got '%s'", u.Name)
	}
}

func TestAddConsumeCoin(t *testing.T) {
	t.Parallel()

	u, _ := NewUser("alice", "", "", 1, 1)

	if err := u.AddCoin(0); err != errs.ErrInvalidCoinDelta {
		t.Fatalf("expected ErrInvalidCoinDelta, got %v", err)
	}
	if err := u.AddCoin(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Coin != 10 {
		t.Fatalf("expected coin=10, got %d", u.Coin)
	}

	if err := u.ConsumeCoin(0); err != errs.ErrInvalidCoinDelta {
		t.Fatalf("expected ErrInvalidCoinDelta, got %v", err)
	}
	if err := u.ConsumeCoin(999); err != errs.ErrInsufficientCoin {
		t.Fatalf("expected ErrInsufficientCoin, got %v", err)
	}
	if err := u.ConsumeCoin(3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Coin != 7 {
		t.Fatalf("expected coin=7, got %d", u.Coin)
	}
}

func TestAddConsumeGachaStone(t *testing.T) {
	t.Parallel()

	u, _ := NewUser("alice", "", "", 1, 1)

	if err := u.AddGachaStone(0); err != errs.ErrInvalidGachaStoneDelta {
		t.Fatalf("expected ErrInvalidGachaStoneDelta, got %v", err)
	}
	if err := u.AddGachaStone(5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.GachaStone != 5 {
		t.Fatalf("expected gachaStone=5, got %d", u.GachaStone)
	}

	if err := u.ConsumeGachaStone(0); err != errs.ErrInvalidGachaStoneDelta {
		t.Fatalf("expected ErrInvalidGachaStoneDelta, got %v", err)
	}
	if err := u.ConsumeGachaStone(999); err != errs.ErrInsufficientGachaStone {
		t.Fatalf("expected ErrInsufficientGachaStone, got %v", err)
	}
	if err := u.ConsumeGachaStone(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.GachaStone != 3 {
		t.Fatalf("expected gachaStone=3, got %d", u.GachaStone)
	}
}

func TestAddConsumeRankPoint(t *testing.T) {
	t.Parallel()

	u, _ := NewUser("alice", "", "", 1, 1)

	if err := u.AddRankPoint(0); err != errs.ErrInvalidRankPointDelta {
		t.Fatalf("expected ErrInvalidRankPointDelta, got %v", err)
	}
	if err := u.AddRankPoint(10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.RankPoint != 10 {
		t.Fatalf("expected rankPoint=10, got %d", u.RankPoint)
	}

	if err := u.ConsumeRankPoint(0); err != errs.ErrInvalidRankPointDelta {
		t.Fatalf("expected ErrInvalidRankPointDelta, got %v", err)
	}

	// 仕様：不足時は 0 にしてエラーにしない
	if err := u.ConsumeRankPoint(999); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if u.RankPoint != 0 {
		t.Fatalf("expected rankPoint=0, got %d", u.RankPoint)
	}
}

func TestLoginDay_Cutoff4AM(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("JST", 9*60*60)

	// 2026-02-24 03:00 JST は 4時間引くと前日 23:00 なので「ログイン日」は 2/23 扱い
	t1 := time.Date(2026, 2, 24, 3, 0, 0, 0, loc)
	d1 := loginDay(t1)
	if d1.Year() != 2026 || d1.Month() != 2 || d1.Day() != 23 {
		t.Fatalf("expected loginDay=2026-02-23, got %v", d1)
	}

	// 2026-02-24 05:00 JST は 4時間引くと 01:00 なので「ログイン日」は 2/24 扱い
	t2 := time.Date(2026, 2, 24, 5, 0, 0, 0, loc)
	d2 := loginDay(t2)
	if d2.Year() != 2026 || d2.Month() != 2 || d2.Day() != 24 {
		t.Fatalf("expected loginDay=2026-02-24, got %v", d2)
	}
}
