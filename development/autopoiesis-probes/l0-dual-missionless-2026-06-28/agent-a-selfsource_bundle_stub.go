//go:build !selfsourcebundle

package quine

// SelfSourceBundle is empty in ordinary package-load builds. Builds that need
// the runtime source-code/ projection generate selfsource_bundle_data.go and set
// the selfsourcebundle build tag.
var SelfSourceBundle []byte
