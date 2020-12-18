function build() {
  VER=$(git describe --tags)
  BUILDDIR=${1:-"."}/${GOOS}-${GOARCH}

  mkdir -p $BUILDDIR
  go build -ldflags "-s -w -X main.Version=${VER}" -o $BUILDDIR/
  echo $VER >${BUILDDIR}/../version
}