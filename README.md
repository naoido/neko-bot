# 🐈 NEKO BOT 🐾

## Git flow
```mermaid
%%{init:{
    'gitGraph':{
        'showCommitLabel': false,
        'mainBranchOrder': 0
}}}%%

gitGraph
    commit
    commit
    branch release/1.0.x
    checkout release/1.0.x
    branch develop order: 1
    checkout develop
    commit
    branch feature/xx order: 2
    checkout feature/xx
    commit
    commit
    checkout develop
    merge feature/xx
    checkout release/1.0.x
    merge develop
    checkout main
    merge release/1.0.x
    checkout develop
    commit
    branch feature/bb order: 4
    checkout feature/bb
    commit
    checkout main
    commit
    branch hotfix/issue-xx
    commit
    checkout main
    merge hotfix/issue-xx
    checkout develop
    merge hotfix/issue-xx
    checkout feature/bb
    merge develop
    commit
    checkout develop
    branch "fix/issue-1" order: 3
    checkout "fix/issue-1"
    commit
    checkout develop
    merge feature/bb
    merge "fix/issue-1"
    checkout release/1.0.x
    merge develop
    checkout main
    merge release/1.0.x
```

## Clone
> git clone --recurse-submodules https://github.com/naoido/neko-bot.git

## .env
```env
DISCORD_TOKEN=YOUR_DISCORD_TOKEN
DISCORD_STATUS_TYPE=online
DISCORD_ACTIVITY_MESSAGE="Just chilling..."
DISCORD_COMMAND_PREFIX="!"
DEVELOPERS="279583321766494208"
```
