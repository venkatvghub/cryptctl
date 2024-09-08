package commands

import (
	"fmt"

	"github.com/opensecrecy/cryptctl/commands/utils"
	"github.com/spf13/cobra"
)

func init() {

	initCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to use (required)")
	initCmd.Flags().StringVarP(&Provider, "provider", "p", "", "provider to use (required)")
	_ = initCmd.MarkFlagRequired("provider")
	initCmd.Flags().StringVarP(&KmsKeyAlias, "kmsalias", "k", "", "KMS key alias to use (required for aws-kms-envelope)")
	initCmd.Flags().StringVarP(&DynamoDbTable, "dynamodb-table", "d", "", "DynamoDB table to use (required for aws-kms-envelope)")
	initCmd.Flags().BoolVarP(&Debug, "debug", "i", false, "Run in Debug Mode Using LocalStack for KMS instead of AWS")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [flags]",
	Short: "init",
	Long:  "Init initializes the encrypted-secrets CLI",
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if Provider == "k8s" && Namespace == "" {
			return fmt.Errorf("namespace is required for k8s provider")
		}
		if Provider == "aws-kms-envelope" {
			if KmsKeyAlias == "" {
				return fmt.Errorf("KMS key alias is required for aws-kms-envelope provider")
			}
			if DynamoDbTable == "" {
				return fmt.Errorf("DynamoDB table is required for aws-kms-envelope provider")
			}
		}
		return nil

	},

	RunE: func(_ *cobra.Command, args []string) error {

		switch Provider {
		case "k8s":
			return utils.InitK8s(Namespace)

		case "aws-kms":
			return utils.InitAwsKms(Namespace)

		case "aws-kms-envelope":
			return utils.InitAwsKmsEnvelope(KmsKeyAlias, DynamoDbTable, Debug)
		}

		return nil

	},
}
