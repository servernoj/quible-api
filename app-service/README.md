# Chat support

## Introduction

All chat related endpoints are combined under `/chat` path. Altogether these endpoints are meant to store, organize and retrieve chat metadata to support client operations with Ably SDK. A client is assumed to act on behalf of an authenticated user, who in turn, will be associated with specific metadata records. 

Chat data is organized into `chat group` and `channels` which are related to each other as `1:M`. In other words every `chat group` can be associated with zero or more `channels`.

Users are associated with `chat groups` and are known as `owners`. Relationship between `channels` and users is implicit, i.e. via `chat groups`. For example, `UserA` can own `chat group A` which is a parent for two channels: `channel A` and `channel B`.

There are 2 types of `chat groups`: **public** and **private**. A **public** `chat group` is one, associated with `channels` that can be freely joined by any user. A **private** `chat group` requires *invitation* from the owner to let a user to join its `channels`.

Both `chat groups` and `channels` can be created and deleted. Public `channels` (ones associated with public `chat group`) can be *listed* (searched), *joined*, and *left* by means of API calls. 

Lastly, having all metadata combined, a special API call can produce `TokenRequest` to be used by the client to initialize Ably SDK on behalf of the authenticated user. Such request (JSON object) would list all `channels` along with associated permissions which the authenticated user can access.

## APi endpoints

### Create `chat group`

Endpoint `POST /chat/groups` 

Exampled request:
```json
{
  "name": "BettingOnly",
  "title": "betting only",
  "summary": "optional summary for the chat group",
  "isPrivate": false
}
```

Exampled response:
```json
{
  "id": "196e445d-a122-45c0-bc20-01e932da0583",
  "resource": "chat:BettingOnly",
  "summary": "optional summary for the chat group",
  "title": "betting only",
  "parent_id": null,
  "is_private": false,
  "owner_id": "9bef41ed-fb10-4791-b02e-96b372c09466"
}
```

Comments:
- `name` field in the request allows only alphabetic characters. It has to be **unique** across other `chat groups` owned by the same user
- `title` field in the request should also be unique across other `chat groups` of the same user
- `isPrivate` field in request is optional and defaults to `false`
- `resource` field in response is a concatenation of the hardcoded string `chat:` and the value of `name` field from the request
- `parent_id` field in response is `null` for `chat groups` and holds ID of the parent `chat group` for `channels`
- `owner_id` in response is the ID of the authenticated user who made this API call

### Create `channel`

Endpoint `POST /chat/channels?chatGroupId=xxx`

Exampled request:
```json
{
  "name": "channel A",
  "title": "channel in public chat group"
}
```

Exampled response:
```json
{
  "id": "39cf7d18-de17-4573-9826-458634ce7ebd",
  "resource": "channel A",
  "summary": null,
  "title": "channel in public chat group",
  "parent_id": "196e445d-a122-45c0-bc20-01e932da0583",
  "is_private": null,
  "owner_id": null
}
```

Comments:
- request may contain `summary` field which is omitted here
- `parent_id` in response is set with the value of request's query parameter `chatGroupId`
- both `is_private` and `owner_id` are set to `null` for `channels`

### Join public channel

Endpoint `POST /chat/channels/{channelId}`

There is no request/response body associated with the endpoint. The path param `channelId` must be set as ID of the channel to which we want to join.

Comments:
- Both `channel` and its holding `chat group` must exist
- If requested channel is associated with a **private** `chat group`, an error will be returned
- An attempt to join a channel associated with the **self-owned** `chat group` will fail with error

### List my `chat groups`

Endpoint `GET /chat/groups`

Exampled response
```json
[
  {
    "id": "8555c8f8-53dc-4a41-a1c1-dc8369fc37f7",
    "resource": "chat:lessie",
    "summary": null,
    "title": "world is good",
    "parent_id": null,
    "is_private": false,
    "owner_id": "9bef41ed-fb10-4791-b02e-96b372c09466"
  },
  {
    "id": "196e445d-a122-45c0-bc20-01e932da0583",
    "resource": "chat:BettingOnly",
    "summary": null,
    "title": "betting only",
    "parent_id": null,
    "is_private": false,
    "owner_id": "9bef41ed-fb10-4791-b02e-96b372c09466"
  }
]
```

Comments:
- The API returns the list of `chat groups` **owned** by the user making the request
- Both **public** and **private** `chat groups` are to be listed in response
- `channels` will not be listed in response

### Search public `chat groups`

Endpoint `GET /chat/groups/search?q=xxx`

The query parameter `q` will be partially matched with all public `chat groups`. When omitted, all public `chat groups` will be listed in response.

Exampled response (`&q=betting`)
```json
[
  {
    "chatGroup": {
      "id": "196e445d-a122-45c0-bc20-01e932da0583",
      "resource": "chat:BettingOnly",
      "summary": null,
      "title": "betting only",
      "parent_id": null,
      "is_private": false,
      "owner_id": "9bef41ed-fb10-4791-b02e-96b372c09466"
    },
    "channels": [
      {
        "id": "39cf7d18-de17-4573-9826-458634ce7ebd",
        "resource": "channel A",
        "summary": null,
        "title": "channel in public chat group",
        "parent_id": "196e445d-a122-45c0-bc20-01e932da0583",
        "is_private": null,
        "owner_id": null
      }
    ]
  }
]
```

Comments:
- The response is an array of records containing matching `chat group` and [an array of] **all** its `channels`
- `chat groups` owned by other users will also be examined and could be returned if title matching is satisfied. 

### Leave `channel`

Endpoint `DELETE /chat/channels/{channelId}`

There is no request/response body associated with the endpoint. The path param `channelId` must be set as the ID of the channel to be left.

Comments:
- An attempt to leave non-existing or not previously joined `channel` will result in an error
- It doesn't matter if the holding `chat group` is private or public -- you can always leave the channel if you have previously were allowed to join it

### Delete one of your owned `chat group`

Endpoint `DELETE /chat/groups/{chatGroupId}`

There is no request/response body associated with the endpoint. The path param `chatGroupId` must be set as the ID of the `chat group` to be removed.

Comments:
- You must own a chat group to be able to delete it
- All `channels` associated with the `chat group` in question will be deleted as well

### Get Ably token token

Endpoint `GET /chat/token`

Exampled response
```json
{
  "ttl": 3600000,
  "capability": "{\"chat:BettingOnly:*\":[\"subscribe\",\"publish\",\"history\"],\"chat:lessie:*\":[\"subscribe\",\"publish\",\"history\"],\"chat:simon:channel_public\":[\"subscribe\",\"publish\",\"history\"]}",
  "clientId": "9bef41ed-fb10-4791-b02e-96b372c09466",
  "timestamp": 1706157590130,
  "keyName": "OzADbA.wQsEWA",
  "nonce": "acf3812506fe9f80f0302bceb52abc5f",
  "mac": "aq5w68PSTFYfIyBT49snOQTsBNEV8LNwPaVRHk6CQhE="
}
```

Comments:
- The response represents `TokenRequest` object described in [https://ably.com/docs/api/realtime-sdk/types#token-request]
- field `capability` represents a JSON object that lists all `channels` and their corresponding access rights for the authenticated user
- this endpoint is meant to be used on the client side to initialize Ably SDK