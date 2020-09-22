package project

import log "github.com/sirupsen/logrus"

//PullPackage pulls package to a library folder and builds them
func PullPackage(libFolder, repo, at string) {
	pullPackages(libFolder, []packageJson{
		{
			Repo: repo,
			At:   at,
		},
	})
	log.Info("Package pulled successfully!")
}
