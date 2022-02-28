$tag = $(git describe --abbrev=0 --tags)
$hash = $(git log -1 --format="%h")  
qtdeploy  -ldflags "-X main.verTag=$tag -X main.verCommitHash=$hash" -debug 2>debug.log