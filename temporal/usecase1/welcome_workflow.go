package usecase1

import (
	"fmt"
	"usecase1/alice"
	"usecase1/bob"

	"go.temporal.io/sdk/workflow"
)

type WelcomeOutput struct {
	Messages []string
}

type WelcomeInput struct {
	Name string
}

// Workflow is a Hello World workflow definition.
func WelcomeWorkflow(ctx workflow.Context, input WelcomeInput) (output WelcomeOutput, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Started, Input: %#v", input))
	var aliceFurture workflow.ChildWorkflowFuture
	{
		cwo := workflow.ChildWorkflowOptions{
			TaskQueue: "alice-hello-world",
		}
		ctx := workflow.WithChildOptions(ctx, cwo)
		aliceHelloIn := alice.HelloIn{
			Name: input.Name,
		}
		aliceFurture = workflow.ExecuteChildWorkflow(ctx, alice.HeloWorkflow, aliceHelloIn)

	}

	var bobFurture workflow.ChildWorkflowFuture
	{
		cwo := workflow.ChildWorkflowOptions{
			TaskQueue: "bob-hello-world",
		}
		ctx := workflow.WithChildOptions(ctx, cwo)
		bobHelloIn := bob.HelloIn{
			Name: input.Name,
		}
		bobFurture = workflow.ExecuteChildWorkflow(ctx, bob.HeloWorkflow, bobHelloIn)

	}

	{
		aliceHelloOut := bob.HelloOut{}
		if err := aliceFurture.Get(ctx, &aliceHelloOut); err != nil {
			logger.Error("alice error", "Error", err)
			return output, err
		}
		output.Messages = append(output.Messages, aliceHelloOut.Message)

		bobHelloOut := bob.HelloOut{}
		if err := bobFurture.Get(ctx, &bobHelloOut); err != nil {
			logger.Error("bob error", "Error", err)
			return output, err
		}
		output.Messages = append(output.Messages, bobHelloOut.Message)
	}

	// // call bob workflow
	// bobHelloIn := bob.HelloIn{
	// 	Name: input.Name,
	// }
	// bobFurture := workflow.ExecuteChildWorkflow(ctx, bob.HeloWorkflow, bobHelloIn)

	// get alice result

	// get bob result
	// bobHelloOut := bob.HelloOut{}
	// if err := bobFurture.Get(ctx, &bobHelloOut); err != nil {
	// 	logger.Error("bob error", "Error", err)
	// 	return output, err
	// }
	// output.Messages = append(output.Messages, aliceHelloOut.Message)

	// logger.Info("Wellcome workflow completed.", "Result", output)
	return output, nil
}

// func HelloActivity(ctx context.Context, name string) (string, error) {
// 	logger := activity.GetLogger(ctx)
// 	logger.Info("Activity", "name", name)
// 	return "Hello " + name + "!", nil
// }
