package domain

import (
	"fmt"
	"time"

	core_errors "github.com/cephalopagus/bkv-golang-todo/internal/core/errors"
)

type TaskPatch struct {
	Title       Nullable[string]
	Descriprion Nullable[string]
	Completed   Nullable[bool]
}

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func NewTaskUninitilized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UnitializedID,
		UnitializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}

func (t *Task) Validate() error {
	titleLenght := len([]rune(t.Title))

	if titleLenght < 1 || titleLenght > 100 {
		return fmt.Errorf(
			"invalide `Title` len: %d:, %w",
			titleLenght,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLenght := len([]rune(*t.Description))

		if descriptionLenght < 1 || descriptionLenght > 1000 {
			return fmt.Errorf(
				"invalide `Description` len: %d:, %w",
				descriptionLenght,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` cannot be `nil` if `Completed` equals `true`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` cannot be before `CreatedAt`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be `nil` if `Completed` equals `false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func NewTaskPatch(
	title Nullable[string],
	descriprion Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Descriprion: descriprion,
		Completed:   completed,
	}
}

func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("'Title' cannot be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf("'Completed' cannot be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (t *Task) ApplyPatch(patch TaskPatch) error {

	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Descriprion.Set {
		tmp.Description = patch.Descriprion.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value

		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp
	return nil
}
