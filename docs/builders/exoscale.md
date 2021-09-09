# Exoscale Builder

Type: `exoscale`

The `exoscale` builder is used to create Exoscale custom templates based on a
Compute instance snapshot.

**Note:** the `exoscale` Packer plugin only supports UNIX-like operating
systems (e.g. GNU/Linux, \*BSD...). To build Exoscale custom templates for
other OS, we recommend using the [QEMU][packerqemu] plugin combined with the
[exoscale-import][exoscale-import] Packer post-processor plugin.


### Required

- `api_key` (string) - The API key used to communicate with Exoscale services.

- `api_secret` (string) - The API secret used to communicate with Exoscale
  services.

- `instance_template` (string) - The name or ID of the template to use when
  creating the Compute instance.

- `template_zone` (string) - The Exoscale [zone][zones] in which to create the
  template.

- `template_name` (string) - The name of the template.


### Optional

- `instance_type` (string) - The instance type of the Compute instance.
  Defaults to `Medium`.

- `instance_name` (string) - The name of the Compute instance.
  Defaults to `packer-<BUILD ID>`.

- `instance_zone` (string) - The Exoscale zone in which to create the Compute
  instance. Defaults to the value of `template_zone`.

- `instance_template_visibility` (string) - The template visibility to specify
  for the `instance_template` parameter. Defaults to `public`.

- `instance_disk_size` (int) - Volume disk size in GB of the Compute instance
  to create. Defaults to `50`.

- `instance_security_groups` (list of strings) - List of Security Groups
  (names) to apply to the Compute instance. Defaults to `["default"]`.

- `instance_private_networks` (list of strings) - List of Private Networks
  (names) to attach to the Compute instance.

- `instance_ssh_key` (string) - Name of the Exoscale SSH key to use with the
  Compute instance. If unset, a throwaway SSH key named `packer-<BUILD ID>`
  will be created before creating the instance, and destroyed after a
  successful build.

- `template_description` (string) - The description of the template.

- `template_username` (string) - An optional username to be used to log into
  Compute instances using this template.

- `template_boot_mode` (string) - The template boot mode. Supported values:
  `legacy` (default), `uefi`.

- `template_disable_password` (boolean) - Whether the template should disable
  Compute instance password reset. Defaults to `false`.

- `template_disable_sshkey` (boolean) - Whether the template should disable
  SSH key installation during Compute instance creation. Defaults to `false`.

In addition to plugin-specific configuration parameters, you can also adjust
the [SSH communicator][packerssh] settings to configure how Packer will log
into the Compute instance.


### Example Usage

```hcl
variable "api_key" { default = "" }
variable "api_secret" { default = "" }

source "exoscale" "my-app" {
  api_key = var.api_key
  api_secret = var.api_secret
  instance_template = "Linux Ubuntu 20.04 LTS 64-bit"
  instance_security_groups = ["packer"]
  template_zone = "ch-gva-2"
  template_name = "my-app"
  template_username = "ubuntu"
  ssh_username = "ubuntu"
}

build {
  sources = ["source.exoscale.test"]

  provisioner "shell" {
    execute_command = "chmod +x {{.Path}}; sudo {{.Path}}"
    scripts = ["install.sh"]
  }
}
```


[packerssh]: https://www.packer.io/docs/communicators/ssh/
[zones]: https://www.exoscale.com/datacenters/
