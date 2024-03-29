package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/sqldb"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMetricToDBErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	ctxB := context.Background()
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, ctxB)

	//nolint:exhaustruct
	metricUnit := model.MetricUnit{
		Name: "test",
	}
	rows := pgxmock.NewRows([]string{"id", "name", "type", "value"}).
		AddRow("errValue", "test", "gauge", "1.0")

	mock.ExpectQuery(
		"select id, name, type, value from metrics where name=\\$1").
		WithArgs("test").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock
	_, err = AddMetricToDB(ctxB, &pgConn, metricUnit)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetMetricByNameFromDBErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	ctxB := context.Background()
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, ctxB)

	//nolint:exhaustruct
	metricUnit := model.MetricUnit{
		Name: "test",
	}
	rows := pgxmock.NewRows([]string{"name", "type", "value"}).
		AddRow("test", "gauge", 1)

	mock.ExpectQuery(
		"select name, type, value from metrics where name=\\$1").
		WithArgs(metricUnit.Name).
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock
	_, err = GetMetricByNameFromDB(ctxB, &pgConn, metricUnit.Name)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllMetricsFromDBErr(t *testing.T) {
	mock, err := pgxmock.NewConn()
	require.NoError(t, err, fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	ctxB := context.Background()
	defer func(mock pgxmock.PgxConnIface, ctx context.Context) {
		mock.ExpectClose()
		err = mock.Close(ctx)
		require.NoError(t, err)
	}(mock, ctxB)

	rows := pgxmock.NewRows([]string{"name", "type", "value"}).
		AddRow("test", "gauge", 1 /*<-error here*/)

	mock.ExpectQuery(
		"select name, type, value from metrics").
		WillReturnRows(rows)

	var pgConn sqldb.PgxIface = mock
	_, err = GetAllMetricsFromDB(ctxB, &pgConn)
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
