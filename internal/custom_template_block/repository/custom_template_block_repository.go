package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
)

// CustomTemplateBlockRepository provides database operations for custom template blocks.
type CustomTemplateBlockRepository struct {
	db *sql.DB
}

// NewCustomTemplateBlockRepository creates a new CustomTemplateBlockRepository.
func NewCustomTemplateBlockRepository(db *sql.DB) *CustomTemplateBlockRepository {
	return &CustomTemplateBlockRepository{db: db}
}

// CreateCustomTemplateBlock inserts a new custom template block and returns its ID.
func (r *CustomTemplateBlockRepository) CreateCustomTemplateBlock(gymID string, block *dto.CreateCustomTemplateBlockDTO) (*string, error) {
	query := `
		INSERT INTO "%s".custom_template_block 
			(template_id, block_name, block_type, block_order, exercise_count, estimated_duration_minutes, instructions, reps, series, rest_time_seconds, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	var id string
	err := r.db.QueryRow(
		fmt.Sprintf(query, gymID),
		block.TemplateID,
		block.BlockName,
		block.BlockType,
		block.BlockOrder,
		block.ExerciseCount,
		block.EstimatedDurationMinutes,
		block.Instructions,
		block.Reps,
		block.Series,
		block.RestTimeSeconds,
		block.CreatedBy,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

// UpdateCustomTemplateBlock updates a custom template block by ID.
func (r *CustomTemplateBlockRepository) UpdateCustomTemplateBlock(gymID, id string, update *dto.UpdateCustomTemplateBlockDTO) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if update.BlockName != nil {
		setParts = append(setParts, fmt.Sprintf("block_name = $%d", argIndex))
		args = append(args, *update.BlockName)
		argIndex++
	}
	if update.BlockType != nil {
		setParts = append(setParts, fmt.Sprintf("block_type = $%d", argIndex))
		args = append(args, *update.BlockType)
		argIndex++
	}
	if update.BlockOrder != nil {
		setParts = append(setParts, fmt.Sprintf("block_order = $%d", argIndex))
		args = append(args, *update.BlockOrder)
		argIndex++
	}
	if update.ExerciseCount != nil {
		setParts = append(setParts, fmt.Sprintf("exercise_count = $%d", argIndex))
		args = append(args, *update.ExerciseCount)
		argIndex++
	}
	if update.EstimatedDurationMinutes != nil {
		setParts = append(setParts, fmt.Sprintf("estimated_duration_minutes = $%d", argIndex))
		args = append(args, *update.EstimatedDurationMinutes)
		argIndex++
	}
	if update.Instructions != nil {
		setParts = append(setParts, fmt.Sprintf("instructions = $%d", argIndex))
		args = append(args, *update.Instructions)
		argIndex++
	}
	if update.Reps != nil {
		setParts = append(setParts, fmt.Sprintf("reps = $%d", argIndex))
		args = append(args, *update.Reps)
		argIndex++
	}
	if update.Series != nil {
		setParts = append(setParts, fmt.Sprintf("series = $%d", argIndex))
		args = append(args, *update.Series)
		argIndex++
	}
	if update.RestTimeSeconds != nil {
		setParts = append(setParts, fmt.Sprintf("rest_time_seconds = $%d", argIndex))
		args = append(args, *update.RestTimeSeconds)
		argIndex++
	}
	if update.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *update.IsActive)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE "%s".custom_template_block
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE id = $%d AND deleted_at IS NULL`, gymID, strings.Join(setParts, ", "), argIndex)

	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetCustomTemplateBlockByID retrieves a custom template block by ID.
func (r *CustomTemplateBlockRepository) GetCustomTemplateBlockByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	query := `
		SELECT id, template_id, block_name, block_type, block_order, exercise_count, 
			   estimated_duration_minutes, instructions, reps, series, rest_time_seconds, 
			   created_at, updated_at, deleted_at, is_active, created_by
		FROM "%s".custom_template_block 
		WHERE id = $1 AND deleted_at IS NULL`

	row := r.db.QueryRow(fmt.Sprintf(query, gymID), id)
	var res dto.ResponseCustomTemplateBlockDTO
	err := row.Scan(
		&res.ID,
		&res.TemplateID,
		&res.BlockName,
		&res.BlockType,
		&res.BlockOrder,
		&res.ExerciseCount,
		&res.EstimatedDurationMinutes,
		&res.Instructions,
		&res.Reps,
		&res.Series,
		&res.RestTimeSeconds,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.DeletedAt,
		&res.IsActive,
		&res.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// ListCustomTemplateBlocksByTemplateID retrieves all custom template blocks for a specific template.
func (r *CustomTemplateBlockRepository) ListCustomTemplateBlocksByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	query := `
		SELECT id, template_id, block_name, block_type, block_order, exercise_count, 
			   estimated_duration_minutes, instructions, reps, series, rest_time_seconds, 
			   created_at, updated_at, deleted_at, is_active, created_by
		FROM "%s".custom_template_block 
		WHERE template_id = $1 AND deleted_at IS NULL 
		ORDER BY block_order ASC`

	rows, err := r.db.Query(fmt.Sprintf(query, gymID), templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.ResponseCustomTemplateBlockDTO
	for rows.Next() {
		var res dto.ResponseCustomTemplateBlockDTO
		if err := rows.Scan(
			&res.ID,
			&res.TemplateID,
			&res.BlockName,
			&res.BlockType,
			&res.BlockOrder,
			&res.ExerciseCount,
			&res.EstimatedDurationMinutes,
			&res.Instructions,
			&res.Reps,
			&res.Series,
			&res.RestTimeSeconds,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
			&res.IsActive,
			&res.CreatedBy,
		); err != nil {
			return nil, err
		}
		result = append(result, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// ListCustomTemplateBlocks retrieves all custom template blocks for a gym.
func (r *CustomTemplateBlockRepository) ListCustomTemplateBlocks(gymID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	query := `
		SELECT id, template_id, block_name, block_type, block_order, exercise_count, 
			   estimated_duration_minutes, instructions, reps, series, rest_time_seconds, 
			   created_at, updated_at, deleted_at, is_active, created_by
		FROM "%s".custom_template_block 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(fmt.Sprintf(query, gymID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.ResponseCustomTemplateBlockDTO
	for rows.Next() {
		var res dto.ResponseCustomTemplateBlockDTO
		if err := rows.Scan(
			&res.ID,
			&res.TemplateID,
			&res.BlockName,
			&res.BlockType,
			&res.BlockOrder,
			&res.ExerciseCount,
			&res.EstimatedDurationMinutes,
			&res.Instructions,
			&res.Reps,
			&res.Series,
			&res.RestTimeSeconds,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.DeletedAt,
			&res.IsActive,
			&res.CreatedBy,
		); err != nil {
			return nil, err
		}
		result = append(result, &res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteCustomTemplateBlock soft-deletes a custom template block by ID.
func (r *CustomTemplateBlockRepository) DeleteCustomTemplateBlock(gymID, id string) error {
	query := `UPDATE "%s".custom_template_block SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(fmt.Sprintf(query, gymID), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
