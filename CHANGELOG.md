# Changelog

## next

### New

- Add support for Build, Maintainer and Version attributes

## 0.2.1

### Bug Fixes

- Fix incomplete HCL2 (config) specification ([#18](https://github.com/exoscale/packer-plugin-exoscale/pull/18))

## 0.2.0

### New

- Allow to create a template in several zones at once ([#15](https://github.com/exoscale/packer-plugin-exoscale/pull/15))

### Changes

- The `template_zone` parameter has been replaced by `template_zones`

## 0.1.3

### Bug Fixes

- Allow undefined `template_description` or `template_username` fields ([#14](https://github.com/exoscale/packer-plugin-exoscale/pull/14))
- Allow custom API timeouts, and set default to 1h in the builder instead of 5mn ([#13](https://github.com/exoscale/packer-plugin-exoscale/pull/13))

## 0.1.2

### Changes

- The `instance_template_filter` parameter has been replaced by `instance_template_visibility`.

### Bug Fixes

- Change SSH RSA key by modern ED25519 key ([#8](https://github.com/exoscale/packer-plugin-exoscale/pull/8))


## 0.1.1

Internal dependencies upgrade


## 0.1.0

Initial release
