package knowledgebase

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Matches the connection info in docker-compose.dev.yml / database.NewClient's
// default, so this runs against the same local pgvector instance the app uses.
const testDBDefaultURL = "host=localhost port=5432 user=user password=pwd dbname=supportops sslmode=disable"

const embeddingDims = 384

func connectTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = testDBDefaultURL
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Skipf("skipping: could not connect to test database at %q: %v", dsn, err)
	}
	return db
}

// vec builds a deterministic embeddingDims-length vector, zero-padded, with
// the given components placed at index 0, 1, 2, ...
func vec(components ...float32) []float32 {
	v := make([]float32, embeddingDims)
	copy(v, components)
	return v
}

// insertFixture creates a knowledge_base row tagged with sourceType (used to
// scope assertions and cleanup to this test only) and returns its source_id.
func insertFixture(t *testing.T, repo *Repository, sourceType string, embedding []float32) string {
	t.Helper()
	sourceID := uuid.NewString()
	_, err := repo.Create(context.Background(), CreateRequest{
		SourceId:        sourceID,
		SourceEventType: sourceType,
		Content:         json.RawMessage(`{"note":"fixture"}`),
		Embedding:       embedding,
	})
	if err != nil {
		t.Fatalf("failed to insert fixture: %v", err)
	}
	return sourceID
}

func cleanupFixtures(t *testing.T, db *sqlx.DB, sourceType string) {
	t.Helper()
	if _, err := db.Exec(`DELETE FROM knowledge_base WHERE source_type = $1`, sourceType); err != nil {
		t.Logf("cleanup failed for source_type=%s: %v", sourceType, err)
	}
}

// TestFind_OrdersByCosineDistance pins down that results come back nearest
// first. query is the unit vector e0; A is identical to it (distance 0), B
// and C lean increasingly away from e0 toward e1, and D is e1 itself
// (orthogonal to the query, i.e. maximally distant).
func TestFind_OrdersByCosineDistance(t *testing.T) {
	db := connectTestDB(t)
	repo := NewRepository(db)
	sourceType := "test_order_" + uuid.NewString()
	t.Cleanup(func() { cleanupFixtures(t, db, sourceType) })

	query := vec(1, 0)
	idA := insertFixture(t, repo, sourceType, vec(1, 0))     // identical to query
	idB := insertFixture(t, repo, sourceType, vec(0.9, 0.1))  // close
	idC := insertFixture(t, repo, sourceType, vec(0.5, 0.5))  // further
	idD := insertFixture(t, repo, sourceType, vec(0, 1))      // orthogonal, farthest

	results, err := repo.Find(context.Background(), searchInput{
		Embedding:  query,
		SourceType: &sourceType,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}

	want := []string{idA, idB, idC, idD}
	if len(results) != len(want) {
		t.Fatalf("got %d results, want %d: %+v", len(results), len(want), results)
	}
	for i, row := range results {
		if row.SourceId != want[i] {
			t.Errorf("position %d: got source_id %s, want %s", i, row.SourceId, want[i])
		}
	}
}

// TestFind_PaginationOffsetRegression guards against the off-by-one where
// getOffset() computed page*limit instead of (page-1)*limit, which caused
// the first page to skip `limit` rows and return nothing whenever the table
// had fewer than `limit` matching rows.
func TestFind_PaginationOffsetRegression(t *testing.T) {
	db := connectTestDB(t)
	repo := NewRepository(db)
	sourceType := "test_page_" + uuid.NewString()
	t.Cleanup(func() { cleanupFixtures(t, db, sourceType) })

	query := vec(1, 0)
	idClosest := insertFixture(t, repo, sourceType, vec(1, 0))
	idMiddle := insertFixture(t, repo, sourceType, vec(0.9, 0.1))
	idFarthest := insertFixture(t, repo, sourceType, vec(0, 1))

	firstPage, err := repo.Find(context.Background(), searchInput{
		Embedding:  query,
		SourceType: &sourceType,
		Limit:      2,
		Page:       1,
	})
	if err != nil {
		t.Fatalf("Find (page 1) returned error: %v", err)
	}
	if len(firstPage) != 2 {
		t.Fatalf("page 1: got %d results, want 2 (the two nearest fixtures): %+v", len(firstPage), firstPage)
	}
	if firstPage[0].SourceId != idClosest || firstPage[1].SourceId != idMiddle {
		t.Errorf("page 1: got order [%s, %s], want [%s, %s]",
			firstPage[0].SourceId, firstPage[1].SourceId, idClosest, idMiddle)
	}

	secondPage, err := repo.Find(context.Background(), searchInput{
		Embedding:  query,
		SourceType: &sourceType,
		Limit:      2,
		Page:       2,
	})
	if err != nil {
		t.Fatalf("Find (page 2) returned error: %v", err)
	}
	if len(secondPage) != 1 {
		t.Fatalf("page 2: got %d results, want 1 (the remaining fixture): %+v", len(secondPage), secondPage)
	}
	if secondPage[0].SourceId != idFarthest {
		t.Errorf("page 2: got source_id %s, want %s", secondPage[0].SourceId, idFarthest)
	}
}

// TestFind_NoMatchesReturnsEmptyNotNil guards against Find returning a nil
// slice on zero matches, which serializes to JSON null and breaks callers
// (e.g. the ml-worker) that iterate over the response directly.
func TestFind_NoMatchesReturnsEmptyNotNil(t *testing.T) {
	db := connectTestDB(t)
	repo := NewRepository(db)
	sourceType := "test_empty_" + uuid.NewString()

	results, err := repo.Find(context.Background(), searchInput{
		Embedding:  vec(1, 0),
		SourceType: &sourceType,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("Find returned error: %v", err)
	}
	if results == nil {
		t.Error("Find returned nil slice on zero matches, want empty non-nil slice")
	}
	if len(results) != 0 {
		t.Errorf("got %d results for a source_type with no fixtures, want 0", len(results))
	}
}
