/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package subtaskmetaSorter

import (
	"errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type DependencySorter struct {
	metas []*plugin.SubTaskMeta
}

func GetDependencySorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &DependencySorter{metas: metas}
}

func (d *DependencySorter) Sort() ([]plugin.SubTaskMeta, error) {
	// create a map to store the dependencies of each subtaskmeta
	dependencyMap := make(map[string][]*plugin.SubTaskMeta)
	for _, meta := range d.metas {
		dependencyMap[meta.Name] = meta.Dependencies
	}

	// create a map to store the visited status of each subtaskmeta
	visitedMap := make(map[string]bool)

	// create a slice to store the sorted subtaskmetas
	var sortedMetas []plugin.SubTaskMeta

	// visit each subtaskmeta
	for _, meta := range d.metas {
		err := visit(meta.Name, dependencyMap, visitedMap, &sortedMetas)
		if err != nil {
			return nil, err
		}
	}

	return sortedMetas, nil
}

func visit(name string, dependencyMap map[string][]*plugin.SubTaskMeta, visitedMap map[string]bool, sortedMetas *[]plugin.SubTaskMeta) error {
	// if the subtaskmeta has already been visited, return nil
	if visitedMap[name] {
		return nil
	}

	// mark the subtaskmeta as visited
	visitedMap[name] = true

	// visit each dependency of the subtaskmeta
	for _, dependency := range dependencyMap[name] {
		err := visit(dependency.Name, dependencyMap, visitedMap, sortedMetas)
		if err != nil {
			return err
		}
	}

	// add the subtaskmeta to the sorted slice
	*sortedMetas = append(*sortedMetas, plugin.SubTaskMeta{Name: name})

	return nil
}

func (d *DependencySorter) DetectLoop() error {
	// create a map to store the dependencies of each subtaskmeta
	dependencyMap := make(map[string][]*plugin.SubTaskMeta)
	for _, meta := range d.metas {
		dependencyMap[meta.Name] = meta.Dependencies
	}

	// create a map to store the visited status of each subtaskmeta
	visitedMap := make(map[string]bool)

	// visit each subtaskmeta
	for _, meta := range d.metas {
		err := detectLoop(meta.Name, dependencyMap, visitedMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func detectLoop(name string, dependencyMap map[string][]*plugin.SubTaskMeta, visitedMap map[string]bool) error {
	// if the subtaskmeta has already been visited, return an error
	if visitedMap[name] {
		return errors.New("loop detected")
	}

	// mark the subtaskmeta as visited
	visitedMap[name] = true

	// visit each dependency of the subtaskmeta
	for _, dependency := range dependencyMap[name] {
		err := detectLoop(dependency.Name, dependencyMap, visitedMap)
		if err != nil {
			return err
		}
	}

	// mark the subtaskmeta as unvisited
	visitedMap[name] = false

	return nil
}
