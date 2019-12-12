package cache

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nbedos/citop/text"
	"github.com/nbedos/citop/utils"
)

type task struct {
	key         PipelineStepKey
	ref         GitRef
	type_       string
	state       State
	name        string
	provider    string
	prefix      string
	createdAt   utils.NullTime
	startedAt   utils.NullTime
	finishedAt  utils.NullTime
	updatedAt   utils.NullTime
	duration    utils.NullDuration
	children    []*task
	traversable bool
	url         string
}

func (t task) Diff(other task) string {
	options := cmp.AllowUnexported(PipelineStepKey{}, task{})
	return cmp.Diff(t, other, options)
}

func (t task) Traversable() bool {
	return t.traversable
}

func (t task) Children() []utils.TreeNode {
	children := make([]utils.TreeNode, len(t.children))
	for i := range t.children {
		children[i] = t.children[i]
	}
	return children
}

func (t task) Tabular(loc *time.Location) map[string]text.StyledString {
	const nullPlaceholder = "-"

	nullTimeToString := func(t utils.NullTime) text.StyledString {
		s := nullPlaceholder
		if t.Valid {
			s = t.Time.In(loc).Truncate(time.Second).Format("Jan 2 15:04")
		}
		return text.NewStyledString(s)
	}

	state := text.NewStyledString(string(t.state))
	switch t.state {
	case Failed, Canceled:
		state.Add(text.StatusFailed)
	case Passed:
		state.Add(text.StatusPassed)
	case Running:
		state.Add(text.StatusRunning)
	case Pending, Skipped, Manual:
		state.Add(text.StatusSkipped)
	}

	name := text.NewStyledString(t.prefix)
	if t.type_ == "P" {
		name.Append(t.provider, text.Provider)
	} else {
		name.Append(t.name)
	}

	pipeline := t.key.ID
	if _, err := strconv.Atoi(t.key.ID); err == nil {
		pipeline = "#" + pipeline
	}

	refClass := text.GitBranch
	if t.ref.IsTag {
		refClass = text.GitTag
	}

	return map[string]text.StyledString{
		"REF":      text.NewStyledString(t.ref.Ref, refClass),
		"PIPELINE": text.NewStyledString(pipeline),
		"TYPE":     text.NewStyledString(t.type_),
		"STATE":    state,
		"NAME":     name,
		"CREATED":  nullTimeToString(t.createdAt),
		"STARTED":  nullTimeToString(t.startedAt),
		"FINISHED": nullTimeToString(t.finishedAt),
		"UPDATED":  nullTimeToString(t.updatedAt),
		"DURATION": text.NewStyledString(t.duration.String()),
	}
}

func (t task) Key() interface{} {
	return t.key
}

func (t task) URL() string {
	return t.url
}

func (t *task) SetTraversable(traversable bool, recursive bool) {
	t.traversable = traversable
	if recursive {
		for _, child := range t.children {
			child.SetTraversable(traversable, recursive)
		}
	}
}

func (t *task) SetPrefix(s string) {
	t.prefix = s
}

/*
func ref(ref string, tag bool) string {
	if tag {
		return fmt.Sprintf("tag: %s", ref)
	}
	return ref
}

func taskFromBuild(b Build) task {
	ref := ref(b.Ref, b.IsTag)
	row := task{
		key: taskKey{
			ref:       ref,
			SHA:       b.SHA,
			accountID: b.Repository.Provider.ID,
			buildID:   b.ID,
		},
		type_:      "P",
		state:      b.State,
		createdAt:  b.CreatedAt,
		startedAt:  b.StartedAt,
		finishedAt: b.FinishedAt,
		updatedAt:  utils.NullTime{Time: b.UpdatedAt, Valid: true},
		url:        b.WebURL,
		duration:   b.Duration,
		provider:   b.Repository.Provider.Name,
	}

	// Prefix only numeric IDs with hash
	if _, err := strconv.Atoi(b.ID); err == nil {
		row.name = fmt.Sprintf("#%s", b.ID)
	} else {
		row.name = b.ID
	}

	for _, job := range b.Jobs {
		child := taskFromJob(b.Repository.Provider, b.SHA, ref, b.ID, 0, *job)
		row.children = append(row.children, &child)
	}

	if b.Stages != nil {
		stageIDs := make([]int, 0, len(b.Stages))
		for stageID := range b.Stages {
			stageIDs = append(stageIDs, stageID)
		}
		sort.Ints(stageIDs)
		for _, stageID := range stageIDs {
			child := taskFromStage(b.Repository.Provider, b.SHA, ref, b.ID, b.WebURL, *b.Stages[stageID])
			row.children = append(row.children, &child)
		}
	}

	return row
}

func taskFromStage(provider Provider, SHA string, ref string, buildID string, webURL string, s Stage) task {
	row := task{
		key: taskKey{
			ref:       ref,
			SHA:       SHA,
			accountID: provider.ID,
			buildID:   buildID,
			stageID:   s.ID,
		},
		type_:    "S",
		state:    s.State,
		name:     s.Name,
		url:      webURL,
		provider: provider.Name,
	}

	// We aggregate jobs by name and only keep the most recent to weed out previous runs of the job.
	// This is mainly for GitLab which keeps jobs after they are restarted.
	jobByName := make(map[string]*Job, len(s.Jobs))
	for _, job := range s.Jobs {
		namedJob, exists := jobByName[job.Name]
		if !exists || job.CreatedAt.Valid && job.CreatedAt.Time.After(namedJob.CreatedAt.Time) {
			jobByName[job.Name] = job
		}
	}

	for _, job := range jobByName {
		row.createdAt = utils.MinNullTime(row.createdAt, job.CreatedAt)
		row.startedAt = utils.MinNullTime(row.startedAt, job.StartedAt)
		row.finishedAt = utils.MaxNullTime(row.finishedAt, job.FinishedAt)
		row.updatedAt = utils.MaxNullTime(row.updatedAt, job.FinishedAt, job.StartedAt, job.CreatedAt)
	}

	row.duration = utils.NullSub(row.finishedAt, row.startedAt)

	for _, job := range s.Jobs {
		child := taskFromJob(provider, SHA, ref, buildID, s.ID, *job)
		row.children = append(row.children, &child)
	}

	return row
}

func taskFromJob(provider Provider, SHA string, ref string, buildID string, stageID int, j Job) task {
	name := j.Name
	if name == "" {
		name = j.ID
	}
	return task{
		key: taskKey{
			ref:       ref,
			SHA:       SHA,
			accountID: provider.ID,
			buildID:   buildID,
			stageID:   stageID,
			jobID:     j.ID,
		},
		type_:      "J",
		state:      j.State,
		name:       name,
		createdAt:  j.CreatedAt,
		startedAt:  j.StartedAt,
		finishedAt: j.FinishedAt,
		updatedAt:  utils.MaxNullTime(j.FinishedAt, j.StartedAt, j.CreatedAt),
		url:        j.WebURL,
		duration:   j.Duration,
		provider:   provider.Name,
	}
}*/

type BuildsByCommit struct {
	cache Cache
	ref   string
}

func (c Cache) BuildsOfRef(ref string) HierarchicalTabularDataSource {
	return BuildsByCommit{
		cache: c,
		ref:   ref,
	}
}

func (s BuildsByCommit) Headers() []string {
	return []string{"REF", "PIPELINE", "TYPE", "STATE", "CREATED", "DURATION", "NAME"}
}

func (s BuildsByCommit) Alignment() map[string]text.Alignment {
	return map[string]text.Alignment{
		"REF":      text.Left,
		"PIPELINE": text.Right,
		"TYPE":     text.Right,
		"STATE":    text.Left,
		"CREATED":  text.Left,
		"STARTED":  text.Left,
		"UPDATED":  text.Left,
		"DURATION": text.Right,
		"NAME":     text.Left,
	}
}

func (s BuildsByCommit) Rows() []HierarchicalTabularSourceRow {
	rows := make([]HierarchicalTabularSourceRow, 0)
	//for _, pipeline := range s.cache.PipelinesByRef(s.ref) {
	for range s.cache.PipelinesByRef(s.ref) {
		// FIXME convert pipeline to task
		row := task{
			key:         PipelineStepKey{},
			type_:       "",
			state:       "",
			name:        "",
			provider:    "",
			prefix:      "",
			createdAt:   utils.NullTime{},
			startedAt:   utils.NullTime{},
			finishedAt:  utils.NullTime{},
			updatedAt:   utils.NullTime{},
			duration:    utils.NullDuration{},
			children:    nil,
			traversable: false,
			url:         "",
		}
		rows = append(rows, &row)
	}

	sort.Slice(rows, func(i, j int) bool {
		ri, rj := rows[i].(*task), rows[j].(*task)
		ti := utils.MinNullTime(
			ri.createdAt,
			ri.startedAt,
			ri.updatedAt,
			ri.finishedAt)

		tj := utils.MinNullTime(
			rj.createdAt,
			rj.startedAt,
			rj.updatedAt,
			rj.finishedAt)

		return ti.Time.Before(tj.Time)
	})

	return rows
}

var ErrNoLogHere = errors.New("no log is associated to this row")

func (s BuildsByCommit) WriteToDisk(ctx context.Context, key interface{}, dir string) (string, error) {
	// TODO Allow filtering for errored jobs
	stepKey, ok := key.(PipelineStepKey)
	if !ok {
		return "", fmt.Errorf("key conversion to taskKey failed: '%v'", key)
	}

	file, err := ioutil.TempFile(dir, "step_*.log")
	w := utils.NewANSIStripper(file)
	defer w.Close()
	if err != nil {
		return "", err
	}
	logPath := path.Join(dir, filepath.Base(file.Name()))

	err = s.cache.WriteLog(ctx, stepKey, w)
	return logPath, err
}
