package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
)

type TemplateBlockRepository struct {
	db *sql.DB
}

func NewTemplateBlockRepository(db *sql.DB) *TemplateBlockRepository {
	return &TemplateBlockRepository{db: db}
}

func (r *TemplateBlockRepository) Create(block dto.CreateTemplateBlockDTO) error {
	query := `
		INSERT INTO public.template_block 
			(template_id, name, type, "order", exercise_count, estimated_duration_minutes, instructions)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	_, err := r.db.Exec(
		query,
		block.TemplateID,
		block.Name,
		block.Type,
		block.Order,
		block.ExerciseCount,
		block.EstimatedDurationMinutes,
		block.Instructions,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *TemplateBlockRepository) GetByID(id string) (dto.TemplateBlockDTO, error) {
	var block dto.TemplateBlockDTO
	query := `
		SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, created_at
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
		return dto.TemplateBlockDTO{}, err
	}
	return block, nil
}

func (r *TemplateBlockRepository) GetByTemplateID(templateID string) ([]dto.TemplateBlockDTO, error) {
	query := `
		SELECT id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, created_at
		FROM public.template_block WHERE template_id = $1 ORDER BY "order"`
	rows, err := r.db.Query(query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []dto.TemplateBlockDTO
	for rows.Next() {
		var block dto.TemplateBlockDTO
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

func (r *TemplateBlockRepository) Update(id string, block dto.TemplateBlockDTO) (dto.TemplateBlockDTO, error) {
	query := `
		UPDATE public.template_block
		SET template_id = $1, name = $2, type = $3, "order" = $4, exercise_count = $5, estimated_duration_minutes = $6, instructions = $7
		WHERE id = $8
		RETURNING id, template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, created_at`
	var updatedBlock dto.TemplateBlockDTO
	err := r.db.QueryRow(
		query,
		block.TemplateID,
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
		return dto.TemplateBlockDTO{}, err
	}
	return updatedBlock, nil
}

func (r *TemplateBlockRepository) Delete(id string) error {
	query := `DELETE FROM public.template_block WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
