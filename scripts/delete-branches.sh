#!/bin/bash
CURRENT_BRANCH=`git rev-parse --abbrev-ref HEAD`
declare -a git_branches
git_branches=`git for-each-ref --format='%(refname:short)' refs/heads/`

for v in ${git_branches[*]}
do 
    if [[ "$v" != "staging" && "$v" != "master" && "$v" != $CURRENT_BRANCH ]]; then
        git branch -D $v
    fi
done