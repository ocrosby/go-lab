// Package utils provides common utilities for the User Management API.
package utils

import (
	"errors"
	"strconv"

	"github.com/ocrosby/go-lab/projects/api/internal/config"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Limit  int
	Offset int
}

// NewPaginationParams creates new pagination parameters with validation
func NewPaginationParams(limit, offset int) PaginationParams {
	if limit < 1 {
		limit = config.DefaultPaginationLimit
	}
	if offset < 0 {
		offset = config.DefaultPaginationOffset
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

// ParsePaginationFromQuery extracts pagination parameters from query strings
func ParsePaginationFromQuery(limitStr, offsetStr string) PaginationParams {
	limit := config.DefaultPaginationLimit
	offset := config.DefaultPaginationOffset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return PaginationParams{
		Limit:  limit,
		Offset: offset,
	}
}

// ValidatePagination validates pagination parameters
func (p PaginationParams) Validate() error {
	if p.Limit < 1 {
		return ErrInvalidLimit
	}
	if p.Offset < 0 {
		return ErrInvalidOffset
	}
	return nil
}

// HasNext checks if there might be more results after the current page
func (p PaginationParams) HasNext(currentCount int) bool {
	return currentCount == p.Limit
}

// CalculateSkip calculates how many records to skip for the current page
func (p PaginationParams) CalculateSkip() int {
	return p.Offset
}

// CalculatePage calculates the current page number (1-based)
func (p PaginationParams) CalculatePage() int {
	if p.Limit == 0 {
		return 1
	}
	return (p.Offset / p.Limit) + 1
}

// PaginationResponse represents a paginated response metadata
type PaginationResponse struct {
	Limit    int  `json:"limit"`
	Offset   int  `json:"offset"`
	Total    *int `json:"total,omitempty"` // Optional total count
	HasNext  bool `json:"has_next"`
	HasPrev  bool `json:"has_prev"`
	Page     int  `json:"page"`
	NextPage *int `json:"next_page,omitempty"`
	PrevPage *int `json:"prev_page,omitempty"`
}

// NewPaginationResponse creates pagination metadata for responses
func NewPaginationResponse(params PaginationParams, currentCount int, total *int) PaginationResponse {
	hasNext := params.HasNext(currentCount)
	hasPrev := params.Offset > 0
	page := params.CalculatePage()

	response := PaginationResponse{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Total:   total,
		HasNext: hasNext,
		HasPrev: hasPrev,
		Page:    page,
	}

	if hasNext {
		nextPage := page + 1
		response.NextPage = &nextPage
	}

	if hasPrev && page > 1 {
		prevPage := page - 1
		response.PrevPage = &prevPage
	}

	return response
}

// Common pagination errors
var (
	ErrInvalidLimit  = errors.New("invalid limit: must be greater than 0")
	ErrInvalidOffset = errors.New("invalid offset: must be greater than or equal to 0")
)
