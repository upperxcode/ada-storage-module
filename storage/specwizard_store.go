package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type SpecWizard struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Description               string    `json:"description"`
	ExpertLanguagePlugin      string    `json:"expert_language_plugin"`
	PRD                       string    `json:"prd"`
	FunctionalRequirements    string    `json:"functional_requirements"`
	NonFunctionalRequirements string    `json:"non_functional_requirements"`
	Persistence               string    `json:"persistence"`
	Architecture              string    `json:"architecture"`
	EngineeringPhilosophies   string    `json:"engineering_philosophies"`
	DesignPatterns            string    `json:"design_patterns"`
	DataPatterns              string    `json:"data_patterns"`
	StackConfig               string    `json:"stack_config"`
	BusinessStateManagement   string    `json:"business_state_management"`
	BusinessAPIContract       string    `json:"business_api_contract"`
	BusinessCustomizationDetails string `json:"business_customization_details"`
	BusinessFinalAdjustments  string    `json:"business_final_adjustments"`
	BusinessArchitectureRecommendations string `json:"business_architecture_recommendations"`
	Color                     string    `json:"color"`
	Icon                      string    `json:"icon"`
	ArchitectureHealth        int       `json:"architecture_health"`
	DependencyManifest        string    `json:"dependency_manifest"`
	StackPlugin               string    `json:"stack_plugin"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type SpecWizardStore struct {
	db *sql.DB
}

func NewSpecWizardStore(db *sql.DB) *SpecWizardStore {
	return &SpecWizardStore{db: db}
}

func (s *SpecWizardStore) Create(ctx context.Context, w *SpecWizard) error {
	query := `INSERT INTO spec_wizards (
		id, name, description, expert_language_plugin, prd,
		functional_requirements, non_functional_requirements,
		persistence, architecture, engineering_philosophies,
		design_patterns, data_patterns, stack_config,
		business_state_management, business_api_contract,
		business_customization_details, business_final_adjustments,
		business_architecture_recommendations, color, icon,
		architecture_health, dependency_manifest, stack_plugin,
		created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		w.ID, w.Name, w.Description, w.ExpertLanguagePlugin, w.PRD,
		w.FunctionalRequirements, w.NonFunctionalRequirements,
		w.Persistence, w.Architecture, w.EngineeringPhilosophies,
		w.DesignPatterns, w.DataPatterns, w.StackConfig,
		w.BusinessStateManagement, w.BusinessAPIContract,
		w.BusinessCustomizationDetails, w.BusinessFinalAdjustments,
		w.BusinessArchitectureRecommendations, w.Color, w.Icon,
		w.ArchitectureHealth, w.DependencyManifest, w.StackPlugin,
		w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create spec wizard: %w", err)
	}
	return nil
}

func (s *SpecWizardStore) Get(ctx context.Context, id string) (*SpecWizard, error) {
	query := `SELECT id, name, description, expert_language_plugin, prd,
		functional_requirements, non_functional_requirements,
		persistence, architecture, engineering_philosophies,
		design_patterns, data_patterns, stack_config,
		business_state_management, business_api_contract,
		business_customization_details, business_final_adjustments,
		business_architecture_recommendations, color, icon,
		architecture_health, dependency_manifest, stack_plugin,
		created_at, updated_at
		FROM spec_wizards WHERE id = ?`

	var w SpecWizard
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&w.ID, &w.Name, &w.Description, &w.ExpertLanguagePlugin, &w.PRD,
		&w.FunctionalRequirements, &w.NonFunctionalRequirements,
		&w.Persistence, &w.Architecture, &w.EngineeringPhilosophies,
		&w.DesignPatterns, &w.DataPatterns, &w.StackConfig,
		&w.BusinessStateManagement, &w.BusinessAPIContract,
		&w.BusinessCustomizationDetails, &w.BusinessFinalAdjustments,
		&w.BusinessArchitectureRecommendations, &w.Color, &w.Icon,
		&w.ArchitectureHealth, &w.DependencyManifest, &w.StackPlugin,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get spec wizard: %w", err)
	}
	return &w, nil
}

func (s *SpecWizardStore) List(ctx context.Context) ([]SpecWizard, error) {
	query := `SELECT id, name, description, expert_language_plugin, prd,
		functional_requirements, non_functional_requirements,
		persistence, architecture, engineering_philosophies,
		design_patterns, data_patterns, stack_config,
		business_state_management, business_api_contract,
		business_customization_details, business_final_adjustments,
		business_architecture_recommendations, color, icon,
		architecture_health, dependency_manifest, stack_plugin,
		created_at, updated_at
		FROM spec_wizards ORDER BY updated_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query spec wizards: %w", err)
	}
	defer rows.Close()

	var wizards []SpecWizard
	for rows.Next() {
		var w SpecWizard
		if err := rows.Scan(
			&w.ID, &w.Name, &w.Description, &w.ExpertLanguagePlugin, &w.PRD,
			&w.FunctionalRequirements, &w.NonFunctionalRequirements,
			&w.Persistence, &w.Architecture, &w.EngineeringPhilosophies,
			&w.DesignPatterns, &w.DataPatterns, &w.StackConfig,
			&w.BusinessStateManagement, &w.BusinessAPIContract,
			&w.BusinessCustomizationDetails, &w.BusinessFinalAdjustments,
			&w.BusinessArchitectureRecommendations, &w.Color, &w.Icon,
			&w.ArchitectureHealth, &w.DependencyManifest, &w.StackPlugin,
			&w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan spec wizard: %w", err)
		}
		wizards = append(wizards, w)
	}
	return wizards, nil
}

func (s *SpecWizardStore) Update(ctx context.Context, w *SpecWizard) error {
	query := `UPDATE spec_wizards SET
		name = ?, description = ?, expert_language_plugin = ?, prd = ?,
		functional_requirements = ?, non_functional_requirements = ?,
		persistence = ?, architecture = ?, engineering_philosophies = ?,
		design_patterns = ?, data_patterns = ?, stack_config = ?,
		business_state_management = ?, business_api_contract = ?,
		business_customization_details = ?, business_final_adjustments = ?,
		business_architecture_recommendations = ?, color = ?, icon = ?,
		architecture_health = ?, dependency_manifest = ?, stack_plugin = ?,
		updated_at = ?
		WHERE id = ?`

	_, err := s.db.ExecContext(ctx, query,
		w.Name, w.Description, w.ExpertLanguagePlugin, w.PRD,
		w.FunctionalRequirements, w.NonFunctionalRequirements,
		w.Persistence, w.Architecture, w.EngineeringPhilosophies,
		w.DesignPatterns, w.DataPatterns, w.StackConfig,
		w.BusinessStateManagement, w.BusinessAPIContract,
		w.BusinessCustomizationDetails, w.BusinessFinalAdjustments,
		w.BusinessArchitectureRecommendations, w.Color, w.Icon,
		w.ArchitectureHealth, w.DependencyManifest, w.StackPlugin,
		w.UpdatedAt, w.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update spec wizard: %w", err)
	}
	return nil
}

func (s *SpecWizardStore) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM spec_wizards WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete spec wizard: %w", err)
	}
	return nil
}

func MarshalStringSlice(v []string) string {
	if v == nil {
		return "[]"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func UnmarshalStringSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return []string{}
	}
	return result
}

func MarshalStackConfig(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func UnmarshalStackConfig(s string) []map[string]interface{} {
	if s == "" {
		return []map[string]interface{}{}
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return []map[string]interface{}{}
	}
	return result
}
