# Rules on authoring templates

A template should

- extend default layout `../../../../assets/acorn/layout.pug`
- define block `filter` to specify generated **renderer function** signature, e.g. `:go:func Activation(first string, second int)`
- define block `content` and use mixins or other markup while referencing params of the **renderer function**
- optionally define block `preheader` (see https://litmus.com/blog/the-ultimate-guide-to-preview-text-support)

Example:
```jade
extends ../../../../assets/acorn/layout.pug

block filter
  :go:func Activation(first string, second int)

block content
  +Row
    +C1-of-3
      .one hello, #{first}
    +C2-of-3
      .one #{second + 5}
```