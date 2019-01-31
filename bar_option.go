package mpb

import (
	"github.com/vbauerster/mpb/v4/decor"
)

// BarOption is a function option which changes the default behavior of a bar.
type BarOption func(*bState)

// AppendDecorators let you inject decorators to the bar's right side.
func AppendDecorators(appenders ...decor.Decorator) BarOption {
	return func(s *bState) {
		for _, decorator := range appenders {
			if ar, ok := decorator.(decor.AmountReceiver); ok {
				s.amountReceivers = append(s.amountReceivers, ar)
			}
			if sl, ok := decorator.(decor.ShutdownListener); ok {
				s.shutdownListeners = append(s.shutdownListeners, sl)
			}
			s.aDecorators = append(s.aDecorators, decorator)
		}
	}
}

// PrependDecorators let you inject decorators to the bar's left side.
func PrependDecorators(prependers ...decor.Decorator) BarOption {
	return func(s *bState) {
		for _, decorator := range prependers {
			if ar, ok := decorator.(decor.AmountReceiver); ok {
				s.amountReceivers = append(s.amountReceivers, ar)
			}
			if sl, ok := decorator.(decor.ShutdownListener); ok {
				s.shutdownListeners = append(s.shutdownListeners, sl)
			}
			s.pDecorators = append(s.pDecorators, decorator)
		}
	}
}

// BarID sets bar id.
func BarID(id int) BarOption {
	return func(s *bState) {
		s.id = id
	}
}

// BarWidth sets bar width independent of the container.
func BarWidth(width int) BarOption {
	return func(s *bState) {
		s.width = width
	}
}

// BarRemoveOnComplete is a flag, if set whole bar line will be removed
// on complete event. If both BarRemoveOnComplete and BarClearOnComplete
// are set, first bar section gets cleared and then whole bar line
// gets removed completely.
func BarRemoveOnComplete() BarOption {
	return func(s *bState) {
		s.removeOnComplete = true
	}
}

// BarReplaceOnComplete is indicator for delayed bar start, after the
// `runningBar` is complete. To achieve bar replacement effect,
// `runningBar` should has its `BarRemoveOnComplete` option set.
func BarReplaceOnComplete(runningBar *Bar) BarOption {
	return func(s *bState) {
		s.runningBar = runningBar
	}
}

// BarClearOnComplete is a flag, if set will clear bar section on
// complete event. If you need to remove a whole bar line, refer to
// BarRemoveOnComplete.
func BarClearOnComplete() BarOption {
	return func(s *bState) {
		s.barClearOnComplete = true
	}
}

// BarPriority sets bar's priority. Zero is highest priority, i.e. bar
// will be on top. If `BarReplaceOnComplete` option is supplied, this
// option is ignored.
func BarPriority(priority int) BarOption {
	return func(s *bState) {
		s.priority = priority
	}
}

// BarExtender is an option to extend bar to the next new line, with
// arbitrary output.
func BarExtender(extender Filler) BarOption {
	return func(s *bState) {
		s.extender = extender
	}
}

// TrimSpace trims bar's edge spaces.
func TrimSpace() BarOption {
	return func(s *bState) {
		s.trimSpace = true
	}
}

// BarStyle sets custom bar style.
// Effective when Filler type is bar.
func BarStyle(style string) BarOption {
	chk := func(filler Filler) (interface{}, bool) {
		if style == "" {
			return nil, false
		}
		t, ok := filler.(*barFiller)
		return t, ok
	}
	cb := func(t interface{}) {
		t.(*barFiller).setStyle(style)
	}
	return MakeFillerTypeSpecificBarOption(chk, cb)
}

// SpinnerStyle sets custom spinner style.
// Effective when Filler type is spinner.
func SpinnerStyle(frames []string) BarOption {
	chk := func(filler Filler) (interface{}, bool) {
		if len(frames) == 0 {
			return nil, false
		}
		t, ok := filler.(*spinnerFiller)
		return t, ok
	}
	cb := func(t interface{}) {
		t.(*spinnerFiller).frames = frames
	}
	return MakeFillerTypeSpecificBarOption(chk, cb)
}

// MakeFillerTypeSpecificBarOption makes BarOption specific to Filler's
// actual type. If you implement your own Filler, so most probably
// you'll need this. See BarStyle or SpinnerStyle for example.
func MakeFillerTypeSpecificBarOption(
	typeChecker func(Filler) (interface{}, bool),
	cb func(interface{}),
) BarOption {
	return func(s *bState) {
		if t, ok := typeChecker(s.filler); ok {
			cb(t)
		}
	}
}

// BarOptOnCond returns option when condition evaluates to true.
func BarOptOnCond(option BarOption, condition func() bool) BarOption {
	if condition() {
		return option
	}
	return nil
}
