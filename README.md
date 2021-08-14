# golang-bbs
Yet another old-school BBS written by Golang

## Settings
golang-bbs reads following envvars
* GOLANG_BBS_WORKDIR
* GOLANG_BBS_TMPLDIR
* PORT

### GOLANG_BBS_WORKDIR
A working directory which golang-bbs stores data files.
If not set, automatically set to cwd.

### GOLANG_BBS_TMPLDIR
A directory which contain 'template' dir. Please note that this is NOT template directory it self, but its parent directory.
If not set, automatically set to cwd.

### PORT
A port number which golang-bbs will listen to.
If not set, automatically set to 8080.
