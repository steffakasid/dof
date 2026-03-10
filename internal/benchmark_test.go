package internal

import (
	"testing"

	"github.com/spf13/viper"
)

func BenchmarkStatus(b *testing.B) {

	viper.Set("LogLevel", "info")

	workDir := "/Users/sid/"
	repoFolderName := ".dof"

	dofRepo, err := OpenDofRepo(workDir, repoFolderName)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		_, err := dofRepo.Status()
		if err != nil {
			b.Fatal(err)
		}
	}
}
