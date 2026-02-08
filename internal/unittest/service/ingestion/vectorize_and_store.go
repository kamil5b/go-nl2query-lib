package ingestion

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/kamil5b/go-nl-sql/internal/domain/service"
	"github.com/kamil5b/go-nl-sql/mocks"
	"github.com/stretchr/testify/require"
)

// @example usage:
//
//	func TestIngestionService_VectorizeAndStore(t *testing.T) {
//	    ingestion.UnitTestVectorizeAndStore(t, New(embedderRepo, vectorStoreRepo, statusRepo))
//	}
func UnitTestVectorizeAndStore(
	t *testing.T,
	svcImp func(
		embedderRepo *mocks.MockEmbedderRepository,
		vectorStoreRepo *mocks.MockVectorStoreRepository,
		statusRepo *mocks.MockStatusRepository,
	) service.IngestionService,
) {
	var (
		mockEmbedderRepo    *mocks.MockEmbedderRepository
		mockVectorStoreRepo *mocks.MockVectorStoreRepository
		mockStatusRepo      *mocks.MockStatusRepository
	)

	mockMetaData := &model.DatabaseMetadata{
		TenantID: "tenant_abc",
		Checksum: "sha256:7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
		Tables: []model.Table{
			{
				// CASE 1: Table with Composite Primary Keys & Self-Reference
				Name: "employees",
				Columns: []model.Column{
					{Name: "org_id", Type: "INT", Nullable: false, IsPrimaryKey: true},
					{Name: "emp_id", Type: "INT", Nullable: false, IsPrimaryKey: true},
					{Name: "manager_id", Type: "INT", Nullable: true, IsForeignKey: true, Comments: "Self-reference"},
					{Name: "email", Type: "VARCHAR(255)", Nullable: false},
					{Name: "status", Type: "VARCHAR(20)", Default: "'active'"},
				},
				Indexes: []model.Index{
					{
						Name:    "idx_org_email",
						Columns: []string{"org_id", "email"}, // Multi-column index
						Unique:  true,
					},
				},
				Constraints: []model.Constraint{
					{
						Name:      "fk_manager",
						Type:      "FOREIGN KEY",
						Columns:   []string{"manager_id"},
						Reference: "employees(emp_id)",
					},
				},
				Comments: "Employee records with multi-tenant composite keys",
			},
			{
				// CASE 2: Join Table (Many-to-Many resolution)
				Name: "project_assignments",
				Columns: []model.Column{
					{Name: "project_id", Type: "INT", IsPrimaryKey: true, IsForeignKey: true},
					{Name: "emp_id", Type: "INT", IsPrimaryKey: true, IsForeignKey: true},
					{Name: "assigned_at", Type: "TIMESTAMP", Default: "CURRENT_TIMESTAMP"},
				},
			},
		},
		Relations: []model.Relation{
			{
				// Self-referential Relation
				SourceTable:  "employees",
				SourceColumn: "manager_id",
				TargetTable:  "employees",
				TargetColumn: "emp_id",
				RelationType: "MANY_TO_ONE",
			},
			{
				// Standard Foreign Relation
				SourceTable:  "project_assignments",
				SourceColumn: "emp_id",
				TargetTable:  "employees",
				TargetColumn: "emp_id",
				RelationType: "MANY_TO_ONE",
			},
		},
	}

	mockVector := [][]float32{{0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}}

	tests := []struct {
		name        string
		metadata    *model.DatabaseMetadata
		prepareMock func()
		expectError error
	}{
		{
			name:     "success",
			metadata: mockMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), mockMetaData.TenantID).
					Return(nil)

				mockEmbedderRepo.EXPECT().
					EmbedBatch(gomock.Any(), gomock.Any()).
					Return(mockVector, nil)

				mockVectorStoreRepo.EXPECT().
					Upsert(gomock.Any(), mockMetaData.TenantID, mockVector).
					Return(nil)

				mockStatusRepo.EXPECT().
					SetDone(gomock.Any(), mockMetaData.TenantID).
					Return(nil)
			},
			expectError: nil,
		},
		{
			name:     "error set status done",
			metadata: mockMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), mockMetaData.TenantID).
					Return(nil)

				mockEmbedderRepo.EXPECT().
					EmbedBatch(gomock.Any(), gomock.Any()).
					Return(mockVector, nil)

				mockVectorStoreRepo.EXPECT().
					Upsert(gomock.Any(), mockMetaData.TenantID, mockVector).
					Return(nil)

				mockStatusRepo.EXPECT().
					SetDone(gomock.Any(), mockMetaData.TenantID).
					Return(errors.New("some error"))
			},
			expectError: errors.New("some error"),
		},
		{
			name:     "error upsert vector",
			metadata: mockMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), mockMetaData.TenantID).
					Return(nil)

				mockEmbedderRepo.EXPECT().
					EmbedBatch(gomock.Any(), gomock.Any()).
					Return(mockVector, nil)

				mockVectorStoreRepo.EXPECT().
					Upsert(gomock.Any(), mockMetaData.TenantID, mockVector).
					Return(errors.New("some error"))

				mockStatusRepo.EXPECT().
					SetError(gomock.Any(), mockMetaData.TenantID, errors.New("some error").Error()).
					Return(nil)
			},
			expectError: errors.New("some error"),
		},
		{
			name:     "error embed vector",
			metadata: mockMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), mockMetaData.TenantID).
					Return(nil)

				mockEmbedderRepo.EXPECT().
					EmbedBatch(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("some error"))

				mockStatusRepo.EXPECT().
					SetError(gomock.Any(), mockMetaData.TenantID, errors.New("some error").Error()).
					Return(nil)
			},
			expectError: errors.New("some error"),
		},
		{
			name:     "error status in progress",
			metadata: mockMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), mockMetaData.TenantID).
					Return(errors.New("some error"))
			},
			expectError: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmbedderRepo = mocks.NewMockEmbedderRepository(ctrl)
			mockVectorStoreRepo = mocks.NewMockVectorStoreRepository(ctrl)
			mockStatusRepo = mocks.NewMockStatusRepository(ctrl)

			svc := svcImp(
				mockEmbedderRepo,
				mockVectorStoreRepo,
				mockStatusRepo,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			err := svc.VectorizeAndStore(context.Background(), tt.metadata)

			if tt.expectError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
