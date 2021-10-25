// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yamlmeta

import "github.com/k14s/ytt/pkg/filepos"

// NewDocumentSet creates a DocumentSet containing `items`
func NewDocumentSet(items ...*Document) *DocumentSet {
    return &DocumentSet{
        Items:    items,
        Position: filepos.NewUnknownPosition(),
    }
}

// DocumentBuilder is builder used to configure and produce a Document.
// Use NewDocumentBuilder() to create instances.
type DocumentBuilder struct {
    newDoc *Document
}

// NewDocumentBuilder initializes a builder that produces a Document
func NewDocumentBuilder() *DocumentBuilder {
    return &DocumentBuilder{newDoc: &Document{}}
}

// Comment adds a comment to the Document
func (b *DocumentBuilder) Comment(data string, position *filepos.Position) *DocumentBuilder {
    comment := &Comment{
        Data:     data,
        Position: position,
    }
    b.newDoc.Comments = append(b.newDoc.Comments, comment)
    return b
}

// Value sets the value of the Document
func (b *DocumentBuilder) Value(val interface{}) *DocumentBuilder {
    b.newDoc.Value = val
    return b
}

// Position sets the filepos.Position of the Document
func (b *DocumentBuilder) Position(p *filepos.Position) *DocumentBuilder {
    b.newDoc.Position = p
    return b
}

// Build yields a Document as configured on this builder.
// If no position was configured (i.e. Position()), the position is set to "unknown"
func (b *DocumentBuilder) Build() *Document {
    if b.newDoc.Position == nil {
        b.newDoc.Position = filepos.NewUnknownPosition()
    }
    return b.newDoc.DeepCopy()
}

// BuildInDocumentSet wraps the configured Document as the one item in a new DocumentSet.
// Optionally, one can specify the filepos.Position of the DocumentSet.
// If no position is configured, the position is set to "unknown"
func (b *DocumentBuilder) BuildInDocumentSet(position ...*filepos.Position) *DocumentSet {
    docSetPosition := filepos.NewUnknownPosition()
    if len(position) > 0 {
        docSetPosition = position[0]
    }
    return &DocumentSet{
        Items:    []*Document{b.newDoc},
        Position: docSetPosition,
    }
}

// MapBuilder is builder used to configure and produce a Map (including its items).
// Use NewMapBuilder() to create instances.
type MapBuilder struct {
    newMap *Map
}

// NewMapBuilder initializes a builder that produces a Map (with items).
func NewMapBuilder() *MapBuilder {
    return &MapBuilder{newMap: &Map{}}
}

// Comment adds a comment to the Map, itself
func (b *MapBuilder) Comment(data string, position *filepos.Position) *MapBuilder {
    comment := &Comment{
        Data:     data,
        Position: position,
    }
    b.newMap.Comments = append(b.newMap.Comments, comment)
    return b
}

// Items adds `items` to the Map.
// This method is useful when the caller needs to control how the MapItems are created.
func (b *MapBuilder) Items(items ...*MapItem) *MapBuilder {
    b.newMap.Items = append(b.newMap.Items, items...)
    return b
}

// Item constructs and adds a MapItem to the Map.
// This method is useful when the MapItem has no other associated data than its key and value.
func (b *MapBuilder) Item(key interface{}, value interface{}, position *filepos.Position) *MapBuilder {
    newMapItem := &MapItem{
        Key:      key,
        Value:    value,
        Position: position,
    }
    b.newMap.Items = append(b.newMap.Items, newMapItem)
    return b
}

// Position sets the filepos.Position of the Map
func (b *MapBuilder) Position(position *filepos.Position) *MapBuilder {
    b.newMap.Position = position
    return b
}

// Build yields a Map as configured on this builder.
// If no position was configured (i.e. Position()), the position is set to "unknown"
func (b *MapBuilder) Build() *Map {
    if b.newMap.Position == nil {
        b.newMap.Position = filepos.NewUnknownPosition()
    }
    return b.newMap.DeepCopy()
}

// MapItemBuilder is builder used to configure and produce a MapItem
// Use NewMapItemBuilder() to create instances.
type MapItemBuilder struct {
    newItem *MapItem
}

// NewMapItemBuilder initializes a builder that produces a MapItem
func NewMapItemBuilder() *MapItemBuilder {
    return &MapItemBuilder{newItem: &MapItem{}}
}

// Comment adds a comment to the MapItem
func (b *MapItemBuilder) Comment(data string, position *filepos.Position) *MapItemBuilder {
    comment := &Comment{
        Data:     data,
        Position: position,
    }
    b.newItem.Comments = append(b.newItem.Comments, comment)
    return b
}

// Key sets the key of the MapItem
func (b *MapItemBuilder) Key(key interface{}) *MapItemBuilder {
    b.newItem.Key = key
    return b
}

// Value sets the value of the MapItem
func (b *MapItemBuilder) Value(value interface{}) *MapItemBuilder {
    b.newItem.Value = value
    return b
}

// Position sets the filepos.Position of the MapItem
func (b *MapItemBuilder) Position(position *filepos.Position) *MapItemBuilder {
    b.newItem.Position = position
    return b
}

// Build yields a MapItem as configured on this builder.
// If no position was configured (i.e. Position()), the position is set to "unknown"
func (b *MapItemBuilder) Build() *MapItem {
    if b.newItem.Position == nil {
        b.newItem.Position = filepos.NewUnknownPosition()
    }
    return b.newItem.DeepCopy()
}

// ArrayBuilder is builder used to configure and produce an Array (including its items).
// Use NewArrayBuilder() to create instances.
type ArrayBuilder struct {
    newArray *Array
}

// NewArrayBuilder initializes a builder that produces an Array (with items).
func NewArrayBuilder() *ArrayBuilder {
    return &ArrayBuilder{newArray: &Array{}}
}

// Position sets the filepos.Position of the Array
func (b *ArrayBuilder) Position(position *filepos.Position) *ArrayBuilder {
    b.newArray.Position = position
    return b
}

// Items adds `items` to the Array.
// This method is useful when the caller needs to control how the ArrayItems are created.
func (b *ArrayBuilder) Items(items ...*ArrayItem) *ArrayBuilder {
    b.newArray.Items = append(b.newArray.Items, items...)
    return b
}

// Item constructs and adds a ArrayItem to the Array.
// This method is useful when the ArrayItem has no other associated data than its key and value.
func (b *ArrayBuilder) Item(value interface{}, position *filepos.Position) *ArrayBuilder {
    newArrayItem := &ArrayItem{
        Value:    value,
        Position: position,
    }
    b.newArray.Items = append(b.newArray.Items, newArrayItem)
    return b
}

// Build yields an Array as configured on this builder.
// If no position was configured (i.e. Position()), the position is set to "unknown"
func (b *ArrayBuilder) Build() *Array {
    if b.newArray.Position == nil {
        b.newArray.Position = filepos.NewUnknownPosition()
    }
    return b.newArray
}


// ArrayItemBuilder is builder used to configure and produce an ArrayItem
// Use NewArrayItemBuilder() to create instances.
type ArrayItemBuilder struct {
    newItem *ArrayItem
}

// NewArrayItemBuilder initializes a builder that produces an ArrayItem
func NewArrayItemBuilder() *ArrayItemBuilder {
    return &ArrayItemBuilder{newItem: &ArrayItem{}}
}

// Comment adds a comment to the ArrayItem
func (b *ArrayItemBuilder) Comment(data string, position *filepos.Position) *ArrayItemBuilder {
    comment := &Comment{
        Data:     data,
        Position: position,
    }
    b.newItem.Comments = append(b.newItem.Comments, comment)
    return b
}

// Value sets the value of the ArrayItem
func (b *ArrayItemBuilder) Value(value interface{}) *ArrayItemBuilder {
    b.newItem.Value = value
    return b
}

// Position sets the filepos.Position of the ArrayItem
func (b *ArrayItemBuilder) Position(position *filepos.Position) *ArrayItemBuilder {
    b.newItem.Position = position
    return b
}

// Build yields an ArrayItem as configured on this builder.
// If no position was configured (i.e. Position()), the position is set to "unknown"
func (b *ArrayItemBuilder) Build() *ArrayItem {
    if b.newItem.Position == nil {
        b.newItem.Position = filepos.NewUnknownPosition()
    }
    return b.newItem
}