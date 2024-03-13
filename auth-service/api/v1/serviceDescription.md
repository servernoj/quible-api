# Purpose

The `auth-service` provides the utility for working with the following entities
- Users' records
- Authorization tokens

Such entities are handled by set of operations (see below) allowing for
- Creation/Registration of new users
- Updating existing users
- Logging in with credentials associated with one of the existing users
- Resetting user password
- Retrieving complete user record for the currently logged in user
- Retrieving public user record (a.k.a. user profile) of an arbitrary user identified by their `id`
- Storing/retrieving user profile image

Every user record has a number of *required* and *optional* fields. Among those **required** we list the following
- Email
- Phone (at least 10 characters describing phone number in arbitrary format)
- Username (a.k.a. nickname)
- Full name (free format)

# Operations

The defined set of specific operations can be split into subsets addressing higher level goals. For example, user registration can be seen as based on at least 2 operations:
- Create user record
- Activate user record

In other words **operation** is a lower-level abstraction that can be combined with other operations to implement a higher-level **flow**. 