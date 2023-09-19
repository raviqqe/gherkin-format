package main

import (
	"strings"
	"testing"

	"github.com/cucumber/gherkin/go/v27"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderer(t *testing.T) {
	newRenderer()
}

func TestRendererRender(t *testing.T) {
	for _, s := range []string{
		"Feature: Foo",
		`
Feature: Foo
  Scenario: Bar
    Given that
    When I do something
    Then something happens
		`,
		`
Feature: Foo
  Scenario: Bar
    When I do something:
    """sh
    foo
    """
		`,
		`
Feature: Foo

  bar
		`,
		`
Feature: Foo
  Scenario: Bar

    baz
		`,
		`
Feature: Foo
  Background: Bar
    When I do something`, `
# Foo

## Background (Bar)

_When_ I do something.
		`,
		`
Feature: Foo
  Background: Bar
  Given Baz:
    | foo |
    | bar |
		`,
		`
Feature: Foo
  Scenario Outline: Bar
    When <someone> does <something>.
    Examples:
      | someone | something |
      | I       | cooking   |
      | You     | coding    |
			`,
		`
Feature: Foo
  Scenario Outline: Bar
    When <someone> does <something>.
    Examples: Baz
      | someone | something |
      | I       | cooking   |
      | You     | coding    |
			`,
		`
Feature: Foo
  Scenario Outline: Bar
    When <someone> does <something>.
    Examples: Baz
      foo bar baz.

      | someone | something |
      | I       | cooking   |
      | You     | coding    |
			`,
		`
Feature: Foo
  Scenario Outline: Bar
    When <someone> does <something>.
    Examples: Baz
      | someone |
      | I       |
      | You     |
    Examples: Blah
      | something |
      | cooking   |
      | coding    |
`, `
Feature: Foo
  Rule: Bar
    Example: Baz
      When qux
`,
	} {
		d, err := gherkin.ParseGherkinDocument(strings.NewReader(s), func() string { return "" })

		assert.Nil(t, err)
		assert.Equal(t, strings.TrimSpace(s)+"\n", newRenderer().Render(d))
	}
}
