#!/bin/bash
bookie-cli() {
	tmpfile=$(mktemp)
	bookie --open-with-type bash --open-with-file "$tmpfile" "$@"
    # shellcheck disable=SC1090
	source "$tmpfile"
	rm "$tmpfile"
}
