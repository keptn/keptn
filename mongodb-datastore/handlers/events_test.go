package handlers

import (
	"go.mongodb.org/mongo-driver/bson"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/magiconair/properties/assert"
)

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
func TestFlattenRecursivelyNestedDocuments(t *testing.T) {
	logger := keptnutils.NewLogger("", "", "mongodb-service")

	grandchild := bson.D{{"apple", "red"}, {"orange", "orange"}}
	child := bson.D{{"foo", "bar"}, {"grandchild", grandchild}}
	parent := bson.D{{"hello", "world"}, {"child", child}}

	// checks:
	flattened, _ := flattenRecursively(parent, logger)
	parentMap, _ := flattened.(map[string]interface{})
	assert.Equal(t, parentMap["hello"], "world", "flatting failed")

	childMap := parentMap["child"].(map[string]interface{})
	assert.Equal(t, childMap["foo"], "bar", "flatting failed")

	grandchildMap := childMap["grandchild"].(map[string]interface{})
	assert.Equal(t, grandchildMap["orange"], "orange", "flatting failed")
}

// TestFlattenRecursivelyNestedDocuments checks whether the flattening works with nested bson.D (documents)
// and bson.A (arrays)
func TestFlattenRecursivelyNestedDocumentsWithArray(t *testing.T) {
	logger := keptnutils.NewLogger("", "", "mongodb-service")

	grandchild := bson.D{{"apple", "red"}, {"orange", "orange"}}
	child := bson.A{grandchild, "foo", "bar"}
	parent := bson.D{{"hello", "world"}, {"child", child}}

	// checks:
	flattened, _ := flattenRecursively(parent, logger)
	parentMap, _ := flattened.(map[string]interface{})
	assert.Equal(t, parentMap["hello"], "world", "flatting failed")

	childMap := parentMap["child"].(bson.A)
	assert.Equal(t, len(childMap), 3, "flatting failed")

	grandchildMap := childMap[0].(map[string]interface{})
	assert.Equal(t, grandchildMap["apple"], "red", "flatting failed")
}
