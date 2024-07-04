package service

import (
	"context"
	"time"

	pb "visualization/api/proto"
	"visualization/internal/processor"
	"visualization/renderers"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VisualizationStore defines the interface for storing and retrieving visualizations.
type VisualizationStore interface {
	// Store saves a visualization and returns its ID.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - visualization: The visualization to be stored.
	//
	// Returns:
	//   - string: The ID of the stored visualization.
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   id, err := store.Store(ctx, &pb.VisualizationResponse{
	//     VisualizationData: []byte("..."),
	//     Metadata: map[string]string{"type": "bar_chart"},
	//   })
	//   if err != nil {
	//     log.Printf("Failed to store visualization: %v", err)
	//     return
	//   }
	//   log.Printf("Stored visualization with ID: %s", id)
	Store(ctx context.Context, visualization *pb.VisualizationResponse) (string, error)

	// CreateVisualization creates a new visualization based on the provided request.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - req: The request containing the data and metadata for the visualization.
	//
	// Returns:
	//   - *pb.VisualizationResponse: The created visualization.
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   viz, err := store.CreateVisualization(ctx, &pb.CreateVisualizationRequest{
	//     DataSourceId: "123e4567-e89b-12d3-a456-426614174000",
	//     Dimensions: []string{"country"},
	//     Measures: []string{"revenue"},
	//     VisualizationType: "bar_chart",
	//   })
	CreateVisualization(ctx context.Context, req *pb.CreateVisualizationRequest) (*pb.VisualizationResponse, error)

	// UpdateVisualization updates an existing visualization.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - req: The request containing the ID of the visualization to update and the new data.
	//
	// Returns:
	//   - *pb.VisualizationResponse: The updated visualization.
	//   - error: An error if the operation failed, or nil if it succeeded.
	// Example:
	//   viz, err := store.UpdateVisualization(ctx, &pb.UpdateVisualizationRequest{
	//     VisualizationId: "123e4567-e89b-12d3-a456-426614174000",
	//     UpdateData: &pb.CreateVisualizationRequest{
	//       DataSourceId: "123e4567-e89b-12d3-a456-426614174000",
	//       Dimensions: []string{"country"},
	//       Measures: []string{"revenue"},
	//       VisualizationType: "bar_chart",
	//     },
	//   })
	UpdateVisualization(ctx context.Context, req *pb.UpdateVisualizationRequest) (*pb.VisualizationResponse, error)

	// Get retrieves a visualization by its ID.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - id: The ID of the visualization to retrieve.
	//
	// Returns:
	//   - *pb.VisualizationResponse: The retrieved visualization, or nil if not found.
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   viz, err := store.Get(ctx, "123e4567-e89b-12d3-a456-426614174000")
	//   if err != nil {
	//     log.Printf("Failed to retrieve visualization: %v", err)
	//     return
	//   }
	//   if viz == nil {
	//     log.Printf("Visualization not found")
	//   } else {
	//     log.Printf("Retrieved visualization of type: %s", viz.Metadata["type"])
	//   }
	Get(ctx context.Context, id string) (*pb.VisualizationResponse, error)

	// Update updates an existing visualization.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - visualization: The updated visualization. The ID field must be set.
	//
	// Returns:
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   err := store.Update(ctx, &pb.VisualizationResponse{
	//     VisualizationId: "123e4567-e89b-12d3-a456-426614174000",
	//     VisualizationData: []byte("..."),
	//     Metadata: map[string]string{"type": "line_chart"},
	//   })
	//   if err != nil {
	//     log.Printf("Failed to update visualization: %v", err)
	//     return
	//   }
	//   log.Printf("Successfully updated visualization")
	Update(ctx context.Context, visualization *pb.VisualizationResponse) error

	// Delete removes a visualization by its ID.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - id: The ID of the visualization to delete.
	//
	// Returns:
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   err := store.Delete(ctx, "123e4567-e89b-12d3-a456-426614174000")
	//   if err != nil {
	//     log.Printf("Failed to delete visualization: %v", err)
	//     return
	//   }
	//   log.Printf("Successfully deleted visualization")
	Delete(ctx context.Context, id string) error

	// List retrieves a list of visualizations with pagination.
	//
	// Parameters:
	//   - ctx: The context for the operation, which can be used for cancellation or timeout.
	//   - page: The page number to retrieve (1-based).
	//   - pageSize: The number of items per page.
	//
	// Returns:
	//   - []*pb.VisualizationResponse: A slice of visualizations for the requested page.
	//   - int32: The total count of visualizations (across all pages).
	//   - error: An error if the operation failed, or nil if it succeeded.
	//
	// Example:
	//   vizs, total, err := store.List(ctx, 1, 10)
	//   if err != nil {
	//     log.Printf("Failed to list visualizations: %v", err)
	//     return
	//   }
	//   log.Printf("Retrieved %d visualizations (total: %d)", len(vizs), total)
	//   for _, viz := range vizs {
	//     log.Printf("Visualization ID: %s, Type: %s", viz.VisualizationId, viz.Metadata["type"])
	//   }
	List(ctx context.Context, page, pageSize int32) ([]*pb.VisualizationResponse, int32, error)
}


// VisualizationService implements the VisualizationService gRPC service.
type VisualizationService struct {
	pb.UnimplementedVisualizationServiceServer
	processor          *processor.DataProcessor
	visualizationStore VisualizationStore
	logger             *zap.Logger
}

// NewVisualizationService creates a new VisualizationService with the given dependencies.
func NewVisualizationService(processor *processor.DataProcessor, store VisualizationStore, logger *zap.Logger) *VisualizationService {
	return &VisualizationService{
		processor:          processor,
		visualizationStore: store,
		logger:             logger,
	}
}

// CreateVisualization creates a new visualization based on the provided request.
func (s *VisualizationService) CreateVisualization(ctx context.Context, req *pb.CreateVisualizationRequest) (*pb.VisualizationResponse, error) {
	s.logger.Info("Creating new visualization", zap.String("type", req.VisualizationType))

	processedData, err := s.processor.ProcessData(ctx, req)
	if err != nil {
		s.logger.Error("Failed to process data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to process data: %v", err)
	}

	renderer, err := renderers.RendererFactory(req.VisualizationType)
	if err != nil {
		s.logger.Error("Failed to create renderer", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to create renderer: %v", err)
	}

	visualizationData, err := renderer.Render(processedData)
	if err != nil {
		s.logger.Error("Failed to render visualization", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to render visualization: %v", err)
	}

	viz := &pb.VisualizationResponse{
		VisualizationId:   uuid.New().String(),
		VisualizationData: []byte(visualizationData),
		Metadata: map[string]string{
			"type":       req.VisualizationType,
			"created_at": time.Now().UTC().Format(time.RFC3339),
		},
	}

	id, err := s.visualizationStore.Store(ctx, viz)
	if err != nil {
		s.logger.Error("Failed to store visualization", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to store visualization: %v", err)
	}

	viz.VisualizationId = id
	s.logger.Info("Successfully created visualization", zap.String("id", id))
	return viz, nil
}

// UpdateVisualization updates an existing visualization.
func (s *VisualizationService) UpdateVisualization(ctx context.Context, req *pb.UpdateVisualizationRequest) (*pb.VisualizationResponse, error) {
	s.logger.Info("Updating visualization", zap.String("id", req.VisualizationId))

	existing, err := s.visualizationStore.Get(ctx, req.VisualizationId)
	if err != nil {
		s.logger.Error("Failed to retrieve existing visualization", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "visualization not found: %v", err)
	}

	processedData, err := s.processor.ProcessData(ctx, req.UpdateData)
	if err != nil {
		s.logger.Error("Failed to process updated data", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to process updated data: %v", err)
	}

	renderer, err := renderers.RendererFactory(req.UpdateData.VisualizationType)
	if err != nil {
		s.logger.Error("Failed to create renderer", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to create renderer: %v", err)
	}

	visualizationData, err := renderer.Render(processedData)
	if err != nil {
		s.logger.Error("Failed to render updated visualization", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to render updated visualization: %v", err)
	}

	existing.VisualizationData = []byte(visualizationData)
	existing.Metadata["type"] = req.UpdateData.VisualizationType
	existing.Metadata["updated_at"] = time.Now().UTC().Format(time.RFC3339)

	err = s.visualizationStore.Update(ctx, existing)
	if err != nil {
		s.logger.Error("Failed to store updated visualization", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update visualization: %v", err)
	}

	s.logger.Info("Successfully updated visualization", zap.String("id", req.VisualizationId))
	return existing, nil
}

// GetVisualization retrieves a specific visualization by its ID.
func (s *VisualizationService) GetVisualization(ctx context.Context, req *pb.GetVisualizationRequest) (*pb.VisualizationResponse, error) {
	s.logger.Info("Retrieving visualization", zap.String("id", req.VisualizationId))

	visualization, err := s.visualizationStore.Get(ctx, req.VisualizationId)
	if err != nil {
		s.logger.Error("Failed to retrieve visualization", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "visualization not found: %v", err)
	}

	s.logger.Info("Successfully retrieved visualization", zap.String("id", req.VisualizationId))
	return visualization, nil
}

// ListVisualizations retrieves a list of visualizations with pagination.
func (s *VisualizationService) ListVisualizations(ctx context.Context, req *pb.ListVisualizationsRequest) (*pb.ListVisualizationsResponse, error) {
	s.logger.Info("Listing visualizations", zap.Int32("page", req.Page), zap.Int32("pageSize", req.PageSize))

	visualizations, totalCount, err := s.visualizationStore.List(ctx, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Failed to list visualizations", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list visualizations: %v", err)
	}

	s.logger.Info("Successfully listed visualizations", zap.Int("count", len(visualizations)), zap.Int32("totalCount", totalCount))
	return &pb.ListVisualizationsResponse{
		Visualizations: visualizations,
		TotalCount:     totalCount,
	}, nil
}

// DeleteVisualization removes a specific visualization by its ID.
func (s *VisualizationService) DeleteVisualization(ctx context.Context, req *pb.DeleteVisualizationRequest) (*pb.DeleteVisualizationResponse, error) {
	s.logger.Info("Deleting visualization", zap.String("id", req.VisualizationId))

	err := s.visualizationStore.Delete(ctx, req.VisualizationId)
	if err != nil {
		s.logger.Error("Failed to delete visualization", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete visualization: %v", err)
	}

	s.logger.Info("Successfully deleted visualization", zap.String("id", req.VisualizationId))
	return &pb.DeleteVisualizationResponse{Success: true}, nil
}

// ExportVisualization exports a visualization in the specified format.
func (s *VisualizationService) ExportVisualization(ctx context.Context, req *pb.ExportVisualizationRequest) (*pb.ExportVisualizationResponse, error) {
	s.logger.Info("Exporting visualization", zap.String("id", req.VisualizationId), zap.String("format", req.Format))

	visualization, err := s.visualizationStore.Get(ctx, req.VisualizationId)
	if err != nil {
		s.logger.Error("Failed to retrieve visualization for export", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "visualization not found: %v", err)
	}

	var exportedData []byte
	switch req.Format {
	case "raw":
		exportedData = visualization.VisualizationData
	case "png", "svg", "pdf":
		// TODO: Implement conversion to these formats
		return nil, status.Errorf(codes.Unimplemented, "export format %s not yet implemented", req.Format)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported export format: %s", req.Format)
	}

	s.logger.Info("Successfully exported visualization", zap.String("id", req.VisualizationId), zap.String("format", req.Format))
	return &pb.ExportVisualizationResponse{
		ExportedData: exportedData,
		Format:       req.Format,
	}, nil
}

// Close closes any resources held by the VisualizationService.
func (s *VisualizationService) Close() error {
	return s.processor.Close()
}