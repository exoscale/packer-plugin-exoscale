# Changelog

## 0.5.2

### New

- Added arm64 builds #45 (thanks @nhedger!)

## 0.5.1

### Bug fixes

- Fix bug with user_data & user_data_file in builder #44

### Improvements

- go.mk: lint with staticcheck #41 
- go.mk: upgrade to v2.0.3 #43 

## 0.5.0

### Improvements

- builder: use egoscale GetTemplateByName on instance creation #37
- Bump golang.org/x/net to v0.17.0 #38

## 0.4.1

### Improvements

- automate release with Exoscale Tooling GPG key #36

## 0.4.0

### New

- Add support for user-data when launching instance ([#33](https://github.com/exoscale/packer-plugin-exoscale/pull/33))
- go.mk: standardize CI with other Go repos ([#32](https://github.com/exoscale/packer-plugin-exoscale/pull/32/))

## 0.3.2

### Bug fixes

- Fix `panic: ConfigSpec failed: gob: type cty.Type has no exported fields` ([#29](https://github.com/exoscale/packer-plugin-exoscale/pull/29))

## 0.3.1

### New

- Add support for downloading an exported instance snapshot ([#27](https://github.com/exoscale/packer-plugin-exoscale/pull/27), thanks to [sternik](https://github.com/sternik))

## 0.3.0

### New

- Add support for Build, Maintainer and Version attributes ([#23](https://github.com/exoscale/packer-plugin-exoscale/pull/23))

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
