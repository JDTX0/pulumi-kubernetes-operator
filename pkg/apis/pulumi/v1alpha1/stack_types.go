package v1alpha1

import (
	"context"

	"github.com/pulumi/pulumi/sdk/v2/go/x/auto"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StackSpec defines the desired state of Pulumi Stack being managed by this operator.
type StackSpec struct {
	// Auth info:

	// (optional) AccessTokenSecret is the name of a secret containing the PULUMI_ACCESS_TOKEN for Pulumi access.
	// Deprecated: use SecretEnvsFromPath with PULUMI_ACCESS_TOKEN as key or SecretEnvs with a secret entry key
	// PULUMI_ACCESS_KEY instead.
	AccessTokenSecret string `json:"accessTokenSecret,omitempty"`

	// (optional) Envs is an optional array of config maps containing environment variables to set.
	// Deprecated: use EnvRefs instead.
	Envs []string `json:"envs,omitempty"`

	// (optional) EnvRefs is an optional map containing environment variables as keys and stores descriptors to where
	// the variables' values should be loaded from (one of literal, environment variable, file on the
	// filesystem, or Kubernetes secret) as values.
	EnvRefs map[string]ResourceRef `json:"envRefs,omitempty"`

	// (optional) SecretEnvs is an optional array of secret names containing environment variables to set.
	// Deprecated: use EnvRefs instead.
	SecretEnvs []string `json:"envSecrets,omitempty"`

	// (optional) Backend is an optional backend URL to use for all Pulumi operations.
	// Examples:
	//   - Pulumi Service:              "https://app.pulumi.com" (default)
	//   - Self-managed Pulumi Service: "https://pulumi.acmecorp.com"
	//   - Local:                       "file://./einstein"
	//   - AWS:                         "s3://<my-pulumi-state-bucket>"
	//   - Azure:                       "azblob://<my-pulumi-state-bucket>"
	//   - GCP:                         "gs://<my-pulumi-state-bucket>"
	// See: https://www.pulumi.com/docs/intro/concepts/state/
	Backend string `json:"backend,omitempty"`

	// Stack identity:

	// Stack is the fully qualified name of the stack to deploy (<org>/<stack>).
	Stack string `json:"stack"`
	// (optional) Config is the configuration for this stack, which can be optionally specified inline. If this
	// is omitted, configuration is assumed to be checked in and taken from the source repository.
	Config map[string]string `json:"config,omitempty"`
	// (optional) Secrets is the secret configuration for this stack, which can be optionally specified inline. If this
	// is omitted, secrets configuration is assumed to be checked in and taken from the source repository.
	// Deprecated: use SecretRefs instead.
	Secrets map[string]string `json:"secrets,omitempty"`

	// (optional) SecretRefs is the secret configuration for this stack which can be specified through ResourceRef.
	// If this is omitted, secrets configuration is assumed to be checked in and taken from the source repository.
	SecretRefs map[string]ResourceRef `json:"secretsRef,omitempty"`
	// (optional) SecretsProvider is used to initialize a Stack with alternative encryption.
	// Examples:
	//   - AWS:   "awskms://arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34bc-56ef-1234567890ab?region=us-east-1"
	//   - Azure: "azurekeyvault://acmecorpvault.vault.azure.net/keys/mykeyname"
	//   - GCP:   "gcpkms://projects/MYPROJECT/locations/MYLOCATION/keyRings/MYKEYRING/cryptoKeys/MYKEY"
	// See: https://www.pulumi.com/docs/intro/concepts/config/#initializing-a-stack-with-alternative-encryption
	SecretsProvider string `json:"secretsProvider,omitempty"`

	// Source control:

	// ProjectRepo is the git source control repository from which we fetch the project code and configuration.
	ProjectRepo string `json:"projectRepo"`
	// (optional) GitAuthSecret is the the name of a secret containing an
	// authentication option for the git repository.
	// There are 3 different authentication options:
	//   * Personal access token
	//   * SSH private key (and it's optional password)
	//   * Basic auth username and password
	// Only 1 authentication path is valid.
	GitAuthSecret string `json:"gitAuthSecret,omitempty"`
	// (optional) RepoDir is the directory to work from in the project's source repository
	// where Pulumi.yaml is located. It is used in case Pulumi.yaml is not
	// in the project source root.
	RepoDir string `json:"repoDir,omitempty"`
	// (optional) Commit is the hash of the commit to deploy. If used, HEAD will be in detached mode. This
	// is mutually exclusive with the Branch setting. If both are empty, the `master` branch is deployed.
	Commit string `json:"commit,omitempty"`
	// (optional) Branch is the branch name to deploy, either the simple or fully qualified ref name. This
	// is mutually exclusive with the Commit setting. If both are empty, the `master` branch is deployed.
	Branch string `json:"branch,omitempty"`

	// Lifecycle:

	// (optional) Refresh can be set to true to refresh the stack before it is updated.
	Refresh bool `json:"refresh,omitempty"`
	// (optional) ExpectNoRefreshChanges can be set to true if a stack is not expected to have
	// changes during a refresh before the update is run.
	// This could occur, for example, is a resource's state is changing outside of Pulumi
	// (e.g., metadata, timestamps).
	ExpectNoRefreshChanges bool `json:"expectNoRefreshChanges,omitempty"`
	// (optional) DestroyOnFinalize can be set to true to destroy the stack completely upon deletion of the CRD.
	DestroyOnFinalize bool `json:"destroyOnFinalize,omitempty"`
	// (optional) RetryOnUpdateConflict issues a stack update retry reconciliation loop
	// in the event that the update hits a HTTP 409 conflict due to
	// another update in progress.
	// This is only recommended if you are sure that the stack updates are
	// idempotent, and if you are willing to accept retry loops until
	// all spawned retries succeed. This will also create a more populated,
	// and randomized activity timeline for the stack in the Pulumi Service.
	RetryOnUpdateConflict bool `json:"retryOnUpdateConflict,omitempty"`
}

// ResourceRef identifies a resource from which information can be loaded.
// Environment variables, files on the filesystem, Kubernetes secrets and literal
// strings are currently supported.
type ResourceRef struct {
	// SelectorType is required and signifies the type of selector. Must be one of:
	// Env, FS, Secret, Literal
	SelectorType     ResourceSelectorType `json:"type"`
	ResourceSelector `json:",inline"`
}

// NewEnvResourceRef creates a new environment variable resource ref.
func NewEnvResourceRef(envVarName string) ResourceRef {
	return ResourceRef{
		SelectorType: ResourceSelectorEnv,
		ResourceSelector: ResourceSelector{
			Env: &EnvSelector{
				Name: envVarName,
			},
		},
	}
}

// NewFileSystemResourceRef creates a new file system resource ref.
func NewFileSystemResourceRef(path string) ResourceRef {
	return ResourceRef{
		SelectorType: ResourceSelectorFS,
		ResourceSelector: ResourceSelector{
			FileSystem: &FSSelector{
				Path: path,
			},
		},
	}
}

// NewSecretResourceRef creates a new secret resource ref.
func NewSecretResourceRef(namespace, name, key string) ResourceRef {
	return ResourceRef{
		SelectorType: ResourceSelectorSecret,
		ResourceSelector: ResourceSelector{
			SecretRef: &SecretSelector{
				Namespace: namespace,
				Name:      name,
				Key:       key,
			},
		},
	}
}

// NewLiteralResourceRef creates a new literal resource ref.
func NewLiteralResourceRef(value string) ResourceRef {
	return ResourceRef{
		SelectorType: ResourceSelectorLiteral,
		ResourceSelector: ResourceSelector{
			LiteralRef: &LiteralRef{
				Value: value,
			},
		},
	}
}

// ResourceSelectorType identifies the type of the resource reference in
type ResourceSelectorType string

const (
	// ResourceSelectorEnv indicates the resource is an environment variable
	ResourceSelectorEnv = ResourceSelectorType("Env")
	// ResourceSelectorFS indicates the resource is on the filesystem
	ResourceSelectorFS = ResourceSelectorType("FS")
	// ResourceSelectorSecret indicates the resource is a Kubernetes secret
	ResourceSelectorSecret = ResourceSelectorType("Secret")
	// ResourceSelectorLiteral indicates the resource is a literal
	ResourceSelectorLiteral = ResourceSelectorType("Literal")
)

// ResourceSelector is a union over resource selectors supporting one of
// filesystem, environment variable, Kubernetes Secret and literal values.
type ResourceSelector struct {
	// FileSystem selects a file on the operator's file system
	FileSystem *FSSelector `json:"filesystem,omitempty"`
	// Env selects an environment variable set on the operator process
	Env *EnvSelector `json:"env,omitempty"`
	// SecretRef refers to a Kubernetes secret
	SecretRef *SecretSelector `json:"secret,omitempty"`
	// LiteralRef refers to a literal value
	LiteralRef *LiteralRef `json:"literal,omitempty"`
}

// FSSelector identifies the path to load information from.
type FSSelector struct {
	// Path on the filesystem to use to load information from.
	Path string `json:"path"`
}

// EnvSelector identifies the environment variable to load information from.
type EnvSelector struct {
	// Name of the environment variable
	Name string `json:"name"`
}

// SecretSelector identifies the information to load from a Kubernetes secret.
type SecretSelector struct {
	// Namespace where the secret is stored. Defaults to 'default' if omitted.
	Namespace string `json:"namespace,omitempty"`
	// Name of the secret
	Name string `json:"name"`
	// Key within the secret to use.
	Key string `json:"key"`
}

// LiteralRef identifies a literal value to load.
type LiteralRef struct {
	// Value to load
	Value string `json:"value"`
}

// StackStatus defines the observed state of Stack
type StackStatus struct {
	// Outputs contains the exported stack output variables resulting from a deployment.
	Outputs StackOutputs `json:"outputs,omitempty"`
	// LastUpdate contains details of the status of the last update.
	LastUpdate *StackUpdateState `json:"lastUpdate,omitempty"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

type StackOutputs map[string]apiextensionsv1.JSON

// StackUpdateState is the status of a stack update
type StackUpdateState struct {
	// State is the state of the stack update - one of `succeeded` or `failed`
	State StackUpdateStateMessage `json:"state,omitempty"`
	// Last commit attempted
	LastAttemptedCommit string `json:"lastAttemptedCommit,omitempty"`
	// Last commit successfully applied
	LastSuccessfulCommit string `json:"lastSuccessfulCommit,omitempty"`
	// Permalink is the Pulumi Console URL of the stack operation.
	Permalink Permalink `json:"permalink,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Stack is the Schema for the stacks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=stacks,scope=Namespaced
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSpec   `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StackList contains a list of Stack
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

// StackUpdateStatus is the status code for the result of a Stack Update run.
type StackUpdateStatus int

const (
	// StackUpdateSucceeded indicates that the stack update completed successfully.
	StackUpdateSucceeded StackUpdateStatus = 0
	// StackUpdateFailed indicates that the stack update failed to complete.
	StackUpdateFailed StackUpdateStatus = 1
	// StackUpdateConflict indicates that the stack update failed to complete due
	// to a conflicting stack update run that is in progress.
	StackUpdateConflict StackUpdateStatus = 2
	// StackUpdatePendingOperations indicates that the stack update failed to complete due
	// to pending operations halting the stack update run.
	StackUpdatePendingOperations StackUpdateStatus = 3
	// StackNotFound indicates that the stack update failed to complete due
	// to stack not being found (HTTP 404) in the Pulumi Service.
	StackNotFound StackUpdateStatus = 4
)

type StackUpdateStateMessage string

const (
	// SucceededStackStateMessage is a const to indicate success in stack status state.
	SucceededStackStateMessage StackUpdateStateMessage = "succeeded"
	// FailedStackStateMessage is a const to indicate stack failure in stack status state.
	FailedStackStateMessage StackUpdateStateMessage = "failed"
)

// Permalink is the Pulumi Service URL of the stack operation.
type Permalink string

// StackController contains methods to operate a Pulumi Project and Stack in an update.
//
// Ignoring operator codegen of interface as it is an API contract for implementation,
// not a type that is used in kubernetes.
// +kubebuilder:object:generate=false
type StackController interface {
	// Project setup:

	// InstallProjectDependencies installs the package manager dependencies for the project's language.
	InstallProjectDependencies(ctx context.Context, workspace auto.Workspace) error
	// SetEnvs populates the environment of the stack run with values
	// from an array of Kubernetes ConfigMaps in a Namespace.
	SetEnvs(configMapNames []string, namespace string) error
	// SetSecretEnvs populates the environment of the stack run with values
	// from an array of Kubernetes Secrets in a Namespace.
	SetSecretEnvs(secretNames []string, namespace string) error

	// Lifecycle:

	// UpdateConfig updates the stack configuration values and secret values by
	// combining any configuration values checked into the source repository with
	// the Config values provided in the Stack, overriding values that match and exist.
	UpdateConfig() error
	// RefreshStack refreshes the stack before the update step is run, and
	// errors the run if changes were not expected but found after the refresh.
	RefreshStack(expectNoChanges bool) (Permalink, error)
	// UpdateStack deploys the stack's resources, computes the new desired
	// state, and returns the update's status.
	UpdateStack() (StackUpdateStatus, Permalink, *auto.UpResult, error)
	// GetStackOutputs returns all of the the stack's output properties.
	GetStackOutputs(outputs auto.OutputMap) (StackOutputs, error)
	// DestroyStack destroys the stack's resources and state, and the stack itself.
	DestroyStack() error
}

func init() {
	SchemeBuilder.Register(&Stack{}, &StackList{})
}
