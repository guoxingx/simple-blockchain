package main

import (
    "crypto/sha256"
)

// merkle tree struct
type MerkleTree struct {
    RootNode *MerkleNode
}

// merkle tree node
type MerkleNode struct {
    Left *MerkleNode
    Right *MerkleNode
    Data []byte
}

// create a new merkle node
// if both left and right was provided, param data will be ignored.
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
    mNode := MerkleNode{}

    if left == nil && right == nil {
        hash := sha256.Sum256(data)
        mNode.Data = hash[:]
    } else {
        // get data from left and right node.
        prevHashes := append(left.Data, right.Data...)
        hash := sha256.Sum256(prevHashes)
        mNode.Data = hash[:]
    }

    mNode.Left = left
    mNode.Right = right

    return &mNode
}

// create a new merkle tree
func NewMerkleTree(data [][]byte) *MerkleTree {
    var nodes []*MerkleNode

    for _, datum := range data {
        node := NewMerkleNode(nil, nil, datum)
        nodes = append(nodes, node)
    }

    return &MerkleTree{NewRootMerkleNode(nodes)}
}

// generate a root node by given nodes.
func NewRootMerkleNode(nodes []*MerkleNode) *MerkleNode {
    if len(nodes) == 1 {
        return nodes[0]
    }

    // to pair
    if len(nodes) % 2 != 0 {
        copy_node := NewMerkleNode(nil, nil, nodes[len(nodes) - 1].Data)
        nodes = append(nodes, copy_node)
    }

    var newLevel []*MerkleNode
    for i := 0; i < len(nodes); i += 2 {
        node := NewMerkleNode(nodes[i], nodes[i + 1], nil)
        newLevel = append(newLevel, node)
    }

    return NewRootMerkleNode(newLevel)
}
