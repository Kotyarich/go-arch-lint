package check

import (
	"fmt"

	"github.com/fe3dback/go-arch-lint/internal/models"
	"github.com/fe3dback/go-arch-lint/internal/models/speca"
)

type Service struct {
	specAssembler   SpecAssembler
	specChecker     SpecChecker
	referenceRender ReferenceRender
}

func NewService(
	specAssembler SpecAssembler,
	specChecker SpecChecker,
	referenceRender ReferenceRender,
) *Service {
	return &Service{
		specAssembler:   specAssembler,
		specChecker:     specChecker,
		referenceRender: referenceRender,
	}
}

func (s *Service) Behave(maxWarnings int) (models.Check, error) {
	spec, err := s.specAssembler.Assemble()
	if err != nil {
		return models.Check{}, fmt.Errorf("failed to assemble spec: %w", err)
	}

	result := s.specChecker.Check(spec)
	result = s.limitResults(result, maxWarnings)

	model := models.Check{
		ModuleName:             spec.ModuleName.Value(),
		DocumentNotices:        s.assembleNotice(spec.Integrity),
		ArchHasWarnings:        s.resultsHasWarnings(result),
		ArchWarningsDependency: result.DependencyWarnings,
		ArchWarningsMatch:      result.MatchWarnings,
	}

	return model, nil
}

func (s *Service) limitResults(result models.CheckResult, maxWarnings int) models.CheckResult {
	totalCount := 0
	limitedResults := models.CheckResult{
		DependencyWarnings: []models.CheckArchWarningDependency{},
		MatchWarnings:      []models.CheckArchWarningMatch{},
	}

	// append deps
	for _, notice := range result.DependencyWarnings {
		if totalCount >= maxWarnings {
			break
		}

		limitedResults.DependencyWarnings = append(limitedResults.DependencyWarnings, notice)
		totalCount++
	}

	// append not matched
	for _, notice := range result.MatchWarnings {
		if totalCount >= maxWarnings {
			break
		}

		limitedResults.MatchWarnings = append(limitedResults.MatchWarnings, notice)
		totalCount++
	}

	return limitedResults
}

func (s *Service) resultsHasWarnings(result models.CheckResult) bool {
	if len(result.DependencyWarnings) > 0 {
		return true
	}

	if len(result.MatchWarnings) > 0 {
		return true
	}

	return false
}

func (s *Service) assembleNotice(integrity speca.Integrity) []models.CheckNotice {
	notices := make([]speca.Notice, 0)
	notices = append(notices, integrity.DocumentNotices...)
	notices = append(notices, integrity.SpecNotices...)

	results := make([]models.CheckNotice, 0)
	for _, notice := range notices {
		results = append(results, models.CheckNotice{
			Text:              fmt.Sprintf("%s", notice.Notice),
			File:              notice.Ref.File,
			Line:              notice.Ref.Line,
			Offset:            notice.Ref.Offset,
			SourceCodePreview: s.referenceRender.SourceCode(notice.Ref, 1, true),
		})
	}

	return results
}
