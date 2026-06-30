SELECT 
    COUNT(*) as tickets,
    EXTRACT(hour FROM g.day) as hour
FROM 
    generate_series($1::TIMESTAMP, $2::TIMESTAMP, '1 hour'::INTERVAL) AS g(day)
LEFT JOIN 
    tickets t 
ON 
    date_trunc('day', t.created_at) = date_trunc('day', g.day)
GROUP BY g.day
ORDER BY g.day ASC