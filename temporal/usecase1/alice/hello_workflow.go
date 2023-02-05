package alice

import (
	"fmt"

	"go.temporal.io/sdk/workflow"
)

type HelloOut struct {
	Message string
}

type HelloIn struct {
	Name string
}

// Workflow is a Hello World workflow definition.
func HeloWorkflow(ctx workflow.Context, input HelloIn) (HelloOut, error) {
	result := HelloOut{
		Message: fmt.Sprintf("Alice says hello %s", input.Name),
	}
	return result, nil
}

// func HelloActivity(ctx context.Context, name string) (string, error) {
// 	logger := activity.GetLogger(ctx)
// 	logger.Info("Activity", "name", name)
// 	return "Hello " + name + "!", nil
// }
