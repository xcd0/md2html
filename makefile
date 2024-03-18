all:
	# release-build
	make -C src all

build:
	make -C src build

cp:
	bash -c 'if [ "$(uname -o)" == "Msys" ]; then cp "build/md2html_windows.exe" "`which md2html`"; else cp "build/md2html_linux" "`which md2html`"; fi'

