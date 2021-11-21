package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var traceLogger *Logger

func init() {
	traceLogger = NewTraceLogger(logrus.DebugLevel, 1)
}

type dofRepo struct {
	workDirPath    string
	repoFolderPath string
	*git.Repository
}

func OpenDofRepo(workDir, repoFolder string) (*dofRepo, error) {
	traceLogger.SetLevel(viper.GetString("LogLevel"))
	traceLogger.Debugf("Open dof repo workDir %s repoFolder %s", workDir, repoFolder)
	wt := osfs.New(workDir)
	dot, err := wt.Chroot(repoFolder)
	if err != nil {
		return nil, err
	}

	traceLogger.Debug("dot.Root()", dot.Root())
	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())
	repo, err := git.Open(s, wt)
	if err != nil {
		return nil, err
	}
	dof := &dofRepo{workDirPath: workDir, repoFolderPath: path.Join(workDir, repoFolder), Repository: repo}

	dof.doNotShowUntrackedFiles()

	return dof, nil
}

func CheckoutDofRepo(workDir, repoFolder, repoUrl, branch string) (*dofRepo, error) {
	traceLogger.Debugf("CheckoutDofRepo workDir %s, repoFolder %s, repoUrl %s, branch %s", workDir, repoFolder, repoUrl, branch)
	worktree := osfs.New(workDir)
	dot, err := worktree.Chroot(repoFolder)
	if err != nil {
		return nil, err
	}

	traceLogger.Debug("dot.Root()", dot.Root())

	opts := &git.CloneOptions{
		URL:      repoUrl,
		Progress: os.Stdout,
	}

	_, err = git.PlainClone(path.Join(workDir, repoFolder), true, opts)
	if err != nil {
		traceLogger.Error(err)
		return nil, err
	}

	dofRepo, err := OpenDofRepo(workDir, repoFolder)
	if err != nil {
		traceLogger.Error(err)
		return nil, err
	}

	dofRepo.doNotShowUntrackedFiles()

	dofRepo.renameOldFiles()

	wt, err := dofRepo.Worktree()
	if err != nil {
		traceLogger.Error(err)
		return nil, err
	}

	checkoutOpts := &git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
		Keep:   true,
	}
	err = wt.Checkout(checkoutOpts)
	if err != nil {
		traceLogger.Error(err)
		return nil, err
	}

	return dofRepo, nil
}

func InitNewDofRepo(workDir, repoFolder, branch string) (*dofRepo, error) {
	// git init --bare $HOME/.cfg
	traceLogger.Debugf("Checkout using workdir %s repoFolder %s", workDir, repoFolder)
	_, err := git.PlainInit(path.Join(workDir, repoFolder), true)
	if err != nil {
		return nil, err
	}

	dofRepo, err := OpenDofRepo(workDir, repoFolder)
	if err != nil {
		return nil, err
	}

	// TODO: do we need to set the remote directly here?
	traceLogger.Debugf("Checkout %s branch\n", branch)
	br := &config.Branch{
		Name:   branch,
		Remote: "origin",
		Rebase: "true",
	}
	err = dofRepo.CreateBranch(br)
	if err != nil {
		return nil, err
	}

	dofRepo.doNotShowUntrackedFiles()

	err = dofRepo.addGitIgnore()
	if err != nil {
		return nil, err
	}

	return dofRepo, nil
}

func (dof *dofRepo) AddFile(file string) error {
	wt, err := dof.Worktree()
	if err != nil {
		traceLogger.Error(err)
		return err
	}

	_, err = wt.Add(file)
	if err != nil {
		traceLogger.Error(err)
		return err
	}

	opts := &git.CommitOptions{
		All: true,
	}
	_, err = wt.Commit(fmt.Sprintf("Add %s", file), opts)
	return err
}

func (dof *dofRepo) Status() ([]byte, error) {
	status := exec.Command("git", fmt.Sprintf("--git-dir=%s", dof.repoFolderPath), fmt.Sprintf("--work-tree=%s", dof.workDirPath), "status", "-s")

	return status.Output()
}

func (dof *dofRepo) addGitIgnore() error {
	gitIgnore := path.Join(dof.workDirPath, ".gitignore")
	file, err := os.Create(gitIgnore)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(file)

	linesToWrite := []string{dof.repoFolderPath}
	for _, line := range linesToWrite {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	writer.Flush()

	return dof.AddFile(gitIgnore)
}

func (dof *dofRepo) doNotShowUntrackedFiles() {
	// alias config='/usr/bin/git --git-dir=$HOME/.cfg/ --work-tree=$HOME'
	// config config --local status.showUntrackedFiles no
	cfg, err := dof.Repository.Config()
	if err != nil {
		panic(err)
	}
	cfg.Raw.SetOption("status", "", "showuntrackedfiles", "no")
	err = dof.Repository.SetConfig(cfg)
	if err != nil {
		panic(err)
	}
}

func (dof *dofRepo) renameOldFiles() {
	err := os.Chdir(dof.workDirPath)
	if err != nil {
		panic(err)
	}

	rev := plumbing.Revision("HEAD")
	hash, err := dof.ResolveRevision(rev)
	if err != nil {
		panic(err)
	}

	commit, err := dof.CommitObject(*hash)
	if err != nil {
		panic(err)
	}

	tree, err := commit.Tree()
	if err != nil {
		panic(err)
	}

	for _, entry := range tree.Entries {
		file := path.Join(dof.workDirPath, entry.Name)
		newName := path.Join(dof.workDirPath, entry.Name+"_before_dof")
		if _, err := os.Stat(file); err == nil {
			traceLogger.Debugf("Rename oldfile %s to newname %s.", entry.Name, newName)
			err := os.Rename(file, newName)
			if err != nil {
				traceLogger.Error(err)
			}
		}
	}
}
