package repository

import (
	"database/sql"

	"github.com/alejandro-albiol/athenai/internal/custom_template_block/dto"
)

type CustomTemplateBlockRepositoryImpl struct {
	DB *sql.DB
}

func NewCustomTemplateBlockRepository(db *sql.DB) *CustomTemplateBlockRepositoryImpl {
	return &CustomTemplateBlockRepositoryImpl{DB: db}
}

func (r *CustomTemplateBlockRepositoryImpl) Create(gymID string, block *dto.CreateCustomTemplateBlockDTO) (string, error) {
	// SQL: INSERT INTO <gymID>.custom_template_block ...
	return "", nil
}

func (r *CustomTemplateBlockRepositoryImpl) GetByID(gymID, id string) (*dto.ResponseCustomTemplateBlockDTO, error) {
	// SQL: SELECT ... FROM <gymID>.custom_template_block WHERE id = $1
	return nil, nil
}

func (r *CustomTemplateBlockRepositoryImpl) ListByTemplateID(gymID, templateID string) ([]*dto.ResponseCustomTemplateBlockDTO, error) {
	// SQL: SELECT ... FROM <gymID>.custom_template_block WHERE template_id = $1
	return nil, nil
}

func (r *CustomTemplateBlockRepositoryImpl) Update(gymID string, block *dto.UpdateCustomTemplateBlockDTO) error {
	// SQL: UPDATE <gymID>.custom_template_block SET ... WHERE id = $1
	return nil
}

func (r *CustomTemplateBlockRepositoryImpl) Delete(gymID, id string) error {
	// SQL: DELETE FROM <gymID>.custom_template_block WHERE id = $1
	return nil
}
