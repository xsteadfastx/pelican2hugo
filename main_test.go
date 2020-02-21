package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArticleParse(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		path    string
		article Article
	}{
		{
			filepath.Join("testdata", "1up-berlin.md"),
			Article{
				Meta{
					Title:  "1UP berlin",
					Date:   "2011-12-20T12:46:00+01:00",
					Slug:   "1up-berlin",
					Tags:   []string{"art", "berlin", "documentary", "graffiti"},
					Author: "marvin",
					Draft:  true,
				},
				"Jeder der einmal durch Berlin gelaufen ist kennt sie.\n\n{% youtube QXxXoSTPivA %}",
				filepath.Join("testdata", "1up-berlin.md"),
				"",
			},
		},
		{
			filepath.Join("testdata", "zwei-neue-american-football-songs.md"),
			Article{
				Meta{
					Title:  "Zwei neue American Football Songs",
					Date:   "2019-01-24T14:52:00+01:00",
					Slug:   "zwei-neue-american-football-songs",
					Tags:   []string{"americanfootball", "emo"},
					Author: "marvin",
					Draft:  false,
				},
				"Als ich als 16 jähriges LiveJournal Einträge durchforstete.\n\n{% youtube CaZUVZ2F_Dc %}\n\n{% youtube q1XUaXk92KA %}", // nolint:lll
				filepath.Join("testdata", "zwei-neue-american-football-songs.md"),
				"",
			},
		},
	}

	for _, table := range tables {
		a := Article{}
		a.Path = table.path
		a.Parse("marvin")
		assert.Equal(table.article, a)
	}
}

func TestArticleMetaYAML(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		yaml string
		path string
	}{
		{
			`title: 1UP berlin
slug: 1up-berlin
tags:
- art
- berlin
- documentary
- graffiti
date: "2011-12-20T12:46:00+01:00"
author: marvin
draft: true
`,
			filepath.Join("testdata", "1up-berlin.md"),
		},
	}

	for _, table := range tables {
		a := Article{}
		a.Path = table.path
		a.Parse("marvin")
		assert.Equal(table.yaml, a.MetaYAML())
	}
}

func TestArticleCreate(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		path string
		a    string
	}{
		{
			filepath.Join("testdata", "1up-berlin.md"),
			`---
title: 1UP berlin
slug: 1up-berlin
tags:
- art
- berlin
- documentary
- graffiti
date: "2011-12-20T12:46:00+01:00"
author: marvin
draft: true
---
Jeder der einmal durch Berlin gelaufen ist kennt sie.

{% youtube QXxXoSTPivA %}`,
		},
	}

	for _, table := range tables {
		a := Article{}
		a.Path = table.path
		a.Parse("marvin")
		assert.Equal(table.a, a.Create())
	}
}

func TestArticleCLean(t *testing.T) {
	assert := assert.New(t)
	tables := []struct {
		old string
		new string
	}{
		{
			"{% youtube IvTNBbFkq4w %}",
			"{{< youtube IvTNBbFkq4w >}}",
		},
		{
			"{% vimeo 28938294 %}",
			"{{< vimeo 28938294 >}}",
		},
		{
			"{% youtube foo %}\n{% youtube bar %}",
			"{{< youtube foo >}}\n{{< youtube bar >}}",
		},
		{
			"[![cc-by-sa Santaduck]({static}/images/Obeyshepard2.jpg)](https://en.wikipedia.org/wiki/File:Obeyshepard2.jpg)",
			"[![cc-by-sa Santaduck](/images/Obeyshepard2.jpg)](https://en.wikipedia.org/wiki/File:Obeyshepard2.jpg)",
		},
		{
			"[Artikel]({static}/posts/meine-neue-shell-xonsh.md)\n[postmodernen Neubaugebiet]({static}/posts/kerksiek-006.md)",
			"[Artikel]({{< ref \"/posts/meine-neue-shell-xonsh.md\" >}})\n[postmodernen Neubaugebiet]({{< ref \"/posts/kerksiek-006.md\" >}})", // nolint:lll
		},
	}

	for _, table := range tables {
		a := Article{}
		a.Text = table.old
		a.Clean()
		assert.Equal(table.new, a.Text)
	}
}
