package allure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	allureDir = "./allure-results"
)

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	require.NotNil(t, container)
	require.Zero(t, container.Start)
	require.Zero(t, container.Stop)
	require.NotNil(t, container.UUID)
}

func TestContainer_AddChild(t *testing.T) {
	container := NewContainer()
	result := NewResult("Mock Name", "Full Mock Name")
	container.AddChild(result.UUID)
	require.NotEmpty(t, container.Children)
	require.Len(t, container.Children, 1)
	require.Equal(t, result.UUID, container.Children[0])
}

func TestContainer_Begin(t *testing.T) {
	container := NewContainer()
	begin := GetNow()
	container.Begin()
	assert.Equal(t, container.Start, begin)
}

func TestContainer_Finish(t *testing.T) {
	container := NewContainer()
	finish := GetNow()
	container.Finish()
	assert.Equal(t, container.Stop, finish)
}

func TestContainer_Print(t *testing.T) {
	container := NewContainer()
	containerBefore := NewSimpleStep("Before")
	containerAfter := NewSimpleStep("After")
	container.Befores = append(container.Befores, containerBefore)
	container.Afters = append(container.Afters, containerAfter)

	_ = container.Print()
	defer os.RemoveAll(allureDir)
	require.DirExists(t, allureDir)
	files, _ := ioutil.ReadDir(allureDir)
	require.Len(t, files, 1)
	var jsonFile *os.File
	defer jsonFile.Close()

	f := files[0]
	emptyContainer := &Container{}
	jsonFile, _ = os.Open(fmt.Sprintf("%s/%s", allureDir, f.Name()))
	bytes, readErr := ioutil.ReadAll(jsonFile)
	require.NoError(t, readErr)
	unMarshallErr := json.Unmarshal(bytes, emptyContainer)
	require.NoError(t, unMarshallErr)
	require.Equal(t, container.UUID, emptyContainer.UUID)

	require.NotEmpty(t, emptyContainer.Befores)
	require.NotEmpty(t, emptyContainer.Afters)
	require.Len(t, emptyContainer.Befores, 1)
	require.Len(t, emptyContainer.Afters, 1)

	before := emptyContainer.Befores[0]
	after := emptyContainer.Afters[0]

	require.Equal(t, containerBefore, before)
	require.Equal(t, containerAfter, after)
}
