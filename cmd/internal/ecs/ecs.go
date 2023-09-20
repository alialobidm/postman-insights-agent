package ecs

import (
	"fmt"

	"github.com/akitasoftware/akita-cli/cmd/internal/cmderr"
	"github.com/akitasoftware/akita-cli/rest"
	"github.com/akitasoftware/akita-cli/telemetry"
	"github.com/akitasoftware/akita-cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Mandatory flag: Postman collection id
	collectionId string

	// Any of these will be interactively prompted if not given on the command line.
	// On the other hand, to run non-interactively then all of them *must* be given.
	awsProfileFlag        string
	awsRegionFlag         string
	ecsClusterFlag        string
	ecsServiceFlag        string
	ecsTaskDefinitionFlag string

	// Location of credentials file.
	awsCredentialsFlag string

	// Print out the steps that would be taken, but do not do them
	dryRunFlag bool
)

var Cmd = &cobra.Command{
	Use:   "ecs",
	Short: "Add the Postman Live Collections Agent to AWS ECS.",
	Long:  "The CLI will collect information from you and add the Postman Live Collections Agent container to an ECS Task.",
	// N.B.: this is useless because the root command makes its own determination,
	// need to return AkitaErr to not show the usage.
	SilenceUsage: true,
	RunE:         addAgentToECS,
}

// 'postman-lc-agent ecs' should default to 'postman-lc-agent ecs add'
var AddToECSCmd = &cobra.Command{
	Use:          "add",
	Short:        Cmd.Short,
	Long:         Cmd.Long,
	SilenceUsage: true,
	RunE:         addAgentToECS,
}

var RemoveFromECSCmd = &cobra.Command{
	Use:          "remove",
	Short:        "Remove the Postman Live Collections Agent from AWS ECS.",
	Long:         "Remove a previously installed Postman container from an ECS Task.",
	SilenceUsage: true,
	RunE:         removeAgentFromECS,

	// Temporarily hide from users until complete
	Hidden: true,
}

func init() {
	// TODO: add the ability to specify the credentials directly instead of via an AWS profile?
	Cmd.PersistentFlags().StringVar(&collectionId, "collection", "", "Your Postman collection ID")
	Cmd.PersistentFlags().StringVar(&awsProfileFlag, "profile", "", "Which of your AWS profiles to use to access ECS.")
	Cmd.PersistentFlags().StringVar(&awsRegionFlag, "region", "", "The AWS region in which your ECS cluster resides.")
	Cmd.PersistentFlags().StringVar(&ecsClusterFlag, "cluster", "", "The name or ARN of your ECS cluster.")
	Cmd.PersistentFlags().StringVar(&ecsServiceFlag, "service", "", "The name or ARN of your ECS service.")
	Cmd.PersistentFlags().StringVar(
		&ecsTaskDefinitionFlag,
		"task",
		"",
		"The name of your ECS task definition to modify.",
	)
	Cmd.PersistentFlags().BoolVar(
		&dryRunFlag,
		"dry-run",
		false,
		"Perform a dry run: show what will be done, but do not modify ECS.",
	)

	// Support for credentials in a nonstandard location
	Cmd.PersistentFlags().StringVar(&awsCredentialsFlag, "aws-credentials", "", "Location of AWS credentials file.")
	Cmd.PersistentFlags().MarkHidden("aws-credentials")

	Cmd.AddCommand(AddToECSCmd)
	Cmd.AddCommand(RemoveFromECSCmd)
}

func addAgentToECS(cmd *cobra.Command, args []string) error {
	// Check for API key
	_, err := cmderr.RequirePostmanAPICredentials("The Postman Live Collections Agent must have an API key in order to capture traces.")
	if err != nil {
		return err
	}

	// Check collecton Id's existence
	if collectionId == "" {
		return errors.New("Must specify the ID of your collection with the --collection flag.")
	}
	frontClient := rest.NewFrontClient(rest.Domain, telemetry.GetClientID())
	_, err = util.GetOrCreateServiceIDByPostmanCollectionID(frontClient, collectionId)
	if err != nil {
		return err
	}

	return RunAddWorkflow()
}

func removeAgentFromECS(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("this command is not yet implemented")
}
