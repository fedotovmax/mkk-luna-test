package teams_test

import (
	"context"
	"testing"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/mysql/teams"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTeam_Create(t *testing.T) {
	ctx := context.Background()

	mockExec := &db.MockStdSQLDriver{}
	mockExtractor := &db.MockExtractor{Exec: mockExec}

	teamDB := teams.New(mockExtractor)

	ownerID := uuid.New().String()

	teamName := "test team"

	id, err := teamDB.Create(ctx, ownerID, teamName)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	require.True(t, mockExec.Called)

	require.Equal(t, 3, len(mockExec.Args))
	require.Equal(t, id, mockExec.Args[0])
	require.Equal(t, teamName, mockExec.Args[1])
	require.Equal(t, ownerID, mockExec.Args[2])
}
