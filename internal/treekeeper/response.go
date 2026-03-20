package treekeeper

type ResponseKind int

const (
	ResponseBranchCreate ResponseKind = iota
	ResponseBranchDelete
	ResponseCheckout
	ResponseClone
	ResponseList
	ResponseDoctor
	ResponsePrune
	ResponseSync
	ResponseSetup
	ResponsePassThrough
	ResponseSyncAll
)

type Response struct {
	Kind         ResponseKind
	BranchCreate *BranchCreateOutput
	BranchDelete *BranchDeleteOutput
	Checkout     *CheckoutOutput
	Clone        *CloneOutput
	Worktrees    []WorktreeInfo
	Doctor       []DoctorInfo
	Prune        *PruneResult
	Sync         *SyncResult
	SyncAll      *SyncAllResult
	Setup        *SetupResult
	PassThrough  *PassThroughResult
}

type BranchCreateOutput struct {
	Branch       string
	Base         string
	WorktreePath string
}

type BranchDeleteOutput struct {
	Branch        string
	WorktreePath  string
	RemoteDeleted bool
	RemoteName    string
}

type CloneOutput struct {
	RepoURL       string
	DefaultBranch string
	WorktreePath  string
}

type CheckoutOutput struct {
	Branch       string
	WorktreePath string
}

type PruneResult struct {
	DryRun           bool
	PrunedWorktrees  []PrunedWorktree
	PrunedBranches   []PrunedBranch
	SkippedWorktrees []SkippedWorktree
	SkippedBranches  []SkippedBranch
}

type SyncResult struct {
	Branch        string
	WorktreePath  string
	Remote        string
	RemoteBranch  string
	DryRun        bool
	AddedUpstream bool
	UpstreamName  string
	UpstreamURL   string
	SetUpstream   bool
	PushRemote    string
	FetchOutput   []string
	MergeOutput   []string
}

type SetupResult struct {
	Branch         string
	UpstreamName   string
	UpstreamURL    string
	OriginName     string
	AddedUpstream  bool
	SetUpstream    bool
	SetPushRemote  bool
	HooksInstalled bool
	HooksPath      string
	DryRun         bool
}

type PrunedWorktree struct {
	Branch string
	Path   string
}

type SkippedWorktree struct {
	Branch string
	Path   string
	Reason string
}

type PrunedBranch struct {
	Branch string
}

type SkippedBranch struct {
	Branch string
	Reason string
}

type SyncAllResult struct {
	Results []SyncResult
	Skipped []SkippedSync
}

type SkippedSync struct {
	Branch string
	Path   string
	Reason string
}

type PassThroughResult struct {
	Args []string
}
