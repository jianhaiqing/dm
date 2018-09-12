// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package master

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/pingcap/tidb-enterprise-tools/dm/ctl/common"
	"github.com/pingcap/tidb-enterprise-tools/dm/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewUpdateTaskCmd creates a UpdateTask command
func NewUpdateTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-task [-w worker ...] <config_file>",
		Short: "update a task's config for routes, filters, column-mappings, black-white-list",
		Run:   updateTaskFunc,
	}
	return cmd
}

// updateTaskFunc does update task request
func updateTaskFunc(cmd *cobra.Command, _ []string) {
	if len(cmd.Flags().Args()) != 1 {
		fmt.Println(cmd.Usage())
		return
	}
	content, err := common.GetFileContent(cmd.Flags().Arg(0))
	if err != nil {
		common.PrintLines("get file content error:\n%v", errors.ErrorStack(err))
		return
	}

	workers, err := common.GetWorkerArgs(cmd)
	if err != nil {
		common.PrintLines("%s", errors.ErrorStack(err))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// NOTE: do whole check now, refine to do TablesChecker and ShardingTablesCheck ?
	err = checkTask(ctx, string(content))
	if err != nil {
		common.PrintLines("precheck failed %s", errors.ErrorStack(err))
		return
	}

	// update task
	cli := common.MasterClient()
	resp, err := cli.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Task:    string(content),
		Workers: workers,
	})
	if err != nil {
		common.PrintLines("can not update task:\n%v", errors.ErrorStack(err))
		return
	}

	common.PrettyPrintResponse(resp)
}