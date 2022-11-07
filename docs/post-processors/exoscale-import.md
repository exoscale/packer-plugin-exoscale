# Exoscale Import Post-processor

Type: `exoscale-import`

The `exoscale-import` post-processor is used to import Exoscale custom
templates from disk image files (supported formats: `QCOW2`), e.g. an artifact
built locally by the [`qemu`][packer-doc-builder-qemu] builder.


### Required

- `api_key` (string) - The API key used to communicate with Exoscale
  services.

- `api_secret` (string) - The API secret used to communicate with Exoscale
  services.

- `image_bucket` (string) - The name of the bucket in which to upload the
  template image to SOS. The bucket must exist when the post-processor is
  run.

- `template_zones` (list of strings) - The Exoscale [zones][zones] in which to create the
  template.

- `template_zone` (string) - The Exoscale [zone][zones] in which to create the
  template. **DEPRECATED** (use `template_zones` instead).

- `template_name` (string) - The name to be used for registering the template.


### Optional

- `api_timeout` (int) - The maximum API async operations waiting time in seconds.
  Defaults to `3600`.

- `image_zone` (string) - The SOS Exoscale [zone][zones] in which to upload the template image.
  Defaults to the first of `template_zones`.

- `sos_endpoint` (string) - The endpoint used to communicate with SOS.
  Defaults to `https://sos-<image_zone>.exo.io`.

- `template_description` (string) - The description of the registered template.

- `template_username` (string) - An optional username to be used to log into
  Compute instances using this template.

- `template_boot_mode` (string) - The template boot mode.
  Supported values: `legacy` (default), `uefi`.

- `template_disable_password` (boolean) - Whether the registered template
  should disable Compute instance password reset. Defaults to `false`.

- `template_disable_sshkey` (boolean) - Whether the registered template
  should disable SSH key installation during Compute instance creation.
  Defaults to `false`.

- `skip_clean` (boolean) - Whether we should skip removing the image file
  uploaded to SOS after the import process has completed. "true" means that
  we should leave it in the bucket, "false" means deleting it.
  Defaults to `false`.


### Example Usage

```hcl
variable "api_key" { default = "" }
variable "api_secret" { default = "" }
variable "exoscale_zone" { default = "ch-gva-2" }

locals {
  image_name        = "base"
  image_format      = "qcow2"
  image_output_dir  = "output-qemu"
  image_username    = "ubuntu"
  image_output_file = "${local.image_output_dir}/${local.image_name}.${local.image_format}"
}

source "qemu" "base" {
  qemuargs = [
    ["-drive", "file=${local.image_output_file},format=${local.image_format},if=virtio"],
    ["-drive", "file=seed.img,format=raw,if=virtio"]
  ]
  cpus              = 4
  memory            = 4096
  vm_name           = "${local.image_name}.${local.image_format}"
  iso_url           = "https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img"
  iso_checksum      = "file:https://cloud-images.ubuntu.com/focal/current/SHA256SUMS"
  format            = local.image_format
  output_directory  = local.image_output_dir

  # ...
}

build {
  sources = ["source.qemu.base"]

  provisioner "shell" {
    environment_vars = ["DEBIAN_FRONTEND=noninteractive"]
    inline = [
      "sudo apt-get update && sudo apt-get upgrade -y",
      "sudo apt-get install --no-install-recommends -y ansible",
    ]
  }
}

source "file" "base" {
  source = local.image_output_file
  target = "${local.image_name}.${local.image_format}"
}

build {
  sources = ["source.file.base"]

  post-processor "exoscale-import" {
    api_key           = var.api_key
    api_secret        = var.api_secret
    image_bucket      = "my-templates-${var.exoscale_zone}"
    template_zones    = [var.exoscale_zone]
    template_name     = local.image_name
    template_username = local.image_username
  }
}
```


[packer-doc-builder-qemu]: https://www.packer.io/docs/builders/qemu
[zones]: https://www.exoscale.com/datacenters/
