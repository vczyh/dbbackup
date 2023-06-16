package s3client

import (
	"context"
	"testing"
)

func TestClient_ListBuckets(t *testing.T) {
	c := New(
		"QTBELHBAPSf3un1m57mG",
		"EXTw1meYdwhqZQEEpBDA9vDDOmQVF4dwlV69mbBb",
		"http://127.0.0.1:9000",
		WithForcePathStyle(true),
		WithRegion("test"),
	)

	ctx := context.TODO()
	buckets, err := c.ListBuckets(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buckets)
}
