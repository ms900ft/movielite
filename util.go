package movielight

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

var defaultRegexes = map[string]string{
	"111.OTR": `^(.+?)_\d{2}\.\d{2}\.\d{2}_\d{2}-\d{2}_`,
	"1.OTR":   `^(.+?)_Top_Tipp_+\d{2}\.\d{2}\.\d{2}_\d{2}-\d{2}_`,
	"22.MT":   `^.+?-(.+?)-\d+-\d+.+\.\w{3}$`,
	"21.MT":   `^.+?_-_.+?-(.+)-\d+-\d+.+\.\w{3}$`,
	//	"MZ2":           `^_.+?-(.+)`,
	"1.ARTE_FERNSEHFILM": `.+?-_Fernsehfilme-(.+?)-\d+-\d+.+\.\w{3}$`,
	"3.MT_ZDF":           `^.+?-(.+?)_-`,
	"4.MT_ZDF_SEASON":    `^.+?-(.+?)_\((\d+)\)-\d+_\w+_\d+`,
	"5.MT_ARD":           `^.+?-(.+?)-\d+_\w+_\d+`,
	"11111.MT_ARD":       `^Filme_im_Ersten-(.+?)-\d+`,
	"6.MT_RBB":           `^.+?rbb-(.+?)-[0-9A-Fa-f\-]+_\d+`,
	"7.3_SAT":            `^(.+)_-_.+`,
	"998.Simple":         `^(.+)\s-\s.+\.\w\w\w$`,
	"999.Simple":         `^(.+)\.\w\w\w$`,
}

func Translatename(filename string, regexes map[string]string) string {
	if len(regexes) > 0 {
		log.Debug("merging regexes\n")
	}
	ret := ""
	keys := make([]string, 0)
	for k := range defaultRegexes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		//fmt.Println(k, defaultRegexes[k])
		name := k
		//for name, regex := range defaultRegexes {

		re := regexp.MustCompile(defaultRegexes[k])
		match := re.FindStringSubmatch(filename)
		if len(match) != 0 {
			i := len(match)
			match2 := match[1:i]
			ret = strings.Join(match2, " ")
			log.Debugf("Found %s Name: %s\n", name, ret)
			break
		}
	}
	ret = strings.Replace(ret, "_", " ", -1)
	ret = strings.TrimSpace(ret)
	return ret
}

func URLEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

const raw = `
on run argv
  tell application "Finder"
    repeat with f in argv
      move (f as POSIX file) to trash
    end repeat
  end tell
end run
`

//var trash = filepath.Join("/Volumes/Transcend/.Trashes", "501")

func Trash(f string, trash string) (trashcan string, err error) {
	bin, err := exec.LookPath("osascript")
	if err != nil {
		err = fmt.Errorf("not yet supported")
		return
	}

	if _, err = os.Stat(trash); err != nil {
		err = fmt.Errorf("trash not found")
		return
	}

	path, err := filepath.Abs(f)
	if err != nil {
		return
	}
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	_ = name

	dest := filepath.Join(trash, base)
	if _, err = os.Stat(dest); err == nil {
		err = fmt.Errorf("already exists")
		return
	}
	trashcan = dest
	log.Debug(path)
	params := append([]string{"-e", raw}, path)
	cmd := exec.Command(bin, params...)
	log.Debugf("%+v", params)
	if err = cmd.Run(); err != nil {
		log.Error(err)
		return
	}

	if _, err = os.Stat(trashcan); err != nil {
		trashcan = ""
	}

	return trashcan, err
}

// func toJSON(value []byte) (JSONB, error) {
// 	jsonRes, _ := json.MarshalIndent(value, "", "  ")
// 	var mj JSONB
// 	err := json.Unmarshal(jsonRes, &mj)
// 	return mj, err
//}
