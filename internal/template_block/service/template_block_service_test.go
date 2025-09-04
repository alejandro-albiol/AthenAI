package service_test

import (
	"errors"
	"testing"

	"github.com/alejandro-albiol/athenai/internal/template_block/dto"
	"github.com/alejandro-albiol/athenai/internal/template_block/service"
	"github.com/alejandro-albiol/athenai/pkg/apierror"
)

type mockRepository struct {
	blocks    map[string]*dto.TemplateBlockDTO
	createErr error
	getErr    error
	listErr   error
	updateErr error
	deleteErr error
}

func (m *mockRepository) CreateTemplateBlock(block *dto.CreateTemplateBlockDTO) (*string, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	id := "new-id"
	m.blocks[id] = &dto.TemplateBlockDTO{
		ID:         id,
		Name:       block.Name,
		TemplateID: block.TemplateID,
	}
	return &id, nil
}

func (m *mockRepository) GetTemplateBlockByTemplateIDAndName(templateID, name string) (*dto.TemplateBlockDTO, error) {
	for _, b := range m.blocks {
		if b.TemplateID == templateID && b.Name == name {
			return b, nil
		}
	}
	return &dto.TemplateBlockDTO{}, nil
}

func (m *mockRepository) GetTemplateBlockByID(id string) (*dto.TemplateBlockDTO, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	b, ok := m.blocks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return b, nil
}

func (m *mockRepository) GetTemplateBlocksByTemplateID(templateID string) ([]*dto.TemplateBlockDTO, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	var result []*dto.TemplateBlockDTO
	for _, b := range m.blocks {
		if b.TemplateID == templateID {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *mockRepository) UpdateTemplateBlock(id string, update *dto.UpdateTemplateBlockDTO) (*dto.TemplateBlockDTO, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	b, ok := m.blocks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	if update.Name != nil {
		b.Name = *update.Name
	}
	return b, nil
}

func (m *mockRepository) DeleteTemplateBlock(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	if _, ok := m.blocks[id]; !ok {
		return errors.New("not found")
	}
	delete(m.blocks, id)
	return nil
}

func TestTemplateBlockService_CreateTemplateBlock(t *testing.T) {
	mock := &mockRepository{blocks: map[string]*dto.TemplateBlockDTO{}}
	service := service.NewTemplateBlockService(mock)
	block := &dto.CreateTemplateBlockDTO{Name: "Block1", TemplateID: "T1"}
	id, err := service.CreateTemplateBlock(block)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || *id == "" {
		t.Fatalf("expected id, got empty string")
	}

	// Duplicate name
	mock.blocks[*id] = &dto.TemplateBlockDTO{ID: *id, Name: "Block1", TemplateID: "T1"}
	_, err = service.CreateTemplateBlock(block)
	if err == nil {
		t.Fatalf("expected conflict error, got nil")
	}

	if apierr, ok := err.(*apierror.APIError); !ok || apierr.Code != "CONFLICT" {
		t.Fatalf("expected conflict error code, got %v", err)
	}
}

func TestTemplateBlockService_GetTemplateBlockByID(t *testing.T) {
	mock := &mockRepository{blocks: map[string]*dto.TemplateBlockDTO{"id1": {ID: "id1", Name: "Block1", TemplateID: "T1"}}}
	service := service.NewTemplateBlockService(mock)
	block, err := service.GetTemplateBlockByID("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.ID != "id1" {
		t.Fatalf("expected id1, got %v", block.ID)
	}

	// Not found
	_, err = service.GetTemplateBlockByID("notfound")
	if err == nil {
		t.Fatalf("expected not found error, got nil")
	}
	if apierr, ok := err.(*apierror.APIError); !ok || apierr.Code != "NOT_FOUND" {
		t.Fatalf("expected not_found error code, got %v", err)
	}
}

func TestTemplateBlockService_ListTemplateBlocksByTemplateID(t *testing.T) {
	mock := &mockRepository{blocks: map[string]*dto.TemplateBlockDTO{
		"id1": {ID: "id1", Name: "Block1", TemplateID: "T1"},
		"id2": {ID: "id2", Name: "Block2", TemplateID: "T1"},
		"id3": {ID: "id3", Name: "Block3", TemplateID: "T2"},
	}}
	service := service.NewTemplateBlockService(mock)
	blocks, err := service.ListTemplateBlocksByTemplateID("T1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	mock.listErr = errors.New("fail")
	_, err = service.ListTemplateBlocksByTemplateID("T1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTemplateBlockService_UpdateTemplateBlock(t *testing.T) {
	mock := &mockRepository{blocks: map[string]*dto.TemplateBlockDTO{"id1": {ID: "id1", Name: "Block1", TemplateID: "T1"}}}
	service := service.NewTemplateBlockService(mock)
	newName := "Updated"
	update := &dto.UpdateTemplateBlockDTO{Name: &newName}
	block, err := service.UpdateTemplateBlock("id1", update)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Name != "Updated" {
		t.Fatalf("expected Updated, got %v", block.Name)
	}

	mock.updateErr = errors.New("fail")
	_, err = service.UpdateTemplateBlock("id1", update)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestTemplateBlockService_DeleteTemplateBlock(t *testing.T) {
	mock := &mockRepository{blocks: map[string]*dto.TemplateBlockDTO{"id1": {ID: "id1", Name: "Block1", TemplateID: "T1"}}}
	service := service.NewTemplateBlockService(mock)
	err := service.DeleteTemplateBlock("id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := mock.blocks["id1"]; ok {
		t.Fatalf("expected block to be deleted")
	}

	mock.deleteErr = errors.New("fail")
	err = service.DeleteTemplateBlock("id1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
