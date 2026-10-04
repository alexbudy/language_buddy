package tui

import (
	"strconv"
	"strings"
)

const linebreak = "     ---------" // visual line break

// choiceItem is a single, optionally selectable entry in a choiceList.
type choiceItem struct {
	label      string
	value      string
	selectable bool // is this item selectable or not
}

func newChoiceItem(label string, value string) choiceItem {
	return choiceItem{label: label, value: value, selectable: true}
}

func newSeparatorItem() choiceItem {
	return choiceItem{label: linebreak, selectable: false}
}

// choiceList is a small carousel-style menu: up/down (with wraparound) moves
// the cursor, enter selects. Skips non-selectable items.
type choiceList struct {
	title  string
	items  []choiceItem
	cursor int
}

func newChoiceList(title string, items []choiceItem) choiceList {
	return choiceList{title: title, items: items, cursor: 0}
}

// scroll up, skip selectable items
func (c *choiceList) up() {
	if len(c.items) == 0 {
		return
	}

	for {
		c.cursor--

		if c.cursor < 0 {
			c.cursor = len(c.items) - 1
		}

		if c.items[c.cursor].selectable {
			return
		}
	}
}

// scroll down, skip selectable items
func (c *choiceList) down() {
	if len(c.items) == 0 {
		return
	}

	for {
		c.cursor = (c.cursor + 1) % len(c.items)
		if c.items[c.cursor].selectable {
			return
		}

		if c.items[c.cursor].selectable {
			return
		}
	}
}

// return the selected choiceItem
func (c choiceList) selected() choiceItem {
	return c.items[c.cursor]
}

func (c choiceList) view() string {
	var b strings.Builder

	if c.title != "" {
		b.WriteString(titleStyle.Render(c.title))
	}
	b.WriteString("\n")

	number := 1
	for i, item := range c.items {
		if item.selectable == false { // new line entries should be skipped
			b.WriteString(item.label + "\n")
			continue
		}

		itemIdxToDisplay := strconv.Itoa(number) + ". "
		number++
		if i == c.cursor {
			b.WriteString(cursorStyle.Render("> "))
			b.WriteString(selectedStyle.Render(itemIdxToDisplay + item.label))
		} else {
			b.WriteString("  " + itemIdxToDisplay)
			b.WriteString(item.label)
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	return b.String()
}

func labeledItems(labels ...string) []choiceItem {
	items := make([]choiceItem, len(labels))
	for i, l := range labels {
		items[i] = newChoiceItem(l, l)
	}
	return items
}