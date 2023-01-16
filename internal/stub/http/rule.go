package http

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	urlPrefix      = []byte("URL:")
	delayPrefix    = []byte("DELAY:")
	headersPrefix  = []byte("HEADERS:")
	httpCodePrefix = []byte("HTTP_CODE")
	bodyPrefix     = []byte("BODY:")

	UnknownRuleError = errors.New("unknown rule for path")
)

type RuleSet struct {
	mtx     sync.RWMutex
	rules   map[string]*Rule
	regexps []regexpPathMap
}

type regexpPathMap struct {
	regexp *regexp.Regexp
	key    string
}

func newRuleSet() *RuleSet {
	return &RuleSet{
		mtx: sync.RWMutex{},
		//TODO: Think about hold unused data on disk
		rules:   make(map[string]*Rule, 0),
		regexps: make([]regexpPathMap, 0),
	}
}

func (rs *RuleSet) FindPath(path []byte) (*Rule, error) {
	if rule, exists := rs.rules[string(path)]; exists {
		return rule, nil
	}
	for _, regs := range rs.regexps {
		if regs.regexp.Match(path) {
			return rs.rules[regs.key], nil
		}
	}

	return nil, UnknownRuleError
}

func (rs *RuleSet) ScanDir(directory string) {
	fileCh := make(chan string, 100)
	ruleCh := make(chan *Rule, 100)
	parserWg := &sync.WaitGroup{}
	for c := 0; c < runtime.NumCPU(); c++ {
		parserWg.Add(1)
		go func(fileCh <-chan string, ruleCh chan<- *Rule, wg *sync.WaitGroup) {
			for file := range fileCh {
				data, readErr := os.ReadFile(file)
				if readErr != nil {
					logrus.Errorf("cannot read rule from file: %s. err: %s", file, readErr.Error())
				}
				ext := filepath.Ext(file)
				ruleCh <- parseRule(data, ext)
			}
			wg.Done()
		}(fileCh, ruleCh, parserWg)
	}
	readDir(directory, fileCh)
	close(fileCh)

	ruleFillWg := &sync.WaitGroup{}
	ruleFillWg.Add(1)
	go func(wg *sync.WaitGroup) {
		for rule := range ruleCh {
			rs.mtx.Lock()
			rs.rules[string(rule.Url)] = rule
			rgxp, regErr := regexp.Compile(string(rule.Url))
			if regErr == nil {
				rs.regexps = append(rs.regexps, regexpPathMap{
					regexp: rgxp,
					key:    string(rule.Url),
				})
			}

			rs.mtx.Unlock()
		}
		ruleFillWg.Done()
	}(ruleFillWg)

	parserWg.Wait()
	close(ruleCh)

	ruleFillWg.Wait()
}

func readDir(dir string, fileCh chan<- string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		logrus.Errorf("Cannot read dir: %s\n", dir)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			readDir(filepath.Join(dir, entry.Name()), fileCh)
		} else {
			fileCh <- filepath.Join(dir, entry.Name())
		}
	}
}

type Rule struct {
	Url      string            `yaml:"url" json:"url"`
	Delay    time.Duration     `yaml:"delay" json:"delay"`
	HttpCode string            `yaml:"httpCode" json:"httpCode"`
	Headers  map[string]string `yaml:"headers" json:"headers"`
	Body     string            `yaml:"body" json:"body"`
}

func parseRule(data []byte, ext string) *Rule {
	rule := &Rule{
		Headers:  make(map[string]string),
		HttpCode: "200",
	}

	switch ext {
	case ".yml":
		fallthrough
	case ".yaml":
		uErr := yaml.Unmarshal(data, rule)
		if uErr != nil {
			logrus.Errorf("Fail to unmarshal. Err = %s", uErr.Error())
		}
	default:
		parseOwnRule(data, rule)
	}

	return rule
}

// TODO: Think about parser modification for fast scan
func parseOwnRule(ruleData []byte, rule *Rule) *Rule {
	parts := bytes.Split(ruleData, []byte("#---#"))
	for _, part := range parts {
		part = bytes.TrimSpace(part)
		if bytes.HasPrefix(part, urlPrefix) {
			rule.Url = string(bytes.TrimSpace(bytes.TrimPrefix(part, urlPrefix)))
			continue
		}
		if bytes.HasPrefix(part, delayPrefix) {
			var pDelayErr error
			rule.Delay, pDelayErr = time.ParseDuration(string(bytes.TrimSpace(bytes.TrimPrefix(part, delayPrefix))))
			if pDelayErr != nil {
				rule.Delay = time.Duration(0)
			}
			continue
		}
		if bytes.HasPrefix(part, bodyPrefix) {
			rule.Body = string(bytes.TrimSpace(bytes.TrimPrefix(part, bodyPrefix)))
			continue
		}
		if bytes.HasPrefix(part, httpCodePrefix) {
			rule.HttpCode = string(bytes.TrimSpace(bytes.TrimPrefix(part, bodyPrefix)))
			continue
		}
		if bytes.HasPrefix(part, headersPrefix) {
			headersStr := bytes.TrimSpace(bytes.TrimPrefix(part, headersPrefix))
			scanner := bufio.NewScanner(bytes.NewReader(headersStr))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				kvArr := strings.Split(scanner.Text(), ":")
				rule.Headers[strings.TrimSpace(kvArr[0])] = strings.TrimSpace(kvArr[1])
			}
			continue
		}
	}

	return rule
}
