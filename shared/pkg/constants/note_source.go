package constants

type NoteSourceID int

func (p NoteSourceID) Int() int {
	return int(p)
}

func (p NoteSourceID) Name() NoteSourceName {
	if name, exist := noteSourceNameMap[p]; exist {
		return name
	}

	return NoteSourceNameUnknown
}

type NoteSourceName string

func (p NoteSourceName) String() string {
	return string(p)
}

var (
	NoteSourceIDUnknown             NoteSourceID = -1
	NoteSourceIDCrossbell           NoteSourceID = 0
	NoteSourceIDEthereumNFT         NoteSourceID = 1
	NoteSourceIDSolanaNFT           NoteSourceID = 2
	NoteSourceIDFlowNFT             NoteSourceID = 3
	NoteSourceIDMirrorEntry         NoteSourceID = 4
	NoteSourceIDGitcoinContribution NoteSourceID = 5
	NoteSourceIDTwitterTweet        NoteSourceID = 6
	NoteSourceIDMisskeyNote         NoteSourceID = 7
	NoteSourceIDJikePost            NoteSourceID = 8
	NoteSourceIDEthereumERC20       NoteSourceID = 9

	NoteSourceNameUnknown             NoteSourceName = "Unknown"
	NoteSourceNameCrossbell           NoteSourceName = "Crossbell"
	NoteSourceNameEthereumNFT         NoteSourceName = "Ethereum NFT"
	NoteSourceNameEthereumERC20       NoteSourceName = "Ethereum ERC20"
	NoteSourceNameSolanaNFT           NoteSourceName = "Solana NFT"
	NoteSourceNameFlowNFT             NoteSourceName = "Flow NFT"
	NoteSourceNameMirrorEntry         NoteSourceName = "Mirror Entry"
	NoteSourceNameGitcoinContribution NoteSourceName = "Gitcoin Contribution"
	NoteSourceNameTwitterTweet        NoteSourceName = "Twitter Tweet"
	NoteSourceNameMisskeyNote         NoteSourceName = "Misskey Note"
	NoteSourceNameJikePost            NoteSourceName = "Jike Post"

	noteSourceNameMap = map[NoteSourceID]NoteSourceName{
		NoteSourceIDUnknown:             NoteSourceNameUnknown,
		NoteSourceIDCrossbell:           NoteSourceNameCrossbell,
		NoteSourceIDEthereumNFT:         NoteSourceNameEthereumNFT,
		NoteSourceIDSolanaNFT:           NoteSourceNameSolanaNFT,
		NoteSourceIDFlowNFT:             NoteSourceNameFlowNFT,
		NoteSourceIDMirrorEntry:         NoteSourceNameMirrorEntry,
		NoteSourceIDGitcoinContribution: NoteSourceNameGitcoinContribution,
		NoteSourceIDTwitterTweet:        NoteSourceNameTwitterTweet,
		NoteSourceIDMisskeyNote:         NoteSourceNameMisskeyNote,
		NoteSourceIDJikePost:            NoteSourceNameJikePost,
	}
	noteSourceIDMap = map[NoteSourceName]NoteSourceID{}
)

func init() {
	for id, name := range noteSourceNameMap {
		noteSourceIDMap[name] = id
	}
}
