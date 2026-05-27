UPDATE team 
SET department = COALESCE($1::TEXT, department)
integrations = COALESCE(integrations,'{}'::jsonb) || $2::jsonb 