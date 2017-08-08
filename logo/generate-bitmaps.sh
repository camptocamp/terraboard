#!/bin/sh
INKSCAPE="$(which inkscape)"
$INKSCAPE ./terraboard_logo.svg      --export-png ./terraboard_logo.png        -w300  
$INKSCAPE ./terraboard_logo_only.svg --export-png ./terraboard_logo_only.png   -w300
$INKSCAPE ./terraboard_logo_only.svg --export-png ./terraboard_logo_small.png  -h150
