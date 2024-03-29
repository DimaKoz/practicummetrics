package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/sqldb"
	"github.com/jackc/pgx/v5"
)

// AddMetricTxToDB adds model.MetricUnit to a db
// returns updated model.MetricUnit after that.
func AddMetricTxToDB(ctx context.Context, dbConn *pgx.Tx, metricUnit model.MetricUnit) (model.MetricUnit, error) {
	isInsert := false
	var nameM, typeM, valueM string
	var idM int64
	row := (*dbConn).QueryRow(ctx, "select id, name, type, value from metrics where name=$1", metricUnit.Name)
	err := row.Scan(&idM, &nameM, &typeM, &valueM)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		isInsert = true
	} else if err != nil {
		return model.EmptyMetric, fmt.Errorf("failed to scan a row: %w", err)
	}

	if !isInsert && metricUnit.Type == model.MetricTypeCounter {
		var foundValue int64
		foundValue, err = strconv.ParseInt(valueM, 10, 64)
		if err != nil {
			return model.EmptyMetric, fmt.Errorf("failed to parse int: %w", err)
		}
		metricUnit.ValueInt += foundValue
		metricUnit.Value = strconv.FormatInt(metricUnit.ValueInt, 10)
	}

	if isInsert {
		_, err = (*dbConn).Exec(ctx,
			"insert into metrics(name, type, value) values($1, $2, $3)",
			metricUnit.Name, metricUnit.Type, metricUnit.Value)
	} else {
		_, err = (*dbConn).Exec(ctx,
			"UPDATE metrics SET name = $1, type = $2, value = $3 where id = $4",
			metricUnit.Name, metricUnit.Type, metricUnit.Value, idM)
	}
	if err != nil {
		return model.EmptyMetric, fmt.Errorf("failed to save a metric by: %w", err)
	}

	return metricUnit, nil
}

// AddMetricToDB adds model.MetricUnit to a db
// returns updated model.MetricUnit after that.
func AddMetricToDB(ctx context.Context, pgConn *sqldb.PgxIface, metricUnit model.MetricUnit) (model.MetricUnit, error) {
	isInsert := false
	var nameM, typeM, valueM string
	var idM int64
	row := (*pgConn).QueryRow(ctx, "select id, name, type, value from metrics where name=$1", metricUnit.Name)
	err := row.Scan(&idM, &nameM, &typeM, &valueM)

	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		isInsert = true
	} else if err != nil {
		return model.EmptyMetric, fmt.Errorf("failed to scan a row: %w", err)
	}

	if !isInsert && metricUnit.Type == model.MetricTypeCounter {
		var foundValue int64
		foundValue, err = strconv.ParseInt(valueM, 10, 64)
		if err != nil {
			return model.EmptyMetric, fmt.Errorf("failed to parse int: %w", err)
		}
		metricUnit.ValueInt += foundValue
		metricUnit.Value = strconv.FormatInt(metricUnit.ValueInt, 10)
	}
	if true {
		return metricUnit, nil
	}
	if isInsert {
		_, err = (*pgConn).Exec(ctx,
			"insert into metrics(name, type, value) values($1, $2, $3)",
			metricUnit.Name, metricUnit.Type, metricUnit.Value)
	} else {
		_, err = (*pgConn).Exec(ctx,
			"UPDATE metrics SET name = $1, type = $2, value = $3 where id = $4",
			metricUnit.Name, metricUnit.Type, metricUnit.Value, idM)
	}
	if err != nil {
		return model.EmptyMetric, fmt.Errorf("failed to save a metric by: %w", err)
	}

	return metricUnit, nil
}

// GetMetricByNameFromDB returns a model.MetricUnit and nil error if found
// or model.EmptyMetric and an error.
func GetMetricByNameFromDB(ctx context.Context, pgConn *sqldb.PgxIface, name string) (model.MetricUnit, error) {
	var nameM, typeM, valueM string
	row := (*pgConn).QueryRow(ctx, "select name, type, value from metrics where name=$1", name)
	err := row.Scan(&nameM, &typeM, &valueM)
	if err != nil {
		return model.EmptyMetric, fmt.Errorf("failed to scan a row: %w", err)
	}
	var result model.MetricUnit
	result, err = model.NewMetricUnit(typeM, nameM, valueM)
	if err != nil {
		err = fmt.Errorf("failed to call model.NewMetricUnit: %w", err)
	}

	return result, err
}

// GetAllMetricsFromDB returns a list of model.MetricUnit from the storage.
func GetAllMetricsFromDB(ctx context.Context, pgConn *sqldb.PgxIface) ([]model.MetricUnit, error) {
	result := make([]model.MetricUnit, 0)
	rows, err := (*pgConn).Query(ctx, "select name, type, value from metrics")
	if err != nil {
		return result, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var nameM, typeM, valueM string
		err = rows.Scan(&nameM, &typeM, &valueM)
		if err != nil {
			return result, fmt.Errorf("failed to scan a row: %w", err)
		}
		var metricUnit model.MetricUnit
		if metricUnit, err = model.NewMetricUnit(typeM, nameM, valueM); err == nil {
			result = append(result, metricUnit)
		} else {
			return result, fmt.Errorf("failed to call model.NewMetricUnit: %w", err)
		}
	}

	return result, nil
}
