package evdev

/*
 #include <linux/input.h>
 #include <linux/uinput.h>
 // wrap c macro function
 static int _EVIOCGNAME(int len) {return EVIOCGNAME(len);}
 static int _EVIOCGPHYS(int len) {return EVIOCGPHYS(len);}
 static int _EVIOCGUNIQ(int len) {return EVIOCGUNIQ(len);}
 static int _EVIOCGPROP(int len) {return EVIOCGPROP(len);}

 static int _EVIOCGKEY(int len) {return EVIOCGKEY(len);}
 static int _EVIOCGLED(int len) {return EVIOCGLED(len);}
 static int _EVIOCGSND(int len) {return EVIOCGSND(len);}
 static int _EVIOCGSW(int len)  {return EVIOCGSW(len);}

 static int _EVIOCGBIT(int ev, int len) {return EVIOCGBIT(ev, len);}
 static int _EVIOCGABS(int abs)    {return EVIOCGABS(abs);}
 static int _EVIOCSABS(int abs)    {return EVIOCSABS(abs);}

 static int _UI_GET_SYSNAME(int len) { return UI_GET_SYSNAME(len);}
*/
import "C"
import (
	"syscall"
	"unsafe"
)

type _InputEvent C.struct_input_event
type _InputAbsinfo C.struct_input_absinfo
type _InputId C.struct_input_id
type _InputKeymapEntry C.struct_input_keymap_entry

const (
	sizeofInputAbsinfo     = C.sizeof_struct_input_absinfo
	sizeofInputId          = C.sizeof_struct_input_id
	sizeofInputKeymapEntry = C.sizeof_struct_input_keymap_entry
)

const MAX_NAME_SIZE = 256

const (
	EVIOCGID      = C.EVIOCGID      // get device ID
	EVIOCGVERSION = C.EVIOCGVERSION // get driver version
	EVIOCGREP     = C.EVIOCGREP     // get repeat settings
	EVIOCSREP     = C.EVIOCSREP     // set repeat settings

	EVIOCGKEYCODE    = C.EVIOCGKEYCODE    // get keycode
	EVIOCGKEYCODE_V2 = C.EVIOCGKEYCODE_V2 // get keycode

	EVIOCSKEYCODE    = C.EVIOCSKEYCODE    // set keycode
	EVIOCSKEYCODE_V2 = C.EVIOCSKEYCODE_V2 // set keycode

	EVIOCSFF      = C.EVIOCSFF      // send a force effect to a force feedback device
	EVIOCRMFF     = C.EVIOCRMFF     // erase a force effect
	EVIOCGEFFECTS = C.EVIOCGEFFECTS // report number of effects playable at the same time

	EVIOCGRAB     = C.EVIOCGRAB     // grab/release device
	EVIOCSCLOCKID = C.EVIOCSCLOCKID // set clockid to be used for timestamps
)

var EVIOCGNAME = C._EVIOCGNAME(MAX_NAME_SIZE) // get device name
var EVIOCGPHYS = C._EVIOCGPHYS(MAX_NAME_SIZE) // get physical location
var EVIOCGUNIQ = C._EVIOCGUNIQ(MAX_NAME_SIZE) // get unique identifier
var EVIOCGPROP = C._EVIOCGPROP(MAX_NAME_SIZE) // get device properties

var EVIOCGKEY = C._EVIOCGKEY(MAX_NAME_SIZE) // get global key state
var EVIOCGLED = C._EVIOCGLED(MAX_NAME_SIZE) // get all LEDs
var EVIOCGSND = C._EVIOCGSND(MAX_NAME_SIZE) // get all sounds status
var EVIOCGSW = C._EVIOCGSW(MAX_NAME_SIZE)   // get all switch states

func EVIOCGBIT(ev, l int) int { return int(C._EVIOCGBIT(C.int(ev), C.int(l))) } // get event bitsb
func EVIOCGABS(abs int) int   { return int(C._EVIOCGABS(C.int(abs))) }          // get abs bits
func EVIOCSABS(abs int) int   { return int(C._EVIOCSABS(C.int(abs))) }          // set abs bits

// uinput

const (
	UI_DEV_CREATE  = C.UI_DEV_CREATE  // create device
	UI_DEV_DESTROY = C.UI_DEV_DESTROY // destroy device

	UI_DEV_SETUP = C.UI_DEV_SETUP // setup device
	UI_ABS_SETUP = C.UI_ABS_SETUP //

	UI_SET_EVBIT   = C.UI_SET_EVBIT
	UI_SET_KEYBIT  = C.UI_SET_KEYBIT
	UI_SET_RELBIT  = C.UI_SET_RELBIT
	UI_SET_ABSBIT  = C.UI_SET_ABSBIT
	UI_SET_MSCBIT  = C.UI_SET_MSCBIT
	UI_SET_LEDBIT  = C.UI_SET_LEDBIT
	UI_SET_SNDBIT  = C.UI_SET_SNDBIT
	UI_SET_FFBIT   = C.UI_SET_FFBIT
	UI_SET_PHYS    = C.UI_SET_PHYS
	UI_SET_SWBIT   = C.UI_SET_SWBIT
	UI_SET_PROPBIT = C.UI_SET_PROPBIT
	UI_GET_VERSION = C.UI_GET_VERSION
)

var UI_GET_SYSNAME = C._UI_GET_SYSNAME(MAX_NAME_SIZE)

type _Uinput_setup = C.struct_uinput_setup

const uinput_name_len = 80

func ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
