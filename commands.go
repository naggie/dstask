package dstask

// CommandNext ...
func CommandNext(repoPath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(repoPath, WithStatuses(NON_RESOLVED_STATUSES...))
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterOutStatus(STATUS_TEMPLATE)
	ts.SortByPriority()
	ctx.PrintContextDescription()
	ts.DisplayByNext(true)
	ts.DisplayCriticalTaskWarning()

	return nil
}

// CommandShowOpen ...
func CommandShowOpen(repoPath string, ctx, cmdLine CmdLine) error {
	ts, err := NewTaskSet(repoPath, WithStatuses(NON_RESOLVED_STATUSES...))
	if err != nil {
		return err
	}
	ts.Filter(ctx)
	ts.Filter(cmdLine)
	ts.FilterOutStatus(STATUS_TEMPLATE)
	ts.SortByPriority()
	ctx.PrintContextDescription()
	ts.DisplayByNext(false)
	ts.DisplayCriticalTaskWarning()
	return nil
}
