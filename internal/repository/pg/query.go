package pg

const (
	updateMetricsQuery = `
		INSERT INTO metrics (name, type, delta, value, hash)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (name) 
        DO UPDATE SET
            type  = EXCLUDED.type,
            delta = metrics.delta + EXCLUDED.delta,
            value = EXCLUDED.value,
            hash = EXCLUDED.hash,
            updated_at = NOW();
		`
	getMetricByID = `
		SELECT name, type, delta, value, hash 
		FROM metrics 
		WHERE name = $1;`
	getAllMetrics = `
		SELECT name, type, delta, value, hash 
		FROM metrics`
)
