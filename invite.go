package disgord

import (
	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

// PartialInvite ...
// {
//    "code": "abc"
// }
type PartialInvite = Invite

// Invite Represents a code that when used, adds a user to a guild.
// https://discordapp.com/developers/docs/resources/invite#invite-object
// Reviewed: 2018-06-10
type Invite struct {
	Lockable `json:"-"`

	// Code the invite code (unique Snowflake)
	Code string `json:"code"`

	// Guild the guild this invite is for
	Guild *PartialGuild `json:"guild"`

	// Channel the channel this invite is for
	Channel *PartialChannel `json:"channel"`

	// ApproximatePresenceCount approximate count of online members
	ApproximatePresenceCount int `json:"approximate_presence_count,omitempty"`

	// ApproximatePresenceCount approximate count of total members
	ApproximateMemberCount int `json:"approximate_member_count,omitempty"`
}

var _ DeepCopier = (*Invite)(nil)
var _ discordDeleter = (*Invite)(nil)

func (i *Invite) deleteFromDiscord(s Session, flags ...Flag) error {
	if i.Code == "" {
		return &ErrorEmptyValue{info: "can not delete invite without the code field populate"}
	}

	_, err := s.DeleteInvite(i.Code, flags...)
	return err
}

// DeepCopy see interface at struct.go#DeepCopier
func (i *Invite) DeepCopy() (copy interface{}) {
	copy = &Invite{}
	i.CopyOverTo(copy)

	return
}

// InviteMetadata Object
// https://discordapp.com/developers/docs/resources/invite#invite-metadata-object
// Reviewed: 2018-06-10
type InviteMetadata struct {
	Lockable `json:"-"`

	// Inviter user who created the invite
	Inviter *User `json:"inviter"`

	// Uses number of times this invite has been used
	Uses int `json:"uses"`

	// MaxUses max number of times this invite can be used
	MaxUses int `json:"max_uses"`

	// MaxAge duration (in seconds) after which the invite expires
	MaxAge int `json:"max_age"`

	// Temporary whether this invite only grants temporary membership
	Temporary bool `json:"temporary"`

	// CreatedAt when this invite was created
	CreatedAt Time `json:"created_at"`

	// Revoked whether this invite is revoked
	Revoked bool `json:"revoked"`
}

var _ DeepCopier = (*InviteMetadata)(nil)

// voiceRegionsFactory temporary until flyweight is implemented
func inviteFactory() interface{} {
	return &Invite{}
}

type GetInviteParams struct {
	WithMemberCount bool `urlparam:"with_count,omitempty"`
}

var _ URLQueryStringer = (*GetInviteParams)(nil)

// GetInvite [REST] Returns an invite object for the given code.
//  Method                  GET
//  Endpoint                /invites/{invite.code}
//  Discord documentation   https://discordapp.com/developers/docs/resources/invite#get-invite
//  Reviewed                2018-06-10
//  Comment                 -
//  withMemberCount: whether or not the invite should contain the approximate number of members
func (c *Client) GetInvite(inviteCode string, params URLQueryStringer, flags ...Flag) (invite *Invite, err error) {
	if params == nil {
		params = &GetInviteParams{}
	}

	r := c.newRESTRequest(&httd.Request{
		Endpoint: endpoint.Invite(inviteCode) + params.URLQueryString(),
	}, flags)
	r.factory = inviteFactory

	return getInvite(r.Execute)
}

// DeleteInvite [REST] Delete an invite. Requires the MANAGE_CHANNELS permission. Returns an invite object on success.
//  Method                  DELETE
//  Endpoint                /invites/{invite.code}
//  Discord documentation   https://discordapp.com/developers/docs/resources/invite#delete-invite
//  Reviewed                2018-06-10
//  Comment                 -
func (c *Client) DeleteInvite(inviteCode string, flags ...Flag) (deleted *Invite, err error) {
	r := c.newRESTRequest(&httd.Request{
		Method:   httd.MethodDelete,
		Endpoint: endpoint.Invite(inviteCode),
	}, flags)
	r.factory = inviteFactory

	return getInvite(r.Execute)
}
