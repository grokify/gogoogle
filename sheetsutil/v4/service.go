package sheetsutil

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Service wraps the Google Sheets API service.
type Service struct {
	httpClient *http.Client
	// SheetsService is the underlying Sheets API service.
	SheetsService *sheets.Service
	// SpreadsheetsService is SheetsService.Spreadsheets.
	SpreadsheetsService *sheets.SpreadsheetsService
}

// NewService creates a Sheets service from an authenticated HTTP client.
// Additional options (for example option.WithEndpoint) are applied after the
// HTTP client.
func NewService(ctx context.Context, c *http.Client, opts ...option.ClientOption) (*Service, error) {
	all := make([]option.ClientOption, 0, len(opts)+1)
	all = append(all, option.WithHTTPClient(c))
	all = append(all, opts...)
	svc, err := sheets.NewService(ctx, all...)
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: new service: %w", err)
	}
	return &Service{httpClient: c, SheetsService: svc, SpreadsheetsService: svc.Spreadsheets}, nil
}

// GetSpreadsheet retrieves spreadsheet metadata by ID.
func (s *Service) GetSpreadsheet(ctx context.Context, spreadsheetID string) (*sheets.Spreadsheet, error) {
	ss, err := s.SpreadsheetsService.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: get spreadsheet %s: %w", spreadsheetID, err)
	}
	return ss, nil
}

// GetValues reads one A1 range. renderOption is "FORMATTED_VALUE" (the
// default when empty), "UNFORMATTED_VALUE", or "FORMULA".
func (s *Service) GetValues(ctx context.Context, spreadsheetID, a1Range, renderOption string) (*sheets.ValueRange, error) {
	call := s.SpreadsheetsService.Values.Get(spreadsheetID, a1Range).Context(ctx)
	if renderOption != "" {
		call = call.ValueRenderOption(renderOption)
	}
	vr, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: get values %s %q: %w", spreadsheetID, a1Range, err)
	}
	return vr, nil
}

// BatchGetValues reads several A1 ranges in one request. renderOption is as
// for GetValues.
func (s *Service) BatchGetValues(ctx context.Context, spreadsheetID string, ranges []string, renderOption string) (*sheets.BatchGetValuesResponse, error) {
	call := s.SpreadsheetsService.Values.BatchGet(spreadsheetID).Ranges(ranges...).Context(ctx)
	if renderOption != "" {
		call = call.ValueRenderOption(renderOption)
	}
	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("sheetsutil: batch get values %s: %w", spreadsheetID, err)
	}
	return resp, nil
}
