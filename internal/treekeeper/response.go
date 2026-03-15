package treekeeper

type ResponseKind int

const (
	ResponseBranchCreate ResponseKind = iota
	ResponseBranchDelete
	ResponseCheckout
	ResponseClone
	ResponseList
	ResponseDoctor
)

type Response struct {
	Kind         ResponseKind
	BranchCreate *BranchCreateOutput
	BranchDelete *BranchDeleteOutput
	Checkout     *CheckoutOutput
	Clone        *CloneOutput
	Worktrees    []WorktreeInfo
	Doctor       []DoctorInfo
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
