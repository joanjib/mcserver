// Mcserver
// Copyright (C) 2023  JUAN JOSÉ IGLESIAS BLANCH

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package files

import (
	"io/ioutil"
	"log"
	"os"
)

func LoadFile(file string) []byte {

	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err)
	}
	return []byte(fileContents)

}

func CreateFile(file string) *os.File {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	return f
}
func AppendToFile(f *os.File, text string) {

	if _, err := f.WriteString(text); err != nil {
		log.Panic(err)

	}

}
