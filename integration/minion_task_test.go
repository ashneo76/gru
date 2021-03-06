package integration

import (
	"reflect"
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

func TestMinionTaskBacklog(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-task-backlog")
	defer tc.recorder.Stop()

	// Setup our minion
	minionName := "Kevin"
	cfg := &minion.EtcdMinionConfig{
		Name:       minionName,
		EtcdConfig: tc.config,
	}

	testMinion, err := minion.NewEtcdMinion(cfg)
	if err != nil {
		t.Fatal(err)
	}

	minionID := testMinion.ID()

	err = testMinion.SetName(minionName)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test task and submit it
	wantTask := task.New("foo", "bar")
	wantTask.ID = uuid.Parse("e6d2bebd-2219-4a8c-9d30-a861097c147e")

	err = tc.client.MinionSubmitTask(minionID, wantTask)
	if err != nil {
		t.Fatal(err)
	}

	// Get pending tasks and verify the task we sent is the task we get
	backlog, err := tc.client.MinionTaskQueue(minionID)
	if err != nil {
		t.Fatal(err)
	}

	if len(backlog) != 1 {
		t.Errorf("want 1 backlog task, got %d backlog tasks", len(backlog))
	}

	gotTask := backlog[0]
	if !reflect.DeepEqual(wantTask, gotTask) {
		t.Errorf("want %q task, got %q task", wantTask.ID, gotTask.ID)
	}
}
