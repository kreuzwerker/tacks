
# Tacks

[![Build Status](https://magnum.travis-ci.com/kreuzwerker/tacks.svg?token=1DipDsTWTNsz2XQgckKa&branch=master)](https://magnum.travis-ci.com/kreuzwerker/tacks)

Tacks makes it convenient to build and maintain AWS [CloudFormation](http://aws.amazon.com/cloudformation/) exectuable templates and snippets without resorting to an additional toolchain.

If offers the following features (in no particular order):

* Combine metadata (e.g. tags, the name of stack, rollback behaviour) with stack resources in one file
* Support for environments (e.g. production, staging ...) and environment-specific settings
* Executable (`+x`) templates
* Stack event viewing
* Interchangeably use YAML or JSON for describing a stack (since JSON is YAML [starting with YAML 1.2](http://www.yaml.org/spec/1.2/spec.html))

See `examples/` for full-blown examples on how to use `tacks`.
