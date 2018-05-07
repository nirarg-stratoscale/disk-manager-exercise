%define longhash %(git log | head -1 | awk '{print $2}')
%define shorthash %(echo %{longhash} | dd bs=1 count=12)

Name:    disk-manager-exercise-deploy
Version: 1.0
Release: 1.strato.%{shorthash}
Summary: disk-manager-exercise service deploy
Packager: Stratoscale Ltd
Vendor: Stratoscale Ltd
URL: http://www.stratoscale.com
#Source0: THIS_GIT_COMMIT
License: Strato

%define __strip /bin/true
%define __spec_install_port /usr/lib/rpm/brp-compress

%description
disk-manager-exercise service deployment

%build
cp %{_srcdir}/deploy/systemd.service .
cp %{_srcdir}/deploy/docker-compose.yml .
cp %{_srcdir}/deploy/monitor.json .
cp %{_srcdir}/deploy/nginx.conf .
echo '{"origin": "the service RPM", "version": "'%{longhash}'", "format": "v1"}' > installed-version.json

%install
install -p -D -m 644 systemd.service $RPM_BUILD_ROOT/usr/lib/systemd/system/disk-manager-exercise.service
install -p -D -m 644 docker-compose.yml $RPM_BUILD_ROOT/etc/stratoscale/compose/rootfs-star/disk-manager-exercise.yml
install -p -D -m 644 monitor.json $RPM_BUILD_ROOT/etc/stratoscale/clustermanager/services/control/disk-manager-exercise.service
install -p -D -m 644 nginx.conf $RPM_BUILD_ROOT/etc/nginx/conf.d/services/disk-manager-exercise.conf
install -p -D -m 644 installed-version.json $RPM_BUILD_ROOT/etc/stratoscale/clustermanager/service-versions/disk-manager-exercise.installed

%files
/usr/lib/systemd/system/disk-manager-exercise.service
/etc/stratoscale/compose/rootfs-star/disk-manager-exercise.yml
/etc/stratoscale/clustermanager/services/control/disk-manager-exercise.service
/etc/nginx/conf.d/services/disk-manager-exercise.conf
/etc/stratoscale/clustermanager/service-versions/disk-manager-exercise.installed
