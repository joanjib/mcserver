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
package config

type Config struct {
	CommitInterval     int    `yaml:"commit-interval"`
	DatabaseFile       string `yaml:"database-file"`
	SqlStatementsDir   string `yaml:"sql-statements-dir"`
	ExecutionQueueSize int    `yaml:"execution-queue-size"`
	SchemaFile         string `yaml:"schema-file"`
}
