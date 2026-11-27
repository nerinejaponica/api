#!/bin/sh
# PRINT A JEST-LIKE COVERAGE TABLE FROM A GO COVERAGE PROFILE
# USAGE: scripts/coverage-table.sh coverage.out

PROFILE="${1:-coverage.out}"
MODULE="github.com/nerinejaponica/api/"

# STEP 1: MERGE DUPLICATE BLOCKS (ONE PER TEST BINARY), THEN SORT BY FILE AND LINE
# PROFILE LINE FORMAT: file.go:startLine.col,endLine.col statements count
# SKIP sqlc GENERATED CODE, IT IS EXCLUDED FROM COVERAGE ON SONAR TOO
grep -v '^mode:' "$PROFILE" | grep -v '/internal/db/' | awk -v mod="$MODULE" '
{
	key = $1
	stmts[key] = $2
	count[key] += $3
}
END {
	for (key in stmts) {
		split(key, a, ":")
		split(a[2], r, ",")
		split(r[1], s, ".")
		split(r[2], e, ".")
		file = a[1]
		sub(mod, "", file)
		print file, s[1], e[1], stmts[key], count[key]
	}
}' | sort -k1,1 -k2,2n | awk '
# STEP 2: SUM STATEMENTS PER FILE AND COLLECT UNCOVERED LINE RANGES
{
	file = $1
	if (!(file in total)) order[++files] = file
	total[file] += $4
	if ($5 > 0) {
		covered[file] += $4
		next
	}
	# MERGE THIS BLOCK INTO THE PREVIOUS RANGE WHEN THEY TOUCH OR OVERLAP
	if ((file in curEnd) && $2 <= curEnd[file] + 1) {
		if ($3 > curEnd[file]) curEnd[file] = $3
		next
	}
	# OTHERWISE CLOSE THE PREVIOUS RANGE AND START A NEW ONE
	if (file in curEnd) flush(file)
	curStart[file] = $2
	curEnd[file] = $3
}
function flush(f,   range) {
	range = (curStart[f] == curEnd[f]) ? curStart[f] : curStart[f] "-" curEnd[f]
	uncovered[f] = (f in uncovered) ? uncovered[f] "," range : range
	delete curEnd[f]
}
END {
	for (f in curEnd) flush(f)
	line = "-----------------------------------+---------+-------+--------------------"
	print line
	printf "%-34s | %7s | %5s | %s\n", "File", "% Stmts", "Stmts", "Uncovered Lines"
	print line
	for (i = 1; i <= files; i++) {
		f = order[i]
		pct = total[f] ? covered[f] * 100 / total[f] : 100
		printf "%-34s | %6.1f%% | %5d | %s\n", f, pct, total[f], uncovered[f]
		allTotal += total[f]
		allCovered += covered[f]
	}
	print line
	printf "%-34s | %6.1f%% | %5d |\n", "All files", allTotal ? allCovered * 100 / allTotal : 100, allTotal
	print line
	# FILES THAT NO TEST BINARY LOADED ARE NOT IN THE PROFILE, SHOW THEM AS UNTESTED
	for (i = 1; i <= files; i++) seen[order[i]] = 1
	cmd = "find internal -name \"*.go\" ! -path \"internal/db/*\" ! -path \"internal/constant/*\" | sort"
	while ((cmd | getline f) > 0) {
		if (!(f in seen)) print "NOT LOADED BY ANY TEST: " f
	}
}'
