package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

const (
	example = `
	# print object yaml
	%[1]s get-yaml <object type> <object name>
`
)

// ErrInsufficientArgs is thrown if arg len <1 or >2
var ErrInsufficientArgs = fmt.Errorf("\nincorrect number or arguments, see --help for usage instructions")

// HiddenField - metadata field to hide
var HiddenField = "managedFields"

// CommandOpts is the struct holding common properties
type CommandOpts struct {
	customNamespace string
	customContext   string
	kubeConfig      string
	objectType      string
	objectName      string
}

// NewCmdGetYaml creates the cobra command to be executed
func NewCmdGetYaml() *cobra.Command {
	res := &CommandOpts{}

	cmd := &cobra.Command{
		Use:          "get-yaml [object type] [object name]",
		Short:        "Print yaml specification without ManagedFields",
		Example:      fmt.Sprintf(example, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := res.Validate(args); err != nil {
				return err
			}
			if err := res.Retrieve(c); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&res.customNamespace, "namespace", "n", res.customNamespace, "override the namespace defined in the current context")
	cmd.Flags().StringVarP(&res.customContext, "context", "c", res.customContext, "override the current context")
	cmd.Flags().StringVarP(&res.kubeConfig, "kubeconfig", "k", res.kubeConfig, "explicitly provide the kubeconfig to use")
	cmd.Flags().StringVar(&res.objectType, "objectType", res.objectType, "type of object")
	cmd.Flags().StringVar(&res.objectName, "objectName", res.objectName, "object name")

	return cmd
}

// Validate ensures proper command usage
func (c *CommandOpts) Validate(args []string) error {
	argLen := len(args)
	if argLen < 1 || argLen > 2 {
		return ErrInsufficientArgs
	}

	c.objectType = args[0]
	if argLen == 2 {
		c.objectName = args[1]
	}

	return nil
}

// Retrieve reads the kubeconfig and get object's yaml
func (c *CommandOpts) Retrieve(cmd *cobra.Command) error {
	nsOverride, _ := cmd.Flags().GetString("namespace")
	ctxOverride, _ := cmd.Flags().GetString("context")
	kubeConfigOverride, _ := cmd.Flags().GetString("kubeconfig")

	var res, cmdErr bytes.Buffer
	commandArgs := []string{"get", c.objectType, c.objectName, "-o", "yaml"}
	if nsOverride != "" {
		commandArgs = append(commandArgs, "-n", nsOverride)
	}

	if ctxOverride != "" {
		commandArgs = append(commandArgs, "--context", ctxOverride)
	}

	if kubeConfigOverride != "" {
		commandArgs = append(commandArgs, "--kubeconfig", kubeConfigOverride)
	}

	out := exec.Command("kubectl", commandArgs...)
	out.Stdout = &res
	out.Stderr = &cmdErr
	err := out.Run()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, cmdErr.String())
		return nil
	}

	var object map[string]interface{}
	if err := yaml.Unmarshal(res.Bytes(), &object); err != nil {
		return err
	}

	return ClearYAML(os.Stdout, os.Stderr, object)
}

// ClearYAML takes the object's json and delete ManagedFields
func ClearYAML(outWriter, errWriter io.Writer, object map[string]interface{}) error {
	metadataRaw, err := yaml.Marshal(object["metadata"])
	if err != nil {
		return err
	}
	var metadataYAML map[string]interface{}
	if err := yaml.Unmarshal(metadataRaw, &metadataYAML); err != nil {
		return err
	}
	delete(metadataYAML, HiddenField)
	object["metadata"] = metadataYAML
	cleanObject, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	fmt.Print(string(cleanObject))
	return nil
}
