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

package subtaskmeta_sorter

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/plugin"
)

type TableSorter struct {
	metas []*plugin.SubTaskMeta
}

func NewTableSorter(metas []*plugin.SubTaskMeta) SubTaskMetaSorter {
	return &TableSorter{metas: metas}
}

func (d *TableSorter) Sort() ([]plugin.SubTaskMeta, error) {
	return dbTableTopologicalSort(d.metas)
}

const (
	prefixCollect string = "collect"
	prefixExtract string = "extract"
	prefixConvert string = "convert"
)

func genClassNameByMetaName(rawName string) (string, error) {
	if len(rawName) > 7 {
		return rawName[7:], nil
	}
	return "", fmt.Errorf("got illeagal raw name = %s", rawName)
}

// stable topological sort
func dbTableTopologicalSort(metas []*plugin.SubTaskMeta) ([]plugin.SubTaskMeta, error) {
	// TODO 1. can i use reflect to realize collect, extractor, converter ?
	// first process same class data
	// suppose different class has no dependency relation
	// construct class name list and table list meta
	// sort different metas
	// add list by convert and

	// 1. construct data
	subtaskMap := make(map[string][]*plugin.SubTaskMeta)       // use subtask class name to get metalist
	subtaskTableMap := make(map[string][]string)               // use class name get meta name list
	subtaskMetaNameMap := make(map[string]*plugin.SubTaskMeta) // use name to get meta
	for _, metaItem := range metas {
		taskClassName, err := genClassNameByMetaName(metaItem.Name)
		if err != nil {
			return nil, err
		}
		if value, ok := subtaskMap[taskClassName]; ok {
			subtaskMap[taskClassName] = append(value, metaItem)
		} else {
			subtaskMap[taskClassName] = []*plugin.SubTaskMeta{metaItem}
		}
		if value, ok := subtaskTableMap[taskClassName]; ok {
			// check if subtask in one class has different tables define
			if len(value) != len(metaItem.Tables) {
				return nil, fmt.Errorf("got different table list in class %s", taskClassName)
			}
			// TODO check list item in value and metaItem.Tables, make sure it's equal
		} else {
			subtaskTableMap[taskClassName] = metaItem.Tables
		}
		subtaskMetaNameMap[metaItem.Name] = metaItem
	}

	// 2. sort
	sortedNameList, err := topologicalSortDifferentElements(subtaskTableMap)
	if err != nil {
		return nil, err
	}

	// 3. gen subtaskmeta list by sorted data and return
	sortedSubtaskMetaList := make([]plugin.SubTaskMeta, 0)
	for _, nameItem := range sortedNameList {
		value, ok := subtaskMap[nameItem]
		if !ok {
			return nil, fmt.Errorf("failed get subtask list by class name = %s", nameItem)
		}
		tmpList := make([]plugin.SubTaskMeta, len(value))
		for _, subtaskItem := range value {
			if len(value) >= 2 && len(subtaskItem.Name) > 7 {
				switch subtaskItem.Name[:7] {
				case prefixCollect:
					tmpList[0] = *subtaskItem
				case prefixExtract:
					tmpList[1] = *subtaskItem
				case prefixConvert:
					if len(value) == 3 {
						tmpList[2] = *subtaskItem
					} else {
						return nil, fmt.Errorf("got wrong length of list with extrac subtask")
					}
				default:
					return nil, fmt.Errorf("got wrong length of subtask %v", subtaskItem)
				}
			}
		}
		sortedSubtaskMetaList = append(sortedSubtaskMetaList, tmpList...)
	}
	return sortedSubtaskMetaList, nil
}

// TODO get subtask class list, different class can task concurrency
func GetSortedClassName() []string {
	return nil
}

// TODO get subtask list by class name, this subtask list should run sequentially
func GetSubtaskMetasByClassName(className string) []*plugin.SubTaskMeta {
	return nil
}
