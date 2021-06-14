# octant-plugin-crd-example

This project is an example of building dashboards for CRDs using
[octant][octant] plugins. The CRDs are generated using kubebuilder tools with
some controller-manager code to optionally write controllers for the CRDs.
Directory `plugin/` contains two plugin packages for creating a user and an
admin dashboard for the CRDs.

The user plugin and the admin plugin provide different set of controls in the
octant dashboard. The user plugin is supposed to be used by a restricted user
for viewing and updating limited options. The admin plugin is supposed to be
used by high privilege users who can modify all the resources. Shortcuts and
power user controls can be provided in the admin plugin for various
administrative tasks.

[octant]: https://octant.dev/

## Setup/Development

1. Install the CRDs in a cluster with `make install`.
2. Download and install octant.
3. Build the octant plugins with `make build-plugins`.
4. Install the octant plugins with `make install-plugins`.
5. Run octant. It should start and register the two plugins.

The example CRDs in this project are `Boxes` and `BoxRecords`. The example
configurations refer to a namespace called `box`. Create the namespace `box`
and install the sample instances of these CRDs with
`kubectl apply -k config/samples/`.

In the octant dashboard, opening the CRD instances should show extra tabs and
cards, created by the plugins.

From the left navigation bar, icons for the user and/or admin plugin can be
used to open the respective dashboards.

When developing, to rebuild and install the plugins run
`make build-install-plugins`. If octant is already using the old binaries, the
installation will fail with `Text file busy` error. Stop octant before
installing a new version of the plugin.

Refer https://reference.octant.dev/ for docs and examples of octant's UI
components.

### Restricted User

__NOTE__: Remove the admin plugin binary from `~/.config/octant/plugins/`
before trying as a restricted user to avoid seeing the admin plugin options.

To try a restricted user experience, create a service account and roles for a
restricted user with `kubectl apply -k plugin/octant-box-user/config/`. This
creates Role, ClusterRoles, RoleBinding, ClusterRoleBinding and a
ServiceAccount, in the `box` namespace with read only access for octant. Run
`bash hack/user-kubeconfig-gen.sh` to generate a restricted kubeconfig
(`box-user.kubeconfig`) using the created ServiceAccount. Run octant with the
created kubeconfig:

```console
$ octant --kubeconfig box-user.kubeconfig
```

This should open the octant dashboard with access to the `box` namespace only.
