builder := go
builddir := bin
exe := ctfinstancer
path := /usr/local/bin
importdir := middlewares models router utils
instdir := /usr/local/share/$(exe)
instrelativedir := $(subst /usr/local,..,$(instdir))
systemddir := /etc/systemd/system
config := .env
systemd := $(exe).service
tags := release
ldflags := -s -w

all: $(builddir)/$(exe) $(builddir)/template

$(builddir)/$(exe): main.go go.mod go.sum $(importdir)
		$(builder) build -o $(builddir)/$(exe) -tags $(tags) -ldflags "$(ldflags)" $<

$(builddir)/template: template
		cp -r template $(builddir)/template

install: $(path)/$(exe) $(systemddir)/$(systemd)

$(path)/$(exe): $(instdir)/$(exe) $(instdir)/$(config) $(instdir)/template
		ln -s $(instrelativedir)/$(exe) $(path)/$(exe)

$(instdir): 
		mkdir $(instdir)

$(instdir)/$(exe): $(instdir) $(builddir)/$(exe)
		cp $(builddir)/$(exe) $(instdir)/$(exe)
		chown root:root $(instdir)/$(exe)
		chmod 4755 $(instdir)/$(exe)

$(instdir)/$(config): $(instdir) $(config).sample
		cp $(config).sample $(instdir)/$(config)

$(instdir)/template: $(instdir) $(builddir)/template
		cp -r $(builddir)/template $(instdir)/template

$(systemddir)/$(systemd): $(systemd)
		cp $(systemd) $(systemddir)/$(systemd)

uninstall:
		rm -rf $(path)/$(exe)
		rm -rf $(instdir)
		rm -rf $(systemddir)/$(systemd)
clean: 
		rm -rf $(builddir)
