package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"jobflow/internal/jobs"
	"jobflow/pkg"
)

const jobColumns = "id,body,status,attempts,result,error"

func (s *Store) CreateJob(ctx context.Context, j jobs.Job) (jobs.Job, bool, error) {
	tag, err := s.pool.Exec(ctx, "INSERT INTO jobs(id,body) VALUES($1,$2) ON CONFLICT(id) DO NOTHING", j.ID, j.Text)
	if err != nil {
		return jobs.Job{}, false, translate(err)
	}
	existing, err := s.GetJob(ctx, j.ID)
	if err != nil {
		return jobs.Job{}, false, err
	}
	if existing.Text != j.Text {
		return jobs.Job{}, false, pkg.ErrConflict
	}
	return existing, tag.RowsAffected() == 1, nil
}
func (s *Store) GetJob(ctx context.Context, id string) (jobs.Job, error) {
	var j jobs.Job
	err := s.pool.QueryRow(ctx, "SELECT "+jobColumns+" FROM jobs WHERE id=$1", id).Scan(&j.ID, &j.Text, &j.Status, &j.Attempts, &j.Result, &j.Error)
	return j, translate(err)
}
func (s *Store) Claim(ctx context.Context) (jobs.Job, bool, error) {
	var j jobs.Job
	err := s.pool.QueryRow(ctx, `UPDATE jobs SET status='running',attempts=attempts+1
 WHERE id=(SELECT id FROM jobs WHERE status='queued' AND attempts<3 AND available_at<=now() ORDER BY id LIMIT 1)
 AND status='queued' AND attempts<3 RETURNING `+jobColumns).Scan(&j.ID, &j.Text, &j.Status, &j.Attempts, &j.Result, &j.Error)
	if errors.Is(err, pgx.ErrNoRows) {
		return jobs.Job{}, false, nil
	}
	return j, err == nil, translate(err)
}
func (s *Store) Finish(ctx context.Context, j jobs.Job, result string, cause error) error {
	status := "done"
	message := ""
	if cause != nil {
		status = "failed"
		message = cause.Error()
		result = ""
		if j.Attempts < 3 && (errors.Is(cause, jobs.ErrTemporary) || errors.Is(cause, context.Canceled) || errors.Is(cause, context.DeadlineExceeded)) {
			status = "queued"
		}
	}
	tag, err := s.pool.Exec(ctx, "UPDATE jobs SET status=$2,result=$3,error=$4,available_at=now()+interval '100 milliseconds' WHERE id=$1 AND status='running' AND attempts=$5", j.ID, status, result, message, j.Attempts)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return pkg.ErrConflict
	}
	return nil
}
func (s *Store) Recover(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, "UPDATE jobs SET status=CASE WHEN attempts<3 THEN 'queued' ELSE 'failed' END,error='процесс был остановлен',available_at=now() WHERE status='running'")
	return err
}
