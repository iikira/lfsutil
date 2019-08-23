# lfsutil
A simple tool for accessing git lfs server.

# Install
```
go get -u -v github.com/iikira/lfsutil
```

# Usage

Setup environment `LFS_REPO_URL` and `LFS_AUTH`
```
export LFS_REPO_URL=https://github.com/user/repo.git/info/lfs
export LFS_AUTH=Basic dXNlcjpwYXNzd29yZA==
```

## Get object infomation and download link

### By oid and size
```
lfsutil go oid:size
lfsutil go oid
```

#### example
```
lfsutil go e6fd9c1b536033f3346b32c391bd58587ea9f549cab7839cf8a1dbc62a739825:3862852
lfsutil go e6fd9c1b536033f3346b32c391bd58587ea9f549cab7839cf8a1dbc62a739825
```

### By local file
```
lfsutil go -by file /path/to/file
```

### By pointer file
Specification: https://github.com/git-lfs/git-lfs/blob/master/docs/spec.md
```
lfsutil go -by ptr /path/to/ptr/file
lfsutil go -by ptr http://example.com/ptr/file
```

## Upload local file
```
lfsutil uo /path/to/file
```
