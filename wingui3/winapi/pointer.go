package winapi

import (
	"image"
	"strings"
	"time"
)

// Event is a pointer event.
type Event struct {
	SWin   *Window
	Kind   Kind
	Source Source
	// PointerID is the id for the pointer and can be used
	// to track a particular pointer from Press to
	// Release or Cancel.
	PointerID ID
	// Priority is the priority of the receiving handler
	// for this event.
	Priority Priority
	// Time is when the event was received. The
	// timestamp is relative to an undefined base.
	Time time.Duration
	// Buttons are the set of pressed mouse buttons for this event.
	Buttons Buttons
	// Position is the coordinates of the event in the local coordinate
	// system of the receiving tag. The transformation from global window
	// coordinates to local coordinates is performed by the inverse of
	// the effective transformation of the tag.
	Position image.Point
	// Scroll is the scroll amount, if any.
	Scroll image.Point
	// Modifiers is the set of active modifiers when
	// the mouse button was pressed.
	Modifiers Modifiers
	Name      string
	State     Kind
}

// PassOp sets the pass-through mode. InputOps added while the pass-through
// mode is set don't block events to siblings.
type PassOp struct {
}

type ID uint16

// Kind of an Event.
type Kind uint

// Priority of an Event.
type Priority uint8

// Source of an Event.
type Source uint8

// Buttons is a set of mouse buttons
type Buttons uint8

// Cursor denotes a pre-defined cursor shape. Its Add method adds an
// operation that sets the cursor shape for the current clip area.
type Cursor byte

// The cursors correspond to CSS pointer naming.
const (
	// CursorDefault is the default cursor.
	CursorDefault Cursor = iota
	// CursorNone hides the cursor. To show it again, use any other cursor.
	CursorNone
	// CursorText is for selecting and inserting text.
	CursorText
	// CursorVerticalText is for selecting and inserting vertical text.
	CursorVerticalText
	// CursorPointer is for a link.
	// Usually displayed as a pointing hand.
	CursorPointer
	// CursorCrosshair is for a precise location.
	CursorCrosshair
	// CursorAllScroll is for indicating scrolling in all directions.
	// Usually displayed as arrows to all four directions.
	CursorAllScroll
	// CursorColResize is for vertical resize.
	// Usually displayed as a vertical bar with arrows pointing east and west.
	CursorColResize
	// CursorRowResize is for horizontal resize.
	// Usually displayed as a horizontal bar with arrows pointing north and south.
	CursorRowResize
	// CursorGrab is for content that can be grabbed (dragged to be moved).
	// Usually displayed as an open hand.
	CursorGrab
	// CursorGrabbing is for content that is being grabbed (dragged to be moved).
	// Usually displayed as a closed hand.
	CursorGrabbing
	// CursorNotAllowed is shown when the request action cannot be carried out.
	// Usually displayed as a circle with a line through.
	CursorNotAllowed
	// CursorWait is shown when the program is busy and user cannot interact.
	// Usually displayed as a hourglass or the system equivalent.
	CursorWait
	// CursorProgress is shown when the program is busy, but the user can still interact.
	// Usually displayed as a default cursor with a hourglass.
	CursorProgress
	// CursorNorthWestResize is for top-left corner resizing.
	// Usually displayed as an arrow towards north-west.
	CursorNorthWestResize
	// CursorNorthEastResize is for top-right corner resizing.
	// Usually displayed as an arrow towards north-east.
	CursorNorthEastResize
	// CursorSouthWestResize is for bottom-left corner resizing.
	// Usually displayed as an arrow towards south-west.
	CursorSouthWestResize
	// CursorSouthEastResize is for bottom-right corner resizing.
	// Usually displayed as an arrow towards south-east.
	CursorSouthEastResize
	// CursorNorthSouth is for top-bottom resizing.
	// Usually displayed as a bi-directional arrow towards north-south.
	CursorNorthSouthResize
	// CursorEastWestResize is for left-right resizing.
	// Usually displayed as a bi-directional arrow towards east-west.
	CursorEastWestResize
	// CursorWestResize is for left resizing.
	// Usually displayed as an arrow towards west.
	CursorWestResize
	// CursorEastResize is for right resizing.
	// Usually displayed as an arrow towards east.
	CursorEastResize
	// CursorNorthResize is for top resizing.
	// Usually displayed as an arrow towards north.
	CursorNorthResize
	// CursorSouthResize is for bottom resizing.
	// Usually displayed as an arrow towards south.
	CursorSouthResize
	// CursorNorthEastSouthWestResize is for top-right to bottom-left diagonal resizing.
	// Usually displayed as a double ended arrow on the corresponding diagonal.
	CursorNorthEastSouthWestResize
	// CursorNorthWestSouthEastResize is for top-left to bottom-right diagonal resizing.
	// Usually displayed as a double ended arrow on the corresponding diagonal.
	CursorNorthWestSouthEastResize
)

const (
	// A Cancel event is generated when the current gesture is
	// interrupted by other handlers or the system.
	Cancel Kind = (1 << iota) >> 1
	// Press of a pointer.
	Press
	// Release of a pointer.
	Release
	// Move of a pointer.
	Move
	// Drag of a pointer.
	Drag
	// Pointer enters an area watching for pointer input
	Enter
	// Pointer leaves an area watching for pointer input
	Leave
	// Scroll of a pointer.
	Scroll
)

const (
	// Mouse generated event.
	Mouse Source = iota
	// Touch generated event.
	Touch
)

const (
	// Shared priority is for handlers that
	// are part of a matching set larger than 1.
	Shared Priority = iota
	// Foremost priority is like Shared, but the
	// handler is the foremost of the matching set.
	Foremost
	// Grabbed is used for matching sets of size 1.
	Grabbed
)

const (
	// ButtonPrimary is the primary button, usually the left button for a
	// right-handed user.
	ButtonPrimary Buttons = 1 << iota
	// ButtonSecondary is the secondary button, usually the right button for a
	// right-handed user.
	ButtonSecondary
	// ButtonTertiary is the tertiary button, usually the middle button.
	ButtonTertiary
)

func (t Kind) String() string {
	if t == Cancel {
		return "Cancel"
	}
	var buf strings.Builder
	for tt := Kind(1); tt > 0; tt <<= 1 {
		if t&tt > 0 {
			if buf.Len() > 0 {
				buf.WriteByte('|')
			}
			buf.WriteString((t & tt).string())
		}
	}
	return buf.String()
}

func (t Kind) string() string {
	switch t {
	case Press:
		return "Press"
	case Release:
		return "Release"
	case Cancel:
		return "Cancel"
	case Move:
		return "Move"
	case Drag:
		return "Drag"
	case Enter:
		return "Enter"
	case Leave:
		return "Leave"
	case Scroll:
		return "Scroll"
	default:
		panic("unknown Type")
	}
}

func (p Priority) String() string {
	switch p {
	case Shared:
		return "Shared"
	case Foremost:
		return "Foremost"
	case Grabbed:
		return "Grabbed"
	default:
		panic("unknown priority")
	}
}

func (s Source) String() string {
	switch s {
	case Mouse:
		return "Mouse"
	case Touch:
		return "Touch"
	default:
		panic("unknown source")
	}
}

// Contain reports whether the set b contains
// all of the buttons.
func (b Buttons) Contain(buttons Buttons) bool {
	return b&buttons == buttons
}

func (b Buttons) String() string {
	var strs []string
	if b.Contain(ButtonPrimary) {
		strs = append(strs, "ButtonPrimary")
	}
	if b.Contain(ButtonSecondary) {
		strs = append(strs, "ButtonSecondary")
	}
	if b.Contain(ButtonTertiary) {
		strs = append(strs, "ButtonTertiary")
	}
	return strings.Join(strs, "|")
}

func (c Cursor) String() string {
	switch c {
	case CursorDefault:
		return "Default"
	case CursorNone:
		return "None"
	case CursorText:
		return "Text"
	case CursorVerticalText:
		return "VerticalText"
	case CursorPointer:
		return "Pointer"
	case CursorCrosshair:
		return "Crosshair"
	case CursorAllScroll:
		return "AllScroll"
	case CursorColResize:
		return "ColResize"
	case CursorRowResize:
		return "RowResize"
	case CursorGrab:
		return "Grab"
	case CursorGrabbing:
		return "Grabbing"
	case CursorNotAllowed:
		return "NotAllowed"
	case CursorWait:
		return "Wait"
	case CursorProgress:
		return "Progress"
	case CursorNorthWestResize:
		return "NorthWestResize"
	case CursorNorthEastResize:
		return "NorthEastResize"
	case CursorSouthWestResize:
		return "SouthWestResize"
	case CursorSouthEastResize:
		return "SouthEastResize"
	case CursorNorthSouthResize:
		return "NorthSouthResize"
	case CursorEastWestResize:
		return "EastWestResize"
	case CursorWestResize:
		return "WestResize"
	case CursorEastResize:
		return "EastResize"
	case CursorNorthResize:
		return "NorthResize"
	case CursorSouthResize:
		return "SouthResize"
	case CursorNorthEastSouthWestResize:
		return "NorthEastSouthWestResize"
	case CursorNorthWestSouthEastResize:
		return "NorthWestSouthEastResize"
	default:
		panic("unknown Type")
	}
}

func (Event) ImplementsEvent() {}
