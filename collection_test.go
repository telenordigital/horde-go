package nbiot

import (
	"net/http"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	collection, err := client.CreateCollection(Collection{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteCollection(collection.ID)

	tagKey := "test key"
	tagValue := "test value"
	collection.Tags = map[string]string{tagKey: tagValue}
	collection, err = client.UpdateCollection(collection)
	if err != nil {
		t.Fatal(err)
	}
	if len(collection.Tags) != 1 || collection.Tags[tagKey] != tagValue {
		t.Fatal("unexpected tags:", collection.Tags)
	}

	collections, err := client.Collections()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, c := range collections {
		if c.ID == collection.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("collection %v not found in %v", collection, collections)
	}

	if _, err := client.Collection(collection.ID); err != nil {
		t.Fatal(err)
	}

	data, err := client.CollectionData(collection.ID, time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Fatalf("collection %v should be empty, has %d elements", collection.ID, len(data))
	}

	if err := client.DeleteCollection(collection.ID); err != nil {
		t.Fatal(err)
	}
	err = client.DeleteCollection(collection.ID)
	if cerr, ok := err.(ClientError); !ok || cerr.HTTPStatusCode != http.StatusNotFound {
		t.Fatal(err)
	}

}
