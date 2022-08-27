package server

import "github.com/st2projects/ssh-sentinel-server/model/db"

func CheckPrincipals(allowed []db.Principal, requested []string) bool {

	var allowedPrincipals = true

	for _, request := range requested {

		var requestAllowed = false

		for _, allow := range allowed {

			if allow.Principal == request {
				requestAllowed = true
				break
			}

		}
		allowedPrincipals = allowedPrincipals && requestAllowed
	}

	return allowedPrincipals
}
