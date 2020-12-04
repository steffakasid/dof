# dof

This little tool can be used to manage dot file based on a git bare repository. The basic manual process is described e.g. https://www.atlassian.com/git/tutorials/dotfiles

## todo

* make it easier to intialize a new repo e.g. use an argument to directly set the remote

## how to use

### Initialize a new repository

If you want to initialize a new repo you could just run `dof init`. This will basically setup a local git repository which can be used to add the dot files via `dof add .dotfilename`. Right now you have to manually add a remote to be able to publish the dot file repository e.g. to github. You can do this with the following command `dof alias remote add <path-to-git-remote-repo>`. Afterward you can run `dof sync` to push all add files to the remote.

### Checkout and setup an existing dot file repository

If you want to checkout an existing repository you can just run `dof checkout <git-remote-repo>`. This command will checkout the repository as a bare repo. Afterwards it will identify all included dot files, rename them and run a checkout on the bare repo. Afterwards you will have your dotfiles setup in your home directory from the dot file repo.

## github actions

This repo uses https://github.com/marketplace/actions/go-release-binaries to build go binaries
