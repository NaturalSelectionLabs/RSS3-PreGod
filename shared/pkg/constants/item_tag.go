package constants

type itemTag string

const (
	itemTagUnknown itemTag = "UNKNOWN"

	itemTagNFT         itemTag = "NFT"
	itemTagPOAP        itemTag = "POAP"
	itemTagMirrorEntry itemTag = "Mirror Entry"
	itemTagDonation    itemTag = "Donation"
	itemTagGitcoin     itemTag = "Gitcoin"
	itemTagTweet       itemTag = "Tweet"
	itemTagMisskeyNote itemTag = "Misskey Note"
	itemTagJikePost    itemTag = "Jike Post"
)

// See https://rss3.io/protocol/RIPs/RIP-4.html
type ItemTags []itemTag

var (
	ItemTagsUnknown ItemTags = []itemTag{itemTagUnknown}

	ItemTagsNFT             ItemTags = []itemTag{itemTagNFT}
	ItemTagsNFTPOAP         ItemTags = []itemTag{itemTagNFT, itemTagPOAP}
	ItemTagsMirrorEntry     ItemTags = []itemTag{itemTagMirrorEntry}
	ItemTagsDonationGitcoin ItemTags = []itemTag{itemTagDonation, itemTagGitcoin}
	ItemTagsTweet           ItemTags = []itemTag{itemTagTweet}
	ItemTagsMisskeyNote     ItemTags = []itemTag{itemTagMisskeyNote}
	ItemTagsJikePost        ItemTags = []itemTag{itemTagJikePost}
)
