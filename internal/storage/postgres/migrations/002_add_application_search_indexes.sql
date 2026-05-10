CREATE INDEX IF NOT EXISTS idx_applications_source_lower
    ON applications (LOWER(source));

CREATE INDEX IF NOT EXISTS idx_applications_company_name_lower
    ON applications (LOWER(company_name));

CREATE INDEX IF NOT EXISTS idx_applications_title_lower
    ON applications (LOWER(title));