# Exoscale Packer Plugin

[![Actions Status](https://github.com/exoscale/packer-plugin-exoscale/workflows/CI/badge.svg?branch=main)](https://github.com/exoscale/packer-plugin-exoscale/actions?query=workflow%3ACI+branch%3Amain)

The `exoscale` multi-plugin can be used with HashiCorp [Packer][packer]
to create [Compute instance custom templates][exo-doc-custom-templates].


## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation][packer-doc-init] for more information


#### Manual installation

You can find pre-built binary releases of the plugin [here][releases].
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin][packer-doc-plugins].


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `make build` from the root of the sources
directory. Upon successful compilation, a `packer-plugin-exoscale` plugin
binary file can be found in the `bin/` directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin][packer-doc-plugins].


### Configuration

For more information on how to configure the plugin, please read the
documentation located in the [`docs/`](docs) directory.


## Contributing

* If you think you've found a bug in the code or you have a question regarding
  the usage of this software, please reach out to us by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.


[packer-doc-plugins]: https://www.packer.io/docs/extending/plugins/#installing-plugins
[exo-doc-custom-templates]: https://community.exoscale.com/documentation/compute/custom-templates/
[packer-doc-init]: https://www.packer.io/docs/commands/init
[packer-doc-plugins]: https://www.packer.io/docs/extending/plugins/#installing-plugins
[packer]: https://www.packer.io/
[releases]: https://github.com/exoscale/packer-plugin-exoscale/releases
