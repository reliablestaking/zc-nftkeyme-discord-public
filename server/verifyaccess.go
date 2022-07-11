package server

import (
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// VerifyAccess rechecks that users are allowed access
func (s Server) VerifyAccess() {
	for true {
		logrus.Info("Verifying access...")
		discordUsers, err := s.Store.GetAllDiscordUsers()
		if err != nil {
			logrus.WithError(err).Error("Error getting all users")
		}

		for _, discordUser := range discordUsers {
			logrus.Infof("Verifying access for user %s", discordUser.DiscordUserID)
			t := oauth2.Token{
				RefreshToken: discordUser.NftkeymeRefreshToken.String,
			}

			tokenSource := s.NftkeymeOauthConfig.TokenSource(oauth2.NoContext, &t)
			newToken, err := tokenSource.Token()
			if err != nil {
				logrus.WithError(err).Error("Error getting token")
				continue
			}

			if newToken.AccessToken != discordUser.NftkeymeAccessToken.String {
				logrus.Infof("Updating discord user %s with new token", discordUser.DiscordUserID)
				err = s.Store.UpdateDiscordUser(discordUser.DiscordUserID, newToken.AccessToken, newToken.RefreshToken)
				if err != nil {
					logrus.WithError(err).Error("Error updating discord user")
					continue
				}
			}

			err = s.assignRoles(*newToken, discordUser.DiscordUserID)
			if err != nil {
				logrus.WithError(err).Error("Error assigning roles")
				continue
			}

			time.Sleep(5 * time.Second)
		}

		time.Sleep(24 * time.Hour)
	}
}

func (s Server) assignRoles(token oauth2.Token, discordUserID string) error {
	assets, err := s.NftkeymeClient.GetAssetsForUser(token.AccessToken, s.PolicyIDCheck)
	if err != nil {
		logrus.WithError(err).Error("Error getting assets")
		return err
	}
	logrus.Infof("Found %d chains", len(assets))

	assetsHunters, err := s.NftkeymeClient.GetAssetsForUser(token.AccessToken, s.PolicyIDCheckHunters)
	if err != nil {
		logrus.WithError(err).Error("Error getting assets")
		return err
	}

	logrus.Infof("Found %d hunters", len(assetsHunters))
	assets = append(assets, assetsHunters...)

	logrus.Infof("Found %d total assets for user %s", len(assets), discordUserID)

	//check for policy id
	numAssets := len(assets)
	logrus.Infof("Found %d assets for policy id %s", numAssets, s.PolicyIDCheck)

	// update num roles
	err = s.Store.UpdateDiscordUserNumAssets(discordUserID, numAssets)
	if err != nil {
		logrus.WithError(err).Error("Error updating number of assets")
		return err
	}

	// manage roles
	keys := make([]int, 0)
	for k, _ := range s.RoleMap {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	// add to first matching role, remove from rest
	roleFound := false
	for _, k := range keys {
		if numAssets >= k && !roleFound {
			logrus.Infof("Adding user %s to role %s", discordUserID, s.RoleMap[k])
			err = s.DiscordSession.GuildMemberRoleAdd(s.DiscordServerID, discordUserID, s.RoleMap[k])
			if err != nil {
				logrus.WithError(err).Error("Error adding user to role")
				return err
			}
			roleFound = true
		} else {
			logrus.Infof("Removing user %s from role %s", discordUserID, s.RoleMap[k])
			err = s.DiscordSession.GuildMemberRoleRemove(s.DiscordServerID, discordUserID, s.RoleMap[k])
			if err != nil {
				logrus.WithError(err).Error("Error removing user from role")
				return err
			}
		}
	}

	return nil
}
