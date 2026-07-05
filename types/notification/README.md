# Notification type

This type defines a notification: a message to be delivered to a destination with a severity level. It is the shared contract used by services that produce and consume notifications.

Unlike the config types in this repo, `Notification` is not read from environment variables. It is a plain data type built through its constructor, which validates every field so an invalid notification cannot exist.

## Required properties

Both are constructor arguments and cannot be empty:

* **destination** defines where the notification is delivered
* **message** defines the notification content

## Severity level

The level is optional and set through the `WithLevel` functional option. The available levels, from least to most severe, are `Info`, `Warning` and `Error`. When no level is given, it defaults to `Info`.

`WithLevel` rejects values outside the defined levels, so a `Notification` can never hold an invalid level.

## Errors

The constructor returns sentinel errors that can be detected with `errors.Is`:

* `ErrEmptyProperty` when destination or message is empty; the error message names the offending property
* `ErrInvalidLevel` when the level passed to `WithLevel` is not one of the defined levels

## Usage

```go
// Default level (Info)
notification, err := notification.NewNotification("ops-team", "backup completed")
if err != nil {
    log.Fatal(err)
}

// Explicit level
notification, err = notification.NewNotification("ops-team", "disk almost full", notification.WithLevel(notification.Warning))
if err != nil {
    log.Fatal(err)
}

notification.Destination() // "ops-team"
notification.Message()     // "disk almost full"
notification.Level()       // notification.Warning
```
