package fod

import (
	tea "charm.land/bubbletea/v2"
)

// Dialog execute dialog
func Dialog(opt Option) ([]string, ResultCode) {

	// create selector
	var selector Selector
	var err error
	if selector, err = newSelector(opt.Mode, opt.Multi); err != nil {
		return []string{}, ResultNone
	}

	model := newDialogModel(selector, opt.Base)
	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		return []string{}, ResultNone
	}

	dialog, ok := finalModel.(dialogModel)
	if !ok {
		return []string{}, ResultNone
	}
	return dialog.selector.result()
}
