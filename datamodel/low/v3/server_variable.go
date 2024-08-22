package v3

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"gopkg.in/yaml.v3"
)

// ServerVariable represents a low-level OpenAPI 3+ ServerVariable object.
//
// ServerVariable is an object representing a Server Variable for server URL template substitution.
// - https://spec.openapis.org/oas/v3.1.0#server-variable-object
//
// This is the only struct that is not Buildable, it's not used by anything other than a Server instance,
// and it has nothing to build that requires it to be buildable.
type ServerVariable struct {
	Enum        []low.NodeReference[string]
	Default     low.NodeReference[string]
	Description low.NodeReference[string]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode     *yaml.Node
	RootNode    *yaml.Node
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the ServerVariable object.
func (s *ServerVariable) GetRootNode() *yaml.Node {
	return s.RootNode
}

// GetExtensions returns all Paths extensions and satisfies the low.HasExtensions interface.
func (s *ServerVariable) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return s.Extensions
}

// GetKeyNode returns the key yaml node of the ServerVariable object.
func (s *ServerVariable) GetKeyNode() *yaml.Node {
	return s.RootNode
}

// Build will extensions.
func (p *ServerVariable) Build(ctx context.Context, keyNode, root *yaml.Node) error {
	root = utils.NodeAlias(root)
	p.KeyNode = keyNode
	p.RootNode = root
	utils.CheckForMergeNodes(root)
	p.Reference = new(low.Reference)
	p.Nodes = low.ExtractNodes(ctx, root)
	p.Extensions = low.ExtractExtensions(root)
	low.ExtractExtensionNodes(ctx, p.Extensions, p.Nodes)
	return nil
}

// Hash will return a consistent SHA256 Hash of the ServerVariable object
func (s *ServerVariable) Hash() [32]byte {
	var f []string
	keys := make([]string, len(s.Enum))
	z := 0
	for k := range s.Enum {
		keys[z] = fmt.Sprint(s.Enum[k].Value)
		z++
	}
	sort.Strings(keys)
	f = append(f, keys...)
	if !s.Default.IsEmpty() {
		f = append(f, s.Default.Value)
	}
	if !s.Description.IsEmpty() {
		f = append(f, s.Description.Value)
	}
	f = append(f, low.HashExtensions(s.Extensions)...)
	return sha256.Sum256([]byte(strings.Join(f, "|")))
}
