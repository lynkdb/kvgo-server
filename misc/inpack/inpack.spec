[project]
name = kvgo-server
version = 0.9.0
vendor = lynkdb.com
homepage = https://github.com/lynkdb/kvgo-server
groups = dev/sys-srv

%build
export PATH=$PATH:/usr/local/go/bin:/opt/gopath/bin
export GOPATH=/opt/gopath

mkdir -p {{.buildroot}}/etc
mkdir -p {{.buildroot}}/bin
mkdir -p {{.buildroot}}/misc
mkdir -p {{.buildroot}}/var/log
mkdir -p {{.buildroot}}/var/data

go build -ldflags "-X main.version={{.project__version}} -X main.release={{.project__release}}" -o {{.buildroot}}/bin/kvgo-server cmd/server/main.go

rm -rf /tmp/rpmbuild/*
mkdir -p /tmp/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS,BUILDROOT}

mkdir -p /tmp/rpmbuild/SOURCES/kvgo-server-{{.project__version}}/
rsync -av {{.buildroot}}/* /tmp/rpmbuild/SOURCES/kvgo-server-{{.project__version}}/

sed -i 's/__version__/{{.project__version}}/g' /tmp/rpmbuild/SOURCES/kvgo-server-{{.project__version}}/misc/rpm/rpm.spec
sed -i 's/__release__/{{.project__release}}/g' /tmp/rpmbuild/SOURCES/kvgo-server-{{.project__version}}/misc/rpm/rpm.spec

cd /tmp/rpmbuild/SOURCES/
tar zcf kvgo-server-{{.project__version}}.tar.gz kvgo-server-{{.project__version}}

rpmbuild --define "debug_package %{nil}" -ba /tmp/rpmbuild/SOURCES/kvgo-server-{{.project__version}}/misc/rpm/rpm.spec \
  --define='_tmppath /tmp/rpmbuild' \
  --define='_builddir /tmp/rpmbuild/BUILD' \
  --define='_topdir /tmp/rpmbuild' \
  --define='dist .{{.project__dist}}'

find /tmp/rpmbuild/RPMS/ -type f -name *.rpm |xargs  -I '{}' mv {} {{.inpack__pack_dir}} 

%files
misc/


