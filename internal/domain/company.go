package domain

// -----------------------------------------------------------------------------
// Company
//
// Represents an organization connected to one or more job applications.
// -----------------------------------------------------------------------------
type Company struct {
	Name    string
	Website string
}

// -----------------------------------------------------------------------------
// NewCompany
//
// Creates a company after applying basic domain validation.
// -----------------------------------------------------------------------------
func NewCompany(name string, website string) (Company, error) {
	company := Company{
		Name:    name,
		Website: website,
	}

	if err := company.Validate(); err != nil {
		return Company{}, err
	}

	return company, nil
}

// -----------------------------------------------------------------------------
// Validate
//
// Verifies that a company has the minimum information required by the domain.
// -----------------------------------------------------------------------------
func (company Company) Validate() error {
	return requireNonEmptyField("company name", company.Name)
}
