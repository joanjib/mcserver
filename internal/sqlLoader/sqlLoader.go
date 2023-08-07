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
package sqlLoader

import (
	//"fmt"
	"github.com/joanjib/mcserver/internal/files"
	"io/ioutil"
	"log"
)

func LoadStrStatements(dirName string) map[string]string {
	strStatements := make(map[string]string)
	dir, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		fileName := file.Name()[0 : len(file.Name())-4]
		strStatements[fileName] = string(files.LoadFile(dirName + fileName + ".sql"))
	}
	/* printing for testing purposes
	   for k,s := range strStatements {
	       fmt.Println("----------------------------------------")
	       fmt.Println(k,s)
	   }
	*/
	return strStatements

}
