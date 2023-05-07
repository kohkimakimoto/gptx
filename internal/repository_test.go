package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestRepositoryAwareAction(t *testing.T) {
	t.Run("inject repository", func(t *testing.T) {
		app := testNewApp(t)
		app.Commands = append(app.Commands, &cli.Command{
			Name: "repository-test-command",
			Action: repositoryAwareAction(func(c *cli.Context, r *Repository) error {
				assert.NotNil(t, r)
				return nil
			}),
		})
		err := app.Run([]string{"repository-test-command"})
		assert.NoError(t, err)
	})
}

func TestHomeDirNotFoundError_Error(t *testing.T) {
	err := HomeDirNotFoundError{Dir: "/tmp/gptx"}
	assert.Equal(t, "The gptx configuration has not been initialized yet in the directory '/tmp/gptx'. Please run 'gptx init'.", err.Error())
}

func TestRepository_NewChatService(t *testing.T) {
	// just run
	app := testNewApp(t)
	r, err := getRepository(app)
	assert.NoError(t, err)
	cs, err := r.NewChatService(app.Writer)
	assert.NoError(t, err)
	assert.NotNil(t, cs)
}
