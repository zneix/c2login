package main

var scopes = []string{
	"user_subscriptions",
	"user_blocks_edit",          // deprecated, replaced with "user:manage:blocked_users"
	"user_blocks_read",          // deprecated, replaced with "user:read:blocked_users"
	"user_follows_edit",         // deprecated, soon to be removed later since we now use "user:edit:follows"
	"channel_editor",            // for /raid
	"channel:moderate",          //
	"channel:read:redemptions",  //
	"chat:edit",                 //
	"chat:read",                 //
	"whispers:read",             //
	"whispers:edit",             //
	"channel_commercial",        // for /commercial
	"channel:edit:commercial",   // in case twitch upgrades things in the future (and this scope is required)
	"user:edit:follows",         // for (un)following
	"clips:edit",                // for clip creation
	"channel:manage:broadcast",  // for creating stream markers with /marker command, and for the /settitle and /setgame commands
	"user:read:blocked_users",   // for getting list of blocked users
	"user:manage:blocked_users", // for blocking/unblocking other users
	"moderator:manage:automod",  // for approving/denying automod messages
}
