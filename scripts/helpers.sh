set -e -x
set -o pipefail

function build() {
  VER=$(git describe --tags)
  BUILDDIR=${1:-"."}/${GOOS}-${GOARCH}

  mkdir -p $BUILDDIR
  go build -ldflags "-s -w -X main.Version=${VER}" -o $BUILDDIR/
  echo $VER >${BUILDDIR}/../version
}

function pkg_zip() {
  cd $1
  upx *
  chmod +x *
  zip x32-osc-bridge_$(cat ../version)_$(basename $PWD | tr '-' '_').zip *
  mv *.zip ..
  cd ..
}