SELECT
    COUNT(*) FILTER (
        WHERE created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS total_tickets_today,

    COUNT(*) FILTER (
        WHERE created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS total_tickets_yday,

    COUNT(*) FILTER (
        WHERE status = 'PENDING'
          AND created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS pending_tickets_today,

    COUNT(*) FILTER (
        WHERE status = 'PROCESSING'
          AND created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS processing_tickets_today,

    COUNT(*) FILTER (
        WHERE status = 'PROCESSED'
          AND created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS processed_tickets_today,

    COUNT(*) FILTER (
        WHERE status = 'FAILED'
          AND created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS failed_tickets_today,

    COUNT(*) FILTER (
        WHERE status = 'RESOLVED'
          AND created_at >= CURRENT_DATE
          AND created_at < CURRENT_DATE + INTERVAL '1 day'
    ) AS resolved_tickets_today,

    COUNT(*) FILTER (
        WHERE status = 'PENDING'
          AND created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS pending_tickets_yday,

    COUNT(*) FILTER (
        WHERE status = 'PROCESSING'
          AND created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS processing_tickets_yday,

    COUNT(*) FILTER (
        WHERE status = 'PROCESSED'
          AND created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS processed_tickets_yday,

    COUNT(*) FILTER (
        WHERE status = 'FAILED'
          AND created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS failed_tickets_yday,

    COUNT(*) FILTER (
        WHERE status = 'RESOLVED'
          AND created_at >= CURRENT_DATE - INTERVAL '1 day'
          AND created_at < CURRENT_DATE
    ) AS resolved_tickets_yday
FROM tickets;
