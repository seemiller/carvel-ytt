// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yamlmeta_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/k14s/difflib"
	"github.com/k14s/ytt/pkg/filepos"
	"github.com/k14s/ytt/pkg/yamlmeta"
)

var _ = fmt.Sprintf

func TestParserDocSetEmpty(t *testing.T) {
	const data = ""

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().Position(filepos.NewPosition(1)).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetNewline(t *testing.T) {
	const data = "\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().Position(filepos.NewPosition(1)).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserOnlyComment(t *testing.T) {
	const data = "#"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().Position(filepos.NewPosition(1)).Build(),
		yamlmeta.NewDocumentBuilder().Comment("", filepos.NewPosition(1)).Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	if parsedValStr != expectedValStr {
		t.Fatalf("not equal\nparsed:\n>>>%s<<<expected:\n>>>%s<<<", parsedValStr, expectedValStr)
	}
}

func TestParserDoc(t *testing.T) {
	const data = "---\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().Position(filepos.NewPosition(1)).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocWithoutDashes(t *testing.T) {
	const data = "key: 1\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(yamlmeta.NewMapBuilder().
			Position(filepos.NewPosition(1)).
			Item("key", 1, filepos.NewPosition(1)).
			Build()).
		BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserRootValue(t *testing.T) {
	parserExamples{
		{Description: "string", Data: "abc",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value("abc").
				BuildInDocumentSet(),
		},
		{Description: "integer", Data: "1",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value(1).
				BuildInDocumentSet(),
		},
		{Description: "float", Data: "2000.1",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value(2000.1).
				BuildInDocumentSet(),
		},
		{Description: "float (exponent)", Data: "9e3",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value(9e3).
				BuildInDocumentSet(),
		},
		{Description: "array", Data: "- 1",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value(yamlmeta.NewArrayBuilder().
					Position(filepos.NewPosition(1)).
					Item(1, filepos.NewPosition(1)).Build()).
				BuildInDocumentSet(),
		},
		{Description: "map", Data: "key: val",
			Expected: yamlmeta.NewDocumentBuilder().
				Position(filepos.NewPosition(1)).
				Value(yamlmeta.NewMapBuilder().
					Position(filepos.NewPosition(1)).
					Item("key", "val", filepos.NewPosition(1)).Build()).
				BuildInDocumentSet(),
		},
	}.Check(t)
}

func TestParserRootString(t *testing.T) {
	expectedVal := yamlmeta.NewDocumentBuilder().
		Comment(" comment", filepos.NewPosition(1)).
		Position(filepos.NewPosition(1)).
		Value("abc").
		BuildInDocumentSet()

	parserExamples{
		// TODO should really be owned by abc
		{Description: "single line", Data: "--- abc # comment", Expected: expectedVal},
		{Description: "common on doc", Data: "--- # comment\nabc", Expected: expectedVal},
		// TODO add *yamlmeta.Value
		// {"comment on value", "---\nabc # comment", expectedVal},
	}.Check(t)
}

func TestParserMapArray(t *testing.T) {
	const data = `---
array:
- 1
- 2
- key: value
`

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(yamlmeta.NewMapBuilder().
			Position(filepos.NewPosition(1)).
			Items(
				yamlmeta.NewMapItemBuilder().
					Key("array").
					Position(filepos.NewPosition(2)).
					Value(
						yamlmeta.NewArrayBuilder().
							Position(filepos.NewPosition(2)).
							Item(1, filepos.NewPosition(3)).
							Item(2, filepos.NewPosition(4)).
							Item(
								yamlmeta.NewMapBuilder().
									Position(filepos.NewPosition(5)).
									Item("key", "value", filepos.NewPosition(5)).Build(),
								filepos.NewPosition(5)).
							Build(),
					).Build(),
			).Build(),
		).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserMapComments(t *testing.T) {
	const data = `---
# before-map
map:
  # before-key1
  key1: val1 # inline-key1
  # after-key1
  # before-key2
  key2: val2
`

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(yamlmeta.NewMapBuilder().
			Position(filepos.NewPosition(1)).
			Items(
				yamlmeta.NewMapItemBuilder().
					Comment(" before-map", filepos.NewPosition(2)).
					Position(filepos.NewPosition(3)).
					Key("map").
					Value(yamlmeta.NewMapBuilder().
						Position(filepos.NewPosition(3)).
						Items(
							yamlmeta.NewMapItemBuilder().
								Comment(" before-key1", filepos.NewPosition(4)).
								Key("key1").
								Value("val1").
								Position(filepos.NewPosition(5)).
								Comment(" inline-key1", filepos.NewPosition(5)).
								Build(),
							yamlmeta.NewMapItemBuilder().
								Comment(" after-key1", filepos.NewPosition(6)).
								Comment(" before-key2", filepos.NewPosition(7)).
								Key("key2").
								Value("val2").
								Position(filepos.NewPosition(8)).
								Build(),
						).Build(),
					).Build(),
			).Build(),
		).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserArrayComments(t *testing.T) {
	const data = `---
array:
# before-1
- 1 # inline-1
# after-1
# before-2
- 2
- 3
- # empty
- 
  # on-map
  key: value
# on-array-item-with-map
- key: value
`

	// TODO comment on top of scalar

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(
			yamlmeta.NewMapBuilder().
				Position(filepos.NewPosition(1)).
				Items(
					yamlmeta.NewMapItemBuilder().
						Position(filepos.NewPosition(2)).
						Key("array").
						Value(
							yamlmeta.NewArrayBuilder().
								Position(filepos.NewPosition(2)).
								Items(
									yamlmeta.NewArrayItemBuilder().
										Comment(" before-1", filepos.NewPosition(3)).
										Position(filepos.NewPosition(4)).
										Value(1).
										Comment(" inline-1", filepos.NewPosition(4)).
										Build(),
									yamlmeta.NewArrayItemBuilder().
										Comment(" after-1", filepos.NewPosition(5)).
										Comment(" before-2", filepos.NewPosition(6)).
										Position(filepos.NewPosition(7)).
										Value(2).
										Build(),
									yamlmeta.NewArrayItemBuilder().
										Position(filepos.NewPosition(8)).
										Value(3).
										Build(),
									yamlmeta.NewArrayItemBuilder().
										Position(filepos.NewPosition(9)).
										Comment(" empty", filepos.NewPosition(9)).
										Value(nil).
										Build(),
									yamlmeta.NewArrayItemBuilder().
										Position(filepos.NewPosition(10)).
										Value(yamlmeta.NewMapBuilder().
											Position(filepos.NewPosition(10)).
											Items(yamlmeta.NewMapItemBuilder().
												Position(filepos.NewPosition(12)).
												Comment(" on-map", filepos.NewPosition(11)).
												Key("key").
												Value("value").
												Build(),
											).Build(),
										).Build(),
									yamlmeta.NewArrayItemBuilder().
										Comment(" on-array-item-with-map", filepos.NewPosition(13)).
										Position(filepos.NewPosition(14)).
										Value(yamlmeta.NewMapBuilder().
											Position(filepos.NewPosition(14)).
											Items(yamlmeta.NewMapItemBuilder().
												Position(filepos.NewPosition(14)).
												Key("key").
												Value("value").
												Build(),
											).Build(),
										).Build(),
								).Build(),
						).Build(),
				).Build(),
		).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetComments(t *testing.T) {
	const data = `---
# comment-first
---
---
# comment-second
`

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(1)).
			Value(nil).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-first", filepos.NewPosition(2)).
			Position(filepos.NewPosition(3)).
			Value(nil).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(4)).
			Value(nil).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-second", filepos.NewPosition(5)).
			Value(nil).
			Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetOnlyComments2(t *testing.T) {
	const data = "---\n# comment-first\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(1)).
			Value(nil).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-first", filepos.NewPosition(2)).
			Position(filepos.NewUnknownPosition()).
			Value(nil).
			Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetOnlyComments3(t *testing.T) {
	const data = "--- # comment\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().
			Comment(" comment", filepos.NewPosition(1)).
			Position(filepos.NewPosition(1)).
			Value(nil).
			Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetOnlyComments(t *testing.T) {
	const data = "# comment-first\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(1)).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-first", filepos.NewPosition(1)).
			Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserDocSetCommentsNoFirstDashes(t *testing.T) {
	const data = `# comment-first
---
---
# comment-second
`

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentSet(
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(1)).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-first", filepos.NewPosition(1)).
			Position(filepos.NewPosition(2)).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Position(filepos.NewPosition(3)).
			Build(),
		yamlmeta.NewDocumentBuilder().
			Comment(" comment-second", filepos.NewPosition(4)).
			Build(),
	)

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserUnindentedComment(t *testing.T) {
	const data = `---
key:
  nested: true
# comment
  nested: true
`

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(
			yamlmeta.NewMapBuilder().
				Position(filepos.NewPosition(1)).
				Items(
					yamlmeta.NewMapItemBuilder().
						Position(filepos.NewPosition(2)).
						Key("key").
						Value(
							yamlmeta.NewMapBuilder().
								Position(filepos.NewPosition(2)).
								Item("nested", true, filepos.NewPosition(3)).
								Items(
									yamlmeta.NewMapItemBuilder().
										Comment(" comment", filepos.NewPosition(4)).
										Position(filepos.NewPosition(5)).
										Key("nested").
										Value(true).
										Build(),
								).Build(),
						).Build(),
				).Build(),
		).BuildInDocumentSet()

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(expectedVal)

	assertEqual(t, parsedValStr, expectedValStr)
}

func TestParserInvalidDoc(t *testing.T) {
	parserExamples{
		{Description: "no doc marker",
			Data:        "apiVersion: @123",
			ExpectedErr: "yaml: line 1: found character that cannot start any token",
		},
		{Description: "doc marker",
			Data:        "---\napiVersion: @123",
			ExpectedErr: "yaml: line 2: found character that cannot start any token",
		},
		{Description: "space before",
			Data:        "\n\n\napiVersion: @123",
			ExpectedErr: "yaml: line 4: found character that cannot start any token",
		},
		{Description: "doc marker with space",
			Data:        "\n\n---\napiVersion: @123",
			ExpectedErr: "yaml: line 4: found character that cannot start any token",
		},
	}.Check(t)
}

func TestParserAnchors(t *testing.T) {
	data := `
#@ variable = 123
value: &value
  path: #@ variable
  #@annotation
  args:
  - 1
  - 2
anchored_value: *value
`

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(
			yamlmeta.NewMapBuilder().
				Position(filepos.NewPosition(1)).
				Items(
					yamlmeta.NewMapItemBuilder().
						Comment("@ variable = 123", filepos.NewPosition(2)).
						Position(filepos.NewPosition(3)).
						Key("value").
						Value(
							yamlmeta.NewMapBuilder().
								Position(filepos.NewPosition(3)).
								Items(
									yamlmeta.NewMapItemBuilder().
										Position(filepos.NewPosition(4)).
										// TODO: should be here as well
										// Comment("@ variable", filepos.NewPosition(4)).
										Key("path").
										Value(nil).Build(),
									yamlmeta.NewMapItemBuilder().
										Comment("@annotation", filepos.NewPosition(5)).
										Position(filepos.NewPosition(6)).
										Key("args").
										Value(
											yamlmeta.NewArrayBuilder().
												Position(filepos.NewPosition(6)).
												Item(1, filepos.NewPosition(7)).
												Item(2, filepos.NewPosition(8)).
												Build(),
										).Build(),
								).Build(),
						).Build(),
					yamlmeta.NewMapItemBuilder().
						Position(filepos.NewPosition(9)).
						Key("anchored_value").
						Value(
							yamlmeta.NewMapBuilder().
								Position(filepos.NewPosition(9)).
								Items(
									yamlmeta.NewMapItemBuilder().
										Position(filepos.NewPosition(4)).
										Comment("@ variable", filepos.NewPosition(4)).
										Key("path").
										Value(nil).Build(),
									yamlmeta.NewMapItemBuilder().
										// TODO: should be here as well
										// Comment("@annotation", filepos.NewPosition(5)).
										Position(filepos.NewPosition(6)).
										Key("args").
										Value(
											yamlmeta.NewArrayBuilder().
												Position(filepos.NewPosition(6)).
												Item(1, filepos.NewPosition(7)).
												Item(2, filepos.NewPosition(8)).
												Build(),
										).Build(),
								).Build(),
						).Build(),
				).Build(),
		).BuildInDocumentSet()

	// TODO annotations are not properly assigned
	parserExamples{{Description: "with seq inside anchored data", Data: data, Expected: expectedVal}}.Check(t)
}

func TestParserMergeOp(t *testing.T) {
	data := `
#@ variable = 123
value: &value
  path: #@ variable
  #@annotation
  args:
  - 1
  - 2
merged_value:
  <<: *value
  other: true
`

	expectedVal := yamlmeta.NewDocumentBuilder().
		Position(filepos.NewPosition(1)).
		Value(
			yamlmeta.NewMapBuilder().
				Position(filepos.NewPosition(1)).
				Items(
					yamlmeta.NewMapItemBuilder().
						Comment("@ variable = 123", filepos.NewPosition(2)).
						Position(filepos.NewPosition(3)).
						Key("value").
						Value(
							yamlmeta.NewMapBuilder().
								Position(filepos.NewPosition(3)).
								Items(
									yamlmeta.NewMapItemBuilder().
										Position(filepos.NewPosition(4)).
										Key("path").
										// TODO: should be here as well
										// Comment("@ variable", filepos.NewPosition(4)).
										Build(),
									yamlmeta.NewMapItemBuilder().
										Comment("@annotation", filepos.NewPosition(5)).
										Position(filepos.NewPosition(6)).
										Key("args").
										Value(
											yamlmeta.NewArrayBuilder().
												Position(filepos.NewPosition(6)).
												Item(1, filepos.NewPosition(7)).
												Item(2, filepos.NewPosition(8)).
												Build(),
										).Build(),
								).Build(),
						).Build(),
					yamlmeta.NewMapItemBuilder().
						Position(filepos.NewPosition(9)).
						Key("merged_value").
						Value(
							yamlmeta.NewMapBuilder().
								Position(filepos.NewPosition(9)).
								Items(
									yamlmeta.NewMapItemBuilder().
										Position(filepos.NewPosition(4)).
										Key("path").
										Comment("@ variable", filepos.NewPosition(4)).
										Build(),
									yamlmeta.NewMapItemBuilder().
										// TODO: should be here as well
										// Comment("@annotation", filepos.NewPosition(5)).
										Position(filepos.NewPosition(6)).
										Key("args").
										Value(
											yamlmeta.NewArrayBuilder().
												Position(filepos.NewPosition(6)).
												Item(1, filepos.NewPosition(7)).
												Item(2, filepos.NewPosition(8)).
												Build(),
										).Build(),
								).
								Item("other", true, filepos.NewPosition(11)).
								Build(),
						).Build(),
				).
				Build(),
		).BuildInDocumentSet()

	// TODO annotations are not properly assigned
	parserExamples{{Description: "merge", Data: data, Expected: expectedVal}}.Check(t)
}

func TestParserDocWithoutDashesPosition(t *testing.T) {
	const data = "key: 1\n"

	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(data), "data.yml")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	parsedPosStr := parsedVal.Items[0].Position.AsString()
	expectedPosStr := "line data.yml:1"

	if parsedPosStr != expectedPosStr {
		t.Fatalf("not equal\nparsed...: %s\nexpected.: %s\n", parsedPosStr, expectedPosStr)
	}
}

type parserExamples []parserExample

func (exs parserExamples) Check(t *testing.T) {
	for _, ex := range exs {
		ex.Check(t)
	}
}

type parserExample struct {
	Description string
	Data        string
	Expected    *yamlmeta.DocumentSet
	ExpectedErr string
}

func (ex parserExample) Check(t *testing.T) {
	parsedVal, err := yamlmeta.NewParser(yamlmeta.ParserOpts{WithoutComments: false}).ParseBytes([]byte(ex.Data), "")
	if len(ex.ExpectedErr) == 0 {
		ex.checkDocSet(t, parsedVal, err)
	} else {
		ex.checkErr(t, err)
	}
}

func (ex parserExample) checkDocSet(t *testing.T, parsedVal *yamlmeta.DocumentSet, err error) {
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	printer := yamlmeta.NewPrinterWithOpts(os.Stdout, yamlmeta.PrinterOpts{ExcludeRefs: true})

	parsedValStr := printer.PrintStr(parsedVal)
	expectedValStr := printer.PrintStr(ex.Expected)

	assertEqual(t, parsedValStr, expectedValStr)
}

func (ex parserExample) checkErr(t *testing.T, err error) {
	if err == nil {
		t.Fatalf("expected error")
	}

	parsedValStr := err.Error()
	expectedValStr := ex.ExpectedErr

	assertEqual(t, parsedValStr, expectedValStr)
}

func assertEqual(t *testing.T, parsedValStr string, expectedValStr string) {
	t.Helper()
	if parsedValStr != expectedValStr {
		t.Fatalf("Not equal; -actual, +expected:\n%v\n", difflib.PPDiff(strings.Split(parsedValStr, "\n"), strings.Split(expectedValStr, "\n")))
	}
}
