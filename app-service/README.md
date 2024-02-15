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
- The response represents `TokenRequest` object described in https://ably.com/docs/api/realtime-sdk/types#token-request
- field `capability` represents a JSON object that lists all `channels` and their corresponding access rights for the authenticated user
- this endpoint is meant to be used on the client side to initialize Ably SDK

### Get my channels grouped or as a flat list

The term "my channels" means those `channels` which you, as a user, can have access to. This API is needed for UI to render possible options for the user to act upon. Effectively the list reflects the same information that is encoded in `TokenRequest` object returned by `/chat/token` but is presented in a human-readable format.

#### Grouped 

Endpoint `GET /chat/channels`

Exampled response:
```json
[
  {
    "id": "8482ba32-840b-4ccd-8d0f-ab5f6628bbcf",
    "title": "betting only",
    "summary": null,
    "resources": [
      {
        "id": "d0d784df-092f-465f-a479-9523a61ddb53",
        "title": "betting one",
        "resource": "chat:BettingOnly:betting one",
        "read-only": false
      },
      {
        "id": "0ea83a0c-02a8-4415-939b-2fe1a99bbcb5",
        "title": "betting two",
        "resource": "chat:BettingOnly:betting two",
        "read-only": false
      }
    ]
  },
  {
    "id": "cd54bef3-5793-4b3f-809d-8b73fad05f4c",
    "title": "nothing serious",
    "summary": null,
    "resources": [
      {
        "id": "f67b76ad-a313-4a6a-be6c-f389e20809f0",
        "title": "private one",
        "resource": "chat:NothingSerious:private one",
        "read-only": false
      }
    ]
  }
]
```

Comments:
- The top-level array shows `chat groups` and their corresponding `channels` if the user has permission to work with them
- A user who owns a `chat group` will see all its channels without explicitly joining to them
- Both **public** and **private** channels are shown as long as the user has corresponding right to work with that channels

#### Flat list

Endpoint `GET /chat/channels/list`

Same set of channels, but without group info. Exampled response:
```json
[
  {
    "id": "d0d784df-092f-465f-a479-9523a61ddb53",
    "title": "betting one",
    "resource": "chat:BettingOnly:bettingOne",
    "read-only": false
  },
  {
    "id": "0ea83a0c-02a8-4415-939b-2fe1a99bbcb5",
    "title": "betting two",
    "resource": "chat:BettingOnly:bettingTwo",
    "read-only": false
  },
  {
    "id": "f67b76ad-a313-4a6a-be6c-f389e20809f0",
    "title": "private one",
    "resource": "chat:NothingSerious:privateOne",
    "read-only": false
  }
]
```

### Invite user to join private channel

A user owning some **private** `chat group` can invite other existing users (one by one) to join any channel associated with this group. During the invitation process an email will be sent to invitee's email address and that email will contain a link to be followed to finalize the invitation. 

Endpoint `POST /chat/channels/{channelId}/invite`

Exampled request body:
```json
{
  "email": "abcdy@gmail.com"
}
```

### Accept invitation to join private channel

A user invited to join private channel will receive an email with a link. Once clicked, the link will be opened in the default web browser and a request to this API will be made under the hood. The link itself contains authentication token. The token allows to **re-identify** *invitor*, *invitee*, and *channel*, which upon successful validation, will activate access to the private channel for the invitee.

Endpoint `POST /chat/channels/accept`

Exampled request body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY0MDgyODUsInVzZXJJZCI6ImM2MTc0ZThhLWUxMmYtNGQ2NC1hNGZlLWEzYjBjMDgxYmQzMSIsImFjdGlvbiI6Ikludml0YXRpb25Ub1ByaXZhdGVDaGF0IiwiZXh0cmFDbGFpbXMiOnsiY2hhbm5lbElkIjoiZjY3Yjc2YWQtYTMxMy00YTZhLWJlNmMtZjM4OWUyMDgwOWYwIiwiaW52aXRlZUlkIjoiOWJlZjQxZWQtZmIxMC00NzkxLWIwMmUtOTZiMzcyYzA5NDY2In19.-TmdBdIe132bacWphpKdXAfMrx5OEup57Fdfyi8GD1k"
}
```

Comments:
- The token cannot be used for anything else except to accept invitation.
- It has expiration time of 24 hours.
- The invitor can repeat invitation process if token is expired.
- The original link in the email contains `token` as a query param, and it is responsibility of the web client to pass it over into POST request body, served by this API.
