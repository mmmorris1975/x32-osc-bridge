set -e -x
set -o pipefail

NAME="x32-osc-bridge"

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
  zip ${NAME}_$(cat ../version)_$(basename $PWD | tr '-' '_').zip *
  mv *.zip ../artifacts/
  cd ..
}

function pkg_rpm() {
  RPM_ARCH="x86_64"
  if [[ `echo $d | grep -sc arm64$` > 0 ]]
  then
    RPM_ARCH="aarch64"
  elif [[ `echo $d | grep -sc arm$` > 0 ]]
  then
    RPM_ARCH="armv7hl"
  fi

  cd $1
  fpm --verbose -s dir -t rpm --name $NAME --version $(cat ../version) --license MIT --architecture $RPM_ARCH \
    --provides $NAME --description $NAME --url "https://github.com/mmmorris1975/$NAME" --maintainer 'mmmorris1975@github' \
    --rpm-user bin --rpm-group bin --rpm-digest sha1 --prefix /usr/local/bin *
  mv *.rpm ../artifacts/
  cd ..
}

function pkg_deb() {
  DEB_ARCH="amd64"
  if [[ `echo $d | grep -sc arm64$` > 0 ]]
  then
    DEB_ARCH="arm64"
  elif [[ `echo $d | grep -sc arm$` > 0 ]]
  then
    DEB_ARCH="armhf"
  fi

  cd $1
  fpm --verbose -s dir -t deb --name $NAME --version $(cat ../version) --license MIT --architecture $DEB_ARCH \
    --provides $NAME --description $NAME --url "https://github.com/mmmorris1975/$NAME" --maintainer 'mmmorris1975@github' \
    --deb-user bin --deb-group bin --prefix /usr/local/bin *
  mv *.deb ../artifacts/
  cd ..
}