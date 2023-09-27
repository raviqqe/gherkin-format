package main_test

import (
	"strings"
	"testing"

	"github.com/cucumber/gherkin/go/v27"
	"github.com/raviqqe/gherkin-format"
	"github.com/stretchr/testify/assert"
)

func TestNewRenderer(t *testing.T) {
	main.NewRenderer()
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

  Scenario: Bar
		`,
		`
Feature: Foo
  Scenario: Bar
    baz

    Given blah
		`,
		`
Feature: Foo
  Background: Bar
    When I do something`,
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
		`,
		`
Feature: Highlander
  Rule: There can be only One
    Scenario: Only One -- More than one alive
      Given there are 3 ninjas
      And there are more than one ninja alive
      When 2 ninjas meet, they will fight
      Then one ninja dies (but not me)
      And there is one ninja less alive

    Scenario: Only One -- One alive
      Given there is only 1 ninja alive
      Then he (or she) will live forever ;-)

  Rule: There can be Two (in some cases)
    Scenario: Two -- Dead and Reborn as Phoenix
		`,
		`
Feature: Foo
  @foo @bar
  Scenario: Bar
		`,
		`
# foo
Feature: Foo
		`,
		`
# foo
# bar
Feature: Foo
		`,
		`
Feature: Foo
  # foo
  Scenario: Bar
		`,
		`
Feature: Foo
  # foo
  # bar
  Scenario: Bar
		`,
		`
Feature: Foo
  Scenario: Bar
    # foo
    Given Baz
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
    # foo
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
    # foo
    # bar
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
    """
    """
    # foo
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
    """
    foo
    """
    # foo
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
    """
    """
    # foo

  # bar
  Scenario: Bar
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
      | foo |
      | bar |
    # foo
		`,
		`
Feature: Foo
  Scenario: Bar
    Given Baz
      | foo |
      | bar |
    # foo

  # bar
  Scenario: Bar
		`,
	} {
		s := strings.TrimSpace(s)

		t.Run(strings.ReplaceAll(s, "\n", " "), func(t *testing.T) {
			d, err := gherkin.ParseGherkinDocument(strings.NewReader(s), func() string { return "" })

			assert.Nil(t, err)
			assert.Equal(t, strings.TrimSpace(s)+"\n", main.NewRenderer().Render(d))
		})
	}
}

func TestRendererRenderCodeBlockMultipleTimes(t *testing.T) {
	s := strings.TrimSpace(`
Feature: Foo
  Scenario: Bar
    Given Baz
     """foo
    bar
    """
  `)
	u := strings.TrimSpace(`
Feature: Foo
  Scenario: Bar
    Given Baz
    """foo
    bar
    """
  `)

	d, err := gherkin.ParseGherkinDocument(strings.NewReader(s), func() string { return "" })

	assert.Nil(t, err)
	assert.Equal(t, u+"\n", main.NewRenderer().Render(d))
}

func TestRendererRenderTrimSpace(t *testing.T) {
	d, err := gherkin.ParseGherkinDocument(strings.NewReader("Feature:  foo  bar\tbaz"), func() string { return "" })

	assert.Nil(t, err)
	assert.Equal(t, "Feature: foo bar baz\n", main.NewRenderer().Render(d))
}

func TestRendererRenderNoDuplicateCommentAmongScenarioAndStep(t *testing.T) {
	s := strings.TrimSpace(`
Feature: Foo
  # foo
  Scenario: Bar
    Given Baz
  `)

	d, err := gherkin.ParseGherkinDocument(strings.NewReader(s), func() string { return "" })

	assert.Nil(t, err)
	assert.Equal(t, s+"\n", main.NewRenderer().Render(d))
}

func TestRendererRenderEscapedCharacters(t *testing.T) {
	s := strings.TrimSpace(`
Feature: Foo
  Scenario Outline: Bar
    Given Put <value>

    Examples:
      | value  |
      | \n     |
      | \\     |
      | \\t    |
      | \\r    |
  `)

	d, err := gherkin.ParseGherkinDocument(strings.NewReader(s), func() string { return "" })

	assert.Nil(t, err)
	assert.Equal(t, s+"\n", main.NewRenderer().Render(d))
}
