# anaconda-tools
Tools and scripts I use to complete tasks

## `bin` folder

|script                                                 | description |
|:------------------------------------------------------|:----------|
| [bin/c3i-one-off](./bin/c3i-one-off)                  | simple script for submitting my pipelines to Concourse; wraps [`c3i`](https://github.com/anaconda-distribution/conda-concourse-ci). |
| [bin/copy-from-aarch64](./bin/copy-from-aarch64)      | script to move manual packages from the **linux-aarch64** worker to zeus. |
| [bin/copy-from-s390x](./bin/copy-from-s390x)          | script to move manual packages from the **linux-s390x** worker to zeus. |
| [bin/git-fetch-master](./bin/git-fetch-master)        | fetches + pulls origin/{master,main} |
| [bin/git-prune-branches](./bin/git-prune-branches)    | fetches + pulls origin/{master,main} and prunes any local branches that have been deleted on origin |
| [bin/prefect-rsync](./bin/prefect-rsync)              | script that relocated Prefect built packages into current directory. Used to combine Concourse + Prefect artifacts for one copy. |
| [bin/promote-to-main](./bin/promote-to-main)          | copies local `conda` packages to `/www/pkgs/main`, preserving parent (subdir) |
