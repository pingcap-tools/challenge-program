# Community Tools

## Syncer

A clean tool to sync github pull requests by polling. Without webhook setting, it runs individually.

## PCP

Performance Challenge Program backend with github API support.

### API

#### Rank

* `/api/rank` query rank of this season

* `/api/rank/season/:season` query rank of given season

Response:

```json
[
  {
    "rank": 1,
    "type": "team",
    "name": "PingCAP",
    "community": false,
    "url": "https://github.com/tidb-perf-challenge/pcp/issues/1",
    "score": 1140,
    "doing-task":"https://github.com/pingcap/tidb/issues/14486"
  },
  ...
]
```

#### Tasks

* `/api/task` query all tasks

* `/api/task/level/:level` query tasks of some level(easy/medium/hard)

* `/api/task/owner/:owner/repo/:repo` query tasks of some repo, eg. `/api/tasks/owner/pingcap/repo/tidb`

Response:

```json
[
  {
    "season": 2,
    "inprogress-user": null,
    "complete-user": {
      "ID": 0,
      "user": "you06",
      "avatar": "https://avatars3.githubusercontent.com/u/9587680?s=460&v=4",
      "github": "https://github.com/you06/"
    },
    "inprogress-team": null,
    "complete-team": null,
    "owner": "pingcap",
    "repo": "tidb",
    "title": "task 1",
    "issue": "https://github.com/pingcap/tidb/issues/10467",
    "level": "medium",
    "score": 1000,
    "status": "success"
  },
  {
    "season": 2,
    "inprogress-user": [
      {
        "ID": 0,
        "user": "you06",
        "avatar": "https://avatars3.githubusercontent.com/u/9587680?s=460&v=4",
        "github": "https://github.com/you06/"
      },
      {
        "ID": 0,
        "user": "illyrix",
        "avatar": "https://avatars3.githubusercontent.com/u/12008675?s=460&v=4",
        "github": "https://github.com/illyrix/"
      }
    ],
    "complete-user": null,
    "inprogress-team": null,
    "complete-team": null,
    "owner": "pingcap",
    "repo": "tidb",
    "title": "task 2",
    "issue": "https://github.com/pingcap/tidb/issues/7546",
    "level": "hard",
    "score": 6000,
    "status": "success"
  }
]
```

#### Task Groups

* `/api/taskgroup` query task groups

Response:

```json
[
  {
    "season": 2,
    "owner": "pingcap",
    "repo": "tidb",
    "title": "task group 1",
    "issue-number": 14486,
    "issue": "https://github.com/pingcap/tidb/issues/14486",
    "progress": 60,
    "doing-users": [
      {
        "ID": 0,
        "user": "you06",
        "avatar": "https://avatars3.githubusercontent.com/u/9587680?s=460u0026v=4",
        "github": "https://github.com/you06/"
      },
      {
        "ID": 0,
        "user": "illyrix",
        "avatar": "https://avatars3.githubusercontent.com/u/12008675?s=460u0026v=4",
        "github": "https://github.com/illyrix/"
      }
    ]
  },
  {
    "season": 2,
    "owner": "tikv",
    "repo": "tikv",
    "title": "task group 2",
    "issue-number": 6519,
    "issue": "https://github.com/tikv/tikv/issues/6519",
    "progress": 30,
    "doing-users": [
      {
        "ID": 0,
        "user": "you06",
        "avatar": "https://avatars3.githubusercontent.com/u/9587680?s=460u0026v=4",
        "github":"https://github.com/you06/"
      }
    ]
  }
]
```