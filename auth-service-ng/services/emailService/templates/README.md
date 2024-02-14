# Rules on authoring templates

A template should

- extend default layout `../../../../assets/acorn/layout.pug`
- define block `filter` to specify generated **renderer function** signature, e.g. `:go:func Demo(first string, second int)`
- define block `content` while using Pug mixins or other markup and reference params of the **renderer function**
- optionally define block `preheader` (see https://litmus.com/blog/the-ultimate-guide-to-preview-text-support)

Example:
```jade
extends ../../../../assets/acorn/layout.pug

block filter
  :go:func Demo(first string, second int)

block content
  +Row
    +C1-of-3
      .one hello, #{first}
    +C2-of-3
      .one #{second + 5}
```