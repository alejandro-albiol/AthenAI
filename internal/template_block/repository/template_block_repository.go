package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/interfaces"
)

// TemplateBlockRepository provides database operations for template blocks.
type TemplateBlockRepository struct {
	db *sql.DB
}

// NewTemplateBlockRepository creates a new TemplateBlockRepository.
func NewTemplateBlockRepository(db *sql.DB) interfaces.TemplateBlockRepository {
	return &TemplateBlockRepository{db: db}
}

// CreateTemplateBlock inserts a new template block and returns its ID.
func (r *TemplateBlockRepository) CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (*string, error) {
	query := `
		INSERT INTO public.template_block 
			(template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	var id string
	err := r.db.QueryRow(
		query,
		block.TemplateID,
		block.Name,
		block.Type,
		block.Order,
		block.ExerciseCount,
		block.EstimatedDurationMinutes,
		block.Instructions,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// GetTemplateBlockByID retrieves a template block by its ID.
func (r *TemplateBlockRepository) GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error) {
	block := &dto.TemplateBlockDTO{}
	query := `
		SELECT id, template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions, created_at
		FROM public.template_block WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&block.ID,
		&block.TemplateID,
		&block.Name,
		&block.Type,
		&block.Order,
		&block.ExerciseCount,
		&block.EstimatedDurationMinutes,
		&block.Instructions,
		&block.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// GetTemplateBlocksByTemplateID retrieves all template blocks for a given template ID.
func (r *TemplateBlockRepository) GetTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error) {
	query := `
		SELECT id, template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions, created_at
		FROM public.template_block WHERE template_id = $1 ORDER BY "order"`
	rows, err := r.db.Query(query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*dto.TemplateBlockDTO
	for rows.Next() {
		block := &dto.TemplateBlockDTO{}
		err := rows.Scan(
			&block.ID,
			&block.TemplateID,
			&block.Name,
			&block.Type,
			&block.Order,
			&block.ExerciseCount,
			&block.EstimatedDurationMinutes,
			&block.Instructions,
			&block.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// GetTemplateBlockByTemplateIDAndName retrieves a template block by template ID and name.
func (r *TemplateBlockRepository) GetTemplateBlockByTemplateIDAndName(templateID string, name string) (*dto.TemplateBlockDTO, error) {
	block := &dto.TemplateBlockDTO{}
	query := `SELECT id, template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions, created_at FROM public.template_block WHERE template_id = $1 AND name = $2`
	err := r.db.QueryRow(query, templateID, name).Scan(
		&block.ID,
		&block.TemplateID,
		&block.Name,
		&block.Type,
		&block.Order,
		&block.ExerciseCount,
		&block.EstimatedDurationMinutes,
		&block.Instructions,
		&block.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return block, nil
}

// Update modifies an existing template block and returns the updated block.
func (r *TemplateBlockRepository) UpdateTemplateBlock(id string, block *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error) {
	query := `
		UPDATE public.template_block
		SET name = $1, type = $2, "order" = $3, exercise_count = $4, estimated_duration_minutes = $5, instructions = $6
		WHERE id = $7
		RETURNING id, template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions, created_at`
	updatedBlock := &dto.TemplateBlockDTO{}
	err := r.db.QueryRow(
		query,
		block.Name,
		block.Type,
		block.Order,
		block.ExerciseCount,
		block.EstimatedDurationMinutes,
		block.Instructions,
		id,
	).Scan(
		&updatedBlock.ID,
		&updatedBlock.TemplateID,
		&updatedBlock.Name,
		&updatedBlock.Type,
		&updatedBlock.Order,
		&updatedBlock.ExerciseCount,
		&updatedBlock.EstimatedDurationMinutes,
		&updatedBlock.Instructions,
		&updatedBlock.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return updatedBlock, nil
}

// DeleteTemplateBlock removes a template block by its ID.
func (r *TemplateBlockRepository) DeleteTemplateBlock(id string) error {
	query := `DELETE FROM public.template_block WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
