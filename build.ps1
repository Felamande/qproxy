taskkill /f /im qproxy.exe
$tag = $(git describe --abbrev=0 --tags)
$hash = $(git log -1 --format="%h")  
qtdeploy  -ldflags "-X main.verTag=$tag -X main.verCommitHash=$hash"
copy -Path psh_host.dll -Destination deploy\windows\psh_host.dll