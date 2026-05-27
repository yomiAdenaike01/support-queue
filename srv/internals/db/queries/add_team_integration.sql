UPDATE team 
SET integrations = COALESCE(integrations, '{}'::jsonb) || $1
WHERE id = $1
RETURNING id, pk, integrations