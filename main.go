package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Meta struct {
	Title  string   `yaml:"title"`
	Slug   string   `yaml:"slug"`
	Tags   []string `yaml:"tags"`
	Date   string   `yaml:"date"`
	Author string   `yaml:"author"`
	Draft  bool     `yaml:"draft"`
}

type Article struct {
	Meta Meta
	Text string
	Path string
}

func main() {
	log.SetLevel(log.DebugLevel)

	fpath := flag.String("path", "", "Path with article files.")
	fauthor := flag.String("author", "marvin", "Default Author.")

	flag.Parse()

	articles := files(*fpath)

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, a := range articles {
		wg.Add(1) // nolint:gmnd
		go worker(a, *fauthor, &wg)
	}
}

func worker(art string, aut string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Infof("work on %s", art)

	a := Article{}
	a.Path = art
	a.Parse(aut)

	fmt.Print(a.Create())
}

// Parse reads file and parses it components
// nolint:funlen
func (a *Article) Parse(aut string) {
	dat, err := ioutil.ReadFile(a.Path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(dat)))
	text := []string{}

	// nolint:gmnd
	for scanner.Scan() {
		t := scanner.Text()

		title := regexp.MustCompile(`^Title:\s(.+)$`).FindStringSubmatch(t)
		if len(title) >= 2 {
			a.Meta.Title = title[1]
			continue
		}

		date := regexp.MustCompile(`^Date:\s(.+)$`).FindStringSubmatch(t)
		if len(date) >= 2 {
			d := date[1] + " CET"
			conTime, err := time.Parse("2006-01-02 15:04 MST", d)
			if err != nil {
				log.Fatal(err)
			}
			a.Meta.Date = conTime.Format(time.RFC3339)

			continue
		}

		slug := regexp.MustCompile(`^Slug:\s(.+)$`).FindStringSubmatch(t)
		if len(slug) >= 2 {
			a.Meta.Slug = slug[1]
			continue
		}

		author := regexp.MustCompile(`^Author:\s(.+)$`).FindStringSubmatch(t)
		if len(author) >= 2 {
			a.Meta.Author = author[1]
			continue
		} else {
			a.Meta.Author = aut
		}

		tags := regexp.MustCompile(`^Tags:\s(.+)$`).FindStringSubmatch(t)
		if len(tags) >= 2 {
			a.Meta.Tags = strings.Split(strings.ReplaceAll(tags[1], " ", ""), ",")
			continue
		}

		category := regexp.MustCompile(`^Category:\s(.+)$`).FindStringSubmatch(t)
		if len(category) >= 2 {
			continue
		}

		status := regexp.MustCompile(`^Status:\s(draft)$`).FindStringSubmatch(t)
		if len(status) >= 2 {
			a.Meta.Draft = true
			continue
		}

		// it must be text if nothing elses matches
		text = append(text, t)
	}

	if text[0] == "" {
		text = text[1:]
	}

	if text[len(text)-1] == "" {
		text = text[:len(text)-1]
	}

	a.Text = strings.Join(text, "\n")
}

// MetaYAML renders a YAML string out of the metadata struct.
func (a *Article) MetaYAML() string {
	out, err := yaml.Marshal(a.Meta)
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

// Create creates a new hugo styled article markdown file.
func (a *Article) Create() string {
	return "---\n" + a.MetaYAML() + "---\n" + a.Text
}

// Write writes the file back to path.
func (a *Article) Write() {
	if err := ioutil.WriteFile(a.Path, []byte(a.Create()), 0644); err != nil {
		log.Fatal(err)
	}
}

// Clean writes some tags new and cleans stuff out.
// TODO: giphy
// TODO: internal links
// TODO: soundcloud
// nolint: godox
func (a *Article) Clean() {
	// youtube
	yRe := regexp.MustCompile(`({%\syoutube\s(.+)\s%})`)
	if yRe.MatchString(a.Text) {
		new := yRe.ReplaceAllString(a.Text, "{{ <youtube $2> }}")
		a.Text = new
	}

	// vimeo
	vRe := regexp.MustCompile(`({%\svimeo\s(.+)\s%})`)
	if vRe.MatchString(a.Text) {
		new := vRe.ReplaceAllString(a.Text, "{{ <vimeo $2> }}")
		a.Text = new
	}
}

// files list a sorted slice of files in a directory.
func files(d string) []string {
	f := make([]string, 0)
	files, err := ioutil.ReadDir(d)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		matched, err := filepath.Match("*.md", file.Name())

		if err != nil {
			log.Fatal(err)
		}

		if matched {
			abs, err := filepath.Abs(filepath.Join(d, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			f = append(f, abs)
		}
	}

	return f
}
