#!/usr/bin/env sh

ffmpeg -framerate 12 \
	-i './out%4d.png' \
	-filter_complex "color=black,format=rgb24[c];[c][0]scale2ref[c][i];[c][i]overlay=format=auto:shortest=1,setsar=1" \
	-pix_fmt yuva420p \
	out.gif
