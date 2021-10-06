## Application Programming Interfaces in ytt

(advice and rationale as to prefer CLI)


### Command-Line Interface

- `ytt` as a binary
- upgradable _in the field_ — decouple versioning of your tool and ytt (dynamic vs. static linking)
- enumerates directories into contained files (and directories, recursively)
- you don't have to worry about bumping Go dependencies (module management effort)
- harder to screw this up

ex:
- kapp-controller
- terraform provider
- 


### As a Go Module
Template (In Process): `template.Options.RunWithFiles()`

- `ytt` as a Go module
- upgrade of `ytt` implies upgrade of your app/tool
- only accepts files
- your users don't have to install/manage a separate binary

To evaluate templates:
1. add `ytt` in your `go.mod`
   ```go.mod
   ...
   require  (
       ...
       github.com/k14s/ytt v0.36.0
       ...
   )
   ```
   _(note: this is `ytt`'s original module name and will be renamed in the future.)_
2. create and populate an instance of `template.Options`:
   ```go
   opts := template.NewOptions()
   
   // equivalent to `--data-values-file`
   opts.DataValuesFlags.FromFiles = []string{"values.yml"}
   
   // equivalent to `--file`
   opts.
   ```
3. invoke `template.Options.RunWithFiles()` to templatize

- example: https://github.com/vmware-tanzu/carvel-kapp/blob/8e1d1f706da29d9f31e003dcf7a7a413f54de75e/pkg/kapp/yttresmod/overlay_contract_v1_mod.go#L40

TODO: note that we're commiting to backwards compatible behavior, not necessarily API, right now.

---

Aspirations:

- attach this to the README.md (or better how does this work in pkg.go.dev for other packages?)
- convert the examples above to an official example (compiled and tested during the build)
- an API that looks and feels closer to the CLI interface
  - for physical files
  - for in-memory files