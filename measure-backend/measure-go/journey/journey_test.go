package journey

import (
	"encoding/json"
	"measure-backend/measure-go/event"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/yourbasic/graph"
)

func readEvents(path string) (events []event.EventField, err error) {
	bytes, err := os.ReadFile(path)
	events = []event.EventField{}
	json.Unmarshal(bytes, &events)
	return
}

func TestNewJourneyAndroidOne(t *testing.T) {
	events, err := readEvents("events_one.json")
	if err != nil {
		panic(err)
	}

	journey := NewJourneyAndroid(events)

	expectedOrder := 4
	gotOrder := journey.Graph.Order()

	if expectedOrder != gotOrder {
		t.Errorf("Expected %d order, but got %d", expectedOrder, gotOrder)
	}

	expectedString := "4 [{0 1} {0 2} {0 3}]"
	gotString := journey.Graph.String()

	if expectedString != gotString {
		t.Errorf("Expected %q, got %q", expectedString, gotString)
	}

	// forward direction
	{
		sessionIds := journey.GetEdgeSessions(0, 1)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("a3d629f5-6bab-4a43-8e75-fa5d6b539d33"),
			uuid.MustParse("58e94ae9-a084-479f-9049-2c5135f6090f"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 2)
		expectedLen := 3
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("a3d629f5-6bab-4a43-8e75-fa5d6b539d33"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 3)
		expectedLen := 2
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
	}

	// reverse direction
	{
		sessionIds := journey.GetEdgeSessions(1, 0)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("a3d629f5-6bab-4a43-8e75-fa5d6b539d33"),
			uuid.MustParse("58e94ae9-a084-479f-9049-2c5135f6090f"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(2, 0)
		expectedLen := 3
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("a3d629f5-6bab-4a43-8e75-fa5d6b539d33"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(3, 0)
		expectedLen := 2
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("9e44aa3a-3d67-4a56-8a76-a9fff7e2aae9"),
			uuid.MustParse("460765ab-1834-454e-b207-d8235b2160d9"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
	}

	{
		expected := graph.Stats{
			Size:     6,
			Multi:    0,
			Weighted: 0,
			Loops:    0,
			Isolated: 0,
		}
		got := graph.Check(journey.Graph)

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected %v graph stats, but got %v", expected, got)
		}
	}
}

func TestNewJourneyAndroidTwo(t *testing.T) {
	events, err := readEvents("events_two.json")
	if err != nil {
		panic(err)
	}

	journey := NewJourneyAndroid(events)

	expectedOrder := 5
	gotOrder := journey.Graph.Order()

	if expectedOrder != gotOrder {
		t.Errorf("Expected %d order, but got %d", expectedOrder, gotOrder)
	}

	expectedString := "5 [(0 1) (0 2) (0 3) (4 0)]"
	gotString := journey.Graph.String()

	if expectedString != gotString {
		t.Errorf("Expected %q, got %q", expectedString, gotString)
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 1)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("4339f2be-ec13-4858-9b7f-322e5ddf55f4"),
			uuid.MustParse("65aaf877-e000-4ff3-9f8f-a0dbb10e9b00"),
			uuid.MustParse("1755de51-18c8-4c14-a58d-ad677485130e"),
			uuid.MustParse("bcafd264-43eb-433b-8851-00306ecc2706"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 2)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("4339f2be-ec13-4858-9b7f-322e5ddf55f4"),
			uuid.MustParse("65aaf877-e000-4ff3-9f8f-a0dbb10e9b00"),
			uuid.MustParse("1755de51-18c8-4c14-a58d-ad677485130e"),
			uuid.MustParse("bcafd264-43eb-433b-8851-00306ecc2706"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 3)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("4339f2be-ec13-4858-9b7f-322e5ddf55f4"),
			uuid.MustParse("65aaf877-e000-4ff3-9f8f-a0dbb10e9b00"),
			uuid.MustParse("1755de51-18c8-4c14-a58d-ad677485130e"),
			uuid.MustParse("bcafd264-43eb-433b-8851-00306ecc2706"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(0, 3)
		expectedLen := 4
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("4339f2be-ec13-4858-9b7f-322e5ddf55f4"),
			uuid.MustParse("65aaf877-e000-4ff3-9f8f-a0dbb10e9b00"),
			uuid.MustParse("1755de51-18c8-4c14-a58d-ad677485130e"),
			uuid.MustParse("bcafd264-43eb-433b-8851-00306ecc2706"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
		if expected[1] != sessionIds[1] {
			t.Errorf("Expected %v, but got %v", expected[1], sessionIds[1])
		}
		if expected[2] != sessionIds[2] {
			t.Errorf("Expected %v, but got %v", expected[2], sessionIds[2])
		}
		if expected[3] != sessionIds[3] {
			t.Errorf("Expected %v, but got %v", expected[3], sessionIds[3])
		}
	}

	{
		sessionIds := journey.GetEdgeSessions(4, 0)
		expectedLen := 1
		gotLen := len(sessionIds)
		if expectedLen != gotLen {
			t.Errorf("Expected %d length, got %d", expectedLen, gotLen)
		}

		expected := []uuid.UUID{
			uuid.MustParse("65aaf877-e000-4ff3-9f8f-a0dbb10e9b00"),
		}
		if expected[0] != sessionIds[0] {
			t.Errorf("Expected %v, but got %v", expected[0], sessionIds[0])
		}
	}

	{
		expected := graph.Stats{
			Size:     4,
			Multi:    0,
			Weighted: 0,
			Loops:    0,
			Isolated: 3,
		}
		got := graph.Check(journey.Graph)

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected %v graph stats, but got %v", expected, got)
		}
	}
}
