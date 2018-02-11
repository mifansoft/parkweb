// Copyright 2016 The mifanpark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

// Ini type configuration file
type INIConfig struct {
	val  map[string]string
	lock *sync.RWMutex
	file string
}

// key value pair
type iniLineOne struct {
	line   string
	index  int
	offset int
}

func createINIConfig(path string) (*INIConfig, error) {
	conf := &INIConfig{}
	conf.val = make(map[string]string)
	conf.file = path
	conf.lock = new(sync.RWMutex)

	err := conf.getResource(path)
	if err != nil {
		return nil, err
	} else {
		return conf, nil
	}
}

func (c *INIConfig) getResource(dir string) error {
	cont, err := ioutil.ReadFile(dir)

	if err != nil {
		return err
	}
	conf := strings.Split(string(cont), "\n")
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, val := range conf {
		val = strings.TrimSpace(val)
		val = strings.TrimRight(val, "\r")
		a := c.trimSpace(val)
		if len(a) > 0 && a[0] == '#' {
			continue
		}
		key, keyVal, err := c.splitEqual(a)
		if err == nil {
			c.val[key] = keyVal
		}
	}
	return nil
}

func (c *INIConfig) trimSpace(str string) string {
	var rst []byte

	by := []byte(str)
	bgn := 0
	for _, val := range by {
		if val == '"' && bgn == 0 {
			bgn = 1
		} else if val == '"' && bgn == 1 {
			bgn = 0
		}

		if bgn == 0 && val == ' ' {

		} else {
			rst = append(rst, val)
		}
	}

	return string(rst)
}

// Modify the value of key in the configuration file.
// If the key does not exist, the key value pair will be added to the configuration file
func (c *INIConfig) Set(key, value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.val[key]; ok {

		// modify configuration variable values.
		fd, err := os.OpenFile(c.file, os.O_RDWR, 666)
		if err != nil {
			return err
		}
		newValue := key + "=" + value

		rd := make([]byte, 1024)

		var lines []iniLineOne
		var tmp iniLineOne
		var line []byte
		var index = -1
		var offset = 0
		var lstIndex = 0
		for {
			size, err := fd.Read(rd)
			if err != nil && err != io.EOF {
				return err
			}

			if size == 0 {
				if len(line) > 0 {
					tmp.index = lstIndex
					tmp.line = string(line)
					tmp.offset = offset + 1
					lines = append(lines, tmp)
				}
				break
			}

			for _, val := range rd[0:size] {
				if val == '\x0a' {
					index++
					offset++
					tmp.index = lstIndex
					lstIndex = index
					tmp.offset = offset
					if len(line) == 0 {
						continue
					}
					if line[len(line)-1] == '\x0d' {
						tmp.line = string(line[:len(line)-1])
						tmp.offset -= 1
					} else {
						tmp.line = string(line)
					}
					lines = append(lines, tmp)
					line = make([]byte, 0)
					offset = 0
				} else {
					index++
					offset++
					line = append(line, val)
				}
			}
		}

		for _, val := range lines {
			t1 := c.trimSpace(val.line)
			mkey, _, _ := c.splitEqual(t1)

			if mkey == key {
				var aprst []byte
				fd.Seek(int64(val.index+val.offset), 0)
				for {
					tmp := make([]byte, 1024)
					n, err := fd.Read(tmp)
					if err != io.EOF && err != nil {
						return err
					}
					if n == 0 {
						break
					}
					aprst = append(aprst, tmp[0:n]...)
				}
				if val.index == 0 {
					fd.Seek(int64(val.index), 0)
				} else {
					fd.Seek(int64(val.index+1), 0)
				}

				if len(newValue) < val.offset {
					if val.offset-len(newValue) > 1 {
						var tb []byte = make([]byte, val.offset-len(newValue)-1)
						for i := 0; i < val.offset-len(newValue)-1; i++ {
							tb[i] = '\x20'
						}
						newValue = newValue + string(tb) + "\n"
					} else {
						newValue = newValue + "\n"
					}
				} else {
					newValue = newValue + "\n"
				}

				if len(aprst) > 0 && aprst[0] == '\x0a' {
					fd.WriteString(newValue + string(aprst[1:]))
				} else {
					fd.WriteString(newValue + string(aprst))
				}
				fd.Close()
			}
		}
		return nil
	} else {
		// add configuration variable values.
		op := key + "=" + value + "\n"
		fd, err := os.OpenFile(c.file, os.O_APPEND, 666)
		if err != nil {
			return err
		}
		defer fd.Close()
		// 读取最后文本最后一个自己，如果不是换行符，则添加一个换行符
		fd.Seek(-1, 2)
		var b []byte = make([]byte, 1)
		fd.Read(b)
		if b[0] != '\n' {
			op = "\n" + op
		}
		_, err = fd.WriteString(op)
		return err
	}
	c.getResource(c.file)
	return nil
}

// read key's value from the configuration file
// if the key does not exist, return key's value is dirty data, error is not nil.
// if the key exist, return the key's value, error is nil
func (c *INIConfig) Get(key string) (string, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if val, ok := c.val[key]; ok {
		return val, nil
	} else {
		return "", errors.New("The key [" + key + "] doesn't exist")
	}
}

func (c *INIConfig) splitEqual(str string) (string, string, error) {
	if len(str) == 0 {
		return "", "", errors.New("empty value")
	}
	bgn := 0
	end := 0
	key := ""
	keyVal := ""

	for _, val := range str {
		if val == '"' && bgn == 0 {
			bgn = 1
		} else if val == '"' && bgn == 1 {
			bgn = 0
		}

		if bgn == 0 && val == '=' && end == 0 {
			end = 1
		} else if end == 0 {
			key += string(val)
		} else {
			keyVal += string(val)
		}
	}
	if keyVal == "" || key == "" {
		return "", "", errors.New("empty value")
	}
	return strings.Trim(key, "\""), strings.Trim(keyVal, "\""), nil
}
