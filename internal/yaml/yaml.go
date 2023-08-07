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
package yaml

import (
	"github.com/joanjib/mcserver/internal/files"
	"gopkg.in/yaml.v3"
	"log"
)

// LoadYml load a yml file into a struct passed at data
func LoadYml(file string, data interface{}) {
	fileBytes := files.LoadFile(file)

	if err := yaml.Unmarshal(fileBytes, data); err != nil {
		log.Panicf("Error reading yml file %s with error: %v", file, err)
	}

}
