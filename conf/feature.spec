Name: feature
Version: 1.0.0
Release: 1.el7
Summary: http web server and cron scheduler
License: GPLv2
Group: Applications/System
Distribution: Linux
Packager: zilchzhong  <zilchzhong@163.com>
Buildarch: x86_64

%description
feature is a http web server and cron scheduler

%prep
%build
%pre
if [[ -f /opt/feature/conf/conf.toml ]]; then
    cp -a /opt/feature/conf/conf.toml{,.rpmsave}
fi

%post
systemctl daemon-reload
systemctl restart feature
systemctl enable feature

%postun
rm -rf /opt/feature

%files
/opt/feature/feature
/opt/feature/conf/conf.toml
/usr/lib/systemd/system/feature.service
/etc/logrotate.d/feature