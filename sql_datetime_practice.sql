-- PostgreSQL Date/Time Practice
--
-- Workflow:
-- 1. Solve only the current challenge.
-- 2. Tell Codex "done" or "complete".
-- 3. Codex will review this file and append the next challenge.
--
-- Training goal:
-- Build a strong mental model for PostgreSQL date, time, timestamp, and
-- timestamptz from beginner to advanced analyst use.
--
-- Core mental model:
-- A table is an array of tuples.
--
--   events = [
--     (1, 'ticket_created', DATE '2026-05-30', TIME '09:15', TIMESTAMP '2026-05-30 09:15:00'),
--     ...
--   ]
--
-- Date/time values are still just column values inside each tuple. The hard
-- part is knowing what each type means:
--
-- - date: a calendar day, no time of day, no timezone.
-- - time: a time of day, no date, no timezone.
-- - timestamp: date + time, but no timezone.
-- - timestamptz: an absolute instant in time, displayed using the session
--   timezone.
-- - interval: a duration, such as 3 hours, 2 days, or 1 month.
--
-- Difficulty plan:
-- 1. Literal values and type differences
-- 2. Casting and extracting parts
-- 3. Comparing dates and timestamps
-- 4. Time ranges and inclusive/exclusive bounds
-- 5. Intervals and durations
-- 6. date_trunc and time bucketing
-- 7. Timezones and timestamptz
-- 8. Business reporting by day/week/month
-- 9. Cohorts and retention over time
-- 10. Performance, indexes, and query plans for time filters

-- Challenge 1 of 12
--
-- Concept:
-- Basic date/time types as columns in a row set.
--
-- Mental model:
-- You will create a temporary row source from arrays. Each produced row is a
-- tuple with four columns:
--
--   (ticket_id, created_date, created_time, created_at_local)
--
-- The values are different types:
--
-- - created_date is a date
-- - created_time is a time
-- - created_at_local is a timestamp without timezone
--
-- Use these arrays:
--
--   ticket_ids = [101, 102, 103]
--   created_dates = ['2026-05-28', '2026-05-29', '2026-05-30']
--   created_times = ['09:15', '14:45', '18:05']
--   created_at_local = [
--     '2026-05-28 09:15:00',
--     '2026-05-29 14:45:00',
--     '2026-05-30 18:05:00'
--   ]
--
-- Expected output shape:
--
--   ticket_id  created_date  created_time  created_at_local
--   ---------  ------------  ------------  -------------------
--   101        2026-05-28    09:15:00      2026-05-28 09:15:00
--   102        2026-05-29    14:45:00      2026-05-29 14:45:00
--   103        2026-05-30    18:05:00      2026-05-30 18:05:00
--
-- Requirements:
-- - Use multi-array unnest.
-- - Cast each array to the correct PostgreSQL array type:
--   - integer[]
--   - date[]
--   - time[]
--   - timestamp[]
-- - Alias the produced columns as:
--   ticket_id, created_date, created_time, created_at_local
-- - Return all four columns.
-- - Order by ticket_id.
-- - End the query with a semicolon.
--
-- Write your answer below:


-- Period Comparison Mini-Track
--
-- These drills are about questions like:
--
--   "Show me tickets from last week and today so I can compare them."
--
-- Core pattern:
--
--   1. Define the time windows.
--   2. Filter rows that belong to either window.
--   3. Label each row with the window it belongs to.
--   4. Aggregate if you want counts, totals, averages, or rates.
--
-- For these exercises, pretend "today" is fixed as:
--
--   DATE '2026-05-30'
--
-- Do not use CURRENT_DATE yet. Fixed dates make the expected output stable
-- while you learn the pattern.

-- Challenge 4 of 12
--
-- Concept:
-- Select rows from two separate time windows and label each row.
--
-- Mental model:
-- You have one array of ticket tuples:
--
--   [
--     (ticket_id, created_at),
--     ...
--   ]
--
-- You only want tuples that are inside either:
--
--   last_week = [2026-05-18 00:00, 2026-05-25 00:00)
--   today     = [2026-05-30 00:00, 2026-05-31 00:00)
--
-- The square bracket [ means included.
-- The round bracket ) means excluded.
--
-- Use this CTE:
--
--   ticket_events(ticket_id, created_at)
--     (201, '2026-05-17 23:59:59')
--     (202, '2026-05-18 00:00:00')
--     (203, '2026-05-20 10:30:00')
--     (204, '2026-05-24 23:59:59')
--     (205, '2026-05-25 00:00:00')
--     (206, '2026-05-30 00:00:00')
--     (207, '2026-05-30 15:45:00')
--     (208, '2026-05-31 00:00:00')
--
-- Task:
-- Return only rows from last_week or today, and add a period label.
--
-- Expected output shape:
--
--   period     ticket_id  created_at
--   ---------  ---------  -------------------
--   last_week  202        2026-05-18 00:00:00
--   last_week  203        2026-05-20 10:30:00
--   last_week  204        2026-05-24 23:59:59
--   today      206        2026-05-30 00:00:00
--   today      207        2026-05-30 15:45:00
--
-- Requirements:
-- - Use a CTE named ticket_events.
-- - Use VALUES inside the CTE.
-- - Cast created_at values to timestamp.
-- - Use CASE to create a period column:
--   - 'last_week'
--   - 'today'
-- - Use half-open ranges:
--   - created_at >= TIMESTAMP '2026-05-18 00:00:00'
--   - created_at <  TIMESTAMP '2026-05-25 00:00:00'
--   - created_at >= TIMESTAMP '2026-05-30 00:00:00'
--   - created_at <  TIMESTAMP '2026-05-31 00:00:00'
-- - Do not use BETWEEN.
-- - Do not use created_at::date in the WHERE clause.
-- - Return period, ticket_id, created_at.
-- - Order by period, ticket_id.
-- - End with a semicolon.
--
-- Write your answer below:


-- Challenge 5 of 12
--
-- Concept:
-- Count rows per comparison period.
--
-- Do this after Challenge 4 has been reviewed.
--
-- Task:
-- Using the same ticket_events rows from Challenge 4, return one row per
-- comparison period:
--
-- Expected output shape:
--
--   period     ticket_count
--   ---------  ------------
--   last_week  3
--   today      2
--
-- Requirements:
-- - Start from the same ticket_events CTE.
-- - Label rows with CASE.
-- - Filter to only last_week and today rows.
-- - GROUP BY period.
-- - COUNT(*) as ticket_count.
-- - End with a semicolon.
--
-- Write your answer below after Challenge 4 is reviewed.


-- Challenge 6 of 12
--
-- Concept:
-- Put comparison periods side by side using FILTER.
--
-- Do this after Challenge 5 has been reviewed.
--
-- Task:
-- Using the same ticket_events rows from Challenge 4, return one row with two
-- count columns:
--
-- Expected output shape:
--
--   last_week_count  today_count
--   ---------------  -----------
--   3                2
--
-- Requirements:
-- - Use COUNT(*) FILTER (WHERE ...).
-- - Use half-open timestamp ranges inside each FILTER.
-- - Return exactly one row.
-- - Do not use GROUP BY.
-- - End with a semicolon.
--
-- Write your answer below after Challenge 5 is reviewed.
SELECT ticket_id, created_date, created_time, created_at_local
FROM UNNEST(
    ARRAY[101,102,103],
    ARRAY[2026-05-28,2026-05-29,2026-05-30],
    ARRAY['09:15:00','14:45:00','18:05:00'],
    ARRAY['2026-05-28 09:15:00','2026-05-29 14:45:00','2026-05-30 18:05:00']
) AS tickets(ticket_id, created_date, created_time, created_at_local)

-- Review:
-- Your mental model is close: one multi-array UNNEST creates one temporary
-- row source with four columns.
--
-- What needs fixing:
--
-- 1. The date values need quotes and a date[] cast.
--
--      ARRAY[2026-05-28, 2026-05-29]
--
--    PostgreSQL reads 2026-05-28 as arithmetic:
--
--      2026 - 5 - 28 = 1993
--
--    Use string-looking date literals with a date[] cast:
--
--      ARRAY['2026-05-28', '2026-05-29']::date[]
--
-- 2. The time and timestamp arrays should also be cast to the target types.
--    Without casts, PostgreSQL may treat them as text first.
--
-- 3. The query needs ORDER BY ticket_id.
--
-- 4. The query needs a semicolon.
--
-- Result grain:
-- - Correct intended grain: one row per ticket_id.
--
-- Double-counting:
-- - No double-counting. Multi-array UNNEST zips the arrays by position.
--
-- NULL handling:
-- - No NULLs in this drill.
--
-- Performance:
-- - Fine for a tiny generated row set. This is a learning query, not a table
--   scan.
--
-- Corrected version:
SELECT
  ticket_id,
  created_date,
  created_time,
  created_at_local
FROM UNNEST(
  ARRAY[101, 102, 103]::integer[],
  ARRAY['2026-05-28', '2026-05-29', '2026-05-30']::date[],
  ARRAY['09:15', '14:45', '18:05']::time[],
  ARRAY[
    '2026-05-28 09:15:00',
    '2026-05-29 14:45:00',
    '2026-05-30 18:05:00'
  ]::timestamp[]
) AS tickets(ticket_id, created_date, created_time, created_at_local)
ORDER BY ticket_id;

-- Mental model:
-- Think of each typed array as a column vector:
--
--   ticket_id column:        [101, 102, 103]
--   created_date column:     [2026-05-28, 2026-05-29, 2026-05-30]
--   created_time column:     [09:15, 14:45, 18:05]
--   created_at_local column: [2026-05-28 09:15, ...]
--
-- Multi-array UNNEST zips those column vectors into tuple rows:
--
--   (101, 2026-05-28, 09:15, 2026-05-28 09:15)
--   (102, 2026-05-29, 14:45, 2026-05-29 14:45)
--   (103, 2026-05-30, 18:05, 2026-05-30 18:05)

-- Challenge 2 of 12
--
-- Concept:
-- Extracting parts from date/time/timestamp values.
--
-- Mental model:
-- A timestamp is one column value, but PostgreSQL can derive new values from
-- it. For example, from:
--
--   TIMESTAMP '2026-05-30 18:05:00'
--
-- PostgreSQL can derive:
--
--   date part:       2026-05-30
--   hour part:       18
--   day of week:     6
--
-- Use this CTE:
--
--   ticket_events(ticket_id, created_at_local)
--     (101, '2026-05-28 09:15:00')
--     (102, '2026-05-29 14:45:00')
--     (103, '2026-05-30 18:05:00')
--
-- Expected output shape:
--
--   ticket_id  created_at_local     created_date  created_hour  day_of_week
--   ---------  -------------------  ------------  ------------  -----------
--   101        2026-05-28 09:15:00  2026-05-28    9             4
--   102        2026-05-29 14:45:00  2026-05-29    14            5
--   103        2026-05-30 18:05:00  2026-05-30    18            6
--
-- Notes:
-- - EXTRACT(DOW FROM timestamp) returns 0 for Sunday, 1 for Monday, ...,
--   6 for Saturday.
-- - EXTRACT returns a numeric value, so cast extracted values to integer for
--   this challenge.
--
-- Requirements:
-- - Use a CTE named ticket_events.
-- - Use VALUES inside the CTE.
-- - Cast created_at_local values to timestamp.
-- - Return:
--   - ticket_id
--   - created_at_local
--   - created_at_local::date AS created_date
--   - EXTRACT(HOUR FROM created_at_local)::integer AS created_hour
--   - EXTRACT(DOW FROM created_at_local)::integer AS day_of_week
-- - Order by ticket_id.
-- - End with a semicolon.
--
-- Write your answer below:
WITH tickets(ticket_id, created_at_local)
VALUES 
    (101, '2026-05-28 09:15:00'),
    (102, '2026-05-29 14:45:00'),
    (103, '2026-05-30 18:05:00')
SELECT 
    ticket_id,
    created_at_local::date AS created_date,
    EXTRACT(HOUR FROM created_at_local)::integer as created_hour,
    EXTRACT(DOW FROM created_at_local)::integer as day_of_week,
FROM tickets
ORDER BY ticket_id;

-- Review:
-- The mental model is right: start with one timestamp column, then derive
-- extra columns from it.
--
-- What needs fixing:
--
-- 1. CTE syntax needs AS (...):
--
--      WITH ticket_events(ticket_id, created_at_local) AS (
--        VALUES ...
--      )
--
-- 2. The challenge asked for the CTE name ticket_events, not tickets.
--
-- 3. The timestamp strings need a timestamp cast. A good pattern with VALUES
--    is to cast each literal:
--
--      '2026-05-28 09:15:00'::timestamp
--
-- 4. The output is missing created_at_local itself.
--
-- 5. There is an extra comma after day_of_week. In SQL, the final selected
--    column does not have a trailing comma.
--
-- Result grain:
-- - Correct intended grain: one row per ticket event.
--
-- Double-counting:
-- - No double-counting. There is one source row per ticket_id and no joins.
--
-- NULL handling:
-- - No NULLs in this drill. If created_at_local were NULL, the derived date,
--   hour, and day_of_week would also be NULL.
--
-- Performance:
-- - Fine. On a real table, extracting parts in SELECT is usually okay. The
--   performance risk appears when you put functions around indexed timestamp
--   columns in WHERE clauses. We will practice that later.
--
-- Corrected version:
WITH ticket_events(ticket_id, created_at_local) AS (
  VALUES
    (101, '2026-05-28 09:15:00'::timestamp),
    (102, '2026-05-29 14:45:00'::timestamp),
    (103, '2026-05-30 18:05:00'::timestamp)
)
SELECT
  ticket_id,
  created_at_local,
  created_at_local::date AS created_date,
  EXTRACT(HOUR FROM created_at_local)::integer AS created_hour,
  EXTRACT(DOW FROM created_at_local)::integer AS day_of_week
FROM ticket_events
ORDER BY ticket_id;

-- Mental model:
-- The CTE creates an array of tuples:
--
--   [
--     (101, TIMESTAMP '2026-05-28 09:15:00'),
--     (102, TIMESTAMP '2026-05-29 14:45:00'),
--     (103, TIMESTAMP '2026-05-30 18:05:00')
--   ]
--
-- The SELECT keeps the original tuple values and adds derived values. It does
-- not create more rows. It widens each row with extra columns.

-- Challenge 3 of 12
--
-- Concept:
-- Filtering timestamps with a half-open time range.
--
-- Mental model:
-- For datetime filtering, prefer:
--
--   created_at >= start_timestamp
--   AND created_at < end_timestamp
--
-- This is called a half-open range. It includes the start and excludes the
-- end. It is safer than BETWEEN for full-day timestamp filters.
--
-- Use this CTE:
--
--   ticket_events(ticket_id, created_at)
--     (101, '2026-05-29 00:00:00')
--     (102, '2026-05-29 13:30:00')
--     (103, '2026-05-29 23:59:59')
--     (104, '2026-05-30 00:00:00')
--     (105, '2026-05-30 09:10:00')
--
-- Task:
-- Return only tickets created on calendar day 2026-05-29.
--
-- Expected output shape:
--
--   ticket_id  created_at
--   ---------  -------------------
--   101        2026-05-29 00:00:00
--   102        2026-05-29 13:30:00
--   103        2026-05-29 23:59:59
--
-- Requirements:
-- - Use a CTE named ticket_events.
-- - Use VALUES inside the CTE.
-- - Cast created_at values to timestamp.
-- - Return ticket_id and created_at.
-- - Filter using:
--   - created_at >= TIMESTAMP '2026-05-29 00:00:00'
--   - created_at < TIMESTAMP '2026-05-30 00:00:00'
-- - Do not use BETWEEN.
-- - Do not use created_at::date in the WHERE clause.
-- - Order by ticket_id.
-- - End with a semicolon.
--
-- Write your answer below:
    
