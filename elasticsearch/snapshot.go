package elasticsearch

import (
	"context"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

// CreateSnapshot creates a DB snapshot as backup (or use for expansion packs)
func CreateSnapshot() error {
	ctx := context.Background()

	client, err := elastic.NewSimpleClient()
	if err != nil {
		return err
	}

	// create snapshot repo
	repo := elastic.
		NewSnapshotCreateRepositoryService(client).
		Repository("repo1").
		Verify(false).
		BodyJson(map[string]interface{}{
			"type": "fs",
			"settings": map[string]interface{}{
				"location": "/mount/backups/snapshots",
			},
		})

	repoResponse, err := repo.Do(ctx)
	if err != nil {
		panic(err)
	}

	if !repoResponse.Acknowledged {
		panic("could not ack repository create")
	}

	// create snapshot
	service := elastic.
		NewSnapshotCreateService(client).
		Repository("repo1").
		Snapshot(time.Now().Format("2006-01-02-15-04-05")).
		WaitForCompletion(true).
		BodyJson(map[string]interface{}{
			"indices":              "scifgif",
			"ignore_unavailable":   false,
			"include_global_state": true,
		})

	response, err := service.Do(ctx)
	if err != nil {
		panic(err)
	}

	if response.Snapshot.State != "SUCCESS" {
		panic("elasticsearch snapshot state != SUCCESS: " + response.Snapshot.State)
	}
	return nil
}
