package models

type (
	Check struct {
		DocumentNotices        []CheckNotice                `json:"ExecutionWarnings"`
		ArchHasWarnings        bool                         `json:"ArchHasWarnings"`
		ArchWarningsDependency []CheckArchWarningDependency `json:"ArchWarningsDeps"`
		ArchWarningsMatch      []CheckArchWarningMatch      `json:"ArchWarningsNotMatched"`
		ArchWarningsDeepScan   []CheckArchWarningDeepscan   `json:"ArchWarningsDeepScan"`
		OmittedCount           int                          `json:"OmittedCount"`
		ModuleName             string                       `json:"ModuleName"`
	}

	CheckNotice struct {
		Text              string `json:"Text"`
		File              string `json:"File"`
		Line              int    `json:"Line"`
		Offset            int    `json:"Offset"`
		SourceCodePreview []byte `json:"-"`
	}

	CheckArchWarningDependency struct {
		ComponentName      string    `json:"ComponentName"`
		FileRelativePath   string    `json:"FileRelativePath"`
		FileAbsolutePath   string    `json:"FileAbsolutePath"`
		ResolvedImportName string    `json:"ResolvedImportName"`
		Reference          Reference `json:"-"`
		SourceCodePreview  []byte    `json:"-"`
	}

	CheckArchWarningMatch struct {
		FileRelativePath  string    `json:"FileRelativePath"`
		FileAbsolutePath  string    `json:"FileAbsolutePath"`
		Reference         Reference `json:"-"`
		SourceCodePreview []byte    `json:"-"`
	}

	CheckArchWarningDeepscan struct {
		Gate       DeepscanWarningGate
		Dependency DeepscanWarningDependency
		Target     DeepscanWarningTarget
	}

	DeepscanWarningGate struct {
		ComponentName string    // operations
		MethodName    string    // NewOperation
		Definition    Reference // internal/glue/code/line_count.go:54
		RelativePath  string    `json:"-"` // internal/glue/code/line_count.go:54
	}

	DeepscanWarningDependency struct {
		ComponentName     string    // repository
		Name              string    // micro.ViewRepository
		InjectionAST      string    // c.provideMicroViewRepository()
		Injection         Reference // internal/app/internal/container/cmd_mapping.go:15
		InjectionPath     string    `json:"-"` // internal/app/internal/container/cmd_mapping.go:15
		SourceCodePreview []byte    `json:"-"`
	}

	DeepscanWarningTarget struct {
		RelativePath string `json:"-"` // internal/app/internal/container/cmd_mapping.go:15
	}

	CheckResult struct {
		DependencyWarnings []CheckArchWarningDependency
		MatchWarnings      []CheckArchWarningMatch
		DeepscanWarnings   []CheckArchWarningDeepscan
	}
)

func (cr *CheckResult) Append(another CheckResult) {
	cr.DependencyWarnings = append(cr.DependencyWarnings, another.DependencyWarnings...)
	cr.MatchWarnings = append(cr.MatchWarnings, another.MatchWarnings...)
	cr.DeepscanWarnings = append(cr.DeepscanWarnings, another.DeepscanWarnings...)
}

func (cr *CheckResult) HasNotices() bool {
	if len(cr.DependencyWarnings) > 0 {
		return true
	}
	if len(cr.MatchWarnings) > 0 {
		return true
	}
	if len(cr.DeepscanWarnings) > 0 {
		return true
	}

	return false
}
