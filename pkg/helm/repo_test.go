package helm

import (
	"github.com/gridworkz/kato/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoAdd(t *testing.T) {
	repo := NewRepo(
		"/tmp/helm/repoName/repositories.yaml",
		"/tmp/helm/cache")
	err := repo.Add(util.NewUUID(), "https://openchart.gridworkz.com/gridworkz/kato", "", "")
	assert.Nil(t, err)
}
